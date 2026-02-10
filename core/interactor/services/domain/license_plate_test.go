package domain

import "testing"

func TestLicense_Plate_ToLogger(t *testing.T) {
	t.Run("returns expected log fields", func(t *testing.T) {
		ar := LicensePlate{
			ID:              "uuid-123",
			LicensePlate:    "HK-5432",
			AircraftModelID: "model-uuid",
			AirlineID:       "airline-uuid",
			ModelName:       "Boeing 737",
			AirlineName:     "Avianca",
		}

		fields := ar.ToLogger()
		if len(fields) == 0 {
			t.Fatal("expected non-empty log fields")
		}
	})

	t.Run("handles empty fields", func(t *testing.T) {
		ar := LicensePlate{}
		fields := ar.ToLogger()
		if fields == nil {
			t.Fatal("expected non-nil log fields even with empty struct")
		}
	})
}
