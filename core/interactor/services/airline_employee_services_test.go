package services

import (
	"context"
	"errors"
	"testing"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
)

// Mock airline employee repository
type mockAirlineEmployeeRepo struct {
	getByIDFn      func(ctx context.Context, id string) (*domain.AirlineEmployee, error)
	beginTxFn      func(ctx context.Context) (output.Tx, error)
	addFn          func(ctx context.Context, tx output.Tx, employee domain.AirlineEmployee) error
	updateFn       func(ctx context.Context, tx output.Tx, employee domain.AirlineEmployee) error
	updateStatusFn func(ctx context.Context, tx output.Tx, id string, active bool) error
}

func (m *mockAirlineEmployeeRepo) GetAirlineEmployeeByID(ctx context.Context, id string) (*domain.AirlineEmployee, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *mockAirlineEmployeeRepo) BeginTx(ctx context.Context) (output.Tx, error) {
	if m.beginTxFn != nil {
		return m.beginTxFn(ctx)
	}
	return &mockTx{}, nil
}

func (m *mockAirlineEmployeeRepo) AddAirlineEmployee(ctx context.Context, tx output.Tx, employee domain.AirlineEmployee) error {
	if m.addFn != nil {
		return m.addFn(ctx, tx, employee)
	}
	return nil
}

func (m *mockAirlineEmployeeRepo) UpdateAirlineEmployee(ctx context.Context, tx output.Tx, employee domain.AirlineEmployee) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, tx, employee)
	}
	return nil
}

func (m *mockAirlineEmployeeRepo) UpdateAirlineEmployeeStatus(ctx context.Context, tx output.Tx, id string, active bool) error {
	if m.updateStatusFn != nil {
		return m.updateStatusFn(ctx, tx, id, active)
	}
	return nil
}

func TestNewAirlineEmployeeService(t *testing.T) {
	t.Run("creates airline employee service", func(t *testing.T) {
		repo := &mockAirlineEmployeeRepo{}
		service := NewAirlineEmployeeService(repo)

		if service == nil {
			t.Error("expected non-nil service")
		}
	})
}

func TestAirlineEmployeeService_GetAirlineEmployeeByID(t *testing.T) {
	t.Run("returns employee when found", func(t *testing.T) {
		expected := &domain.AirlineEmployee{
			ID:        "ae-123",
			AirlineID: "airline-1",
		}

		repo := &mockAirlineEmployeeRepo{
			getByIDFn: func(ctx context.Context, id string) (*domain.AirlineEmployee, error) {
				return expected, nil
			},
		}

		service := NewAirlineEmployeeService(repo)
		result, err := service.GetAirlineEmployeeByID(context.Background(), "ae-123")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.ID != "ae-123" {
			t.Errorf("expected ID 'ae-123', got %q", result.ID)
		}
	})

	t.Run("returns error when not found", func(t *testing.T) {
		repo := &mockAirlineEmployeeRepo{
			getByIDFn: func(ctx context.Context, id string) (*domain.AirlineEmployee, error) {
				return nil, errors.New("not found")
			},
		}

		service := NewAirlineEmployeeService(repo)
		_, err := service.GetAirlineEmployeeByID(context.Background(), "non-existent")

		if err == nil {
			t.Error("expected error for non-existent employee")
		}
	})
}

func TestAirlineEmployeeService_AddAirlineEmployee(t *testing.T) {
	t.Run("adds employee successfully", func(t *testing.T) {
		repo := &mockAirlineEmployeeRepo{
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return &mockTx{}, nil
			},
			addFn: func(ctx context.Context, tx output.Tx, employee domain.AirlineEmployee) error {
				return nil
			},
		}

		service := NewAirlineEmployeeService(repo)
		err := service.AddAirlineEmployee(context.Background(), domain.AirlineEmployee{ID: "new-ae"})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("returns error when transaction fails", func(t *testing.T) {
		repo := &mockAirlineEmployeeRepo{
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return nil, errors.New("tx failed")
			},
		}

		service := NewAirlineEmployeeService(repo)
		err := service.AddAirlineEmployee(context.Background(), domain.AirlineEmployee{})

		if err == nil {
			t.Error("expected error when transaction fails")
		}
	})

	t.Run("returns error when add fails", func(t *testing.T) {
		repo := &mockAirlineEmployeeRepo{
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return &mockTx{}, nil
			},
			addFn: func(ctx context.Context, tx output.Tx, employee domain.AirlineEmployee) error {
				return errors.New("add failed")
			},
		}

		service := NewAirlineEmployeeService(repo)
		err := service.AddAirlineEmployee(context.Background(), domain.AirlineEmployee{})

		if err == nil {
			t.Error("expected error when add fails")
		}
	})
}

