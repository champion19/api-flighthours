//go:build integration
// +build integration

package employee

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

	if err := testContainer.SetupEmployeeSchema(ctx); err != nil {
		panic("Failed to setup schema: " + err.Error())
	}

	m.Run()

	testContainer.Stop(ctx)
}

func TestNewClientRepository_Integration(t *testing.T) {
	if testContainer == nil {
		t.Skip("Integration tests require docker")
	}

	t.Run("creates repository successfully", func(t *testing.T) {
		repo, err := NewClientRepository(testContainer.DB)
		if err != nil {
			t.Fatalf("failed to create repository: %v", err)
		}
		if repo == nil {
			t.Error("expected non-nil repository")
		}
	})

	t.Run("fails with nil db", func(t *testing.T) {
		_, err := NewClientRepository(nil)
		if err == nil {
			t.Error("expected error with nil db")
		}
	})
}

func TestRepository_GetEmployeeByID_Integration(t *testing.T) {
	if testContainer == nil {
		t.Skip("Integration tests require docker")
	}
	ctx := context.Background()

	testContainer.CleanEmployeeTable(ctx)
	testContainer.InsertEmployee(ctx, "emp-123", "Juan Piloto", "juan@test.com", "123456789", "pilot", true)

	repo, _ := NewClientRepository(testContainer.DB)

	t.Run("finds existing employee by ID", func(t *testing.T) {
		emp, err := repo.GetEmployeeByID(ctx, "emp-123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if emp == nil {
			t.Fatal("expected non-nil employee")
		}
		if emp.Name != "Juan Piloto" {
			t.Errorf("expected 'Juan Piloto', got %q", emp.Name)
		}
	})

	t.Run("returns error for non-existent ID", func(t *testing.T) {
		_, err := repo.GetEmployeeByID(ctx, "non-existent")
		if err == nil {
			t.Error("expected error for non-existent employee")
		}
	})
}

func TestRepository_GetEmployeeByEmail_Integration(t *testing.T) {
	if testContainer == nil {
		t.Skip("Integration tests require docker")
	}
	ctx := context.Background()

	testContainer.CleanEmployeeTable(ctx)
	testContainer.InsertEmployee(ctx, "emp-456", "Maria Admin", "maria@test.com", "987654321", "admin", true)

	repo, _ := NewClientRepository(testContainer.DB)

	t.Run("finds existing employee by email", func(t *testing.T) {
		emp, err := repo.GetEmployeeByEmail(ctx, "maria@test.com")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if emp == nil {
			t.Fatal("expected non-nil employee")
		}
		if emp.Email != "maria@test.com" {
			t.Errorf("expected 'maria@test.com', got %q", emp.Email)
		}
	})
}

func TestRepository_GetEmployeesByRole_Integration(t *testing.T) {
	if testContainer == nil {
		t.Skip("Integration tests require docker")
	}
	ctx := context.Background()

	testContainer.CleanEmployeeTable(ctx)
	testContainer.InsertEmployee(ctx, "p1", "Pilot1", "pilot1@test.com", "111111111", "pilot", true)
	testContainer.InsertEmployee(ctx, "p2", "Pilot2", "pilot2@test.com", "222222222", "pilot", true)
	testContainer.InsertEmployee(ctx, "a1", "Admin1", "admin1@test.com", "333333333", "admin", true)

	repo, _ := NewClientRepository(testContainer.DB)

	t.Run("filters employees by role", func(t *testing.T) {
		pilots, err := repo.GetEmployeesByRole(ctx, "pilot")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(pilots) != 2 {
			t.Errorf("expected 2 pilots, got %d", len(pilots))
		}
	})

	t.Run("returns error for non-existent role", func(t *testing.T) {
		_, err := repo.GetEmployeesByRole(ctx, "nonexistent")
		if err == nil {
			t.Error("expected error for non-existent role")
		}
	})
}

func TestRepository_BeginTx_Integration(t *testing.T) {
	if testContainer == nil {
		t.Skip("Integration tests require docker")
	}
	ctx := context.Background()

	repo, _ := NewClientRepository(testContainer.DB)

	t.Run("begins transaction successfully", func(t *testing.T) {
		tx, err := repo.BeginTx(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tx == nil {
			t.Error("expected non-nil transaction")
		}
		tx.Rollback()
	})
}
