//go:build integration
// +build integration

package engine

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

	if err := testContainer.SetupEngineSchema(ctx); err != nil {
		panic("Failed to setup schema: " + err.Error())
	}

	m.Run()

	testContainer.Stop(ctx)
}

func TestNewEngineRepository_Integration(t *testing.T) {
	if testContainer == nil {
		t.Skip("Integration tests require docker")
	}

	t.Run("creates repository successfully", func(t *testing.T) {
		repo, err := NewEngineRepository(testContainer.DB)
		if err != nil {
			t.Fatalf("failed to create repository: %v", err)
		}
		if repo == nil {
			t.Error("expected non-nil repository")
		}
	})

	t.Run("fails with nil db", func(t *testing.T) {
		_, err := NewEngineRepository(nil)
		if err == nil {
			t.Error("expected error with nil db")
		}
	})
}

func TestRepository_GetEngineByID_Integration(t *testing.T) {
	if testContainer == nil {
		t.Skip("Integration tests require docker")
	}
	ctx := context.Background()

	testContainer.CleanEngineTable(ctx)
	testContainer.InsertEngine(ctx, "engine-1", "CFM")

	repo, _ := NewEngineRepository(testContainer.DB)

	t.Run("finds existing engine", func(t *testing.T) {
		engine, err := repo.GetEngineByID(ctx, "engine-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if engine == nil {
			t.Fatal("expected non-nil engine")
		}
		if engine.Name != "CFM" {
			t.Errorf("expected 'CFM', got %q", engine.Name)
		}
	})

	t.Run("returns error for non-existent engine", func(t *testing.T) {
		_, err := repo.GetEngineByID(ctx, "non-existent")
		if err == nil {
			t.Error("expected error for non-existent engine")
		}
	})
}

func TestRepository_ListEngines_Integration(t *testing.T) {
	if testContainer == nil {
		t.Skip("Integration tests require docker")
	}
	ctx := context.Background()

	testContainer.CleanEngineTable(ctx)
	testContainer.InsertEngine(ctx, "e1", "CFM")
	testContainer.InsertEngine(ctx, "e2", "GE9")
	testContainer.InsertEngine(ctx, "e3", "LEA")

	repo, _ := NewEngineRepository(testContainer.DB)

	t.Run("lists all engines", func(t *testing.T) {
		engines, err := repo.ListEngines(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(engines) != 3 {
			t.Errorf("expected 3 engines, got %d", len(engines))
		}
	})
}
