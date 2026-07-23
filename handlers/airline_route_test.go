package handlers

import (
	"errors"
	"testing"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

func TestFromDomainAirlineRoute(t *testing.T) {
	t.Run("converts domain airline route to response", func(t *testing.T) {
		ar := &domain.AirlineRoute{
			ID:                     "raw-id",
			RouteID:                "route-id",
			AirlineID:              "airline-id",
			Status:                 domain.AirlineRouteStatusActive,
			AirlineCode:            "AV",
			AirlineName:            "Avianca",
			OriginIataCode:         "BOG",
			DestinationIataCode:    "MDE",
			RouteCode:              "BOG-MDE",
			OriginAirportName:      "El Dorado",
			DestinationAirportName: "Jose Maria Cordova",
			AirportType:            "domestic",
			EstimatedFlightTime:    "1h",
		}

		result := FromDomainAirlineRoute(ar, "enc-id", "enc-route", "enc-airline")

		if result.ID != "enc-id" {
			t.Errorf("expected 'enc-id', got %q", result.ID)
		}
		if result.RouteID != "enc-route" {
			t.Errorf("expected 'enc-route', got %q", result.RouteID)
		}
		if result.AirlineID != "enc-airline" {
			t.Errorf("expected 'enc-airline', got %q", result.AirlineID)
		}
		if result.AirlineCode != "AV" {
			t.Errorf("expected 'AV', got %q", result.AirlineCode)
		}
		if result.Status != domain.AirlineRouteStatusActive {
			t.Error("expected Status to be active")
		}
	})
}

func TestToAirlineRouteListResponse(t *testing.T) {
	t.Run("converts slice with links", func(t *testing.T) {
		routes := []domain.AirlineRoute{
			{
				ID:             "ar1",
				RouteID:        "r1",
				AirlineID:      "a1",
				Status:         domain.AirlineRouteStatusActive,
				OriginIataCode: "BOG",
			},
			{
				ID:             "ar2",
				RouteID:        "r2",
				AirlineID:      "a2",
				Status:         domain.AirlineRouteStatusInactive,
				OriginIataCode: "MDE",
			},
		}

		encodeFunc := func(id string) (string, error) {
			return "enc-" + id, nil
		}

		result := ToAirlineRouteListResponse(routes, encodeFunc, "http://api.test")

		if result.Total != 2 {
			t.Errorf("expected total 2, got %d", result.Total)
		}
		if len(result.AirlineRoutes) != 2 {
			t.Errorf("expected 2 routes, got %d", len(result.AirlineRoutes))
		}
		if result.AirlineRoutes[0].ID != "enc-ar1" {
			t.Errorf("expected 'enc-ar1', got %q", result.AirlineRoutes[0].ID)
		}
		// Should have HATEOAS links
		if len(result.AirlineRoutes[0].Links) == 0 {
			t.Error("expected HATEOAS links on items")
		}
		if len(result.Links) == 0 {
			t.Error("expected collection-level links")
		}
	})

	t.Run("handles empty slice", func(t *testing.T) {
		routes := []domain.AirlineRoute{}
		encodeFunc := func(id string) (string, error) {
			return "enc-" + id, nil
		}

		result := ToAirlineRouteListResponse(routes, encodeFunc, "http://api.test")

		if result.Total != 0 {
			t.Errorf("expected total 0, got %d", result.Total)
		}
	})

	t.Run("uses original IDs when encoding fails", func(t *testing.T) {
		routes := []domain.AirlineRoute{
			{ID: "fallback-id", RouteID: "fallback-route", AirlineID: "fallback-airline"},
		}

		encodeFunc := func(id string) (string, error) {
			return "", errors.New("encoding failed")
		}

		result := ToAirlineRouteListResponse(routes, encodeFunc, "http://api.test")

		if result.AirlineRoutes[0].ID != "fallback-id" {
			t.Errorf("expected 'fallback-id', got %q", result.AirlineRoutes[0].ID)
		}
	})

	t.Run("works without base URL", func(t *testing.T) {
		routes := []domain.AirlineRoute{
			{ID: "id1", RouteID: "r1", AirlineID: "a1"},
		}

		encodeFunc := func(id string) (string, error) {
			return "enc-" + id, nil
		}

		result := ToAirlineRouteListResponse(routes, encodeFunc, "")

		if len(result.AirlineRoutes[0].Links) != 0 {
			t.Error("expected no links with empty baseURL")
		}
	})
}
