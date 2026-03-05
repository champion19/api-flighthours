package services

import (
	"context"
	"strings"

	"github.com/Nerzal/gocloak/v13"
	"github.com/champion19/api-flighthours/core/interactor/dto"
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/input"
	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/jwt"
	"github.com/champion19/api-flighthours/platform/logger"
)

type service struct {
	repository output.Repository
	keycloak   output.AuthClient
}

func NewService(repository output.Repository, keycloak output.AuthClient) input.Service {
	return &service{
		repository: repository,
		keycloak:   keycloak,
	}
}
func (s service) GetEmployeeByEmail(ctx context.Context, email string) (*domain.Employee, error) {
	log.Debug(logger.LogEmployeeServiceSearchByEmail, "email", email)
	employee, err := s.repository.GetEmployeeByEmail(ctx, email)
	if err != nil {
		log.Error(logger.LogEmployeeServiceErrorByEmail, "email", email, "error", err)
		return nil, err
	}
	log.Debug(logger.LogEmployeeServiceFoundByEmail, "email", email, "employee_id", employee.ID)
	return employee, nil
}

func (s service) GetEmployeeByID(ctx context.Context, id string) (*domain.Employee, error) {
	log.Debug(logger.LogEmployeeServiceSearchByID, "employee_id", id)
	employee, err := s.repository.GetEmployeeByID(ctx, id)
	if err != nil {
		log.Error(logger.LogEmployeeServiceErrorByID, "employee_id", id, "error", err)
		return nil, err
	}
	log.Debug(logger.LogEmployeeServiceFoundByID, "employee_id", id, "email", employee.Email)
	return employee, nil
}

func (s service) GetEmployeeByKeycloakID(ctx context.Context, keycloakUserID string) (*domain.Employee, error) {
	log.Debug(logger.LogEmployeeServiceSearchByID, "keycloak_user_id", keycloakUserID)
	employee, err := s.repository.GetEmployeeByKeycloakID(ctx, keycloakUserID)
	if err != nil {
		log.Error(logger.LogEmployeeServiceErrorByID, "keycloak_user_id", keycloakUserID, "error", err)
		return nil, err
	}
	log.Debug(logger.LogEmployeeServiceFoundByID, "keycloak_user_id", keycloakUserID, "email", employee.Email)
	return employee, nil
}

func (s service) BeginTx(ctx context.Context) (output.Tx, error) {
	return s.repository.BeginTx(ctx)
}

func (s service) RegisterEmployee(ctx context.Context, employee domain.Employee) (*dto.RegisterEmployee, error) {
	log.Info(logger.LogEmployeeServiceValidationStart, employee.ToLogger())
	log.Debug(logger.LogDualSystemCheck, "email", employee.Email)

	existingEmployee, errDB := s.repository.GetEmployeeByEmail(ctx, employee.Email)
	if err := checkSystemAvailability(errDB, employee.Email, logger.LogDatabaseUnavailable); err != nil {
		return nil, err
	}

	dbExists := errDB == nil && existingEmployee != nil

	keycloakUser, errKC := s.keycloak.GetUserByEmail(ctx, employee.Email)
	if err := checkSystemAvailability(errKC, employee.Email, logger.LogKeycloakUnavailable); err != nil {
		return nil, err
	}

	kcExists := errKC == nil && keycloakUser != nil

	if dbExists && kcExists {
		log.Warn(logger.LogUserExistsInBoth, "email", employee.Email)
		return nil, domain.ErrDuplicateUser
	}

	if dbExists && !kcExists {
		log.Warn(logger.LogUserExistsOnlyInDB,
			"email", employee.Email,
			"employee_id", existingEmployee.ID,
			"action", "will be cleaned")
		return nil, domain.ErrIncompleteRegistration
	}

	if !dbExists && kcExists {
		log.Warn(logger.LogUserExistsOnlyInKeycloak,
			"email", employee.Email,
			"keycloak_id", *keycloakUser.ID,
			"action", "will be cleaned")
		return nil, domain.ErrIncompleteRegistration
	}

	log.Debug(logger.LogUserNotFoundInEither, "email", employee.Email)
	log.Info(logger.LogEmployeeServiceValidationComplete, employee.ToLogger())
	return &dto.RegisterEmployee{
		Employee: employee,
		Message:  "Validaciones exitosas",
	}, nil
}

