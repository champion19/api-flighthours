package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/champion19/api-flighthours/core/interactor"
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/input"
	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/middleware"
	"github.com/champion19/api-flighthours/platform/cache/messaging"
	cachetypes "github.com/champion19/api-flighthours/platform/cache/types"
	"github.com/champion19/api-flighthours/tools/idencoder"
	"github.com/gin-gonic/gin"
)

type fakeAirlineRouteService struct {
	getByIDFn                  func(ctx context.Context, id string) (*domain.AirlineRoute, error)
	listFn                     func(ctx context.Context, filters map[string]interface{}) ([]domain.AirlineRoute, error)
	listByAirlineIDFn          func(ctx context.Context, airlineID string) ([]domain.AirlineRoute, error)
	activateFn                 func(ctx context.Context, id string) error
	deactivateFn               func(ctx context.Context, id string) error
	beginTxFn                  func(ctx context.Context) (output.Tx, error)
	updateAirlineRouteStatusFn func(ctx context.Context, tx output.Tx, id string, status bool) error
}

var _ input.AirlineRouteService = (*fakeAirlineRouteService)(nil)

func (f *fakeAirlineRouteService) BeginTx(ctx context.Context) (output.Tx, error) {
	if f.beginTxFn != nil {
		return f.beginTxFn(ctx)
	}
	return fakeTx{}, nil
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
	return errors.New("not implemented")
}

func (f *fakeAirlineRouteService) DeactivateAirlineRoute(ctx context.Context, id string) error {
	if f.deactivateFn != nil {
		return f.deactivateFn(ctx, id)
	}
	return errors.New("not implemented")
}
func (f *fakeAirlineRouteService) ActivateAirlineRouteTx(ctx context.Context, tx output.Tx, id string) error {
	if f.activateFn != nil {
		return f.activateFn(ctx, id)
	}
	return nil
}
func (f *fakeAirlineRouteService) DeactivateAirlineRouteTx(ctx context.Context, tx output.Tx, id string) error {
	if f.deactivateFn != nil {
		return f.deactivateFn(ctx, id)
	}
	return nil
}

func newTestAirlineRouteMessageCache(t *testing.T) *messaging.MessageCache {
	t.Helper()

	repo := fakeMessageCacheRepo{messages: []cachetypes.CachedMessage{
		{Code: domain.MsgAirlineRouteListOK, Type: cachetypes.TypeSuccess, Content: "airline routes listed successfully"},
		{Code: domain.MsgAirlineRouteActivateOK, Type: cachetypes.TypeSuccess, Content: "airline route activated"},
		{Code: domain.MsgAirlineRouteDeactivateOK, Type: cachetypes.TypeSuccess, Content: "airline route deactivated"},
		{Code: domain.MsgAirlineRouteNotFound, Type: cachetypes.TypeError, Content: "airline route not found"},
		{Code: domain.MsgServerError, Type: cachetypes.TypeError, Content: "internal server error"},
		{Code: domain.MsgValIDInvalid, Type: cachetypes.TypeError, Content: "invalid id"},
	}}
	cache := messaging.NewMessageCache(repo, 0)
	if err := cache.LoadMessages(context.Background()); err != nil {
		t.Fatalf("failed to load message cache: %v", err)
	}
	return cache
}

