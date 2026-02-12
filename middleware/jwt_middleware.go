package middleware

import (
	"errors"
	"strings"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/input"
	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/cache/messaging"
	"github.com/champion19/api-flighthours/platform/jwt"
	"github.com/gin-gonic/gin"
)

func RequireAuth(employeeService input.Service, msgCache *messaging.MessageCache, jwtValidator output.TokenValidator) gin.HandlerFunc {
	tokenParser := jwt.NewTokenParser()
	_ = tokenParser

	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Error(domain.ErrInvalidToken)
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Error(domain.ErrInvalidToken)
			c.Abort()
			return
		}

		token := parts[1]

		var claims map[string]interface{}
		var err error

		if jwtValidator != nil {

			claims, err = jwtValidator.ValidateToken(token)
			if err != nil {

				switch {
				case errors.Is(err, jwt.ErrTokenExpired):
					c.Error(domain.ErrTokenExpired)
				case errors.Is(err, jwt.ErrInvalidSignature):
					c.Error(domain.ErrInvalidToken)
				case errors.Is(err, jwt.ErrInvalidIssuer):
					c.Error(domain.ErrInvalidToken)
				default:
					c.Error(domain.ErrInvalidToken)
				}
				c.Abort()
				return
			}
		} else {
			claims, err = tokenParser.ExtractClaimsFromToken(token)
			if err != nil {
				c.Error(domain.ErrInvalidToken)
				c.Abort()
				return
			}
		}

		keycloakUserID, ok := claims["sub"].(string)
		if !ok || keycloakUserID == "" {
			c.Error(domain.ErrInvalidToken)
			c.Abort()
			return
		}

		employee, err := employeeService.GetEmployeeByKeycloakID(c.Request.Context(), keycloakUserID)
		if err != nil {
			c.Error(domain.ErrUserNotFound)
			c.Abort()
			return
		}

		c.Set("authenticated_user", employee)

		c.Next()
	}
}

func GetAuthenticatedUser(c *gin.Context) (*domain.Employee, bool) {
	user, exists := c.Get("authenticated_user")
	if !exists {
		return nil, false
	}

	employee, ok := user.(*domain.Employee)
	return employee, ok
}

func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		employee, exists := GetAuthenticatedUser(c)
		if !exists {
			c.Error(domain.ErrUserNotFound)
			c.Abort()
			return
		}

		for _, role := range allowedRoles {
			if employee.Role == role {
				c.Next()
				return
			}
		}

		c.Error(domain.ErrRoleRequired)
		c.Abort()
	}
}
