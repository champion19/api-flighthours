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
	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/middleware"
	"github.com/champion19/api-flighthours/platform/cache/messaging"
	cachetypes "github.com/champion19/api-flighthours/platform/cache/types"
	"github.com/champion19/api-flighthours/tools/idencoder"
	"github.com/gin-gonic/gin"
)

// ── Fake AirlineEmployeeService ────────────────────────────────────────

type fakeAirlineEmployeeService struct {
	getByIDFn      func(ctx context.Context, id string) (*domain.AirlineEmployee, error)
	addFn          func(ctx context.Context, employee domain.AirlineEmployee) error
	updateFn       func(ctx context.Context, employee domain.AirlineEmployee) error
	activateAEFn   func(ctx context.Context, id string) error
	deactivateAEFn func(ctx context.Context, id string) error
}

func (f *fakeAirlineEmployeeService) BeginTx(ctx context.Context) (output.Tx, error) {
	return fakeTx{}, nil
}
func (f *fakeAirlineEmployeeService) GetAirlineEmployeeByID(ctx context.Context, id string) (*domain.AirlineEmployee, error) {
	if f.getByIDFn != nil {
		return f.getByIDFn(ctx, id)
	}
	return nil, nil
}
func (f *fakeAirlineEmployeeService) AddAirlineEmployee(ctx context.Context, employee domain.AirlineEmployee) error {
	if f.addFn != nil {
		return f.addFn(ctx, employee)
	}
	return nil
}
func (f *fakeAirlineEmployeeService) UpdateAirlineEmployee(ctx context.Context, employee domain.AirlineEmployee) error {
	if f.updateFn != nil {
		return f.updateFn(ctx, employee)
	}
	return nil
}
func (f *fakeAirlineEmployeeService) ActivateAirlineEmployee(ctx context.Context, id string) error {
	if f.activateAEFn != nil {
		return f.activateAEFn(ctx, id)
	}
	return nil
}
func (f *fakeAirlineEmployeeService) DeactivateAirlineEmployee(ctx context.Context, id string) error {
	if f.deactivateAEFn != nil {
		return f.deactivateAEFn(ctx, id)
	}
	return nil
}
func (f *fakeAirlineEmployeeService) AddAirlineEmployeeTx(ctx context.Context, tx output.Tx, employee domain.AirlineEmployee) error {
	if f.addFn != nil {
		return f.addFn(ctx, employee)
	}
	return nil
}
func (f *fakeAirlineEmployeeService) UpdateAirlineEmployeeTx(ctx context.Context, tx output.Tx, employee domain.AirlineEmployee) error {
	if f.updateFn != nil {
		return f.updateFn(ctx, employee)
	}
	return nil
}
func (f *fakeAirlineEmployeeService) ActivateAirlineEmployeeTx(ctx context.Context, tx output.Tx, id string) error {
	if f.activateAEFn != nil {
		return f.activateAEFn(ctx, id)
	}
	return nil
}
func (f *fakeAirlineEmployeeService) DeactivateAirlineEmployeeTx(ctx context.Context, tx output.Tx, id string) error {
	if f.deactivateAEFn != nil {
		return f.deactivateAEFn(ctx, id)
	}
	return nil
}

// fakeAirlineService is already defined in airline_http_test.go

// ── Test setup ─────────────────────────────────────────────────────────

