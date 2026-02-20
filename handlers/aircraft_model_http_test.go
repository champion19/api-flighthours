package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/champion19/api-flighthours/core/interactor"
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/input"
	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/middleware"
	"github.com/champion19/api-flighthours/tools/idencoder"
	"github.com/gin-gonic/gin"
)

// fakeAircraftModelServiceForHandler implements input.AircraftModelService
type fakeAircraftModelServiceForHandler struct {
	getByIDFn     func(ctx context.Context, id string) (*domain.AircraftModel, error)
	listFn        func(ctx context.Context, filters map[string]interface{}) ([]domain.AircraftModel, error)
	getByFamilyFn func(ctx context.Context, family string) ([]domain.AircraftModel, error)
	activateFn    func(ctx context.Context, id string) error
	deactivateFn  func(ctx context.Context, id string) error
}

var _ input.AircraftModelService = (*fakeAircraftModelServiceForHandler)(nil)

func (f *fakeAircraftModelServiceForHandler) BeginTx(ctx context.Context) (output.Tx, error) {
	return fakeTx{}, nil
}
func (f *fakeAircraftModelServiceForHandler) GetAircraftModelByID(ctx context.Context, id string) (*domain.AircraftModel, error) {
	if f.getByIDFn != nil {
		return f.getByIDFn(ctx, id)
	}
	return nil, errors.New("not implemented")
}
func (f *fakeAircraftModelServiceForHandler) ListAircraftModels(ctx context.Context, filters map[string]interface{}) ([]domain.AircraftModel, error) {
	if f.listFn != nil {
		return f.listFn(ctx, filters)
	}
	return nil, errors.New("not implemented")
}
func (f *fakeAircraftModelServiceForHandler) GetAircraftModelsByFamily(ctx context.Context, family string) ([]domain.AircraftModel, error) {
	if f.getByFamilyFn != nil {
		return f.getByFamilyFn(ctx, family)
	}
	return nil, errors.New("not implemented")
}
func (f *fakeAircraftModelServiceForHandler) ActivateAircraftModel(ctx context.Context, id string) error {
	if f.activateFn != nil {
		return f.activateFn(ctx, id)
	}
	return nil
}
func (f *fakeAircraftModelServiceForHandler) DeactivateAircraftModel(ctx context.Context, id string) error {
	if f.deactivateFn != nil {
		return f.deactivateFn(ctx, id)
	}
	return nil
}
func (f *fakeAircraftModelServiceForHandler) ActivateAircraftModelTx(ctx context.Context, tx output.Tx, id string) error {
	if f.activateFn != nil {
		return f.activateFn(ctx, id)
	}
	return nil
}
func (f *fakeAircraftModelServiceForHandler) DeactivateAircraftModelTx(ctx context.Context, tx output.Tx, id string) error {
	if f.deactivateFn != nil {
		return f.deactivateFn(ctx, id)
	}
	return nil
}

func newAircraftModelTestRouter(svc input.AircraftModelService, enc *idencoder.HashidsEncoder, resp *middleware.ResponseHandler, errHandler *middleware.ErrorHandler) *gin.Engine {
	aircraftModelInteractor := interactor.NewAircraftModelInteractor(svc)
	h := New(nil, &fakeEmployeeInteractor{}, enc, resp, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, aircraftModelInteractor, nil)

	r := gin.New()
	r.Use(middleware.RequestID())
	r.Use(errHandler.Handle())
	r.GET("/aircraft-models/:id", h.GetAircraftModelByID())
	r.GET("/aircraft-models", h.ListAircraftModels())
	r.GET("/aircraft-families/:family", h.GetAircraftModelsByFamily())
	r.PATCH("/aircraft-models/:id/activate", h.ActivateAircraftModel())
	r.PATCH("/aircraft-models/:id/deactivate", h.DeactivateAircraftModel())
	return r
}