func TestHTTP_ListAirlineRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := newTestAirlineRouteMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)

	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	newRouter := func(svc input.AirlineRouteService) *gin.Engine {
		airlineRouteInteractor := interactor.NewAirlineRouteInteractor(svc)
		h := New(nil, nil, enc, resp, nil, nil, nil, nil, nil, nil, nil, nil, airlineRouteInteractor, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middleware.RequestID())
		r.Use(errHandler.Handle())
		r.GET("/airline-routes", h.ListAirlineRoutes())
		return r
	}

	t.Run("success - list all", func(t *testing.T) {
		expectedRoutes := []domain.AirlineRoute{
			{
				ID:                  "route-uuid-1",
				RouteID:             "base-route-1",
				AirlineID:           "airline-uuid-1",
				Status:              true,
				AirlineCode:         "TST",
				AirlineName:         "Test Airlines",
				OriginIataCode:      "BOG",
				DestinationIataCode: "MDE",
				RouteCode:           "BOG-MDE",
				EstimatedFlightTime: "60",
			},
		}

		svc := &fakeAirlineRouteService{
			listFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.AirlineRoute, error) {
				return expectedRoutes, nil
			},
		}

		r := newRouter(svc)
		req := httptest.NewRequest(http.MethodGet, "/airline-routes", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusOK, w.Code, w.Body.String())
		}

		var out middleware.APIResponse
		if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
			t.Fatalf("invalid json response: %v; body=%s", err, w.Body.String())
		}
		if !out.Success {
			t.Fatalf("expected success=true, got false")
		}
		if out.Code != domain.MsgAirlineRouteListOK {
			t.Fatalf("expected code %q, got %q", domain.MsgAirlineRouteListOK, out.Code)
		}
	})

	t.Run("success - list with airline_code filter", func(t *testing.T) {
		filterCalled := false
		svc := &fakeAirlineRouteService{
			listFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.AirlineRoute, error) {
				filterCalled = true
				if filters["airline_code"] != "TST" {
					t.Errorf("expected airline_code filter 'TST', got %v", filters["airline_code"])
				}
				return []domain.AirlineRoute{}, nil
			},
		}

		r := newRouter(svc)
		req := httptest.NewRequest(http.MethodGet, "/airline-routes?airline_code=TST", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusOK, w.Code, w.Body.String())
		}
		if !filterCalled {
			t.Fatal("expected filter function to be called")
		}
	})

	t.Run("success - list with status filter", func(t *testing.T) {
		svc := &fakeAirlineRouteService{
			listFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.AirlineRoute, error) {
				if filters["status"] != true {
					t.Errorf("expected status filter true, got %v", filters["status"])
				}
				return []domain.AirlineRoute{}, nil
			},
		}

		r := newRouter(svc)
		req := httptest.NewRequest(http.MethodGet, "/airline-routes?status=true", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusOK, w.Code, w.Body.String())
		}
	})

	t.Run("service error returns 500", func(t *testing.T) {
		svc := &fakeAirlineRouteService{
			listFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.AirlineRoute, error) {
				return nil, errors.New("database error")
			},
		}

		r := newRouter(svc)
		req := httptest.NewRequest(http.MethodGet, "/airline-routes", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusInternalServerError, w.Code, w.Body.String())
		}
	})

	t.Run("success - list with status=false filter", func(t *testing.T) {
		svc := &fakeAirlineRouteService{
			listFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.AirlineRoute, error) {
				if filters["status"] != false {
					t.Errorf("expected status filter false, got %v", filters["status"])
				}
				return []domain.AirlineRoute{}, nil
			},
		}

		r := newRouter(svc)
		req := httptest.NewRequest(http.MethodGet, "/airline-routes?status=false", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusOK, w.Code, w.Body.String())
		}
	})

	t.Run("success - list with status=0 filter", func(t *testing.T) {
		svc := &fakeAirlineRouteService{
			listFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.AirlineRoute, error) {
				if filters["status"] != false {
					t.Errorf("expected status filter false (from 0), got %v", filters["status"])
				}
				return []domain.AirlineRoute{}, nil
			},
		}

		r := newRouter(svc)
		req := httptest.NewRequest(http.MethodGet, "/airline-routes?status=0", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusOK, w.Code, w.Body.String())
		}
	})

	t.Run("success - list with status=1 filter", func(t *testing.T) {
		svc := &fakeAirlineRouteService{
			listFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.AirlineRoute, error) {
				if filters["status"] != true {
					t.Errorf("expected status filter true (from 1), got %v", filters["status"])
				}
				return []domain.AirlineRoute{}, nil
			},
		}

		r := newRouter(svc)
		req := httptest.NewRequest(http.MethodGet, "/airline-routes?status=1", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusOK, w.Code, w.Body.String())
		}
	})
}

