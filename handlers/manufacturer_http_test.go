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

type fakeManufacturerService struct {
	getByIDFn func(ctx context.Context, id string) (*domain.Manufacturer, error)
	listFn    func(ctx context.Context) ([]domain.Manufacturer, error)
}

var _ input.ManufacturerService = (*fakeManufacturerService)(nil)

func (f *fakeManufacturerService) GetManufacturerByID(ctx context.Context, id string) (*domain.Manufacturer, error) {
	if f.getByIDFn != nil {
		return f.getByIDFn(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (f *fakeManufacturerService) ListManufacturers(ctx context.Context) ([]domain.Manufacturer, error) {
	if f.listFn != nil {
		return f.listFn(ctx)
	}
	return nil, errors.New("not implemented")
}

func newTestManufacturerMessageCache(t *testing.T) *messaging.MessageCache {
	t.Helper()

	repo := fakeMessageCacheRepo{messages: []cachetypes.CachedMessage{
		{Code: domain.MsgManufacturerGetOK, Type: cachetypes.TypeSuccess, Content: "manufacturer retrieved successfully"},
		{Code: domain.MsgManufacturerNotFound, Type: cachetypes.TypeError, Content: "manufacturer not found"},
		{Code: domain.MsgServerError, Type: cachetypes.TypeError, Content: "internal server error"},
		{Code: domain.MsgValIDInvalid, Type: cachetypes.TypeError, Content: "invalid id"},
	}}
	cache := messaging.NewMessageCache(repo, 0)
	if err := cache.LoadMessages(context.Background()); err != nil {
		t.Fatalf("failed to load message cache: %v", err)
	}
	return cache
}

func TestHTTP_GetManufacturerByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := newTestManufacturerMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)

	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	newRouter := func(svc input.ManufacturerService) *gin.Engine {
		manufacturerInteractor := interactor.NewManufacturerInteractor(svc)
		h := New(nil, nil, enc, resp, nil, nil, nil, nil, nil, nil, manufacturerInteractor, nil, nil, nil)

		r := gin.New()
		r.Use(middleware.RequestID())
		r.Use(errHandler.Handle())
		r.GET("/manufacturers/:id", h.GetManufacturerByID())
		return r
	}

	t.Run("success with obfuscated ID", func(t *testing.T) {
		manufacturerUUID := "550e8400-e29b-41d4-a716-446655440000"
		encodedID, err := enc.Encode(manufacturerUUID)
		if err != nil {
			t.Fatalf("failed to encode ID: %v", err)
		}

		expectedManufacturer := &domain.Manufacturer{
			ID:   manufacturerUUID,
			Name: "Boeing",
		}

		svc := &fakeManufacturerService{
			getByIDFn: func(ctx context.Context, id string) (*domain.Manufacturer, error) {
				if id != manufacturerUUID {
					t.Errorf("expected id %s, got %s", manufacturerUUID, id)
				}
				return expectedManufacturer, nil
			},
		}

		router := newRouter(svc)

		req := httptest.NewRequest(http.MethodGet, "/manufacturers/"+encodedID, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d. body=%s", w.Code, w.Body.String())
		}

		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if response["success"] != true {
			t.Errorf("expected success=true, got %v", response["success"])
		}
	})

	t.Run("manufacturer not found", func(t *testing.T) {
		manufacturerUUID := "550e8400-e29b-41d4-a716-446655440000"
		encodedID, _ := enc.Encode(manufacturerUUID)

		svc := &fakeManufacturerService{
			getByIDFn: func(ctx context.Context, id string) (*domain.Manufacturer, error) {
				return nil, domain.ErrManufacturerNotFound
			},
		}

		router := newRouter(svc)

		req := httptest.NewRequest(http.MethodGet, "/manufacturers/"+encodedID, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if response["success"] != false {
			t.Errorf("expected success=false, got %v", response["success"])
		}
	})

	t.Run("invalid ID format => 400", func(t *testing.T) {
		svc := &fakeManufacturerService{}
		router := newRouter(svc)

		req := httptest.NewRequest(http.MethodGet, "/manufacturers/invalid-id!!!", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d. body=%s", w.Code, w.Body.String())
		}
	})

	t.Run("generic service error => 500", func(t *testing.T) {
		manufacturerUUID := "550e8400-e29b-41d4-a716-446655440000"
		encodedID, _ := enc.Encode(manufacturerUUID)

		svc := &fakeManufacturerService{
			getByIDFn: func(ctx context.Context, id string) (*domain.Manufacturer, error) {
				return nil, errors.New("database connection failed")
			},
		}

		router := newRouter(svc)
		req := httptest.NewRequest(http.MethodGet, "/manufacturers/"+encodedID, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d. body=%s", w.Code, w.Body.String())
		}
	})

	t.Run("empty ID returns 404", func(t *testing.T) {
		svc := &fakeManufacturerService{}
		router := newRouter(svc)

		req := httptest.NewRequest(http.MethodGet, "/manufacturers/", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Should return 404 because the route doesn't match
		if w.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", w.Code)
		}
	})
}

func TestHTTP_ListManufacturers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := newTestManufacturerMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)

	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	newRouter := func(svc input.ManufacturerService) *gin.Engine {
		manufacturerInteractor := interactor.NewManufacturerInteractor(svc)
		h := New(nil, nil, enc, resp, nil, nil, nil, nil, nil, nil, manufacturerInteractor, nil, nil, nil)

		r := gin.New()
		r.Use(middleware.RequestID())
		r.Use(errHandler.Handle())
		r.GET("/manufacturers", h.ListManufacturers())
		return r
	}

	t.Run("success - returns list of manufacturers", func(t *testing.T) {
		expectedManufacturers := []domain.Manufacturer{
			{ID: "uuid-1", Name: "Boeing"},
			{ID: "uuid-2", Name: "Airbus"},
			{ID: "uuid-3", Name: "Embraer"},
		}

		svc := &fakeManufacturerService{
			listFn: func(ctx context.Context) ([]domain.Manufacturer, error) {
				return expectedManufacturers, nil
			},
		}

		router := newRouter(svc)
		req := httptest.NewRequest(http.MethodGet, "/manufacturers", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d. body=%s", w.Code, w.Body.String())
		}
	})

	t.Run("success - empty list", func(t *testing.T) {
		svc := &fakeManufacturerService{
			listFn: func(ctx context.Context) ([]domain.Manufacturer, error) {
				return []domain.Manufacturer{}, nil
			},
		}

		router := newRouter(svc)
		req := httptest.NewRequest(http.MethodGet, "/manufacturers", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("service error => 500", func(t *testing.T) {
		svc := &fakeManufacturerService{
			listFn: func(ctx context.Context) ([]domain.Manufacturer, error) {
				return nil, errors.New("database error")
			},
		}

		router := newRouter(svc)
		req := httptest.NewRequest(http.MethodGet, "/manufacturers", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", w.Code)
		}
	})
}
