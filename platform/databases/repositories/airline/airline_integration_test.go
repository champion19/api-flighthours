//go:build integration
// +build integration

package airline

import (
	"context"
	"testing"

	"github.com/champion19/api-flighthours/platform/databases/testhelper"
)

var testContainer *testhelper.MySQLContainer

func TestMain(m *testing.M) {
	// Skip if not running integration tests
	ctx := context.Background()

	container, err := testhelper.StartMySQL(ctx)
	if err != nil {
		panic("Failed to start MySQL container: " + err.Error())
	}
	testContainer = container

	// Setup schema
	if err := testContainer.SetupAirlineSchema(ctx); err != nil {
		panic("Failed to setup schema: " + err.Error())
	}

	// Run tests
	m.Run()

	// Cleanup
	testContainer.Stop(ctx)
}

func TestNewAirlineRepository_Integration(t *testing.T) {
	if testContainer == nil {
		t.Skip("Integration tests require docker")
	}

	t.Run("creates repository successfully", func(t *testing.T) {
		repo, err := NewAirlineRepository(testContainer.DB)
		if err != nil {
			t.Fatalf("failed to create repository: %v", err)
		}
		if repo == nil {
			t.Error("expected non-nil repository")
		}
	})

	t.Run("fails with nil db", func(t *testing.T) {
		_, err := NewAirlineRepository(nil)
		if err == nil {
			t.Error("expected error with nil db")
		}
	})
}

func TestRepository_GetAirlineByID_Integration(t *testing.T) {
	if testContainer == nil {
		t.Skip("Integration tests require docker")
	}
	ctx := context.Background()

	// Clean and seed data
	testContainer.CleanAirlineTable(ctx)
	testContainer.InsertAirline(ctx, "test-id-123", "Avianca", "AV", true)

	repo, _ := NewAirlineRepository(testContainer.DB)

	t.Run("finds existing airline", func(t *testing.T) {
		airline, err := repo.GetAirlineByID(ctx, "test-id-123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if airline == nil {
			t.Fatal("expected non-nil airline")
		}
		if airline.AirlineName != "Avianca" {
			t.Errorf("expected Avianca, got %s", airline.AirlineName)
		}
	})

	t.Run("returns error for non-existent airline", func(t *testing.T) {
		_, err := repo.GetAirlineByID(ctx, "non-existent")
		if err == nil {
			t.Error("expected error for non-existent airline")
		}
	})
}

func TestRepository_ListAirlines_Integration(t *testing.T) {
	if testContainer == nil {
		t.Skip("Integration tests require docker")
	}
	ctx := context.Background()

	// Clean and seed data
	testContainer.CleanAirlineTable(ctx)
	testContainer.InsertAirline(ctx, "id-1", "Avianca", "AV", true)
	testContainer.InsertAirline(ctx, "id-2", "LATAM", "LA", true)
	testContainer.InsertAirline(ctx, "id-3", "Aerocivil", "AC", false)

	repo, _ := NewAirlineRepository(testContainer.DB)

	t.Run("lists all airlines", func(t *testing.T) {
		airlines, err := repo.ListAirlines(ctx, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(airlines) != 3 {
			t.Errorf("expected 3 airlines, got %d", len(airlines))
		}
	})

	t.Run("filters by active status", func(t *testing.T) {
		filters := map[string]interface{}{"status": true}
		airlines, err := repo.ListAirlines(ctx, filters)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(airlines) != 2 {
			t.Errorf("expected 2 active airlines, got %d", len(airlines))
		}
	})

	t.Run("filters by inactive status", func(t *testing.T) {
		filters := map[string]interface{}{"status": false}
		airlines, err := repo.ListAirlines(ctx, filters)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(airlines) != 1 {
			t.Errorf("expected 1 inactive airline, got %d", len(airlines))
		}
	})
}

func TestRepository_BeginTx_Integration(t *testing.T) {
	if testContainer == nil {
		t.Skip("Integration tests require docker")
	}
	ctx := context.Background()

	repo, _ := NewAirlineRepository(testContainer.DB)

	t.Run("begins transaction successfully", func(t *testing.T) {
		tx, err := repo.BeginTx(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tx == nil {
			t.Error("expected non-nil transaction")
		}
	})
}
