package utils

import (
	"testing"
)

func TestGenerate(t *testing.T) {
	t.Run("generates valid UUID", func(t *testing.T) {
		id := Generate()
		if id == "" {
			t.Error("expected non-empty UUID")
		}
		if !IsValid(id) {
			t.Errorf("generated ID %q is not a valid UUID", id)
		}
	})

	t.Run("generates unique UUIDs", func(t *testing.T) {
		id1 := Generate()
		id2 := Generate()
		if id1 == id2 {
			t.Error("expected unique UUIDs, got same value")
		}
	})
}

func TestIsValid(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		expected bool
	}{
		{"valid UUID", "550e8400-e29b-41d4-a716-446655440000", true},
		{"valid UUID lowercase", "550e8400-e29b-41d4-a716-446655440000", true},
		{"valid UUID uppercase", "550E8400-E29B-41D4-A716-446655440000", true},
		{"empty string", "", false},
		{"random string", "not-a-uuid", false},
		{"too short", "550e8400-e29b-41d4", false},
		{"invalid characters", "550e8400-e29b-41d4-a716-44665544ZZZZ", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValid(tt.id)
			if result != tt.expected {
				t.Errorf("IsValid(%q) = %v, expected %v", tt.id, result, tt.expected)
			}
		})
	}
}

func TestParse(t *testing.T) {
	t.Run("valid UUID", func(t *testing.T) {
		id := "550e8400-e29b-41d4-a716-446655440000"
		parsed, err := Parse(id)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if parsed.String() != id {
			t.Errorf("expected %q, got %q", id, parsed.String())
		}
	})

	t.Run("invalid UUID returns error", func(t *testing.T) {
		_, err := Parse("invalid-uuid")
		if err == nil {
			t.Error("expected error for invalid UUID")
		}
	})
}

func TestMustParse(t *testing.T) {
	t.Run("valid UUID", func(t *testing.T) {
		id := "550e8400-e29b-41d4-a716-446655440000"
		parsed := MustParse(id)
		if parsed.String() != id {
			t.Errorf("expected %q, got %q", id, parsed.String())
		}
	})

	t.Run("invalid UUID panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for invalid UUID")
			}
		}()
		MustParse("invalid-uuid")
	})
}
