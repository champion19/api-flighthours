package services

import (
	"context"
	"errors"
	"testing"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
)

// mockLicensePlateRepo implements output.LicensePlateRepository for testing
type mockLicensePlateRepo struct {
	getByIDFn func(ctx context.Context, id string) (*domain.LicensePlate, error)
	listFn    func(ctx context.Context, filters map[string]interface{}) ([]domain.LicensePlate, error)
	saveFn    func(ctx context.Context, tx output.Tx, registration domain.LicensePlate) error
	updateFn  func(ctx context.Context, tx output.Tx, registration domain.LicensePlate) error
	beginTxFn func(ctx context.Context) (output.Tx, error)
}

func (m *mockLicensePlateRepo) GetLicensePlateByID(ctx context.Context, id string) (*domain.LicensePlate, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *mockLicensePlateRepo) ListLicensePlates(ctx context.Context, filters map[string]interface{}) ([]domain.LicensePlate, error) {
	if m.listFn != nil {
		return m.listFn(ctx, filters)
	}
	return nil, nil
}

func (m *mockLicensePlateRepo) SaveLicensePlate(ctx context.Context, tx output.Tx, registration domain.LicensePlate) error {
	if m.saveFn != nil {
		return m.saveFn(ctx, tx, registration)
	}
	return nil
}

func (m *mockLicensePlateRepo) UpdateLicensePlate(ctx context.Context, tx output.Tx, registration domain.LicensePlate) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, tx, registration)
	}
	return nil
}

func (m *mockLicensePlateRepo) BeginTx(ctx context.Context) (output.Tx, error) {
	if m.beginTxFn != nil {
		return m.beginTxFn(ctx)
	}
	return &mockTx{}, nil
}

func (m *mockLicensePlateRepo) GetLicensePlateByPlate(ctx context.Context, plate string) (*domain.LicensePlate, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, plate)
	}
	return nil, nil
}

func TestNewLicensePlateService(t *testing.T) {
	t.Run("creates service", func(t *testing.T) {
		repo := &mockLicensePlateRepo{}
		service := NewLicensePlateService(repo)
		if service == nil {
			t.Error("expected non-nil service")
		}
	})
}

func TestLicensePlateService_GetByID(t *testing.T) {
	t.Run("returns registration when found", func(t *testing.T) {
		expected := &domain.LicensePlate{
			ID:           "ar-123",
			LicensePlate: "HK-5432",
		}
		repo := &mockLicensePlateRepo{
			getByIDFn: func(ctx context.Context, id string) (*domain.LicensePlate, error) {
				return expected, nil
			},
		}

		service := NewLicensePlateService(repo)
		result, err := service.GetLicensePlateByID(context.Background(), "ar-123")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.ID != "ar-123" {
			t.Errorf("expected ID 'ar-123', got %q", result.ID)
		}
	})

	t.Run("returns error when not found", func(t *testing.T) {
		repo := &mockLicensePlateRepo{
			getByIDFn: func(ctx context.Context, id string) (*domain.LicensePlate, error) {
				return nil, domain.ErrLicensePlateNotFound
			},
		}

		service := NewLicensePlateService(repo)
		_, err := service.GetLicensePlateByID(context.Background(), "non-existent")

		if err == nil {
			t.Error("expected error for non-existent registration")
		}
	})
}

func TestLicensePlateService_List(t *testing.T) {
	t.Run("returns list successfully", func(t *testing.T) {
		expected := []domain.LicensePlate{
			{ID: "ar-1", LicensePlate: "HK-5432"},
			{ID: "ar-2", LicensePlate: "CC-BFA"},
		}
		repo := &mockLicensePlateRepo{
			listFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.LicensePlate, error) {
				return expected, nil
			},
		}

		service := NewLicensePlateService(repo)
		result, err := service.ListLicensePlates(context.Background(), nil)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 2 {
			t.Errorf("expected 2 registrations, got %d", len(result))
		}
	})

	t.Run("returns error on failure", func(t *testing.T) {
		repo := &mockLicensePlateRepo{
			listFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.LicensePlate, error) {
				return nil, errors.New("database error")
			},
		}

		service := NewLicensePlateService(repo)
		_, err := service.ListLicensePlates(context.Background(), nil)

		if err == nil {
			t.Error("expected error")
		}
	})
}

