package services

import (
	"context"
	"errors"
	"testing"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/platform/logger"
)

// Mock route repository
type mockRouteRepo struct {
	getByIDFn func(ctx context.Context, id string) (*domain.Route, error)
	listFn    func(ctx context.Context, filters map[string]interface{}) ([]domain.Route, error)
}

func (m *mockRouteRepo) GetRouteByID(ctx context.Context, id string) (*domain.Route, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *mockRouteRepo) ListRoutes(ctx context.Context, filters map[string]interface{}) ([]domain.Route, error) {
	if m.listFn != nil {
		return m.listFn(ctx, filters)
	}
	return nil, nil
}

func (m *mockRouteRepo) BeginTx(ctx context.Context) (interface{}, error) {
	return nil, nil
}

func TestNewRouteService(t *testing.T) {
	t.Run("creates route service", func(t *testing.T) {
		repo := &mockRouteRepo{}
		log := logger.NewSlogLogger()
		service := NewRouteService(repo, log)

		if service == nil {
			t.Error("expected non-nil service")
		}
	})
}

func TestRouteService_GetRouteByID(t *testing.T) {
	t.Run("returns route when found", func(t *testing.T) {
		expected := &domain.Route{
			ID: "route-123",
		}

		repo := &mockRouteRepo{
			getByIDFn: func(ctx context.Context, id string) (*domain.Route, error) {
				return expected, nil
			},
		}

		log := logger.NewSlogLogger()
		service := NewRouteService(repo, log)
		result, err := service.GetRouteByID(context.Background(), "route-123")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.ID != "route-123" {
			t.Errorf("expected ID 'route-123', got %q", result.ID)
		}
	})

	t.Run("returns error when not found", func(t *testing.T) {
		repo := &mockRouteRepo{
			getByIDFn: func(ctx context.Context, id string) (*domain.Route, error) {
				return nil, errors.New("not found")
			},
		}

		log := logger.NewSlogLogger()
		service := NewRouteService(repo, log)
		_, err := service.GetRouteByID(context.Background(), "non-existent")

		if err == nil {
			t.Error("expected error for non-existent route")
		}
	})
}

func TestRouteService_ListRoutes(t *testing.T) {
	t.Run("returns list of routes", func(t *testing.T) {
		expected := []domain.Route{
			{ID: "route-1"},
			{ID: "route-2"},
		}

		repo := &mockRouteRepo{
			listFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.Route, error) {
				return expected, nil
			},
		}

		log := logger.NewSlogLogger()
		service := NewRouteService(repo, log)
		result, err := service.ListRoutes(context.Background(), nil)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 2 {
			t.Errorf("expected 2 routes, got %d", len(result))
		}
	})

	t.Run("passes filters to repository", func(t *testing.T) {
		var passedFilters map[string]interface{}
		repo := &mockRouteRepo{
			listFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.Route, error) {
				passedFilters = filters
				return nil, nil
			},
		}

		log := logger.NewSlogLogger()
		service := NewRouteService(repo, log)
		filters := map[string]interface{}{"active": true}
		_, _ = service.ListRoutes(context.Background(), filters)

		if passedFilters["active"] != true {
			t.Error("expected filters to be passed to repository")
		}
	})

	t.Run("returns error on failure", func(t *testing.T) {
		repo := &mockRouteRepo{
			listFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.Route, error) {
				return nil, errors.New("database error")
			},
		}

		log := logger.NewSlogLogger()
		service := NewRouteService(repo, log)
		_, err := service.ListRoutes(context.Background(), nil)

		if err == nil {
			t.Error("expected error on database failure")
		}
	})
}
