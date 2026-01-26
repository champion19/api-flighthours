package services

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
	"github.com/champion19/flighthours-api/core/interactor/dto"
	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/input"
	"github.com/champion19/flighthours-api/core/ports/output"
	"github.com/champion19/flighthours-api/platform/jwt"
	"github.com/champion19/flighthours-api/platform/logger"
)

type service struct {
	repository output.Repository
	keycloak   output.AuthClient
	logger     logger.Logger
}

func NewService(repository output.Repository, keycloak output.AuthClient, logger logger.Logger) input.Service {
	return &service{
		repository: repository,
		keycloak:   keycloak,
		logger:     logger,
	}
}
func (s service) GetEmployeeByEmail(ctx context.Context, email string) (*domain.Employee, error) {
	s.logger.Debug(logger.LogEmployeeServiceSearchByEmail, "email", email)
	employee, err := s.repository.GetEmployeeByEmail(ctx, email)
	if err != nil {
		s.logger.Error(logger.LogEmployeeServiceErrorByEmail, "email", email, "error", err)
		return nil, err
	}
	s.logger.Debug(logger.LogEmployeeServiceFoundByEmail, "email", email, "employee_id", employee.ID)
	return employee, nil
}

func (s service) GetEmployeeByID(ctx context.Context, id string) (*domain.Employee, error) {
	s.logger.Debug(logger.LogEmployeeServiceSearchByID, "employee_id", id)
	employee, err := s.repository.GetEmployeeByID(ctx, id)
	if err != nil {
		s.logger.Error(logger.LogEmployeeServiceErrorByID, "employee_id", id, "error", err)
		return nil, err
	}
	s.logger.Debug(logger.LogEmployeeServiceFoundByID, "employee_id", id, "email", employee.Email)
	return employee, nil
}

func (s service) GetEmployeeByKeycloakID(ctx context.Context, keycloakUserID string) (*domain.Employee, error) {
	s.logger.Debug(logger.LogEmployeeServiceSearchByID, "keycloak_user_id", keycloakUserID)
	employee, err := s.repository.GetEmployeeByKeycloakID(ctx, keycloakUserID)
	if err != nil {
		s.logger.Error(logger.LogEmployeeServiceErrorByID, "keycloak_user_id", keycloakUserID, "error", err)
		return nil, err
	}
	s.logger.Debug(logger.LogEmployeeServiceFoundByID, "keycloak_user_id", keycloakUserID, "email", employee.Email)
	return employee, nil
}

func (s service) BeginTx(ctx context.Context) (output.Tx, error) {
	return s.repository.BeginTx(ctx)
}

func (s service) RegisterEmployee(ctx context.Context, employee domain.Employee) (*dto.RegisterEmployee, error) {
	s.logger.Info(logger.LogEmployeeServiceValidationStart, employee.ToLogger())
	s.logger.Debug(logger.LogDualSystemCheck, "email", employee.Email)

	existingEmployee, errDB := s.repository.GetEmployeeByEmail(ctx, employee.Email)
	if errDB != nil {
		if isConnectionError(errDB) || isTimeoutError(errDB) {
			//TODO: Agregar mensaje de log aquí
			s.logger.Error(logger.LogDatabaseUnavailable,
				"email", employee.Email,
				"error", errDB,
				"error_type", "connection")
			return nil, domain.ErrDatabaseUnavailable
		}

	}

	dbExists := errDB == nil && existingEmployee != nil

	keycloakUser, errKC := s.keycloak.GetUserByEmail(ctx, employee.Email)

	if errKC != nil {
		if isConnectionError(errKC) || isTimeoutError(errKC) {
			s.logger.Error(logger.LogKeycloakUnavailable,
				"email", employee.Email,
				"error", errKC,
				"error_type", "connection")
			return nil, domain.ErrKeycloakUnavailable
		}
	}

	kcExists := errKC == nil && keycloakUser != nil

	if dbExists && kcExists {
		s.logger.Warn(logger.LogUserExistsInBoth, "email", employee.Email)
		return nil, domain.ErrDuplicateUser
	}

	if dbExists && !kcExists {
		s.logger.Warn(logger.LogUserExistsOnlyInDB,
			"email", employee.Email,
			"employee_id", existingEmployee.ID,
			"action", "will be cleaned")
		return nil, domain.ErrIncompleteRegistration
	}

	if !dbExists && kcExists {
		s.logger.Warn(logger.LogUserExistsOnlyInKeycloak,
			"email", employee.Email,
			"keycloak_id", *keycloakUser.ID,
			"action", "will be cleaned")
		return nil, domain.ErrIncompleteRegistration
	}

	s.logger.Debug(logger.LogUserNotFoundInEither, "email", employee.Email)
	s.logger.Info(logger.LogEmployeeServiceValidationComplete, employee.ToLogger())
	return &dto.RegisterEmployee{
		Employee: employee,
		Message:  "Validaciones exitosas",
	}, nil
}

