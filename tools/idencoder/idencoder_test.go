package idencoder

import (
	"testing"

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
