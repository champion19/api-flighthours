package jwt
// Package jwt provides JWT token parsing and validation utilities.
// This package handles the validation of Keycloak-issued JWTs using JWKS.

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
)

// Custom errors for token validation
var (
	ErrTokenExpired     = errors.New("token has expired")
	ErrTokenNotValidYet = errors.New("token is not valid yet")
	ErrInvalidSignature = errors.New("token signature is invalid")
	ErrInvalidIssuer    = errors.New("token issuer is invalid")
	ErrInvalidClaims    = errors.New("token claims are invalid")
	ErrJWKSUnavailable  = errors.New("JWKS endpoint is unavailable")
	ErrTokenMalformed   = errors.New("token is malformed")
)

// JWKSValidator provides JWT validation using Keycloak's JWKS endpoint
type JWKSValidator struct {
	jwks           *keyfunc.JWKS
	expectedIssuer string
}

// JWKSConfig holds configuration for JWKS validator
type JWKSConfig struct {
	JWKSURL         string        // JWKS endpoint URL (e.g., http://localhost:8080/realms/myrealm/protocol/openid-connect/certs)
	Issuer          string        // Expected issuer (e.g., http://localhost:8080/realms/myrealm)
	RefreshInterval time.Duration // How often to refresh JWKS (default: 1 hour)
}

// NewJWKSValidator creates a new JWKS-based JWT validator
// It fetches the public keys from Keycloak's JWKS endpoint and caches them
func NewJWKSValidator(ctx context.Context, config JWKSConfig) (*JWKSValidator, error) {
	if config.JWKSURL == "" {
		return nil, errors.New("JWKS URL cannot be empty")
	}

	if config.RefreshInterval == 0 {
		config.RefreshInterval = time.Hour
	}

	// Configure JWKS options
	options := keyfunc.Options{
		Ctx: ctx,
		RefreshErrorHandler: func(err error) {
			// Log JWKS refresh errors (silent fail, use cached keys)
			fmt.Printf("JWKS refresh error: %v\n", err)
		},
		RefreshInterval:   config.RefreshInterval,
		RefreshRateLimit:  time.Minute * 5,
		RefreshTimeout:    time.Second * 10,
		RefreshUnknownKID: true,
	}

	// Fetch JWKS from Keycloak
	jwks, err := keyfunc.Get(config.JWKSURL, options)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrJWKSUnavailable, err)
	}

	return &JWKSValidator{
		jwks:           jwks,
		expectedIssuer: config.Issuer,
	}, nil
}

// ValidateToken validates a JWT token and returns the claims if valid
// It checks:
// - Token signature using Keycloak's public keys
// - Token expiration (exp claim)
// - Token issuer (iss claim) matches expected issuer
// - Token not-before time (nbf claim) if present
func (v *JWKSValidator) ValidateToken(tokenString string) (map[string]interface{}, error) {
	// Parse and validate token
	token, err := jwt.Parse(tokenString, v.jwks.Keyfunc)
	if err != nil {
		// Check for specific error types
		var validationErr *jwt.ValidationError
		if errors.As(err, &validationErr) {
			switch {
			case validationErr.Errors&jwt.ValidationErrorExpired != 0:
				return nil, ErrTokenExpired
			case validationErr.Errors&jwt.ValidationErrorNotValidYet != 0:
				return nil, ErrTokenNotValidYet
			case validationErr.Errors&jwt.ValidationErrorSignatureInvalid != 0:
				return nil, ErrInvalidSignature
			case validationErr.Errors&jwt.ValidationErrorMalformed != 0:
				return nil, ErrTokenMalformed
			}
		}
		return nil, fmt.Errorf("%w: %v", ErrInvalidClaims, err)
	}

	if !token.Valid {
		return nil, ErrInvalidSignature
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidClaims
	}

	// Validate issuer if configured
	if v.expectedIssuer != "" {
		issuer, ok := claims["iss"].(string)
		if !ok || issuer != v.expectedIssuer {
			return nil, ErrInvalidIssuer
		}
	}

	// Convert to map[string]interface{} for compatibility
	result := make(map[string]interface{})
	for k, val := range claims {
		result[k] = val
	}

	return result, nil
}

// Close releases resources used by the JWKS validator
func (v *JWKSValidator) Close() {
	if v.jwks != nil {
		v.jwks.EndBackground()
	}
}
