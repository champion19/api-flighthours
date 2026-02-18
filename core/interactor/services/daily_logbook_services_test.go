package services

import (
	"context"
	"errors"
	"testing"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
)

// mock daily logbook repository
type mockDailyLogbookRepo struct {
	getByIDFn      func(ctx context.Context, id string) (*domain.DailyLogbook, error)
	listFn         func(ctx context.Context, employeeID string, filters map[string]interface{}) ([]domain.DailyLogbook, error)
	saveFn         func(ctx context.Context, tx output.Tx, logbook domain.DailyLogbook) error
	updateFn       func(ctx context.Context, tx output.Tx, logbook domain.DailyLogbook) error
	beginTxFn      func(ctx context.Context) (output.Tx, error)
	updateStatusFn func(ctx context.Context, tx output.Tx, id string, status bool) error
	deleteFn       func(ctx context.Context, tx output.Tx, id string) error
}

func (m *mockDailyLogbookRepo) GetDailyLogbookByID(ctx context.Context, id string) (*domain.DailyLogbook, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *mockDailyLogbookRepo) ListDailyLogbooksByEmployee(ctx context.Context, employeeID string, filters map[string]interface{}) ([]domain.DailyLogbook, error) {
	if m.listFn != nil {
		return m.listFn(ctx, employeeID, filters)
	}
	return nil, nil
}

func (m *mockDailyLogbookRepo) SaveDailyLogbook(ctx context.Context, tx output.Tx, logbook domain.DailyLogbook) error {
	if m.saveFn != nil {
		return m.saveFn(ctx, tx, logbook)
	}
	return nil
}

func (m *mockDailyLogbookRepo) UpdateDailyLogbook(ctx context.Context, tx output.Tx, logbook domain.DailyLogbook) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, tx, logbook)
	}
	return nil
}

func (m *mockDailyLogbookRepo) BeginTx(ctx context.Context) (output.Tx, error) {
	if m.beginTxFn != nil {
		return m.beginTxFn(ctx)
	}
	return &mockTx{}, nil
}

func (m *mockDailyLogbookRepo) UpdateDailyLogbookStatus(ctx context.Context, tx output.Tx, id string, status bool) error {
	if m.updateStatusFn != nil {
		return m.updateStatusFn(ctx, tx, id, status)
	}
	return nil
}

func (m *mockDailyLogbookRepo) DeleteDailyLogbook(ctx context.Context, tx output.Tx, id string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, tx, id)
	}
	return nil
}

func TestNewDailyLogbookService(t *testing.T) {
	repo := &mockDailyLogbookRepo{}
	svc := NewDailyLogbookService(repo)
	if svc == nil {
		t.Error("expected non-nil service")
	}
}

func TestDailyLogbookService_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expected := &domain.DailyLogbook{ID: "lb-1", EmployeeID: "emp-1"}
		repo := &mockDailyLogbookRepo{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return expected, nil
			},
		}
		svc := NewDailyLogbookService(repo)
		result, err := svc.GetDailyLogbookByID(context.Background(), "lb-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.ID != "lb-1" {
			t.Errorf("expected ID 'lb-1', got %q", result.ID)
		}
	})

	t.Run("error", func(t *testing.T) {
		repo := &mockDailyLogbookRepo{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return nil, errors.New("not found")
			},
		}
		svc := NewDailyLogbookService(repo)
		_, err := svc.GetDailyLogbookByID(context.Background(), "lb-1")
		if err == nil {
			t.Error("expected error")
		}
	})
}

func TestDailyLogbookService_ListByEmployee(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expected := []domain.DailyLogbook{
			{ID: "lb-1", EmployeeID: "emp-1"},
			{ID: "lb-2", EmployeeID: "emp-1"},
		}
		repo := &mockDailyLogbookRepo{
			listFn: func(ctx context.Context, employeeID string, filters map[string]interface{}) ([]domain.DailyLogbook, error) {
				return expected, nil
			},
		}
		svc := NewDailyLogbookService(repo)
		result, err := svc.ListDailyLogbooksByEmployee(context.Background(), "emp-1", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 2 {
			t.Errorf("expected 2 logbooks, got %d", len(result))
		}
	})
}

