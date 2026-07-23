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
	getByIDFn              func(ctx context.Context, id string) (*domain.AirlineRoute, error)
	listFn                 func(ctx context.Context, filters map[string]interface{}) ([]domain.AirlineRoute, error)
	listByAirlineIDFn      func(ctx context.Context, airlineID string) ([]domain.AirlineRoute, error)
	activateFn             func(ctx context.Context, id string) error
	deactivateFn           func(ctx context.Context, id string) error
	activateTxFn           func(ctx context.Context, tx output.Tx, id string) error
	deactivateTxFn         func(ctx context.Context, tx output.Tx, id string) error
	getByRouteAndAirlineFn func(ctx context.Context, routeID, airlineID string) (*domain.AirlineRoute, error)
	saveTxFn               func(ctx context.Context, tx output.Tx, airlineRoute domain.AirlineRoute) error
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

func (f *fakeAirlineRouteService) ActivateAirlineRouteTx(ctx context.Context, tx output.Tx, id string) error {
	if f.activateTxFn != nil {
		return f.activateTxFn(ctx, tx, id)
	}
	return nil
}

func (f *fakeAirlineRouteService) DeactivateAirlineRouteTx(ctx context.Context, tx output.Tx, id string) error {
	if f.deactivateTxFn != nil {
		return f.deactivateTxFn(ctx, tx, id)
	}
	return nil
}

func (f *fakeAirlineRouteService) GetAirlineRouteByRouteAndAirline(ctx context.Context, routeID, airlineID string) (*domain.AirlineRoute, error) {
	if f.getByRouteAndAirlineFn != nil {
		return f.getByRouteAndAirlineFn(ctx, routeID, airlineID)
	}
	return nil, errors.New("not implemented")
}

func (f *fakeAirlineRouteService) SaveAirlineRouteTx(ctx context.Context, tx output.Tx, airlineRoute domain.AirlineRoute) error {
	if f.saveTxFn != nil {
		return f.saveTxFn(ctx, tx, airlineRoute)
	}
	return nil
}