func newAirlineEmployeeMessageCache(t *testing.T) *messaging.MessageCache {
	t.Helper()
	repo := fakeMessageCacheRepo{messages: []cachetypes.CachedMessage{
		{Code: domain.MsgAirlineEmployeeGetOK, Type: cachetypes.TypeSuccess, Content: "ok"},
		{Code: domain.MsgAirlineEmployeeNotFound, Type: cachetypes.TypeError, Content: "not found"},
		{Code: domain.MsgAirlineEmployeeGetErr, Type: cachetypes.TypeError, Content: "get error"},
		{Code: domain.MsgAirlineEmployeeCreated, Type: cachetypes.TypeSuccess, Content: "created"},
		{Code: domain.MsgAirlineEmployeeSaveError, Type: cachetypes.TypeError, Content: "save err"},
		{Code: domain.MsgAirlineEmployeeDuplicate, Type: cachetypes.TypeError, Content: "duplicate"},
		{Code: domain.MsgAirlineEmployeeUpdated, Type: cachetypes.TypeSuccess, Content: "updated"},
		{Code: domain.MsgAirlineEmployeeUpdateError, Type: cachetypes.TypeError, Content: "update err"},
		{Code: domain.MsgAirlineEmployeeActivateOK, Type: cachetypes.TypeSuccess, Content: "activated"},
		{Code: domain.MsgAirlineEmployeeActivateErr, Type: cachetypes.TypeError, Content: "activate err"},
		{Code: domain.MsgAirlineEmployeeDeactivateOK, Type: cachetypes.TypeSuccess, Content: "deactivated"},
		{Code: domain.MsgAirlineEmployeeDeactivateErr, Type: cachetypes.TypeError, Content: "deactivate err"},
		{Code: domain.MsgAirlineEmployeeInvalidAirline, Type: cachetypes.TypeError, Content: "invalid airline"},
		{Code: domain.MsgUnauthorized, Type: cachetypes.TypeError, Content: "unauthorized"},
		{Code: domain.MsgValInvalidReq, Type: cachetypes.TypeError, Content: "invalid req"},
		{Code: domain.MsgValInvalidDateFormat, Type: cachetypes.TypeError, Content: "invalid date"},
		{Code: domain.MsgServerError, Type: cachetypes.TypeError, Content: "server error"},
	}}
	cache := messaging.NewMessageCache(repo, 0)
	if err := cache.LoadMessages(context.Background()); err != nil {
		t.Fatalf("failed to load cache: %v", err)
	}
	return cache
}

func newAirlineEmployeeTestRouter(
	aeService *fakeAirlineEmployeeService,
	airlineSvc *fakeAirlineService,
	enc *idencoder.HashidsEncoder,
	resp *middleware.ResponseHandler,
	errHandler *middleware.ErrorHandler,
	authUser *domain.Employee,
) *gin.Engine {
	aeInteractor := interactor.NewAirlineEmployeeInteractor(aeService)
	airlineInteractor := interactor.NewAirlineInteractor(airlineSvc)
	h := New(HandlerDeps{
		EmployeeInteractor: &fakeEmployeeInteractor{},
		IDEncoder: enc,
		Response: resp,
		AirlineInteractor: airlineInteractor,
		AirlineEmployeeInteractor: aeInteractor,
		})

	r := gin.New()
	r.Use(middleware.RequestID())
	r.Use(errHandler.Handle())
	r.Use(func(c *gin.Context) {
		if authUser != nil {
			c.Set("authenticated_user", authUser)
		}
		c.Next()
	})
	r.GET("/employees/airline", h.GetEmployeeAirlineInfo())
	r.PUT("/employees/airline", h.AddEmployeeAirlineInfo())
	r.PUT("/employees/airline-info", h.UpdateEmployeeAirlineInfo())
	r.PATCH("/employees/airline/activate", h.ActivateEmployeeAirlineInfo())
	r.PATCH("/employees/airline/deactivate", h.DeactivateEmployeeAirlineInfo())
	return r
}

// ── TestHTTP_GetEmployeeAirlineInfo ────────────────────────────────────

