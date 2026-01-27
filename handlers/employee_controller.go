package handlers

import (
	domain "github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/middleware"
	"github.com/champion19/flighthours-api/platform/logger"
	"github.com/gin-gonic/gin"
)

// RegisterEmployee godoc
// @Summary      Registrar nueva cuenta de empleado
// @Description  Crea una nueva cuenta de empleado en el sistema con sincronización a Keycloak
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        request body handlers.EmployeeRequest true "Datos del empleado a registrar"
// @Success      201 {object} middleware.APIResponse{data=handlers.RegisterEmployeeResponse} "Cuenta creada exitosamente"
// @Failure      400 {object} middleware.APIResponse "Error de validación - Datos inválidos"
// @Failure      409 {object} middleware.APIResponse "Conflicto - Email o número de identidad ya registrado"
// @Failure      500 {object} middleware.APIResponse "Error interno del servidor"
// @Router       /register [post]
func (h handler) RegisterEmployee() func(c *gin.Context) {

	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogRegRequestReceived,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		var employeeRequest EmployeeRequest
		if err := c.ShouldBindJSON(&employeeRequest); err != nil {
			log.Error(logger.LogRegJSONParseError, "error", err, "client_ip", c.ClientIP())
			c.Error(domain.ErrInvalidJSONFormat)
			return
		}

		employeeRequest.Sanitize()

		log.Info(logger.LogRegProcessing,
			"email", employeeRequest.Email,
			"role", employeeRequest.Role)
		employeeDomain, err := employeeRequest.ToDomain()
		if err != nil {
			log.Error(logger.LogRegProcessError,
				"email", employeeRequest.Email,
				"error", err,
				"client_ip", c.ClientIP())
			switch err {
			case domain.ErrStartDateAfterEndDate:
				h.Response.Error(c, domain.MsgValStartDateAfterEndDate)
			case domain.ErrInvalidDateFormat:
				h.Response.Error(c, domain.MsgValInvalidDateFormat)
			default:
				h.Response.Error(c, domain.MsgValBadFormat)
			}
			return
		}

		result, err := h.Interactor.RegisterEmployee(c, employeeDomain)
		if err != nil {
			log.Error(logger.LogRegProcessError,
				"email", employeeRequest.Email,
				"error", err,
				"client_ip", c.ClientIP())
			c.Error(err)
			return
		}

		encodedID, err := h.EncodeID(result.Employee.ID)
		if err != nil {
			h.HandleIDEncodingError(c, result.Employee.ID, err)
			return
		}

		c.SetCookie(
			"employee_id",
			result.Employee.ID,
			3600,
			"/",
			c.Request.Host,
			c.Request.TLS != nil,
			true,
		)

		log.Success(logger.LogEmployeeRegisterSuccessLog,
			result.Employee.ToLogger(),
			"encoded_id", encodedID,
			"client_ip", c.ClientIP())

		middleware.RecordEmployeeRegistration()
		h.Response.Success(c, domain.MsgUserRegistered)
	}
}

// ResendVerificationEmail godoc
// @Summary      Reenviar email de verificación
// @Description  Reenvía el email de verificación a un usuario registrado que no ha verificado su cuenta
// @Tags         authentication
// @Accept       json
// @Produce      json
// @Param        request body handlers.ResendVerificationEmailRequest true "Email del usuario"
// @Success      200 {object} middleware.APIResponse{data=handlers.ResendVerificationEmailResponse} "Email reenviado exitosamente"
// @Failure      400 {object} middleware.APIResponse "Error de validación - Email inválido"
// @Failure      404 {object} middleware.APIResponse "Usuario no encontrado"
// @Failure      409 {object} middleware.APIResponse "Email ya verificado"
// @Failure      500 {object} middleware.APIResponse "Error interno del servidor"
// @Router       /auth/resend-verification [post]
func (h handler) ResendVerificationEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ResendVerificationEmailRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			h.Response.Error(c, domain.MsgValBadFormat)
			return
		}

		req.Sanitize()

		err := h.Interactor.ResendVerificationEmail(c, req.Email)
		if err != nil {
			switch err {
			case domain.ErrUserNotFound:
				h.Response.Error(c, domain.MsgKCUserNotFound)
			case domain.ErrEmailAlreadyVerified:
				h.Response.Warning(c, domain.MsgKCEmailAlreadyVerified)
			default:
				h.Response.Error(c, domain.MsgKCVerifEmailError)
			}
			return
		}
		h.Response.Success(c, domain.MsgKCVerifEmailResent, req.Email)
	}
}