func (s service) SaveEmployeeToDB(ctx context.Context, tx output.Tx, employee domain.Employee) error {
	s.logger.Info(logger.LogEmployeeServiceSavingToDB, employee.ToLogger())
	err := s.repository.Save(ctx, tx, employee)
	if err != nil {
		s.logger.Error(logger.LogEmployeeServiceSaveError, employee.ToLogger(), "error", err)
		return err
	}
	s.logger.Success(logger.LogEmployeeServiceSavedToDB, employee.ToLogger())
	return nil
}

func (s service) CreateUserInKeycloak(ctx context.Context, employee *domain.Employee) (string, error) {
	s.logger.Info(logger.LogEmployeeServiceCreatingKeycloak, employee.ToLogger())
	keycloakUserID, err := s.keycloak.CreateUser(ctx, employee)
	if err != nil {
		if isConnectionError(err) || isTimeoutError(err) {
			s.logger.Error(logger.LogKeycloakUnavailable,
				employee.ToLogger(),
				"error", err,
				"error_type", "connection")
			return "", domain.ErrKeycloakUnavailable
		}
		s.logger.Error(logger.LogEmployeeServiceKeycloakError, employee.ToLogger(), "error", err)
		return "", domain.ErrKeycloakUserCreationFailed
	}
	s.logger.Success(logger.LogEmployeeServiceCreatedKeycloak, employee.ToLogger(), "Keycloak_user_id", keycloakUserID)
	return keycloakUserID, nil
}

func (s service) SetUserPassword(ctx context.Context, userID string, password string) error {
	s.logger.Debug(logger.LogEmployeeServicePasswordSet, "keycloak_user_id", userID)
	err := s.keycloak.SetPassword(ctx, userID, password, false)
	if err != nil {
		s.logger.Error(logger.LogEmployeeServicePasswordError, "keycloak_user_id", userID, "error", err)
		return err
	}

	s.logger.Success(logger.LogEmployeeServicePasswordSetOK, "keycloak_user_id", userID)
	return nil
}

func (s service) AssignUserRole(ctx context.Context, userID string, role string) error {
	s.logger.Info(logger.LogEmployeeServiceRoleAssigning, "keycloak_user_id", userID, "role", role)
	err := s.keycloak.AssignRole(ctx, userID, role)
	if err != nil {
		s.logger.Error(logger.LogEmployeeServiceRoleError, "keycloak_user_id", userID, "role", role, "error", err)
		return err
	}
	s.logger.Success(logger.LogEmployeeServiceRoleAssigned, "keycloak_user_id", userID, "role", role)
	return nil
}

func (s service) UpdateEmployeeKeycloakID(ctx context.Context, tx output.Tx, employeeID string, keycloakUserID string) error {
	s.logger.Debug(logger.LogEmployeeServiceKeycloakIDUpdate, "employee_id", employeeID, "keycloak_user_id", keycloakUserID)
	err := s.repository.PatchEmployee(ctx, tx, employeeID, keycloakUserID)
	if err != nil {
		s.logger.Error(logger.LogEmployeeServiceKeycloakIDUpdateError, "employee_id", employeeID, "error", err)
		return err
	}
	s.logger.Success(logger.LogEmployeeServiceKeycloakIDUpdated, "employee_id", employeeID, "keycloak_user_id", keycloakUserID)
	return nil
}

func (s service) RollbackEmployee(ctx context.Context, employeeID string) error {
	s.logger.Warn(logger.LogEmployeeServiceRollbackEmployee, "employee_id", employeeID)
	err := s.repository.DeleteEmployee(ctx, nil, employeeID)
	if err != nil {
		s.logger.Error(logger.LogEmployeeServiceRollbackEmployeeError, "employee_id", employeeID, "error", err)
		return err
	}
	s.logger.Info(logger.LogEmployeeServiceRollbackEmployeeComplete, "employee_id", employeeID)
	return nil
}