func TestHTTP_ActivateAirlineRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := newTestAirlineRouteMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)

	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	newRouter := func(svc input.AirlineRouteService) *gin.Engine {
		airlineRouteInteractor := interactor.NewAirlineRouteInteractor(svc)
		h := New(nil, nil, enc, resp, nil, nil, nil, nil, nil, nil, nil, nil, airlineRouteInteractor, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middleware.RequestID())
		r.Use(errHandler.Handle())
		r.PATCH("/airline-routes/:id/activate", h.ActivateAirlineRoute())
		return r
	}

	t.Run("success", func(t *testing.T) {
		routeUUID := "550e8400-e29b-41d4-a716-446655440000"
		encodedID, _ := enc.Encode(routeUUID)
		activateCalled := false

		svc := &fakeAirlineRouteService{
			getByIDFn: func(ctx context.Context, id string) (*domain.AirlineRoute, error) {
				return &domain.AirlineRoute{
					ID:        routeUUID,
					AirlineID: "airline-uuid",
					Status:    false,
				}, nil
			},
			activateFn: func(ctx context.Context, id string) error {
				activateCalled = true
				return nil
			},
		}

		r := newRouter(svc)
		req := httptest.NewRequest(http.MethodPatch, "/airline-routes/"+encodedID+"/activate", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusOK, w.Code, w.Body.String())
		}
		if !activateCalled {
			t.Fatal("expected activate to be called")
		}

		var out middleware.APIResponse
		if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
			t.Fatalf("invalid json response: %v; body=%s", err, w.Body.String())
		}
		if !out.Success {
			t.Fatalf("expected success=true, got false")
		}
	})

	t.Run("already active - idempotent success", func(t *testing.T) {
		routeUUID := "550e8400-e29b-41d4-a716-446655440001"
		encodedID, _ := enc.Encode(routeUUID)

		svc := &fakeAirlineRouteService{
			getByIDFn: func(ctx context.Context, id string) (*domain.AirlineRoute, error) {
				return &domain.AirlineRoute{
					ID:        routeUUID,
					AirlineID: "airline-uuid",
					Status:    true,
				}, nil
			},
			activateFn: func(ctx context.Context, id string) error {
				return domain.ErrAirlineRouteAlreadyActive
			},
		}

		r := newRouter(svc)
		req := httptest.NewRequest(http.MethodPatch, "/airline-routes/"+encodedID+"/activate", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected status %d (idempotent), got %d. body=%s", http.StatusOK, w.Code, w.Body.String())
		}
	})

	t.Run("airline route not found => 404", func(t *testing.T) {
		routeUUID := "550e8400-e29b-41d4-a716-446655440002"
		encodedID, _ := enc.Encode(routeUUID)

		svc := &fakeAirlineRouteService{
			activateFn: func(ctx context.Context, id string) error {
				return domain.ErrAirlineRouteNotFound
			},
		}

		r := newRouter(svc)
		req := httptest.NewRequest(http.MethodPatch, "/airline-routes/"+encodedID+"/activate", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusNotFound {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusNotFound, w.Code, w.Body.String())
		}
	})

	t.Run("invalid ID => 400", func(t *testing.T) {
		svc := &fakeAirlineRouteService{}

		r := newRouter(svc)
		req := httptest.NewRequest(http.MethodPatch, "/airline-routes/invalid-id!!!/activate", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusBadRequest, w.Code, w.Body.String())
		}
	})

	t.Run("generic error => 500", func(t *testing.T) {
		routeUUID := "550e8400-e29b-41d4-a716-446655440003"
		encodedID, _ := enc.Encode(routeUUID)

		svc := &fakeAirlineRouteService{
			activateFn: func(ctx context.Context, id string) error {
				return errors.New("database error")
			},
		}

		r := newRouter(svc)
		req := httptest.NewRequest(http.MethodPatch, "/airline-routes/"+encodedID+"/activate", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusInternalServerError, w.Code, w.Body.String())
		}
	})
}

