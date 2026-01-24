package domain

import (
	"testing"
)

func TestEmployee_SetID(t *testing.T) {
	emp := &Employee{}

	if emp.ID != "" {
		t.Fatal("expected ID to be empty initially")
	}

	emp.SetID()

	if emp.ID == "" {
		t.Fatal("expected ID to be set after SetID()")
	}

	// Verify UUID format (basic check)
	if len(emp.ID) < 30 {
		t.Fatalf("expected UUID format, got: %s", emp.ID)
	}
}

func TestEmployee_ToLogger(t *testing.T) {
	emp := &Employee{
		ID:    "emp-123",
		Email: "test@example.com",
		Role:  "admin",
	}

	result := emp.ToLogger()

	if len(result) != 3 {
		t.Fatalf("expected 3 elements, got %d", len(result))
	}

	expectedID := "id:emp-123"
	if result[0] != expectedID {
		t.Fatalf("expected %s, got %s", expectedID, result[0])
	}

	expectedEmail := "email:test@example.com"
	if result[1] != expectedEmail {
		t.Fatalf("expected %s, got %s", expectedEmail, result[1])
	}

	expectedRole := "role:admin"
	if result[2] != expectedRole {
		t.Fatalf("expected %s, got %s", expectedRole, result[2])
	}
}

func TestEmployee_ToLogger_EmptyFields(t *testing.T) {
	emp := &Employee{}

	result := emp.ToLogger()

	if len(result) != 3 {
		t.Fatalf("expected 3 elements, got %d", len(result))
	}

	if result[0] != "id:" {
		t.Fatalf("expected 'id:', got %s", result[0])
	}
	if result[1] != "email:" {
		t.Fatalf("expected 'email:', got %s", result[1])
	}
	if result[2] != "role:" {
		t.Fatalf("expected 'role:', got %s", result[2])
	}
}
