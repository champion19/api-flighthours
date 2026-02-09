package handlers

import (
	"testing"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
)

func TestCreateLicensePlateRequest_Sanitize(t *testing.T) {
	t.Run("trims all fields", func(t *testing.T) {
		req := &CreateLicensePlateRequest{
			LicensePlate:    "  HK-5432  ",
			AircraftModelID: "  model-123  ",
			AirlineID:       "  airline-456  ",
		}
		req.Sanitize()

		if req.LicensePlate != "HK-5432" {
			t.Errorf("expected 'HK-5432', got %q", req.LicensePlate)
		}
		if req.AircraftModelID != "model-123" {
			t.Errorf("expected 'model-123', got %q", req.AircraftModelID)
		}
		if req.AirlineID != "airline-456" {
			t.Errorf("expected 'airline-456', got %q", req.AirlineID)
		}
	})
}

func TestUpdateLicensePlateRequest_Sanitize(t *testing.T) {
	t.Run("trims all fields", func(t *testing.T) {
		req := &UpdateLicensePlateRequest{
			LicensePlate:    "  CC-BFA  ",
			AircraftModelID: "  model-789  ",
			AirlineID:       "  airline-012  ",
		}
		req.Sanitize()

		if req.LicensePlate != "CC-BFA" {
			t.Errorf("expected 'CC-BFA', got %q", req.LicensePlate)
		}
		if req.AircraftModelID != "model-789" {
			t.Errorf("expected 'model-789', got %q", req.AircraftModelID)
		}
		if req.AirlineID != "airline-012" {
			t.Errorf("expected 'airline-012', got %q", req.AirlineID)
		}
	})
}

func TestCreateLicensePlateRequest_ToDomain(t *testing.T) {
	t.Run("converts valid request to domain", func(t *testing.T) {
		req := &CreateLicensePlateRequest{
			LicensePlate:    "HK-5432",
			AircraftModelID: "model-uuid",
			AirlineID:       "airline-uuid",
		}

		result := req.ToDomain()

		if result.ID == "" {
			t.Error("expected non-empty ID (auto-generated UUID)")
		}
		if result.LicensePlate != "HK-5432" {
			t.Errorf("expected LicensePlate 'HK-5432', got %q", result.LicensePlate)
		}
		if result.AircraftModelID != "model-uuid" {
			t.Errorf("expected AircraftModelID 'model-uuid', got %q", result.AircraftModelID)
		}
		if result.AirlineID != "airline-uuid" {
			t.Errorf("expected AirlineID 'airline-uuid', got %q", result.AirlineID)
		}
	})
}

func TestUpdateLicensePlateRequest_ToDomain(t *testing.T) {
	t.Run("converts valid request to domain", func(t *testing.T) {
		req := &UpdateLicensePlateRequest{
			LicensePlate:    "CC-BFA",
			AircraftModelID: "model-uuid",
			AirlineID:       "airline-uuid",
		}

		result := req.ToDomain("existing-uuid")

		if result.ID != "existing-uuid" {
			t.Errorf("expected ID 'existing-uuid', got %q", result.ID)
		}
		if result.LicensePlate != "CC-BFA" {
			t.Errorf("expected LicensePlate 'CC-BFA', got %q", result.LicensePlate)
		}
	})
}

func TestFromDomainLicensePlate(t *testing.T) {
	t.Run("converts domain to response with encoded IDs", func(t *testing.T) {
		registration := &domain.LicensePlate{
			ID:              "raw-uuid",
			LicensePlate:    "HK-5432",
			AircraftModelID: "raw-model-uuid",
			AirlineID:       "raw-airline-uuid",
			ModelName:       "Boeing 737",
			AirlineName:     "Avianca",
		}

		result := FromDomainLicensePlate(registration, "encoded-id", "encoded-model-id", "encoded-airline-id")

		if result.ID != "encoded-id" {
			t.Errorf("expected 'encoded-id', got %q", result.ID)
		}
		if result.LicensePlate != "HK-5432" {
			t.Errorf("expected 'HK-5432', got %q", result.LicensePlate)
		}
		if result.AircraftModelID != "encoded-model-id" {
			t.Errorf("expected 'encoded-model-id', got %q", result.AircraftModelID)
		}
		if result.AirlineID != "encoded-airline-id" {
			t.Errorf("expected 'encoded-airline-id', got %q", result.AirlineID)
		}
		if result.ModelName != "Boeing 737" {
			t.Errorf("expected 'Boeing 737', got %q", result.ModelName)
		}
		if result.AirlineName != "Avianca" {
			t.Errorf("expected 'Avianca', got %q", result.AirlineName)
		}
	})
}

func TestToLicensePlateListResponse(t *testing.T) {
	t.Run("converts list to response", func(t *testing.T) {
		registrations := []domain.LicensePlate{
			{
				ID:              "uuid-1",
				LicensePlate:    "HK-5432",
				AircraftModelID: "model-1",
				AirlineID:       "airline-1",
				ModelName:       "Boeing 737",
				AirlineName:     "Avianca",
			},
			{
				ID:              "uuid-2",
				LicensePlate:    "CC-BFA",
				AircraftModelID: "model-2",
				AirlineID:       "airline-2",
				ModelName:       "Airbus A320",
				AirlineName:     "LATAM",
			},
		}

		encodeFunc := func(uuid string) (string, error) {
			return "enc-" + uuid, nil
		}

		result := ToLicensePlateListResponse(registrations, encodeFunc, "http://localhost:8080")

		if len(result.Registrations) != 2 {
			t.Fatalf("expected 2 items, got %d", len(result.Registrations))
		}
		if result.Registrations[0].LicensePlate != "HK-5432" {
			t.Errorf("expected 'HK-5432', got %q", result.Registrations[0].LicensePlate)
		}
		if result.Registrations[1].LicensePlate != "CC-BFA" {
			t.Errorf("expected 'CC-BFA', got %q", result.Registrations[1].LicensePlate)
		}
	})

	t.Run("empty list returns empty response", func(t *testing.T) {
		encodeFunc := func(uuid string) (string, error) {
			return "enc-" + uuid, nil
		}

		result := ToLicensePlateListResponse([]domain.LicensePlate{}, encodeFunc, "http://localhost:8080")

		if len(result.Registrations) != 0 {
			t.Errorf("expected 0 items, got %d", len(result.Registrations))
		}
	})
}