// Login godoc
// @Summary      Login de usuario
// @Description  Autentica un usuario y retorna tokens de acceso
// @Tags         authentication
// @Accept       json
// @Produce      json
// @Param        request body handlers.LoginRequest true "Credenciales de login"
// @Success      200 {object} middleware.APIResponse{data=handlers.LoginResponse} "Login exitoso"
// @Failure      400 {object} middleware.APIResponse "Credenciales inválidas"
// @Failure      401 {object} middleware.APIResponse "Email no verificado o credenciales incorrectas"
// @Failure      500 {object} middleware.APIResponse "Error interno del servidor"
// @Router       /login [post]
func (h handler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error(logger.LogRegJSONParseError, "error", err)
			h.Response.Error(c, domain.MsgValBadFormat)
			return
		}

		req.Sanitize()

		log.Info(logger.LogKeycloakUserLogin, "email", req.Email, "client_ip", c.ClientIP())

		token, err := h.Interactor.Login(c, req.Email, req.Password)
		if err != nil {
			log.Error(logger.LogKeycloakUserLoginError, "email", req.Email, "error", err, "client_ip", c.ClientIP())

			switch err {
			case domain.ErrorEmailNotVerified:
				h.Response.Error(c, domain.MsgKCLoginEmailNotVerified)
			case domain.ErrUserNotFound:
				h.Response.Error(c, domain.MsgUnauthorized)
			default:
				h.Response.Error(c, domain.MsgUnauthorized)
			}
			return
		}

		response := LoginResponse{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			ExpiresIn:    token.ExpiresIn,
			TokenType:    token.TokenType,
		}

		log.Success(logger.LogKeycloakUserLoginOK, "email", req.Email, "client_ip", c.ClientIP())
		middleware.RecordEmployeeRegistration()
		h.Response.SuccessWithData(c, domain.MsgKCLoginSuccess, response)
	}
}

// VerifyEmailByToken godoc
// @Summary      Verificar email de usuario
// @Description  Verifica el email de un usuario usando un token JWT proxy
// @Tags         authentication
// @Accept       json
// @Produce      json
// @Param        request body handlers.VerifyEmailRequest true "Token de verificación"
// @Success      200 {object} middleware.APIResponse{data=handlers.VerifyEmailResponse} "Email verificado exitosamente"
// @Failure      400 {object} middleware.APIResponse "Token inválido o expirado"
// @Failure      404 {object} middleware.APIResponse "Usuario no encontrado"
// @Failure      409 {object} middleware.APIResponse "Email ya estaba verificado"
// @Failure      500 {object} middleware.APIResponse "Error interno del servidor"
// @Router       /auth/verify-email [post]
func (h handler) VerifyEmailByToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		var req VerifyEmailRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error(logger.LogRegJSONParseError, "error", err)
			h.Response.Error(c, domain.MsgValBadFormat)
			return
		}

		log.Info(logger.LogKeycloakEmailVerify, "client_ip", c.ClientIP())

		email, err := h.Interactor.VerifyEmailByToken(c, req.Token)
		if err != nil {
			switch err {
			case domain.ErrInvalidToken:
				h.Response.Error(c, domain.MsgKCInvalidToken)
			case domain.ErrUserNotFound:
				h.Response.Error(c, domain.MsgKCUserNotFound)
			case domain.ErrEmailAlreadyVerified:
				h.Response.Warning(c, domain.MsgKCEmailAlreadyVerified)
			default:
				h.Response.Error(c, domain.MsgKCEmailVerifyError)
			}
			return
		}

		response := VerifyEmailResponse{
			Verified: true,
			Email:    email,
		}

		log.Success(logger.LogKeycloakEmailVerifyOK, "email", email, "client_ip", c.ClientIP())
		h.Response.SuccessWithData(c, domain.MsgKCEmailVerified, response)
	}
}

