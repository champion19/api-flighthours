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
	"github.com/champion19/api-flighthours/middleware"
	"github.com/champion19/api-flighthours/platform/cache/messaging"
	cachetypes "github.com/champion19/api-flighthours/platform/cache/types"
	"github.com/champion19/api-flighthours/tools/idencoder"
	"github.com/gin-gonic/gin"
)

type fakeRouteService struct {
	getByIDFn func(ctx context.Context, id string) (*domain.Route, error)
	listFn    func(ctx context.Context, filters map[string]interface{}) ([]domain.Route, error)
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

func newTestRouteMessageCache(t *testing.T) *messaging.MessageCache {
	t.Helper()

	repo := fakeMessageCacheRepo{messages: []cachetypes.CachedMessage{
		{Code: domain.MsgRouteGetOK, Type: cachetypes.TypeSuccess, Content: "route retrieved successfully"},
		{Code: domain.MsgRouteNotFound, Type: cachetypes.TypeError, Content: "route not found"},
		{Code: domain.MsgRouteListOK, Type: cachetypes.TypeSuccess, Content: "routes listed successfully"},
		{Code: domain.MsgServerError, Type: cachetypes.TypeError, Content: "internal server error"},
		{Code: domain.MsgValIDInvalid, Type: cachetypes.TypeError, Content: "invalid id"},
	}}
	cache := messaging.NewMessageCache(repo, 0)
	if err := cache.LoadMessages(context.Background()); err != nil {
		t.Fatalf("failed to load message cache: %v", err)
	}
	return cache
}

func TestHTTP_GetRouteByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := newTestRouteMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)

	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	newRouter := func(svc input.RouteService) *gin.Engine {
		routeInteractor := interactor.NewRouteInteractor(svc, noopLogger{})
		h := New(nil, nil, enc, resp, nil, nil, nil, nil, nil, routeInteractor, nil)

		r := gin.New()
		r.Use(middleware.RequestID())
		r.Use(errHandler.Handle())
		r.GET("/routes/:id", h.GetRouteByID())
		return r
	}

	t.Run("success with obfuscated ID", func(t *testing.T) {
		routeUUID := "550e8400-e29b-41d4-a716-446655440000"
		encodedID, _ := enc.Encode(routeUUID)
		expectedRoute := &domain.Route{
			ID:                   routeUUID,
			OriginIataCode:       "BOG",
			DestinationIataCode:  "MDE",
			RouteCode:            "BOG-MDE",
			OriginAirportID:      "origin-uuid",
			DestinationAirportID: "dest-uuid",
		}

		svc := &fakeRouteService{
			getByIDFn: func(ctx context.Context, id string) (*domain.Route, error) {
				if id != routeUUID {
					t.Errorf("expected id %s, got %s", routeUUID, id)
				}
				return expectedRoute, nil
			},
		}

		r := newRouter(svc)
		req := httptest.NewRequest(http.MethodGet, "/routes/"+encodedID, nil)
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
		if out.Code != domain.MsgRouteGetOK {
			t.Fatalf("expected code %q, got %q", domain.MsgRouteGetOK, out.Code)
		}
	})

	t.Run("route not found => 404", func(t *testing.T) {
		routeUUID := "550e8400-e29b-41d4-a716-446655440001"
		encodedID, _ := enc.Encode(routeUUID)

		svc := &fakeRouteService{
			getByIDFn: func(context.Context, string) (*domain.Route, error) {
				return nil, domain.ErrRouteNotFound
			},
		}

		r := newRouter(svc)
		req := httptest.NewRequest(http.MethodGet, "/routes/"+encodedID, nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusNotFound {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusNotFound, w.Code, w.Body.String())
		}
	})

	t.Run("invalid ID => 400", func(t *testing.T) {
		svc := &fakeRouteService{}

		r := newRouter(svc)
		req := httptest.NewRequest(http.MethodGet, "/routes/invalid-id-format", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusBadRequest, w.Code, w.Body.String())
		}
	})

	t.Run("interactor returns generic error => 500", func(t *testing.T) {
		routeUUID := "550e8400-e29b-41d4-a716-446655440002"
		encodedID, _ := enc.Encode(routeUUID)

		svc := &fakeRouteService{
			getByIDFn: func(context.Context, string) (*domain.Route, error) {
				return nil, errors.New("database connection error")
			},
		}

		r := newRouter(svc)
		req := httptest.NewRequest(http.MethodGet, "/routes/"+encodedID, nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusInternalServerError, w.Code, w.Body.String())
		}
	})
}

func TestHTTP_ListRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := newTestRouteMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)

	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	newRouter := func(svc input.RouteService) *gin.Engine {
		routeInteractor := interactor.NewRouteInteractor(svc, noopLogger{})
		h := New(nil, nil, enc, resp, nil, nil, nil, nil, nil, routeInteractor, nil)

		r := gin.New()
		r.Use(middleware.RequestID())
		r.Use(errHandler.Handle())
		r.GET("/routes", h.ListRoutes())
		return r
	}

	t.Run("success - list all routes", func(t *testing.T) {
		expectedRoutes := []domain.Route{
			{ID: "route-1", RouteCode: "BOG-MDE", OriginIataCode: "BOG", DestinationIataCode: "MDE"},
			{ID: "route-2", RouteCode: "MDE-BOG", OriginIataCode: "MDE", DestinationIataCode: "BOG"},
		}

		svc := &fakeRouteService{
			listFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.Route, error) {
				return expectedRoutes, nil
			},
		}

		r := newRouter(svc)
		req := httptest.NewRequest(http.MethodGet, "/routes", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusOK, w.Code, w.Body.String())
		}
	})

	t.Run("success - list with airport_type filter", func(t *testing.T) {
		filterCalled := false
		svc := &fakeRouteService{
			listFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.Route, error) {
				filterCalled = true
				if filters["airport_type"] != "Nacional" {
					t.Errorf("expected airport_type filter 'Nacional', got %v", filters["airport_type"])
				}
				return []domain.Route{}, nil
			},
		}

		r := newRouter(svc)
		req := httptest.NewRequest(http.MethodGet, "/routes?airport_type=Nacional", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusOK, w.Code, w.Body.String())
		}
		if !filterCalled {
			t.Fatal("expected filter function to be called")
		}
	})

	t.Run("success - list with origin_country filter", func(t *testing.T) {
		svc := &fakeRouteService{
			listFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.Route, error) {
				if filters["origin_country"] != "Colombia" {
					t.Errorf("expected origin_country filter 'Colombia', got %v", filters["origin_country"])
				}
				return []domain.Route{}, nil
			},
		}

		r := newRouter(svc)
		req := httptest.NewRequest(http.MethodGet, "/routes?origin_country=Colombia", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusOK, w.Code, w.Body.String())
		}
	})

	t.Run("service error => 500", func(t *testing.T) {
		svc := &fakeRouteService{
			listFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.Route, error) {
				return nil, errors.New("database error")
			},
		}

		r := newRouter(svc)
		req := httptest.NewRequest(http.MethodGet, "/routes", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusInternalServerError, w.Code, w.Body.String())
		}
	})
}