func TestLicensePlateService_Create(t *testing.T) {
	t.Run("creates successfully", func(t *testing.T) {
		repo := &mockLicensePlateRepo{
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return &mockTx{}, nil
			},
			saveFn: func(ctx context.Context, tx output.Tx, registration domain.LicensePlate) error {
				return nil
			},
		}

		service := NewLicensePlateService(repo)
		err := service.CreateLicensePlate(context.Background(), domain.LicensePlate{ID: "new-ar"})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("returns error when transaction fails", func(t *testing.T) {
		repo := &mockLicensePlateRepo{
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return nil, errors.New("tx failed")
			},
		}

		service := NewLicensePlateService(repo)
		err := service.CreateLicensePlate(context.Background(), domain.LicensePlate{})

		if err == nil {
			t.Error("expected error when transaction fails")
		}
	})

	t.Run("returns error when save fails", func(t *testing.T) {
		repo := &mockLicensePlateRepo{
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return &mockTx{}, nil
			},
			saveFn: func(ctx context.Context, tx output.Tx, registration domain.LicensePlate) error {
				return errors.New("save failed")
			},
		}

		service := NewLicensePlateService(repo)
		err := service.CreateLicensePlate(context.Background(), domain.LicensePlate{})

		if err == nil {
			t.Error("expected error when save fails")
		}
	})

	t.Run("returns error when commit fails", func(t *testing.T) {
		repo := &mockLicensePlateRepo{
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return &mockTx{commitFn: func() error { return errors.New("commit failed") }}, nil
			},
			saveFn: func(ctx context.Context, tx output.Tx, registration domain.LicensePlate) error {
				return nil
			},
		}

		service := NewLicensePlateService(repo)
		err := service.CreateLicensePlate(context.Background(), domain.LicensePlate{})

		if err == nil {
			t.Error("expected error when commit fails")
		}
	})
}

func TestLicensePlateService_Update(t *testing.T) {
	t.Run("updates successfully", func(t *testing.T) {
		updateCalled := false
		repo := &mockLicensePlateRepo{
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return &mockTx{}, nil
			},
			updateFn: func(ctx context.Context, tx output.Tx, registration domain.LicensePlate) error {
				updateCalled = true
				if registration.ID != "ar-123" {
					t.Errorf("expected ID 'ar-123', got %q", registration.ID)
				}
				return nil
			},
		}

		service := NewLicensePlateService(repo)
		err := service.UpdateLicensePlate(context.Background(), domain.LicensePlate{ID: "ar-123"})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !updateCalled {
			t.Error("expected update to be called")
		}
	})

	t.Run("returns error when transaction fails", func(t *testing.T) {
		repo := &mockLicensePlateRepo{
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return nil, errors.New("tx failed")
			},
		}

		service := NewLicensePlateService(repo)
		err := service.UpdateLicensePlate(context.Background(), domain.LicensePlate{})

		if err == nil {
			t.Error("expected error when transaction fails")
		}
	})

	t.Run("returns error when update fails", func(t *testing.T) {
		repo := &mockLicensePlateRepo{
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return &mockTx{}, nil
			},
			updateFn: func(ctx context.Context, tx output.Tx, registration domain.LicensePlate) error {
				return errors.New("update failed")
			},
		}

		service := NewLicensePlateService(repo)
		err := service.UpdateLicensePlate(context.Background(), domain.LicensePlate{})

		if err == nil {
			t.Error("expected error when update fails")
		}
	})

	t.Run("returns error when commit fails", func(t *testing.T) {
		repo := &mockLicensePlateRepo{
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return &mockTx{commitFn: func() error { return errors.New("commit failed") }}, nil
			},
			updateFn: func(ctx context.Context, tx output.Tx, registration domain.LicensePlate) error {
				return nil
			},
		}

		service := NewLicensePlateService(repo)
		err := service.UpdateLicensePlate(context.Background(), domain.LicensePlate{})

		if err == nil {
			t.Error("expected error when commit fails")
		}
	})
}

func TestLicensePlateService_BeginTx(t *testing.T) {
	t.Run("begins transaction successfully", func(t *testing.T) {
		repo := &mockLicensePlateRepo{
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return &mockTx{}, nil
			},
		}

		service := NewLicensePlateService(repo)
		tx, err := service.BeginTx(context.Background())

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tx == nil {
			t.Error("expected non-nil transaction")
		}
	})
}
