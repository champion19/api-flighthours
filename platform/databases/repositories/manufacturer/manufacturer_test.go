package manufacturer

import (
	"testing"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

func TestManufacturer_ToDomain(t *testing.T) {
	m := &Manufacturer{ID: "mfr-1", Name: "Boeing"}
	result := m.ToDomain()

	if result.ID != "mfr-1" {
		t.Errorf("expected ID 'mfr-1', got %q", result.ID)
	}
	if result.Name != "Boeing" {
		t.Errorf("expected Name 'Boeing', got %q", result.Name)
	}
}

func TestFromDomain(t *testing.T) {
	dm := &domain.Manufacturer{ID: "mfr-2", Name: "Airbus"}
	result := FromDomain(dm)

	if result.ID != "mfr-2" {
		t.Errorf("expected ID 'mfr-2', got %q", result.ID)
	}
	if result.Name != "Airbus" {
		t.Errorf("expected Name 'Airbus', got %q", result.Name)
	}
}

func TestManufacturer_RoundTrip(t *testing.T) {
	original := &domain.Manufacturer{ID: "mfr-3", Name: "Embraer"}
	restored := FromDomain(original).ToDomain()

	if restored.ID != original.ID {
		t.Errorf("ID mismatch: expected %s, got %s", original.ID, restored.ID)
	}
	if restored.Name != original.Name {
		t.Errorf("Name mismatch: expected %s, got %s", original.Name, restored.Name)
	}
}
