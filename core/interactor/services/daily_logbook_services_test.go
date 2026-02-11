package services

import (
	"context"
	"errors"
	"testing"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/logger"
)

// mock daily logbook repository
type mockDailyLogbookRepo struct {
	getByIDFn func(ctx context.Context, id string) (*domain.DailyLogbook, error)
	listFn    func(ctx context.Context, employeeID string, filters map[string]interface{}) ([]domain.DailyLogbook, error)
	saveFn    func(ctx context.Context, tx output.Tx, logbook domain.DailyLogbook) error
	beginTxFn func(ctx context.Context) (output.Tx, error)
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

func (m *mockDailyLogbookRepo) BeginTx(ctx context.Context) (output.Tx, error) {
	if m.beginTxFn != nil {
		return m.beginTxFn(ctx)
	}
	return &mockTx{}, nil
}

// noopDailyLogbookLogger satisfies logger.Logger for tests
type noopDailyLogbookLogger struct{}

func (n noopDailyLogbookLogger) Info(msg string, args ...interface{})     {}
func (n noopDailyLogbookLogger) Debug(msg string, args ...interface{})    {}
func (n noopDailyLogbookLogger) Error(msg string, args ...interface{})    {}
func (n noopDailyLogbookLogger) Warn(msg string, args ...interface{})     {}
func (n noopDailyLogbookLogger) Success(msg string, args ...interface{})  {}
func (n noopDailyLogbookLogger) Fatal(msg string, args ...interface{})    {}
func (n noopDailyLogbookLogger) Panic(msg string, args ...interface{})    {}
func (n noopDailyLogbookLogger) WithTraceID(traceID string) logger.Logger { return n }

func TestNewDailyLogbookService(t *testing.T) {
	repo := &mockDailyLogbookRepo{}
	svc := NewDailyLogbookService(repo, noopDailyLogbookLogger{})
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
		svc := NewDailyLogbookService(repo, noopDailyLogbookLogger{})
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
		svc := NewDailyLogbookService(repo, noopDailyLogbookLogger{})
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
		svc := NewDailyLogbookService(repo, noopDailyLogbookLogger{})
		result, err := svc.ListDailyLogbooksByEmployee(context.Background(), "emp-1", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 2 {
			t.Errorf("expected 2 logbooks, got %d", len(result))
		}
	})
}

func TestDailyLogbookService_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &mockDailyLogbookRepo{
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return &mockTx{}, nil
			},
			saveFn: func(ctx context.Context, tx output.Tx, logbook domain.DailyLogbook) error {
				return nil
			},
		}
		svc := NewDailyLogbookService(repo, noopDailyLogbookLogger{})
		err := svc.CreateDailyLogbook(context.Background(), domain.DailyLogbook{EmployeeID: "emp-1"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("begin tx error", func(t *testing.T) {
		repo := &mockDailyLogbookRepo{
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return nil, errors.New("tx error")
			},
		}
		svc := NewDailyLogbookService(repo, noopDailyLogbookLogger{})
		err := svc.CreateDailyLogbook(context.Background(), domain.DailyLogbook{EmployeeID: "emp-1"})
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("save error rolls back", func(t *testing.T) {
		rollbackCalled := false
		repo := &mockDailyLogbookRepo{
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return &mockTx{
					rollbackFn: func() error { rollbackCalled = true; return nil },
				}, nil
			},
			saveFn: func(ctx context.Context, tx output.Tx, logbook domain.DailyLogbook) error {
				return errors.New("save error")
			},
		}
		svc := NewDailyLogbookService(repo, noopDailyLogbookLogger{})
		err := svc.CreateDailyLogbook(context.Background(), domain.DailyLogbook{EmployeeID: "emp-1"})
		if err == nil {
			t.Error("expected error")
		}
		if !rollbackCalled {
			t.Error("expected rollback to be called")
		}
	})
}

func TestDailyLogbookService_BeginTx(t *testing.T) {
	repo := &mockDailyLogbookRepo{
		beginTxFn: func(ctx context.Context) (output.Tx, error) {
			return &mockTx{}, nil
		},
	}
	svc := NewDailyLogbookService(repo, noopDailyLogbookLogger{})
	tx, err := svc.BeginTx(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tx == nil {
		t.Error("expected non-nil tx")
	}
}