func TestHTTP_GetAircraftModelByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := newTestEmployeeMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)
	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	t.Run("success - returns aircraft model", func(t *testing.T) {
		testUUID := "550e8400-e29b-41d4-a716-446655440000"
		encodedID, _ := enc.Encode(testUUID)

		svc := &fakeAircraftModelServiceForHandler{
			getByIDFn: func(ctx context.Context, id string) (*domain.AircraftModel, error) {
				return &domain.AircraftModel{
					ID:               testUUID,
					ModelName:        "Boeing 737-800",
					AircraftTypeName: "Narrow Body",
					EngineTypeName:   "JET",
					Family:           "737",
					Manufacturer:     "Boeing",
					Status:           true,
				}, nil
			},
		}

		router := newAircraftModelTestRouter(svc, enc, resp, errHandler)
		req := httptest.NewRequest(http.MethodGet, "/aircraft-models/"+encodedID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d. body=%s", w.Code, w.Body.String())
		}
	})

	t.Run("not found - returns error", func(t *testing.T) {
		testUUID := "550e8400-e29b-41d4-a716-446655440000"
		encodedID, _ := enc.Encode(testUUID)

		svc := &fakeAircraftModelServiceForHandler{
			getByIDFn: func(ctx context.Context, id string) (*domain.AircraftModel, error) {
				return nil, domain.ErrAircraftModelNotFound
			},
		}

		router := newAircraftModelTestRouter(svc, enc, resp, errHandler)
		req := httptest.NewRequest(http.MethodGet, "/aircraft-models/"+encodedID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			t.Errorf("expected error status, got 200")
		}
	})

	t.Run("service error - returns 500", func(t *testing.T) {
		testUUID := "550e8400-e29b-41d4-a716-446655440000"
		encodedID, _ := enc.Encode(testUUID)

		svc := &fakeAircraftModelServiceForHandler{
			getByIDFn: func(ctx context.Context, id string) (*domain.AircraftModel, error) {
				return nil, errors.New("db connection error")
			},
		}

		router := newAircraftModelTestRouter(svc, enc, resp, errHandler)
		req := httptest.NewRequest(http.MethodGet, "/aircraft-models/"+encodedID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			t.Errorf("expected error status, got 200")
		}
	})

	t.Run("invalid ID - returns error", func(t *testing.T) {
		svc := &fakeAircraftModelServiceForHandler{}
		router := newAircraftModelTestRouter(svc, enc, resp, errHandler)
		req := httptest.NewRequest(http.MethodGet, "/aircraft-models/invalid!!!", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			t.Errorf("expected error status, got 200")
		}
	})
}

func TestHTTP_ListAircraftModels(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := newTestEmployeeMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)
	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	t.Run("success - returns list", func(t *testing.T) {
		svc := &fakeAircraftModelServiceForHandler{
			listFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.AircraftModel, error) {
				return []domain.AircraftModel{
					{ID: "550e8400-e29b-41d4-a716-446655440001", ModelName: "737-800", Family: "737", Status: true},
					{ID: "550e8400-e29b-41d4-a716-446655440002", ModelName: "A320", Family: "A320", Status: true},
				}, nil
			},
		}

		router := newAircraftModelTestRouter(svc, enc, resp, errHandler)
		req := httptest.NewRequest(http.MethodGet, "/aircraft-models", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d. body=%s", w.Code, w.Body.String())
		}
	})

	t.Run("error - returns error", func(t *testing.T) {
		svc := &fakeAircraftModelServiceForHandler{
			listFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.AircraftModel, error) {
				return nil, errors.New("db error")
			},
		}

		router := newAircraftModelTestRouter(svc, enc, resp, errHandler)
		req := httptest.NewRequest(http.MethodGet, "/aircraft-models", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			t.Errorf("expected error status, got 200")
		}
	})
}

func TestHTTP_GetAircraftModelsByFamily(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := newTestEmployeeMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)
	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	t.Run("success - returns models for family", func(t *testing.T) {
		svc := &fakeAircraftModelServiceForHandler{
			getByFamilyFn: func(ctx context.Context, family string) ([]domain.AircraftModel, error) {
				return []domain.AircraftModel{
					{ID: "550e8400-e29b-41d4-a716-446655440001", ModelName: "737-700", Family: "737"},
					{ID: "550e8400-e29b-41d4-a716-446655440002", ModelName: "737-800", Family: "737"},
				}, nil
			},
		}

		router := newAircraftModelTestRouter(svc, enc, resp, errHandler)
		req := httptest.NewRequest(http.MethodGet, "/aircraft-families/737", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d. body=%s", w.Code, w.Body.String())
		}
	})

	t.Run("empty family - returns not found", func(t *testing.T) {
		svc := &fakeAircraftModelServiceForHandler{
			getByFamilyFn: func(ctx context.Context, family string) ([]domain.AircraftModel, error) {
				return []domain.AircraftModel{}, nil
			},
		}

		router := newAircraftModelTestRouter(svc, enc, resp, errHandler)
		req := httptest.NewRequest(http.MethodGet, "/aircraft-families/nonexistent", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			t.Errorf("expected error status for empty family, got 200")
		}
	})

	t.Run("service error - returns error", func(t *testing.T) {
		svc := &fakeAircraftModelServiceForHandler{
			getByFamilyFn: func(ctx context.Context, family string) ([]domain.AircraftModel, error) {
				return nil, errors.New("db error")
			},
		}

		router := newAircraftModelTestRouter(svc, enc, resp, errHandler)
		req := httptest.NewRequest(http.MethodGet, "/aircraft-families/737", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			t.Errorf("expected error status, got 200")
		}
	})
}

