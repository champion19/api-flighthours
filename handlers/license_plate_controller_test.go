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

// fakeLicensePlateService implements input.LicensePlateService
type fakeLicensePlateService struct {
	getByIDFn func(ctx context.Context, id string) (*domain.LicensePlate, error)
	listFn    func(ctx context.Context, filters map[string]interface{}) ([]domain.LicensePlate, error)
	createFn  func(ctx context.Context, registration domain.LicensePlate) error
	updateFn  func(ctx context.Context, registration domain.LicensePlate) error
}

var _ input.LicensePlateService = (*fakeLicensePlateService)(nil)

func (f *fakeLicensePlateService) BeginTx(ctx context.Context) (output.Tx, error) {
	return fakeTx{}, nil
}
func (f *fakeLicensePlateService) GetLicensePlateByID(ctx context.Context, id string) (*domain.LicensePlate, error) {
	if f.getByIDFn != nil {
		return f.getByIDFn(ctx, id)
	}
	return nil, errors.New("not implemented")
}
func (f *fakeLicensePlateService) ListLicensePlates(ctx context.Context, filters map[string]interface{}) ([]domain.LicensePlate, error) {
	if f.listFn != nil {
		return f.listFn(ctx, filters)
	}
	return nil, errors.New("not implemented")
}
func (f *fakeLicensePlateService) CreateLicensePlate(ctx context.Context, registration domain.LicensePlate) error {
	if f.createFn != nil {
		return f.createFn(ctx, registration)
	}
	return nil
}
func (f *fakeLicensePlateService) UpdateLicensePlate(ctx context.Context, registration domain.LicensePlate) error {
	if f.updateFn != nil {
		return f.updateFn(ctx, registration)
	}
	return nil
}

func (f *fakeLicensePlateService) GetLicensePlateByPlate(ctx context.Context, plate string) (*domain.LicensePlate, error) {
	if f.getByIDFn != nil {
		return f.getByIDFn(ctx, plate)
	}
	return nil, errors.New("not implemented")
}

func newTestLicensePlateMessageCache(t *testing.T) *messaging.MessageCache {
	t.Helper()

	repo := fakeMessageCacheRepo{messages: []cachetypes.CachedMessage{
		{Code: domain.MsgLicensePlateGetOK, Type: cachetypes.TypeSuccess, Content: "license plate found"},
		{Code: domain.MsgLicensePlateNotFound, Type: cachetypes.TypeError, Content: "license plate not found"},
		{Code: domain.MsgLicensePlateListOK, Type: cachetypes.TypeSuccess, Content: "license plates listed"},
		{Code: domain.MsgLicensePlateCreated, Type: cachetypes.TypeSuccess, Content: "license plate created"},
		{Code: domain.MsgLicensePlateUpdated, Type: cachetypes.TypeSuccess, Content: "license plate updated"},
		{Code: domain.MsgServerError, Type: cachetypes.TypeError, Content: "internal server error"},
		{Code: domain.MsgValIDInvalid, Type: cachetypes.TypeError, Content: "invalid id"},
	}}
	cache := messaging.NewMessageCache(repo, 0)
	if err := cache.LoadMessages(context.Background()); err != nil {
		t.Fatalf("failed to load message cache: %v", err)
	}
	return cache
}

func newLicensePlateTestRouter(svc input.LicensePlateService, enc *idencoder.HashidsEncoder, resp *middleware.ResponseHandler, errHandler *middleware.ErrorHandler) *gin.Engine {
	licensePlateInteractor := interactor.NewLicensePlateInteractor(svc, Logger)
	h := New(nil, &fakeEmployeeInteractor{}, enc, resp, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, licensePlateInteractor)

	r := gin.New()
	r.Use(middleware.RequestID())
	r.Use(errHandler.Handle())
	r.GET("/license-plates/:plate", h.GetLicensePlateByPlate())
	r.GET("/license-plates", h.ListLicensePlates())
	r.POST("/license-plates", h.CreateLicensePlate())
	r.PUT("/license-plates/:id", h.UpdateLicensePlate())
	return r
}

