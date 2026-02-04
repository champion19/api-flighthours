package interactor

import (
	"context"
	"errors"
	"testing"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/input"
	"github.com/champion19/api-flighthours/core/ports/output"
)

// fakeAirlineRouteService implements input.AirlineRouteService for testing
type fakeAirlineRouteService struct {
	getByIDFn         func(ctx context.Context, id string) (*domain.AirlineRoute, error)
	listFn            func(ctx context.Context, filters map[string]interface{}) ([]domain.AirlineRoute, error)
	listByAirlineIDFn func(ctx context.Context, airlineID string) ([]domain.AirlineRoute, error)
	activateFn        func(ctx context.Context, id string) error
	deactivateFn      func(ctx context.Context, id string) error
}

var _ input.AirlineRouteService = (*fakeAirlineRouteService)(nil)

func (f *fakeAirlineRouteService) BeginTx(ctx context.Context) (output.Tx, error) {
	return &fakeTx{}, nil
}

func (f *fakeAirlineRouteService) GetAirlineRouteByID(ctx context.Context, id string) (*domain.AirlineRoute, error) {
	if f.getByIDFn != nil {
		return f.getByIDFn(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (f *fakeAirlineRouteService) ListAirlineRoutes(ctx context.Context, filters map[string]interface{}) ([]domain.AirlineRoute, error) {
	if f.listFn != nil {
		return f.listFn(ctx, filters)
	}
	return nil, errors.New("not implemented")
}

func (f *fakeAirlineRouteService) ListAirlineRoutesByAirlineID(ctx context.Context, airlineID string) ([]domain.AirlineRoute, error) {
	if f.listByAirlineIDFn != nil {
		return f.listByAirlineIDFn(ctx, airlineID)
	}
	return nil, errors.New("not implemented")
}

func (f *fakeAirlineRouteService) ActivateAirlineRoute(ctx context.Context, id string) error {
	if f.activateFn != nil {
		return f.activateFn(ctx, id)
	}
	return nil
}

func (f *fakeAirlineRouteService) DeactivateAirlineRoute(ctx context.Context, id string) error {
	if f.deactivateFn != nil {
		return f.deactivateFn(ctx, id)
	}
	return nil
}

func TestNewAirlineRouteInteractor(t *testing.T) {
	svc := &fakeAirlineRouteService{}
	interactor := NewAirlineRouteInteractor(svc)
	if interactor == nil {
		t.Error("expected non-nil AirlineRouteInteractor")
	}
}

func TestAirlineRouteInteractor_GetAirlineRouteByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expectedRoute := &domain.AirlineRoute{
			ID:        "route-123",
			AirlineID: "airline-123",
			RouteID:   "route-1",
			RouteCode: "BOG-CLO",
		}
		svc := &fakeAirlineRouteService{
			getByIDFn: func(ctx context.Context, id string) (*domain.AirlineRoute, error) {
				return expectedRoute, nil
			},
		}
		interactor := NewAirlineRouteInteractor(svc)

		result, err := interactor.GetAirlineRouteByID(context.Background(), "trace-123", "route-123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.ID != expectedRoute.ID {
			t.Errorf("expected ID %q, got %q", expectedRoute.ID, result.ID)
		}
	})

	t.Run("not found", func(t *testing.T) {
		svc := &fakeAirlineRouteService{
			getByIDFn: func(ctx context.Context, id string) (*domain.AirlineRoute, error) {
				return nil, domain.ErrAirlineRouteNotFound
			},
		}
		interactor := NewAirlineRouteInteractor(svc)

		_, err := interactor.GetAirlineRouteByID(context.Background(), "trace-123", "nonexistent")
		if err != domain.ErrAirlineRouteNotFound {
			t.Errorf("expected ErrAirlineRouteNotFound, got %v", err)
		}
	})

	t.Run("service error", func(t *testing.T) {
		svc := &fakeAirlineRouteService{
			getByIDFn: func(ctx context.Context, id string) (*domain.AirlineRoute, error) {
				return nil, errors.New("database error")
			},
		}
		interactor := NewAirlineRouteInteractor(svc)

		_, err := interactor.GetAirlineRouteByID(context.Background(), "trace-123", "route-123")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestAirlineRouteInteractor_ListAirlineRoutes(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expectedRoutes := []domain.AirlineRoute{
			{ID: "route-1", AirlineID: "airline-1"},
			{ID: "route-2", AirlineID: "airline-1"},
		}
		svc := &fakeAirlineRouteService{
			listFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.AirlineRoute, error) {
				return expectedRoutes, nil
			},
		}
		interactor := NewAirlineRouteInteractor(svc)

		result, err := interactor.ListAirlineRoutes(context.Background(), "trace-123", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 2 {
			t.Errorf("expected 2 routes, got %d", len(result))
		}
	})

	t.Run("service error", func(t *testing.T) {
		svc := &fakeAirlineRouteService{
			listFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.AirlineRoute, error) {
				return nil, errors.New("database error")
			},
		}
		interactor := NewAirlineRouteInteractor(svc)

		_, err := interactor.ListAirlineRoutes(context.Background(), "trace-123", nil)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestAirlineRouteInteractor_ListMyAirlineRoutes(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expectedRoutes := []domain.AirlineRoute{
			{ID: "route-1", AirlineID: "airline-123"},
		}
		svc := &fakeAirlineRouteService{
			listByAirlineIDFn: func(ctx context.Context, airlineID string) ([]domain.AirlineRoute, error) {
				if airlineID != "airline-123" {
					t.Errorf("expected airlineID 'airline-123', got %q", airlineID)
				}
				return expectedRoutes, nil
			},
		}
		interactor := NewAirlineRouteInteractor(svc)

		result, err := interactor.ListMyAirlineRoutes(context.Background(), "trace-123", "airline-123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 1 {
			t.Errorf("expected 1 route, got %d", len(result))
		}
	})

	t.Run("service error", func(t *testing.T) {
		svc := &fakeAirlineRouteService{
			listByAirlineIDFn: func(ctx context.Context, airlineID string) ([]domain.AirlineRoute, error) {
				return nil, errors.New("database error")
			},
		}
		interactor := NewAirlineRouteInteractor(svc)

		_, err := interactor.ListMyAirlineRoutes(context.Background(), "trace-123", "airline-123")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestAirlineRouteInteractor_ActivateAirlineRoute(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &fakeAirlineRouteService{
			activateFn: func(ctx context.Context, id string) error {
				return nil
			},
		}
		interactor := NewAirlineRouteInteractor(svc)

		err := interactor.ActivateAirlineRoute(context.Background(), "trace-123", "route-123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		svc := &fakeAirlineRouteService{
			activateFn: func(ctx context.Context, id string) error {
				return domain.ErrAirlineRouteNotFound
			},
		}
		interactor := NewAirlineRouteInteractor(svc)

		err := interactor.ActivateAirlineRoute(context.Background(), "trace-123", "nonexistent")
		if err != domain.ErrAirlineRouteNotFound {
			t.Errorf("expected ErrAirlineRouteNotFound, got %v", err)
		}
	})

	t.Run("already active - idempotent", func(t *testing.T) {
		svc := &fakeAirlineRouteService{
			activateFn: func(ctx context.Context, id string) error {
				return domain.ErrAirlineRouteAlreadyActive
			},
		}
		interactor := NewAirlineRouteInteractor(svc)

		err := interactor.ActivateAirlineRoute(context.Background(), "trace-123", "route-123")
		if err != domain.ErrAirlineRouteAlreadyActive {
			t.Errorf("expected ErrAirlineRouteAlreadyActive, got %v", err)
		}
	})

	t.Run("service error", func(t *testing.T) {
		svc := &fakeAirlineRouteService{
			activateFn: func(ctx context.Context, id string) error {
				return errors.New("activation failed")
			},
		}
		interactor := NewAirlineRouteInteractor(svc)

		err := interactor.ActivateAirlineRoute(context.Background(), "trace-123", "route-123")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestAirlineRouteInteractor_DeactivateAirlineRoute(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &fakeAirlineRouteService{
			deactivateFn: func(ctx context.Context, id string) error {
				return nil
			},
		}
		interactor := NewAirlineRouteInteractor(svc)

		err := interactor.DeactivateAirlineRoute(context.Background(), "trace-123", "route-123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		svc := &fakeAirlineRouteService{
			deactivateFn: func(ctx context.Context, id string) error {
				return domain.ErrAirlineRouteNotFound
			},
		}
		interactor := NewAirlineRouteInteractor(svc)

		err := interactor.DeactivateAirlineRoute(context.Background(), "trace-123", "nonexistent")
		if err != domain.ErrAirlineRouteNotFound {
			t.Errorf("expected ErrAirlineRouteNotFound, got %v", err)
		}
	})

	t.Run("already inactive - idempotent", func(t *testing.T) {
		svc := &fakeAirlineRouteService{
			deactivateFn: func(ctx context.Context, id string) error {
				return domain.ErrAirlineRouteAlreadyInactive
			},
		}
		interactor := NewAirlineRouteInteractor(svc)

		err := interactor.DeactivateAirlineRoute(context.Background(), "trace-123", "route-123")
		if err != domain.ErrAirlineRouteAlreadyInactive {
			t.Errorf("expected ErrAirlineRouteAlreadyInactive, got %v", err)
		}
	})

	t.Run("service error", func(t *testing.T) {
		svc := &fakeAirlineRouteService{
			deactivateFn: func(ctx context.Context, id string) error {
				return errors.New("deactivation failed")
			},
		}
		interactor := NewAirlineRouteInteractor(svc)

		err := interactor.DeactivateAirlineRoute(context.Background(), "trace-123", "route-123")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