func TestAirlineEmployeeService_UpdateAirlineEmployee(t *testing.T) {
	t.Run("updates employee successfully", func(t *testing.T) {
		updateCalled := false
		repo := &mockAirlineEmployeeRepo{
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return &mockTx{}, nil
			},
			updateFn: func(ctx context.Context, tx output.Tx, employee domain.AirlineEmployee) error {
				updateCalled = true
				if employee.ID != "ae-123" {
					t.Errorf("expected ID 'ae-123', got %q", employee.ID)
				}
				return nil
			},
		}

		service := NewAirlineEmployeeService(repo)
		err := service.UpdateAirlineEmployee(context.Background(), domain.AirlineEmployee{ID: "ae-123"})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !updateCalled {
			t.Error("expected update to be called")
		}
	})

	t.Run("returns error when transaction fails", func(t *testing.T) {
		repo := &mockAirlineEmployeeRepo{
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return nil, errors.New("tx failed")
			},
		}

		service := NewAirlineEmployeeService(repo)
		err := service.UpdateAirlineEmployee(context.Background(), domain.AirlineEmployee{})

		if err == nil {
			t.Error("expected error when transaction fails")
		}
	})

	t.Run("returns error and rolls back when update fails", func(t *testing.T) {
		repo := &mockAirlineEmployeeRepo{
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return &mockTx{}, nil
			},
			updateFn: func(ctx context.Context, tx output.Tx, employee domain.AirlineEmployee) error {
				return errors.New("update failed")
			},
		}

		service := NewAirlineEmployeeService(repo)
		err := service.UpdateAirlineEmployee(context.Background(), domain.AirlineEmployee{})

		if err == nil {
			t.Error("expected error when update fails")
		}
	})
}

func TestAirlineEmployeeService_ActivateAirlineEmployee(t *testing.T) {
	t.Run("activates employee successfully", func(t *testing.T) {
		repo := &mockAirlineEmployeeRepo{
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return &mockTx{}, nil
			},
			updateStatusFn: func(ctx context.Context, tx output.Tx, id string, active bool) error {
				if !active {
					t.Error("expected active to be true for activation")
				}
				return nil
			},
		}

		service := NewAirlineEmployeeService(repo)
		err := service.ActivateAirlineEmployee(context.Background(), "ae-123")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("returns error when transaction fails", func(t *testing.T) {
		repo := &mockAirlineEmployeeRepo{
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return nil, errors.New("tx failed")
			},
		}

		service := NewAirlineEmployeeService(repo)
		err := service.ActivateAirlineEmployee(context.Background(), "ae-123")

		if err == nil {
			t.Error("expected error when transaction fails")
		}
	})

	t.Run("returns error and rolls back when update fails", func(t *testing.T) {
		repo := &mockAirlineEmployeeRepo{
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return &mockTx{}, nil
			},
			updateStatusFn: func(ctx context.Context, tx output.Tx, id string, active bool) error {
				return errors.New("update failed")
			},
		}

		service := NewAirlineEmployeeService(repo)
		err := service.ActivateAirlineEmployee(context.Background(), "ae-123")

		if err == nil {
			t.Error("expected error when update fails")
		}
	})

	t.Run("returns error when commit fails", func(t *testing.T) {
		repo := &mockAirlineEmployeeRepo{
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return &mockTx{commitFn: func() error { return errors.New("commit failed") }}, nil
			},
			updateStatusFn: func(ctx context.Context, tx output.Tx, id string, active bool) error {
				return nil
			},
		}

		service := NewAirlineEmployeeService(repo)
		err := service.ActivateAirlineEmployee(context.Background(), "ae-123")

		if err == nil {
			t.Error("expected error when commit fails")
		}
	})
}

func TestAirlineEmployeeService_DeactivateAirlineEmployee(t *testing.T) {
	t.Run("deactivates employee successfully", func(t *testing.T) {
		repo := &mockAirlineEmployeeRepo{
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return &mockTx{}, nil
			},
			updateStatusFn: func(ctx context.Context, tx output.Tx, id string, active bool) error {
				if active {
					t.Error("expected active to be false for deactivation")
				}
				return nil
			},
		}

		service := NewAirlineEmployeeService(repo)
		err := service.DeactivateAirlineEmployee(context.Background(), "ae-123")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("returns error when transaction fails", func(t *testing.T) {
		repo := &mockAirlineEmployeeRepo{
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return nil, errors.New("tx failed")
			},
		}

		service := NewAirlineEmployeeService(repo)
		err := service.DeactivateAirlineEmployee(context.Background(), "ae-123")

		if err == nil {
			t.Error("expected error when transaction fails")
		}
	})

	t.Run("returns error and rolls back when update fails", func(t *testing.T) {
		repo := &mockAirlineEmployeeRepo{
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return &mockTx{}, nil
			},
			updateStatusFn: func(ctx context.Context, tx output.Tx, id string, active bool) error {
				return errors.New("update failed")
			},
		}

		service := NewAirlineEmployeeService(repo)
		err := service.DeactivateAirlineEmployee(context.Background(), "ae-123")

		if err == nil {
			t.Error("expected error when update fails")
		}
	})

	t.Run("returns error when commit fails", func(t *testing.T) {
		repo := &mockAirlineEmployeeRepo{
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return &mockTx{commitFn: func() error { return errors.New("commit failed") }}, nil
			},
			updateStatusFn: func(ctx context.Context, tx output.Tx, id string, active bool) error {
				return nil
			},
		}

		service := NewAirlineEmployeeService(repo)
		err := service.DeactivateAirlineEmployee(context.Background(), "ae-123")

		if err == nil {
			t.Error("expected error when commit fails")
		}
	})
}

func TestAirlineEmployeeService_BeginTx(t *testing.T) {
	t.Run("begins transaction successfully", func(t *testing.T) {
		repo := &mockAirlineEmployeeRepo{
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return &mockTx{}, nil
			},
		}

		service := NewAirlineEmployeeService(repo)
		tx, err := service.BeginTx(context.Background())

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tx == nil {
			t.Error("expected non-nil transaction")
		}
	})
}
