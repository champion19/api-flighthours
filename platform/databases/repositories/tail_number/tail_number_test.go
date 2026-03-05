package tailnumber

import (
	"testing"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

func TestTailNumber_ToDomain(t *testing.T) {
	lp := &TailNumber{
		ID:              "lp1",
		TailNumber:    "HK-1234",
		AircraftModelID: "am1",
		AirlineID:       "a1",
		ModelName:       "Boeing 737",
		AirlineName:     "Test Air",
	}

	result := lp.ToDomain()
	if result.ID != "lp1" {
		t.Errorf("expected ID 'lp1', got %q", result.ID)
	}
	if result.TailNumber != "HK-1234" {
		t.Errorf("expected TailNumber 'HK-1234', got %q", result.TailNumber)
	}
	if result.ModelName != "Boeing 737" {
		t.Errorf("expected ModelName 'Boeing 737', got %q", result.ModelName)
	}
}

func TestTailNumber_FromDomain(t *testing.T) {
	dm := &domain.TailNumber{
		ID:              "lp2",
		TailNumber:    "HK-5678",
		AircraftModelID: "am2",
		AirlineID:       "a2",
		ModelName:       "Airbus A320",
		AirlineName:     "Domain Air",
	}

	result := FromDomain(dm)
	if result.ID != "lp2" {
		t.Errorf("expected ID 'lp2', got %q", result.ID)
	}
	if result.TailNumber != "HK-5678" {
		t.Errorf("expected TailNumber 'HK-5678', got %q", result.TailNumber)
	}
}

func TestTailNumber_RoundTrip(t *testing.T) {
	original := &domain.TailNumber{
		ID: "lp3", TailNumber: "HK-9999", AircraftModelID: "am3",
		AirlineID: "a3", ModelName: "Embraer E190", AirlineName: "Round Air",
	}
	restored := FromDomain(original).ToDomain()
	if restored.ID != original.ID {
		t.Errorf("ID mismatch: expected %s, got %s", original.ID, restored.ID)
	}
	if restored.TailNumber != original.TailNumber {
		t.Errorf("TailNumber mismatch: expected %s, got %s", original.TailNumber, restored.TailNumber)
	}
}
