package licenseplate

import (
	"testing"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

func TestLicensePlate_ToDomain(t *testing.T) {
	lp := &LicensePlate{
		ID:              "lp1",
		LicensePlate:    "HK-1234",
		AircraftModelID: "am1",
		AirlineID:       "a1",
		ModelName:       "Boeing 737",
		AirlineName:     "Test Air",
	}

	result := lp.ToDomain()
	if result.ID != "lp1" {
		t.Errorf("expected ID 'lp1', got %q", result.ID)
	}
	if result.LicensePlate != "HK-1234" {
		t.Errorf("expected LicensePlate 'HK-1234', got %q", result.LicensePlate)
	}
	if result.ModelName != "Boeing 737" {
		t.Errorf("expected ModelName 'Boeing 737', got %q", result.ModelName)
	}
}

func TestLicensePlate_FromDomain(t *testing.T) {
	dm := &domain.LicensePlate{
		ID:              "lp2",
		LicensePlate:    "HK-5678",
		AircraftModelID: "am2",
		AirlineID:       "a2",
		ModelName:       "Airbus A320",
		AirlineName:     "Domain Air",
	}

	result := FromDomain(dm)
	if result.ID != "lp2" {
		t.Errorf("expected ID 'lp2', got %q", result.ID)
	}
	if result.LicensePlate != "HK-5678" {
		t.Errorf("expected LicensePlate 'HK-5678', got %q", result.LicensePlate)
	}
}

func TestLicensePlate_RoundTrip(t *testing.T) {
	original := &domain.LicensePlate{
		ID: "lp3", LicensePlate: "HK-9999", AircraftModelID: "am3",
		AirlineID: "a3", ModelName: "Embraer E190", AirlineName: "Round Air",
	}
	restored := FromDomain(original).ToDomain()
	if restored.ID != original.ID {
		t.Errorf("ID mismatch: expected %s, got %s", original.ID, restored.ID)
	}
	if restored.LicensePlate != original.LicensePlate {
		t.Errorf("LicensePlate mismatch: expected %s, got %s", original.LicensePlate, restored.LicensePlate)
	}
}
