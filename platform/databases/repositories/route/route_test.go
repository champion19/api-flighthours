package route

import (
	"testing"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

func TestRoute_ToDomain(t *testing.T) {
	t.Run("converts all fields", func(t *testing.T) {
		r := &Route{
			ID:                     "route-123",
			OriginAirportID:        "airport-001",
			OriginIataCode:         "BOG",
			OriginAirportName:      "El Dorado",
			DestinationAirportID:   "airport-002",
			DestinationIataCode:    "JFK",
			DestinationAirportName: "John F Kennedy",
			OriginCountry:          "Colombia",
			DestinationCountry:     "USA",
			AirportType:            "International",
			EstimatedFlightTime:    "5h 30m",
			RouteCode:              "BOG-JFK",
		}

		result := r.ToDomain()

		if result.ID != "route-123" {
			t.Errorf("expected ID 'route-123', got %q", result.ID)
		}
		if result.OriginAirportID != "airport-001" {
			t.Errorf("expected OriginAirportID 'airport-001', got %q", result.OriginAirportID)
		}
		if result.OriginIataCode != "BOG" {
			t.Errorf("expected OriginIataCode 'BOG', got %q", result.OriginIataCode)
		}
		if result.OriginAirportName != "El Dorado" {
			t.Errorf("expected OriginAirportName 'El Dorado', got %q", result.OriginAirportName)
		}
		if result.DestinationIataCode != "JFK" {
			t.Errorf("expected DestinationIataCode 'JFK', got %q", result.DestinationIataCode)
		}
		if result.OriginCountry != "Colombia" {
			t.Errorf("expected OriginCountry 'Colombia', got %q", result.OriginCountry)
		}
		if result.DestinationCountry != "USA" {
			t.Errorf("expected DestinationCountry 'USA', got %q", result.DestinationCountry)
		}
		if result.RouteCode != "BOG-JFK" {
			t.Errorf("expected RouteCode 'BOG-JFK', got %q", result.RouteCode)
		}
	})

	t.Run("handles empty fields", func(t *testing.T) {
		r := &Route{ID: "route-456"}

		result := r.ToDomain()

		if result.ID != "route-456" {
			t.Errorf("expected ID 'route-456', got %q", result.ID)
		}
		if result.OriginIataCode != "" {
			t.Error("expected empty OriginIataCode")
		}
	})
}

func TestFromDomain(t *testing.T) {
	t.Run("converts all fields", func(t *testing.T) {
		dm := &domain.Route{
			ID:                     "route-789",
			OriginAirportID:        "airport-003",
			OriginIataCode:         "MDE",
			OriginAirportName:      "Jose Maria Cordova",
			DestinationAirportID:   "airport-004",
			DestinationIataCode:    "MIA",
			DestinationAirportName: "Miami International",
			OriginCountry:          "Colombia",
			DestinationCountry:     "USA",
			AirportType:            "International",
			EstimatedFlightTime:    "4h 15m",
			RouteCode:              "MDE-MIA",
		}

		result := FromDomain(dm)

		if result.ID != "route-789" {
			t.Errorf("expected ID 'route-789', got %q", result.ID)
		}
		if result.OriginAirportID != "airport-003" {
			t.Errorf("expected OriginAirportID 'airport-003', got %q", result.OriginAirportID)
		}
		if result.OriginIataCode != "MDE" {
			t.Errorf("expected OriginIataCode 'MDE', got %q", result.OriginIataCode)
		}
		if result.DestinationIataCode != "MIA" {
			t.Errorf("expected DestinationIataCode 'MIA', got %q", result.DestinationIataCode)
		}
		if result.RouteCode != "MDE-MIA" {
			t.Errorf("expected RouteCode 'MDE-MIA', got %q", result.RouteCode)
		}
	})

	t.Run("roundtrip preserves data", func(t *testing.T) {
		original := &domain.Route{
			ID:                     "route-trip",
			OriginAirportID:        "origin-id",
			OriginIataCode:         "CLO",
			OriginAirportName:      "Alfonso Bonilla Aragon",
			DestinationAirportID:   "dest-id",
			DestinationIataCode:    "LAX",
			DestinationAirportName: "Los Angeles International",
			OriginCountry:          "Colombia",
			DestinationCountry:     "USA",
			AirportType:            "International",
			EstimatedFlightTime:    "7h 00m",
			RouteCode:              "CLO-LAX",
		}

		dbEntity := FromDomain(original)
		restored := dbEntity.ToDomain()

		if restored.ID != original.ID {
			t.Errorf("ID mismatch: expected %s, got %s", original.ID, restored.ID)
		}
		if restored.RouteCode != original.RouteCode {
			t.Errorf("RouteCode mismatch: expected %s, got %s", original.RouteCode, restored.RouteCode)
		}
		if restored.OriginCountry != original.OriginCountry {
			t.Errorf("OriginCountry mismatch: expected %s, got %s", original.OriginCountry, restored.OriginCountry)
		}
		if restored.EstimatedFlightTime != original.EstimatedFlightTime {
			t.Errorf("EstimatedFlightTime mismatch: expected %s, got %s", original.EstimatedFlightTime, restored.EstimatedFlightTime)
		}
	})
}
