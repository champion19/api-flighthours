package interactor

import (
	"context"
	"errors"
	"testing"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/input"
)

// fakeRouteService implements input.RouteService for testing
type fakeRouteService struct {
	getByIDFn       func(ctx context.Context, id string) (*domain.Route, error)
	listFn          func(ctx context.Context, filters map[string]interface{}) ([]domain.Route, error)
	getByAirportsFn func(ctx context.Context, originAirportID, destinationAirportID string) (*domain.Route, error)
}

var _ input.RouteService = (*fakeRouteService)(nil)

func (f *fakeRouteService) GetRouteByID(ctx context.Context, id string) (*domain.Route, error) {
	if f.getByIDFn != nil {
		return f.getByIDFn(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (f *fakeRouteService) ListRoutes(ctx context.Context, filters map[string]interface{}) ([]domain.Route, error) {
	if f.listFn != nil {
		return f.listFn(ctx, filters)
	}
	return nil, errors.New("not implemented")
}

func (f *fakeRouteService) GetRouteByAirports(ctx context.Context, originAirportID, destinationAirportID string) (*domain.Route, error) {
	if f.getByAirportsFn != nil {
		return f.getByAirportsFn(ctx, originAirportID, destinationAirportID)
	}
	return nil, errors.New("not implemented")
}

func TestNewRouteInteractor(t *testing.T) {
	svc := &fakeRouteService{}
	interactor := NewRouteInteractor(svc)
	if interactor == nil {
		t.Error("expected non-nil RouteInteractor")
	}
}

func TestRouteInteractor_GetRouteByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expectedRoute := &domain.Route{
			ID:                  "route-123",
			OriginIataCode:      "BOG",
			DestinationIataCode: "CLO",
			RouteCode:           "BOG-CLO",
		}
		svc := &fakeRouteService{
			getByIDFn: func(ctx context.Context, id string) (*domain.Route, error) {
				return expectedRoute, nil
			},
		}
		interactor := NewRouteInteractor(svc)

		result, err := interactor.GetRouteByID(context.Background(), "route-123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.ID != expectedRoute.ID {
			t.Errorf("expected ID %q, got %q", expectedRoute.ID, result.ID)
		}
		if result.RouteCode != expectedRoute.RouteCode {
			t.Errorf("expected RouteCode %q, got %q", expectedRoute.RouteCode, result.RouteCode)
		}
	})

	t.Run("not found", func(t *testing.T) {
		svc := &fakeRouteService{
			getByIDFn: func(ctx context.Context, id string) (*domain.Route, error) {
				return nil, domain.ErrRouteNotFound
			},
		}
		interactor := NewRouteInteractor(svc)

		_, err := interactor.GetRouteByID(context.Background(), "nonexistent")
		if err != domain.ErrRouteNotFound {
			t.Errorf("expected ErrRouteNotFound, got %v", err)
		}
	})

	t.Run("service error", func(t *testing.T) {
		svc := &fakeRouteService{
			getByIDFn: func(ctx context.Context, id string) (*domain.Route, error) {
				return nil, errors.New("database error")
			},
		}
		interactor := NewRouteInteractor(svc)

		_, err := interactor.GetRouteByID(context.Background(), "route-123")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestRouteInteractor_ListRoutes(t *testing.T) {
	t.Run("success without filters", func(t *testing.T) {
		expectedRoutes := []domain.Route{
			{ID: "route-1", RouteCode: "BOG-CLO"},
			{ID: "route-2", RouteCode: "MDE-BOG"},
		}
		svc := &fakeRouteService{
			listFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.Route, error) {
				return expectedRoutes, nil
			},
		}
		interactor := NewRouteInteractor(svc)

		result, err := interactor.ListRoutes(context.Background(), nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 2 {
			t.Errorf("expected 2 routes, got %d", len(result))
		}
	})

	t.Run("success with filters", func(t *testing.T) {
		svc := &fakeRouteService{
			listFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.Route, error) {
				if _, ok := filters["origin"]; !ok {
					t.Error("expected 'origin' filter")
				}
				return []domain.Route{{ID: "route-1"}}, nil
			},
		}
		interactor := NewRouteInteractor(svc)

		filters := map[string]interface{}{"origin": "BOG"}
		result, err := interactor.ListRoutes(context.Background(), filters)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 1 {
			t.Errorf("expected 1 route, got %d", len(result))
		}
	})

	t.Run("empty list", func(t *testing.T) {
		svc := &fakeRouteService{
			listFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.Route, error) {
				return []domain.Route{}, nil
			},
		}
		interactor := NewRouteInteractor(svc)

		result, err := interactor.ListRoutes(context.Background(), nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 0 {
			t.Errorf("expected 0 routes, got %d", len(result))
		}
	})

	t.Run("service error", func(t *testing.T) {
		svc := &fakeRouteService{
			listFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.Route, error) {
				return nil, errors.New("database error")
			},
		}
		interactor := NewRouteInteractor(svc)

		_, err := interactor.ListRoutes(context.Background(), nil)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
