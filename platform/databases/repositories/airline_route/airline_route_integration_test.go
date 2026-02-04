//go:build integration
// +build integration

package airline_route

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

	// Setup all required schemas in dependency order
	if err := testContainer.SetupAirlineRouteSchema(ctx); err != nil {
		panic("Failed to setup airline_route schema: " + err.Error())
	}

	m.Run()

	testContainer.Stop(ctx)
}

func setupTestData(ctx context.Context, t *testing.T) {
	t.Helper()
	// Clean all tables in reverse dependency order
	testContainer.CleanAirlineRouteTable(ctx)
	testContainer.CleanRouteTable(ctx)
	testContainer.CleanAirportTable(ctx)
	testContainer.CleanAirlineTable(ctx)

	// Insert airports
	testContainer.InsertAirport(ctx, "airport-bog", "El Dorado International", "Bogota", "Colombia", "BOG", "International", true)
	testContainer.InsertAirport(ctx, "airport-jfk", "John F Kennedy International", "New York", "USA", "JFK", "International", true)
	testContainer.InsertAirport(ctx, "airport-mde", "Jose Maria Cordova", "Medellin", "Colombia", "MDE", "Domestic", true)

	// Insert airline
	testContainer.InsertAirline(ctx, "airline-av", "Avianca", "AV", true)

	// Insert routes
	testContainer.InsertRoute(ctx, "route-bog-jfk", "airport-bog", "airport-jfk", "Colombia", "USA", "International", "05:30:00")
	testContainer.InsertRoute(ctx, "route-bog-mde", "airport-bog", "airport-mde", "Colombia", "Colombia", "Domestic", "00:45:00")
}

func TestNewAirlineRouteRepository_Integration(t *testing.T) {
	if testContainer == nil {
		t.Skip("Integration tests require docker")
	}

	t.Run("creates repository successfully", func(t *testing.T) {
		repo, err := NewAirlineRouteRepository(testContainer.DB)
		if err != nil {
			t.Fatalf("failed to create repository: %v", err)
		}
		if repo == nil {
			t.Error("expected non-nil repository")
		}
	})

	t.Run("fails with nil db", func(t *testing.T) {
		_, err := NewAirlineRouteRepository(nil)
		if err == nil {
			t.Error("expected error with nil db")
		}
	})
}

func TestRepository_ListAirlineRoutes_Integration(t *testing.T) {
	if testContainer == nil {
		t.Skip("Integration tests require docker")
	}
	ctx := context.Background()

	setupTestData(ctx, t)

	// Insert airline routes
	testContainer.InsertAirlineRoute(ctx, "ar-1", "route-bog-jfk", "airline-av", true)
	testContainer.InsertAirlineRoute(ctx, "ar-2", "route-bog-mde", "airline-av", true)

	repo, _ := NewAirlineRouteRepository(testContainer.DB)

	t.Run("lists all airline routes", func(t *testing.T) {
		routes, err := repo.ListAirlineRoutes(ctx, map[string]interface{}{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(routes) < 2 {
			t.Errorf("expected at least 2 routes, got %d", len(routes))
		}
	})

	t.Run("filters by airline_id", func(t *testing.T) {
		routes, err := repo.ListAirlineRoutes(ctx, map[string]interface{}{"airline_id": "airline-av"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// All routes should belong to airline-av
		for _, r := range routes {
			if r.AirlineID != "airline-av" {
				t.Errorf("expected airline_id 'airline-av', got %q", r.AirlineID)
			}
		}
	})

	t.Run("filters by status", func(t *testing.T) {
		routes, err := repo.ListAirlineRoutes(ctx, map[string]interface{}{"status": true})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		for _, r := range routes {
			if !r.Status {
				t.Error("expected only active routes")
			}
		}
	})
}

func TestRepository_GetAirlineRouteByID_Integration(t *testing.T) {
	if testContainer == nil {
		t.Skip("Integration tests require docker")
	}
	ctx := context.Background()

	setupTestData(ctx, t)
	testContainer.InsertAirlineRoute(ctx, "ar-get-1", "route-bog-jfk", "airline-av", true)

	repo, _ := NewAirlineRouteRepository(testContainer.DB)

	t.Run("finds existing airline route", func(t *testing.T) {
		route, err := repo.GetAirlineRouteByID(ctx, "ar-get-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if route == nil {
			t.Fatal("expected non-nil route")
		}
		if route.ID != "ar-get-1" {
			t.Errorf("expected ID 'ar-get-1', got %q", route.ID)
		}
		// Verify joined data comes through
		if route.AirlineCode == "" {
			t.Error("expected AirlineCode from join")
		}
		if route.OriginIataCode == "" {
			t.Error("expected OriginIataCode from join")
		}
	})

	t.Run("returns error for non-existent route", func(t *testing.T) {
		_, err := repo.GetAirlineRouteByID(ctx, "non-existent")
		if err == nil {
			t.Error("expected error for non-existent route")
		}
	})
}

// Note: UpdateAirlineRouteStatus requires a transaction (output.Tx)
// which is more complex to test in isolation. The coverage for this method
// is better achieved through service-level integration tests.

func TestRepository_ListAirlineRoutesByAirlineID_Integration(t *testing.T) {
	if testContainer == nil {
		t.Skip("Integration tests require docker")
	}
	ctx := context.Background()

	setupTestData(ctx, t)
	// Insert active and inactive routes
	testContainer.InsertAirlineRoute(ctx, "ar-by-airline-1", "route-bog-jfk", "airline-av", true)
	testContainer.InsertAirlineRoute(ctx, "ar-by-airline-2", "route-bog-mde", "airline-av", false) // inactive

	repo, _ := NewAirlineRouteRepository(testContainer.DB)

	t.Run("lists only active routes by airline ID", func(t *testing.T) {
		routes, err := repo.ListAirlineRoutesByAirlineID(ctx, "airline-av")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// Should only return active routes
		for _, r := range routes {
			if !r.Status {
				t.Errorf("expected only active routes, got inactive route: %s", r.ID)
			}
		}
	})

	t.Run("returns empty for non-existent airline", func(t *testing.T) {
		routes, err := repo.ListAirlineRoutesByAirlineID(ctx, "non-existent-airline")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(routes) != 0 {
			t.Errorf("expected 0 routes, got %d", len(routes))
		}
	})
}