func TestHTTP_GetLicensePlateByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := newTestLicensePlateMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)
	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	t.Run("success - returns license plate", func(t *testing.T) {
		testUUID := "550e8400-e29b-41d4-a716-446655440000"
		modelUUID := "550e8400-e29b-41d4-a716-446655440001"
		airlineUUID := "550e8400-e29b-41d4-a716-446655440002"
		encodedID, _ := enc.Encode(testUUID)

		svc := &fakeLicensePlateService{
			getByIDFn: func(ctx context.Context, id string) (*domain.LicensePlate, error) {
				return &domain.LicensePlate{
					ID:              testUUID,
					LicensePlate:    "HK-5432",
					AircraftModelID: modelUUID,
					AirlineID:       airlineUUID,
					ModelName:       "Boeing 737",
					AirlineName:     "Avianca",
				}, nil
			},
		}

		router := newLicensePlateTestRouter(svc, enc, resp, errHandler)
		req := httptest.NewRequest(http.MethodGet, "/license-plates/"+encodedID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d. body=%s", w.Code, w.Body.String())
		}
	})

	t.Run("not found - returns error", func(t *testing.T) {
		testUUID := "550e8400-e29b-41d4-a716-446655440000"
		encodedID, _ := enc.Encode(testUUID)

		svc := &fakeLicensePlateService{
			getByIDFn: func(ctx context.Context, id string) (*domain.LicensePlate, error) {
				return nil, domain.ErrLicensePlateNotFound
			},
		}

		router := newLicensePlateTestRouter(svc, enc, resp, errHandler)
		req := httptest.NewRequest(http.MethodGet, "/license-plates/"+encodedID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			t.Errorf("expected error status, got 200")
		}
	})
}

func TestHTTP_ListLicensePlates(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := newTestLicensePlateMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)
	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	t.Run("success - returns list", func(t *testing.T) {
		svc := &fakeLicensePlateService{
			listFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.LicensePlate, error) {
				return []domain.LicensePlate{
					{ID: "550e8400-e29b-41d4-a716-446655440001", LicensePlate: "HK-5432", AircraftModelID: "m1", AirlineID: "a1"},
					{ID: "550e8400-e29b-41d4-a716-446655440002", LicensePlate: "CC-BFA", AircraftModelID: "m2", AirlineID: "a2"},
				}, nil
			},
		}

		router := newLicensePlateTestRouter(svc, enc, resp, errHandler)
		req := httptest.NewRequest(http.MethodGet, "/license-plates", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d. body=%s", w.Code, w.Body.String())
		}
	})

	t.Run("error - returns error", func(t *testing.T) {
		svc := &fakeLicensePlateService{
			listFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.LicensePlate, error) {
				return nil, errors.New("db error")
			},
		}

		router := newLicensePlateTestRouter(svc, enc, resp, errHandler)
		req := httptest.NewRequest(http.MethodGet, "/license-plates", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			t.Errorf("expected error status, got 200")
		}
	})
}

