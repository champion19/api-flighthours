package interactor

import (
	"context"
	"errors"
	"testing"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/input"
)

// fakeEngineService implements input.EngineService for testing
type fakeEngineService struct {
	getByIDFn func(ctx context.Context, id string) (*domain.Engine, error)
	listFn    func(ctx context.Context) ([]domain.Engine, error)
}

var _ input.EngineService = (*fakeEngineService)(nil)

func (f *fakeEngineService) GetEngineByID(ctx context.Context, id string) (*domain.Engine, error) {
	if f.getByIDFn != nil {
		return f.getByIDFn(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (f *fakeEngineService) ListEngines(ctx context.Context) ([]domain.Engine, error) {
	if f.listFn != nil {
		return f.listFn(ctx)
	}
	return nil, errors.New("not implemented")
}

func TestNewEngineInteractor(t *testing.T) {
	svc := &fakeEngineService{}
	interactor := NewEngineInteractor(svc)
	if interactor == nil {
		t.Error("expected non-nil EngineInteractor")
	}
}

func TestEngineInteractor_GetEngineByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expectedEngine := &domain.Engine{
			ID:   "engine-123",
			Name: "Turbofan XL",
		}
		svc := &fakeEngineService{
			getByIDFn: func(ctx context.Context, id string) (*domain.Engine, error) {
				return expectedEngine, nil
			},
		}
		interactor := NewEngineInteractor(svc)

		result, err := interactor.GetEngineByID(context.Background(), "engine-123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.ID != expectedEngine.ID {
			t.Errorf("expected ID %q, got %q", expectedEngine.ID, result.ID)
		}
		if result.Name != expectedEngine.Name {
			t.Errorf("expected Name %q, got %q", expectedEngine.Name, result.Name)
		}
	})

	t.Run("not found", func(t *testing.T) {
		svc := &fakeEngineService{
			getByIDFn: func(ctx context.Context, id string) (*domain.Engine, error) {
				return nil, domain.ErrEngineNotFound
			},
		}
		interactor := NewEngineInteractor(svc)

		_, err := interactor.GetEngineByID(context.Background(), "nonexistent")
		if err != domain.ErrEngineNotFound {
			t.Errorf("expected ErrEngineNotFound, got %v", err)
		}
	})

	t.Run("service error", func(t *testing.T) {
		svc := &fakeEngineService{
			getByIDFn: func(ctx context.Context, id string) (*domain.Engine, error) {
				return nil, errors.New("database error")
			},
		}
		interactor := NewEngineInteractor(svc)

		_, err := interactor.GetEngineByID(context.Background(), "engine-123")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestEngineInteractor_ListEngines(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expectedEngines := []domain.Engine{
			{ID: "engine-1", Name: "Turbofan"},
			{ID: "engine-2", Name: "Turboprop"},
		}
		svc := &fakeEngineService{
			listFn: func(ctx context.Context) ([]domain.Engine, error) {
				return expectedEngines, nil
			},
		}
		interactor := NewEngineInteractor(svc)

		result, err := interactor.ListEngines(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 2 {
			t.Errorf("expected 2 engines, got %d", len(result))
		}
	})

	t.Run("empty list", func(t *testing.T) {
		svc := &fakeEngineService{
			listFn: func(ctx context.Context) ([]domain.Engine, error) {
				return []domain.Engine{}, nil
			},
		}
		interactor := NewEngineInteractor(svc)

		result, err := interactor.ListEngines(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 0 {
			t.Errorf("expected 0 engines, got %d", len(result))
		}
	})

	t.Run("service error", func(t *testing.T) {
		svc := &fakeEngineService{
			listFn: func(ctx context.Context) ([]domain.Engine, error) {
				return nil, errors.New("database error")
			},
		}
		interactor := NewEngineInteractor(svc)

		_, err := interactor.ListEngines(context.Background())
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
