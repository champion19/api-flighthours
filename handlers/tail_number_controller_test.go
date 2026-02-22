package handlers

import (
	"bytes"
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
	"github.com/champion19/api-flighthours/platform/cache/messaging"
	cachetypes "github.com/champion19/api-flighthours/platform/cache/types"
	"github.com/champion19/api-flighthours/tools/idencoder"
	"github.com/gin-gonic/gin"
)

// fakeTailNumberService implements input.TailNumberService
type fakeTailNumberService struct {
	getByIDFn    func(ctx context.Context, id string) (*domain.TailNumber, error)
	getByPlateFn func(ctx context.Context, plate string) (*domain.TailNumber, error)
	listFn       func(ctx context.Context, filters map[string]interface{}) ([]domain.TailNumber, error)
	createFn     func(ctx context.Context, registration domain.TailNumber) error
	updateFn     func(ctx context.Context, registration domain.TailNumber) error
}

var _ input.TailNumberService = (*fakeTailNumberService)(nil)

func (f *fakeTailNumberService) BeginTx(ctx context.Context) (output.Tx, error) {
	return fakeTx{}, nil
}
func (f *fakeTailNumberService) GetTailNumberByID(ctx context.Context, id string) (*domain.TailNumber, error) {
	if f.getByIDFn != nil {
		return f.getByIDFn(ctx, id)
	}
	return nil, errors.New("not implemented")
}
func (f *fakeTailNumberService) ListTailNumbers(ctx context.Context, filters map[string]interface{}) ([]domain.TailNumber, error) {
	if f.listFn != nil {
		return f.listFn(ctx, filters)
	}
	return nil, errors.New("not implemented")
}
func (f *fakeTailNumberService) CreateTailNumber(ctx context.Context, registration domain.TailNumber) error {
	if f.createFn != nil {
		return f.createFn(ctx, registration)
	}
	return nil
}
func (f *fakeTailNumberService) UpdateTailNumber(ctx context.Context, registration domain.TailNumber) error {
	if f.updateFn != nil {
		return f.updateFn(ctx, registration)
	}
	return nil
}

func (f *fakeTailNumberService) GetTailNumberByPlate(ctx context.Context, plate string) (*domain.TailNumber, error) {
	if f.getByPlateFn != nil {
		return f.getByPlateFn(ctx, plate)
	}
	return nil, errors.New("not implemented")
}
func (f *fakeTailNumberService) CreateTailNumberTx(ctx context.Context, tx output.Tx, registration domain.TailNumber) error {
	if f.createFn != nil {
		return f.createFn(ctx, registration)
	}
	return nil
}
func (f *fakeTailNumberService) UpdateTailNumberTx(ctx context.Context, tx output.Tx, registration domain.TailNumber) error {
	if f.updateFn != nil {
		return f.updateFn(ctx, registration)
	}
	return nil
}

func newTestTailNumberMessageCache(t *testing.T) *messaging.MessageCache {
	t.Helper()

	repo := fakeMessageCacheRepo{messages: []cachetypes.CachedMessage{
		{Code: domain.MsgTailNumberGetOK, Type: cachetypes.TypeSuccess, Content: "license plate found"},
		{Code: domain.MsgTailNumberNotFound, Type: cachetypes.TypeError, Content: "license plate not found"},
		{Code: domain.MsgTailNumberListOK, Type: cachetypes.TypeSuccess, Content: "license plates listed"},
		{Code: domain.MsgTailNumberCreated, Type: cachetypes.TypeSuccess, Content: "license plate created"},
		{Code: domain.MsgTailNumberUpdated, Type: cachetypes.TypeSuccess, Content: "license plate updated"},
		{Code: domain.MsgServerError, Type: cachetypes.TypeError, Content: "internal server error"},
		{Code: domain.MsgValIDInvalid, Type: cachetypes.TypeError, Content: "invalid id"},
	}}
	cache := messaging.NewMessageCache(repo, 0)
	if err := cache.LoadMessages(context.Background()); err != nil {
		t.Fatalf("failed to load message cache: %v", err)
	}
	return cache
}

