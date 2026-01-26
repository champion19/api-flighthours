package jwt

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
)

var (
	ErrInvalidTokenFormat = errors.New("invalid token format: expected 3 parts")
	ErrPayloadDecode      = errors.New("failed to decode token payload")
	ErrClaimsParse        = errors.New("failed to parse token claims")
	ErrEmailNotFound      = errors.New("email not found in token claims")
)

type TokenParser struct{}

func NewTokenParser() *TokenParser {
	return &TokenParser{}
}



func (tp *TokenParser) ExtractEmailFromToken(token string) (string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", ErrInvalidTokenFormat
	}

	payload, err := base64URLDecode(parts[1])
	if err != nil {
		return "", ErrPayloadDecode
	}

	var claims map[string]interface{}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return "", ErrClaimsParse
	}


	if email, ok := claims["eml"].(string); ok && email != "" {
		return email, nil
	}


	if email, ok := claims["email"].(string); ok && email != "" {
		return email, nil
	}


	if sub, ok := claims["sub"].(string); ok && sub != "" {
		if isValidEmail(sub) {
			return sub, nil
		}
	}

	return "", ErrEmailNotFound
}

func (tp *TokenParser) ExtractClaimsFromToken(token string) (map[string]interface{}, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidTokenFormat
	}

	payload, err := base64URLDecode(parts[1])
	if err != nil {
		return nil, ErrPayloadDecode
	}
	var claims map[string]interface{}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, ErrClaimsParse
	}

	return claims, nil
}

func base64URLDecode(s string) ([]byte, error) {
	switch len(s) % 4 {
	case 2:
		s += "=="
	case 3:
		s += "="
	}


	return base64.URLEncoding.DecodeString(s)
}


func isValidEmail(s string) bool {
	atIndex := strings.Index(s, "@")
	if atIndex <= 0 || atIndex >= len(s)-1 {
		return false
	}

	return strings.Count(s, "@") == 1
}
