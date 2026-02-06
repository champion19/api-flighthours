//go:build integration
// +build integration

package route

import (
	"context"
	"testing"

	"github.com/champion19/api-flighthours/platform/databases/testhelper"
)

var testContainer *testhelper.MySQLContainer

func TestMain(m *testing.M) {
	ctx := context.Background()

	container, err := testhelper.StartMySQL(ctx)
	if err != nil {
		panic("Failed to start MySQL container: " + err.Error())
	}
	testContainer = container

	// Setup route schema (also sets up airport)
	if err := testContainer.SetupRouteSchema(ctx); err != nil {
		panic("Failed to setup route schema: " + err.Error())
	}

	m.Run()

	testContainer.Stop(ctx)
}

func setupTestData(ctx context.Context, t *testing.T) {
	t.Helper()
	// Clean tables in reverse dependency order
	testContainer.CleanRouteTable(ctx)
	testContainer.CleanAirportTable(ctx)

	// Insert airports
	testContainer.InsertAirport(ctx, "airport-bog", "El Dorado International", "BOG", "International", true)
	testContainer.InsertAirport(ctx, "airport-jfk", "John F Kennedy International", "JFK", "International", true)
	testContainer.InsertAirport(ctx, "airport-mde", "Jose Maria Cordova", "MDE", "Domestic", true)
	testContainer.InsertAirport(ctx, "airport-clo", "Alfonso Bonilla Aragon", "CLO", "Domestic", true)
}

func TestNewRouteRepository_Integration(t *testing.T) {
	if testContainer == nil {
		t.Skip("Integration tests require docker")
	}

	t.Run("creates repository successfully", func(t *testing.T) {
		repo, err := NewRouteRepository(testContainer.DB)
		if err != nil {
			t.Fatalf("failed to create repository: %v", err)
		}
		if repo == nil {
			t.Error("expected non-nil repository")
		}
	})

	t.Run("fails with nil db", func(t *testing.T) {
		_, err := NewRouteRepository(nil)
		if err == nil {
			t.Error("expected error with nil db")
		}
	})
}

func TestRepository_ListRoutes_Integration(t *testing.T) {
	if testContainer == nil {
		t.Skip("Integration tests require docker")
	}
	ctx := context.Background()

	setupTestData(ctx, t)

	// Insert routes
	testContainer.InsertRoute(ctx, "route-bog-jfk", "airport-bog", "airport-jfk", "International", "05:30:00")
	testContainer.InsertRoute(ctx, "route-bog-mde", "airport-bog", "airport-mde", "Domestic", "00:45:00")
	testContainer.InsertRoute(ctx, "route-mde-clo", "airport-mde", "airport-clo", "Domestic", "00:35:00")

	repo, _ := NewRouteRepository(testContainer.DB)

	t.Run("lists all routes", func(t *testing.T) {
		routes, err := repo.ListRoutes(ctx, map[string]interface{}{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(routes) < 3 {
			t.Errorf("expected at least 3 routes, got %d", len(routes))
		}
	})

	t.Run("filters by airport_type International", func(t *testing.T) {
		routes, err := repo.ListRoutes(ctx, map[string]interface{}{"airport_type": "International"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		for _, r := range routes {
			if r.AirportType != "International" {
				t.Errorf("expected airport_type 'International', got %q", r.AirportType)
			}
		}
	})

	t.Run("filters by airport_type Domestic", func(t *testing.T) {
		routes, err := repo.ListRoutes(ctx, map[string]interface{}{"airport_type": "Domestic"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(routes) < 2 {
			t.Errorf("expected at least 2 domestic routes, got %d", len(routes))
		}
		for _, r := range routes {
			if r.AirportType != "Domestic" {
				t.Errorf("expected airport_type 'Domestic', got %q", r.AirportType)
			}
		}
	})
}

func TestRepository_GetRouteByID_Integration(t *testing.T) {
	if testContainer == nil {
		t.Skip("Integration tests require docker")
	}
	ctx := context.Background()

	setupTestData(ctx, t)
	testContainer.InsertRoute(ctx, "route-get-1", "airport-bog", "airport-jfk", "International", "05:30:00")

	repo, _ := NewRouteRepository(testContainer.DB)

	t.Run("finds existing route", func(t *testing.T) {
		route, err := repo.GetRouteByID(ctx, "route-get-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if route == nil {
			t.Fatal("expected non-nil route")
		}
		if route.ID != "route-get-1" {
			t.Errorf("expected ID 'route-get-1', got %q", route.ID)
		}
		// Verify joined data
		if route.OriginIataCode != "BOG" {
			t.Errorf("expected OriginIataCode 'BOG', got %q", route.OriginIataCode)
		}
		if route.DestinationIataCode != "JFK" {
			t.Errorf("expected DestinationIataCode 'JFK', got %q", route.DestinationIataCode)
		}
		if route.RouteCode != "BOG-JFK" {
			t.Errorf("expected RouteCode 'BOG-JFK', got %q", route.RouteCode)
		}
	})

	t.Run("returns error for non-existent route", func(t *testing.T) {
		_, err := repo.GetRouteByID(ctx, "non-existent")
		if err == nil {
			t.Error("expected error for non-existent route")
		}
	})
}