func newTailNumberTestRouter(svc input.TailNumberService, enc *idencoder.HashidsEncoder, resp *middleware.ResponseHandler, errHandler *middleware.ErrorHandler) *gin.Engine {
	tailNumberInteractor := interactor.NewTailNumberInteractor(svc, log)
	h := New(HandlerDeps{
		EmployeeInteractor: &fakeEmployeeInteractor{},
		IDEncoder: enc,
		Response: resp,
		TailNumberInteractor: tailNumberInteractor,
		})

	r := gin.New()
	r.Use(middleware.RequestID())
	r.Use(errHandler.Handle())
	r.GET("/tail-numbers/:plate", h.GetTailNumberByPlate())
	r.GET("/tail-numbers", h.ListTailNumbers())
	r.POST("/tail-numbers", h.CreateTailNumber())
	r.PUT("/tail-numbers/:id", h.UpdateTailNumber())
	return r
}

func TestHTTP_GetTailNumberByPlate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := newTestTailNumberMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)
	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	t.Run("success - returns license plate by plate number", func(t *testing.T) {
		testUUID := "550e8400-e29b-41d4-a716-446655440000"
		modelUUID := "550e8400-e29b-41d4-a716-446655440001"
		airlineUUID := "550e8400-e29b-41d4-a716-446655440002"

		svc := &fakeTailNumberService{
			getByPlateFn: func(ctx context.Context, plate string) (*domain.TailNumber, error) {
				return &domain.TailNumber{
					ID:              testUUID,
					TailNumber:    "HK-5432",
					AircraftModelID: modelUUID,
					AirlineID:       airlineUUID,
					ModelName:       "Boeing 737",
					AirlineName:     "Avianca",
				}, nil
			},
		}

		router := newTailNumberTestRouter(svc, enc, resp, errHandler)
		req := httptest.NewRequest(http.MethodGet, "/tail-numbers/HK-5432", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d. body=%s", w.Code, w.Body.String())
		}
	})

	t.Run("not found - returns error", func(t *testing.T) {
		svc := &fakeTailNumberService{
			getByPlateFn: func(ctx context.Context, plate string) (*domain.TailNumber, error) {
				return nil, domain.ErrTailNumberNotFound
			},
		}

		router := newTailNumberTestRouter(svc, enc, resp, errHandler)
		req := httptest.NewRequest(http.MethodGet, "/tail-numbers/INVALID", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			t.Errorf("expected error status, got 200")
		}
	})

	t.Run("service error - returns error", func(t *testing.T) {
		svc := &fakeTailNumberService{
			getByPlateFn: func(ctx context.Context, plate string) (*domain.TailNumber, error) {
				return nil, errors.New("database connection error")
			},
		}

		router := newTailNumberTestRouter(svc, enc, resp, errHandler)
		req := httptest.NewRequest(http.MethodGet, "/tail-numbers/HK-9999", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			t.Errorf("expected error status, got 200")
		}
	})
}

func TestHTTP_ListTailNumbers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := newTestTailNumberMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)
	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	t.Run("success - returns list", func(t *testing.T) {
		svc := &fakeTailNumberService{
			listFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.TailNumber, error) {
				return []domain.TailNumber{
					{ID: "550e8400-e29b-41d4-a716-446655440001", TailNumber: "HK-5432", AircraftModelID: "m1", AirlineID: "a1"},
					{ID: "550e8400-e29b-41d4-a716-446655440002", TailNumber: "CC-BFA", AircraftModelID: "m2", AirlineID: "a2"},
				}, nil
			},
		}

		router := newTailNumberTestRouter(svc, enc, resp, errHandler)
		req := httptest.NewRequest(http.MethodGet, "/tail-numbers", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d. body=%s", w.Code, w.Body.String())
		}
	})

	t.Run("error - returns error", func(t *testing.T) {
		svc := &fakeTailNumberService{
			listFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.TailNumber, error) {
				return nil, errors.New("db error")
			},
		}

		router := newTailNumberTestRouter(svc, enc, resp, errHandler)
		req := httptest.NewRequest(http.MethodGet, "/tail-numbers", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			t.Errorf("expected error status, got 200")
		}
	})

	t.Run("success - with tail_number filter", func(t *testing.T) {
		svc := &fakeTailNumberService{
			listFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.TailNumber, error) {
				if filters["tail_number"] != "HK-5432" {
					t.Errorf("expected tail_number filter HK-5432, got %v", filters["tail_number"])
				}
				return []domain.TailNumber{
					{ID: "uuid-1", TailNumber: "HK-5432", AircraftModelID: "m1", AirlineID: "a1"},
				}, nil
			},
		}

		router := newTailNumberTestRouter(svc, enc, resp, errHandler)
		req := httptest.NewRequest(http.MethodGet, "/tail-numbers?tail_number=HK-5432", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d. body=%s", w.Code, w.Body.String())
		}
	})

	t.Run("success - with airline_id filter", func(t *testing.T) {
		airlineUUID := "550e8400-e29b-41d4-a716-446655440010"
		encodedAirlineID, _ := enc.Encode(airlineUUID)

		svc := &fakeTailNumberService{
			listFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.TailNumber, error) {
				if _, ok := filters["airline_id"]; !ok {
					t.Error("expected airline_id filter to be set")
				}
				return []domain.TailNumber{
					{ID: "uuid-1", TailNumber: "HK-5432", AircraftModelID: "m1", AirlineID: airlineUUID},
				}, nil
			},
		}

		router := newTailNumberTestRouter(svc, enc, resp, errHandler)
		req := httptest.NewRequest(http.MethodGet, "/tail-numbers?airline_id="+encodedAirlineID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d. body=%s", w.Code, w.Body.String())
		}
	})
}

