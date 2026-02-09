package domain

import "testing"

func TestAircraftFamily_ToLogger(t *testing.T) {
	af := &AircraftFamily{
		Family:       "A320",
		Manufacturer: "Airbus",
	}

	result := af.ToLogger()

	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	if result[0] != "family:A320" {
		t.Errorf("expected 'family:A320', got %q", result[0])
	}
	if result[1] != "manufacturer:Airbus" {
		t.Errorf("expected 'manufacturer:Airbus', got %q", result[1])
	}
}