// checkSystemAvailability returns a domain error if the given error is a connection or timeout error.
func checkSystemAvailability(err error, email, logMsg string) error {
	if err == nil {
		return nil
	}
	if isConnectionError(err) || isTimeoutError(err) {
		log.Error(logMsg, "email", email, "error", err, "error_type", "connection")
		if logMsg == logger.LogDatabaseUnavailable {
			return domain.ErrDatabaseUnavailable
		}
		return domain.ErrKeycloakUnavailable
	}
	return nil
}

func (s service) SaveEmployeeToDB(ctx context.Context, tx output.Tx, employee domain.Employee) error {
	log.Info(logger.LogEmployeeServiceSavingToDB, employee.ToLogger())
	err := s.repository.Save(ctx, tx, employee)
	if err != nil {
		log.Error(logger.LogEmployeeServiceSaveError, employee.ToLogger(), "error", err)
		return err
	}
	log.Success(logger.LogEmployeeServiceSavedToDB, employee.ToLogger())
	return nil
}

func (s service) UpdateEmployee(ctx context.Context, tx output.Tx, employee domain.Employee) error {
	log.Info(logger.LogEmployeeUpdating, employee.ToLogger())
	err := s.repository.UpdateEmployee(ctx, tx, employee)
	if err != nil {
		log.Error(logger.LogEmployeeUpdateError, employee.ToLogger(), "error", err)
		return err
	}
	log.Success(logger.LogEmployeeUpdated, employee.ToLogger())
	return nil
}

func (s service) CreateUserInKeycloak(ctx context.Context, employee *domain.Employee) (string, error) {
	log.Info(logger.LogEmployeeServiceCreatingKeycloak, employee.ToLogger())
	keycloakUserID, err := s.keycloak.CreateUser(ctx, employee)
	if err != nil {
		if isConnectionError(err) || isTimeoutError(err) {
			log.Error(logger.LogKeycloakUnavailable,
				employee.ToLogger(),
				"error", err,
				"error_type", "connection")
			return "", domain.ErrKeycloakUnavailable
		}
		log.Error(logger.LogEmployeeServiceKeycloakError, employee.ToLogger(), "error", err)
		return "", domain.ErrKeycloakUserCreationFailed
	}
	log.Success(logger.LogEmployeeServiceCreatedKeycloak, employee.ToLogger(), "Keycloak_user_id", keycloakUserID)
	return keycloakUserID, nil
}

func (s service) SetUserPassword(ctx context.Context, userID string, password string) error {
	log.Debug(logger.LogEmployeeServicePasswordSet, "keycloak_user_id", userID)
	err := s.keycloak.SetPassword(ctx, userID, password, false)
	if err != nil {
		log.Error(logger.LogEmployeeServicePasswordError, "keycloak_user_id", userID, "error", err)
		return err
	}

	log.Success(logger.LogEmployeeServicePasswordSetOK, "keycloak_user_id", userID)
	return nil
}

func (s service) AssignUserRole(ctx context.Context, userID string, role string) error {
	log.Info(logger.LogEmployeeServiceRoleAssigning, "keycloak_user_id", userID, "role", role)
	err := s.keycloak.AssignRole(ctx, userID, role)
	if err != nil {
		log.Error(logger.LogEmployeeServiceRoleError, "keycloak_user_id", userID, "role", role, "error", err)
		return err
	}
	log.Success(logger.LogEmployeeServiceRoleAssigned, "keycloak_user_id", userID, "role", role)
	return nil
}

func (s service) UpdateEmployeeKeycloakID(ctx context.Context, tx output.Tx, employeeID string, keycloakUserID string) error {
	log.Debug(logger.LogEmployeeServiceKeycloakIDUpdate, "employee_id", employeeID, "keycloak_user_id", keycloakUserID)
	err := s.repository.PatchEmployee(ctx, tx, employeeID, keycloakUserID)
	if err != nil {
		log.Error(logger.LogEmployeeServiceKeycloakIDUpdateError, "employee_id", employeeID, "error", err)
		return err
	}
	log.Success(logger.LogEmployeeServiceKeycloakIDUpdated, "employee_id", employeeID, "keycloak_user_id", keycloakUserID)
	return nil
}

