package middleware

import (
	"errors"
	"strings"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/input"
	"github.com/champion19/flighthours-api/platform/cache/messaging"
	"github.com/champion19/flighthours-api/platform/jwt"
	"github.com/gin-gonic/gin"
)

// RequireAuth creates a middleware that validates JWT tokens from Keycloak
// and injects the authenticated user into the Gin context
func RequireAuth(employeeService input.Service, msgCache *messaging.MessageCache,jwtValidator *jwt.JWKSValidator) gin.HandlerFunc {
	tokenParser := jwt.NewTokenParser()
	_=tokenParser

	return func(c *gin.Context) {
		// Extract Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Error(domain.ErrInvalidToken)
			c.Abort()
			return
		}

		// Check Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Error(domain.ErrInvalidToken)
			c.Abort()
			return
		}

		token := parts[1]

		// Validate JWT using JWKS (signature, expiration, issuer)
		var claims map[string]interface{}
		var err error

		if jwtValidator != nil {
			// Use JWKS validation (secure path)
			claims, err = jwtValidator.ValidateToken(token)
			if err != nil {
				// Map specific JWT errors to domain errors
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
			// Fallback to simple parsing (NOT RECOMMENDED - no validation)
			// This path should only be used if JWKS initialization fails
			claims, err = tokenParser.ExtractClaimsFromToken(token)
			if err != nil {
				c.Error(domain.ErrInvalidToken)
				c.Abort()
				return
			}
		}

		// Extract Keycloak User ID from "sub" claim
		keycloakUserID, ok := claims["sub"].(string)
		if !ok || keycloakUserID == "" {
			c.Error(domain.ErrInvalidToken)
			c.Abort()
			return
		}

		// Find user in database by Keycloak ID
		employee, err := employeeService.GetEmployeeByKeycloakID(c.Request.Context(), keycloakUserID)
		if err != nil {
			// User not found in our database
			c.Error(domain.ErrUserNotFound)
			c.Abort()
			return
		}

		// Inject authenticated user into context
		c.Set("authenticated_user", employee)

		c.Next()
	}
}

// GetAuthenticatedUser extracts the authenticated user from the Gin context
func GetAuthenticatedUser(c *gin.Context) (*domain.Employee, bool) {
	user, exists := c.Get("authenticated_user")
	if !exists {
		return nil, false
	}

	employee, ok := user.(*domain.Employee)
	return employee, ok
}
// RequireRole creates a middleware that validates the user has the required role
// Must be used AFTER RequireAuth middleware
// Example usage: router.POST("/branches", RequireRole(domain.RoleRepresentative), handler.RegisterBranch())
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		employee, exists := GetAuthenticatedUser(c)
		if !exists {
			c.Error(domain.ErrUserNotFound)
			c.Abort()
			return
		}

		// Check if user's role is in the allowed roles
		for _, role := range allowedRoles {
			if employee.Role == role {
				c.Next()
				return
			}
		}

		// Role not allowed
		c.Error(domain.ErrRoleRequired)
		c.Abort()
	}
}
