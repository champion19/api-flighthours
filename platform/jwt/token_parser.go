package jwt

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
)

// Common errors for token parsing
var (
	ErrInvalidTokenFormat = errors.New("invalid token format: expected 3 parts")
	ErrPayloadDecode      = errors.New("failed to decode token payload")
	ErrClaimsParse        = errors.New("failed to parse token claims")
	ErrEmailNotFound      = errors.New("email not found in token claims")
)

// TokenParser provides JWT token parsing functionality
type TokenParser struct{}

// NewTokenParser creates a new TokenParser instance
func NewTokenParser() *TokenParser {
	return &TokenParser{}
}

// ExtractEmailFromToken extracts the email from a Keycloak JWT action token.
// The token has format: header.payload.signature
// The payload (second part) contains the field "eml" with the email for verify-email action tokens.
//
// Fields checked (in order):
//  1. "eml" - Used by Keycloak for verify-email action tokens
//  2. "email" - Used in standard access tokens
//  3. "sub" - If it has email format (some configurations use email as subject)
func (tp *TokenParser) ExtractEmailFromToken(token string) (string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", ErrInvalidTokenFormat
	}

	// Decode the payload (second part) using base64url
	payload, err := base64URLDecode(parts[1])
	if err != nil {
		return "", ErrPayloadDecode
	}

	// Parse as JSON
	var claims map[string]interface{}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return "", ErrClaimsParse
	}

	// Check "eml" field (Keycloak verify-email action tokens)
	if email, ok := claims["eml"].(string); ok && email != "" {
		return email, nil
	}

	// Fallback: check "email" field (standard access tokens)
	if email, ok := claims["email"].(string); ok && email != "" {
		return email, nil
	}

	// Fallback: check "sub" field if it looks like an email
	if sub, ok := claims["sub"].(string); ok && sub != "" {
		if isValidEmail(sub) {
			return sub, nil
		}
	}

	return "", ErrEmailNotFound
}

// ExtractClaimsFromToken extracts all claims from a Keycloak JWT token.
// This is used for authentication middleware to get user information like "sub" (Keycloak User ID).
// Returns the claims map for flexible access to any claim.
func (tp *TokenParser) ExtractClaimsFromToken(token string) (map[string]interface{}, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidTokenFormat
	}

	// Decode the payload (second part) using base64url
	payload, err := base64URLDecode(parts[1])
	if err != nil {
		return nil, ErrPayloadDecode
	}

	// Parse as JSON
	var claims map[string]interface{}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, ErrClaimsParse
	}

	return claims, nil
}

// base64URLDecode decodes base64url encoded string (without padding)
func base64URLDecode(s string) ([]byte, error) {
	// Add padding if necessary
	switch len(s) % 4 {
	case 2:
		s += "=="
	case 3:
		s += "="
	}

	// Use base64.URLEncoding for URL-safe base64
	return base64.URLEncoding.DecodeString(s)
}

// isValidEmail performs a basic email format validation
// Checks that the string contains exactly one @ with content before and after
func isValidEmail(s string) bool {
	atIndex := strings.Index(s, "@")
	if atIndex <= 0 || atIndex >= len(s)-1 {
		return false
	}
	// Ensure there's no second @
	return strings.Count(s, "@") == 1
}