func TestHTTP_GetEmployeeAirlineInfo(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := newAirlineEmployeeMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errH := middleware.NewErrorHandler(cache)
	enc, _ := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	empID := "550e8400-e29b-41d4-a716-446655440001"
	airlineID := "550e8400-e29b-41d4-a716-446655440099"

	t.Run("success with airline info", func(t *testing.T) {
		aeSvc := &fakeAirlineEmployeeService{
			getByIDFn: func(_ context.Context, id string) (*domain.AirlineEmployee, error) {
				return &domain.AirlineEmployee{ID: empID, AirlineID: airlineID, Bp: "BP1", Active: true}, nil
			},
		}
		airSvc := &fakeAirlineService{
			getByIDFn: func(_ context.Context, id string) (*domain.Airline, error) {
				return &domain.Airline{ID: airlineID, AirlineName: "Avianca", AirlineCode: "AV"}, nil
			},
		}
		router := newAirlineEmployeeTestRouter(aeSvc, airSvc, enc, resp, errH, &domain.Employee{ID: empID})
		req := httptest.NewRequest(http.MethodGet, "/employees/airline", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		t.Logf("status: %d", w.Code)
	})

	t.Run("success without airline info", func(t *testing.T) {
		aeSvc := &fakeAirlineEmployeeService{
			getByIDFn: func(_ context.Context, id string) (*domain.AirlineEmployee, error) {
				return nil, nil // no airline info
			},
		}
		airSvc := &fakeAirlineService{}
		router := newAirlineEmployeeTestRouter(aeSvc, airSvc, enc, resp, errH, &domain.Employee{ID: empID})
		req := httptest.NewRequest(http.MethodGet, "/employees/airline", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		t.Logf("status: %d", w.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		aeSvc := &fakeAirlineEmployeeService{}
		airSvc := &fakeAirlineService{}
		router := newAirlineEmployeeTestRouter(aeSvc, airSvc, enc, resp, errH, nil)
		req := httptest.NewRequest(http.MethodGet, "/employees/airline", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		t.Logf("status: %d", w.Code)
	})
}

// ── TestHTTP_AddEmployeeAirlineInfo ────────────────────────────────────

func TestHTTP_AddEmployeeAirlineInfo(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := newAirlineEmployeeMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errH := middleware.NewErrorHandler(cache)
	enc, _ := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	empID := "550e8400-e29b-41d4-a716-446655440001"
	airlineID := "550e8400-e29b-41d4-a716-446655440099"
	encodedAirlineID, _ := enc.Encode(airlineID)

	validBody := `{
		"airline_id":"` + encodedAirlineID + `",
		"bp":"BP1",
		"start_date":"2025-01-01",
		"end_date":"2025-12-31"
	}`

	t.Run("success", func(t *testing.T) {
		aeSvc := &fakeAirlineEmployeeService{}
		airSvc := &fakeAirlineService{
			getByIDFn: func(_ context.Context, id string) (*domain.Airline, error) {
				return &domain.Airline{ID: airlineID, AirlineName: "Avianca", AirlineCode: "AV"}, nil
			},
		}
		router := newAirlineEmployeeTestRouter(aeSvc, airSvc, enc, resp, errH, &domain.Employee{ID: empID})
		req := httptest.NewRequest(http.MethodPut, "/employees/airline", bytes.NewBufferString(validBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		t.Logf("status: %d", w.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		aeSvc := &fakeAirlineEmployeeService{}
		airSvc := &fakeAirlineService{}
		router := newAirlineEmployeeTestRouter(aeSvc, airSvc, enc, resp, errH, nil)
		req := httptest.NewRequest(http.MethodPut, "/employees/airline", bytes.NewBufferString(validBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		t.Logf("status: %d", w.Code)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		aeSvc := &fakeAirlineEmployeeService{}
		airSvc := &fakeAirlineService{}
		router := newAirlineEmployeeTestRouter(aeSvc, airSvc, enc, resp, errH, &domain.Employee{ID: empID})
		req := httptest.NewRequest(http.MethodPut, "/employees/airline", bytes.NewBufferString("{bad"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		t.Logf("status: %d", w.Code)
	})

	t.Run("invalid airline ID", func(t *testing.T) {
		aeSvc := &fakeAirlineEmployeeService{}
		airSvc := &fakeAirlineService{}
		router := newAirlineEmployeeTestRouter(aeSvc, airSvc, enc, resp, errH, &domain.Employee{ID: empID})
		body := `{"airline_id":"invalid!!!","bp":"BP1","start_date":"2025-01-01"}`
		req := httptest.NewRequest(http.MethodPut, "/employees/airline", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		t.Logf("status: %d", w.Code)
	})

	t.Run("airline not found", func(t *testing.T) {
		aeSvc := &fakeAirlineEmployeeService{}
		airSvc := &fakeAirlineService{
			getByIDFn: func(_ context.Context, id string) (*domain.Airline, error) {
				return nil, domain.ErrAirlineNotFound
			},
		}
		router := newAirlineEmployeeTestRouter(aeSvc, airSvc, enc, resp, errH, &domain.Employee{ID: empID})
		req := httptest.NewRequest(http.MethodPut, "/employees/airline", bytes.NewBufferString(validBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		t.Logf("status: %d", w.Code)
	})

	t.Run("invalid start_date", func(t *testing.T) {
		aeSvc := &fakeAirlineEmployeeService{}
		airSvc := &fakeAirlineService{
			getByIDFn: func(_ context.Context, id string) (*domain.Airline, error) {
				return &domain.Airline{ID: airlineID, AirlineName: "Avianca"}, nil
			},
		}
		router := newAirlineEmployeeTestRouter(aeSvc, airSvc, enc, resp, errH, &domain.Employee{ID: empID})
		body := `{"airline_id":"` + encodedAirlineID + `","bp":"BP1","start_date":"bad-date"}`
		req := httptest.NewRequest(http.MethodPut, "/employees/airline", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		t.Logf("status: %d", w.Code)
	})

	t.Run("invalid end_date", func(t *testing.T) {
		aeSvc := &fakeAirlineEmployeeService{}
		airSvc := &fakeAirlineService{
			getByIDFn: func(_ context.Context, id string) (*domain.Airline, error) {
				return &domain.Airline{ID: airlineID, AirlineName: "Avianca"}, nil
			},
		}
		router := newAirlineEmployeeTestRouter(aeSvc, airSvc, enc, resp, errH, &domain.Employee{ID: empID})
		body := `{"airline_id":"` + encodedAirlineID + `","bp":"BP1","start_date":"2025-01-01","end_date":"bad-date"}`
		req := httptest.NewRequest(http.MethodPut, "/employees/airline", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		t.Logf("status: %d", w.Code)
	})

	t.Run("service save error", func(t *testing.T) {
		aeSvc := &fakeAirlineEmployeeService{
			addFn: func(_ context.Context, _ domain.AirlineEmployee) error { return errors.New("save error") },
		}
		airSvc := &fakeAirlineService{
			getByIDFn: func(_ context.Context, id string) (*domain.Airline, error) {
				return &domain.Airline{ID: airlineID, AirlineName: "Avianca"}, nil
			},
		}
		router := newAirlineEmployeeTestRouter(aeSvc, airSvc, enc, resp, errH, &domain.Employee{ID: empID})
		req := httptest.NewRequest(http.MethodPut, "/employees/airline", bytes.NewBufferString(validBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		t.Logf("status: %d", w.Code)
	})

	t.Run("service foreign key error", func(t *testing.T) {
		aeSvc := &fakeAirlineEmployeeService{
			addFn: func(_ context.Context, _ domain.AirlineEmployee) error { return domain.ErrInvalidForeignKey },
		}
		airSvc := &fakeAirlineService{
			getByIDFn: func(_ context.Context, id string) (*domain.Airline, error) {
				return &domain.Airline{ID: airlineID, AirlineName: "Avianca"}, nil
			},
		}
		router := newAirlineEmployeeTestRouter(aeSvc, airSvc, enc, resp, errH, &domain.Employee{ID: empID})
		req := httptest.NewRequest(http.MethodPut, "/employees/airline", bytes.NewBufferString(validBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		t.Logf("status: %d", w.Code)
	})
}

// ── TestHTTP_UpdateEmployeeAirlineInfo ─────────────────────────────────

func TestHTTP_UpdateEmployeeAirlineInfo(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := newAirlineEmployeeMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errH := middleware.NewErrorHandler(cache)
	enc, _ := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	empID := "550e8400-e29b-41d4-a716-446655440001"
	airlineID := "550e8400-e29b-41d4-a716-446655440099"
	encodedAirlineID, _ := enc.Encode(airlineID)

	validBody := `{
		"airline_id":"` + encodedAirlineID + `",
		"bp":"BP2",
		"start_date":"2025-02-01",
		"end_date":"2026-01-31"
	}`

	t.Run("success", func(t *testing.T) {
		aeSvc := &fakeAirlineEmployeeService{
			getByIDFn: func(_ context.Context, id string) (*domain.AirlineEmployee, error) {
				return &domain.AirlineEmployee{ID: empID, AirlineID: airlineID, Active: true}, nil
			},
		}
		airSvc := &fakeAirlineService{
			getByIDFn: func(_ context.Context, id string) (*domain.Airline, error) {
				return &domain.Airline{ID: airlineID, AirlineName: "Avianca", AirlineCode: "AV"}, nil
			},
		}
		router := newAirlineEmployeeTestRouter(aeSvc, airSvc, enc, resp, errH, &domain.Employee{ID: empID})
		req := httptest.NewRequest(http.MethodPut, "/employees/airline-info", bytes.NewBufferString(validBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		t.Logf("status: %d", w.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		aeSvc := &fakeAirlineEmployeeService{}
		airSvc := &fakeAirlineService{}
		router := newAirlineEmployeeTestRouter(aeSvc, airSvc, enc, resp, errH, nil)
		req := httptest.NewRequest(http.MethodPut, "/employees/airline-info", bytes.NewBufferString(validBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		t.Logf("status: %d", w.Code)
	})

	t.Run("no existing airline info", func(t *testing.T) {
		aeSvc := &fakeAirlineEmployeeService{
			getByIDFn: func(_ context.Context, id string) (*domain.AirlineEmployee, error) {
				return nil, nil
			},
		}
		airSvc := &fakeAirlineService{}
		router := newAirlineEmployeeTestRouter(aeSvc, airSvc, enc, resp, errH, &domain.Employee{ID: empID})
		req := httptest.NewRequest(http.MethodPut, "/employees/airline-info", bytes.NewBufferString(validBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		t.Logf("status: %d", w.Code)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		aeSvc := &fakeAirlineEmployeeService{
			getByIDFn: func(_ context.Context, id string) (*domain.AirlineEmployee, error) {
				return &domain.AirlineEmployee{ID: empID, AirlineID: airlineID, Active: true}, nil
			},
		}
		airSvc := &fakeAirlineService{}
		router := newAirlineEmployeeTestRouter(aeSvc, airSvc, enc, resp, errH, &domain.Employee{ID: empID})
		req := httptest.NewRequest(http.MethodPut, "/employees/airline-info", bytes.NewBufferString("{bad"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		t.Logf("status: %d", w.Code)
	})

	t.Run("invalid airline ID", func(t *testing.T) {
		aeSvc := &fakeAirlineEmployeeService{
			getByIDFn: func(_ context.Context, id string) (*domain.AirlineEmployee, error) {
				return &domain.AirlineEmployee{ID: empID, AirlineID: airlineID, Active: true}, nil
			},
		}
		airSvc := &fakeAirlineService{}
		router := newAirlineEmployeeTestRouter(aeSvc, airSvc, enc, resp, errH, &domain.Employee{ID: empID})
		body := `{"airline_id":"invalid!!!","bp":"BP2","start_date":"2025-02-01"}`
		req := httptest.NewRequest(http.MethodPut, "/employees/airline-info", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		t.Logf("status: %d", w.Code)
	})

	t.Run("airline not found", func(t *testing.T) {
		aeSvc := &fakeAirlineEmployeeService{
			getByIDFn: func(_ context.Context, id string) (*domain.AirlineEmployee, error) {
				return &domain.AirlineEmployee{ID: empID, AirlineID: airlineID, Active: true}, nil
			},
		}
		airSvc := &fakeAirlineService{
			getByIDFn: func(_ context.Context, id string) (*domain.Airline, error) { return nil, domain.ErrAirlineNotFound },
		}
		router := newAirlineEmployeeTestRouter(aeSvc, airSvc, enc, resp, errH, &domain.Employee{ID: empID})
		req := httptest.NewRequest(http.MethodPut, "/employees/airline-info", bytes.NewBufferString(validBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		t.Logf("status: %d", w.Code)
	})

	t.Run("invalid start_date", func(t *testing.T) {
		aeSvc := &fakeAirlineEmployeeService{
			getByIDFn: func(_ context.Context, id string) (*domain.AirlineEmployee, error) {
				return &domain.AirlineEmployee{ID: empID, AirlineID: airlineID, Active: true}, nil
			},
		}
		airSvc := &fakeAirlineService{
			getByIDFn: func(_ context.Context, id string) (*domain.Airline, error) {
				return &domain.Airline{ID: airlineID, AirlineName: "Avianca"}, nil
			},
		}
		router := newAirlineEmployeeTestRouter(aeSvc, airSvc, enc, resp, errH, &domain.Employee{ID: empID})
		body := `{"airline_id":"` + encodedAirlineID + `","bp":"BP2","start_date":"bad-date"}`
		req := httptest.NewRequest(http.MethodPut, "/employees/airline-info", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		t.Logf("status: %d", w.Code)
	})

	t.Run("service update error", func(t *testing.T) {
		aeSvc := &fakeAirlineEmployeeService{
			getByIDFn: func(_ context.Context, id string) (*domain.AirlineEmployee, error) {
				return &domain.AirlineEmployee{ID: empID, AirlineID: airlineID, Active: true}, nil
			},
			updateFn: func(_ context.Context, _ domain.AirlineEmployee) error { return errors.New("db error") },
		}
		airSvc := &fakeAirlineService{
			getByIDFn: func(_ context.Context, id string) (*domain.Airline, error) {
				return &domain.Airline{ID: airlineID, AirlineName: "Avianca"}, nil
			},
		}
		router := newAirlineEmployeeTestRouter(aeSvc, airSvc, enc, resp, errH, &domain.Employee{ID: empID})
		req := httptest.NewRequest(http.MethodPut, "/employees/airline-info", bytes.NewBufferString(validBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		t.Logf("status: %d", w.Code)
	})

	t.Run("service foreign key error", func(t *testing.T) {
		aeSvc := &fakeAirlineEmployeeService{
			getByIDFn: func(_ context.Context, id string) (*domain.AirlineEmployee, error) {
				return &domain.AirlineEmployee{ID: empID, AirlineID: airlineID, Active: true}, nil
			},
			updateFn: func(_ context.Context, _ domain.AirlineEmployee) error { return domain.ErrInvalidForeignKey },
		}
		airSvc := &fakeAirlineService{
			getByIDFn: func(_ context.Context, id string) (*domain.Airline, error) {
				return &domain.Airline{ID: airlineID, AirlineName: "Avianca"}, nil
			},
		}
		router := newAirlineEmployeeTestRouter(aeSvc, airSvc, enc, resp, errH, &domain.Employee{ID: empID})
		req := httptest.NewRequest(http.MethodPut, "/employees/airline-info", bytes.NewBufferString(validBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		t.Logf("status: %d", w.Code)
	})
}

// ── TestHTTP_ActivateEmployeeAirlineInfo ───────────────────────────────

func TestHTTP_ActivateEmployeeAirlineInfo(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := newAirlineEmployeeMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errH := middleware.NewErrorHandler(cache)
	enc, _ := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	empID := "550e8400-e29b-41d4-a716-446655440001"
	airlineID := "550e8400-e29b-41d4-a716-446655440099"

	t.Run("success", func(t *testing.T) {
		aeSvc := &fakeAirlineEmployeeService{
			getByIDFn: func(_ context.Context, id string) (*domain.AirlineEmployee, error) {
				return &domain.AirlineEmployee{ID: empID, AirlineID: airlineID, Active: false}, nil
			},
		}
		airSvc := &fakeAirlineService{}
		router := newAirlineEmployeeTestRouter(aeSvc, airSvc, enc, resp, errH, &domain.Employee{ID: empID})
		req := httptest.NewRequest(http.MethodPatch, "/employees/airline/activate", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		t.Logf("status: %d", w.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		aeSvc := &fakeAirlineEmployeeService{}
		airSvc := &fakeAirlineService{}
		router := newAirlineEmployeeTestRouter(aeSvc, airSvc, enc, resp, errH, nil)
		req := httptest.NewRequest(http.MethodPatch, "/employees/airline/activate", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		t.Logf("status: %d", w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		aeSvc := &fakeAirlineEmployeeService{
			getByIDFn: func(_ context.Context, id string) (*domain.AirlineEmployee, error) {
				return nil, nil // triggers ErrAirlineEmployeeNotFound
			},
		}
		airSvc := &fakeAirlineService{}
		router := newAirlineEmployeeTestRouter(aeSvc, airSvc, enc, resp, errH, &domain.Employee{ID: empID})
		req := httptest.NewRequest(http.MethodPatch, "/employees/airline/activate", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		t.Logf("status: %d", w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		aeSvc := &fakeAirlineEmployeeService{
			getByIDFn: func(_ context.Context, id string) (*domain.AirlineEmployee, error) {
				return &domain.AirlineEmployee{ID: empID, AirlineID: airlineID}, nil
			},
			activateAEFn: func(_ context.Context, id string) error { return errors.New("db error") },
		}
		airSvc := &fakeAirlineService{}
		router := newAirlineEmployeeTestRouter(aeSvc, airSvc, enc, resp, errH, &domain.Employee{ID: empID})
		req := httptest.NewRequest(http.MethodPatch, "/employees/airline/activate", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		t.Logf("status: %d", w.Code)
	})
}

// ── TestHTTP_DeactivateEmployeeAirlineInfo ─────────────────────────────

func TestHTTP_DeactivateEmployeeAirlineInfo(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := newAirlineEmployeeMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errH := middleware.NewErrorHandler(cache)
	enc, _ := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	empID := "550e8400-e29b-41d4-a716-446655440001"
	airlineID := "550e8400-e29b-41d4-a716-446655440099"

	t.Run("success", func(t *testing.T) {
		aeSvc := &fakeAirlineEmployeeService{
			getByIDFn: func(_ context.Context, id string) (*domain.AirlineEmployee, error) {
				return &domain.AirlineEmployee{ID: empID, AirlineID: airlineID, Active: true}, nil
			},
		}
		airSvc := &fakeAirlineService{}
		router := newAirlineEmployeeTestRouter(aeSvc, airSvc, enc, resp, errH, &domain.Employee{ID: empID})
		req := httptest.NewRequest(http.MethodPatch, "/employees/airline/deactivate", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		t.Logf("status: %d", w.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		aeSvc := &fakeAirlineEmployeeService{}
		airSvc := &fakeAirlineService{}
		router := newAirlineEmployeeTestRouter(aeSvc, airSvc, enc, resp, errH, nil)
		req := httptest.NewRequest(http.MethodPatch, "/employees/airline/deactivate", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		t.Logf("status: %d", w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		aeSvc := &fakeAirlineEmployeeService{
			getByIDFn: func(_ context.Context, id string) (*domain.AirlineEmployee, error) {
				return nil, nil
			},
		}
		airSvc := &fakeAirlineService{}
		router := newAirlineEmployeeTestRouter(aeSvc, airSvc, enc, resp, errH, &domain.Employee{ID: empID})
		req := httptest.NewRequest(http.MethodPatch, "/employees/airline/deactivate", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		t.Logf("status: %d", w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		aeSvc := &fakeAirlineEmployeeService{
			getByIDFn: func(_ context.Context, id string) (*domain.AirlineEmployee, error) {
				return &domain.AirlineEmployee{ID: empID, AirlineID: airlineID}, nil
			},
			deactivateAEFn: func(_ context.Context, id string) error { return errors.New("db error") },
		}
		airSvc := &fakeAirlineService{}
		router := newAirlineEmployeeTestRouter(aeSvc, airSvc, enc, resp, errH, &domain.Employee{ID: empID})
		req := httptest.NewRequest(http.MethodPatch, "/employees/airline/deactivate", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		t.Logf("status: %d", w.Code)
	})
}