func (s service) RollbackEmployee(ctx context.Context, employeeID string) error {
	log.Warn(logger.LogEmployeeServiceRollbackEmployee, "employee_id", employeeID)
	err := s.repository.DeleteEmployee(ctx, nil, employeeID)
	if err != nil {
		log.Error(logger.LogEmployeeServiceRollbackEmployeeError, "employee_id", employeeID, "error", err)
		return err
	}
	log.Info(logger.LogEmployeeServiceRollbackEmployeeComplete, "employee_id", employeeID)
	return nil
}

func (s service) RollbackKeycloakUser(ctx context.Context, KeycloakUserID string) error {
	log.Warn(logger.LogEmployeeServiceRollbackKeycloak, "keycloak_user_id", KeycloakUserID)
	err := s.keycloak.DeleteUser(ctx, KeycloakUserID)
	if err != nil {
		log.Error(logger.LogEmployeeServiceRollbackKeycloakError, "keycloak_user_id", KeycloakUserID, "error", err)
		return err
	}
	log.Info(logger.LogEmployeeServiceRollbackKeycloakComplete, "keycloak_user_id", KeycloakUserID)
	return nil
}

func (s service) CheckAndCleanInconsistentState(ctx context.Context, email string) error {
	log.Debug(logger.LogDualSystemCheck, "email", email)

	employeeInDB, errDB := s.repository.GetEmployeeByEmail(ctx, email)
	dbExists := errDB == nil && employeeInDB != nil

	keycloakUser, errKC := s.keycloak.GetUserByEmail(ctx, email)
	kcExists := errKC == nil && keycloakUser != nil

	if dbExists == kcExists {
		if dbExists {
			log.Debug(logger.LogUserExistsInBoth, "email", email)
		} else {
			log.Debug(logger.LogUserNotFoundInEither, "email", email)
		}
		return nil
	}

	log.Warn(logger.LogInconsistentStateDetect,
		"email", email,
		"in_database", dbExists,
		"in_keycloak", kcExists,
		"db_person_id", getIDOrNA(dbExists, func() string { return employeeInDB.ID }),
		"kc_user_id", getIDOrNA(kcExists, func() string { return *keycloakUser.ID }))

	if !dbExists && kcExists {
		return s.cleanOrphanKeycloak(ctx, email, *keycloakUser.ID)
	}

	return s.cleanOrphanDB(ctx, email, employeeInDB.ID)
}

func getIDOrNA(exists bool, idFunc func() string) string {
	if exists {
		return idFunc()
	}
	return "N/A"
}

func (s service) cleanOrphanKeycloak(ctx context.Context, email, keycloakUserID string) error {
	log.Info(logger.LogEmployeeServiceCleaningOrphan,
		"email", email,
		"source", "keycloak",
		"keycloak_user_id", keycloakUserID,
		"reason", "missing in business database")

	if err := s.keycloak.DeleteUser(ctx, keycloakUserID); err != nil {
		log.Error(logger.LogEmployeeServiceOrphanCleanError,
			"email", email,
			"source", "keycloak",
			"keycloak_user_id", keycloakUserID,
			"error", err)
		return domain.ErrKeycloakCleanupFailed
	}

	log.Success(logger.LogEmployeeServiceOrphanCleaned,
		"email", email,
		"source", "keycloak",
		"action", "deleted from Keycloak")
	return nil
}

