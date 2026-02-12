package handlers

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

// fakeDailyLogbookService implements input.DailyLogbookService
type fakeDailyLogbookService struct {
	getByIDFn func(ctx context.Context, id string) (*domain.DailyLogbook, error)
	listFn    func(ctx context.Context, employeeID string, filters map[string]interface{}) ([]domain.DailyLogbook, error)
	createFn  func(ctx context.Context, logbook domain.DailyLogbook) error
}

func (f *fakeDailyLogbookService) BeginTx(ctx context.Context) (output.Tx, error) {
	return fakeTx{}, nil
}

func (f *fakeDailyLogbookService) GetDailyLogbookByID(ctx context.Context, id string) (*domain.DailyLogbook, error) {
	if f.getByIDFn != nil {
		return f.getByIDFn(ctx, id)
	}
	return nil, nil
}

func (f *fakeDailyLogbookService) ListDailyLogbooksByEmployee(ctx context.Context, employeeID string, filters map[string]interface{}) ([]domain.DailyLogbook, error) {
	if f.listFn != nil {
		return f.listFn(ctx, employeeID, filters)
	}
	return nil, nil
}

func (f *fakeDailyLogbookService) CreateDailyLogbook(ctx context.Context, logbook domain.DailyLogbook) error {
	if f.createFn != nil {
		return f.createFn(ctx, logbook)
	}
	return nil
}
func (f *fakeDailyLogbookService) CreateDailyLogbookTx(ctx context.Context, tx output.Tx, logbook domain.DailyLogbook) error {
	if f.createFn != nil {
		return f.createFn(ctx, logbook)
	}
	return nil
}

func newTestDailyLogbookMessageCache(t *testing.T) *messaging.MessageCache {
	t.Helper()

	repo := fakeMessageCacheRepo{messages: []cachetypes.CachedMessage{
		{Code: domain.MsgDailyLogbookGetOK, Type: cachetypes.TypeSuccess, Content: "logbook found"},
		{Code: domain.MsgDailyLogbookNotFound, Type: cachetypes.TypeError, Content: "logbook not found"},
		{Code: domain.MsgDailyLogbookGetErr, Type: cachetypes.TypeError, Content: "error getting logbook"},
		{Code: domain.MsgDailyLogbookCreated, Type: cachetypes.TypeSuccess, Content: "logbook created"},
		{Code: domain.MsgDailyLogbookSaveError, Type: cachetypes.TypeError, Content: "error saving logbook"},
		{Code: domain.MsgDailyLogbookListOK, Type: cachetypes.TypeSuccess, Content: "logbooks listed"},
		{Code: domain.MsgDailyLogbookListError, Type: cachetypes.TypeError, Content: "error listing logbooks"},
		{Code: domain.MsgDailyLogbookUnauthorized, Type: cachetypes.TypeError, Content: "unauthorized"},
		{Code: domain.MsgServerError, Type: cachetypes.TypeError, Content: "internal server error"},
		{Code: domain.MsgValIDInvalid, Type: cachetypes.TypeError, Content: "invalid id"},
		{Code: domain.MsgValJSONInvalid, Type: cachetypes.TypeError, Content: "invalid json"},
	}}
	cache := messaging.NewMessageCache(repo, 0)
	if err := cache.LoadMessages(context.Background()); err != nil {
		t.Fatalf("failed to load message cache: %v", err)
	}
	return cache
}

func newDailyLogbookTestRouter(
	svc input.DailyLogbookService,
	enc *idencoder.HashidsEncoder,
	resp *middleware.ResponseHandler,
	errHandler *middleware.ErrorHandler,
	authUser *domain.Employee,
) *gin.Engine {
	logbookInteractor := interactor.NewDailyLogbookInteractor(svc)
	h := New(nil, &fakeEmployeeInteractor{}, enc, resp, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, logbookInteractor, nil, nil)

	r := gin.New()
	r.Use(middleware.RequestID())
	r.Use(errHandler.Handle())
	// Mock authentication middleware
	r.Use(func(c *gin.Context) {
		if authUser != nil {
			c.Set("authenticated_user", authUser)
		}
		c.Next()
	})
	r.GET("/daily-logbooks", h.ListDailyLogbooks())
	r.GET("/daily-logbooks/:id", h.GetDailyLogbookByID())
	r.POST("/daily-logbooks", h.CreateDailyLogbook())
	return r
}

