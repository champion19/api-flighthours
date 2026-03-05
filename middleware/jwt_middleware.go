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

	return func(c *gin.Context) {
		token, err := extractBearerToken(c)
		if err != nil {
			c.Error(domain.ErrInvalidToken)
			c.Abort()
			return
		}

		claims, err := resolveTokenClaims(token, jwtValidator, tokenParser)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
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

func extractBearerToken(c *gin.Context) (string, error) {
	// 1. Try Authorization header first (mobile / API clients)
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1], nil
		}
		return "", domain.ErrInvalidToken
	}

	// 2. Fallback: try fh_access_token cookie (web clients)
	if cookieToken, err := c.Cookie("fh_access_token"); err == nil && cookieToken != "" {
		return cookieToken, nil
	}

	return "", domain.ErrInvalidToken
}

func resolveTokenClaims(token string, jwtValidator output.TokenValidator, tokenParser *jwt.TokenParser) (map[string]interface{}, error) {
	if jwtValidator == nil {
		claims, err := tokenParser.ExtractClaimsFromToken(token)
		if err != nil {
			return nil, domain.ErrInvalidToken
		}
		return claims, nil
	}

	claims, err := jwtValidator.ValidateToken(token)
	if err != nil {
		return nil, mapTokenError(err)
	}
	return claims, nil
}

func mapTokenError(err error) error {
	if errors.Is(err, jwt.ErrTokenExpired) {
		return domain.ErrTokenExpired
	}
	return domain.ErrInvalidToken
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