func (s service) cleanOrphanDB(ctx context.Context, email, employeeID string) error {
	log.Info(logger.LogEmployeeServiceCleaningOrphan,
		"email", email,
		"source", "database",
		"employee_id", employeeID,
		"reason", "missing in Keycloak")

	if err := s.repository.DeleteEmployee(ctx, nil, employeeID); err != nil {
		log.Error(logger.LogEmployeeServiceOrphanCleanError,
			"email", email,
			"source", "database",
			"employee_id", employeeID,
			"error", err)
		return domain.ErrKeycloakCleanupFailed
	}

	log.Success(logger.LogEmployeeServiceOrphanCleaned,
		"email", email,
		"source", "database",
		"action", "deleted from business database")
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
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

func (s service) GetUserByEmail(ctx context.Context, email string) (*gocloak.User, error) {
	log.Debug(logger.LogKeycloakSearchUserByEmail, "email", email)
	user, err := s.keycloak.GetUserByEmail(ctx, email)
	if err != nil {
		log.Error(logger.LogKeycloakUserNotFound, "email", email, "error", err)
		return nil, err
	}
	log.Debug(logger.LogKeycloakSearchUserByEmailOK, "email", email, "user_id", *user.ID)
	return user, nil
}

func (s service) SendVerificationEmail(ctx context.Context, userID string) error {
	log.Debug(logger.LogKeycloakSendVerificationEmail, "user_id", userID)

	err := s.keycloak.SendVerificationEmail(ctx, userID)
	if err != nil {
		log.Error(logger.LogKeycloakSendVerificationEmailError, "user_id", userID, "error", err)
		return err
	}
	log.Success(logger.LogKeycloakSendVerificationEmailOK, "user_id", userID)
	return nil
}

func (s service) SendPasswordResetEmail(ctx context.Context, email string) error {
	log.Debug(logger.LogKeycloakSendPasswordReset, "email", email)

	err := s.keycloak.SendPasswordResetEmail(ctx, email)
	if err != nil {
		log.Error(logger.LogKeycloakSendPasswordResetError, "email", email, "error", err)
		return err
	}

	log.Success(logger.LogKeycloakSendPasswordResetOK, "email", email)
	return nil
}

func (s service) Login(ctx context.Context, email, password string) (*gocloak.JWT, error) {
	log.Debug(logger.LogKeycloakLoginCheckingVerification, "email", email)

	user, err := s.keycloak.GetUserByEmail(ctx, email)
	if err != nil {
		log.Error(logger.LogKeycloakUserNotFound, "email", email, "error", err)
		return nil, domain.ErrUserNotFound
	}

	if user.EmailVerified == nil || !*user.EmailVerified {
		log.Warn(logger.LogKeycloakLoginEmailNotVerified, "email", email, "user_id", *user.ID)

		log.Info(logger.LogKeycloakLoginResendingVerification, "email", email, "user_id", *user.ID)
		if sendErr := s.keycloak.SendVerificationEmail(ctx, *user.ID); sendErr != nil {
			log.Error(logger.LogKeycloakLoginResendVerificationError,
				"email", email,
				"user_id", *user.ID,
				"error", sendErr)
		} else {
			log.Success(logger.LogKeycloakLoginResendVerificationOK, "email", email, "user_id", *user.ID)
		}

		return nil, domain.ErrorEmailNotVerified
	}

	log.Debug(logger.LogKeycloakLoginEmailVerified, "email", email, "user_id", *user.ID)

	log.Debug(logger.LogKeycloakUserLogin, "email", email)
	token, err := s.keycloak.LoginUser(ctx, email, password)
	if err != nil {
		log.Error(logger.LogKeycloakUserLoginError, "email", email, "error", err)
		return nil, err
	}

	log.Success(logger.LogKeycloakUserLoginOK, "email", email)
	return token, nil
}

func (s service) RefreshToken(ctx context.Context, refreshToken string) (*gocloak.JWT, error) {
	log.Info(logger.LogKeycloakUserTokenRefresh)

	if refreshToken == "" {
		log.Warn(logger.LogKeycloakRefreshTokenEmpty)
		return nil, domain.ErrRefreshTokenFailed
	}

	token, err := s.keycloak.RefreshToken(ctx, refreshToken)
	if err != nil {
		log.Error(logger.LogKeycloakUserTokenRefreshErr, "error", err)
		return nil, domain.ErrRefreshTokenFailed
	}

	log.Success(logger.LogKeycloakUserTokenRefreshOK)
	return token, nil
}

func (s service) VerifyEmailByToken(ctx context.Context, token string) (string, error) {
	log.Info(logger.LogKeycloakEmailVerify)

	tokenParser := jwt.NewTokenParser()
	email, err := tokenParser.ExtractEmailFromToken(token)
	if err != nil {
		log.Error(logger.LogKeycloakEmailVerifyError, "error", err, "reason", "failed to extract email from token")
		return "", domain.ErrInvalidToken
	}

	log.Debug(logger.LogEmployeeServiceEmailExtracted, "email", email)

	user, err := s.keycloak.GetUserByEmail(ctx, email)
	if err != nil {
		log.Error(logger.LogKeycloakUserNotFound, "email", email, "error", err)
		return "", domain.ErrUserNotFound
	}

	if user.EmailVerified != nil && *user.EmailVerified {
		log.Warn(logger.LogKeycloakEmailAlreadyVerified, "email", email, "user_id", *user.ID)
		return email, domain.ErrEmailAlreadyVerified
	}
	if err := s.keycloak.VerifyEmail(ctx, *user.ID); err != nil {
		log.Error(logger.LogKeycloakEmailVerifyError, "email", email, "user_id", *user.ID, "error", err)
		return "", err
	}

	log.Success(logger.LogKeycloakEmailVerifyOK, "email", email, "user_id", *user.ID)
	return email, nil
}

func (s service) LocateEmployee(ctx context.Context, id string) (*dto.RegisterEmployee, error) {
	employee, err := s.repository.GetEmployeeByID(ctx, id)
	if err != nil {
		log.Error(logger.LogEmployeeGetByIDError, "error", err)
		return nil, err
	}

	if employee == nil {
		log.Warn(logger.LogEmployeeNotFound, "id", id)
		return nil, domain.ErrPersonNotFound
	}

	return &dto.RegisterEmployee{
		Employee: *employee,
		Message:  "Employee located successfully",
	}, nil
}

func (s service) UpdatePassword(ctx context.Context, token, newPassword string) (string, error) {
	log.Info(logger.LogKeycloakPasswordUpdate)
	log.Debug(logger.LogKeycloakPasswordTokenValidation)
	userID, email, err := s.keycloak.ValidateActionToken(ctx, token)
	if err != nil {
		log.Error(logger.LogKeycloakPasswordTokenInvalid, "error", err)
		return "", domain.ErrInvalidToken
	}
	log.Debug(logger.LogKeycloakPasswordTokenValidOK, "user_id", userID, "email", email)

	if err := s.keycloak.SetPassword(ctx, userID, newPassword, false); err != nil {
		log.Error(logger.LogKeycloakPasswordUpdateError, "user_id", userID, "error", err)
		return "", domain.ErrPasswordUpdateFailed
	}

	log.Success(logger.LogKeycloakPasswordUpdateOK, "user_id", userID, "email", email)
	return email, nil
}

func (s service) ChangePassword(ctx context.Context, email, currentPassword, newPassword string) (string, error) {
	log.Info(logger.LogKeycloakChangePassword, "email", email)
	log.Debug(logger.LogKeycloakChangePasswordValidating, "email", email)
	user, err := s.keycloak.GetUserByEmail(ctx, email)
	if err != nil {
		log.Error(logger.LogKeycloakUserNotFound, "email", email, "error", err)
		return "", domain.ErrUserNotFound
	}

	_, err = s.keycloak.LoginUser(ctx, email, currentPassword)
	if err != nil {
		log.Warn(logger.LogKeycloakChangePasswordInvalid, "email", email, "error", err)
		return "", domain.ErrInvalidCurrentPassword
	}

	if err := s.keycloak.SetPassword(ctx, *user.ID, newPassword, false); err != nil {
		log.Error(logger.LogKeycloakChangePasswordError, "email", email, "user_id", *user.ID, "error", err)
		return "", domain.ErrPasswordUpdateFailed
	}

	log.Success(logger.LogKeycloakChangePasswordOK, "email", email, "user_id", *user.ID)
	return email, nil
}

func (s service) DeleteEmployee(ctx context.Context, employeeID string, keycloakUserID string) error {
	log.Info(logger.LogEmployeeDeleting, "employee_id", employeeID, "keycloak_user_id", keycloakUserID)
	if keycloakUserID != "" {
		log.Debug(logger.LogEmployeeDeletingKeycloak, "keycloak_user_id", keycloakUserID)
		if err := s.keycloak.DeleteUser(ctx, keycloakUserID); err != nil {
			log.Error(logger.LogEmployeeDeleteKeycloakError, "keycloak_user_id", keycloakUserID, "error", err)
			return domain.ErrUserCannotDelete
		}
		log.Success(logger.LogEmployeeDeletedKeycloak, "keycloak_user_id", keycloakUserID)
	}

	log.Debug(logger.LogEmployeeDeletingDB, "employee_id", employeeID)
	if err := s.repository.DeleteEmployee(ctx, nil, employeeID); err != nil {
		log.Error(logger.LogEmployeeDeleteDBError, "employee_id", employeeID, "error", err)
		return domain.ErrUserCannotDelete
	}

	log.Success(logger.LogEmployeeDeleteComplete, "employee_id", employeeID)
	return nil
}