func (s service) RollbackKeycloakUser(ctx context.Context, KeycloakUserID string) error {
	s.logger.Warn(logger.LogEmployeeServiceRollbackKeycloak, "keycloak_user_id", KeycloakUserID)
	err := s.keycloak.DeleteUser(ctx, KeycloakUserID)
	if err != nil {
		s.logger.Error(logger.LogEmployeeServiceRollbackKeycloakError, "keycloak_user_id", KeycloakUserID, "error", err)
		return err
	}
	s.logger.Info(logger.LogEmployeeServiceRollbackKeycloakComplete, "keycloak_user_id", KeycloakUserID)
	return nil
}

func (s service) CheckAndCleanInconsistentState(ctx context.Context, email string) error {
	s.logger.Debug(logger.LogDualSystemCheck, "email", email)

	employeeInDB, errDB := s.repository.GetEmployeeByEmail(ctx, email)
	dbExists := errDB == nil && employeeInDB != nil

	keycloakUser, errKC := s.keycloak.GetUserByEmail(ctx, email)
	kcExists := errKC == nil && keycloakUser != nil

	if (dbExists && kcExists) || (!dbExists && !kcExists) {
		if dbExists && kcExists {
			s.logger.Debug(logger.LogUserExistsInBoth, "email", email)
		} else {
			s.logger.Debug(logger.LogUserNotFoundInEither, "email", email)
		}
		return nil
	}

	s.logger.Warn(logger.LogInconsistentStateDetect,
		"email", email,
		"in_database", dbExists,
		"in_keycloak", kcExists,
		"db_person_id", func() string {
			if dbExists {
				return employeeInDB.ID
			}
			return "N/A"
		}(),
		"kc_user_id", func() string {
			if kcExists {
				return *keycloakUser.ID
			}
			return "N/A"
		}())

	if !dbExists && kcExists {
		s.logger.Info(logger.LogEmployeeServiceCleaningOrphan,
			"email", email,
			"source", "keycloak",
			"keycloak_user_id", *keycloakUser.ID,
			"reason", "missing in business database")

		if err := s.keycloak.DeleteUser(ctx, *keycloakUser.ID); err != nil {
			s.logger.Error(logger.LogEmployeeServiceOrphanCleanError,
				"email", email,
				"source", "keycloak",
				"keycloak_user_id", *keycloakUser.ID,
				"error", err)
			return domain.ErrKeycloakCleanupFailed
		}

		s.logger.Success(logger.LogEmployeeServiceOrphanCleaned,
			"email", email,
			"source", "keycloak",
			"action", "deleted from Keycloak")
		return nil
	}

	if dbExists && !kcExists {
		s.logger.Info(logger.LogEmployeeServiceCleaningOrphan,
			"email", email,
			"source", "database",
			"employee_id", employeeInDB.ID,
			"reason", "missing in Keycloak")

		if err := s.repository.DeleteEmployee(ctx, nil, employeeInDB.ID); err != nil {
			s.logger.Error(logger.LogEmployeeServiceOrphanCleanError,
				"email", email,
				"source", "database",
				"employee_id", employeeInDB.ID,
				"error", err)
			return domain.ErrKeycloakCleanupFailed
		}

		s.logger.Success(logger.LogEmployeeServiceOrphanCleaned,
			"email", email,
			"source", "database",
			"action", "deleted from business database")
		return nil
	}

	return nil
}
func isConnectionError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return contains(errStr, "connection refused") ||
		contains(errStr, "no such host") ||
		contains(errStr, "connection reset") ||
		contains(errStr, "network is unreachable") ||
		contains(errStr, "connect: connection refused")
}

func isTimeoutError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return contains(errStr, "timeout") ||
		contains(errStr, "deadline exceeded") ||
		contains(errStr, "context deadline exceeded")
}

