package handlers

import (
	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/middleware"
	"github.com/champion19/api-flighthours/platform/logger"
	"github.com/gin-gonic/gin"
)

// RegisterEmployee godoc
// @Summary      Register a new employee
// @Description  Registers a new employee in the system and sends a verification email
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        employee body EmployeeRequest true "Employee registration data"
// @Success      201  {object}  middleware.APIResponse "Employee registered successfully"
// @Failure      400  {object}  middleware.ErrorResponse "Invalid request data"
// @Failure      409  {object}  middleware.ErrorResponse "Email or ID already exists"
// @Failure      500  {object}  middleware.ErrorResponse "Internal server error"
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
		employeeDomain := employeeRequest.ToDomain()

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
// @Summary      Resend verification email
// @Description  Resends the verification email to a registered user who hasn't verified their account
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body ResendVerificationEmailRequest true "Email to resend verification"
// @Success      200  {object}  ResendVerificationEmailResponse "Email sent"
// @Failure      400  {object}  middleware.ErrorResponse "Invalid email"
// @Failure      404  {object}  middleware.ErrorResponse "User not found"
// @Failure      409  {object}  middleware.ErrorResponse "Email already verified"
// @Failure      500  {object}  middleware.ErrorResponse "Internal server error"
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
// @Summary      User login
// @Description  Authenticates a user and returns access tokens
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        credentials body LoginRequest true "Login credentials"
// @Success      200  {object}  LoginResponse "Login successful"
// @Failure      400  {object}  middleware.ErrorResponse "Invalid credentials format"
// @Failure      401  {object}  middleware.ErrorResponse "Email not verified or incorrect credentials"
// @Failure      500  {object}  middleware.ErrorResponse "Internal server error"
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
// @Summary      Verify user email
// @Description  Verifies the user's email using a JWT proxy token
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body VerifyEmailRequest true "Verification token"
// @Success      200  {object}  VerifyEmailResponse "Email verified"
// @Failure      400  {object}  middleware.ErrorResponse "Invalid or expired token"
// @Failure      404  {object}  middleware.ErrorResponse "User not found"
// @Failure      409  {object}  middleware.ErrorResponse "Email already verified"
// @Failure      500  {object}  middleware.ErrorResponse "Internal server error"
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
// @Summary      Get authenticated user profile
// @Description  Returns the profile information of the authenticated employee
// @Tags         Employees
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  EmployeeResponse "User profile"
// @Failure      401  {object}  middleware.ErrorResponse "Not authenticated"
// @Failure      500  {object}  middleware.ErrorResponse "Internal server error"
// @Router       /employees [get]
func (h handler) GetEmployee() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogEmployeeGetByID, "endpoint", "GET /employees", "client_ip", c.ClientIP())

		employee, exists := middleware.GetAuthenticatedUser(c)
		if !exists {
			log.Error(logger.LogEmployeeNotFound, "error", logger.LogErrAuthUserNotInContext, "client_ip", c.ClientIP())
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
// @Summary      Change authenticated user password
// @Description  Allows an authenticated user to change their password by providing the current password
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body ChangePasswordRequest true "Current and new password"
// @Success      200  {object}  ChangePasswordResponse "Password changed"
// @Failure      400  {object}  middleware.ErrorResponse "Passwords don't match"
// @Failure      401  {object}  middleware.ErrorResponse "Current password incorrect"
// @Failure      404  {object}  middleware.ErrorResponse "User not found"
// @Failure      500  {object}  middleware.ErrorResponse "Internal server error"
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
// @Summary      Update authenticated user information
// @Description  Updates the basic information of the authenticated employee. Preserves: email, password, airline, bp.
// @Tags         Employees
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body UpdateEmployeeRequest true "Employee data to update"
// @Success      200  {object}  UpdateEmployeeResponse "Employee updated"
// @Failure      400  {object}  middleware.ErrorResponse "Invalid data"
// @Failure      401  {object}  middleware.ErrorResponse "Not authenticated"
// @Failure      404  {object}  middleware.ErrorResponse "Employee not found"
// @Failure      500  {object}  middleware.ErrorResponse "Internal server error"
// @Router       /employees [put]
func (h handler) UpdateEmployee() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogEmployeeUpdateRequest, "endpoint", "PUT /employees", "client_ip", c.ClientIP())

		existingEmployee, exists := middleware.GetAuthenticatedUser(c)
		if !exists {
			log.Error(logger.LogEmployeeNotFound, "error", logger.LogErrAuthUserNotInContext, "client_ip", c.ClientIP())
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

		employeeDomain := req.ToUpdateData(existingEmployee)

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
// @Summary      Delete authenticated user account
// @Description  Permanently deletes the authenticated employee's account from the system (Keycloak + DB)
// @Tags         Employees
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  middleware.APIResponse "Account deleted"
// @Failure      401  {object}  middleware.ErrorResponse "Not authenticated"
// @Failure      404  {object}  middleware.ErrorResponse "Employee not found"
// @Failure      500  {object}  middleware.ErrorResponse "Internal server error"
// @Router       /employees [delete]
func (h handler) DeleteEmployee() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogEmployeeDeleting, "endpoint", "DELETE /employees", "client_ip", c.ClientIP())

		employee, exists := middleware.GetAuthenticatedUser(c)
		if !exists {
			log.Error(logger.LogEmployeeNotFound, "error", logger.LogErrAuthUserNotInContext, "client_ip", c.ClientIP())
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

// PasswordReset godoc
// @Summary      Request password reset
// @Description  Sends an email with instructions to reset the password. Always returns success for security.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body PasswordResetRequest true "Email for password reset"
// @Success      200  {object}  PasswordResetResponse "Reset email sent"
// @Failure      400  {object}  middleware.ErrorResponse "Invalid email"
// @Failure      500  {object}  middleware.ErrorResponse "Internal server error"
// @Router       /auth/password-reset [post]
func (h handler) PasswordReset() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		var req PasswordResetRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error(logger.LogRegJSONParseError, "error", err)
			h.Response.Error(c, domain.MsgValBadFormat)
			return
		}

		req.Sanitize()

		log.Info(logger.LogKeycloakSendPasswordReset, "email", req.Email, "client_ip", c.ClientIP())

		_ = h.Interactor.RequestPasswordReset(c, req.Email)

		response := PasswordResetResponse{
			Sent: true,
		}

		log.Success(logger.LogKeycloakSendPasswordResetOK, "email", req.Email, "client_ip", c.ClientIP())
		h.Response.SuccessWithData(c, domain.MsgKCPwdResetSent, response)
	}
}

