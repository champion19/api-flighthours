package services

import (
	"context"
	"errors"
	"testing"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

type fakeManufacturerRepo struct {
	getByIDFn func(ctx context.Context, id string) (*domain.Manufacturer, error)
	listFn    func(ctx context.Context) ([]domain.Manufacturer, error)
}

func (f fakeManufacturerRepo) GetManufacturerByID(ctx context.Context, id string) (*domain.Manufacturer, error) {
	if f.getByIDFn != nil {
		return f.getByIDFn(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (f fakeManufacturerRepo) ListManufacturers(ctx context.Context) ([]domain.Manufacturer, error) {
	if f.listFn != nil {
		return f.listFn(ctx)
	}
	return nil, errors.New("not implemented")
}

func TestManufacturerService_NewManufacturerService(t *testing.T) {
	repo := fakeManufacturerRepo{}
	svc := NewManufacturerService(repo)
	if svc == nil {
		t.Fatal("expected service, got nil")
	}
}

func TestManufacturerService_GetManufacturerByID(t *testing.T) {
	ctx := context.Background()

	t.Run("success => returns manufacturer", func(t *testing.T) {
		expectedManufacturer := &domain.Manufacturer{
			ID:   "manufacturer-123",
			Name: "Boeing",
		}
		svc := NewManufacturerService(fakeManufacturerRepo{
			getByIDFn: func(context.Context, string) (*domain.Manufacturer, error) {
				return expectedManufacturer, nil
			},
		})

		result, err := svc.GetManufacturerByID(ctx, "manufacturer-123")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result.ID != expectedManufacturer.ID {
			t.Fatalf("expected ID %s, got %s", expectedManufacturer.ID, result.ID)
		}
		if result.Name != expectedManufacturer.Name {
			t.Fatalf("expected name %s, got %s", expectedManufacturer.Name, result.Name)
		}
	})

	t.Run("not found => returns error", func(t *testing.T) {
		svc := NewManufacturerService(fakeManufacturerRepo{
			getByIDFn: func(context.Context, string) (*domain.Manufacturer, error) {
				return nil, domain.ErrManufacturerNotFound
			},
		})

		_, err := svc.GetManufacturerByID(ctx, "non-existent")
		if !errors.Is(err, domain.ErrManufacturerNotFound) {
			t.Fatalf("expected %v, got %v", domain.ErrManufacturerNotFound, err)
		}
	})

	t.Run("propagates repository error", func(t *testing.T) {
		repoErr := errors.New("database connection error")
		svc := NewManufacturerService(fakeManufacturerRepo{
			getByIDFn: func(context.Context, string) (*domain.Manufacturer, error) {
				return nil, repoErr
			},
		})

		_, err := svc.GetManufacturerByID(ctx, "manufacturer-123")
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected %v, got %v", repoErr, err)
		}
	})
}

func TestManufacturerService_ListManufacturers(t *testing.T) {
	ctx := context.Background()

	t.Run("success => returns list of manufacturers", func(t *testing.T) {
		expectedManufacturers := []domain.Manufacturer{
			{ID: "m-1", Name: "Boeing"},
			{ID: "m-2", Name: "Airbus"},
			{ID: "m-3", Name: "Embraer"},
		}
		svc := NewManufacturerService(fakeManufacturerRepo{
			listFn: func(context.Context) ([]domain.Manufacturer, error) {
				return expectedManufacturers, nil
			},
		})

		result, err := svc.ListManufacturers(ctx)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(result) != 3 {
			t.Errorf("expected 3 manufacturers, got %d", len(result))
		}
	})

	t.Run("empty list => returns empty slice", func(t *testing.T) {
		svc := NewManufacturerService(fakeManufacturerRepo{
			listFn: func(context.Context) ([]domain.Manufacturer, error) {
				return []domain.Manufacturer{}, nil
			},
		})

		result, err := svc.ListManufacturers(ctx)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(result) != 0 {
			t.Errorf("expected 0 manufacturers, got %d", len(result))
		}
	})

	t.Run("propagates repository error", func(t *testing.T) {
		repoErr := errors.New("database error")
		svc := NewManufacturerService(fakeManufacturerRepo{
			listFn: func(context.Context) ([]domain.Manufacturer, error) {
				return nil, repoErr
			},
		})

		_, err := svc.ListManufacturers(ctx)
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected %v, got %v", repoErr, err)
		}
	})
}
