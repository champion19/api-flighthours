package engine

import (
	"testing"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

func TestEngine_ToDomain(t *testing.T) {
	e := &Engine{ID: "e1", Name: "Turbofan"}
	result := e.ToDomain()

	if result.ID != "e1" {
		t.Errorf("expected ID 'e1', got %q", result.ID)
	}
	if result.Name != "Turbofan" {
		t.Errorf("expected Name 'Turbofan', got %q", result.Name)
	}
}

func TestEngine_FromDomain(t *testing.T) {
	de := &domain.Engine{ID: "e2", Name: "Turboprop"}
	result := FromDomain(de)

	if result.ID != "e2" {
		t.Errorf("expected ID 'e2', got %q", result.ID)
	}
	if result.Name != "Turboprop" {
		t.Errorf("expected Name 'Turboprop', got %q", result.Name)
	}
}

func TestEngine_RoundTrip(t *testing.T) {
	original := &domain.Engine{ID: "e3", Name: "Piston"}
	restored := FromDomain(original).ToDomain()

	if restored.ID != original.ID {
		t.Errorf("ID mismatch: expected %s, got %s", original.ID, restored.ID)
	}
	if restored.Name != original.Name {
		t.Errorf("Name mismatch: expected %s, got %s", original.Name, restored.Name)
	}
}