func TestDailyLogbookService_BeginTx(t *testing.T) {
	repo := &mockDailyLogbookRepo{
		beginTxFn: func(ctx context.Context) (output.Tx, error) {
			return &mockTx{}, nil
		},
	}
	svc := NewDailyLogbookService(repo)
	tx, err := svc.BeginTx(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tx == nil {
		t.Error("expected non-nil tx")
	}
}

func TestDailyLogbookService_CreateDailyLogbookTx(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &mockDailyLogbookRepo{
			saveFn: func(ctx context.Context, tx output.Tx, logbook domain.DailyLogbook) error {
				return nil
			},
		}
		svc := NewDailyLogbookService(repo)
		err := svc.CreateDailyLogbookTx(context.Background(), &mockTx{}, domain.DailyLogbook{EmployeeID: "emp-1"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		repo := &mockDailyLogbookRepo{
			saveFn: func(ctx context.Context, tx output.Tx, logbook domain.DailyLogbook) error {
				return errors.New("save failed")
			},
		}
		svc := NewDailyLogbookService(repo)
		err := svc.CreateDailyLogbookTx(context.Background(), &mockTx{}, domain.DailyLogbook{})
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("generates ID when empty", func(t *testing.T) {
		var savedLogbook domain.DailyLogbook
		repo := &mockDailyLogbookRepo{
			saveFn: func(ctx context.Context, tx output.Tx, logbook domain.DailyLogbook) error {
				savedLogbook = logbook
				return nil
			},
		}
		svc := NewDailyLogbookService(repo)
		err := svc.CreateDailyLogbookTx(context.Background(), &mockTx{}, domain.DailyLogbook{EmployeeID: "emp-1"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if savedLogbook.ID == "" {
			t.Error("expected ID to be generated")
		}
	})
}

func TestDailyLogbookService_ActivateDailyLogbookTx(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &mockDailyLogbookRepo{
			updateStatusFn: func(ctx context.Context, tx output.Tx, id string, status bool) error {
				if !status {
					t.Error("expected status to be true for activate")
				}
				return nil
			},
		}
		svc := NewDailyLogbookService(repo)
		err := svc.ActivateDailyLogbookTx(context.Background(), &mockTx{}, "lb-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		repo := &mockDailyLogbookRepo{
			updateStatusFn: func(ctx context.Context, tx output.Tx, id string, status bool) error {
				return errors.New("activate failed")
			},
		}
		svc := NewDailyLogbookService(repo)
		err := svc.ActivateDailyLogbookTx(context.Background(), &mockTx{}, "lb-1")
		if err == nil {
			t.Error("expected error")
		}
	})
}

func TestDailyLogbookService_DeactivateDailyLogbookTx(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &mockDailyLogbookRepo{
			updateStatusFn: func(ctx context.Context, tx output.Tx, id string, status bool) error {
				if status {
					t.Error("expected status to be false for deactivate")
				}
				return nil
			},
		}
		svc := NewDailyLogbookService(repo)
		err := svc.DeactivateDailyLogbookTx(context.Background(), &mockTx{}, "lb-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		repo := &mockDailyLogbookRepo{
			updateStatusFn: func(ctx context.Context, tx output.Tx, id string, status bool) error {
				return errors.New("deactivate failed")
			},
		}
		svc := NewDailyLogbookService(repo)
		err := svc.DeactivateDailyLogbookTx(context.Background(), &mockTx{}, "lb-1")
		if err == nil {
			t.Error("expected error")
		}
	})
}

func TestDailyLogbookService_DeleteDailyLogbookTx(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &mockDailyLogbookRepo{
			deleteFn: func(ctx context.Context, tx output.Tx, id string) error {
				return nil
			},
		}
		svc := NewDailyLogbookService(repo)
		err := svc.DeleteDailyLogbookTx(context.Background(), &mockTx{}, "lb-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		repo := &mockDailyLogbookRepo{
			deleteFn: func(ctx context.Context, tx output.Tx, id string) error {
				return errors.New("delete failed")
			},
		}
		svc := NewDailyLogbookService(repo)
		err := svc.DeleteDailyLogbookTx(context.Background(), &mockTx{}, "lb-1")
		if err == nil {
			t.Error("expected error")
		}
	})
}
