package airline_route

import (
	"testing"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

func TestAirlineRoute_ToDomain(t *testing.T) {
	t.Run("converts all fields", func(t *testing.T) {
		ar := &AirlineRoute{
			ID:                     "ar-123",
			RouteID:                "route-456",
			AirlineID:              "airline-789",
			Status:                 domain.AirlineRouteStatusActive,
			AirlineCode:            "TST",
			AirlineName:            "Test Airlines",
			OriginIataCode:         "BOG",
			DestinationIataCode:    "JFK",
			RouteCode:              "BOG-JFK",
			OriginAirportName:      "El Dorado",
			DestinationAirportName: "John F Kennedy",
			AirportType:            "International",
			EstimatedFlightTime:    "5h 30m",
		}

		result := ar.ToDomain()

		if result.ID != "ar-123" {
			t.Errorf("expected ID 'ar-123', got %q", result.ID)
		}
		if result.RouteID != "route-456" {
			t.Errorf("expected RouteID 'route-456', got %q", result.RouteID)
		}
		if result.AirlineID != "airline-789" {
			t.Errorf("expected AirlineID 'airline-789', got %q", result.AirlineID)
		}
		if !result.IsActive() {
			t.Error("expected Status to be active")
		}
		if result.AirlineCode != "TST" {
			t.Errorf("expected AirlineCode 'TST', got %q", result.AirlineCode)
		}
		if result.AirlineName != "Test Airlines" {
			t.Errorf("expected AirlineName 'Test Airlines', got %q", result.AirlineName)
		}
		if result.OriginIataCode != "BOG" {
			t.Errorf("expected OriginIataCode 'BOG', got %q", result.OriginIataCode)
		}
		if result.DestinationIataCode != "JFK" {
			t.Errorf("expected DestinationIataCode 'JFK', got %q", result.DestinationIataCode)
		}
		if result.RouteCode != "BOG-JFK" {
			t.Errorf("expected RouteCode 'BOG-JFK', got %q", result.RouteCode)
		}
		if result.EstimatedFlightTime != "5h 30m" {
			t.Errorf("expected EstimatedFlightTime '5h 30m', got %q", result.EstimatedFlightTime)
		}
	})

	t.Run("converts inactive route", func(t *testing.T) {
		ar := &AirlineRoute{
			ID:     "ar-456",
			Status: domain.AirlineRouteStatusInactive,
		}

		result := ar.ToDomain()

		if result.IsActive() {
			t.Error("expected Status to be inactive")
		}
	})

	t.Run("converts pending route", func(t *testing.T) {
		ar := &AirlineRoute{
			ID:     "ar-789",
			Status: domain.AirlineRouteStatusPending,
		}

		result := ar.ToDomain()

		if !result.IsPending() {
			t.Error("expected Status to be pending")
		}
	})
}

func TestFromDomain(t *testing.T) {
	t.Run("converts all fields", func(t *testing.T) {
		dm := &domain.AirlineRoute{
			ID:                     "ar-123",
			RouteID:                "route-456",
			AirlineID:              "airline-789",
			Status:                 domain.AirlineRouteStatusActive,
			AirlineCode:            "DOM",
			AirlineName:            "Domain Airlines",
			OriginIataCode:         "MDE",
			DestinationIataCode:    "MIA",
			RouteCode:              "MDE-MIA",
			OriginAirportName:      "Jose Maria Cordova",
			DestinationAirportName: "Miami International",
			AirportType:            "International",
			EstimatedFlightTime:    "4h 15m",
		}

		result := FromDomain(dm)

		if result.ID != "ar-123" {
			t.Errorf("expected ID 'ar-123', got %q", result.ID)
		}
		if result.RouteID != "route-456" {
			t.Errorf("expected RouteID 'route-456', got %q", result.RouteID)
		}
		if result.AirlineID != "airline-789" {
			t.Errorf("expected AirlineID 'airline-789', got %q", result.AirlineID)
		}
		if result.Status != domain.AirlineRouteStatusActive {
			t.Error("expected Status to be active")
		}
		if result.AirlineCode != "DOM" {
			t.Errorf("expected AirlineCode 'DOM', got %q", result.AirlineCode)
		}
		if result.OriginIataCode != "MDE" {
			t.Errorf("expected OriginIataCode 'MDE', got %q", result.OriginIataCode)
		}
	})

	t.Run("roundtrip preserves data", func(t *testing.T) {
		original := &domain.AirlineRoute{
			ID:                     "ar-trip",
			RouteID:                "route-trip",
			AirlineID:              "airline-trip",
			Status:                 domain.AirlineRouteStatusActive,
			AirlineCode:            "RND",
			AirlineName:            "Roundtrip Airlines",
			OriginIataCode:         "CLO",
			DestinationIataCode:    "LAX",
			RouteCode:              "CLO-LAX",
			OriginAirportName:      "Alfonso Bonilla Aragon",
			DestinationAirportName: "Los Angeles International",
			AirportType:            "International",
			EstimatedFlightTime:    "7h 00m",
		}

		dbEntity := FromDomain(original)
		restored := dbEntity.ToDomain()

		if restored.ID != original.ID {
			t.Errorf("ID mismatch: expected %s, got %s", original.ID, restored.ID)
		}
		if restored.RouteCode != original.RouteCode {
			t.Errorf("RouteCode mismatch: expected %s, got %s", original.RouteCode, restored.RouteCode)
		}
		if restored.EstimatedFlightTime != original.EstimatedFlightTime {
			t.Errorf("EstimatedFlightTime mismatch: expected %s, got %s", original.EstimatedFlightTime, restored.EstimatedFlightTime)
		}
	})
}
