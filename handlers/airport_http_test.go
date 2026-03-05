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

type fakeAirportService struct {
	getByIDFn      func(ctx context.Context, id string) (*domain.Airport, error)
	updateStatusFn func(ctx context.Context, id string, status bool) error
	activateFn     func(ctx context.Context, id string) error
	deactivateFn   func(ctx context.Context, id string) error
	beginTxFn      func(ctx context.Context) (output.Tx, error)
	listAirportsFn func(ctx context.Context, filters map[string]interface{}) ([]domain.Airport, error)
	getByTypeFn    func(ctx context.Context, airportType string) ([]domain.Airport, error)
}

var _ input.AirportService = (*fakeAirportService)(nil)

func (f *fakeAirportService) BeginTx(ctx context.Context) (output.Tx, error) {
	if f.beginTxFn != nil {
		return f.beginTxFn(ctx)
	}
	return fakeTx{}, nil
}

func (f *fakeAirportService) GetAirportByID(ctx context.Context, id string) (*domain.Airport, error) {
	if f.getByIDFn != nil {
		return f.getByIDFn(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (f *fakeAirportService) UpdateAirportStatus(ctx context.Context, id string, status bool) error {
	if f.updateStatusFn != nil {
		return f.updateStatusFn(ctx, id, status)
	}
	return errors.New("not implemented")
}

func (f *fakeAirportService) ActivateAirport(ctx context.Context, id string) error {
	if f.activateFn != nil {
		return f.activateFn(ctx, id)
	}
	return errors.New("not implemented")
}

func (f *fakeAirportService) DeactivateAirport(ctx context.Context, id string) error {
	if f.deactivateFn != nil {
		return f.deactivateFn(ctx, id)
	}
	return errors.New("not implemented")
}
func (f *fakeAirportService) ActivateAirportTx(ctx context.Context, tx output.Tx, id string) error {
	if f.activateFn != nil {
		return f.activateFn(ctx, id)
	}
	return nil
}
func (f *fakeAirportService) DeactivateAirportTx(ctx context.Context, tx output.Tx, id string) error {
	if f.deactivateFn != nil {
		return f.deactivateFn(ctx, id)
	}
	return nil
}

func (f *fakeAirportService) ListAirports(ctx context.Context, filters map[string]interface{}) ([]domain.Airport, error) {
	if f.listAirportsFn != nil {
		return f.listAirportsFn(ctx, filters)
	}
	return nil, errors.New("not implemented")
}

func (f *fakeAirportService) GetAirportsByType(ctx context.Context, airportType string) ([]domain.Airport, error) {
	if f.getByTypeFn != nil {
		return f.getByTypeFn(ctx, airportType)
	}
	return nil, errors.New("not implemented")
}

func newTestAirportMessageCache(t *testing.T) *messaging.MessageCache {
	t.Helper()

	repo := fakeMessageCacheRepo{messages: []cachetypes.CachedMessage{
		{Code: domain.MsgAirportGetOK, Type: cachetypes.TypeSuccess, Content: "airport retrieved successfully"},
		{Code: domain.MsgAirportNotFound, Type: cachetypes.TypeError, Content: "airport not found"},
		{Code: domain.MsgAirportActivateOK, Type: cachetypes.TypeSuccess, Content: "airport activated"},
		{Code: domain.MsgAirportDeactivateOK, Type: cachetypes.TypeSuccess, Content: "airport deactivated"},
		{Code: domain.MsgServerError, Type: cachetypes.TypeError, Content: "internal server error"},
		{Code: domain.MsgValIDInvalid, Type: cachetypes.TypeError, Content: "invalid id"},
	}}
	cache := messaging.NewMessageCache(repo, 0)
	if err := cache.LoadMessages(context.Background()); err != nil {
		t.Fatalf("failed to load message cache: %v", err)
	}
	return cache
}

func TestHTTP_GetAirportByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := newTestAirportMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)

	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	newRouter := func(svc input.AirportService) *gin.Engine {
		airportInteractor := interactor.NewAirportInteractor(svc)
		h := New(HandlerDeps{
		IDEncoder: enc,
		Response: resp,
		AirportInteractor: airportInteractor,
		})

		r := gin.New()
		r.Use(middleware.RequestID())
		r.Use(errHandler.Handle())
		r.GET("/airports/:id", h.GetAirportByID())
		return r
	}

	t.Run("success with obfuscated ID", func(t *testing.T) {
		airportUUID := "550e8400-e29b-41d4-a716-446655440000"
		encodedID, err := enc.Encode(airportUUID)
		if err != nil {
			t.Fatalf("failed to encode ID: %v", err)
		}

		expectedAirport := &domain.Airport{
			ID:       airportUUID,
			Name:     "El Dorado International",
			IATACode: "BOG",
			Status:   true,
		}

		svc := &fakeAirportService{
			getByIDFn: func(ctx context.Context, id string) (*domain.Airport, error) {
				if id != airportUUID {
					t.Errorf("expected id %s, got %s", airportUUID, id)
				}
				return expectedAirport, nil
			},
		}

		router := newRouter(svc)

		req := httptest.NewRequest(http.MethodGet, "/airports/"+encodedID, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if response["success"] != true {
			t.Errorf("expected success=true, got %v", response["success"])
		}
	})

	t.Run("airport not found", func(t *testing.T) {
		airportUUID := "550e8400-e29b-41d4-a716-446655440000"
		encodedID, _ := enc.Encode(airportUUID)

		svc := &fakeAirportService{
			getByIDFn: func(ctx context.Context, id string) (*domain.Airport, error) {
				return nil, domain.ErrAirportNotFound
			},
		}

		router := newRouter(svc)

		req := httptest.NewRequest(http.MethodGet, "/airports/"+encodedID, nil)
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

	t.Run("empty id returns error", func(t *testing.T) {
		svc := &fakeAirportService{}
		router := newRouter(svc)

		req := httptest.NewRequest(http.MethodGet, "/airports/", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Should return 404 because the route doesn't match
		if w.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", w.Code)
		}
	})
}

func TestHTTP_DeactivateAirport(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := newTestAirportMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)

	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	newRouter := func(svc input.AirportService) *gin.Engine {
		airportInteractor := interactor.NewAirportInteractor(svc)
		h := New(HandlerDeps{
		IDEncoder: enc,
		Response: resp,
		AirportInteractor: airportInteractor,
		})

		r := gin.New()
		r.Use(middleware.RequestID())
		r.Use(errHandler.Handle())
		r.PATCH("/airports/:id/deactivate", h.DeactivateAirport())
		return r
	}

	t.Run("success", func(t *testing.T) {
		airportUUID := "550e8400-e29b-41d4-a716-446655440000"
		encodedID, _ := enc.Encode(airportUUID)
		deactivateCalled := false

		svc := &fakeAirportService{
			getByIDFn: func(ctx context.Context, id string) (*domain.Airport, error) {
				return &domain.Airport{ID: id, Status: true}, nil
			},
			deactivateFn: func(ctx context.Context, id string) error {
				deactivateCalled = true
				return nil
			},
		}

		router := newRouter(svc)

		req := httptest.NewRequest(http.MethodPatch, "/airports/"+encodedID+"/deactivate", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		if !deactivateCalled {
			t.Error("expected deactivateFn to be called")
		}

		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if response["success"] != true {
			t.Errorf("expected success=true, got %v", response["success"])
		}
	})

	t.Run("airport not found", func(t *testing.T) {
		airportUUID := "550e8400-e29b-41d4-a716-446655440000"
		encodedID, _ := enc.Encode(airportUUID)

		svc := &fakeAirportService{
			getByIDFn: func(ctx context.Context, id string) (*domain.Airport, error) {
				return &domain.Airport{ID: id}, nil
			},
			deactivateFn: func(ctx context.Context, id string) error {
				return domain.ErrAirportNotFound
			},
		}

		router := newRouter(svc)

		req := httptest.NewRequest(http.MethodPatch, "/airports/"+encodedID+"/deactivate", nil)
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

	t.Run("service error", func(t *testing.T) {
		airportUUID := "550e8400-e29b-41d4-a716-446655440000"
		encodedID, _ := enc.Encode(airportUUID)

		svc := &fakeAirportService{
			getByIDFn: func(ctx context.Context, id string) (*domain.Airport, error) {
				return &domain.Airport{ID: id}, nil
			},
			deactivateFn: func(ctx context.Context, id string) error {
				return errors.New("service unavailable")
			},
		}

		router := newRouter(svc)

		req := httptest.NewRequest(http.MethodPatch, "/airports/"+encodedID+"/deactivate", nil)
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

	t.Run("invalid ID => 400", func(t *testing.T) {
		svc := &fakeAirportService{}
		router := newRouter(svc)

		req := httptest.NewRequest(http.MethodPatch, "/airports/invalid-id!!!/deactivate", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}

func TestHTTP_ListAirports(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := newTestAirportMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)

	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	newRouter := func(svc input.AirportService) *gin.Engine {
		airportInteractor := interactor.NewAirportInteractor(svc)
		h := New(HandlerDeps{
		IDEncoder: enc,
		Response: resp,
		AirportInteractor: airportInteractor,
		})

		r := gin.New()
		r.Use(middleware.RequestID())
		r.Use(errHandler.Handle())
		r.GET("/airports", h.ListAirports())
		return r
	}

	t.Run("success - returns list of airports", func(t *testing.T) {
		expectedAirports := []domain.Airport{
			{ID: "uuid-1", Name: "El Dorado", IATACode: "BOG", Status: true},
			{ID: "uuid-2", Name: "Jose Maria Cordova", IATACode: "MDE", Status: true},
		}

		svc := &fakeAirportService{
			listAirportsFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.Airport, error) {
				return expectedAirports, nil
			},
		}

		router := newRouter(svc)
		req := httptest.NewRequest(http.MethodGet, "/airports", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d. body=%s", w.Code, w.Body.String())
		}
	})

	t.Run("success - empty list", func(t *testing.T) {
		svc := &fakeAirportService{
			listAirportsFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.Airport, error) {
				return []domain.Airport{}, nil
			},
		}

		router := newRouter(svc)
		req := httptest.NewRequest(http.MethodGet, "/airports", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("success - with status=true filter", func(t *testing.T) {
		svc := &fakeAirportService{
			listAirportsFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.Airport, error) {
				if filters["status"] != true {
					t.Errorf("expected status filter true, got %v", filters["status"])
				}
				return []domain.Airport{}, nil
			},
		}

		router := newRouter(svc)
		req := httptest.NewRequest(http.MethodGet, "/airports?status=true", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("success - with status=false filter", func(t *testing.T) {
		svc := &fakeAirportService{
			listAirportsFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.Airport, error) {
				if filters["status"] != false {
					t.Errorf("expected status filter false, got %v", filters["status"])
				}
				return []domain.Airport{}, nil
			},
		}

		router := newRouter(svc)
		req := httptest.NewRequest(http.MethodGet, "/airports?status=false", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("success - with status=active filter", func(t *testing.T) {
		svc := &fakeAirportService{
			listAirportsFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.Airport, error) {
				if filters["status"] != true {
					t.Errorf("expected status filter true, got %v", filters["status"])
				}
				return []domain.Airport{}, nil
			},
		}

		router := newRouter(svc)
		req := httptest.NewRequest(http.MethodGet, "/airports?status=active", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("success - with status=inactive filter", func(t *testing.T) {
		svc := &fakeAirportService{
			listAirportsFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.Airport, error) {
				if filters["status"] != false {
					t.Errorf("expected status filter false, got %v", filters["status"])
				}
				return []domain.Airport{}, nil
			},
		}

		router := newRouter(svc)
		req := httptest.NewRequest(http.MethodGet, "/airports?status=inactive", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("success - with status=1 filter", func(t *testing.T) {
		svc := &fakeAirportService{
			listAirportsFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.Airport, error) {
				if filters["status"] != true {
					t.Errorf("expected status filter true, got %v", filters["status"])
				}
				return []domain.Airport{}, nil
			},
		}

		router := newRouter(svc)
		req := httptest.NewRequest(http.MethodGet, "/airports?status=1", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("success - with status=0 filter", func(t *testing.T) {
		svc := &fakeAirportService{
			listAirportsFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.Airport, error) {
				if filters["status"] != false {
					t.Errorf("expected status filter false, got %v", filters["status"])
				}
				return []domain.Airport{}, nil
			},
		}

		router := newRouter(svc)
		req := httptest.NewRequest(http.MethodGet, "/airports?status=0", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("service error => 500", func(t *testing.T) {
		svc := &fakeAirportService{
			listAirportsFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.Airport, error) {
				return nil, errors.New("database error")
			},
		}

		router := newRouter(svc)
		req := httptest.NewRequest(http.MethodGet, "/airports", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", w.Code)
		}
	})
}

func TestHTTP_GetAirportsByType(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := newTestAirportMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)

	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	newRouter := func(svc input.AirportService) *gin.Engine {
		airportInteractor := interactor.NewAirportInteractor(svc)
		h := New(HandlerDeps{
		IDEncoder: enc,
		Response: resp,
		AirportInteractor: airportInteractor,
		})

		r := gin.New()
		r.Use(middleware.RequestID())
		r.Use(errHandler.Handle())
		r.GET("/airport-types/:airport_type", h.GetAirportsByType())
		return r
	}

	t.Run("success - returns airports by type", func(t *testing.T) {
		expectedAirports := []domain.Airport{
			{ID: "uuid-1", Name: "El Dorado", IATACode: "BOG", AirportType: "INTERNACIONAL", Status: true},
			{ID: "uuid-2", Name: "Aeropuerto Alfonso Bonilla", IATACode: "CLO", AirportType: "INTERNACIONAL", Status: true},
		}

		svc := &fakeAirportService{
			getByTypeFn: func(ctx context.Context, airportType string) ([]domain.Airport, error) {
				if airportType != "INTERNACIONAL" {
					t.Errorf("expected type INTERNACIONAL, got %s", airportType)
				}
				return expectedAirports, nil
			},
		}

		router := newRouter(svc)
		req := httptest.NewRequest(http.MethodGet, "/airport-types/INTERNACIONAL", nil)
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

	t.Run("success - empty list for type", func(t *testing.T) {
		svc := &fakeAirportService{
			getByTypeFn: func(ctx context.Context, airportType string) ([]domain.Airport, error) {
				return []domain.Airport{}, nil
			},
		}

		router := newRouter(svc)
		req := httptest.NewRequest(http.MethodGet, "/airport-types/NACIONAL", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("service error", func(t *testing.T) {
		svc := &fakeAirportService{
			getByTypeFn: func(ctx context.Context, airportType string) ([]domain.Airport, error) {
				return nil, errors.New("database error")
			},
		}

		router := newRouter(svc)
		req := httptest.NewRequest(http.MethodGet, "/airport-types/INTERNACIONAL", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Should return error response
		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if response["success"] != false {
			t.Errorf("expected success=false, got %v", response["success"])
		}
	})

	t.Run("empty airport_type - returns error", func(t *testing.T) {
		svc := &fakeAirportService{}

		airportInteractor := interactor.NewAirportInteractor(svc)
		h := New(HandlerDeps{
		IDEncoder: enc,
		Response: resp,
		AirportInteractor: airportInteractor,
		})

		r := gin.New()
		r.Use(middleware.RequestID())
		r.Use(errHandler.Handle())
		// Route that matches empty type
		r.GET("/airport-types/", h.GetAirportsByType())

		req := httptest.NewRequest(http.MethodGet, "/airport-types/", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		// Should return error response for empty type
		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if response["success"] != false {
			t.Errorf("expected success=false for empty type, got %v", response["success"])
		}
	})
}

func TestHTTP_GetAirportByID_EdgeCases(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := newTestAirportMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)

	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	newRouter := func(svc input.AirportService) *gin.Engine {
		airportInteractor := interactor.NewAirportInteractor(svc)
		h := New(HandlerDeps{
		IDEncoder: enc,
		Response: resp,
		AirportInteractor: airportInteractor,
		})

		r := gin.New()
		r.Use(middleware.RequestID())
		r.Use(errHandler.Handle())
		r.GET("/airports/:id", h.GetAirportByID())
		return r
	}

	t.Run("invalid ID format => 400", func(t *testing.T) {
		svc := &fakeAirportService{}
		router := newRouter(svc)

		req := httptest.NewRequest(http.MethodGet, "/airports/invalid-id!!!", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d. body=%s", w.Code, w.Body.String())
		}
	})

	t.Run("generic service error => 500", func(t *testing.T) {
		airportUUID := "550e8400-e29b-41d4-a716-446655440000"
		encodedID, _ := enc.Encode(airportUUID)

		svc := &fakeAirportService{
			getByIDFn: func(ctx context.Context, id string) (*domain.Airport, error) {
				return nil, errors.New("database connection failed")
			},
		}

		router := newRouter(svc)
		req := httptest.NewRequest(http.MethodGet, "/airports/"+encodedID, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d. body=%s", w.Code, w.Body.String())
		}
	})
}

func TestHTTP_ActivateAirport(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := newTestAirportMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)

	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	newRouter := func(svc input.AirportService) *gin.Engine {
		airportInteractor := interactor.NewAirportInteractor(svc)
		h := New(HandlerDeps{
		IDEncoder: enc,
		Response: resp,
		AirportInteractor: airportInteractor,
		})

		r := gin.New()
		r.Use(middleware.RequestID())
		r.Use(errHandler.Handle())
		r.PATCH("/airports/:id/activate", h.ActivateAirport())
		return r
	}

	t.Run("success", func(t *testing.T) {
		airportUUID := "550e8400-e29b-41d4-a716-446655440000"
		encodedID, _ := enc.Encode(airportUUID)
		activateCalled := false

		svc := &fakeAirportService{
			getByIDFn: func(ctx context.Context, id string) (*domain.Airport, error) {
				return &domain.Airport{ID: id, Status: false}, nil
			},
			activateFn: func(ctx context.Context, id string) error {
				activateCalled = true
				return nil
			},
		}

		router := newRouter(svc)

		req := httptest.NewRequest(http.MethodPatch, "/airports/"+encodedID+"/activate", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		if !activateCalled {
			t.Error("expected activateFn to be called")
		}

		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if response["success"] != true {
			t.Errorf("expected success=true, got %v", response["success"])
		}
	})

	t.Run("airport not found", func(t *testing.T) {
		airportUUID := "550e8400-e29b-41d4-a716-446655440000"
		encodedID, _ := enc.Encode(airportUUID)

		svc := &fakeAirportService{
			getByIDFn: func(ctx context.Context, id string) (*domain.Airport, error) {
				return &domain.Airport{ID: id}, nil
			},
			activateFn: func(ctx context.Context, id string) error {
				return domain.ErrAirportNotFound
			},
		}

		router := newRouter(svc)

		req := httptest.NewRequest(http.MethodPatch, "/airports/"+encodedID+"/activate", nil)
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

	t.Run("service error", func(t *testing.T) {
		airportUUID := "550e8400-e29b-41d4-a716-446655440000"
		encodedID, _ := enc.Encode(airportUUID)

		svc := &fakeAirportService{
			activateFn: func(ctx context.Context, id string) error {
				return errors.New("activation failed")
			},
		}

		router := newRouter(svc)

		req := httptest.NewRequest(http.MethodPatch, "/airports/"+encodedID+"/activate", nil)
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

	t.Run("invalid ID => 400", func(t *testing.T) {
		svc := &fakeAirportService{}
		router := newRouter(svc)

		req := httptest.NewRequest(http.MethodPatch, "/airports/invalid-id!!!/activate", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}