func TestHTTP_ActivateAircraftModel(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := newTestEmployeeMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)
	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	t.Run("success - activates model", func(t *testing.T) {
		testUUID := "550e8400-e29b-41d4-a716-446655440000"
		encodedID, _ := enc.Encode(testUUID)

		svc := &fakeAircraftModelServiceForHandler{
			getByIDFn: func(ctx context.Context, id string) (*domain.AircraftModel, error) {
				return &domain.AircraftModel{ID: id, ModelName: "737-800"}, nil
			},
			activateFn: func(ctx context.Context, id string) error {
				return nil
			},
		}

		router := newAircraftModelTestRouter(svc, enc, resp, errHandler)
		req := httptest.NewRequest(http.MethodPatch, "/aircraft-models/"+encodedID+"/activate", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d. body=%s", w.Code, w.Body.String())
		}
	})

	t.Run("not found - returns error", func(t *testing.T) {
		testUUID := "550e8400-e29b-41d4-a716-446655440000"
		encodedID, _ := enc.Encode(testUUID)

		svc := &fakeAircraftModelServiceForHandler{
			getByIDFn: func(ctx context.Context, id string) (*domain.AircraftModel, error) {
				return nil, domain.ErrAircraftModelNotFound
			},
		}

		router := newAircraftModelTestRouter(svc, enc, resp, errHandler)
		req := httptest.NewRequest(http.MethodPatch, "/aircraft-models/"+encodedID+"/activate", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			t.Errorf("expected error status, got 200")
		}
	})

	t.Run("service error - returns error", func(t *testing.T) {
		testUUID := "550e8400-e29b-41d4-a716-446655440000"
		encodedID, _ := enc.Encode(testUUID)

		svc := &fakeAircraftModelServiceForHandler{
			activateFn: func(ctx context.Context, id string) error {
				return errors.New("activation failed")
			},
		}

		router := newAircraftModelTestRouter(svc, enc, resp, errHandler)
		req := httptest.NewRequest(http.MethodPatch, "/aircraft-models/"+encodedID+"/activate", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			t.Errorf("expected error status, got 200")
		}
	})

	t.Run("invalid ID - returns error", func(t *testing.T) {
		svc := &fakeAircraftModelServiceForHandler{}
		router := newAircraftModelTestRouter(svc, enc, resp, errHandler)
		req := httptest.NewRequest(http.MethodPatch, "/aircraft-models/invalid!!!/activate", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			t.Errorf("expected error status, got 200")
		}
	})
}

func TestHTTP_DeactivateAircraftModel(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := newTestEmployeeMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)
	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	t.Run("success - deactivates model", func(t *testing.T) {
		testUUID := "550e8400-e29b-41d4-a716-446655440000"
		encodedID, _ := enc.Encode(testUUID)

		svc := &fakeAircraftModelServiceForHandler{
			getByIDFn: func(ctx context.Context, id string) (*domain.AircraftModel, error) {
				return &domain.AircraftModel{ID: id, ModelName: "737-800"}, nil
			},
			deactivateFn: func(ctx context.Context, id string) error {
				return nil
			},
		}

		router := newAircraftModelTestRouter(svc, enc, resp, errHandler)
		req := httptest.NewRequest(http.MethodPatch, "/aircraft-models/"+encodedID+"/deactivate", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d. body=%s", w.Code, w.Body.String())
		}
	})

	t.Run("not found - returns error", func(t *testing.T) {
		testUUID := "550e8400-e29b-41d4-a716-446655440000"
		encodedID, _ := enc.Encode(testUUID)

		svc := &fakeAircraftModelServiceForHandler{
			getByIDFn: func(ctx context.Context, id string) (*domain.AircraftModel, error) {
				return nil, domain.ErrAircraftModelNotFound
			},
		}

		router := newAircraftModelTestRouter(svc, enc, resp, errHandler)
		req := httptest.NewRequest(http.MethodPatch, "/aircraft-models/"+encodedID+"/deactivate", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			t.Errorf("expected error status, got 200")
		}
	})

	t.Run("service error - returns error", func(t *testing.T) {
		testUUID := "550e8400-e29b-41d4-a716-446655440000"
		encodedID, _ := enc.Encode(testUUID)

		svc := &fakeAircraftModelServiceForHandler{
			deactivateFn: func(ctx context.Context, id string) error {
				return errors.New("deactivation failed")
			},
		}

		router := newAircraftModelTestRouter(svc, enc, resp, errHandler)
		req := httptest.NewRequest(http.MethodPatch, "/aircraft-models/"+encodedID+"/deactivate", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			t.Errorf("expected error status, got 200")
		}
	})

	t.Run("invalid ID - returns error", func(t *testing.T) {
		svc := &fakeAircraftModelServiceForHandler{}
		router := newAircraftModelTestRouter(svc, enc, resp, errHandler)
		req := httptest.NewRequest(http.MethodPatch, "/aircraft-models/invalid!!!/deactivate", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			t.Errorf("expected error status, got 200")
		}
	})
}