func contains(s, substr string) bool {
	for i := 0; i+len(substr) <= len(s); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			c1 := s[i+j]
			c2 := substr[j]
			if c1 >= 'A' && c1 <= 'Z' {
				c1 += 32
			}
			if c2 >= 'A' && c2 <= 'Z' {
				c2 += 32
			}
			if c1 != c2 {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

func (s service) GetUserByEmail(ctx context.Context, email string) (*gocloak.User, error) {
	s.logger.Debug(logger.LogKeycloakSearchUserByEmail, "email", email)
	user, err := s.keycloak.GetUserByEmail(ctx, email)
	if err != nil {
		s.logger.Error(logger.LogKeycloakUserNotFound, "email", email, "error", err)
		return nil, err
	}
	s.logger.Debug(logger.LogKeycloakSearchUserByEmailOK, "email", email, "user_id", *user.ID)
	return user, nil
}

func (s service) SendVerificationEmail(ctx context.Context, userID string) error {
	s.logger.Debug(logger.LogKeycloakSendVerificationEmail, "user_id", userID)

	err := s.keycloak.SendVerificationEmail(ctx, userID)
	if err != nil {
		s.logger.Error(logger.LogKeycloakSendVerificationEmailError, "user_id", userID, "error", err)
		return err
	}
	s.logger.Success(logger.LogKeycloakSendVerificationEmailOK, "user_id", userID)
	return nil
}

func (s service) SendPasswordResetEmail(ctx context.Context, email string) error {
	s.logger.Debug(logger.LogKeycloakSendPasswordReset, "email", email)

	err := s.keycloak.SendPasswordResetEmail(ctx, email)
	if err != nil {
		s.logger.Error(logger.LogKeycloakSendPasswordResetError, "email", email, "error", err)
		return err
	}

	s.logger.Success(logger.LogKeycloakSendPasswordResetOK, "email", email)
	return nil
}

func (s service) Login(ctx context.Context, email, password string) (*gocloak.JWT, error) {
	s.logger.Debug(logger.LogKeycloakLoginCheckingVerification, "email", email)

	user, err := s.keycloak.GetUserByEmail(ctx, email)
	if err != nil {
		s.logger.Error(logger.LogKeycloakUserNotFound, "email", email, "error", err)
		return nil, domain.ErrUserNotFound
	}

	if user.EmailVerified == nil || !*user.EmailVerified {
		s.logger.Warn(logger.LogKeycloakLoginEmailNotVerified, "email", email, "user_id", *user.ID)

		s.logger.Info(logger.LogKeycloakLoginResendingVerification, "email", email, "user_id", *user.ID)
		if sendErr := s.keycloak.SendVerificationEmail(ctx, *user.ID); sendErr != nil {
			s.logger.Error(logger.LogKeycloakLoginResendVerificationError,
				"email", email,
				"user_id", *user.ID,
				"error", sendErr)
		} else {
			s.logger.Success(logger.LogKeycloakLoginResendVerificationOK, "email", email, "user_id", *user.ID)
		}

		return nil, domain.ErrorEmailNotVerified
	}

	s.logger.Debug(logger.LogKeycloakLoginEmailVerified, "email", email, "user_id", *user.ID)

	s.logger.Debug(logger.LogKeycloakUserLogin, "email", email)
	token, err := s.keycloak.LoginUser(ctx, email, password)
	if err != nil {
		s.logger.Error(logger.LogKeycloakUserLoginError, "email", email, "error", err)
		return nil, err
	}

	s.logger.Success(logger.LogKeycloakUserLoginOK, "email", email)
	return token, nil
}

func (s service) VerifyEmailByToken(ctx context.Context, token string) (string, error) {
	s.logger.Info(logger.LogKeycloakEmailVerify)

	tokenParser := jwt.NewTokenParser()
	email, err := tokenParser.ExtractEmailFromToken(token)
	if err != nil {
		s.logger.Error(logger.LogKeycloakEmailVerifyError, "error", err, "reason", "failed to extract email from token")
		return "", domain.ErrInvalidToken
	}

	s.logger.Debug("Email extracted from token", "email", email)

	user, err := s.keycloak.GetUserByEmail(ctx, email)
	if err != nil {
		s.logger.Error(logger.LogKeycloakUserNotFound, "email", email, "error", err)
		return "", domain.ErrUserNotFound
	}

	if user.EmailVerified != nil && *user.EmailVerified {
		s.logger.Warn(logger.LogKeycloakEmailAlreadyVerified, "email", email, "user_id", *user.ID)
		return email, domain.ErrEmailAlreadyVerified
	}
	if err := s.keycloak.VerifyEmail(ctx, *user.ID); err != nil {
		s.logger.Error(logger.LogKeycloakEmailVerifyError, "email", email, "user_id", *user.ID, "error", err)
		return "", err
	}

	s.logger.Success(logger.LogKeycloakEmailVerifyOK, "email", email, "user_id", *user.ID)
	return email, nil
}

func (s service) LocateEmployee(ctx context.Context, id string) (*dto.RegisterEmployee, error) {
	employee, err := s.repository.GetEmployeeByID(ctx, id)
	if err != nil {
		s.logger.Error(logger.LogEmployeeGetByIDError, "error", err)
		return nil, err
	}

	if employee == nil {
		s.logger.Warn(logger.LogEmployeeNotFound, "id", id)
		return nil, domain.ErrPersonNotFound
	}

	return &dto.RegisterEmployee{
		Employee: *employee,
		Message:  "Employee located successfully",
	}, nil
}

func (s service) UpdatePassword(ctx context.Context, token, newPassword string) (string, error) {
	s.logger.Info(logger.LogKeycloakPasswordUpdate)
	s.logger.Debug(logger.LogKeycloakPasswordTokenValidation)
	userID, email, err := s.keycloak.ValidateActionToken(ctx, token)
	if err != nil {
		s.logger.Error(logger.LogKeycloakPasswordTokenInvalid, "error", err)
		return "", domain.ErrInvalidToken
	}
	s.logger.Debug(logger.LogKeycloakPasswordTokenValidOK, "user_id", userID, "email", email)

	if err := s.keycloak.SetPassword(ctx, userID, newPassword, false); err != nil {
		s.logger.Error(logger.LogKeycloakPasswordUpdateError, "user_id", userID, "error", err)
		return "", domain.ErrPasswordUpdateFailed
	}

	s.logger.Success(logger.LogKeycloakPasswordUpdateOK, "user_id", userID, "email", email)
	return email, nil
}

func (s service) ChangePassword(ctx context.Context, email, currentPassword, newPassword string) (string, error) {
	s.logger.Info(logger.LogKeycloakChangePassword, "email", email)
	s.logger.Debug(logger.LogKeycloakChangePasswordValidating, "email", email)
	user, err := s.keycloak.GetUserByEmail(ctx, email)
	if err != nil {
		s.logger.Error(logger.LogKeycloakUserNotFound, "email", email, "error", err)
		return "", domain.ErrUserNotFound
	}

	_, err = s.keycloak.LoginUser(ctx, email, currentPassword)
	if err != nil {
		s.logger.Warn(logger.LogKeycloakChangePasswordInvalid, "email", email, "error", err)
		return "", domain.ErrInvalidCurrentPassword
	}

	if err := s.keycloak.SetPassword(ctx, *user.ID, newPassword, false); err != nil {
		s.logger.Error(logger.LogKeycloakChangePasswordError, "email", email, "user_id", *user.ID, "error", err)
		return "", domain.ErrPasswordUpdateFailed
	}

	s.logger.Success(logger.LogKeycloakChangePasswordOK, "email", email, "user_id", *user.ID)
	return email, nil
}

func (s service) DeleteEmployee(ctx context.Context, employeeID string, keycloakUserID string) error {
	s.logger.Info(logger.LogEmployeeDeleting, "employee_id", employeeID, "keycloak_user_id", keycloakUserID)
	if keycloakUserID != "" {
		s.logger.Debug(logger.LogEmployeeDeletingKeycloak, "keycloak_user_id", keycloakUserID)
		if err := s.keycloak.DeleteUser(ctx, keycloakUserID); err != nil {
			s.logger.Error(logger.LogEmployeeDeleteKeycloakError, "keycloak_user_id", keycloakUserID, "error", err)
			return domain.ErrUserCannotDelete
		}
		s.logger.Success(logger.LogEmployeeDeletedKeycloak, "keycloak_user_id", keycloakUserID)
	}

	s.logger.Debug(logger.LogEmployeeDeletingDB, "employee_id", employeeID)
	if err := s.repository.DeleteEmployee(ctx, nil, employeeID); err != nil {
		s.logger.Error(logger.LogEmployeeDeleteDBError, "employee_id", employeeID, "error", err)
		return domain.ErrUserCannotDelete
	}

	s.logger.Success(logger.LogEmployeeDeleteComplete, "employee_id", employeeID)
	return nil
}

func (s service) GetEmployeesByRole(ctx context.Context, role string) ([]domain.Employee, error) {
	s.logger.Debug(logger.LogCrewMemberTypeGet, "role", role)
	employees, err := s.repository.GetEmployeesByRole(ctx, role)
	if err != nil {
		s.logger.Error(logger.LogCrewMemberTypeGetError, "role", role, "error", err)
		return nil, err
	}
	s.logger.Debug(logger.LogCrewMemberTypeGetOK, "role", role, "count", len(employees))
	return employees, nil
}