func TestNewAirlineRouteInteractor(t *testing.T) {
	svc := &fakeAirlineRouteService{}
	interactor := NewAirlineRouteInteractor(svc, &fakeRouteService{})
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
		interactor := NewAirlineRouteInteractor(svc, &fakeRouteService{})

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
		interactor := NewAirlineRouteInteractor(svc, &fakeRouteService{})

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
		interactor := NewAirlineRouteInteractor(svc, &fakeRouteService{})

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
		interactor := NewAirlineRouteInteractor(svc, &fakeRouteService{})

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
		interactor := NewAirlineRouteInteractor(svc, &fakeRouteService{})

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
		interactor := NewAirlineRouteInteractor(svc, &fakeRouteService{})

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
		interactor := NewAirlineRouteInteractor(svc, &fakeRouteService{})

		_, err := interactor.ListMyAirlineRoutes(context.Background(), "trace-123", "airline-123")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestAirlineRouteInteractor_ActivateAirlineRoute(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &fakeAirlineRouteService{
			activateTxFn: func(ctx context.Context, tx output.Tx, id string) error {
				return nil
			},
		}
		interactor := NewAirlineRouteInteractor(svc, &fakeRouteService{})

		err := interactor.ActivateAirlineRoute(context.Background(), "trace-123", "route-123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		svc := &fakeAirlineRouteService{
			activateTxFn: func(ctx context.Context, tx output.Tx, id string) error {
				return domain.ErrAirlineRouteNotFound
			},
		}
		interactor := NewAirlineRouteInteractor(svc, &fakeRouteService{})

		err := interactor.ActivateAirlineRoute(context.Background(), "trace-123", "nonexistent")
		if err != domain.ErrAirlineRouteNotFound {
			t.Errorf("expected ErrAirlineRouteNotFound, got %v", err)
		}
	})

	t.Run("already active - idempotent", func(t *testing.T) {
		svc := &fakeAirlineRouteService{
			activateTxFn: func(ctx context.Context, tx output.Tx, id string) error {
				return domain.ErrAirlineRouteAlreadyActive
			},
		}
		interactor := NewAirlineRouteInteractor(svc, &fakeRouteService{})

		err := interactor.ActivateAirlineRoute(context.Background(), "trace-123", "route-123")
		if err != domain.ErrAirlineRouteAlreadyActive {
			t.Errorf("expected ErrAirlineRouteAlreadyActive, got %v", err)
		}
	})

	t.Run("service error", func(t *testing.T) {
		svc := &fakeAirlineRouteService{
			activateTxFn: func(ctx context.Context, tx output.Tx, id string) error {
				return errors.New("activation failed")
			},
		}
		interactor := NewAirlineRouteInteractor(svc, &fakeRouteService{})

		err := interactor.ActivateAirlineRoute(context.Background(), "trace-123", "route-123")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestAirlineRouteInteractor_DeactivateAirlineRoute(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &fakeAirlineRouteService{
			deactivateTxFn: func(ctx context.Context, tx output.Tx, id string) error {
				return nil
			},
		}
		interactor := NewAirlineRouteInteractor(svc, &fakeRouteService{})

		err := interactor.DeactivateAirlineRoute(context.Background(), "trace-123", "route-123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		svc := &fakeAirlineRouteService{
			deactivateTxFn: func(ctx context.Context, tx output.Tx, id string) error {
				return domain.ErrAirlineRouteNotFound
			},
		}
		interactor := NewAirlineRouteInteractor(svc, &fakeRouteService{})

		err := interactor.DeactivateAirlineRoute(context.Background(), "trace-123", "nonexistent")
		if err != domain.ErrAirlineRouteNotFound {
			t.Errorf("expected ErrAirlineRouteNotFound, got %v", err)
		}
	})

	t.Run("already inactive - idempotent", func(t *testing.T) {
		svc := &fakeAirlineRouteService{
			deactivateTxFn: func(ctx context.Context, tx output.Tx, id string) error {
				return domain.ErrAirlineRouteAlreadyInactive
			},
		}
		interactor := NewAirlineRouteInteractor(svc, &fakeRouteService{})

		err := interactor.DeactivateAirlineRoute(context.Background(), "trace-123", "route-123")
		if err != domain.ErrAirlineRouteAlreadyInactive {
			t.Errorf("expected ErrAirlineRouteAlreadyInactive, got %v", err)
		}
	})

	t.Run("service error", func(t *testing.T) {
		svc := &fakeAirlineRouteService{
			deactivateTxFn: func(ctx context.Context, tx output.Tx, id string) error {
				return errors.New("deactivation failed")
			},
		}
		interactor := NewAirlineRouteInteractor(svc, &fakeRouteService{})

		err := interactor.DeactivateAirlineRoute(context.Background(), "trace-123", "route-123")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestAirlineRouteInteractor_ResolveOrCreatePendingAirlineRoute(t *testing.T) {
	t.Run("route not configured", func(t *testing.T) {
		routeSvc := &fakeRouteService{
			getByAirportsFn: func(ctx context.Context, originAirportID, destinationAirportID string) (*domain.Route, error) {
				return nil, domain.ErrRouteNotFound
			},
		}
		interactor := NewAirlineRouteInteractor(&fakeAirlineRouteService{}, routeSvc)

		_, created, err := interactor.ResolveOrCreatePendingAirlineRoute(context.Background(), "trace-123", "airline-1", "ap-1", "ap-2")
		if err != domain.ErrRouteNotFound {
			t.Errorf("expected ErrRouteNotFound, got %v", err)
		}
		if created {
			t.Error("expected created=false")
		}
	})

	t.Run("link already exists — returns it without creating anything", func(t *testing.T) {
		existing := &domain.AirlineRoute{ID: "ar-1", RouteID: "route-1", AirlineID: "airline-1", Status: domain.AirlineRouteStatusActive}
		routeSvc := &fakeRouteService{
			getByAirportsFn: func(ctx context.Context, originAirportID, destinationAirportID string) (*domain.Route, error) {
				return &domain.Route{ID: "route-1"}, nil
			},
		}
		saveCalled := false
		svc := &fakeAirlineRouteService{
			getByRouteAndAirlineFn: func(ctx context.Context, routeID, airlineID string) (*domain.AirlineRoute, error) {
				return existing, nil
			},
			saveTxFn: func(ctx context.Context, tx output.Tx, airlineRoute domain.AirlineRoute) error {
				saveCalled = true
				return nil
			},
		}
		interactor := NewAirlineRouteInteractor(svc, routeSvc)

		result, created, err := interactor.ResolveOrCreatePendingAirlineRoute(context.Background(), "trace-123", "airline-1", "ap-1", "ap-2")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if created {
			t.Error("expected created=false")
		}
		if result.ID != "ar-1" {
			t.Errorf("expected existing link 'ar-1', got %q", result.ID)
		}
		if saveCalled {
			t.Error("expected SaveAirlineRouteTx not to be called when a link already exists")
		}
	})

	t.Run("no link yet — creates it as pending", func(t *testing.T) {
		routeSvc := &fakeRouteService{
			getByAirportsFn: func(ctx context.Context, originAirportID, destinationAirportID string) (*domain.Route, error) {
				return &domain.Route{ID: "route-1"}, nil
			},
		}
		lookupCalls := 0
		var savedStatus string
		svc := &fakeAirlineRouteService{
			getByRouteAndAirlineFn: func(ctx context.Context, routeID, airlineID string) (*domain.AirlineRoute, error) {
				lookupCalls++
				if lookupCalls == 1 {
					return nil, domain.ErrAirlineRouteNotFound
				}
				return &domain.AirlineRoute{ID: "ar-new", RouteID: routeID, AirlineID: airlineID, Status: savedStatus}, nil
			},
			saveTxFn: func(ctx context.Context, tx output.Tx, airlineRoute domain.AirlineRoute) error {
				savedStatus = airlineRoute.Status
				return nil
			},
		}
		interactor := NewAirlineRouteInteractor(svc, routeSvc)

		result, created, err := interactor.ResolveOrCreatePendingAirlineRoute(context.Background(), "trace-123", "airline-1", "ap-1", "ap-2")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !created {
			t.Error("expected created=true")
		}
		if result.Status != domain.AirlineRouteStatusPending {
			t.Errorf("expected status=pending, got %q", result.Status)
		}
		if savedStatus != domain.AirlineRouteStatusPending {
			t.Errorf("expected the saved row to use status=pending, got %q", savedStatus)
		}
	})

	t.Run("concurrent creation — re-fetches instead of failing", func(t *testing.T) {
		routeSvc := &fakeRouteService{
			getByAirportsFn: func(ctx context.Context, originAirportID, destinationAirportID string) (*domain.Route, error) {
				return &domain.Route{ID: "route-1"}, nil
			},
		}
		lookupCalls := 0
		svc := &fakeAirlineRouteService{
			getByRouteAndAirlineFn: func(ctx context.Context, routeID, airlineID string) (*domain.AirlineRoute, error) {
				lookupCalls++
				if lookupCalls == 1 {
					return nil, domain.ErrAirlineRouteNotFound
				}
				return &domain.AirlineRoute{ID: "ar-race", RouteID: routeID, AirlineID: airlineID, Status: domain.AirlineRouteStatusPending}, nil
			},
			saveTxFn: func(ctx context.Context, tx output.Tx, airlineRoute domain.AirlineRoute) error {
				return domain.ErrAirlineRouteAlreadyExists
			},
		}
		interactor := NewAirlineRouteInteractor(svc, routeSvc)

		result, _, err := interactor.ResolveOrCreatePendingAirlineRoute(context.Background(), "trace-123", "airline-1", "ap-1", "ap-2")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.ID != "ar-race" {
			t.Errorf("expected the concurrently-created link 'ar-race', got %q", result.ID)
		}
	})

	t.Run("save error propagates", func(t *testing.T) {
		routeSvc := &fakeRouteService{
			getByAirportsFn: func(ctx context.Context, originAirportID, destinationAirportID string) (*domain.Route, error) {
				return &domain.Route{ID: "route-1"}, nil
			},
		}
		svc := &fakeAirlineRouteService{
			getByRouteAndAirlineFn: func(ctx context.Context, routeID, airlineID string) (*domain.AirlineRoute, error) {
				return nil, domain.ErrAirlineRouteNotFound
			},
			saveTxFn: func(ctx context.Context, tx output.Tx, airlineRoute domain.AirlineRoute) error {
				return errors.New("db error")
			},
		}
		interactor := NewAirlineRouteInteractor(svc, routeSvc)

		_, _, err := interactor.ResolveOrCreatePendingAirlineRoute(context.Background(), "trace-123", "airline-1", "ap-1", "ap-2")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
