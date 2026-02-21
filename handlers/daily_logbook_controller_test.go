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
	getByIDFn    func(ctx context.Context, id string) (*domain.DailyLogbook, error)
	listFn       func(ctx context.Context, employeeID string, filters map[string]interface{}) ([]domain.DailyLogbook, error)
	createTxFn   func(ctx context.Context, tx output.Tx, logbook domain.DailyLogbook) error
	updateTxFn   func(ctx context.Context, tx output.Tx, logbook domain.DailyLogbook) error
	activateFn   func(ctx context.Context, tx output.Tx, id string) error
	deactivateFn func(ctx context.Context, tx output.Tx, id string) error
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

func (f *fakeDailyLogbookService) CreateDailyLogbookTx(ctx context.Context, tx output.Tx, logbook domain.DailyLogbook) error {
	if f.createTxFn != nil {
		return f.createTxFn(ctx, tx, logbook)
	}
	return nil
}

func (f *fakeDailyLogbookService) UpdateDailyLogbookTx(ctx context.Context, tx output.Tx, logbook domain.DailyLogbook) error {
	if f.updateTxFn != nil {
		return f.updateTxFn(ctx, tx, logbook)
	}
	return nil
}

func (f *fakeDailyLogbookService) ActivateDailyLogbookTx(ctx context.Context, tx output.Tx, id string) error {
	if f.activateFn != nil {
		return f.activateFn(ctx, tx, id)
	}
	return nil
}

func (f *fakeDailyLogbookService) DeactivateDailyLogbookTx(ctx context.Context, tx output.Tx, id string) error {
	if f.deactivateFn != nil {
		return f.deactivateFn(ctx, tx, id)
	}
	return nil
}

