package handlers

import (
	domain "github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/middleware"
	"github.com/champion19/flighthours-api/platform/logger"
	"github.com/gin-gonic/gin"
)

// RegisterEmployee godoc
// @Summary      Registrar nueva cuenta de empleado
// @Description  Crea una nueva cuenta de empleado en el sistema con sincronización a Keycloak. Incluye validación de datos, verificación de duplicados y creación de usuario en el sistema de autenticación.
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        account  body      EmployeeRequest  true  "Datos del empleado a registrar"
// @Success      201      {object}  middleware.APIResponse{data=RegisterEmployeeResponse}  "Cuenta creada exitosamente"
// @Failure      400      {object}  middleware.APIResponse  "Error de validación - Datos inválidos o incompletos"
// @Failure      409      {object}  middleware.APIResponse  "Conflicto - Email o número de identidad ya registrado"
// @Failure      500      {object}  middleware.APIResponse  "Error interno del servidor"
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

		// Sanitize input data (trim whitespace)
		employeeRequest.Sanitize()

		log.Info(logger.LogRegProcessing,
			"email", employeeRequest.Email,
			"role", employeeRequest.Role)

		// Convert to domain and validate dates
		employeeDomain, err := employeeRequest.ToDomain()
		if err != nil {
			log.Error(logger.LogRegProcessError,
				"email", employeeRequest.Email,
				"error", err,
				"client_ip", c.ClientIP())
			// Handle specific validation errors
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

		// Ofuscar el ID antes de exponerlo en la API
		encodedID, err := h.EncodeID(result.Employee.ID)
		if err != nil {
			h.HandleIDEncodingError(c, result.Employee.ID, err)
			return
		}

		//TODO;TENERLO EN CUENTA, ESTO ES DE COOKIES HTTTPONLY
		c.SetCookie(
			"employee_id",        // name
			result.Employee.ID,   // value
			3600,                 // expira en 1 hora
			"/",                  // path
			c.Request.Host,       // domain
			c.Request.TLS != nil, // secure
			true,                 // httpOnly
		)

		log.Success("register employee success",
			result.Employee.ToLogger(),
			"encoded_id", encodedID,
			"client_ip", c.ClientIP())

		// Record Prometheus metric for employee registration
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
// @Param        request  body      ResendVerificationEmailRequest  true  "Email del usuario"
// @Success      200      {object}  middleware.APIResponse{data=ResendVerificationEmailResponse}  "Email reenviado exitosamente"
// @Failure      400      {object}  middleware.APIResponse  "Error de validación - Email inválido"
// @Failure      404      {object}  middleware.APIResponse  "Usuario no encontrado"
// @Failure      500      {object}  middleware.APIResponse  "Error interno del servidor"
// @Router       /auth/resend-verification [post]
func (h handler) ResendVerificationEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ResendVerificationEmailRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			h.Response.Error(c, domain.MsgValBadFormat)
			return
		}

		// Sanitize input data
		req.Sanitize()

		err := h.Interactor.ResendVerificationEmail(c, req.Email)
		if err != nil {
			// Manejar diferentes tipos de errores
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

// @Summary Login de usuario
// @Description Autentica un usuario y retorna tokens de acceso
// @Tags Autenticación
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Credenciales de login"
// @Success 200 {object} middleware.APIResponse{data=LoginResponse} "Login exitoso"
// @Failure 400 {object} middleware.APIResponse "Credenciales inválidas"
// @Failure 401 {object} middleware.APIResponse "Email no verificado o credenciales incorrectas"
// @Failure 500 {object} middleware.APIResponse "Error interno del servidor"
// @Router /auth/login [post]
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

		// Sanitize input data (email only, password preserved as-is)
		req.Sanitize()

		log.Info(logger.LogKeycloakUserLogin, "email", req.Email, "client_ip", c.ClientIP())

		token, err := h.Interactor.Login(c, req.Email, req.Password)
		if err != nil {
			log.Error(logger.LogKeycloakUserLoginError, "email", req.Email, "error", err, "client_ip", c.ClientIP())

			// Handle specific errors with appropriate messages
			switch err {
			case domain.ErrorEmailNotVerified:
				// Email not verified - verification email was resent automatically
				h.Response.Error(c, domain.MsgKCLoginEmailNotVerified)
			case domain.ErrUserNotFound:
				// User not found - return generic unauthorized for security
				h.Response.Error(c, domain.MsgUnauthorized)
			default:
				// Other errors (invalid credentials, etc.)
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

// @Summary Verificar email de usuario (Proxy)
// @Description Verifica el email de un usuario usando un token JWT. Este endpoint actúa como proxy para no exponer Keycloak directamente.
// @Tags Autenticación
// @Accept json
// @Produce json
// @Param request body VerifyEmailRequest true "Token de verificación del email"
// @Success 200 {object} middleware.APIResponse{data=VerifyEmailResponse} "Email verificado exitosamente"
// @Failure 400 {object} middleware.APIResponse "Token inválido o expirado"
// @Failure 404 {object} middleware.APIResponse "Usuario no encontrado"
// @Failure 409 {object} middleware.APIResponse "Email ya estaba verificado"
// @Failure 500 {object} middleware.APIResponse "Error interno del servidor"
// @Router /auth/verify-email [post]
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

		// Pasar el token al Interactor - la extracción del email se hace en la capa de negocio
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

// GetMe godoc
// @Summary      Obtener perfil del usuario autenticado
// @Description  Obtiene la información del empleado que realiza la petición usando el token JWT. No requiere pasar ID.
// @Tags         employees
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  middleware.APIResponse{data=EmployeeResponse}  "Datos del usuario autenticado"
// @Failure      401  {object}  middleware.APIResponse  "No autenticado"
// @Failure      500  {object}  middleware.APIResponse  "Error interno del servidor"
// @Router       /employee/me [get]
func (h handler) GetMe() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogEmployeeGetByID, "endpoint", "GET /employee/me", "client_ip", c.ClientIP())

		// Extraer usuario autenticado del contexto (ya validado por middleware RequireAuth)
		employee, exists := middleware.GetAuthenticatedUser(c)
		if !exists {
			log.Error(logger.LogEmployeeNotFound, "error", "authenticated user not in context", "client_ip", c.ClientIP())
			c.Error(domain.ErrUserNotFound)
			return
		}

		// Ofuscar ID para la respuesta
		encodedID, err := h.EncodeID(employee.ID)
		if err != nil {
			h.HandleIDEncodingError(c, employee.ID, err)
			return
		}

		// Convertir a EmployeeResponse (sin contraseña)
		response := FromDomain(employee, encodedID)

		// Build HATEOAS links
		baseURL := GetBaseURL(c)
		response.Links = BuildEmployeeLinks(baseURL, encodedID)

		log.Success(logger.LogEmployeeGetByIDOK, "email", employee.Email, "client_ip", c.ClientIP())
		h.Response.SuccessWithData(c, domain.MsgUserFound, response)
	}
}

// ChangePassword godoc
// @Summary      Cambiar contraseña de usuario autenticado
// @Description  Permite a un usuario cambiar su contraseña conociendo la contraseña actual. Este flujo no requiere salir de la API ni tokens por email.
// @Tags         authentication
// @Accept       json
// @Produce      json
// @Param        request  body      ChangePasswordRequest  true  "Email, contraseña actual y nueva contraseña"
// @Success      200      {object}  middleware.APIResponse{data=ChangePasswordResponse}  "Contraseña cambiada exitosamente"
// @Failure      400      {object}  middleware.APIResponse  "Error de validación - Contraseñas no coinciden"
// @Failure      401      {object}  middleware.APIResponse  "Contraseña actual incorrecta"
// @Failure      404      {object}  middleware.APIResponse  "Usuario no encontrado"
// @Failure      500      {object}  middleware.APIResponse  "Error interno del servidor"
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

		// Sanitize input data (email only, passwords preserved as-is)
		req.Sanitize()

		// Validate new passwords match
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
