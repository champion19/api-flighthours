package interactor

import (
	"context"
	"errors"
	"testing"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/input"
	"github.com/champion19/api-flighthours/core/ports/output"
)

// fakeLicensePlateService implements input.LicensePlateService for testing
type fakeLicensePlateService struct {
	getByIDFn    func(ctx context.Context, id string) (*domain.LicensePlate, error)
	getByPlateFn func(ctx context.Context, plate string) (*domain.LicensePlate, error)
	listFn       func(ctx context.Context, filters map[string]interface{}) ([]domain.LicensePlate, error)
	createFn     func(ctx context.Context, registration domain.LicensePlate) error
	updateFn     func(ctx context.Context, registration domain.LicensePlate) error
}

var _ input.LicensePlateService = (*fakeLicensePlateService)(nil)

func (f *fakeLicensePlateService) BeginTx(ctx context.Context) (output.Tx, error) {
	return &fakeTx{}, nil
}

func (f *fakeLicensePlateService) GetLicensePlateByID(ctx context.Context, id string) (*domain.LicensePlate, error) {
	if f.getByIDFn != nil {
		return f.getByIDFn(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (f *fakeLicensePlateService) ListLicensePlates(ctx context.Context, filters map[string]interface{}) ([]domain.LicensePlate, error) {
	if f.listFn != nil {
		return f.listFn(ctx, filters)
	}
	return nil, errors.New("not implemented")
}

func (f *fakeLicensePlateService) CreateLicensePlate(ctx context.Context, registration domain.LicensePlate) error {
	if f.createFn != nil {
		return f.createFn(ctx, registration)
	}
	return nil
}

func (f *fakeLicensePlateService) UpdateLicensePlate(ctx context.Context, registration domain.LicensePlate) error {
	if f.updateFn != nil {
		return f.updateFn(ctx, registration)
	}
	return nil
}

func (f *fakeLicensePlateService) GetLicensePlateByPlate(ctx context.Context, plate string) (*domain.LicensePlate, error) {
	if f.getByPlateFn != nil {
		return f.getByPlateFn(ctx, plate)
	}
	return nil, errors.New("not implemented")
}

func TestNewLicensePlateInteractor(t *testing.T) {
	svc := &fakeLicensePlateService{}
	interactor := NewLicensePlateInteractor(svc, noopLogger{})
	if interactor == nil {
		t.Error("expected non-nil LicensePlateInteractor")
	}
}

func TestLicensePlateInteractor_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expected := &domain.LicensePlate{
			ID:           "ar-123",
			LicensePlate: "HK-5432",
			ModelName:    "Boeing 737",
			AirlineName:  "Avianca",
		}
		svc := &fakeLicensePlateService{
			getByIDFn: func(ctx context.Context, id string) (*domain.LicensePlate, error) {
				return expected, nil
			},
		}
		interactor := NewLicensePlateInteractor(svc, noopLogger{})

		result, err := interactor.GetLicensePlateByID(context.Background(), "ar-123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.ID != expected.ID {
			t.Errorf("expected ID %q, got %q", expected.ID, result.ID)
		}
		if result.LicensePlate != expected.LicensePlate {
			t.Errorf("expected LicensePlate %q, got %q", expected.LicensePlate, result.LicensePlate)
		}
	})

	t.Run("service error", func(t *testing.T) {
		svc := &fakeLicensePlateService{
			getByIDFn: func(ctx context.Context, id string) (*domain.LicensePlate, error) {
				return nil, errors.New("database error")
			},
		}
		interactor := NewLicensePlateInteractor(svc, noopLogger{})

		_, err := interactor.GetLicensePlateByID(context.Background(), "ar-123")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestLicensePlateInteractor_GetByPlate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expected := &domain.LicensePlate{
			ID:           "ar-123",
			LicensePlate: "HK-5432",
			ModelName:    "Boeing 737",
			AirlineName:  "Avianca",
		}
		svc := &fakeLicensePlateService{
			getByPlateFn: func(ctx context.Context, plate string) (*domain.LicensePlate, error) {
				return expected, nil
			},
		}
		interactor := NewLicensePlateInteractor(svc, noopLogger{})

		result, err := interactor.GetLicensePlateByPlate(context.Background(), "HK-5432")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.LicensePlate != expected.LicensePlate {
			t.Errorf("expected LicensePlate %q, got %q", expected.LicensePlate, result.LicensePlate)
		}
	})

	t.Run("service error", func(t *testing.T) {
		svc := &fakeLicensePlateService{
			getByPlateFn: func(ctx context.Context, plate string) (*domain.LicensePlate, error) {
				return nil, errors.New("database error")
			},
		}
		interactor := NewLicensePlateInteractor(svc, noopLogger{})

		_, err := interactor.GetLicensePlateByPlate(context.Background(), "HK-5432")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestLicensePlateInteractor_List(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expected := []domain.LicensePlate{
			{ID: "ar-1", LicensePlate: "HK-5432"},
			{ID: "ar-2", LicensePlate: "CC-BFA"},
		}
		svc := &fakeLicensePlateService{
			listFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.LicensePlate, error) {
				return expected, nil
			},
		}
		interactor := NewLicensePlateInteractor(svc, noopLogger{})

		result, err := interactor.ListLicensePlates(context.Background(), nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 2 {
			t.Errorf("expected 2 registrations, got %d", len(result))
		}
	})

	t.Run("service error", func(t *testing.T) {
		svc := &fakeLicensePlateService{
			listFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.LicensePlate, error) {
				return nil, errors.New("database error")
			},
		}
		interactor := NewLicensePlateInteractor(svc, noopLogger{})

		_, err := interactor.ListLicensePlates(context.Background(), nil)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestLicensePlateInteractor_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &fakeLicensePlateService{
			createFn: func(ctx context.Context, registration domain.LicensePlate) error {
				if registration.LicensePlate != "HK-5432" {
					t.Errorf("expected LicensePlate HK-5432, got %s", registration.LicensePlate)
				}
				return nil
			},
		}
		interactor := NewLicensePlateInteractor(svc, noopLogger{})

		registration := domain.LicensePlate{
			ID:              "new-uuid",
			LicensePlate:    "HK-5432",
			AircraftModelID: "model-uuid",
			AirlineID:       "airline-uuid",
		}
		err := interactor.CreateLicensePlate(context.Background(), registration)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("service error", func(t *testing.T) {
		svc := &fakeLicensePlateService{
			createFn: func(ctx context.Context, registration domain.LicensePlate) error {
				return errors.New("create failed")
			},
		}
		interactor := NewLicensePlateInteractor(svc, noopLogger{})

		err := interactor.CreateLicensePlate(context.Background(), domain.LicensePlate{})
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestLicensePlateInteractor_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &fakeLicensePlateService{
			updateFn: func(ctx context.Context, registration domain.LicensePlate) error {
				return nil
			},
		}
		interactor := NewLicensePlateInteractor(svc, noopLogger{})

		err := interactor.UpdateLicensePlate(context.Background(), domain.LicensePlate{
			ID:              "ar-123",
			LicensePlate:    "HK-9999",
			AircraftModelID: "model-uuid",
			AirlineID:       "airline-uuid",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		svc := &fakeLicensePlateService{
			updateFn: func(ctx context.Context, registration domain.LicensePlate) error {
				return domain.ErrLicensePlateNotFound
			},
		}
		interactor := NewLicensePlateInteractor(svc, noopLogger{})

		err := interactor.UpdateLicensePlate(context.Background(), domain.LicensePlate{ID: "non-existent"})
		if err != domain.ErrLicensePlateNotFound {
			t.Errorf("expected ErrLicensePlateNotFound, got %v", err)
		}
	})

	t.Run("update service error", func(t *testing.T) {
		svc := &fakeLicensePlateService{
			updateFn: func(ctx context.Context, registration domain.LicensePlate) error {
				return errors.New("update failed")
			},
		}
		interactor := NewLicensePlateInteractor(svc, noopLogger{})

		err := interactor.UpdateLicensePlate(context.Background(), domain.LicensePlate{ID: "ar-123"})
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