// UpdatePassword godoc
// @Summary      Update password with reset token
// @Description  Allows updating the password using a token received by email. Token expires in 12 hours.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body UpdatePasswordRequest true "Token and new password"
// @Success      200  {object}  UpdatePasswordResponse "Password updated"
// @Failure      400  {object}  middleware.ErrorResponse "Passwords don't match or invalid token"
// @Failure      401  {object}  middleware.ErrorResponse "Expired or invalid token"
// @Failure      500  {object}  middleware.ErrorResponse "Internal server error"
// @Router       /auth/update-password [post]
func (h handler) UpdatePassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		var req UpdatePasswordRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error(logger.LogRegJSONParseError, "error", err)
			h.Response.Error(c, domain.MsgValBadFormat)
			return
		}

		if req.NewPassword != req.ConfirmPassword {
			log.Warn(logger.LogKeycloakPasswordMismatch, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgKCPwdMismatch)
			return
		}

		log.Info(logger.LogKeycloakPasswordUpdate, "client_ip", c.ClientIP())

		email, err := h.Interactor.UpdatePassword(c, req.Token, req.NewPassword, req.ConfirmPassword)
		if err != nil {
			switch err {
			case domain.ErrInvalidToken, domain.ErrTokenExpired:
				h.Response.Error(c, domain.MsgKCPwdUpdateTokenInvalid)
			case domain.ErrPasswordMismatch:
				h.Response.Error(c, domain.MsgKCPwdMismatch)
			case domain.ErrPasswordUpdateFailed:
				h.Response.Error(c, domain.MsgKCPwdUpdateError)
			default:
				h.Response.Error(c, domain.MsgKCPwdUpdateError)
			}
			return
		}

		response := UpdatePasswordResponse{
			Updated: true,
			Email:   email,
		}

		log.Success(logger.LogKeycloakPasswordUpdateOK, "email", email, "client_ip", c.ClientIP())
		h.Response.SuccessWithData(c, domain.MsgKCPwdUpdated, response)
	}
}
