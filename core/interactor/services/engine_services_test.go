package services

import (
	"context"
	"errors"
	"testing"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

// Mock repository for testing
type mockEngineRepo struct {
	getByIDFn func(ctx context.Context, id string) (*domain.Engine, error)
	listFn    func(ctx context.Context) ([]domain.Engine, error)
}

func (m *mockEngineRepo) GetEngineByID(ctx context.Context, id string) (*domain.Engine, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *mockEngineRepo) ListEngines(ctx context.Context) ([]domain.Engine, error) {
	if m.listFn != nil {
		return m.listFn(ctx)
	}
	return nil, nil
}

func (m *mockEngineRepo) BeginTx(ctx context.Context) (interface{}, error) {
	return nil, nil
}

func TestNewEngineService(t *testing.T) {
	t.Run("creates engine service", func(t *testing.T) {
		repo := &mockEngineRepo{}
		service := NewEngineService(repo)

		if service == nil {
			t.Error("expected non-nil service")
		}
	})
}

func TestEngineService_GetEngineByID(t *testing.T) {
	t.Run("returns engine when found", func(t *testing.T) {
		expected := &domain.Engine{
			ID:   "engine-123",
			Name: "Test Engine",
		}

		repo := &mockEngineRepo{
			getByIDFn: func(ctx context.Context, id string) (*domain.Engine, error) {
				return expected, nil
			},
		}

		service := NewEngineService(repo)
		result, err := service.GetEngineByID(context.Background(), "engine-123")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.ID != "engine-123" {
			t.Errorf("expected ID 'engine-123', got %q", result.ID)
		}
	})

	t.Run("returns error when not found", func(t *testing.T) {
		repo := &mockEngineRepo{
			getByIDFn: func(ctx context.Context, id string) (*domain.Engine, error) {
				return nil, errors.New("not found")
			},
		}

		service := NewEngineService(repo)
		_, err := service.GetEngineByID(context.Background(), "non-existent")

		if err == nil {
			t.Error("expected error for non-existent engine")
		}
	})
}

func TestEngineService_ListEngines(t *testing.T) {
	t.Run("returns list of engines", func(t *testing.T) {
		expected := []domain.Engine{
			{ID: "engine-1", Name: "Engine 1"},
			{ID: "engine-2", Name: "Engine 2"},
		}

		repo := &mockEngineRepo{
			listFn: func(ctx context.Context) ([]domain.Engine, error) {
				return expected, nil
			},
		}

		service := NewEngineService(repo)
		result, err := service.ListEngines(context.Background())

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 2 {
			t.Errorf("expected 2 engines, got %d", len(result))
		}
	})

	t.Run("returns error on failure", func(t *testing.T) {
		repo := &mockEngineRepo{
			listFn: func(ctx context.Context) ([]domain.Engine, error) {
				return nil, errors.New("database error")
			},
		}

		service := NewEngineService(repo)
		_, err := service.ListEngines(context.Background())

		if err == nil {
			t.Error("expected error on database failure")
		}
	})
}
