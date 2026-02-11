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
	getByIDFn func(ctx context.Context, id string) (*domain.DailyLogbook, error)
	listFn    func(ctx context.Context, employeeID string, filters map[string]interface{}) ([]domain.DailyLogbook, error)
	createFn  func(ctx context.Context, logbook domain.DailyLogbook) error
}

func (f *fakeDailyLogbookService) BeginTx(ctx context.Context) (output.Tx, error) {
	return nil, nil
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

func TestNewDailyLogbookInteractor(t *testing.T) {
	svc := &fakeDailyLogbookService{}
	inter := NewDailyLogbookInteractor(svc, noopLogger{})
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
		inter := NewDailyLogbookInteractor(svc, noopLogger{})
		result, err := inter.GetDailyLogbookByID(context.Background(), "lb-1")
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
		inter := NewDailyLogbookInteractor(svc, noopLogger{})
		_, err := inter.GetDailyLogbookByID(context.Background(), "lb-1")
		if err == nil {
			t.Error("expected error")
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
		inter := NewDailyLogbookInteractor(svc, noopLogger{})
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
		inter := NewDailyLogbookInteractor(svc, noopLogger{})
		_, err := inter.ListDailyLogbooksByEmployee(context.Background(), "emp-1", nil)
		if err == nil {
			t.Error("expected error")
		}
	})
}

func TestDailyLogbookInteractor_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &fakeDailyLogbookService{
			createFn: func(ctx context.Context, logbook domain.DailyLogbook) error {
				return nil
			},
		}
		inter := NewDailyLogbookInteractor(svc, noopLogger{})
		err := inter.CreateDailyLogbook(context.Background(), domain.DailyLogbook{EmployeeID: "emp-1"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		svc := &fakeDailyLogbookService{
			createFn: func(ctx context.Context, logbook domain.DailyLogbook) error {
				return errors.New("create error")
			},
		}
		inter := NewDailyLogbookInteractor(svc, noopLogger{})
		err := inter.CreateDailyLogbook(context.Background(), domain.DailyLogbook{EmployeeID: "emp-1"})
		if err == nil {
			t.Error("expected error")
		}
	})
}
