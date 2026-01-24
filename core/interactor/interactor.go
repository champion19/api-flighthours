package interactor

import (
	"context"

	"github.com/champion19/flighthours-api/core/interactor/dto"
	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/input"
	"github.com/champion19/flighthours-api/middleware"
	"github.com/champion19/flighthours-api/platform/logger"
)

type Interactor struct {
	service input.Service
	logger  logger.Logger
}

func NewInteractor(service input.Service, log logger.Logger) *Interactor {
	return &Interactor{
		service: service,
		logger:  log,
	}
}
func (i *Interactor) RegisterEmployee(ctx context.Context, employee domain.Employee) (result *dto.RegisterEmployee, err error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

	log.Info(logger.LogEmployeeInteractorRegStart, employee.ToLogger())
	//paso 1
	result, err = i.service.RegisterEmployee(ctx, employee)
	if err != nil {
		if err == domain.ErrIncompleteRegistration {
			log.Warn(logger.LogEmployeeInteractorIncompleteDetected, "email", employee.Email)

			//intentar limpiar el estado inconsistente antes de reintentar
			if cleanErr := i.service.CheckAndCleanInconsistentState(ctx, employee.Email); cleanErr != nil {
				log.Error(logger.LogEmployeeInteractorCleanup_Error, "email", employee.Email, "error", cleanErr)
				return nil, cleanErr
			}

			log.Success(logger.LogEmployeeInteractorCleanup_OK, "email", employee.Email)
			//retornar el error de registro incompleto para que el cliente sepa que debe reintentar
			return nil, err
		}
		log.Error(logger.LogEmployeeInteractorStep1_Error, "error", err)
		return
	}
	log.Success(logger.LogEmployeeInteractorStep1_OK)

	employee.SetID()
	log.Debug(logger.LogEmployeeInteractorIDGenerated, "employee_id", employee.ID)

	//paso 1.5
	if err = i.service.CheckAndCleanInconsistentState(ctx, employee.Email); err != nil {
		log.Error(logger.LogEmployeeInteractorStep15_Error, "error", err)
		return
	}
	log.Success(logger.LogEmployeeInteractorStep15_OK, "email", employee.Email)

	//paso 2
	tx, err := i.service.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogEmployeeInteractorStep2_Error, "error", err)
		return
	}
	log.Success(logger.LogEmployeeInteractorStep2_OK)

	var keycloakUserID string
	var keycloakCreated bool

	defer func() {
		if err != nil {

			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(logger.LogEmployeeInteractorRollbackDB_Error,
					"rollback_error", rbErr,
					"original_error", err)
			} else {
				log.Warn(logger.LogEmployeeInteractorRollbackDB_OK)
			}

			if keycloakCreated {
				if kcErr := i.service.RollbackKeycloakUser(ctx, keycloakUserID); kcErr != nil {
					log.Error(logger.LogEmployeeInteractorRollbackKeycloak_Err,
						"keycloak_error", kcErr,
						"keycloak_user_id", keycloakUserID)
				} else {
					log.Warn(logger.LogEmployeeInteractorRollbackKeycloak_OK)
				}
			}
		}
	}()
	//paso 3
	if err = i.service.SaveEmployeeToDB(ctx, tx, employee); err != nil {
		log.Error(logger.LogEmployeeInteractorStep3_Error, "error", err)
		return
	}
	log.Success(logger.LogEmployeeInteractorStep3_OK)

	//paso 4
	keycloakUserID, err = i.service.CreateUserInKeycloak(ctx, &employee)
	if err != nil {
		log.Error(logger.LogEmployeeInteractorStep4_Error, "error", err)
		err = domain.ErrKeycloakUserCreationFailed
		return
	}
	keycloakCreated = true
	log.Success(logger.LogEmployeeInteractorStep4_OK, "keycloak_user_id", keycloakUserID)

	//paso 5
	if err = i.service.SetUserPassword(ctx, keycloakUserID, employee.Password); err != nil {
		log.Error(logger.LogEmployeeInteractorStep5_Error, "error", err)
		return
	}
	log.Success(logger.LogEmployeeInteractorStep5_OK)

	//paso 6
	if err = i.service.AssignUserRole(ctx, keycloakUserID, employee.Role); err != nil {
		log.Error(logger.LogEmployeeInteractorStep6_Error, "error", err)
		return
	}
	log.Success(logger.LogEmployeeInteractorStep6_OK, "role", employee.Role)

	//paso 7
	if err = i.service.UpdateEmployeeKeycloakID(ctx, tx, employee.ID, keycloakUserID); err != nil {
		log.Error(logger.LogEmployeeInteractorStep7_Error, "error", err)
		return
	}
	log.Success(logger.LogEmployeeInteractorStep7_OK)

	//paso 8
	if err = tx.Commit(); err != nil {
		log.Error(logger.LogEmployeeInteractorCommit_Error, "error", err)
		return
	}
	log.Success(logger.LogEmployeeInteractorCommit_OK)

	employee.KeycloakUserID = keycloakUserID
	result.Employee = employee
	result.Message = "user registered successfully"

	// Paso 9: Enviar email de verificación
	if sendErr := i.service.SendVerificationEmail(ctx, keycloakUserID); sendErr != nil {
		//log warning pero no falla la creación del usuario
		log.Warn(logger.LogKeycloakSendVerificationEmailError,
			"keycloak_user_id", keycloakUserID,
			"email", employee.Email,
			"error", sendErr)
	} else {
		log.Info(logger.LogKeycloakSendVerificationEmailOK, "keycloak_user_id", keycloakUserID, "email", employee.Email)
	}

	log.Success(logger.LogEmployeeInteractorRegComplete, employee.ToLogger(),
		"Keycloak_id", keycloakUserID)
	err = nil

	return
}