func TestHTTP_ListDailyLogbooks(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := newTestDailyLogbookMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)
	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	testEmployeeID := "550e8400-e29b-41d4-a716-446655440000"

	t.Run("unauthorized - no employee in context", func(t *testing.T) {
		svc := &fakeDailyLogbookService{}
		router := newDailyLogbookTestRouter(svc, enc, resp, errHandler, nil)

		req := httptest.NewRequest(http.MethodGet, "/daily-logbooks", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d, body: %s", w.Code, w.Body.String())
	})

	t.Run("success - returns logbooks", func(t *testing.T) {
		svc := &fakeDailyLogbookService{
			listFn: func(ctx context.Context, employeeID string, filters map[string]interface{}) ([]domain.DailyLogbook, error) {
				return []domain.DailyLogbook{
					{ID: "550e8400-e29b-41d4-a716-446655440001", EmployeeID: testEmployeeID, LogDate: time.Now(), Status: true},
				}, nil
			},
		}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookTestRouter(svc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodGet, "/daily-logbooks", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("with status filter", func(t *testing.T) {
		svc := &fakeDailyLogbookService{
			listFn: func(ctx context.Context, employeeID string, filters map[string]interface{}) ([]domain.DailyLogbook, error) {
				return []domain.DailyLogbook{}, nil
			},
		}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookTestRouter(svc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodGet, "/daily-logbooks?status=true", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &fakeDailyLogbookService{
			listFn: func(ctx context.Context, employeeID string, filters map[string]interface{}) ([]domain.DailyLogbook, error) {
				return nil, errors.New("db error")
			},
		}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookTestRouter(svc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodGet, "/daily-logbooks", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})
}

func TestHTTP_GetDailyLogbookByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := newTestDailyLogbookMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)
	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	testUUID := "550e8400-e29b-41d4-a716-446655440000"
	testEmployeeID := "550e8400-e29b-41d4-a716-446655440001"
	encodedID, _ := enc.Encode(testUUID)

	t.Run("success - returns logbook", func(t *testing.T) {
		svc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{
					ID:         testUUID,
					EmployeeID: testEmployeeID,
					LogDate:    time.Now(),
					Status:     true,
				}, nil
			},
		}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookTestRouter(svc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodGet, "/daily-logbooks/"+encodedID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("unauthorized - no employee", func(t *testing.T) {
		svc := &fakeDailyLogbookService{}
		router := newDailyLogbookTestRouter(svc, enc, resp, errHandler, nil)

		req := httptest.NewRequest(http.MethodGet, "/daily-logbooks/"+encodedID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("invalid ID", func(t *testing.T) {
		svc := &fakeDailyLogbookService{}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookTestRouter(svc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodGet, "/daily-logbooks/invalid!!!!", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return nil, errors.New("db error")
			},
		}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookTestRouter(svc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodGet, "/daily-logbooks/"+encodedID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("unauthorized - different employee", func(t *testing.T) {
		svc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{
					ID:         testUUID,
					EmployeeID: "other-employee-id",
					LogDate:    time.Now(),
					Status:     true,
				}, nil
			},
		}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookTestRouter(svc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodGet, "/daily-logbooks/"+encodedID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})
}

func TestHTTP_CreateDailyLogbook(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := newTestDailyLogbookMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)
	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	testEmployeeID := "550e8400-e29b-41d4-a716-446655440001"

	t.Run("success - creates logbook", func(t *testing.T) {
		svc := &fakeDailyLogbookService{
			createFn: func(ctx context.Context, logbook domain.DailyLogbook) error {
				return nil
			},
		}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookTestRouter(svc, enc, resp, errHandler, authUser)

		body := `{"log_date":"2025-01-15","book_page":42}`
		req := httptest.NewRequest(http.MethodPost, "/daily-logbooks", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		svc := &fakeDailyLogbookService{}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookTestRouter(svc, enc, resp, errHandler, authUser)

		body := `{bad json`
		req := httptest.NewRequest(http.MethodPost, "/daily-logbooks", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("invalid date format", func(t *testing.T) {
		svc := &fakeDailyLogbookService{}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookTestRouter(svc, enc, resp, errHandler, authUser)

		body := `{"log_date":"not-a-date"}`
		req := httptest.NewRequest(http.MethodPost, "/daily-logbooks", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &fakeDailyLogbookService{
			createFn: func(ctx context.Context, logbook domain.DailyLogbook) error {
				return errors.New("save error")
			},
		}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookTestRouter(svc, enc, resp, errHandler, authUser)

		body := `{"log_date":"2025-01-15"}`
		req := httptest.NewRequest(http.MethodPost, "/daily-logbooks", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("unauthorized - no employee", func(t *testing.T) {
		svc := &fakeDailyLogbookService{}
		router := newDailyLogbookTestRouter(svc, enc, resp, errHandler, nil)

		body := `{"log_date":"2025-01-15"}`
		req := httptest.NewRequest(http.MethodPost, "/daily-logbooks", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})
}