// GetEmployee godoc
// @Summary      Obtener perfil del usuario autenticado
// @Description  Obtiene la información del empleado que realiza la petición usando el token JWT
// @Tags         employees
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} middleware.APIResponse{data=handlers.EmployeeResponse} "Datos del usuario autenticado"
// @Failure      401 {object} middleware.APIResponse "No autenticado"
// @Failure      500 {object} middleware.APIResponse "Error interno del servidor"
// @Router       /employee/me [get]
func (h handler) GetEmployee() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogEmployeeGetByID, "endpoint", "GET /employee/me", "client_ip", c.ClientIP())

		employee, exists := middleware.GetAuthenticatedUser(c)
		if !exists {
			log.Error(logger.LogEmployeeNotFound, "error", "authenticated user not in context", "client_ip", c.ClientIP())
			c.Error(domain.ErrUserNotFound)
			return
		}

		encodedID, err := h.EncodeID(employee.ID)
		if err != nil {
			h.HandleIDEncodingError(c, employee.ID, err)
			return
		}

		response := FromDomain(employee, encodedID)

		baseURL := GetBaseURL(c)
		response.Links = BuildEmployeeLinks(baseURL, encodedID)

		log.Success(logger.LogEmployeeGetByIDOK, "email", employee.Email, "client_ip", c.ClientIP())
		h.Response.SuccessWithData(c, domain.MsgUserFound, response)
	}
}

// ChangePassword godoc
// @Summary      Cambiar contraseña de usuario autenticado
// @Description  Permite a un usuario cambiar su contraseña conociendo la contraseña actual
// @Tags         authentication
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body handlers.ChangePasswordRequest true "Email, contraseña actual y nueva"
// @Success      200 {object} middleware.APIResponse{data=handlers.ChangePasswordResponse} "Contraseña cambiada exitosamente"
// @Failure      400 {object} middleware.APIResponse "Error de validación - Contraseñas no coinciden"
// @Failure      401 {object} middleware.APIResponse "Contraseña actual incorrecta"
// @Failure      404 {object} middleware.APIResponse "Usuario no encontrado"
// @Failure      500 {object} middleware.APIResponse "Error interno del servidor"
// @Router       /auth/change-password [post]
func (h handler) ChangePassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		var req ChangePasswordRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error(logger.LogRegJSONParseError, "error", err)
			h.Response.Error(c, domain.MsgValBadFormat)
			return
		}

		req.Sanitize()

		if req.NewPassword != req.ConfirmPassword {
			log.Warn(logger.LogKeycloakChangePasswordMismatch, "email", req.Email, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgKCPwdChangeNewMismatch)
			return
		}

		log.Info(logger.LogKeycloakChangePassword, "email", req.Email, "client_ip", c.ClientIP())

		email, err := h.Interactor.ChangePassword(c, req.Email, req.CurrentPassword, req.NewPassword, req.ConfirmPassword)
		if err != nil {
			switch err {
			case domain.ErrInvalidCurrentPassword:
				h.Response.Error(c, domain.MsgKCPwdCurrentInvalid)
			case domain.ErrUserNotFound:
				h.Response.Error(c, domain.MsgKCUserNotFound)
			case domain.ErrPasswordMismatch:
				h.Response.Error(c, domain.MsgKCPwdChangeNewMismatch)
			case domain.ErrPasswordUpdateFailed:
				h.Response.Error(c, domain.MsgKCPwdChangeError)
			default:
				h.Response.Error(c, domain.MsgKCPwdChangeError)
			}
			return
		}

		response := ChangePasswordResponse{
			Changed: true,
			Email:   email,
		}

		log.Success(logger.LogKeycloakChangePasswordOK, "email", email, "client_ip", c.ClientIP())
		h.Response.SuccessWithData(c, domain.MsgKCPwdChanged, response)
	}
}

