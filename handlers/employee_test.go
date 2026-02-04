package handlers

import (
	"testing"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
)

func TestEmployeeRequest_Sanitize(t *testing.T) {
	t.Run("trims all fields", func(t *testing.T) {
		req := &EmployeeRequest{
			Name:                 "  John Doe  ",
			Email:                "  john@example.com  ",
			Password:             "password123",
			IdentificationNumber: "  ABC123  ",
			Role:                 "  pilot  ",
		}
		req.Sanitize()

		if req.Name != "John Doe" {
			t.Errorf("expected 'John Doe', got %q", req.Name)
		}
		if req.Email != "john@example.com" {
			t.Errorf("expected 'john@example.com', got %q", req.Email)
		}
		if req.IdentificationNumber != "ABC123" {
			t.Errorf("expected 'ABC123', got %q", req.IdentificationNumber)
		}
		if req.Role != "pilot" {
			t.Errorf("expected 'pilot', got %q", req.Role)
		}
	})
}

func TestEmployeeRequest_ToDomain(t *testing.T) {
	t.Run("converts to domain employee", func(t *testing.T) {
		req := EmployeeRequest{
			Name:                 "John Doe",
			Email:                "john@example.com",
			Password:             "secret123",
			IdentificationNumber: "ID123",
			Role:                 "pilot",
		}

		result := req.ToDomain()

		if result.Name != "John Doe" {
			t.Errorf("expected 'John Doe', got %q", result.Name)
		}
		if result.Email != "john@example.com" {
			t.Errorf("expected 'john@example.com', got %q", result.Email)
		}
		if result.Password != "secret123" {
			t.Errorf("expected 'secret123', got %q", result.Password)
		}
		if result.IdentificationNumber != "ID123" {
			t.Errorf("expected 'ID123', got %q", result.IdentificationNumber)
		}
		if result.Role != "pilot" {
			t.Errorf("expected 'pilot', got %q", result.Role)
		}
	})
}

func TestFromDomain(t *testing.T) {
	t.Run("converts domain employee to response", func(t *testing.T) {
		employee := &domain.Employee{
			ID:                   "raw-uuid",
			Name:                 "Jane Doe",
			Email:                "jane@example.com",
			IdentificationNumber: "XYZ789",
			Role:                 "admin",
		}

		result := FromDomain(employee, "encoded-id-123")

		if result.ID != "encoded-id-123" {
			t.Errorf("expected 'encoded-id-123', got %q", result.ID)
		}
		if result.Name != "Jane Doe" {
			t.Errorf("expected 'Jane Doe', got %q", result.Name)
		}
		if result.Email != "jane@example.com" {
			t.Errorf("expected 'jane@example.com', got %q", result.Email)
		}
	})
}

func TestLoginRequest_Sanitize(t *testing.T) {
	t.Run("trims email", func(t *testing.T) {
		req := &LoginRequest{Email: "  user@test.com  ", Password: "pass"}
		req.Sanitize()
		if req.Email != "user@test.com" {
			t.Errorf("expected 'user@test.com', got %q", req.Email)
		}
	})
}

func TestResendVerificationEmailRequest_Sanitize(t *testing.T) {
	t.Run("trims email", func(t *testing.T) {
		req := &ResendVerificationEmailRequest{Email: "  test@email.com  "}
		req.Sanitize()
		if req.Email != "test@email.com" {
			t.Errorf("expected 'test@email.com', got %q", req.Email)
		}
	})
}

func TestPasswordResetRequest_Sanitize(t *testing.T) {
	t.Run("trims email", func(t *testing.T) {
		req := &PasswordResetRequest{Email: "  reset@test.com  "}
		req.Sanitize()
		if req.Email != "reset@test.com" {
			t.Errorf("expected 'reset@test.com', got %q", req.Email)
		}
	})
}

func TestChangePasswordRequest_Sanitize(t *testing.T) {
	t.Run("trims email", func(t *testing.T) {
		req := &ChangePasswordRequest{
			Email:           "  change@test.com  ",
			CurrentPassword: "old",
			NewPassword:     "new",
			ConfirmPassword: "new",
		}
		req.Sanitize()
		if req.Email != "change@test.com" {
			t.Errorf("expected 'change@test.com', got %q", req.Email)
		}
	})
}

func TestUpdateEmployeeRequest_Sanitize(t *testing.T) {
	t.Run("trims name and ID number", func(t *testing.T) {
		req := &UpdateEmployeeRequest{
			Name:                 "  Updated Name  ",
			IdentificationNumber: "  NEW123  ",
		}
		req.Sanitize()
		if req.Name != "Updated Name" {
			t.Errorf("expected 'Updated Name', got %q", req.Name)
		}
		if req.IdentificationNumber != "NEW123" {
			t.Errorf("expected 'NEW123', got %q", req.IdentificationNumber)
		}
	})
}

func TestUpdateEmployeeRequest_ToUpdateData(t *testing.T) {
	t.Run("merges with existing employee preserving immutable fields", func(t *testing.T) {
		existing := &domain.Employee{
			ID:                   "existing-id",
			Name:                 "Old Name",
			Email:                "old@email.com",
			Password:             "hashed-password",
			IdentificationNumber: "OLD123",
			Role:                 "pilot",
			KeycloakUserID:       "kc-123",
		}

		req := UpdateEmployeeRequest{
			Name:                 "New Name",
			IdentificationNumber: "NEW456",
		}

		result := req.ToUpdateData(existing)

		// Mutable fields should be updated
		if result.Name != "New Name" {
			t.Errorf("expected 'New Name', got %q", result.Name)
		}
		if result.IdentificationNumber != "NEW456" {
			t.Errorf("expected 'NEW456', got %q", result.IdentificationNumber)
		}

		// Immutable fields should be preserved
		if result.ID != "existing-id" {
			t.Errorf("expected 'existing-id', got %q", result.ID)
		}
		if result.Email != "old@email.com" {
			t.Errorf("expected 'old@email.com', got %q", result.Email)
		}
		if result.Role != "pilot" {
			t.Errorf("expected 'pilot', got %q", result.Role)
		}
		if result.KeycloakUserID != "kc-123" {
			t.Errorf("expected 'kc-123', got %q", result.KeycloakUserID)
		}
	})
}
