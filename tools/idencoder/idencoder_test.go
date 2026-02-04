package idencoder

import (
	"testing"

	"github.com/champion19/api-flighthours/platform/logger"
	"github.com/google/uuid"
)

func TestHashidsEncoder_RoundTrip(t *testing.T) {
	enc, err := NewHashidsEncoder(Config{Secret: "secret", MinLength: 10}, nil)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	id := uuid.New().String()
	encoded, err := enc.Encode(id)
	if err != nil {
		t.Fatalf("encode err: %v", err)
	}

	decoded, err := enc.Decode(encoded)
	if err != nil {
		t.Fatalf("decode err: %v", err)
	}
	if decoded != id {
		t.Fatalf("expected %s got %s", id, decoded)
	}
}

func TestHashidsEncoder_InvalidInputs(t *testing.T) {
	enc, err := NewHashidsEncoder(Config{Secret: "secret", MinLength: 10}, nil)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if _, err := enc.Encode("not-a-uuid"); err == nil {
		t.Fatalf("expected error for invalid uuid")
	}
	if _, err := enc.Decode(""); err == nil {
		t.Fatalf("expected error for empty encoded")
	}
}

func TestHashidsEncoder_NewHashidsEncoder_RequiresSecret(t *testing.T) {
	if _, err := NewHashidsEncoder(Config{Secret: "", MinLength: 10}, nil); err == nil {
		t.Fatalf("expected error when secret empty")
	}
}

func TestHashidsEncoder_MustEncode(t *testing.T) {
	enc, _ := NewHashidsEncoder(Config{Secret: "secret", MinLength: 10}, nil)

	t.Run("returns encoded string for valid UUID", func(t *testing.T) {
		id := uuid.New().String()
		encoded := enc.MustEncode(id)
		if encoded == "" {
			t.Error("expected non-empty encoded string")
		}
	})

	t.Run("returns empty string for invalid UUID", func(t *testing.T) {
		encoded := enc.MustEncode("not-a-uuid")
		if encoded != "" {
			t.Errorf("expected empty string, got %q", encoded)
		}
	})
}

func TestHashidsEncoder_IsValidEncoded(t *testing.T) {
	enc, _ := NewHashidsEncoder(Config{Secret: "secret", MinLength: 10}, nil)

	t.Run("returns true for valid encoded ID", func(t *testing.T) {
		id := uuid.New().String()
		encoded, _ := enc.Encode(id)
		if !enc.IsValidEncoded(encoded) {
			t.Error("expected IsValidEncoded to return true")
		}
	})

	t.Run("returns false for invalid encoded ID", func(t *testing.T) {
		if enc.IsValidEncoded("invalid-encoded") {
			t.Error("expected IsValidEncoded to return false for invalid input")
		}
	})

	t.Run("returns false for empty string", func(t *testing.T) {
		if enc.IsValidEncoded("") {
			t.Error("expected IsValidEncoded to return false for empty string")
		}
	})
}

func TestIsUUID(t *testing.T) {
	t.Run("returns true for valid UUID", func(t *testing.T) {
		validUUID := uuid.New().String()
		if !IsUUID(validUUID) {
			t.Error("expected true for valid UUID")
		}
	})

	t.Run("returns true for uppercase UUID", func(t *testing.T) {
		validUUID := "550E8400-E29B-41D4-A716-446655440000"
		if !IsUUID(validUUID) {
			t.Error("expected true for uppercase UUID")
		}
	})

	t.Run("returns false for invalid UUID", func(t *testing.T) {
		if IsUUID("not-a-uuid") {
			t.Error("expected false for invalid UUID")
		}
	})

	t.Run("returns false for empty string", func(t *testing.T) {
		if IsUUID("") {
			t.Error("expected false for empty string")
		}
	})
}

// mockLogger for testing logger branches
type mockLogger struct{}

func (m *mockLogger) Info(msg string, args ...interface{})     {}
func (m *mockLogger) Error(msg string, args ...interface{})    {}
func (m *mockLogger) Debug(msg string, args ...interface{})    {}
func (m *mockLogger) Warn(msg string, args ...interface{})     {}
func (m *mockLogger) Success(msg string, args ...interface{})  {}
func (m *mockLogger) Fatal(msg string, args ...interface{})    {}
func (m *mockLogger) Panic(msg string, args ...interface{})    {}
func (m *mockLogger) WithTraceID(traceID string) logger.Logger { return m }

func TestHashidsEncoder_WithLogger(t *testing.T) {
	logger := &mockLogger{}
	enc, err := NewHashidsEncoder(Config{Secret: "secret", MinLength: 10}, logger)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	t.Run("Encode with invalid UUID logs error", func(t *testing.T) {
		_, err := enc.Encode("not-a-uuid")
		if err == nil {
			t.Error("expected error for invalid UUID")
		}
	})

	t.Run("Decode with empty string logs error", func(t *testing.T) {
		_, err := enc.Decode("")
		if err == nil {
			t.Error("expected error for empty encoded")
		}
	})

	t.Run("Decode with invalid encoded logs error", func(t *testing.T) {
		_, err := enc.Decode("invalid-encoded-value!!!")
		if err == nil {
			t.Error("expected error for invalid encoded")
		}
	})

	t.Run("MustEncode with invalid UUID logs error", func(t *testing.T) {
		result := enc.MustEncode("not-a-uuid")
		if result != "" {
			t.Errorf("expected empty string, got %q", result)
		}
	})
}
