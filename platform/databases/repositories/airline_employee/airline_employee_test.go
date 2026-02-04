package airline_employee

import (
	"testing"
	"time"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

func TestAirlineEmployee_ToDomain(t *testing.T) {
	t.Run("converts with nil pointers", func(t *testing.T) {
		ae := &AirlineEmployee{
			ID:        "emp-123",
			AirlineID: nil,
			Bp:        nil,
			StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
			Active:    true,
		}

		result := ae.ToDomain()

		if result.ID != "emp-123" {
			t.Errorf("expected ID 'emp-123', got %q", result.ID)
		}
		if result.AirlineID != "" {
			t.Errorf("expected empty AirlineID, got %q", result.AirlineID)
		}
		if result.Bp != "" {
			t.Errorf("expected empty Bp, got %q", result.Bp)
		}
		if !result.Active {
			t.Error("expected Active to be true")
		}
	})

	t.Run("converts with non-nil pointers", func(t *testing.T) {
		airlineID := "airline-456"
		bp := "BP123"
		ae := &AirlineEmployee{
			ID:        "emp-789",
			AirlineID: &airlineID,
			Bp:        &bp,
			StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
			Active:    false,
		}

		result := ae.ToDomain()

		if result.AirlineID != "airline-456" {
			t.Errorf("expected AirlineID 'airline-456', got %q", result.AirlineID)
		}
		if result.Bp != "BP123" {
			t.Errorf("expected Bp 'BP123', got %q", result.Bp)
		}
		if result.Active {
			t.Error("expected Active to be false")
		}
	})

	t.Run("preserves dates", func(t *testing.T) {
		startDate := time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC)
		ae := &AirlineEmployee{
			ID:        "emp-123",
			StartDate: startDate,
			EndDate:   endDate,
		}

		result := ae.ToDomain()

		if !result.StartDate.Equal(startDate) {
			t.Error("StartDate should be preserved")
		}
		if !result.EndDate.Equal(endDate) {
			t.Error("EndDate should be preserved")
		}
	})
}

func TestFromDomain(t *testing.T) {
	t.Run("converts domain to repository struct", func(t *testing.T) {
		dm := &domain.AirlineEmployee{
			ID:        "emp-123",
			AirlineID: "airline-456",
			Bp:        "BP789",
			StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
			Active:    true,
		}

		result := FromDomain(dm)

		if result.ID != "emp-123" {
			t.Errorf("expected ID 'emp-123', got %q", result.ID)
		}
		if result.AirlineID == nil || *result.AirlineID != "airline-456" {
			t.Error("expected AirlineID to be set to 'airline-456'")
		}
		if result.Bp == nil || *result.Bp != "BP789" {
			t.Error("expected Bp to be set to 'BP789'")
		}
		if !result.Active {
			t.Error("expected Active to be true")
		}
	})

	t.Run("handles empty strings", func(t *testing.T) {
		dm := &domain.AirlineEmployee{
			ID:        "emp-123",
			AirlineID: "",
			Bp:        "",
		}

		result := FromDomain(dm)

		if result.AirlineID != nil {
			t.Error("expected nil AirlineID for empty string")
		}
		if result.Bp != nil {
			t.Error("expected nil Bp for empty string")
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
		result := stringPtrOrNil("test-value")

		if result == nil {
			t.Fatal("expected non-nil result")
		}
		if *result != "test-value" {
			t.Errorf("expected 'test-value', got %q", *result)
		}
	})
}