func (i *Interactor) Locate(ctx context.Context, id string) (*dto.RegisterEmployee, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

	result, err := i.service.LocateEmployee(ctx, id)
	if err != nil {
		log.Error(logger.LogEmployeeLocateError, "error", err)
		return nil, err
	}
	return result, nil
}

// ResendVerificationEmail reenvía el email de verificación a un usuario por email
func (i *Interactor) ResendVerificationEmail(ctx context.Context, email string) error {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

	log.Info(logger.LogKeycloakSendVerificationEmail, "email", email)

	// Buscar empleado por email en la base de datos
	user, err := i.service.GetUserByEmail(ctx, email)
	if err != nil {
		log.Error(logger.LogKeycloakUserNotFound, "email", email, "error", err)
		return domain.ErrUserNotFound
	}

	//verificar si el email ya fue verificado
	if user.EmailVerified != nil && *user.EmailVerified {
		log.Warn(logger.LogKeycloakSendVerificationEmailError, "email", email, "reason", "email already verified")
		return domain.ErrEmailAlreadyVerified
	}

	// Enviar email de verificación
	if err = i.service.SendVerificationEmail(ctx, *user.ID); err != nil {
		log.Error(logger.LogKeycloakSendVerificationEmailError, "email", email, "error", err)
		return err
	}

	log.Success(logger.LogKeycloakSendVerificationEmailOK, "email", email, "user_id", *user.ID)
	return nil
}

func (i *Interactor) RequestPasswordReset(ctx context.Context, email string) error {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

	log.Info(logger.LogKeycloakSendPasswordReset, "email", email)

	// Llamar al servicio que busca el usuario y envía el email
	if err := i.service.SendPasswordResetEmail(ctx, email); err != nil {
		log.Warn(logger.LogKeycloakSendPasswordResetError, "email", email, "error", err)
	} else {
		log.Success(logger.LogKeycloakSendPasswordResetOK, "email", email)
	}

	return nil
}

// VerifyEmailByToken verifica el email de un usuario extrayéndolo del token JWT
// Este método delega al Service que maneja la lógica de negocio (parsing del token y verificación en Keycloak)
func (i *Interactor) VerifyEmailByToken(ctx context.Context, token string) (string, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

	log.Info(logger.LogKeycloakEmailVerify)

	// Delegar toda la lógica al Service (parsing del token + verificación en Keycloak)
	email, err := i.service.VerifyEmailByToken(ctx, token)
	if err != nil {
		switch err {
		case domain.ErrInvalidToken:
			log.Error(logger.LogKeycloakEmailVerifyError, "error", err, "reason", "invalid token")
		case domain.ErrUserNotFound:
			log.Warn(logger.LogKeycloakUserNotFound, "email", email)
		case domain.ErrEmailAlreadyVerified:
			log.Warn(logger.LogKeycloakEmailAlreadyVerified, "email", email)
		default:
			log.Error(logger.LogKeycloakEmailVerifyError, "email", email, "error", err)
		}
		return email, err
	}

	log.Success(logger.LogKeycloakEmailVerifyOK, "email", email)
	return email, nil
}

func (i *Interactor) Login(ctx context.Context, email string, password string) (*dto.TokenResponse, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

	log.Info(logger.LogKeycloakUserLogin, "email", email, "client_ip")

	// Llamar al servicio de autenticación de Keycloak
	token, err := i.service.Login(ctx, email, password)
	if err != nil {
		log.Error(logger.LogKeycloakUserLoginError, "email", email, "error", err, "client_ip")
		return nil, err
	}

	log.Success(logger.LogKeycloakUserLoginOK, "email", email)
	return &dto.TokenResponse{
		ExpiresIn:    token.ExpiresIn,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    token.TokenType,
	}, nil
}

// UpdatePassword validates the action token and updates the user's password
// This method handles password updates from the password reset flow
func (i *Interactor) UpdatePassword(ctx context.Context, token, newPassword, confirmPassword string) (string, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

	log.Info(logger.LogKeycloakPasswordUpdate)

	// Validate passwords match
	if newPassword != confirmPassword {
		log.Warn(logger.LogKeycloakPasswordMismatch)
		return "", domain.ErrPasswordMismatch
	}

	// Delegate to service
	email, err := i.service.UpdatePassword(ctx, token, newPassword)
	if err != nil {
		switch err {
		case domain.ErrInvalidToken:
			log.Error(logger.LogKeycloakPasswordTokenInvalid, "error", err)
		case domain.ErrPasswordUpdateFailed:
			log.Error(logger.LogKeycloakPasswordUpdateError, "error", err)
		default:
			log.Error(logger.LogKeycloakPasswordUpdateError, "error", err)
		}
		return "", err
	}

	log.Success(logger.LogKeycloakPasswordUpdateOK, "email", email)
	return email, nil
}

