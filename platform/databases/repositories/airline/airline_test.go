package airline

import (
	"testing"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

func TestAirline_ToDomain(t *testing.T) {
	t.Run("converts database entity to domain model", func(t *testing.T) {
		dbAirline := &Airline{
			ID:          "airline-uuid-123",
			AirlineName: "Test Airlines",
			AirlineCode: "TST",
			Status:      "active",
		}

		result := dbAirline.ToDomain()
		if result == nil {
			t.Errorf("expected non-nil result, domain airline is nil")
		}

		if result.ID != dbAirline.ID {
			t.Errorf("expected ID %s, got %s", dbAirline.ID, result.ID)
		}
		if result.AirlineName != dbAirline.AirlineName {
			t.Errorf("expected AirlineName %s, got %s", dbAirline.AirlineName, result.AirlineName)
		}
		if result.AirlineCode != dbAirline.AirlineCode {
			t.Errorf("expected AirlineCode %s, got %s", dbAirline.AirlineCode, result.AirlineCode)
		}
		if result.Status != dbAirline.Status {
			t.Errorf("expected Status %s, got %s", dbAirline.Status, result.Status)
		}
	})

	t.Run("handles empty fields", func(t *testing.T) {
		dbAirline := &Airline{}

		result := dbAirline.ToDomain()

		if result == nil {
			t.Errorf("expected non-nil result, domain airline is nil")
		}

		if result.ID != "" {
			t.Errorf("expected empty ID, got %s", result.ID)
		}

	})
}

func TestFromDomain(t *testing.T) {
	t.Run("converts domain model to database entity", func(t *testing.T) {
		domainAirline := &domain.Airline{
			ID:          "airline-uuid-456",
			AirlineName: "Domain Airlines",
			AirlineCode: "DOM",
			Status:      "inactive",
		}

		result := FromDomain(domainAirline)

		if result.ID != domainAirline.ID {
			t.Errorf("expected ID %s, got %s", domainAirline.ID, result.ID)
		}
		if result.AirlineName != domainAirline.AirlineName {
			t.Errorf("expected AirlineName %s, got %s", domainAirline.AirlineName, result.AirlineName)
		}
		if result.AirlineCode != domainAirline.AirlineCode {
			t.Errorf("expected AirlineCode %s, got %s", domainAirline.AirlineCode, result.AirlineCode)
		}
		if result.Status != domainAirline.Status {
			t.Errorf("expected Status %s, got %s", domainAirline.Status, result.Status)
		}
	})

	t.Run("handles empty domain model", func(t *testing.T) {
		domainAirline := &domain.Airline{}

		result := FromDomain(domainAirline)

		if result.ID != "" {
			t.Errorf("expected empty ID, got %s", result.ID)
		}
	})
}

func TestAirline_RoundTrip(t *testing.T) {
	t.Run("domain -> db -> domain preserves data", func(t *testing.T) {
		original := &domain.Airline{
			ID:          "roundtrip-uuid-789",
			AirlineName: "Roundtrip Airlines",
			AirlineCode: "RND",
			Status:      "active",
		}

		dbEntity := FromDomain(original)
		restored := dbEntity.ToDomain()

		if restored.ID != original.ID {
			t.Errorf("ID mismatch: expected %s, got %s", original.ID, restored.ID)
		}
		if restored.AirlineName != original.AirlineName {
			t.Errorf("AirlineName mismatch: expected %s, got %s", original.AirlineName, restored.AirlineName)
		}
		if restored.AirlineCode != original.AirlineCode {
			t.Errorf("AirlineCode mismatch: expected %s, got %s", original.AirlineCode, restored.AirlineCode)
		}
		if restored.Status != original.Status {
			t.Errorf("Status mismatch: expected %v, got %v", original.Status, restored.Status)
		}
	})
}
