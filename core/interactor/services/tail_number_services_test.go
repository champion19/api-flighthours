package services

import (
	"context"
	"errors"
	"testing"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
)

// mockTailNumberRepo implements output.TailNumberRepository for testing
type mockTailNumberRepo struct {
	getByIDFn    func(ctx context.Context, id string) (*domain.TailNumber, error)
	getByPlateFn func(ctx context.Context, plate string) (*domain.TailNumber, error)
	listFn       func(ctx context.Context, filters map[string]interface{}) ([]domain.TailNumber, error)
	saveFn       func(ctx context.Context, tx output.Tx, registration domain.TailNumber) error
	updateFn     func(ctx context.Context, tx output.Tx, registration domain.TailNumber) error
	beginTxFn    func(ctx context.Context) (output.Tx, error)
}

func (m *mockTailNumberRepo) GetTailNumberByID(ctx context.Context, id string) (*domain.TailNumber, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *mockTailNumberRepo) ListTailNumbers(ctx context.Context, filters map[string]interface{}) ([]domain.TailNumber, error) {
	if m.listFn != nil {
		return m.listFn(ctx, filters)
	}
	return nil, nil
}

func (m *mockTailNumberRepo) SaveTailNumber(ctx context.Context, tx output.Tx, registration domain.TailNumber) error {
	if m.saveFn != nil {
		return m.saveFn(ctx, tx, registration)
	}
	return nil
}

func (m *mockTailNumberRepo) UpdateTailNumber(ctx context.Context, tx output.Tx, registration domain.TailNumber) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, tx, registration)
	}
	return nil
}

func (m *mockTailNumberRepo) BeginTx(ctx context.Context) (output.Tx, error) {
	if m.beginTxFn != nil {
		return m.beginTxFn(ctx)
	}
	return &mockTx{}, nil
}

func (m *mockTailNumberRepo) GetTailNumberByPlate(ctx context.Context, plate string) (*domain.TailNumber, error) {
	if m.getByPlateFn != nil {
		return m.getByPlateFn(ctx, plate)
	}
	return nil, nil
}

func TestNewTailNumberService(t *testing.T) {
	t.Run("creates service", func(t *testing.T) {
		repo := &mockTailNumberRepo{}
		service := NewTailNumberService(repo)
		if service == nil {
			t.Error("expected non-nil service")
		}
	})
}

func TestTailNumberService_GetByID(t *testing.T) {
	t.Run("returns registration when found", func(t *testing.T) {
		expected := &domain.TailNumber{
			ID:           "ar-123",
			TailNumber: "HK-5432",
		}
		repo := &mockTailNumberRepo{
			getByIDFn: func(ctx context.Context, id string) (*domain.TailNumber, error) {
				return expected, nil
			},
		}

		service := NewTailNumberService(repo)
		result, err := service.GetTailNumberByID(context.Background(), "ar-123")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.ID != "ar-123" {
			t.Errorf("expected ID 'ar-123', got %q", result.ID)
		}
	})

	t.Run("returns error when not found", func(t *testing.T) {
		repo := &mockTailNumberRepo{
			getByIDFn: func(ctx context.Context, id string) (*domain.TailNumber, error) {
				return nil, domain.ErrTailNumberNotFound
			},
		}

		service := NewTailNumberService(repo)
		_, err := service.GetTailNumberByID(context.Background(), "non-existent")

		if err == nil {
			t.Error("expected error for non-existent registration")
		}
	})
}

func TestTailNumberService_GetByPlate(t *testing.T) {
	t.Run("returns registration when found", func(t *testing.T) {
		expected := &domain.TailNumber{
			ID:           "ar-123",
			TailNumber: "HK-5432",
		}
		repo := &mockTailNumberRepo{
			getByPlateFn: func(ctx context.Context, plate string) (*domain.TailNumber, error) {
				return expected, nil
			},
		}

		service := NewTailNumberService(repo)
		result, err := service.GetTailNumberByPlate(context.Background(), "HK-5432")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.TailNumber != "HK-5432" {
			t.Errorf("expected TailNumber 'HK-5432', got %q", result.TailNumber)
		}
	})

	t.Run("returns error when not found", func(t *testing.T) {
		repo := &mockTailNumberRepo{
			getByPlateFn: func(ctx context.Context, plate string) (*domain.TailNumber, error) {
				return nil, domain.ErrTailNumberNotFound
			},
		}

		service := NewTailNumberService(repo)
		_, err := service.GetTailNumberByPlate(context.Background(), "INVALID")

		if err == nil {
			t.Error("expected error for non-existent registration")
		}
	})
}

func TestTailNumberService_List(t *testing.T) {
	t.Run("returns list successfully", func(t *testing.T) {
		expected := []domain.TailNumber{
			{ID: "ar-1", TailNumber: "HK-5432"},
			{ID: "ar-2", TailNumber: "CC-BFA"},
		}
		repo := &mockTailNumberRepo{
			listFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.TailNumber, error) {
				return expected, nil
			},
		}

		service := NewTailNumberService(repo)
		result, err := service.ListTailNumbers(context.Background(), nil)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 2 {
			t.Errorf("expected 2 registrations, got %d", len(result))
		}
	})

	t.Run("returns error on failure", func(t *testing.T) {
		repo := &mockTailNumberRepo{
			listFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.TailNumber, error) {
				return nil, errors.New("database error")
			},
		}

		service := NewTailNumberService(repo)
		_, err := service.ListTailNumbers(context.Background(), nil)

		if err == nil {
			t.Error("expected error")
		}
	})
}

func TestTailNumberService_BeginTx(t *testing.T) {
	t.Run("begins transaction successfully", func(t *testing.T) {
		repo := &mockTailNumberRepo{
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return &mockTx{}, nil
			},
		}

		service := NewTailNumberService(repo)
		tx, err := service.BeginTx(context.Background())

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tx == nil {
			t.Error("expected non-nil transaction")
		}
	})
}

func TestTailNumberService_CreateTailNumberTx(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &mockTailNumberRepo{
			saveFn: func(ctx context.Context, tx output.Tx, registration domain.TailNumber) error {
				return nil
			},
		}
		service := NewTailNumberService(repo)
		err := service.CreateTailNumberTx(context.Background(), &mockTx{}, domain.TailNumber{ID: "lp-1"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		repo := &mockTailNumberRepo{
			saveFn: func(ctx context.Context, tx output.Tx, registration domain.TailNumber) error {
				return errors.New("save failed")
			},
		}
		service := NewTailNumberService(repo)
		err := service.CreateTailNumberTx(context.Background(), &mockTx{}, domain.TailNumber{})
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestTailNumberService_UpdateTailNumberTx(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &mockTailNumberRepo{
			updateFn: func(ctx context.Context, tx output.Tx, registration domain.TailNumber) error {
				return nil
			},
		}
		service := NewTailNumberService(repo)
		err := service.UpdateTailNumberTx(context.Background(), &mockTx{}, domain.TailNumber{ID: "lp-1"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		repo := &mockTailNumberRepo{
			updateFn: func(ctx context.Context, tx output.Tx, registration domain.TailNumber) error {
				return errors.New("update failed")
			},
		}
		service := NewTailNumberService(repo)
		err := service.UpdateTailNumberTx(context.Background(), &mockTx{}, domain.TailNumber{})
		if err == nil {
			t.Fatal("expected error")
		}
	})
}
