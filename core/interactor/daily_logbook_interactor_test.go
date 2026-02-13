package interactor

import (
	"context"
	"errors"
	"testing"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
)

// fakeDailyLogbookService implements input.DailyLogbookService
type fakeDailyLogbookService struct {
	getByIDFn      func(ctx context.Context, id string) (*domain.DailyLogbook, error)
	listFn         func(ctx context.Context, employeeID string, filters map[string]interface{}) ([]domain.DailyLogbook, error)
	createFn       func(ctx context.Context, logbook domain.DailyLogbook) error
	createTxFn     func(ctx context.Context, tx output.Tx, logbook domain.DailyLogbook) error
	updateTxFn     func(ctx context.Context, tx output.Tx, logbook domain.DailyLogbook) error
	activateTxFn   func(ctx context.Context, tx output.Tx, id string) error
	deactivateTxFn func(ctx context.Context, tx output.Tx, id string) error
	beginTxFn      func(ctx context.Context) (output.Tx, error)
}

func (f *fakeDailyLogbookService) BeginTx(ctx context.Context) (output.Tx, error) {
	if f.beginTxFn != nil {
		return f.beginTxFn(ctx)
	}
	return &fakeTx{}, nil
}

func (f *fakeDailyLogbookService) GetDailyLogbookByID(ctx context.Context, id string) (*domain.DailyLogbook, error) {
	if f.getByIDFn != nil {
		return f.getByIDFn(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (f *fakeDailyLogbookService) ListDailyLogbooksByEmployee(ctx context.Context, employeeID string, filters map[string]interface{}) ([]domain.DailyLogbook, error) {
	if f.listFn != nil {
		return f.listFn(ctx, employeeID, filters)
	}
	return nil, errors.New("not implemented")
}

func (f *fakeDailyLogbookService) CreateDailyLogbook(ctx context.Context, logbook domain.DailyLogbook) error {
	if f.createFn != nil {
		return f.createFn(ctx, logbook)
	}
	return nil
}

func (f *fakeDailyLogbookService) CreateDailyLogbookTx(ctx context.Context, tx output.Tx, logbook domain.DailyLogbook) error {
	if f.createTxFn != nil {
		return f.createTxFn(ctx, tx, logbook)
	}
	return nil
}

func (f *fakeDailyLogbookService) UpdateDailyLogbookTx(ctx context.Context, tx output.Tx, logbook domain.DailyLogbook) error {
	if f.updateTxFn != nil {
		return f.updateTxFn(ctx, tx, logbook)
	}
	return nil
}

func (f *fakeDailyLogbookService) ActivateDailyLogbookTx(ctx context.Context, tx output.Tx, id string) error {
	if f.activateTxFn != nil {
		return f.activateTxFn(ctx, tx, id)
	}
	return nil
}

func (f *fakeDailyLogbookService) DeactivateDailyLogbookTx(ctx context.Context, tx output.Tx, id string) error {
	if f.deactivateTxFn != nil {
		return f.deactivateTxFn(ctx, tx, id)
	}
	return nil
}

func TestNewDailyLogbookInteractor(t *testing.T) {
	svc := &fakeDailyLogbookService{}
	inter := NewDailyLogbookInteractor(svc)
	if inter == nil {
		t.Error("expected non-nil interactor")
	}
}

func TestDailyLogbookInteractor_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expected := &domain.DailyLogbook{ID: "lb-1", EmployeeID: "emp-1"}
		svc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return expected, nil
			},
		}
		inter := NewDailyLogbookInteractor(svc)
		result, err := inter.GetDailyLogbookByID(context.Background(), "lb-1", "emp-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.ID != "lb-1" {
			t.Errorf("expected ID 'lb-1', got %q", result.ID)
		}
	})

	t.Run("error", func(t *testing.T) {
		svc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return nil, errors.New("not found")
			},
		}
		inter := NewDailyLogbookInteractor(svc)
		_, err := inter.GetDailyLogbookByID(context.Background(), "lb-1", "emp-1")
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("unauthorized", func(t *testing.T) {
		svc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: "lb-1", EmployeeID: "emp-other"}, nil
			},
		}
		inter := NewDailyLogbookInteractor(svc)
		_, err := inter.GetDailyLogbookByID(context.Background(), "lb-1", "emp-1")
		if err == nil {
			t.Error("expected unauthorized error")
		}
	})
}

func TestDailyLogbookInteractor_ListByEmployee(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expected := []domain.DailyLogbook{{ID: "lb-1"}, {ID: "lb-2"}}
		svc := &fakeDailyLogbookService{
			listFn: func(ctx context.Context, employeeID string, filters map[string]interface{}) ([]domain.DailyLogbook, error) {
				return expected, nil
			},
		}
		inter := NewDailyLogbookInteractor(svc)
		result, err := inter.ListDailyLogbooksByEmployee(context.Background(), "emp-1", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 2 {
			t.Errorf("expected 2 logbooks, got %d", len(result))
		}
	})

	t.Run("error", func(t *testing.T) {
		svc := &fakeDailyLogbookService{
			listFn: func(ctx context.Context, employeeID string, filters map[string]interface{}) ([]domain.DailyLogbook, error) {
				return nil, errors.New("db error")
			},
		}
		inter := NewDailyLogbookInteractor(svc)
		_, err := inter.ListDailyLogbooksByEmployee(context.Background(), "emp-1", nil)
		if err == nil {
			t.Error("expected error")
		}
	})
}

