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

func newTestEngineMessageCache(t *testing.T) *messaging.MessageCache {
	t.Helper()

	repo := fakeMessageCacheRepo{messages: []cachetypes.CachedMessage{
		{Code: domain.MsgEngineGetOK, Type: cachetypes.TypeSuccess, Content: "engine retrieved successfully"},
		{Code: domain.MsgEngineNotFound, Type: cachetypes.TypeError, Content: "engine not found"},
		{Code: domain.MsgEngineListOK, Type: cachetypes.TypeSuccess, Content: "engines listed successfully"},
		{Code: domain.MsgServerError, Type: cachetypes.TypeError, Content: "internal server error"},
		{Code: domain.MsgValIDInvalid, Type: cachetypes.TypeError, Content: "invalid id"},
	}}
	cache := messaging.NewMessageCache(repo, 0)
	if err := cache.LoadMessages(context.Background()); err != nil {
		t.Fatalf("failed to load message cache: %v", err)
	}
	return cache
}

func TestHTTP_GetEngineByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := newTestEngineMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)

	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	newRouter := func(svc input.EngineService) *gin.Engine {
		engineInteractor := interactor.NewEngineInteractor(svc)
		h := New(HandlerDeps{
		IDEncoder: enc,
		Response: resp,
		EngineInteractor: engineInteractor,
		})

		r := gin.New()
		r.Use(middleware.RequestID())
		r.Use(errHandler.Handle())
		r.GET("/engines/:id", h.GetEngineByID())
		return r
	}

	t.Run("success with obfuscated ID", func(t *testing.T) {
		engineUUID := "550e8400-e29b-41d4-a716-446655440000"
		encodedID, _ := enc.Encode(engineUUID)
		expectedEngine := &domain.Engine{
			ID:   engineUUID,
			Name: "Turbofan",
		}

		svc := &fakeEngineService{
			getByIDFn: func(ctx context.Context, id string) (*domain.Engine, error) {
				if id != engineUUID {
					t.Errorf("expected id %s, got %s", engineUUID, id)
				}
				return expectedEngine, nil
			},
		}

		r := newRouter(svc)
		req := httptest.NewRequest(http.MethodGet, "/engines/"+encodedID, nil)
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
		if out.Code != domain.MsgEngineGetOK {
			t.Fatalf("expected code %q, got %q", domain.MsgEngineGetOK, out.Code)
		}
	})

	t.Run("engine not found => 404", func(t *testing.T) {
		engineUUID := "550e8400-e29b-41d4-a716-446655440001"
		encodedID, _ := enc.Encode(engineUUID)

		svc := &fakeEngineService{
			getByIDFn: func(context.Context, string) (*domain.Engine, error) {
				return nil, domain.ErrEngineNotFound
			},
		}

		r := newRouter(svc)
		req := httptest.NewRequest(http.MethodGet, "/engines/"+encodedID, nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusNotFound {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusNotFound, w.Code, w.Body.String())
		}
	})

	t.Run("invalid ID => 400", func(t *testing.T) {
		svc := &fakeEngineService{}

		r := newRouter(svc)
		req := httptest.NewRequest(http.MethodGet, "/engines/invalid-id-format", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusBadRequest, w.Code, w.Body.String())
		}
	})

	t.Run("interactor returns generic error => 500", func(t *testing.T) {
		engineUUID := "550e8400-e29b-41d4-a716-446655440002"
		encodedID, _ := enc.Encode(engineUUID)

		svc := &fakeEngineService{
			getByIDFn: func(context.Context, string) (*domain.Engine, error) {
				return nil, errors.New("database connection error")
			},
		}

		r := newRouter(svc)
		req := httptest.NewRequest(http.MethodGet, "/engines/"+encodedID, nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusInternalServerError, w.Code, w.Body.String())
		}
	})
}

func TestHTTP_ListEngines(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := newTestEngineMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)

	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	newRouter := func(svc input.EngineService) *gin.Engine {
		engineInteractor := interactor.NewEngineInteractor(svc)
		h := New(HandlerDeps{
		IDEncoder: enc,
		Response: resp,
		EngineInteractor: engineInteractor,
		})

		r := gin.New()
		r.Use(middleware.RequestID())
		r.Use(errHandler.Handle())
		r.GET("/engines", h.ListEngines())
		return r
	}

	t.Run("success - list all engines", func(t *testing.T) {
		expectedEngines := []domain.Engine{
			{ID: "engine-1", Name: "Turbofan"},
			{ID: "engine-2", Name: "Turboprop"},
		}

		svc := &fakeEngineService{
			listFn: func(ctx context.Context) ([]domain.Engine, error) {
				return expectedEngines, nil
			},
		}

		r := newRouter(svc)
		req := httptest.NewRequest(http.MethodGet, "/engines", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusOK, w.Code, w.Body.String())
		}
	})

	t.Run("success - empty list", func(t *testing.T) {
		svc := &fakeEngineService{
			listFn: func(ctx context.Context) ([]domain.Engine, error) {
				return []domain.Engine{}, nil
			},
		}

		r := newRouter(svc)
		req := httptest.NewRequest(http.MethodGet, "/engines", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusOK, w.Code, w.Body.String())
		}
	})

	t.Run("service error => 500", func(t *testing.T) {
		svc := &fakeEngineService{
			listFn: func(ctx context.Context) ([]domain.Engine, error) {
				return nil, errors.New("database error")
			},
		}

		r := newRouter(svc)
		req := httptest.NewRequest(http.MethodGet, "/engines", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusInternalServerError, w.Code, w.Body.String())
		}
	})
}
