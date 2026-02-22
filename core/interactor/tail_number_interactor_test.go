package interactor

import (
	"context"
	"errors"
	"testing"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/input"
	"github.com/champion19/api-flighthours/core/ports/output"
)

// fakeTailNumberService implements input.TailNumberService for testing
type fakeTailNumberService struct {
	getByIDFn    func(ctx context.Context, id string) (*domain.TailNumber, error)
	getByPlateFn func(ctx context.Context, plate string) (*domain.TailNumber, error)
	listFn       func(ctx context.Context, filters map[string]interface{}) ([]domain.TailNumber, error)
	createFn     func(ctx context.Context, registration domain.TailNumber) error
	updateFn     func(ctx context.Context, registration domain.TailNumber) error
	createTxFn   func(ctx context.Context, tx output.Tx, registration domain.TailNumber) error
	updateTxFn   func(ctx context.Context, tx output.Tx, registration domain.TailNumber) error
}

var _ input.TailNumberService = (*fakeTailNumberService)(nil)

func (f *fakeTailNumberService) BeginTx(ctx context.Context) (output.Tx, error) {
	return &fakeTx{}, nil
}

func (f *fakeTailNumberService) GetTailNumberByID(ctx context.Context, id string) (*domain.TailNumber, error) {
	if f.getByIDFn != nil {
		return f.getByIDFn(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (f *fakeTailNumberService) ListTailNumbers(ctx context.Context, filters map[string]interface{}) ([]domain.TailNumber, error) {
	if f.listFn != nil {
		return f.listFn(ctx, filters)
	}
	return nil, errors.New("not implemented")
}

func (f *fakeTailNumberService) CreateTailNumber(ctx context.Context, registration domain.TailNumber) error {
	if f.createFn != nil {
		return f.createFn(ctx, registration)
	}
	return nil
}

func (f *fakeTailNumberService) UpdateTailNumber(ctx context.Context, registration domain.TailNumber) error {
	if f.updateFn != nil {
		return f.updateFn(ctx, registration)
	}
	return nil
}

func (f *fakeTailNumberService) GetTailNumberByPlate(ctx context.Context, plate string) (*domain.TailNumber, error) {
	if f.getByPlateFn != nil {
		return f.getByPlateFn(ctx, plate)
	}
	return nil, errors.New("not implemented")
}

func (f *fakeTailNumberService) CreateTailNumberTx(ctx context.Context, tx output.Tx, registration domain.TailNumber) error {
	if f.createTxFn != nil {
		return f.createTxFn(ctx, tx, registration)
	}
	return nil
}

func (f *fakeTailNumberService) UpdateTailNumberTx(ctx context.Context, tx output.Tx, registration domain.TailNumber) error {
	if f.updateTxFn != nil {
		return f.updateTxFn(ctx, tx, registration)
	}
	return nil
}

func TestNewTailNumberInteractor(t *testing.T) {
	svc := &fakeTailNumberService{}
	interactor := NewTailNumberInteractor(svc, noopLogger{})
	if interactor == nil {
		t.Error("expected non-nil TailNumberInteractor")
	}
}

func TestTailNumberInteractor_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expected := &domain.TailNumber{
			ID:          "ar-123",
			TailNumber:  "HK-5432",
			ModelName:   "Boeing 737",
			AirlineName: "Avianca",
		}
		svc := &fakeTailNumberService{
			getByIDFn: func(ctx context.Context, id string) (*domain.TailNumber, error) {
				return expected, nil
			},
		}
		interactor := NewTailNumberInteractor(svc, noopLogger{})

		result, err := interactor.GetTailNumberByID(context.Background(), "ar-123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.ID != expected.ID {
			t.Errorf("expected ID %q, got %q", expected.ID, result.ID)
		}
		if result.TailNumber != expected.TailNumber {
			t.Errorf("expected TailNumber %q, got %q", expected.TailNumber, result.TailNumber)
		}
	})

	t.Run("service error", func(t *testing.T) {
		svc := &fakeTailNumberService{
			getByIDFn: func(ctx context.Context, id string) (*domain.TailNumber, error) {
				return nil, errors.New("database error")
			},
		}
		interactor := NewTailNumberInteractor(svc, noopLogger{})

		_, err := interactor.GetTailNumberByID(context.Background(), "ar-123")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestTailNumberInteractor_GetByPlate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expected := &domain.TailNumber{
			ID:          "ar-123",
			TailNumber:  "HK-5432",
			ModelName:   "Boeing 737",
			AirlineName: "Avianca",
		}
		svc := &fakeTailNumberService{
			getByPlateFn: func(ctx context.Context, plate string) (*domain.TailNumber, error) {
				return expected, nil
			},
		}
		interactor := NewTailNumberInteractor(svc, noopLogger{})

		result, err := interactor.GetTailNumberByPlate(context.Background(), "HK-5432")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.TailNumber != expected.TailNumber {
			t.Errorf("expected TailNumber %q, got %q", expected.TailNumber, result.TailNumber)
		}
	})

	t.Run("service error", func(t *testing.T) {
		svc := &fakeTailNumberService{
			getByPlateFn: func(ctx context.Context, plate string) (*domain.TailNumber, error) {
				return nil, errors.New("database error")
			},
		}
		interactor := NewTailNumberInteractor(svc, noopLogger{})

		_, err := interactor.GetTailNumberByPlate(context.Background(), "HK-5432")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestTailNumberInteractor_List(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expected := []domain.TailNumber{
			{ID: "ar-1", TailNumber: "HK-5432"},
			{ID: "ar-2", TailNumber: "CC-BFA"},
		}
		svc := &fakeTailNumberService{
			listFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.TailNumber, error) {
				return expected, nil
			},
		}
		interactor := NewTailNumberInteractor(svc, noopLogger{})

		result, err := interactor.ListTailNumbers(context.Background(), nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 2 {
			t.Errorf("expected 2 registrations, got %d", len(result))
		}
	})

	t.Run("service error", func(t *testing.T) {
		svc := &fakeTailNumberService{
			listFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.TailNumber, error) {
				return nil, errors.New("database error")
			},
		}
		interactor := NewTailNumberInteractor(svc, noopLogger{})

		_, err := interactor.ListTailNumbers(context.Background(), nil)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestTailNumberInteractor_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &fakeTailNumberService{
			createTxFn: func(ctx context.Context, tx output.Tx, registration domain.TailNumber) error {
				if registration.TailNumber != "HK-5432" {
					t.Errorf("expected TailNumber HK-5432, got %s", registration.TailNumber)
				}
				return nil
			},
		}
		interactor := NewTailNumberInteractor(svc, noopLogger{})

		registration := domain.TailNumber{
			ID:              "new-uuid",
			TailNumber:      "HK-5432",
			AircraftModelID: "model-uuid",
			AirlineID:       "airline-uuid",
		}
		err := interactor.CreateTailNumber(context.Background(), registration)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("service error", func(t *testing.T) {
		svc := &fakeTailNumberService{
			createTxFn: func(ctx context.Context, tx output.Tx, registration domain.TailNumber) error {
				return errors.New("create failed")
			},
		}
		interactor := NewTailNumberInteractor(svc, noopLogger{})

		err := interactor.CreateTailNumber(context.Background(), domain.TailNumber{})
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestTailNumberInteractor_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &fakeTailNumberService{
			updateTxFn: func(ctx context.Context, tx output.Tx, registration domain.TailNumber) error {
				return nil
			},
		}
		interactor := NewTailNumberInteractor(svc, noopLogger{})

		err := interactor.UpdateTailNumber(context.Background(), domain.TailNumber{
			ID:              "ar-123",
			TailNumber:      "HK-9999",
			AircraftModelID: "model-uuid",
			AirlineID:       "airline-uuid",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		svc := &fakeTailNumberService{
			updateTxFn: func(ctx context.Context, tx output.Tx, registration domain.TailNumber) error {
				return domain.ErrTailNumberNotFound
			},
		}
		interactor := NewTailNumberInteractor(svc, noopLogger{})

		err := interactor.UpdateTailNumber(context.Background(), domain.TailNumber{ID: "non-existent"})
		if err != domain.ErrTailNumberNotFound {
			t.Errorf("expected ErrTailNumberNotFound, got %v", err)
		}
	})

	t.Run("update service error", func(t *testing.T) {
		svc := &fakeTailNumberService{
			updateTxFn: func(ctx context.Context, tx output.Tx, registration domain.TailNumber) error {
				return errors.New("update failed")
			},
		}
		interactor := NewTailNumberInteractor(svc, noopLogger{})

		err := interactor.UpdateTailNumber(context.Background(), domain.TailNumber{ID: "ar-123"})
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
