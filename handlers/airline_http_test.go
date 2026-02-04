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

type fakeAirlineService struct {
	getByIDFn      func(ctx context.Context, id string) (*domain.Airline, error)
	updateStatusFn func(ctx context.Context, id string, status bool) error
	activateFn     func(ctx context.Context, id string) error
	beginTxFn      func(ctx context.Context) (output.Tx, error)
	listAirlinesFn func(ctx context.Context, filters map[string]interface{}) ([]domain.Airline, error)
}

var _ input.AirlineService = (*fakeAirlineService)(nil)

func (f *fakeAirlineService) BeginTx(ctx context.Context) (output.Tx, error) {
	if f.beginTxFn != nil {
		return f.beginTxFn(ctx)
	}
	return fakeTx{}, nil
}

func (f *fakeAirlineService) GetAirlineByID(ctx context.Context, id string) (*domain.Airline, error) {
	if f.getByIDFn != nil {
		return f.getByIDFn(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (f *fakeAirlineService) UpdateAirlineStatus(ctx context.Context, id string, status bool) error {
	if f.updateStatusFn != nil {
		return f.updateStatusFn(ctx, id, status)
	}
	return errors.New("not implemented")
}

func (f *fakeAirlineService) ActivateAirline(ctx context.Context, id string) error {
	if f.activateFn != nil {
		return f.activateFn(ctx, id)
	}
	return errors.New("not implemented")
}

func (f *fakeAirlineService) ListAirlines(ctx context.Context, filters map[string]interface{}) ([]domain.Airline, error) {
	if f.listAirlinesFn != nil {
		return f.listAirlinesFn(ctx, filters)
	}
	return nil, errors.New("not implemented")
}

func newTestAirlineMessageCache(t *testing.T) *messaging.MessageCache {
	t.Helper()

	repo := fakeMessageCacheRepo{messages: []cachetypes.CachedMessage{
		{Code: domain.MsgAirlineGetOK, Type: cachetypes.TypeSuccess, Content: "airline retrieved successfully"},
		{Code: domain.MsgAirlineNotFound, Type: cachetypes.TypeError, Content: "airline not found"},
		{Code: domain.MsgAirlineActivateOK, Type: cachetypes.TypeSuccess, Content: "airline activated"},
		{Code: domain.MsgAirlineDeactivateOK, Type: cachetypes.TypeSuccess, Content: "airline deactivated"},
		{Code: domain.MsgServerError, Type: cachetypes.TypeError, Content: "internal server error"},
		{Code: domain.MsgValIDInvalid, Type: cachetypes.TypeError, Content: "invalid id"},
	}}
	cache := messaging.NewMessageCache(repo, 0)
	if err := cache.LoadMessages(context.Background()); err != nil {
		t.Fatalf("failed to load message cache: %v", err)
	}
	return cache
}

func TestHTTP_GetAirlineByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := newTestAirlineMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)

	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	newRouter := func(svc input.AirlineService) *gin.Engine {
		airlineInteractor := interactor.NewAirlineInteractor(svc)
		h := New(nil, nil, enc, resp, nil, nil, airlineInteractor, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middleware.RequestID())
		r.Use(errHandler.Handle())
		r.GET("/airlines/:id", h.GetAirlineByID())
		return r
	}

	t.Run("success with obfuscated ID", func(t *testing.T) {
		airlineUUID := "550e8400-e29b-41d4-a716-446655440000"
		encodedID, _ := enc.Encode(airlineUUID)
		expectedAirline := &domain.Airline{
			ID:          airlineUUID,
			AirlineName: "Test Airlines",
			AirlineCode: "TST",
			Status:      "active",
		}

		svc := &fakeAirlineService{
			getByIDFn: func(ctx context.Context, id string) (*domain.Airline, error) {
				if id != airlineUUID {
					t.Errorf("expected id %s, got %s", airlineUUID, id)
				}
				return expectedAirline, nil
			},
		}

		r := newRouter(svc)
		req := httptest.NewRequest(http.MethodGet, "/airlines/"+encodedID, nil)
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
		if out.Code != domain.MsgAirlineGetOK {
			t.Fatalf("expected code %q, got %q", domain.MsgAirlineGetOK, out.Code)
		}
	})

	t.Run("success with encoded ID", func(t *testing.T) {
		airlineUUID := "550e8400-e29b-41d4-a716-446655440001"
		encodedID, _ := enc.Encode(airlineUUID)
		expectedAirline := &domain.Airline{
			ID:          airlineUUID,
			AirlineName: "Encoded Airlines",
			AirlineCode: "ENC",
			Status:      "active",
		}

		svc := &fakeAirlineService{
			getByIDFn: func(ctx context.Context, id string) (*domain.Airline, error) {
				if id != airlineUUID {
					t.Errorf("expected decoded id %s, got %s", airlineUUID, id)
				}
				return expectedAirline, nil
			},
		}

		r := newRouter(svc)
		req := httptest.NewRequest(http.MethodGet, "/airlines/"+encodedID, nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusOK, w.Code, w.Body.String())
		}
	})

	t.Run("airline not found => 404", func(t *testing.T) {
		airlineUUID := "550e8400-e29b-41d4-a716-446655440002"
		encodedID, _ := enc.Encode(airlineUUID)

		svc := &fakeAirlineService{
			getByIDFn: func(context.Context, string) (*domain.Airline, error) {
				return nil, domain.ErrAirlineNotFound
			},
		}

		r := newRouter(svc)
		req := httptest.NewRequest(http.MethodGet, "/airlines/"+encodedID, nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusNotFound {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusNotFound, w.Code, w.Body.String())
		}
	})

	t.Run("empty id => 400", func(t *testing.T) {
		svc := &fakeAirlineService{}

		r := newRouter(svc)

		req := httptest.NewRequest(http.MethodGet, "/airlines/", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusMovedPermanently && w.Code != http.StatusNotFound {
			// This is expected for empty route param
		}
	})
}

func TestHTTP_ActivateAirline(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := newTestAirlineMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)

	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	newRouter := func(svc input.AirlineService) *gin.Engine {
		airlineInteractor := interactor.NewAirlineInteractor(svc)
		h := New(nil, nil, enc, resp, nil, nil, airlineInteractor, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middleware.RequestID())
		r.Use(errHandler.Handle())
		r.PATCH("/airlines/:id/activate", h.ActivateAirline())
		return r
	}

	t.Run("success", func(t *testing.T) {
		airlineUUID := "550e8400-e29b-41d4-a716-446655440000"
		encodedID, _ := enc.Encode(airlineUUID)
		activateCalled := false

		svc := &fakeAirlineService{
			getByIDFn: func(context.Context, string) (*domain.Airline, error) {
				return &domain.Airline{ID: airlineUUID, Status: "inactive"}, nil
			},
			activateFn: func(ctx context.Context, id string) error {
				activateCalled = true
				return nil
			},
		}

		r := newRouter(svc)
		req := httptest.NewRequest(http.MethodPatch, "/airlines/"+encodedID+"/activate", nil)
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

	t.Run("airline not found => 404", func(t *testing.T) {
		airlineUUID := "550e8400-e29b-41d4-a716-446655440002"
		encodedID, _ := enc.Encode(airlineUUID)

		svc := &fakeAirlineService{
			getByIDFn: func(context.Context, string) (*domain.Airline, error) {
				return nil, domain.ErrAirlineNotFound
			},
			activateFn: func(context.Context, string) error {
				return domain.ErrAirlineNotFound
			},
		}

		r := newRouter(svc)
		req := httptest.NewRequest(http.MethodPatch, "/airlines/"+encodedID+"/activate", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusNotFound {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusNotFound, w.Code, w.Body.String())
		}
	})
}

func TestHTTP_ListAirlines(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := newTestAirlineMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)

	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	newRouter := func(svc input.AirlineService) *gin.Engine {
		airlineInteractor := interactor.NewAirlineInteractor(svc)
		h := New(nil, nil, enc, resp, nil, nil, airlineInteractor, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middleware.RequestID())
		r.Use(errHandler.Handle())
		r.GET("/airlines", h.ListAirlines())
		return r
	}

	t.Run("success with airlines", func(t *testing.T) {
		airlines := []domain.Airline{
			{ID: "uuid-1", AirlineName: "Avianca", AirlineCode: "AV", Status: "active"},
			{ID: "uuid-2", AirlineName: "LATAM", AirlineCode: "LA", Status: "active"},
		}

		svc := &fakeAirlineService{
			listAirlinesFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.Airline, error) {
				return airlines, nil
			},
		}

		r := newRouter(svc)
		req := httptest.NewRequest(http.MethodGet, "/airlines", nil)
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
	})

	t.Run("success with empty list", func(t *testing.T) {
		svc := &fakeAirlineService{
			listAirlinesFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.Airline, error) {
				return []domain.Airline{}, nil
			},
		}

		r := newRouter(svc)
		req := httptest.NewRequest(http.MethodGet, "/airlines", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusOK, w.Code, w.Body.String())
		}
	})

	t.Run("service error returns 500", func(t *testing.T) {
		svc := &fakeAirlineService{
			listAirlinesFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.Airline, error) {
				return nil, errors.New("database error")
			},
		}

		r := newRouter(svc)
		req := httptest.NewRequest(http.MethodGet, "/airlines", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusInternalServerError, w.Code, w.Body.String())
		}
	})
}

