package output

// TokenValidator defines the contract for JWT token validation.
// The core/middleware depends on this interface, while the concrete
// implementation (JWKSValidator) lives in platform/jwt.
type TokenValidator interface {
	// ValidateToken validates a JWT token string and returns its claims.
	// Returns an error if the token is expired, malformed, or invalid.
	ValidateToken(tokenString string) (map[string]interface{}, error)

	// Close gracefully shuts down the validator (e.g., stops JWKS refresh).
	Close()
}
