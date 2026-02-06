package interactor

import (
	"context"
	"errors"
	"testing"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/input"
)

type fakeManufacturerServiceForInteractor struct {
	getByIDFn func(ctx context.Context, id string) (*domain.Manufacturer, error)
	listFn    func(ctx context.Context) ([]domain.Manufacturer, error)
}

var _ input.ManufacturerService = (*fakeManufacturerServiceForInteractor)(nil)

func (f *fakeManufacturerServiceForInteractor) GetManufacturerByID(ctx context.Context, id string) (*domain.Manufacturer, error) {
	if f.getByIDFn != nil {
		return f.getByIDFn(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (f *fakeManufacturerServiceForInteractor) ListManufacturers(ctx context.Context) ([]domain.Manufacturer, error) {
	if f.listFn != nil {
		return f.listFn(ctx)
	}
	return nil, errors.New("not implemented")
}

func TestManufacturerInteractor_GetManufacturerByID(t *testing.T) {
	ctx := context.Background()

	t.Run("success => returns manufacturer", func(t *testing.T) {
		expectedManufacturer := &domain.Manufacturer{
			ID:   "manufacturer-123",
			Name: "Boeing",
		}
		svc := &fakeManufacturerServiceForInteractor{
			getByIDFn: func(context.Context, string) (*domain.Manufacturer, error) {
				return expectedManufacturer, nil
			},
		}
		interactor := NewManufacturerInteractor(svc)

		result, err := interactor.GetManufacturerByID(ctx, "manufacturer-123")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result == nil {
			t.Fatal("expected manufacturer, got nil")
		}
		if result.ID != expectedManufacturer.ID {
			t.Errorf("expected ID %s, got %s", expectedManufacturer.ID, result.ID)
		}
		if result.Name != expectedManufacturer.Name {
			t.Errorf("expected name %s, got %s", expectedManufacturer.Name, result.Name)
		}
	})

	t.Run("not found => returns ErrManufacturerNotFound", func(t *testing.T) {
		svc := &fakeManufacturerServiceForInteractor{
			getByIDFn: func(context.Context, string) (*domain.Manufacturer, error) {
				return nil, domain.ErrManufacturerNotFound
			},
		}
		interactor := NewManufacturerInteractor(svc)

		_, err := interactor.GetManufacturerByID(ctx, "non-existent")
		if !errors.Is(err, domain.ErrManufacturerNotFound) {
			t.Fatalf("expected %v, got %v", domain.ErrManufacturerNotFound, err)
		}
	})

	t.Run("service error => propagate error", func(t *testing.T) {
		serviceErr := errors.New("service unavailable")
		svc := &fakeManufacturerServiceForInteractor{
			getByIDFn: func(context.Context, string) (*domain.Manufacturer, error) {
				return nil, serviceErr
			},
		}
		interactor := NewManufacturerInteractor(svc)

		_, err := interactor.GetManufacturerByID(ctx, "manufacturer-123")
		if !errors.Is(err, serviceErr) {
			t.Fatalf("expected %v, got %v", serviceErr, err)
		}
	})
}

func TestManufacturerInteractor_ListManufacturers(t *testing.T) {
	ctx := context.Background()

	t.Run("success => returns list of manufacturers", func(t *testing.T) {
		expectedManufacturers := []domain.Manufacturer{
			{ID: "m-1", Name: "Boeing"},
			{ID: "m-2", Name: "Airbus"},
			{ID: "m-3", Name: "Embraer"},
		}
		svc := &fakeManufacturerServiceForInteractor{
			listFn: func(context.Context) ([]domain.Manufacturer, error) {
				return expectedManufacturers, nil
			},
		}
		interactor := NewManufacturerInteractor(svc)

		result, err := interactor.ListManufacturers(ctx)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(result) != 3 {
			t.Errorf("expected 3 manufacturers, got %d", len(result))
		}
	})

	t.Run("empty list => returns empty slice", func(t *testing.T) {
		svc := &fakeManufacturerServiceForInteractor{
			listFn: func(context.Context) ([]domain.Manufacturer, error) {
				return []domain.Manufacturer{}, nil
			},
		}
		interactor := NewManufacturerInteractor(svc)

		result, err := interactor.ListManufacturers(ctx)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(result) != 0 {
			t.Errorf("expected 0 manufacturers, got %d", len(result))
		}
	})

	t.Run("service error => propagate error", func(t *testing.T) {
		serviceErr := errors.New("database error")
		svc := &fakeManufacturerServiceForInteractor{
			listFn: func(context.Context) ([]domain.Manufacturer, error) {
				return nil, serviceErr
			},
		}
		interactor := NewManufacturerInteractor(svc)

		_, err := interactor.ListManufacturers(ctx)
		if !errors.Is(err, serviceErr) {
			t.Fatalf("expected %v, got %v", serviceErr, err)
		}
	})
}
