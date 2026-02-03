package handlers

import (
	"testing"
	"time"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
)

func TestAirlineEmployeeRequest_Sanitize(t *testing.T) {
	t.Run("trims all fields", func(t *testing.T) {
		req := &AirlineEmployeeRequest{
			AirlineID: "  airline-123  ",
			Bp:        "  BP001  ",
			StartDate: "  2024-01-01  ",
			EndDate:   "  2024-12-31  ",
		}
		req.Sanitize()

		if req.AirlineID != "airline-123" {
			t.Errorf("expected 'airline-123', got %q", req.AirlineID)
		}
		if req.Bp != "BP001" {
			t.Errorf("expected 'BP001', got %q", req.Bp)
		}
		if req.StartDate != "2024-01-01" {
			t.Errorf("expected '2024-01-01', got %q", req.StartDate)
		}
		if req.EndDate != "2024-12-31" {
			t.Errorf("expected '2024-12-31', got %q", req.EndDate)
		}
	})
}

func TestAirlineEmployeeRequest_ToDomain(t *testing.T) {
	t.Run("converts valid request to domain", func(t *testing.T) {
		req := &AirlineEmployeeRequest{
			AirlineID: "airline-123",
			Bp:        "BP001",
			StartDate: "2024-01-15",
			EndDate:   "2024-12-31",
			Active:    true,
		}

		result, err := req.ToDomain()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.AirlineID != "airline-123" {
			t.Errorf("expected 'airline-123', got %q", result.AirlineID)
		}
		if result.Bp != "BP001" {
			t.Errorf("expected 'BP001', got %q", result.Bp)
		}
		if !result.Active {
			t.Error("expected Active to be true")
		}
	})

	t.Run("returns error for invalid start date", func(t *testing.T) {
		req := &AirlineEmployeeRequest{
			AirlineID: "airline-123",
			StartDate: "invalid-date",
		}

		_, err := req.ToDomain()
		if err == nil {
			t.Error("expected error for invalid start date")
		}
	})

	t.Run("returns error for invalid end date", func(t *testing.T) {
		req := &AirlineEmployeeRequest{
			AirlineID: "airline-123",
			StartDate: "2024-01-01",
			EndDate:   "not-a-date",
		}

		_, err := req.ToDomain()
		if err == nil {
			t.Error("expected error for invalid end date")
		}
	})

	t.Run("returns error when start date is after end date", func(t *testing.T) {
		req := &AirlineEmployeeRequest{
			AirlineID: "airline-123",
			StartDate: "2024-12-31",
			EndDate:   "2024-01-01",
		}

		_, err := req.ToDomain()
		if err == nil {
			t.Error("expected error for start date after end date")
		}
	})

	t.Run("handles empty end date", func(t *testing.T) {
		req := &AirlineEmployeeRequest{
			AirlineID: "airline-123",
			StartDate: "2024-01-01",
			EndDate:   "",
		}

		result, err := req.ToDomain()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.EndDate.IsZero() {
			t.Error("expected zero EndDate for empty string")
		}
	})
}

func TestFromDomainAirlineEmployee(t *testing.T) {
	t.Run("converts domain to response", func(t *testing.T) {
		startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)

		employee := &domain.AirlineEmployee{
			ID:        "raw-id",
			AirlineID: "raw-airline-id",
			Bp:        "BP001",
			StartDate: startDate,
			EndDate:   endDate,
			Active:    true,
		}

		result := FromDomainAirlineEmployee(employee, "encoded-id", "encoded-airline-id")

		if result.ID != "encoded-id" {
			t.Errorf("expected 'encoded-id', got %q", result.ID)
		}
		if result.AirlineID != "encoded-airline-id" {
			t.Errorf("expected 'encoded-airline-id', got %q", result.AirlineID)
		}
		if result.Bp != "BP001" {
			t.Errorf("expected 'BP001', got %q", result.Bp)
		}
		if !result.Active {
			t.Error("expected Active to be true")
		}
	})
}

func TestAddEmployeeAirlineRequest_Sanitize(t *testing.T) {
	t.Run("trims all fields", func(t *testing.T) {
		req := &AddEmployeeAirlineRequest{
			AirlineID: "  add-airline-123  ",
			Bp:        "  ADD001  ",
			StartDate: "  2024-06-01  ",
			EndDate:   "  2024-06-30  ",
		}
		req.Sanitize()

		if req.AirlineID != "add-airline-123" {
			t.Errorf("expected 'add-airline-123', got %q", req.AirlineID)
		}
		if req.Bp != "ADD001" {
			t.Errorf("expected 'ADD001', got %q", req.Bp)
		}
	})
}

func TestUpdateEmployeeAirlineRequest_Sanitize(t *testing.T) {
	t.Run("trims all fields", func(t *testing.T) {
		req := &UpdateEmployeeAirlineRequest{
			AirlineID: "  update-airline-456  ",
			Bp:        "  UPD001  ",
			StartDate: "  2024-07-01  ",
			EndDate:   "  2024-07-31  ",
		}
		req.Sanitize()

		if req.AirlineID != "update-airline-456" {
			t.Errorf("expected 'update-airline-456', got %q", req.AirlineID)
		}
		if req.Bp != "UPD001" {
			t.Errorf("expected 'UPD001', got %q", req.Bp)
		}
	})
}