func TestDailyLogbookInteractor_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &fakeDailyLogbookService{
			createTxFn: func(ctx context.Context, tx output.Tx, logbook domain.DailyLogbook) error {
				return nil
			},
		}
		inter := NewDailyLogbookInteractor(svc)
		err := inter.CreateDailyLogbook(context.Background(), domain.DailyLogbook{EmployeeID: "emp-1"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		svc := &fakeDailyLogbookService{
			createTxFn: func(ctx context.Context, tx output.Tx, logbook domain.DailyLogbook) error {
				return errors.New("create error")
			},
		}
		inter := NewDailyLogbookInteractor(svc)
		err := inter.CreateDailyLogbook(context.Background(), domain.DailyLogbook{EmployeeID: "emp-1"})
		if err == nil {
			t.Error("expected error")
		}
	})
}

func TestDailyLogbookInteractor_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: id, EmployeeID: "emp-1"}, nil
			},
			updateTxFn: func(ctx context.Context, tx output.Tx, logbook domain.DailyLogbook) error {
				return nil
			},
		}
		inter := NewDailyLogbookInteractor(svc)
		err := inter.UpdateDailyLogbook(context.Background(), domain.DailyLogbook{ID: "lb-1", EmployeeID: "emp-1"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		svc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return nil, errors.New("not found")
			},
		}
		inter := NewDailyLogbookInteractor(svc)
		err := inter.UpdateDailyLogbook(context.Background(), domain.DailyLogbook{ID: "lb-1"})
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("begin tx error", func(t *testing.T) {
		svc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: id}, nil
			},
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return nil, errors.New("tx error")
			},
		}
		inter := NewDailyLogbookInteractor(svc)
		err := inter.UpdateDailyLogbook(context.Background(), domain.DailyLogbook{ID: "lb-1"})
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("update tx error", func(t *testing.T) {
		svc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: id}, nil
			},
			updateTxFn: func(ctx context.Context, tx output.Tx, logbook domain.DailyLogbook) error {
				return errors.New("update failed")
			},
		}
		inter := NewDailyLogbookInteractor(svc)
		err := inter.UpdateDailyLogbook(context.Background(), domain.DailyLogbook{ID: "lb-1"})
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("commit error", func(t *testing.T) {
		svc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: id}, nil
			},
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return &fakeTx{commitFn: func() error { return errors.New("commit error") }}, nil
			},
			updateTxFn: func(ctx context.Context, tx output.Tx, logbook domain.DailyLogbook) error {
				return nil
			},
		}
		inter := NewDailyLogbookInteractor(svc)
		err := inter.UpdateDailyLogbook(context.Background(), domain.DailyLogbook{ID: "lb-1"})
		if err == nil {
			t.Error("expected commit error")
		}
	})
}

func TestDailyLogbookInteractor_Activate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: id}, nil
			},
			activateTxFn: func(ctx context.Context, tx output.Tx, id string) error {
				return nil
			},
		}
		inter := NewDailyLogbookInteractor(svc)
		err := inter.ActivateDailyLogbook(context.Background(), "lb-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		svc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return nil, errors.New("not found")
			},
		}
		inter := NewDailyLogbookInteractor(svc)
		err := inter.ActivateDailyLogbook(context.Background(), "lb-1")
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("begin tx error", func(t *testing.T) {
		svc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: id}, nil
			},
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return nil, errors.New("tx error")
			},
		}
		inter := NewDailyLogbookInteractor(svc)
		err := inter.ActivateDailyLogbook(context.Background(), "lb-1")
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("activate tx error", func(t *testing.T) {
		svc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: id}, nil
			},
			activateTxFn: func(ctx context.Context, tx output.Tx, id string) error {
				return errors.New("activate failed")
			},
		}
		inter := NewDailyLogbookInteractor(svc)
		err := inter.ActivateDailyLogbook(context.Background(), "lb-1")
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("commit error", func(t *testing.T) {
		svc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: id}, nil
			},
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return &fakeTx{commitFn: func() error { return errors.New("commit error") }}, nil
			},
			activateTxFn: func(ctx context.Context, tx output.Tx, id string) error {
				return nil
			},
		}
		inter := NewDailyLogbookInteractor(svc)
		err := inter.ActivateDailyLogbook(context.Background(), "lb-1")
		if err == nil {
			t.Error("expected commit error")
		}
	})
}

func TestDailyLogbookInteractor_Deactivate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: id}, nil
			},
			deactivateTxFn: func(ctx context.Context, tx output.Tx, id string) error {
				return nil
			},
		}
		inter := NewDailyLogbookInteractor(svc)
		err := inter.DeactivateDailyLogbook(context.Background(), "lb-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		svc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return nil, errors.New("not found")
			},
		}
		inter := NewDailyLogbookInteractor(svc)
		err := inter.DeactivateDailyLogbook(context.Background(), "lb-1")
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("begin tx error", func(t *testing.T) {
		svc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: id}, nil
			},
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return nil, errors.New("tx error")
			},
		}
		inter := NewDailyLogbookInteractor(svc)
		err := inter.DeactivateDailyLogbook(context.Background(), "lb-1")
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("deactivate tx error", func(t *testing.T) {
		svc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: id}, nil
			},
			deactivateTxFn: func(ctx context.Context, tx output.Tx, id string) error {
				return errors.New("deactivate failed")
			},
		}
		inter := NewDailyLogbookInteractor(svc)
		err := inter.DeactivateDailyLogbook(context.Background(), "lb-1")
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("commit error", func(t *testing.T) {
		svc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: id}, nil
			},
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return &fakeTx{commitFn: func() error { return errors.New("commit error") }}, nil
			},
			deactivateTxFn: func(ctx context.Context, tx output.Tx, id string) error {
				return nil
			},
		}
		inter := NewDailyLogbookInteractor(svc)
		err := inter.DeactivateDailyLogbook(context.Background(), "lb-1")
		if err == nil {
			t.Error("expected commit error")
		}
	})
}
