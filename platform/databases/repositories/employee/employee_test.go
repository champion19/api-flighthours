package employee

import (
	"testing"
	"time"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

func TestEmployee_ToDomain(t *testing.T) {
	t.Run("converts Employee to domain with nil KeycloakUserID", func(t *testing.T) {
		emp := Employee{
			ID:                   "emp-123",
			Name:                 "John Doe",
			Email:                "john@example.com",
			IdentificationNumber: "ID123456",
			Role:                 "pilot",
			KeycloakUserID:       nil,
		}

		result := emp.ToDomain()

		if result.ID != "emp-123" {
			t.Errorf("expected ID 'emp-123', got %q", result.ID)
		}
		if result.Name != "John Doe" {
			t.Errorf("expected Name 'John Doe', got %q", result.Name)
		}
		if result.Email != "john@example.com" {
			t.Errorf("expected Email 'john@example.com', got %q", result.Email)
		}
		if result.IdentificationNumber != "ID123456" {
			t.Errorf("expected IdentificationNumber 'ID123456', got %q", result.IdentificationNumber)
		}
		if result.Role != "pilot" {
			t.Errorf("expected Role 'pilot', got %q", result.Role)
		}
		if result.KeycloakUserID != "" {
			t.Errorf("expected empty KeycloakUserID, got %q", result.KeycloakUserID)
		}
	})

	t.Run("converts Employee to domain with KeycloakUserID", func(t *testing.T) {
		kcID := "kc-user-456"
		emp := Employee{
			ID:             "emp-123",
			Name:           "Jane Doe",
			Email:          "jane@example.com",
			Role:           "admin",
			KeycloakUserID: &kcID,
		}

		result := emp.ToDomain()

		if result.KeycloakUserID != "kc-user-456" {
			t.Errorf("expected KeycloakUserID 'kc-user-456', got %q", result.KeycloakUserID)
		}
	})
}

func TestFromDomain(t *testing.T) {
	t.Run("converts domain.Employee to Employee", func(t *testing.T) {
		domainEmp := domain.Employee{
			ID:                   "emp-123",
			Name:                 "John Doe",
			Email:                "john@example.com",
			IdentificationNumber: "ID123456",
			Role:                 "pilot",
			KeycloakUserID:       "kc-123",
		}

		result := FromDomain(domainEmp)

		if result.ID != "emp-123" {
			t.Errorf("expected ID 'emp-123', got %q", result.ID)
		}
		if result.Name != "John Doe" {
			t.Errorf("expected Name 'John Doe', got %q", result.Name)
		}
		if result.Email != "john@example.com" {
			t.Errorf("expected Email 'john@example.com', got %q", result.Email)
		}
		if !result.Active {
			// Active should be false by default
		}
		if result.Airline != nil {
			t.Error("expected nil Airline")
		}
		if result.KeycloakUserID == nil || *result.KeycloakUserID != "kc-123" {
			t.Error("expected KeycloakUserID to be set")
		}
	})

	t.Run("handles empty KeycloakUserID", func(t *testing.T) {
		domainEmp := domain.Employee{
			ID:             "emp-123",
			KeycloakUserID: "",
		}

		result := FromDomain(domainEmp)

		if result.KeycloakUserID != nil {
			t.Error("expected nil KeycloakUserID for empty string")
		}
	})
}

func TestStringPtrOrNil(t *testing.T) {
	t.Run("returns nil for empty string", func(t *testing.T) {
		result := stringPtrOrNil("")

		if result != nil {
			t.Error("expected nil for empty string")
		}
	})

	t.Run("returns pointer for non-empty string", func(t *testing.T) {
		result := stringPtrOrNil("test")

		if result == nil {
			t.Fatal("expected non-nil result")
		}
		if *result != "test" {
			t.Errorf("expected 'test', got %q", *result)
		}
	})
}

func TestEmployeeStruct(t *testing.T) {
	t.Run("creates Employee with all fields", func(t *testing.T) {
		airline := "Airline1"
		bp := "BP123"
		startDate := time.Now()
		endDate := startDate.Add(24 * time.Hour)
		kcID := "kc-123"

		emp := Employee{
			ID:                   "emp-123",
			Name:                 "Test User",
			Airline:              &airline,
			Email:                "test@example.com",
			IdentificationNumber: "ID123",
			Bp:                   &bp,
			StartDate:            &startDate,
			EndDate:              &endDate,
			Active:               true,
			Role:                 "pilot",
			KeycloakUserID:       &kcID,
		}

		if emp.ID != "emp-123" {
			t.Errorf("expected ID 'emp-123', got %q", emp.ID)
		}
		if *emp.Airline != "Airline1" {
			t.Errorf("expected Airline 'Airline1', got %q", *emp.Airline)
		}
		if !emp.Active {
			t.Error("expected Active to be true")
		}
	})
}