// UpdateEmployee godoc
// @Summary      Actualizar información del empleado autenticado
// @Description  Actualiza la información básica del empleado autenticado usando el token JWT. Preserva: email, password, airline, bp (estos campos requieren otros endpoints).
// @Tags         employees
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body handlers.UpdateEmployeeRequest true "Datos a actualizar"
// @Success      200 {object} middleware.APIResponse{data=handlers.UpdateEmployeeResponse} "Empleado actualizado exitosamente"
// @Failure      400 {object} middleware.APIResponse "Error de validación - Datos inválidos"
// @Failure      401 {object} middleware.APIResponse "No autenticado"
// @Failure      404 {object} middleware.APIResponse "Empleado no encontrado"
// @Failure      500 {object} middleware.APIResponse "Error interno del servidor"
// @Router       /employees/me [put]
func (h handler) UpdateEmployee() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogEmployeeUpdateRequest, "endpoint", "PUT /employees/me", "client_ip", c.ClientIP())

		existingEmployee, exists := middleware.GetAuthenticatedUser(c)
		if !exists {
			log.Error(logger.LogEmployeeNotFound, "error", "authenticated user not in context", "client_ip", c.ClientIP())
			c.Error(domain.ErrUserNotFound)
			return
		}

		var req UpdateEmployeeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error(logger.LogRegJSONParseError, "error", err, "client_ip", c.ClientIP())
			c.Error(domain.ErrInvalidJSONFormat)
			return
		}

		req.Sanitize()

		employeeDomain, err := req.ToUpdateData(existingEmployee)
		if err != nil {
			log.Error(logger.LogEmployeeUpdateError, "email", existingEmployee.Email, "error", err, "client_ip", c.ClientIP())
			switch err {
			case domain.ErrStartDateAfterEndDate:
				h.Response.Error(c, domain.MsgValStartDateAfterEndDate)
			case domain.ErrInvalidDateFormat:
				h.Response.Error(c, domain.MsgValInvalidDateFormat)
			default:
				h.Response.Error(c, domain.MsgValBadFormat)
			}
			return
		}

		result, err := h.Interactor.UpdateEmployee(c, employeeDomain)
		if err != nil {
			log.Error(logger.LogEmployeeUpdateError, "employee_id", existingEmployee.ID, "error", err, "client_ip", c.ClientIP())
			c.Error(err)
			return
		}

		encodedID, err := h.EncodeID(result.ID)
		if err != nil {
			h.HandleIDEncodingError(c, result.ID, err)
			return
		}

		baseURL := GetBaseURL(c)
		response := UpdateEmployeeResponse{
			ID:      encodedID,
			Updated: result.Updated,
			Links:   BuildEmployeeLinks(baseURL, encodedID),
		}

		log.Success(logger.LogEmployeeUpdateComplete, "employee_id", result.ID, "email", existingEmployee.Email, "client_ip", c.ClientIP())
		h.Response.SuccessWithData(c, domain.MsgUserUpdated, response)
	}
}

// DeleteEmployee godoc
// @Summary      Eliminar cuenta del empleado autenticado
// @Description  Elimina permanentemente la cuenta del empleado autenticado del sistema (Keycloak + DB)
// @Tags         employees
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} middleware.APIResponse "Cuenta eliminada exitosamente"
// @Failure      401 {object} middleware.APIResponse "No autenticado"
// @Failure      404 {object} middleware.APIResponse "Empleado no encontrado"
// @Failure      500 {object} middleware.APIResponse "Error interno del servidor"
// @Router       /employees/me [delete]
func (h handler) DeleteEmployee() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogEmployeeDeleting, "endpoint", "DELETE /employees/me", "client_ip", c.ClientIP())

		employee, exists := middleware.GetAuthenticatedUser(c)
		if !exists {
			log.Error(logger.LogEmployeeNotFound, "error", "authenticated user not in context", "client_ip", c.ClientIP())
			c.Error(domain.ErrUserNotFound)
			return
		}

		if err := h.Interactor.DeleteEmployee(c, employee.ID); err != nil {
			log.Error(logger.LogEmployeeDeleteError, "employee_id", employee.ID, "error", err, "client_ip", c.ClientIP())
			c.Error(err)
			return
		}

		log.Success(logger.LogEmployeeDeleteComplete, "employee_id", employee.ID, "email", employee.Email, "client_ip", c.ClientIP())
		h.Response.Success(c, domain.MsgUserDeleted)
	}
}