func (f *fakeDailyLogbookService) DeleteDailyLogbookTx(ctx context.Context, tx output.Tx, id string) error {
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
		{Code: domain.MsgDailyLogbookActivateOK, Type: cachetypes.TypeSuccess, Content: "logbook activated"},
		{Code: domain.MsgDailyLogbookActivateErr, Type: cachetypes.TypeError, Content: "error activating logbook"},
		{Code: domain.MsgDailyLogbookDeactivateOK, Type: cachetypes.TypeSuccess, Content: "logbook deactivated"},
		{Code: domain.MsgDailyLogbookDeactivateErr, Type: cachetypes.TypeError, Content: "error deactivating logbook"},
		{Code: domain.MsgDailyLogbookDeleted, Type: cachetypes.TypeSuccess, Content: "logbook deleted"},
		{Code: domain.MsgDailyLogbookDeleteError, Type: cachetypes.TypeError, Content: "error deleting logbook"},
		{Code: domain.MsgDailyLogbookGetErr, Type: cachetypes.TypeError, Content: "error getting logbook"},
		{Code: domain.MsgDailyLogbookUpdated, Type: cachetypes.TypeSuccess, Content: "logbook updated"},
		{Code: domain.MsgDailyLogbookUpdateError, Type: cachetypes.TypeError, Content: "error updating logbook"},
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
	h := New(HandlerDeps{
		EmployeeInteractor: &fakeEmployeeInteractor{},
		IDEncoder: enc,
		Response: resp,
		DailyLogbookInteractor: logbookInteractor,
		})

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
	r.PATCH("/daily-logbooks/:id/activate", h.ActivateDailyLogbook())
	r.PATCH("/daily-logbooks/:id/deactivate", h.DeactivateDailyLogbook())
	r.DELETE("/daily-logbooks/:id", h.DeleteDailyLogbook())
	r.PUT("/daily-logbooks/:id", h.UpdateDailyLogbook())
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
			createTxFn: func(ctx context.Context, tx output.Tx, logbook domain.DailyLogbook) error {
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
			createTxFn: func(ctx context.Context, tx output.Tx, logbook domain.DailyLogbook) error {
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

func TestHTTP_ActivateDailyLogbook(t *testing.T) {
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

	t.Run("success", func(t *testing.T) {
		svc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: testUUID, EmployeeID: testEmployeeID, LogDate: time.Now(), Status: false}, nil
			},
		}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookTestRouter(svc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodPatch, "/daily-logbooks/"+encodedID+"/activate", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		svc := &fakeDailyLogbookService{}
		router := newDailyLogbookTestRouter(svc, enc, resp, errHandler, nil)

		req := httptest.NewRequest(http.MethodPatch, "/daily-logbooks/"+encodedID+"/activate", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("invalid ID", func(t *testing.T) {
		svc := &fakeDailyLogbookService{}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookTestRouter(svc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodPatch, "/daily-logbooks/invalid!!!!/activate", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("logbook not found", func(t *testing.T) {
		svc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return nil, errors.New("not found")
			},
		}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookTestRouter(svc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodPatch, "/daily-logbooks/"+encodedID+"/activate", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("unauthorized - different employee", func(t *testing.T) {
		svc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return nil, domain.ErrFlightUnauthorized
			},
		}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookTestRouter(svc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodPatch, "/daily-logbooks/"+encodedID+"/activate", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("activate service error", func(t *testing.T) {
		svc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: testUUID, EmployeeID: testEmployeeID, LogDate: time.Now(), Status: false}, nil
			},
			activateFn: func(ctx context.Context, tx output.Tx, id string) error {
				return errors.New("activate failed")
			},
		}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookTestRouter(svc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodPatch, "/daily-logbooks/"+encodedID+"/activate", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})
}

func TestHTTP_DeactivateDailyLogbook(t *testing.T) {
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

	t.Run("success", func(t *testing.T) {
		svc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: testUUID, EmployeeID: testEmployeeID, LogDate: time.Now(), Status: true}, nil
			},
		}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookTestRouter(svc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodPatch, "/daily-logbooks/"+encodedID+"/deactivate", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		svc := &fakeDailyLogbookService{}
		router := newDailyLogbookTestRouter(svc, enc, resp, errHandler, nil)

		req := httptest.NewRequest(http.MethodPatch, "/daily-logbooks/"+encodedID+"/deactivate", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("invalid ID", func(t *testing.T) {
		svc := &fakeDailyLogbookService{}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookTestRouter(svc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodPatch, "/daily-logbooks/invalid!!!!/deactivate", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("logbook not found", func(t *testing.T) {
		svc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return nil, errors.New("not found")
			},
		}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookTestRouter(svc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodPatch, "/daily-logbooks/"+encodedID+"/deactivate", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("unauthorized - different employee", func(t *testing.T) {
		svc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return nil, domain.ErrFlightUnauthorized
			},
		}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookTestRouter(svc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodPatch, "/daily-logbooks/"+encodedID+"/deactivate", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("deactivate service error", func(t *testing.T) {
		svc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: testUUID, EmployeeID: testEmployeeID, LogDate: time.Now(), Status: true}, nil
			},
			deactivateFn: func(ctx context.Context, tx output.Tx, id string) error {
				return errors.New("deactivate failed")
			},
		}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookTestRouter(svc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodPatch, "/daily-logbooks/"+encodedID+"/deactivate", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})
}

func TestHTTP_DeleteDailyLogbook(t *testing.T) {
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

	t.Run("success", func(t *testing.T) {
		svc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: testUUID, EmployeeID: testEmployeeID, LogDate: time.Now(), Status: true}, nil
			},
		}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookTestRouter(svc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodDelete, "/daily-logbooks/"+encodedID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		svc := &fakeDailyLogbookService{}
		router := newDailyLogbookTestRouter(svc, enc, resp, errHandler, nil)

		req := httptest.NewRequest(http.MethodDelete, "/daily-logbooks/"+encodedID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("invalid ID", func(t *testing.T) {
		svc := &fakeDailyLogbookService{}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookTestRouter(svc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodDelete, "/daily-logbooks/invalid!!!!", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("logbook not found", func(t *testing.T) {
		svc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return nil, errors.New("not found")
			},
		}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookTestRouter(svc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodDelete, "/daily-logbooks/"+encodedID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("unauthorized - different employee", func(t *testing.T) {
		svc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: testUUID, EmployeeID: "other-employee", LogDate: time.Now(), Status: true}, nil
			},
		}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookTestRouter(svc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodDelete, "/daily-logbooks/"+encodedID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})
}

func TestHTTP_UpdateDailyLogbook(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := newTestDailyLogbookMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)
	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	testEmployeeID := "550e8400-e29b-41d4-a716-446655440000"
	testUUID := "660e8400-e29b-41d4-a716-446655440099"
	encodedID, _ := enc.Encode(testUUID)

	t.Run("success - update daily logbook", func(t *testing.T) {
		svc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: testUUID, EmployeeID: testEmployeeID, LogDate: time.Now(), Status: true}, nil
			},
		}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookTestRouter(svc, enc, resp, errHandler, authUser)

		body := `{"log_date":"2025-06-15","status":true}`
		req := httptest.NewRequest(http.MethodPut, "/daily-logbooks/"+encodedID, bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
		t.Logf("status: %d, body: %s", w.Code, w.Body.String())
	})

	t.Run("unauthorized - no employee in context", func(t *testing.T) {
		svc := &fakeDailyLogbookService{}
		router := newDailyLogbookTestRouter(svc, enc, resp, errHandler, nil)

		body := `{"log_date":"2025-06-15","status":true}`
		req := httptest.NewRequest(http.MethodPut, "/daily-logbooks/"+encodedID, bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			t.Errorf("expected non-200 status, got %d", w.Code)
		}
		t.Logf("status: %d, body: %s", w.Code, w.Body.String())
	})

	t.Run("error - invalid ID", func(t *testing.T) {
		svc := &fakeDailyLogbookService{}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookTestRouter(svc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodPut, "/daily-logbooks/invalid-id", bytes.NewBufferString(`{}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			t.Errorf("expected non-200 status, got %d", w.Code)
		}
		t.Logf("status: %d, body: %s", w.Code, w.Body.String())
	})

	t.Run("error - invalid JSON body", func(t *testing.T) {
		svc := &fakeDailyLogbookService{}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookTestRouter(svc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodPut, "/daily-logbooks/"+encodedID, bytes.NewBufferString(`{invalid`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			t.Errorf("expected non-200 status, got %d", w.Code)
		}
		t.Logf("status: %d, body: %s", w.Code, w.Body.String())
	})

	t.Run("error - interactor returns error", func(t *testing.T) {
		svc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: testUUID, EmployeeID: testEmployeeID, LogDate: time.Now(), Status: true}, nil
			},
			updateTxFn: func(ctx context.Context, tx output.Tx, logbook domain.DailyLogbook) error {
				return errors.New("update failed")
			},
		}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookTestRouter(svc, enc, resp, errHandler, authUser)

		body := `{"log_date":"2025-06-15","status":true}`
		req := httptest.NewRequest(http.MethodPut, "/daily-logbooks/"+encodedID, bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			t.Errorf("expected non-200 status, got %d", w.Code)
		}
		t.Logf("status: %d, body: %s", w.Code, w.Body.String())
	})

	t.Run("unauthorized - different employee", func(t *testing.T) {
		svc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: testUUID, EmployeeID: "other-employee", LogDate: time.Now(), Status: true}, nil
			},
		}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookTestRouter(svc, enc, resp, errHandler, authUser)

		body := `{"log_date":"2025-06-15","status":true}`
		req := httptest.NewRequest(http.MethodPut, "/daily-logbooks/"+encodedID, bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			t.Errorf("expected non-200 status, got %d", w.Code)
		}
		t.Logf("status: %d, body: %s", w.Code, w.Body.String())
	})
}
