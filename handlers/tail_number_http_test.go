package handlers

import (
	"testing"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
)

func TestCreateTailNumberRequest_Sanitize(t *testing.T) {
	t.Run("trims all fields", func(t *testing.T) {
		req := &CreateTailNumberRequest{
			TailNumber:    "  HK-5432  ",
			AircraftModelID: "  model-123  ",
			AirlineID:       "  airline-456  ",
		}
		req.Sanitize()

		if req.TailNumber != "HK-5432" {
			t.Errorf("expected 'HK-5432', got %q", req.TailNumber)
		}
		if req.AircraftModelID != "model-123" {
			t.Errorf("expected 'model-123', got %q", req.AircraftModelID)
		}
		if req.AirlineID != "airline-456" {
			t.Errorf("expected 'airline-456', got %q", req.AirlineID)
		}
	})
}

func TestUpdateTailNumberRequest_Sanitize(t *testing.T) {
	t.Run("trims all fields", func(t *testing.T) {
		req := &UpdateTailNumberRequest{
			TailNumber:    "  CC-BFA  ",
			AircraftModelID: "  model-789  ",
			AirlineID:       "  airline-012  ",
		}
		req.Sanitize()

		if req.TailNumber != "CC-BFA" {
			t.Errorf("expected 'CC-BFA', got %q", req.TailNumber)
		}
		if req.AircraftModelID != "model-789" {
			t.Errorf("expected 'model-789', got %q", req.AircraftModelID)
		}
		if req.AirlineID != "airline-012" {
			t.Errorf("expected 'airline-012', got %q", req.AirlineID)
		}
	})
}

func TestCreateTailNumberRequest_ToDomain(t *testing.T) {
	t.Run("converts valid request to domain", func(t *testing.T) {
		req := &CreateTailNumberRequest{
			TailNumber:    "HK-5432",
			AircraftModelID: "model-uuid",
			AirlineID:       "airline-uuid",
		}

		result := req.ToDomain()

		if result.ID == "" {
			t.Error("expected non-empty ID (auto-generated UUID)")
		}
		if result.TailNumber != "HK-5432" {
			t.Errorf("expected TailNumber 'HK-5432', got %q", result.TailNumber)
		}
		if result.AircraftModelID != "model-uuid" {
			t.Errorf("expected AircraftModelID 'model-uuid', got %q", result.AircraftModelID)
		}
		if result.AirlineID != "airline-uuid" {
			t.Errorf("expected AirlineID 'airline-uuid', got %q", result.AirlineID)
		}
	})
}

func TestUpdateTailNumberRequest_ToDomain(t *testing.T) {
	t.Run("converts valid request to domain", func(t *testing.T) {
		req := &UpdateTailNumberRequest{
			TailNumber:    "CC-BFA",
			AircraftModelID: "model-uuid",
			AirlineID:       "airline-uuid",
		}

		result := req.ToDomain("existing-uuid")

		if result.ID != "existing-uuid" {
			t.Errorf("expected ID 'existing-uuid', got %q", result.ID)
		}
		if result.TailNumber != "CC-BFA" {
			t.Errorf("expected TailNumber 'CC-BFA', got %q", result.TailNumber)
		}
	})
}

func TestFromDomainTailNumber(t *testing.T) {
	t.Run("converts domain to response with encoded IDs", func(t *testing.T) {
		registration := &domain.TailNumber{
			ID:              "raw-uuid",
			TailNumber:    "HK-5432",
			AircraftModelID: "raw-model-uuid",
			AirlineID:       "raw-airline-uuid",
			ModelName:       "Boeing 737",
			AirlineName:     "Avianca",
		}

		result := FromDomainTailNumber(registration, "encoded-id", "encoded-model-id", "encoded-airline-id")

		if result.ID != "encoded-id" {
			t.Errorf("expected 'encoded-id', got %q", result.ID)
		}
		if result.TailNumber != "HK-5432" {
			t.Errorf("expected 'HK-5432', got %q", result.TailNumber)
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

func TestToTailNumberListResponse(t *testing.T) {
	t.Run("converts list to response", func(t *testing.T) {
		registrations := []domain.TailNumber{
			{
				ID:              "uuid-1",
				TailNumber:    "HK-5432",
				AircraftModelID: "model-1",
				AirlineID:       "airline-1",
				ModelName:       "Boeing 737",
				AirlineName:     "Avianca",
			},
			{
				ID:              "uuid-2",
				TailNumber:    "CC-BFA",
				AircraftModelID: "model-2",
				AirlineID:       "airline-2",
				ModelName:       "Airbus A320",
				AirlineName:     "LATAM",
			},
		}

		encodeFunc := func(uuid string) (string, error) {
			return "enc-" + uuid, nil
		}

		result := ToTailNumberListResponse(registrations, encodeFunc, "http://localhost:8080")

		if len(result.Registrations) != 2 {
			t.Fatalf("expected 2 items, got %d", len(result.Registrations))
		}
		if result.Registrations[0].TailNumber != "HK-5432" {
			t.Errorf("expected 'HK-5432', got %q", result.Registrations[0].TailNumber)
		}
		if result.Registrations[1].TailNumber != "CC-BFA" {
			t.Errorf("expected 'CC-BFA', got %q", result.Registrations[1].TailNumber)
		}
	})

	t.Run("empty list returns empty response", func(t *testing.T) {
		encodeFunc := func(uuid string) (string, error) {
			return "enc-" + uuid, nil
		}

		result := ToTailNumberListResponse([]domain.TailNumber{}, encodeFunc, "http://localhost:8080")

		if len(result.Registrations) != 0 {
			t.Errorf("expected 0 items, got %d", len(result.Registrations))
		}
	})
}