func TestHTTP_CreateLicensePlate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := newTestLicensePlateMessageCache(t)
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

		svc := &fakeLicensePlateService{
			createFn: func(ctx context.Context, registration domain.LicensePlate) error {
				return nil
			},
			getByIDFn: func(ctx context.Context, id string) (*domain.LicensePlate, error) {
				return &domain.LicensePlate{
					ID:              id,
					LicensePlate:    "HK-5432",
					AircraftModelID: modelUUID,
					AirlineID:       airlineUUID,
					ModelName:       "Boeing 737",
					AirlineName:     "Avianca",
				}, nil
			},
		}

		router := newLicensePlateTestRouter(svc, enc, resp, errHandler)
		body := `{"license_plate":"HK-5432","aircraft_model_id":"` + encodedModelID + `","airline_id":"` + encodedAirlineID + `"}`
		req := httptest.NewRequest(http.MethodPost, "/license-plates", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected 201, got %d. body=%s", w.Code, w.Body.String())
		}
	})

	t.Run("invalid JSON - returns error", func(t *testing.T) {
		svc := &fakeLicensePlateService{}

		router := newLicensePlateTestRouter(svc, enc, resp, errHandler)
		body := `not json`
		req := httptest.NewRequest(http.MethodPost, "/license-plates", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code == http.StatusCreated {
			t.Errorf("expected error status, got 201")
		}
	})

	t.Run("invalid model ID - returns error", func(t *testing.T) {
		svc := &fakeLicensePlateService{}

		router := newLicensePlateTestRouter(svc, enc, resp, errHandler)
		body := `{"license_plate":"HK-5432","aircraft_model_id":"invalid-model","airline_id":"some-airline"}`
		req := httptest.NewRequest(http.MethodPost, "/license-plates", bytes.NewBufferString(body))
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

		svc := &fakeLicensePlateService{}

		router := newLicensePlateTestRouter(svc, enc, resp, errHandler)
		body := `{"license_plate":"HK-5432","aircraft_model_id":"` + encodedModelID + `","airline_id":"invalid-airline"}`
		req := httptest.NewRequest(http.MethodPost, "/license-plates", bytes.NewBufferString(body))
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

		svc := &fakeLicensePlateService{
			createFn: func(ctx context.Context, registration domain.LicensePlate) error {
				return domain.ErrLicensePlateDuplicatePlate
			},
		}

		router := newLicensePlateTestRouter(svc, enc, resp, errHandler)
		body := `{"license_plate":"HK-5432","aircraft_model_id":"` + encodedModelID + `","airline_id":"` + encodedAirlineID + `"}`
		req := httptest.NewRequest(http.MethodPost, "/license-plates", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code == http.StatusCreated {
			t.Errorf("expected error status for duplicate plate, got 201")
		}
	})
}

func TestHTTP_UpdateLicensePlate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := newTestLicensePlateMessageCache(t)
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

		svc := &fakeLicensePlateService{
			getByIDFn: func(ctx context.Context, id string) (*domain.LicensePlate, error) {
				return &domain.LicensePlate{
					ID:              testUUID,
					LicensePlate:    "HK-OLD",
					AircraftModelID: modelUUID,
					AirlineID:       airlineUUID,
					ModelName:       "Boeing 737",
					AirlineName:     "Avianca",
				}, nil
			},
			updateFn: func(ctx context.Context, registration domain.LicensePlate) error {
				return nil
			},
		}

		router := newLicensePlateTestRouter(svc, enc, resp, errHandler)
		body := `{"license_plate":"HK-NEW","aircraft_model_id":"` + encodedModelID + `","airline_id":"` + encodedAirlineID + `"}`
		req := httptest.NewRequest(http.MethodPut, "/license-plates/"+encodedID, bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d. body=%s", w.Code, w.Body.String())
		}
	})

	t.Run("not found - returns error", func(t *testing.T) {
		svc := &fakeLicensePlateService{
			listFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.LicensePlate, error) {
				return []domain.LicensePlate{}, nil
			},
		}

		router := newLicensePlateTestRouter(svc, enc, resp, errHandler)
		body := `{"license_plate":"HK-NEW","aircraft_model_id":"model","airline_id":"airline"}`
		req := httptest.NewRequest(http.MethodPut, "/license-plates/invalid-encoded-id", bytes.NewBufferString(body))
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

		svc := &fakeLicensePlateService{
			getByIDFn: func(ctx context.Context, id string) (*domain.LicensePlate, error) {
				return &domain.LicensePlate{
					ID:              testUUID,
					LicensePlate:    "HK-OLD",
					AircraftModelID: modelUUID,
					AirlineID:       airlineUUID,
				}, nil
			},
			updateFn: func(ctx context.Context, registration domain.LicensePlate) error {
				return domain.ErrLicensePlateDuplicatePlate
			},
		}

		router := newLicensePlateTestRouter(svc, enc, resp, errHandler)
		body := `{"license_plate":"HK-DUP","aircraft_model_id":"` + encodedModelID + `","airline_id":"` + encodedAirlineID + `"}`
		req := httptest.NewRequest(http.MethodPut, "/license-plates/"+encodedID, bytes.NewBufferString(body))
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

		svc := &fakeLicensePlateService{
			getByIDFn: func(ctx context.Context, id string) (*domain.LicensePlate, error) {
				return &domain.LicensePlate{ID: testUUID}, nil
			},
		}

		router := newLicensePlateTestRouter(svc, enc, resp, errHandler)
		body := `not json`
		req := httptest.NewRequest(http.MethodPut, "/license-plates/"+encodedID, bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			t.Errorf("expected error status for invalid JSON, got 200")
		}
	})
}