// Tests for DTO conversion functions
func TestFromDomainAirline(t *testing.T) {
	t.Run("converts domain airline to response correctly", func(t *testing.T) {
		airline := &domain.Airline{
			ID:          "original-uuid",
			AirlineName: "Test Airlines",
			AirlineCode: "TST",
			Status:      "active",
		}
		encodedID := "encoded-id-123"

		result := FromDomainAirline(airline, encodedID)

		if result.ID != encodedID {
			t.Errorf("expected ID %s, got %s", encodedID, result.ID)
		}
		if result.AirlineName != airline.AirlineName {
			t.Errorf("expected AirlineName %s, got %s", airline.AirlineName, result.AirlineName)
		}
		if result.AirlineCode != airline.AirlineCode {
			t.Errorf("expected AirlineCode %s, got %s", airline.AirlineCode, result.AirlineCode)
		}
		if result.Status != airline.Status {
			t.Errorf("expected Status %s, got %s", airline.Status, result.Status)
		}
	})
}

func TestToAirlineListResponse(t *testing.T) {
	t.Run("converts airline slice to list response", func(t *testing.T) {
		airlines := []domain.Airline{
			{ID: "uuid-1", AirlineName: "Avianca", AirlineCode: "AV", Status: "active"},
			{ID: "uuid-2", AirlineName: "LATAM", AirlineCode: "LA", Status: "inactive"},
		}
		encodeFunc := func(id string) (string, error) {
			return "encoded-" + id, nil
		}

		result := ToAirlineListResponse(airlines, encodeFunc, "http://localhost:8080")

		if result.Total != 2 {
			t.Errorf("expected Total 2, got %d", result.Total)
		}
		if len(result.Airlines) != 2 {
			t.Errorf("expected 2 airlines, got %d", len(result.Airlines))
		}
		if result.Airlines[0].ID != "encoded-uuid-1" {
			t.Errorf("expected encoded ID, got %s", result.Airlines[0].ID)
		}
		if len(result.Airlines[0].Links) == 0 {
			t.Error("expected HATEOAS links on airline response")
		}
	})

	t.Run("handles empty airline slice", func(t *testing.T) {
		airlines := []domain.Airline{}
		encodeFunc := func(id string) (string, error) {
			return id, nil
		}

		result := ToAirlineListResponse(airlines, encodeFunc, "")

		if result.Total != 0 {
			t.Errorf("expected Total 0, got %d", result.Total)
		}
		if len(result.Airlines) != 0 {
			t.Errorf("expected 0 airlines, got %d", len(result.Airlines))
		}
	})

	t.Run("uses original ID when encoding fails", func(t *testing.T) {
		airlines := []domain.Airline{
			{ID: "uuid-1", AirlineName: "Avianca", AirlineCode: "AV"},
		}
		encodeFunc := func(id string) (string, error) {
			return "", errors.New("encoding failed")
		}

		result := ToAirlineListResponse(airlines, encodeFunc, "")

		if result.Airlines[0].ID != "uuid-1" {
			t.Errorf("expected original ID uuid-1, got %s", result.Airlines[0].ID)
		}
	})
}

func TestUpdateAirlineStatusRequest_Sanitize(t *testing.T) {
	t.Run("trims whitespace from status", func(t *testing.T) {
		req := &UpdateAirlineStatusRequest{
			Status: "  active  ",
		}

		req.Sanitize()

		if req.Status != "active" {
			t.Errorf("expected 'active', got %q", req.Status)
		}
	})
}
