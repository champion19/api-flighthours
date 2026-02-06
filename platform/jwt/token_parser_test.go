package jwt

import (
	"encoding/base64"
	"encoding/json"
	"testing"
)

func createTestToken(claims map[string]interface{}) string {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payload, _ := json.Marshal(claims)
	payloadEncoded := base64.RawURLEncoding.EncodeToString(payload)
	signature := base64.RawURLEncoding.EncodeToString([]byte("test-signature"))
	return header + "." + payloadEncoded + "." + signature
}

func TestNewTokenParser(t *testing.T) {
	parser := NewTokenParser()
	if parser == nil {
		t.Error("expected non-nil TokenParser")
	}
}

func TestExtractEmailFromToken(t *testing.T) {
	parser := NewTokenParser()

	tests := []struct {
		name        string
		claims      map[string]interface{}
		expected    string
		expectError bool
	}{
		{
			name:     "email from eml claim",
			claims:   map[string]interface{}{"eml": "test@example.com"},
			expected: "test@example.com",
		},
		{
			name:     "email from email claim",
			claims:   map[string]interface{}{"email": "user@domain.com"},
			expected: "user@domain.com",
		},
		{
			name:     "email from sub claim when valid email",
			claims:   map[string]interface{}{"sub": "admin@company.com"},
			expected: "admin@company.com",
		},
		{
			name:        "no email found",
			claims:      map[string]interface{}{"name": "John"},
			expectError: true,
		},
		{
			name:        "sub is not valid email",
			claims:      map[string]interface{}{"sub": "user-id-123"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := createTestToken(tt.claims)
			email, err := parser.ExtractEmailFromToken(token)

			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if email != tt.expected {
					t.Errorf("expected %q, got %q", tt.expected, email)
				}
			}
		})
	}
}

func TestExtractEmailFromToken_InvalidFormat(t *testing.T) {
	parser := NewTokenParser()

	t.Run("invalid token format - missing parts", func(t *testing.T) {
		_, err := parser.ExtractEmailFromToken("invalid.token")
		if err != ErrInvalidTokenFormat {
			t.Errorf("expected ErrInvalidTokenFormat, got %v", err)
		}
	})

	t.Run("invalid token format - single part", func(t *testing.T) {
		_, err := parser.ExtractEmailFromToken("singlepart")
		if err != ErrInvalidTokenFormat {
			t.Errorf("expected ErrInvalidTokenFormat, got %v", err)
		}
	})

	t.Run("invalid base64 payload", func(t *testing.T) {
		_, err := parser.ExtractEmailFromToken("header.!!!invalid!!!.signature")
		if err != ErrPayloadDecode {
			t.Errorf("expected ErrPayloadDecode, got %v", err)
		}
	})

	t.Run("invalid JSON payload", func(t *testing.T) {
		invalidJSON := base64.RawURLEncoding.EncodeToString([]byte("not json"))
		_, err := parser.ExtractEmailFromToken("header." + invalidJSON + ".signature")
		if err != ErrClaimsParse {
			t.Errorf("expected ErrClaimsParse, got %v", err)
		}
	})
}

func TestExtractClaimsFromToken(t *testing.T) {
	parser := NewTokenParser()

	t.Run("extracts all claims", func(t *testing.T) {
		claims := map[string]interface{}{
			"sub":   "user-123",
			"email": "test@example.com",
			"name":  "Test User",
			"role":  "admin",
		}
		token := createTestToken(claims)

		extracted, err := parser.ExtractClaimsFromToken(token)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if extracted["sub"] != claims["sub"] {
			t.Errorf("expected sub %q, got %q", claims["sub"], extracted["sub"])
		}
		if extracted["email"] != claims["email"] {
			t.Errorf("expected email %q, got %q", claims["email"], extracted["email"])
		}
	})

	t.Run("invalid token format", func(t *testing.T) {
		_, err := parser.ExtractClaimsFromToken("invalid")
		if err != ErrInvalidTokenFormat {
			t.Errorf("expected ErrInvalidTokenFormat, got %v", err)
		}
	})

	t.Run("invalid base64 payload", func(t *testing.T) {
		_, err := parser.ExtractClaimsFromToken("header.!!!invalid!!!.signature")
		if err != ErrPayloadDecode {
			t.Errorf("expected ErrPayloadDecode, got %v", err)
		}
	})

	t.Run("invalid JSON payload", func(t *testing.T) {
		invalidJSON := base64.RawURLEncoding.EncodeToString([]byte("not json"))
		_, err := parser.ExtractClaimsFromToken("header." + invalidJSON + ".signature")
		if err != ErrClaimsParse {
			t.Errorf("expected ErrClaimsParse, got %v", err)
		}
	})
}