// UpdateEmployee actualiza la información general de un empleado
// preservando email, password y keycloak_user_id
// Sincroniza cambios de rol y estado active con Keycloak
func (i *Interactor) UpdateEmployee(ctx context.Context, employeeID string, updatedData domain.Employee) error {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

	log.Info(logger.LogEmployeeUpdateRequest, "employee_id", employeeID)

	// Step 1: Get current employee data
	currentEmployee, err := i.service.GetEmployeeByID(ctx, employeeID)
	if err != nil {
		log.Error(logger.LogEmployeeGetByIDError, "employee_id", employeeID, "error", err)
		return domain.ErrPersonNotFound
	}

	if currentEmployee == nil {
		log.Warn(logger.LogEmployeeNotFound, "employee_id", employeeID)
		return domain.ErrPersonNotFound
	}

	// Step 2: Preserve protected fields
	previousActive := currentEmployee.Active
	previousRole := currentEmployee.Role

	// Merge updated data with current data (preserving ID, email, KeycloakUserID)
	updatedData.ID = currentEmployee.ID
	updatedData.Email = currentEmployee.Email
	updatedData.KeycloakUserID = currentEmployee.KeycloakUserID

	// Step 3: Call service to update
	err = i.service.UpdateEmployee(ctx, updatedData, previousActive, previousRole)
	if err != nil {
		log.Error(logger.LogEmployeeUpdateError, "employee_id", employeeID, "error", err)
		return err
	}

	log.Success(logger.LogEmployeeUpdateComplete, "employee_id", employeeID)
	return nil
}

// ChangePassword allows an authenticated user to change their own password
// This method validates the current password and sets a new one without requiring email/token flow
func (i *Interactor) ChangePassword(ctx context.Context, email, currentPassword, newPassword, confirmPassword string) (string, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

	log.Info(logger.LogKeycloakChangePassword, "email", email)

	// Validate new passwords match
	if newPassword != confirmPassword {
		log.Warn(logger.LogKeycloakChangePasswordMismatch, "email", email)
		return "", domain.ErrPasswordMismatch
	}

	// Delegate to service (validates current password and sets new one)
	resultEmail, err := i.service.ChangePassword(ctx, email, currentPassword, newPassword)
	if err != nil {
		switch err {
		case domain.ErrInvalidCurrentPassword:
			log.Warn(logger.LogKeycloakChangePasswordInvalid, "email", email)
		case domain.ErrUserNotFound:
			log.Warn(logger.LogKeycloakUserNotFound, "email", email)
		case domain.ErrPasswordUpdateFailed:
			log.Error(logger.LogKeycloakChangePasswordError, "email", email, "error", err)
		default:
			log.Error(logger.LogKeycloakChangePasswordError, "email", email, "error", err)
		}
		return "", err
	}

	log.Success(logger.LogKeycloakChangePasswordOK, "email", resultEmail)
	return resultEmail, nil
}

// DeleteEmployee removes an employee from the system (DB and Keycloak)
// First retrieves the employee to get the keycloakUserID, then delegates to service
func (i *Interactor) DeleteEmployee(ctx context.Context, employeeID string) error {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

	log.Info(logger.LogEmployeeDeleting, "employee_id", employeeID)

	// Step 1: Get the employee to retrieve keycloakUserID
	employee, err := i.service.GetEmployeeByID(ctx, employeeID)
	if err != nil {
		log.Warn(logger.LogEmployeeNotFound, "employee_id", employeeID, "error", err)
		return domain.ErrPersonNotFound
	}

	// Step 2: Delegate to service to delete from Keycloak and DB
	if err := i.service.DeleteEmployee(ctx, employee.ID, employee.KeycloakUserID); err != nil {
		log.Error(logger.LogEmployeeDeleteError, "employee_id", employeeID, "error", err)
		return err
	}

	log.Success(logger.LogEmployeeDeleteComplete, "employee_id", employeeID, "email", employee.Email)
	return nil
}

// GetEmployeesByRole retrieves all employees for a specific role (HU47 - Virtual Entity pattern)
// This implements the "Crew Member Types" feature - returns employees filtered by role field
func (i *Interactor) GetEmployeesByRole(ctx context.Context, role string) ([]domain.Employee, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

	log.Info(logger.LogCrewMemberTypeGet, "role", role)

	employees, err := i.service.GetEmployeesByRole(ctx, role)
	if err != nil {
		log.Error(logger.LogCrewMemberTypeGetError, "role", role, "error", err)
		return nil, err
	}

	log.Success(logger.LogCrewMemberTypeGetOK, "role", role, "count", len(employees))
	return employees, nil
}