func TestHTTP_DeactivateAirlineRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := newTestAirlineRouteMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)

	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	newRouter := func(svc input.AirlineRouteService) *gin.Engine {
		airlineRouteInteractor := interactor.NewAirlineRouteInteractor(svc)
		h := New(nil, nil, enc, resp, nil, nil, nil, nil, nil, nil, nil, nil, airlineRouteInteractor, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middleware.RequestID())
		r.Use(errHandler.Handle())
		r.PATCH("/airline-routes/:id/deactivate", h.DeactivateAirlineRoute())
		return r
	}

	t.Run("success", func(t *testing.T) {
		routeUUID := "550e8400-e29b-41d4-a716-446655440000"
		encodedID, _ := enc.Encode(routeUUID)
		deactivateCalled := false

		svc := &fakeAirlineRouteService{
			getByIDFn: func(ctx context.Context, id string) (*domain.AirlineRoute, error) {
				return &domain.AirlineRoute{
					ID:        routeUUID,
					AirlineID: "airline-uuid",
					Status:    true,
				}, nil
			},
			deactivateFn: func(ctx context.Context, id string) error {
				deactivateCalled = true
				return nil
			},
		}

		r := newRouter(svc)
		req := httptest.NewRequest(http.MethodPatch, "/airline-routes/"+encodedID+"/deactivate", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusOK, w.Code, w.Body.String())
		}
		if !deactivateCalled {
			t.Fatal("expected deactivate to be called")
		}

		var out middleware.APIResponse
		if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
			t.Fatalf("invalid json response: %v; body=%s", err, w.Body.String())
		}
		if !out.Success {
			t.Fatalf("expected success=true, got false")
		}
	})

	t.Run("already inactive - idempotent success", func(t *testing.T) {
		routeUUID := "550e8400-e29b-41d4-a716-446655440001"
		encodedID, _ := enc.Encode(routeUUID)

		svc := &fakeAirlineRouteService{
			getByIDFn: func(ctx context.Context, id string) (*domain.AirlineRoute, error) {
				return &domain.AirlineRoute{
					ID:        routeUUID,
					AirlineID: "airline-uuid",
					Status:    false,
				}, nil
			},
			deactivateFn: func(ctx context.Context, id string) error {
				return domain.ErrAirlineRouteAlreadyInactive
			},
		}

		r := newRouter(svc)
		req := httptest.NewRequest(http.MethodPatch, "/airline-routes/"+encodedID+"/deactivate", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected status %d (idempotent), got %d. body=%s", http.StatusOK, w.Code, w.Body.String())
		}
	})

	t.Run("airline route not found => 404", func(t *testing.T) {
		routeUUID := "550e8400-e29b-41d4-a716-446655440002"
		encodedID, _ := enc.Encode(routeUUID)

		svc := &fakeAirlineRouteService{
			deactivateFn: func(ctx context.Context, id string) error {
				return domain.ErrAirlineRouteNotFound
			},
		}

		r := newRouter(svc)
		req := httptest.NewRequest(http.MethodPatch, "/airline-routes/"+encodedID+"/deactivate", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusNotFound {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusNotFound, w.Code, w.Body.String())
		}
	})

	t.Run("invalid ID => 400", func(t *testing.T) {
		svc := &fakeAirlineRouteService{}

		r := newRouter(svc)
		req := httptest.NewRequest(http.MethodPatch, "/airline-routes/invalid-id!!!/deactivate", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusBadRequest, w.Code, w.Body.String())
		}
	})

	t.Run("generic error => 500", func(t *testing.T) {
		routeUUID := "550e8400-e29b-41d4-a716-446655440003"
		encodedID, _ := enc.Encode(routeUUID)

		svc := &fakeAirlineRouteService{
			deactivateFn: func(ctx context.Context, id string) error {
				return errors.New("database error")
			},
		}

		r := newRouter(svc)
		req := httptest.NewRequest(http.MethodPatch, "/airline-routes/"+encodedID+"/deactivate", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusInternalServerError, w.Code, w.Body.String())
		}
	})
}