func TestHTTP_CreateTailNumber(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := newTestTailNumberMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)
	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	t.Run("success - creates license plate", func(t *testing.T) {
		modelUUID := "550e8400-e29b-41d4-a716-446655440001"
		airlineUUID := "550e8400-e29b-41d4-a716-446655440002"
		encodedModelID, _ := enc.Encode(modelUUID)
		encodedAirlineID, _ := enc.Encode(airlineUUID)

		svc := &fakeTailNumberService{
			createFn: func(ctx context.Context, registration domain.TailNumber) error {
				return nil
			},
			getByIDFn: func(ctx context.Context, id string) (*domain.TailNumber, error) {
				return &domain.TailNumber{
					ID:              id,
					TailNumber:    "HK-5432",
					AircraftModelID: modelUUID,
					AirlineID:       airlineUUID,
					ModelName:       "Boeing 737",
					AirlineName:     "Avianca",
				}, nil
			},
		}

		router := newTailNumberTestRouter(svc, enc, resp, errHandler)
		body := `{"tail_number":"HK-5432","aircraft_model_id":"` + encodedModelID + `","airline_id":"` + encodedAirlineID + `"}`
		req := httptest.NewRequest(http.MethodPost, "/tail-numbers", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected 201, got %d. body=%s", w.Code, w.Body.String())
		}
	})

	t.Run("invalid JSON - returns error", func(t *testing.T) {
		svc := &fakeTailNumberService{}

		router := newTailNumberTestRouter(svc, enc, resp, errHandler)
		body := `not json`
		req := httptest.NewRequest(http.MethodPost, "/tail-numbers", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code == http.StatusCreated {
			t.Errorf("expected error status, got 201")
		}
	})

	t.Run("invalid model ID - returns error", func(t *testing.T) {
		svc := &fakeTailNumberService{}

		router := newTailNumberTestRouter(svc, enc, resp, errHandler)
		body := `{"tail_number":"HK-5432","aircraft_model_id":"invalid-model","airline_id":"some-airline"}`
		req := httptest.NewRequest(http.MethodPost, "/tail-numbers", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code == http.StatusCreated {
			t.Errorf("expected error status for invalid model ID, got 201")
		}
	})

	t.Run("invalid airline ID - returns error", func(t *testing.T) {
		modelUUID := "550e8400-e29b-41d4-a716-446655440001"
		encodedModelID, _ := enc.Encode(modelUUID)

		svc := &fakeTailNumberService{}

		router := newTailNumberTestRouter(svc, enc, resp, errHandler)
		body := `{"tail_number":"HK-5432","aircraft_model_id":"` + encodedModelID + `","airline_id":"invalid-airline"}`
		req := httptest.NewRequest(http.MethodPost, "/tail-numbers", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code == http.StatusCreated {
			t.Errorf("expected error status for invalid airline ID, got 201")
		}
	})

	t.Run("service error - returns error", func(t *testing.T) {
		modelUUID := "550e8400-e29b-41d4-a716-446655440001"
		airlineUUID := "550e8400-e29b-41d4-a716-446655440002"
		encodedModelID, _ := enc.Encode(modelUUID)
		encodedAirlineID, _ := enc.Encode(airlineUUID)

		svc := &fakeTailNumberService{
			createFn: func(ctx context.Context, registration domain.TailNumber) error {
				return domain.ErrTailNumberDuplicatePlate
			},
		}

		router := newTailNumberTestRouter(svc, enc, resp, errHandler)
		body := `{"tail_number":"HK-5432","aircraft_model_id":"` + encodedModelID + `","airline_id":"` + encodedAirlineID + `"}`
		req := httptest.NewRequest(http.MethodPost, "/tail-numbers", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code == http.StatusCreated {
			t.Errorf("expected error status for duplicate plate, got 201")
		}
	})
}

func TestHTTP_UpdateTailNumber(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := newTestTailNumberMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)
	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	t.Run("success - updates license plate", func(t *testing.T) {
		testUUID := "550e8400-e29b-41d4-a716-446655440000"
		modelUUID := "550e8400-e29b-41d4-a716-446655440001"
		airlineUUID := "550e8400-e29b-41d4-a716-446655440002"
		encodedID, _ := enc.Encode(testUUID)
		encodedModelID, _ := enc.Encode(modelUUID)
		encodedAirlineID, _ := enc.Encode(airlineUUID)

		svc := &fakeTailNumberService{
			getByIDFn: func(ctx context.Context, id string) (*domain.TailNumber, error) {
				return &domain.TailNumber{
					ID:              testUUID,
					TailNumber:    "HK-OLD",
					AircraftModelID: modelUUID,
					AirlineID:       airlineUUID,
					ModelName:       "Boeing 737",
					AirlineName:     "Avianca",
				}, nil
			},
			updateFn: func(ctx context.Context, registration domain.TailNumber) error {
				return nil
			},
		}

		router := newTailNumberTestRouter(svc, enc, resp, errHandler)
		body := `{"tail_number":"HK-NEW","aircraft_model_id":"` + encodedModelID + `","airline_id":"` + encodedAirlineID + `"}`
		req := httptest.NewRequest(http.MethodPut, "/tail-numbers/"+encodedID, bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d. body=%s", w.Code, w.Body.String())
		}
	})

	t.Run("not found - returns error", func(t *testing.T) {
		svc := &fakeTailNumberService{
			listFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.TailNumber, error) {
				return []domain.TailNumber{}, nil
			},
		}

		router := newTailNumberTestRouter(svc, enc, resp, errHandler)
		body := `{"tail_number":"HK-NEW","aircraft_model_id":"model","airline_id":"airline"}`
		req := httptest.NewRequest(http.MethodPut, "/tail-numbers/invalid-encoded-id", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			t.Errorf("expected error status, got 200")
		}
	})

	t.Run("update error - duplicate plate", func(t *testing.T) {
		testUUID := "550e8400-e29b-41d4-a716-446655440000"
		modelUUID := "550e8400-e29b-41d4-a716-446655440001"
		airlineUUID := "550e8400-e29b-41d4-a716-446655440002"
		encodedID, _ := enc.Encode(testUUID)
		encodedModelID, _ := enc.Encode(modelUUID)
		encodedAirlineID, _ := enc.Encode(airlineUUID)

		svc := &fakeTailNumberService{
			getByIDFn: func(ctx context.Context, id string) (*domain.TailNumber, error) {
				return &domain.TailNumber{
					ID:              testUUID,
					TailNumber:    "HK-OLD",
					AircraftModelID: modelUUID,
					AirlineID:       airlineUUID,
				}, nil
			},
			updateFn: func(ctx context.Context, registration domain.TailNumber) error {
				return domain.ErrTailNumberDuplicatePlate
			},
		}

		router := newTailNumberTestRouter(svc, enc, resp, errHandler)
		body := `{"tail_number":"HK-DUP","aircraft_model_id":"` + encodedModelID + `","airline_id":"` + encodedAirlineID + `"}`
		req := httptest.NewRequest(http.MethodPut, "/tail-numbers/"+encodedID, bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			t.Errorf("expected error status for duplicate plate, got 200")
		}
	})

	t.Run("invalid JSON - returns error", func(t *testing.T) {
		testUUID := "550e8400-e29b-41d4-a716-446655440000"
		encodedID, _ := enc.Encode(testUUID)

		svc := &fakeTailNumberService{
			getByIDFn: func(ctx context.Context, id string) (*domain.TailNumber, error) {
				return &domain.TailNumber{ID: testUUID}, nil
			},
		}

		router := newTailNumberTestRouter(svc, enc, resp, errHandler)
		body := `not json`
		req := httptest.NewRequest(http.MethodPut, "/tail-numbers/"+encodedID, bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			t.Errorf("expected error status for invalid JSON, got 200")
		}
	})

	t.Run("invalid model ID in body", func(t *testing.T) {
		testUUID := "550e8400-e29b-41d4-a716-446655440000"
		encodedID, _ := enc.Encode(testUUID)
		airlineUUID := "550e8400-e29b-41d4-a716-446655440002"
		encodedAirlineID, _ := enc.Encode(airlineUUID)

		svc := &fakeTailNumberService{
			getByIDFn: func(ctx context.Context, id string) (*domain.TailNumber, error) {
				return &domain.TailNumber{ID: testUUID}, nil
			},
		}

		router := newTailNumberTestRouter(svc, enc, resp, errHandler)
		body := `{"tail_number":"HK-NEW","aircraft_model_id":"invalid!!!","airline_id":"` + encodedAirlineID + `"}`
		req := httptest.NewRequest(http.MethodPut, "/tail-numbers/"+encodedID, bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			t.Errorf("expected error status for invalid model ID, got 200")
		}
	})

	t.Run("invalid airline ID in body", func(t *testing.T) {
		testUUID := "550e8400-e29b-41d4-a716-446655440000"
		modelUUID := "550e8400-e29b-41d4-a716-446655440001"
		encodedID, _ := enc.Encode(testUUID)
		encodedModelID, _ := enc.Encode(modelUUID)

		svc := &fakeTailNumberService{
			getByIDFn: func(ctx context.Context, id string) (*domain.TailNumber, error) {
				return &domain.TailNumber{ID: testUUID}, nil
			},
		}

		router := newTailNumberTestRouter(svc, enc, resp, errHandler)
		body := `{"tail_number":"HK-NEW","aircraft_model_id":"` + encodedModelID + `","airline_id":"invalid!!!"}`
		req := httptest.NewRequest(http.MethodPut, "/tail-numbers/"+encodedID, bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			t.Errorf("expected error status for invalid airline ID, got 200")
		}
	})

	t.Run("identical data - duplicate error", func(t *testing.T) {
		testUUID := "550e8400-e29b-41d4-a716-446655440000"
		modelUUID := "550e8400-e29b-41d4-a716-446655440001"
		airlineUUID := "550e8400-e29b-41d4-a716-446655440002"
		encodedID, _ := enc.Encode(testUUID)
		encodedModelID, _ := enc.Encode(modelUUID)
		encodedAirlineID, _ := enc.Encode(airlineUUID)

		svc := &fakeTailNumberService{
			getByIDFn: func(ctx context.Context, id string) (*domain.TailNumber, error) {
				// Return data identical to what will be sent in the body
				return &domain.TailNumber{
					ID:              testUUID,
					TailNumber:    "HK-SAME",
					AircraftModelID: modelUUID,
					AirlineID:       airlineUUID,
				}, nil
			},
		}

		router := newTailNumberTestRouter(svc, enc, resp, errHandler)
		body := `{"tail_number":"HK-SAME","aircraft_model_id":"` + encodedModelID + `","airline_id":"` + encodedAirlineID + `"}`
		req := httptest.NewRequest(http.MethodPut, "/tail-numbers/"+encodedID, bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			t.Errorf("expected error status for duplicate data, got 200")
		}
	})

	t.Run("update service error", func(t *testing.T) {
		testUUID := "550e8400-e29b-41d4-a716-446655440000"
		modelUUID := "550e8400-e29b-41d4-a716-446655440001"
		airlineUUID := "550e8400-e29b-41d4-a716-446655440002"
		encodedID, _ := enc.Encode(testUUID)
		encodedModelID, _ := enc.Encode(modelUUID)
		encodedAirlineID, _ := enc.Encode(airlineUUID)

		svc := &fakeTailNumberService{
			getByIDFn: func(ctx context.Context, id string) (*domain.TailNumber, error) {
				return &domain.TailNumber{
					ID:              testUUID,
					TailNumber:    "HK-OLD",
					AircraftModelID: modelUUID,
					AirlineID:       airlineUUID,
				}, nil
			},
			updateFn: func(ctx context.Context, registration domain.TailNumber) error {
				return errors.New("db write error")
			},
		}

		router := newTailNumberTestRouter(svc, enc, resp, errHandler)
		body := `{"tail_number":"HK-NEW","aircraft_model_id":"` + encodedModelID + `","airline_id":"` + encodedAirlineID + `"}`
		req := httptest.NewRequest(http.MethodPut, "/tail-numbers/"+encodedID, bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			t.Errorf("expected error status for update error, got 200")
		}
	})
}
