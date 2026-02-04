package handlers

import (
	"testing"
)

func TestTrimString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"trims leading spaces", "  hello", "hello"},
		{"trims trailing spaces", "world  ", "world"},
		{"trims both sides", "  test  ", "test"},
		{"handles empty string", "", ""},
		{"handles no spaces", "clean", "clean"},
		{"handles tabs and newlines", "\t\nhello\n\t", "hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TrimString(tt.input)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestTrimStringPtr(t *testing.T) {
	t.Run("returns nil for nil input", func(t *testing.T) {
		result := TrimStringPtr(nil)
		if result != nil {
			t.Error("expected nil result")
		}
	})

	t.Run("trims whitespace from pointer", func(t *testing.T) {
		input := "  hello world  "
		result := TrimStringPtr(&input)
		if result == nil {
			t.Fatal("expected non-nil result")
		}
		if *result != "hello world" {
			t.Errorf("expected 'hello world', got %q", *result)
		}
	})

	t.Run("handles empty string pointer", func(t *testing.T) {
		input := ""
		result := TrimStringPtr(&input)
		if result == nil {
			t.Fatal("expected non-nil result")
		}
		if *result != "" {
			t.Errorf("expected empty string, got %q", *result)
		}
	})
}
