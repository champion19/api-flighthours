package dto

import (
	"testing"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/platform/logger"
)

func TestFromDomainToDTO(t *testing.T) {
	t.Run("converts domain employee to DTO", func(t *testing.T) {
		employee := &domain.Employee{
			ID:    "emp-123",
			Name:  "John Doe",
			Email: "john@example.com",
			Role:  "pilot",
		}

		result := FromDomainToDTO(employee)

		if result == nil {
			t.Fatal("expected non-nil result")
		}
		if result.Employee.ID != "emp-123" {
			t.Errorf("expected ID 'emp-123', got %q", result.Employee.ID)
		}
		if result.Employee.Name != "John Doe" {
			t.Errorf("expected Name 'John Doe', got %q", result.Employee.Name)
		}
		if result.Message != logger.DtoMsgEmployeeLocated {
			t.Errorf("expected message %q, got %q", logger.DtoMsgEmployeeLocated, result.Message)
		}
	})
}

func TestRegisterEmployee(t *testing.T) {
	t.Run("creates RegisterEmployee struct", func(t *testing.T) {
		emp := domain.Employee{
			ID:   "test-id",
			Name: "Test Name",
		}

		reg := RegisterEmployee{
			Employee: emp,
			Message:  "test message",
		}

		if reg.Employee.ID != "test-id" {
			t.Errorf("expected ID 'test-id', got %q", reg.Employee.ID)
		}
		if reg.Message != "test message" {
			t.Errorf("expected 'test message', got %q", reg.Message)
		}
	})
}

func TestUserSyncStatus(t *testing.T) {
	t.Run("creates UserSyncStatus struct", func(t *testing.T) {
		sync := UserSyncStatus{
			EmployeeID:     "emp-123",
			KeycloakUserID: "kc-456",
			IsSynced:       true,
			LastSyncAt:     "2024-01-01",
		}

		if sync.EmployeeID != "emp-123" {
			t.Errorf("expected EmployeeID 'emp-123', got %q", sync.EmployeeID)
		}
		if !sync.IsSynced {
			t.Error("expected IsSynced to be true")
		}
	})
}

func TestTokenResponse(t *testing.T) {
	t.Run("creates TokenResponse struct", func(t *testing.T) {
		token := TokenResponse{
			AccessToken:  "access-token",
			RefreshToken: "refresh-token",
			ExpiresIn:    3600,
			TokenType:    "Bearer",
		}

		if token.AccessToken != "access-token" {
			t.Errorf("expected 'access-token', got %q", token.AccessToken)
		}
		if token.ExpiresIn != 3600 {
			t.Errorf("expected 3600, got %d", token.ExpiresIn)
		}
	})
}

func TestUpdateEmployee(t *testing.T) {
	t.Run("creates UpdateEmployee struct", func(t *testing.T) {
		update := UpdateEmployee{
			ID:      "emp-123",
			Updated: true,
			Message: "update successful",
		}

		if update.ID != "emp-123" {
			t.Errorf("expected ID 'emp-123', got %q", update.ID)
		}
		if !update.Updated {
			t.Error("expected Updated to be true")
		}
	})
}
