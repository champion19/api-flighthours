package interactor

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/input"
	"github.com/champion19/api-flighthours/core/ports/output"
)

// fakeAirlineEmployeeService implements input.AirlineEmployeeService for testing
type fakeAirlineEmployeeService struct {
	getByIDFn      func(ctx context.Context, id string) (*domain.AirlineEmployee, error)
	addFn          func(ctx context.Context, employee domain.AirlineEmployee) error
	updateFn       func(ctx context.Context, employee domain.AirlineEmployee) error
	activateFn     func(ctx context.Context, id string) error
	deactivateFn   func(ctx context.Context, id string) error
	addTxFn        func(ctx context.Context, tx output.Tx, employee domain.AirlineEmployee) error
	updateTxFn     func(ctx context.Context, tx output.Tx, employee domain.AirlineEmployee) error
	activateTxFn   func(ctx context.Context, tx output.Tx, id string) error
	deactivateTxFn func(ctx context.Context, tx output.Tx, id string) error
}

var _ input.AirlineEmployeeService = (*fakeAirlineEmployeeService)(nil)

func (f *fakeAirlineEmployeeService) BeginTx(ctx context.Context) (output.Tx, error) {
	return &fakeTx{}, nil
}

func (f *fakeAirlineEmployeeService) GetAirlineEmployeeByID(ctx context.Context, id string) (*domain.AirlineEmployee, error) {
	if f.getByIDFn != nil {
		return f.getByIDFn(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (f *fakeAirlineEmployeeService) AddAirlineEmployee(ctx context.Context, employee domain.AirlineEmployee) error {
	if f.addFn != nil {
		return f.addFn(ctx, employee)
	}
	return nil
}

func (f *fakeAirlineEmployeeService) UpdateAirlineEmployee(ctx context.Context, employee domain.AirlineEmployee) error {
	if f.updateFn != nil {
		return f.updateFn(ctx, employee)
	}
	return nil
}

func (f *fakeAirlineEmployeeService) ActivateAirlineEmployee(ctx context.Context, id string) error {
	if f.activateFn != nil {
		return f.activateFn(ctx, id)
	}
	return nil
}

func (f *fakeAirlineEmployeeService) DeactivateAirlineEmployee(ctx context.Context, id string) error {
	if f.deactivateFn != nil {
		return f.deactivateFn(ctx, id)
	}
	return nil
}

func (f *fakeAirlineEmployeeService) AddAirlineEmployeeTx(ctx context.Context, tx output.Tx, employee domain.AirlineEmployee) error {
	if f.addTxFn != nil {
		return f.addTxFn(ctx, tx, employee)
	}
	return nil
}

func (f *fakeAirlineEmployeeService) UpdateAirlineEmployeeTx(ctx context.Context, tx output.Tx, employee domain.AirlineEmployee) error {
	if f.updateTxFn != nil {
		return f.updateTxFn(ctx, tx, employee)
	}
	return nil
}

func (f *fakeAirlineEmployeeService) ActivateAirlineEmployeeTx(ctx context.Context, tx output.Tx, id string) error {
	if f.activateTxFn != nil {
		return f.activateTxFn(ctx, tx, id)
	}
	return nil
}

func (f *fakeAirlineEmployeeService) DeactivateAirlineEmployeeTx(ctx context.Context, tx output.Tx, id string) error {
	if f.deactivateTxFn != nil {
		return f.deactivateTxFn(ctx, tx, id)
	}
	return nil
}

func TestNewAirlineEmployeeInteractor(t *testing.T) {
	svc := &fakeAirlineEmployeeService{}
	interactor := NewAirlineEmployeeInteractor(svc)
	if interactor == nil {
		t.Error("expected non-nil AirlineEmployeeInteractor")
	}
}

func TestAirlineEmployeeInteractor_GetAirlineEmployeeByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expected := &domain.AirlineEmployee{
			ID:        "emp-123",
			AirlineID: "airline-456",
			Bp:        "BP12345",
			Active:    true,
		}
		svc := &fakeAirlineEmployeeService{
			getByIDFn: func(ctx context.Context, id string) (*domain.AirlineEmployee, error) {
				return expected, nil
			},
		}
		interactor := NewAirlineEmployeeInteractor(svc)

		result, err := interactor.GetAirlineEmployeeByID(context.Background(), "emp-123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.ID != expected.ID {
			t.Errorf("expected ID %q, got %q", expected.ID, result.ID)
		}
	})

	t.Run("service error", func(t *testing.T) {
		svc := &fakeAirlineEmployeeService{
			getByIDFn: func(ctx context.Context, id string) (*domain.AirlineEmployee, error) {
				return nil, errors.New("database error")
			},
		}
		interactor := NewAirlineEmployeeInteractor(svc)

		_, err := interactor.GetAirlineEmployeeByID(context.Background(), "emp-123")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestAirlineEmployeeInteractor_AddAirlineEmployee(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &fakeAirlineEmployeeService{
			addTxFn: func(ctx context.Context, tx output.Tx, employee domain.AirlineEmployee) error {
				if employee.ID != "emp-123" {
					t.Errorf("expected ID emp-123, got %s", employee.ID)
				}
				return nil
			},
		}
		interactor := NewAirlineEmployeeInteractor(svc)

		info := domain.AirlineEmployee{
			AirlineID: "airline-456",
			Bp:        "BP12345",
			StartDate: time.Now(),
		}
		err := interactor.AddAirlineEmployee(context.Background(), "emp-123", info)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("service error", func(t *testing.T) {
		svc := &fakeAirlineEmployeeService{
			addTxFn: func(ctx context.Context, tx output.Tx, employee domain.AirlineEmployee) error {
				return errors.New("add failed")
			},
		}
		interactor := NewAirlineEmployeeInteractor(svc)

		err := interactor.AddAirlineEmployee(context.Background(), "emp-123", domain.AirlineEmployee{})
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestAirlineEmployeeInteractor_UpdateAirlineEmployee(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &fakeAirlineEmployeeService{
			updateTxFn: func(ctx context.Context, tx output.Tx, employee domain.AirlineEmployee) error {
				return nil
			},
		}
		interactor := NewAirlineEmployeeInteractor(svc)

		info := domain.AirlineEmployee{AirlineID: "airline-456"}
		err := interactor.UpdateAirlineEmployee(context.Background(), "emp-123", info)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("service error", func(t *testing.T) {
		svc := &fakeAirlineEmployeeService{
			updateTxFn: func(ctx context.Context, tx output.Tx, employee domain.AirlineEmployee) error {
				return errors.New("update failed")
			},
		}
		interactor := NewAirlineEmployeeInteractor(svc)

		err := interactor.UpdateAirlineEmployee(context.Background(), "emp-123", domain.AirlineEmployee{})
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestAirlineEmployeeInteractor_ActivateAirlineEmployee(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &fakeAirlineEmployeeService{
			getByIDFn: func(ctx context.Context, id string) (*domain.AirlineEmployee, error) {
				return &domain.AirlineEmployee{ID: id, AirlineID: "airline-123"}, nil
			},
			activateTxFn: func(ctx context.Context, tx output.Tx, id string) error {
				return nil
			},
		}
		interactor := NewAirlineEmployeeInteractor(svc)

		err := interactor.ActivateAirlineEmployee(context.Background(), "emp-123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("not found - no airline info", func(t *testing.T) {
		svc := &fakeAirlineEmployeeService{
			getByIDFn: func(ctx context.Context, id string) (*domain.AirlineEmployee, error) {
				return nil, domain.ErrAirlineEmployeeNotFound
			},
		}
		interactor := NewAirlineEmployeeInteractor(svc)

		err := interactor.ActivateAirlineEmployee(context.Background(), "emp-123")
		if err != domain.ErrAirlineEmployeeNotFound {
			t.Errorf("expected ErrAirlineEmployeeNotFound, got %v", err)
		}
	})

	t.Run("activation service error", func(t *testing.T) {
		svc := &fakeAirlineEmployeeService{
			getByIDFn: func(ctx context.Context, id string) (*domain.AirlineEmployee, error) {
				return &domain.AirlineEmployee{ID: id, AirlineID: "airline-123"}, nil
			},
			activateTxFn: func(ctx context.Context, tx output.Tx, id string) error {
				return errors.New("activation failed")
			},
		}
		interactor := NewAirlineEmployeeInteractor(svc)

		err := interactor.ActivateAirlineEmployee(context.Background(), "emp-123")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestAirlineEmployeeInteractor_DeactivateAirlineEmployee(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &fakeAirlineEmployeeService{
			getByIDFn: func(ctx context.Context, id string) (*domain.AirlineEmployee, error) {
				return &domain.AirlineEmployee{ID: id, AirlineID: "airline-123"}, nil
			},
			deactivateTxFn: func(ctx context.Context, tx output.Tx, id string) error {
				return nil
			},
		}
		interactor := NewAirlineEmployeeInteractor(svc)

		err := interactor.DeactivateAirlineEmployee(context.Background(), "emp-123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("not found - no airline info", func(t *testing.T) {
		svc := &fakeAirlineEmployeeService{
			getByIDFn: func(ctx context.Context, id string) (*domain.AirlineEmployee, error) {
				return nil, domain.ErrAirlineEmployeeNotFound
			},
		}
		interactor := NewAirlineEmployeeInteractor(svc)

		err := interactor.DeactivateAirlineEmployee(context.Background(), "emp-123")
		if err != domain.ErrAirlineEmployeeNotFound {
			t.Errorf("expected ErrAirlineEmployeeNotFound, got %v", err)
		}
	})

	t.Run("deactivation service error", func(t *testing.T) {
		svc := &fakeAirlineEmployeeService{
			getByIDFn: func(ctx context.Context, id string) (*domain.AirlineEmployee, error) {
				return &domain.AirlineEmployee{ID: id, AirlineID: "airline-123"}, nil
			},
			deactivateTxFn: func(ctx context.Context, tx output.Tx, id string) error {
				return errors.New("deactivation failed")
			},
		}
		interactor := NewAirlineEmployeeInteractor(svc)

		err := interactor.DeactivateAirlineEmployee(context.Background(), "emp-123")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
