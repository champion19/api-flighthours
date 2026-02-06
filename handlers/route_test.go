package handlers

import (
	"errors"
	"testing"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

func TestFromDomainRoute(t *testing.T) {
	t.Run("converts domain route to response", func(t *testing.T) {
		route := &domain.Route{
			ID:                     "raw-id",
			OriginAirportID:        "origin-airport-id",
			OriginIataCode:         "BOG",
			OriginAirportName:      "El Dorado International",
			DestinationAirportID:   "dest-airport-id",
			DestinationIataCode:    "MDE",
			DestinationAirportName: "José María Córdova",
			AirportType:            "domestic",
			EstimatedFlightTime:    "1h 10m",
			RouteCode:              "BOG-MDE",
		}

		result := FromDomainRoute(route, "encoded-id", "encoded-origin", "encoded-dest")

		if result.ID != "encoded-id" {
			t.Errorf("expected 'encoded-id', got %q", result.ID)
		}
		if result.OriginAirportID != "encoded-origin" {
			t.Errorf("expected 'encoded-origin', got %q", result.OriginAirportID)
		}
		if result.DestinationAirportID != "encoded-dest" {
			t.Errorf("expected 'encoded-dest', got %q", result.DestinationAirportID)
		}
		if result.OriginIataCode != "BOG" {
			t.Errorf("expected 'BOG', got %q", result.OriginIataCode)
		}
		if result.RouteCode != "BOG-MDE" {
			t.Errorf("expected 'BOG-MDE', got %q", result.RouteCode)
		}
	})
}

func TestToRouteListResponse(t *testing.T) {
	t.Run("converts route slice to list response", func(t *testing.T) {
		routes := []domain.Route{
			{
				ID:                   "id1",
				OriginAirportID:      "origin1",
				DestinationAirportID: "dest1",
				OriginIataCode:       "BOG",
				DestinationIataCode:  "MDE",
			},
			{
				ID:                   "id2",
				OriginAirportID:      "origin2",
				DestinationAirportID: "dest2",
				OriginIataCode:       "MDE",
				DestinationIataCode:  "CTG",
			},
		}

		encodeFunc := func(id string) (string, error) {
			return "encoded-" + id, nil
		}

		result := ToRouteListResponse(routes, encodeFunc, "http://api.test")

		if result.Total != 2 {
			t.Errorf("expected total 2, got %d", result.Total)
		}
		if len(result.Routes) != 2 {
			t.Errorf("expected 2 routes, got %d", len(result.Routes))
		}
		if result.Routes[0].ID != "encoded-id1" {
			t.Errorf("expected 'encoded-id1', got %q", result.Routes[0].ID)
		}
	})

	t.Run("handles empty slice", func(t *testing.T) {
		routes := []domain.Route{}
		encodeFunc := func(id string) (string, error) {
			return "encoded-" + id, nil
		}

		result := ToRouteListResponse(routes, encodeFunc, "http://api.test")

		if result.Total != 0 {
			t.Errorf("expected total 0, got %d", result.Total)
		}
	})

	t.Run("uses original ID when encoding fails", func(t *testing.T) {
		routes := []domain.Route{
			{
				ID:                   "fallback-id",
				OriginAirportID:      "origin-fallback",
				DestinationAirportID: "dest-fallback",
			},
		}

		encodeFunc := func(id string) (string, error) {
			return "", errors.New("encoding failed")
		}

		result := ToRouteListResponse(routes, encodeFunc, "http://api.test")

		if result.Routes[0].ID != "fallback-id" {
			t.Errorf("expected 'fallback-id', got %q", result.Routes[0].ID)
		}
	})
}