// ── TestHTTP_ListMyAirlineRoutes ───────────────────────────────────────

func TestHTTP_ListMyAirlineRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := newTestAirlineRouteMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)

	enc, _ := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})

	empID := "550e8400-e29b-41d4-a716-446655440001"
	airlineID := "550e8400-e29b-41d4-a716-446655440099"

	newMyRoutesRouter := func(aeSvc *fakeAirlineEmployeeService, arSvc *fakeAirlineRouteService, authUser *domain.Employee) *gin.Engine {
		aeInteractor := interactor.NewAirlineEmployeeInteractor(aeSvc)
		airlineRouteInteractor := interactor.NewAirlineRouteInteractor(arSvc)
		h := New(nil, &fakeEmployeeInteractor{}, enc, resp, nil, nil, nil, aeInteractor, nil, nil, nil, nil, airlineRouteInteractor, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middleware.RequestID())
		r.Use(errHandler.Handle())
		r.Use(func(c *gin.Context) {
			if authUser != nil {
				c.Set("authenticated_user", authUser)
			}
			c.Next()
		})
		r.GET("/employees/airline-routes", h.ListMyAirlineRoutes())
		return r
	}

	t.Run("success", func(t *testing.T) {
		aeSvc := &fakeAirlineEmployeeService{
			getByIDFn: func(_ context.Context, id string) (*domain.AirlineEmployee, error) {
				return &domain.AirlineEmployee{ID: empID, AirlineID: airlineID, Active: true}, nil
			},
		}
		arSvc := &fakeAirlineRouteService{
			listByAirlineIDFn: func(_ context.Context, airID string) ([]domain.AirlineRoute, error) {
				return []domain.AirlineRoute{
					{ID: "ar-1", RouteID: "r-1", AirlineID: airlineID, Status: true, RouteCode: "BOG-MDE"},
				}, nil
			},
		}
		router := newMyRoutesRouter(aeSvc, arSvc, &domain.Employee{ID: empID})
		req := httptest.NewRequest(http.MethodGet, "/employees/airline-routes", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		t.Logf("status: %d", w.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		aeSvc := &fakeAirlineEmployeeService{}
		arSvc := &fakeAirlineRouteService{}
		router := newMyRoutesRouter(aeSvc, arSvc, nil)
		req := httptest.NewRequest(http.MethodGet, "/employees/airline-routes", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		t.Logf("status: %d", w.Code)
	})

	t.Run("no airline info", func(t *testing.T) {
		aeSvc := &fakeAirlineEmployeeService{
			getByIDFn: func(_ context.Context, id string) (*domain.AirlineEmployee, error) {
				return nil, nil
			},
		}
		arSvc := &fakeAirlineRouteService{}
		router := newMyRoutesRouter(aeSvc, arSvc, &domain.Employee{ID: empID})
		req := httptest.NewRequest(http.MethodGet, "/employees/airline-routes", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		t.Logf("status: %d", w.Code)
	})

	t.Run("route list error", func(t *testing.T) {
		aeSvc := &fakeAirlineEmployeeService{
			getByIDFn: func(_ context.Context, id string) (*domain.AirlineEmployee, error) {
				return &domain.AirlineEmployee{ID: empID, AirlineID: airlineID, Active: true}, nil
			},
		}
		arSvc := &fakeAirlineRouteService{
			listByAirlineIDFn: func(_ context.Context, airID string) ([]domain.AirlineRoute, error) {
				return nil, errors.New("db error")
			},
		}
		router := newMyRoutesRouter(aeSvc, arSvc, &domain.Employee{ID: empID})
		req := httptest.NewRequest(http.MethodGet, "/employees/airline-routes", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		t.Logf("status: %d", w.Code)
	})
}
