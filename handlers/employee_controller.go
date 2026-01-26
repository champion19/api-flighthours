package handlers

import (
	domain "github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/middleware"
	"github.com/champion19/flighthours-api/platform/logger"
	"github.com/gin-gonic/gin"
)

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

		log.Success("register employee success",
			result.Employee.ToLogger(),
			"encoded_id", encodedID,
			"client_ip", c.ClientIP())

		middleware.RecordEmployeeRegistration()
		h.Response.Success(c, domain.MsgUserRegistered)
	}
}

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

// UpdateEmployee - PUT /employees/me - HU23
// Updates the authenticated employee's basic information
// Preserves: email, password, airline, bp, keycloak_user_id
func (h handler) UpdateEmployee() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogEmployeeUpdateRequest, "endpoint", "PUT /employees/me", "client_ip", c.ClientIP())

		// Get authenticated user from context
		existingEmployee, exists := middleware.GetAuthenticatedUser(c)
		if !exists {
			log.Error(logger.LogEmployeeNotFound, "error", "authenticated user not in context", "client_ip", c.ClientIP())
			c.Error(domain.ErrUserNotFound)
			return
		}

		// Bind and sanitize request
		var req UpdateEmployeeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error(logger.LogRegJSONParseError, "error", err, "client_ip", c.ClientIP())
			c.Error(domain.ErrInvalidJSONFormat)
			return
		}

		req.Sanitize()

		// Convert request to domain, preserving immutable fields from existing employee
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

		// Call interactor
		result, err := h.Interactor.UpdateEmployee(c, employeeDomain)
		if err != nil {
			log.Error(logger.LogEmployeeUpdateError, "employee_id", existingEmployee.ID, "error", err, "client_ip", c.ClientIP())
			c.Error(err)
			return
		}

		// Encode ID for response
		encodedID, err := h.EncodeID(result.ID)
		if err != nil {
			h.HandleIDEncodingError(c, result.ID, err)
			return
		}

		// Build response with HATEOAS links
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
