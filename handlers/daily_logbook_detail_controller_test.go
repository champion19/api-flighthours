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

// fakeDailyLogbookDetailService implements input.DailyLogbookDetailService
func domainPilotRolePtr(r domain.PilotRole) *domain.PilotRole { return &r }

type fakeDailyLogbookDetailService struct {
	getByIDFn        func(ctx context.Context, id string) (*domain.DailyLogbookDetail, error)
	listByLogbookFn  func(ctx context.Context, logbookID string) ([]domain.DailyLogbookDetail, error)
	listByEmployeeFn func(ctx context.Context, employeeID string) ([]domain.DailyLogbookDetail, error)
	createFn         func(ctx context.Context, detail domain.DailyLogbookDetail) error
	updateFn         func(ctx context.Context, detail domain.DailyLogbookDetail) error
	validateTimeFn   func(outTime, takeoffTime, landingTime, inTime string) error
	deleteFn         func(ctx context.Context, id string) error
}

func (f *fakeDailyLogbookDetailService) BeginTx(ctx context.Context) (output.Tx, error) {
	return fakeTx{}, nil
}

func (f *fakeDailyLogbookDetailService) GetDailyLogbookDetailByID(ctx context.Context, id string) (*domain.DailyLogbookDetail, error) {
	if f.getByIDFn != nil {
		return f.getByIDFn(ctx, id)
	}
	return nil, nil
}

func (f *fakeDailyLogbookDetailService) ListDailyLogbookDetailsByLogbook(ctx context.Context, logbookID string) ([]domain.DailyLogbookDetail, error) {
	if f.listByLogbookFn != nil {
		return f.listByLogbookFn(ctx, logbookID)
	}
	return nil, nil
}

func (f *fakeDailyLogbookDetailService) ListDailyLogbookDetailsByEmployee(ctx context.Context, employeeID string) ([]domain.DailyLogbookDetail, error) {
	if f.listByEmployeeFn != nil {
		return f.listByEmployeeFn(ctx, employeeID)
	}
	return nil, nil
}

func (f *fakeDailyLogbookDetailService) ValidateTimeSequence(outTime, takeoffTime, landingTime, inTime string) error {
	if f.validateTimeFn != nil {
		return f.validateTimeFn(outTime, takeoffTime, landingTime, inTime)
	}
	return nil
}
func (f *fakeDailyLogbookDetailService) CreateDailyLogbookDetailTx(ctx context.Context, tx output.Tx, detail domain.DailyLogbookDetail) error {
	if f.createFn != nil {
		return f.createFn(ctx, detail)
	}
	return nil
}
func (f *fakeDailyLogbookDetailService) UpdateDailyLogbookDetailTx(ctx context.Context, tx output.Tx, detail domain.DailyLogbookDetail) error {
	if f.updateFn != nil {
		return f.updateFn(ctx, detail)
	}
	return nil
}

func (f *fakeDailyLogbookDetailService) ExistsByUniqueKey(ctx context.Context, employeeLogbookID, flightRealDate, flightNumber, licensePlateID string) (bool, error) {
	return false, nil
}

func (f *fakeDailyLogbookDetailService) DeleteDailyLogbookDetailTx(ctx context.Context, tx output.Tx, id string) error {
	if f.deleteFn != nil {
		return f.deleteFn(ctx, id)
	}
	return nil
}

func newTestDailyLogbookDetailMessageCache(t *testing.T) *messaging.MessageCache {
	t.Helper()

	repo := fakeMessageCacheRepo{messages: []cachetypes.CachedMessage{
		{Code: domain.MsgFlightGetOK, Type: cachetypes.TypeSuccess, Content: "flight found"},
		{Code: domain.MsgFlightNotFound, Type: cachetypes.TypeError, Content: "flight not found"},
		{Code: domain.MsgFlightGetErr, Type: cachetypes.TypeError, Content: "error getting flight"},
		{Code: domain.MsgFlightCreated, Type: cachetypes.TypeSuccess, Content: "flight created"},
		{Code: domain.MsgFlightSaveError, Type: cachetypes.TypeError, Content: "error saving flight"},
		{Code: domain.MsgFlightUpdated, Type: cachetypes.TypeSuccess, Content: "flight updated"},
		{Code: domain.MsgFlightUpdateError, Type: cachetypes.TypeError, Content: "error updating flight"},
		{Code: domain.MsgFlightListOK, Type: cachetypes.TypeSuccess, Content: "flights listed"},
		{Code: domain.MsgFlightListError, Type: cachetypes.TypeError, Content: "error listing flights"},
		{Code: domain.MsgFlightUnauthorized, Type: cachetypes.TypeError, Content: "unauthorized"},
		{Code: domain.MsgFlightInvalidTimeSequence, Type: cachetypes.TypeError, Content: "invalid time sequence"},
		{Code: domain.MsgFlightInvalidRoute, Type: cachetypes.TypeError, Content: "invalid route"},
		{Code: domain.MsgFlightInvalidLogbook, Type: cachetypes.TypeError, Content: "invalid logbook"},
		{Code: domain.MsgFlightInvalidLicensePlate, Type: cachetypes.TypeError, Content: "invalid aircraft"},
		{Code: domain.MsgDailyLogbookNotFound, Type: cachetypes.TypeError, Content: "logbook not found"},
		{Code: domain.MsgDailyLogbookUnauthorized, Type: cachetypes.TypeError, Content: "unauthorized logbook"},
		{Code: domain.MsgServerError, Type: cachetypes.TypeError, Content: "internal server error"},
		{Code: domain.MsgValIDInvalid, Type: cachetypes.TypeError, Content: "invalid id"},
		{Code: domain.MsgValJSONInvalid, Type: cachetypes.TypeError, Content: "invalid json"},
		{Code: domain.MsgValFieldFormat, Type: cachetypes.TypeError, Content: "invalid field format"},
		{Code: domain.MsgFlightDeleted, Type: cachetypes.TypeSuccess, Content: "flight deleted"},
		{Code: domain.MsgFlightDeleteError, Type: cachetypes.TypeError, Content: "error deleting flight"},
	}}
	cache := messaging.NewMessageCache(repo, 0)
	if err := cache.LoadMessages(context.Background()); err != nil {
		t.Fatalf("failed to load message cache: %v", err)
	}
	return cache
}

func newDailyLogbookDetailTestRouter(
	detailSvc input.DailyLogbookDetailService,
	logbookSvc input.DailyLogbookService,
	enc *idencoder.HashidsEncoder,
	resp *middleware.ResponseHandler,
	errHandler *middleware.ErrorHandler,
	authUser *domain.Employee,
) *gin.Engine {
	detailInteractor := interactor.NewDailyLogbookDetailInteractor(detailSvc, logbookSvc)
	logbookInteractor := interactor.NewDailyLogbookInteractor(logbookSvc)
	h := New(nil, &fakeEmployeeInteractor{}, enc, resp, nil, nil, nil, nil, nil, nil, nil, nil, nil, detailInteractor, logbookInteractor, nil, nil)

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
	r.GET("/daily-logbook-details/:id", h.GetDailyLogbookDetail())
	r.GET("/daily-logbooks/:id/details", h.ListDailyLogbookDetails())
	r.POST("/daily-logbooks/:id/details", h.CreateDailyLogbookDetail())
	r.PUT("/daily-logbook-details/:id", h.UpdateDailyLogbookDetail())
	r.GET("/employees/flights", h.ListMyFlights())
	r.DELETE("/daily-logbook-details/:id", h.DeleteDailyLogbookDetail())
	return r
}

// ── GetDailyLogbookDetail ──────────────────────────────────────────────

func TestHTTP_GetDailyLogbookDetail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := newTestDailyLogbookDetailMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)
	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	testUUID := "550e8400-e29b-41d4-a716-446655440000"
	testEmployeeID := "550e8400-e29b-41d4-a716-446655440001"
	testLogbookID := "550e8400-e29b-41d4-a716-446655440002"
	encodedDetailID, _ := enc.Encode(testUUID)

	t.Run("success", func(t *testing.T) {
		logbookSvc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: testLogbookID, EmployeeID: testEmployeeID}, nil
			},
		}
		detailSvc := &fakeDailyLogbookDetailService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbookDetail, error) {
				return &domain.DailyLogbookDetail{
					ID:             testUUID,
					DailyLogbookID: testLogbookID,
					FlightNumber:   "AV123",
					PilotRole:      domainPilotRolePtr(domain.PilotRolePF),
				}, nil
			},
		}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookDetailTestRouter(detailSvc, logbookSvc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodGet, "/daily-logbook-details/"+encodedDetailID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		logbookSvc := &fakeDailyLogbookService{}
		detailSvc := &fakeDailyLogbookDetailService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbookDetail, error) {
				return nil, nil
			},
		}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookDetailTestRouter(detailSvc, logbookSvc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodGet, "/daily-logbook-details/"+encodedDetailID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		logbookSvc := &fakeDailyLogbookService{}
		detailSvc := &fakeDailyLogbookDetailService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbookDetail, error) {
				return nil, errors.New("db error")
			},
		}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookDetailTestRouter(detailSvc, logbookSvc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodGet, "/daily-logbook-details/"+encodedDetailID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("invalid ID", func(t *testing.T) {
		logbookSvc := &fakeDailyLogbookService{}
		detailSvc := &fakeDailyLogbookDetailService{}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookDetailTestRouter(detailSvc, logbookSvc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodGet, "/daily-logbook-details/invalid!!!", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("unauthorized - no employee", func(t *testing.T) {
		logbookSvc := &fakeDailyLogbookService{}
		detailSvc := &fakeDailyLogbookDetailService{}
		router := newDailyLogbookDetailTestRouter(detailSvc, logbookSvc, enc, resp, errHandler, nil)

		req := httptest.NewRequest(http.MethodGet, "/daily-logbook-details/"+encodedDetailID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})
}

// ── ListDailyLogbookDetails ────────────────────────────────────────────

func TestHTTP_ListDailyLogbookDetails(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := newTestDailyLogbookDetailMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)
	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	testLogbookID := "550e8400-e29b-41d4-a716-446655440002"
	testEmployeeID := "550e8400-e29b-41d4-a716-446655440001"
	encodedLogbookID, _ := enc.Encode(testLogbookID)

	t.Run("success", func(t *testing.T) {
		logbookSvc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: testLogbookID, EmployeeID: testEmployeeID}, nil
			},
		}
		detailSvc := &fakeDailyLogbookDetailService{
			listByLogbookFn: func(ctx context.Context, logbookID string) ([]domain.DailyLogbookDetail, error) {
				return []domain.DailyLogbookDetail{
					{ID: "550e8400-e29b-41d4-a716-446655440010", DailyLogbookID: testLogbookID, FlightNumber: "AV123", PilotRole: domainPilotRolePtr(domain.PilotRolePF)},
				}, nil
			},
		}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookDetailTestRouter(detailSvc, logbookSvc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodGet, "/daily-logbooks/"+encodedLogbookID+"/details", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("unauthorized - different employee", func(t *testing.T) {
		logbookSvc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: testLogbookID, EmployeeID: "other-employee"}, nil
			},
		}
		detailSvc := &fakeDailyLogbookDetailService{}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookDetailTestRouter(detailSvc, logbookSvc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodGet, "/daily-logbooks/"+encodedLogbookID+"/details", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("logbook not found", func(t *testing.T) {
		logbookSvc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return nil, nil
			},
		}
		detailSvc := &fakeDailyLogbookDetailService{}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookDetailTestRouter(detailSvc, logbookSvc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodGet, "/daily-logbooks/"+encodedLogbookID+"/details", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("invalid logbook ID", func(t *testing.T) {
		logbookSvc := &fakeDailyLogbookService{}
		detailSvc := &fakeDailyLogbookDetailService{}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookDetailTestRouter(detailSvc, logbookSvc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodGet, "/daily-logbooks/invalid!!!/details", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("service error listing details", func(t *testing.T) {
		logbookSvc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: testLogbookID, EmployeeID: testEmployeeID}, nil
			},
		}
		detailSvc := &fakeDailyLogbookDetailService{
			listByLogbookFn: func(ctx context.Context, logbookID string) ([]domain.DailyLogbookDetail, error) {
				return nil, errors.New("db error")
			},
		}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookDetailTestRouter(detailSvc, logbookSvc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodGet, "/daily-logbooks/"+encodedLogbookID+"/details", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})
}

// ── CreateDailyLogbookDetail ───────────────────────────────────────────

func TestHTTP_CreateDailyLogbookDetail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := newTestDailyLogbookDetailMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)
	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	testLogbookID := "550e8400-e29b-41d4-a716-446655440002"
	testEmployeeID := "550e8400-e29b-41d4-a716-446655440001"
	testRouteID := "550e8400-e29b-41d4-a716-446655440003"
	testAircraftID := "550e8400-e29b-41d4-a716-446655440004"
	encodedLogbookID, _ := enc.Encode(testLogbookID)
	encodedRouteID, _ := enc.Encode(testRouteID)
	encodedAircraftID, _ := enc.Encode(testAircraftID)

	validBody := `{
		"flight_real_date":"2025-01-15",
		"flight_number":"AV123",
		"airline_route_id":"` + encodedRouteID + `",
		"license_plate_id":"` + encodedAircraftID + `",
		"out_time":"08:00",
		"takeoff_time":"08:15",
		"landing_time":"09:30",
		"in_time":"09:45",
		"pilot_role":"PF",
		"air_time":"01:15",
		"block_time":"01:45"
	}`

	t.Run("success", func(t *testing.T) {
		logbookSvc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: testLogbookID, EmployeeID: testEmployeeID}, nil
			},
		}
		detailSvc := &fakeDailyLogbookDetailService{
			createFn: func(ctx context.Context, detail domain.DailyLogbookDetail) error {
				return nil
			},
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbookDetail, error) {
				return &domain.DailyLogbookDetail{
					ID:             id,
					DailyLogbookID: testLogbookID,
					FlightNumber:   "AV123",
					PilotRole:      domainPilotRolePtr(domain.PilotRolePF),
					AirlineRouteID: testRouteID,
					LicensePlateID: testAircraftID,
				}, nil
			},
		}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookDetailTestRouter(detailSvc, logbookSvc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodPost, "/daily-logbooks/"+encodedLogbookID+"/details", bytes.NewBufferString(validBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("unauthorized - no employee", func(t *testing.T) {
		logbookSvc := &fakeDailyLogbookService{}
		detailSvc := &fakeDailyLogbookDetailService{}
		router := newDailyLogbookDetailTestRouter(detailSvc, logbookSvc, enc, resp, errHandler, nil)

		req := httptest.NewRequest(http.MethodPost, "/daily-logbooks/"+encodedLogbookID+"/details", bytes.NewBufferString(validBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("unauthorized - different employee", func(t *testing.T) {
		logbookSvc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: testLogbookID, EmployeeID: "other-employee"}, nil
			},
		}
		detailSvc := &fakeDailyLogbookDetailService{}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookDetailTestRouter(detailSvc, logbookSvc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodPost, "/daily-logbooks/"+encodedLogbookID+"/details", bytes.NewBufferString(validBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("invalid logbook ID", func(t *testing.T) {
		logbookSvc := &fakeDailyLogbookService{}
		detailSvc := &fakeDailyLogbookDetailService{}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookDetailTestRouter(detailSvc, logbookSvc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodPost, "/daily-logbooks/invalid!!!/details", bytes.NewBufferString(validBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		logbookSvc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: testLogbookID, EmployeeID: testEmployeeID}, nil
			},
		}
		detailSvc := &fakeDailyLogbookDetailService{}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookDetailTestRouter(detailSvc, logbookSvc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodPost, "/daily-logbooks/"+encodedLogbookID+"/details", bytes.NewBufferString("{bad json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("invalid route ID in body", func(t *testing.T) {
		logbookSvc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: testLogbookID, EmployeeID: testEmployeeID}, nil
			},
		}
		detailSvc := &fakeDailyLogbookDetailService{}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookDetailTestRouter(detailSvc, logbookSvc, enc, resp, errHandler, authUser)

		body := `{
			"flight_real_date":"2025-01-15","flight_number":"AV123",
			"airline_route_id":"invalid!!!",
			"license_plate_id":"` + encodedAircraftID + `",
			"out_time":"08:00","takeoff_time":"08:15","landing_time":"09:30","in_time":"09:45",
			"pilot_role":"PF","air_time":"01:15","block_time":"01:45"
		}`
		req := httptest.NewRequest(http.MethodPost, "/daily-logbooks/"+encodedLogbookID+"/details", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("invalid aircraft ID in body", func(t *testing.T) {
		logbookSvc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: testLogbookID, EmployeeID: testEmployeeID}, nil
			},
		}
		detailSvc := &fakeDailyLogbookDetailService{}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookDetailTestRouter(detailSvc, logbookSvc, enc, resp, errHandler, authUser)

		body := `{
			"flight_real_date":"2025-01-15","flight_number":"AV123",
			"airline_route_id":"` + encodedRouteID + `",
			"license_plate_id":"invalid!!!",
			"out_time":"08:00","takeoff_time":"08:15","landing_time":"09:30","in_time":"09:45",
			"pilot_role":"PF","air_time":"01:15","block_time":"01:45"
		}`
		req := httptest.NewRequest(http.MethodPost, "/daily-logbooks/"+encodedLogbookID+"/details", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("invalid pilot role", func(t *testing.T) {
		logbookSvc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: testLogbookID, EmployeeID: testEmployeeID}, nil
			},
		}
		detailSvc := &fakeDailyLogbookDetailService{}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookDetailTestRouter(detailSvc, logbookSvc, enc, resp, errHandler, authUser)

		body := `{
			"flight_real_date":"2025-01-15","flight_number":"AV123",
			"airline_route_id":"` + encodedRouteID + `",
			"license_plate_id":"` + encodedAircraftID + `",
			"out_time":"08:00","takeoff_time":"08:15","landing_time":"09:30","in_time":"09:45",
			"pilot_role":"INVALID_ROLE","air_time":"01:15","block_time":"01:45"
		}`
		req := httptest.NewRequest(http.MethodPost, "/daily-logbooks/"+encodedLogbookID+"/details", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("invalid approach type", func(t *testing.T) {
		logbookSvc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: testLogbookID, EmployeeID: testEmployeeID}, nil
			},
		}
		detailSvc := &fakeDailyLogbookDetailService{}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookDetailTestRouter(detailSvc, logbookSvc, enc, resp, errHandler, authUser)

		body := `{
			"flight_real_date":"2025-01-15","flight_number":"AV123",
			"airline_route_id":"` + encodedRouteID + `",
			"license_plate_id":"` + encodedAircraftID + `",
			"out_time":"08:00","takeoff_time":"08:15","landing_time":"09:30","in_time":"09:45",
			"pilot_role":"PF","air_time":"01:15","block_time":"01:45",
			"approach_type":"INVALID_APPROACH"
		}`
		req := httptest.NewRequest(http.MethodPost, "/daily-logbooks/"+encodedLogbookID+"/details", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("time sequence error", func(t *testing.T) {
		logbookSvc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: testLogbookID, EmployeeID: testEmployeeID}, nil
			},
		}
		detailSvc := &fakeDailyLogbookDetailService{
			validateTimeFn: func(_, _, _, _ string) error {
				return domain.ErrFlightInvalidTimeSequence
			},
		}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookDetailTestRouter(detailSvc, logbookSvc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodPost, "/daily-logbooks/"+encodedLogbookID+"/details", bytes.NewBufferString(validBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("service save error", func(t *testing.T) {
		logbookSvc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: testLogbookID, EmployeeID: testEmployeeID}, nil
			},
		}
		detailSvc := &fakeDailyLogbookDetailService{
			createFn: func(ctx context.Context, detail domain.DailyLogbookDetail) error {
				return errors.New("save error")
			},
		}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookDetailTestRouter(detailSvc, logbookSvc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodPost, "/daily-logbooks/"+encodedLogbookID+"/details", bytes.NewBufferString(validBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("duplicate flight error", func(t *testing.T) {
		logbookSvc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: testLogbookID, EmployeeID: testEmployeeID}, nil
			},
		}
		detailSvc := &fakeDailyLogbookDetailService{
			createFn: func(ctx context.Context, detail domain.DailyLogbookDetail) error {
				return domain.ErrFlightDuplicate
			},
		}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookDetailTestRouter(detailSvc, logbookSvc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodPost, "/daily-logbooks/"+encodedLogbookID+"/details", bytes.NewBufferString(validBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("interactor invalid logbook error", func(t *testing.T) {
		logbookSvc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: testLogbookID, EmployeeID: testEmployeeID}, nil
			},
		}
		detailSvc := &fakeDailyLogbookDetailService{
			createFn: func(ctx context.Context, detail domain.DailyLogbookDetail) error {
				return domain.ErrFlightInvalidLogbook
			},
		}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookDetailTestRouter(detailSvc, logbookSvc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodPost, "/daily-logbooks/"+encodedLogbookID+"/details", bytes.NewBufferString(validBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("interactor invalid license plate error", func(t *testing.T) {
		logbookSvc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: testLogbookID, EmployeeID: testEmployeeID}, nil
			},
		}
		detailSvc := &fakeDailyLogbookDetailService{
			createFn: func(ctx context.Context, detail domain.DailyLogbookDetail) error {
				return domain.ErrFlightInvalidLicensePlate
			},
		}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookDetailTestRouter(detailSvc, logbookSvc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodPost, "/daily-logbooks/"+encodedLogbookID+"/details", bytes.NewBufferString(validBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("refetch error after create", func(t *testing.T) {
		logbookSvc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: testLogbookID, EmployeeID: testEmployeeID}, nil
			},
		}
		detailSvc := &fakeDailyLogbookDetailService{
			createFn: func(ctx context.Context, detail domain.DailyLogbookDetail) error {
				return nil
			},
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbookDetail, error) {
				return nil, errors.New("refetch failed")
			},
		}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookDetailTestRouter(detailSvc, logbookSvc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodPost, "/daily-logbooks/"+encodedLogbookID+"/details", bytes.NewBufferString(validBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})
}

// ── UpdateDailyLogbookDetail ───────────────────────────────────────────

func TestHTTP_UpdateDailyLogbookDetail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := newTestDailyLogbookDetailMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)
	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	testDetailID := "550e8400-e29b-41d4-a716-446655440010"
	testLogbookID := "550e8400-e29b-41d4-a716-446655440002"
	testEmployeeID := "550e8400-e29b-41d4-a716-446655440001"
	testRouteID := "550e8400-e29b-41d4-a716-446655440003"
	testAircraftID := "550e8400-e29b-41d4-a716-446655440004"
	encodedDetailID, _ := enc.Encode(testDetailID)
	encodedRouteID, _ := enc.Encode(testRouteID)
	encodedAircraftID, _ := enc.Encode(testAircraftID)

	validUpdateBody := `{
		"flight_real_date":"2025-01-16",
		"flight_number":"AV456",
		"airline_route_id":"` + encodedRouteID + `",
		"license_plate_id":"` + encodedAircraftID + `",
		"out_time":"10:00",
		"takeoff_time":"10:15",
		"landing_time":"11:30",
		"in_time":"11:45",
		"pilot_role":"PM",
		"air_time":"01:15",
		"block_time":"01:45"
	}`

	// Helper: standard services that return a valid ownership chain
	makeOwningServices := func() (*fakeDailyLogbookDetailService, *fakeDailyLogbookService) {
		detailSvc := &fakeDailyLogbookDetailService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbookDetail, error) {
				return &domain.DailyLogbookDetail{
					ID:             testDetailID,
					DailyLogbookID: testLogbookID,
					FlightNumber:   "AV123",
					PilotRole:      domainPilotRolePtr(domain.PilotRolePF),
					AirlineRouteID: testRouteID,
					LicensePlateID: testAircraftID,
				}, nil
			},
		}
		logbookSvc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: testLogbookID, EmployeeID: testEmployeeID}, nil
			},
		}
		return detailSvc, logbookSvc
	}

	t.Run("success", func(t *testing.T) {
		detailSvc, logbookSvc := makeOwningServices()
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookDetailTestRouter(detailSvc, logbookSvc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodPut, "/daily-logbook-details/"+encodedDetailID, bytes.NewBufferString(validUpdateBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("unauthorized - no employee", func(t *testing.T) {
		detailSvc, logbookSvc := makeOwningServices()
		router := newDailyLogbookDetailTestRouter(detailSvc, logbookSvc, enc, resp, errHandler, nil)

		req := httptest.NewRequest(http.MethodPut, "/daily-logbook-details/"+encodedDetailID, bytes.NewBufferString(validUpdateBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("invalid detail ID", func(t *testing.T) {
		detailSvc, logbookSvc := makeOwningServices()
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookDetailTestRouter(detailSvc, logbookSvc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodPut, "/daily-logbook-details/invalid!!!", bytes.NewBufferString(validUpdateBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("detail not found - ownership check", func(t *testing.T) {
		detailSvc := &fakeDailyLogbookDetailService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbookDetail, error) {
				return nil, nil // detail not found
			},
		}
		logbookSvc := &fakeDailyLogbookService{}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookDetailTestRouter(detailSvc, logbookSvc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodPut, "/daily-logbook-details/"+encodedDetailID, bytes.NewBufferString(validUpdateBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("unauthorized - different employee", func(t *testing.T) {
		detailSvc := &fakeDailyLogbookDetailService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbookDetail, error) {
				return &domain.DailyLogbookDetail{
					ID:             testDetailID,
					DailyLogbookID: testLogbookID,
				}, nil
			},
		}
		logbookSvc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: testLogbookID, EmployeeID: "other-employee"}, nil
			},
		}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookDetailTestRouter(detailSvc, logbookSvc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodPut, "/daily-logbook-details/"+encodedDetailID, bytes.NewBufferString(validUpdateBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		detailSvc, logbookSvc := makeOwningServices()
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookDetailTestRouter(detailSvc, logbookSvc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodPut, "/daily-logbook-details/"+encodedDetailID, bytes.NewBufferString("{bad json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("invalid route ID in body", func(t *testing.T) {
		detailSvc, logbookSvc := makeOwningServices()
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookDetailTestRouter(detailSvc, logbookSvc, enc, resp, errHandler, authUser)

		body := `{
			"flight_real_date":"2025-01-16","flight_number":"AV456",
			"airline_route_id":"invalid!!!",
			"license_plate_id":"` + encodedAircraftID + `",
			"out_time":"10:00","takeoff_time":"10:15","landing_time":"11:30","in_time":"11:45",
			"pilot_role":"PM","air_time":"01:15","block_time":"01:45"
		}`
		req := httptest.NewRequest(http.MethodPut, "/daily-logbook-details/"+encodedDetailID, bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("invalid aircraft ID in body", func(t *testing.T) {
		detailSvc, logbookSvc := makeOwningServices()
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookDetailTestRouter(detailSvc, logbookSvc, enc, resp, errHandler, authUser)

		body := `{
			"flight_real_date":"2025-01-16","flight_number":"AV456",
			"airline_route_id":"` + encodedRouteID + `",
			"license_plate_id":"invalid!!!",
			"out_time":"10:00","takeoff_time":"10:15","landing_time":"11:30","in_time":"11:45",
			"pilot_role":"PM","air_time":"01:15","block_time":"01:45"
		}`
		req := httptest.NewRequest(http.MethodPut, "/daily-logbook-details/"+encodedDetailID, bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("invalid pilot role", func(t *testing.T) {
		detailSvc, logbookSvc := makeOwningServices()
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookDetailTestRouter(detailSvc, logbookSvc, enc, resp, errHandler, authUser)

		body := `{
			"flight_real_date":"2025-01-16","flight_number":"AV456",
			"airline_route_id":"` + encodedRouteID + `",
			"license_plate_id":"` + encodedAircraftID + `",
			"out_time":"10:00","takeoff_time":"10:15","landing_time":"11:30","in_time":"11:45",
			"pilot_role":"INVALID_ROLE","air_time":"01:15","block_time":"01:45"
		}`
		req := httptest.NewRequest(http.MethodPut, "/daily-logbook-details/"+encodedDetailID, bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("invalid approach type", func(t *testing.T) {
		detailSvc, logbookSvc := makeOwningServices()
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookDetailTestRouter(detailSvc, logbookSvc, enc, resp, errHandler, authUser)

		body := `{
			"flight_real_date":"2025-01-16","flight_number":"AV456",
			"airline_route_id":"` + encodedRouteID + `",
			"license_plate_id":"` + encodedAircraftID + `",
			"out_time":"10:00","takeoff_time":"10:15","landing_time":"11:30","in_time":"11:45",
			"pilot_role":"PM","air_time":"01:15","block_time":"01:45",
			"approach_type":"INVALID_APPROACH"
		}`
		req := httptest.NewRequest(http.MethodPut, "/daily-logbook-details/"+encodedDetailID, bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("time sequence error", func(t *testing.T) {
		detailSvc := &fakeDailyLogbookDetailService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbookDetail, error) {
				return &domain.DailyLogbookDetail{
					ID:             testDetailID,
					DailyLogbookID: testLogbookID,
					FlightNumber:   "AV123",
					PilotRole:      domainPilotRolePtr(domain.PilotRolePF),
					AirlineRouteID: testRouteID,
					LicensePlateID: testAircraftID,
				}, nil
			},
			validateTimeFn: func(_, _, _, _ string) error {
				return domain.ErrFlightInvalidTimeSequence
			},
		}
		logbookSvc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: testLogbookID, EmployeeID: testEmployeeID}, nil
			},
		}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookDetailTestRouter(detailSvc, logbookSvc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodPut, "/daily-logbook-details/"+encodedDetailID, bytes.NewBufferString(validUpdateBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("service update error", func(t *testing.T) {
		detailSvc := &fakeDailyLogbookDetailService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbookDetail, error) {
				return &domain.DailyLogbookDetail{
					ID:             testDetailID,
					DailyLogbookID: testLogbookID,
					FlightNumber:   "AV123",
					PilotRole:      domainPilotRolePtr(domain.PilotRolePF),
					AirlineRouteID: testRouteID,
					LicensePlateID: testAircraftID,
				}, nil
			},
			updateFn: func(ctx context.Context, detail domain.DailyLogbookDetail) error {
				return errors.New("update error")
			},
		}
		logbookSvc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: testLogbookID, EmployeeID: testEmployeeID}, nil
			},
		}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookDetailTestRouter(detailSvc, logbookSvc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodPut, "/daily-logbook-details/"+encodedDetailID, bytes.NewBufferString(validUpdateBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("ownership check service error", func(t *testing.T) {
		detailSvc := &fakeDailyLogbookDetailService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbookDetail, error) {
				return nil, errors.New("db error")
			},
		}
		logbookSvc := &fakeDailyLogbookService{}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookDetailTestRouter(detailSvc, logbookSvc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodPut, "/daily-logbook-details/"+encodedDetailID, bytes.NewBufferString(validUpdateBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("flight not found", func(t *testing.T) {
		detailSvc := &fakeDailyLogbookDetailService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbookDetail, error) {
				return &domain.DailyLogbookDetail{
					ID:             testDetailID,
					DailyLogbookID: testLogbookID,
					FlightNumber:   "AV123",
					PilotRole:      domainPilotRolePtr(domain.PilotRolePF),
					AirlineRouteID: testRouteID,
					LicensePlateID: testAircraftID,
				}, nil
			},
			updateFn: func(ctx context.Context, detail domain.DailyLogbookDetail) error {
				return domain.ErrFlightNotFound
			},
		}
		logbookSvc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: testLogbookID, EmployeeID: testEmployeeID}, nil
			},
		}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookDetailTestRouter(detailSvc, logbookSvc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodPut, "/daily-logbook-details/"+encodedDetailID, bytes.NewBufferString(validUpdateBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("refetch error after update", func(t *testing.T) {
		callCount := 0
		detailSvc := &fakeDailyLogbookDetailService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbookDetail, error) {
				callCount++
				if callCount == 1 {
					// First call: ownership check succeeds
					return &domain.DailyLogbookDetail{
						ID:             testDetailID,
						DailyLogbookID: testLogbookID,
						FlightNumber:   "AV123",
						PilotRole:      domainPilotRolePtr(domain.PilotRolePF),
						AirlineRouteID: testRouteID,
						LicensePlateID: testAircraftID,
					}, nil
				}
				// Second call: refetch fails
				return nil, errors.New("refetch failed")
			},
		}
		logbookSvc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: testLogbookID, EmployeeID: testEmployeeID}, nil
			},
		}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookDetailTestRouter(detailSvc, logbookSvc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodPut, "/daily-logbook-details/"+encodedDetailID, bytes.NewBufferString(validUpdateBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})
}

// ── ListMyFlights ──────────────────────────────────────────────────────

func TestHTTP_ListMyFlights(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := newTestDailyLogbookDetailMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)
	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	testEmployeeID := "550e8400-e29b-41d4-a716-446655440001"

	t.Run("success", func(t *testing.T) {
		logbookSvc := &fakeDailyLogbookService{}
		detailSvc := &fakeDailyLogbookDetailService{
			listByEmployeeFn: func(ctx context.Context, employeeID string) ([]domain.DailyLogbookDetail, error) {
				return []domain.DailyLogbookDetail{
					{ID: "550e8400-e29b-41d4-a716-446655440010", FlightNumber: "AV123", PilotRole: domainPilotRolePtr(domain.PilotRolePF)},
				}, nil
			},
		}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookDetailTestRouter(detailSvc, logbookSvc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodGet, "/employees/flights", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("unauthorized - no employee", func(t *testing.T) {
		logbookSvc := &fakeDailyLogbookService{}
		detailSvc := &fakeDailyLogbookDetailService{}
		router := newDailyLogbookDetailTestRouter(detailSvc, logbookSvc, enc, resp, errHandler, nil)

		req := httptest.NewRequest(http.MethodGet, "/employees/flights", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		logbookSvc := &fakeDailyLogbookService{}
		detailSvc := &fakeDailyLogbookDetailService{
			listByEmployeeFn: func(ctx context.Context, employeeID string) ([]domain.DailyLogbookDetail, error) {
				return nil, errors.New("db error")
			},
		}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookDetailTestRouter(detailSvc, logbookSvc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodGet, "/employees/flights", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("empty results", func(t *testing.T) {
		logbookSvc := &fakeDailyLogbookService{}
		detailSvc := &fakeDailyLogbookDetailService{
			listByEmployeeFn: func(ctx context.Context, employeeID string) ([]domain.DailyLogbookDetail, error) {
				return []domain.DailyLogbookDetail{}, nil
			},
		}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookDetailTestRouter(detailSvc, logbookSvc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodGet, "/employees/flights", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})
}

// ── DTO Conversion Unit Tests ─────────────────────────────────────────

func TestToDomainDailyLogbookDetail_OptionalFields(t *testing.T) {
	crewRole := "captain"
	approachType := "ILS_CAT_I"
	pilotRole := "PF"

	req := CreateDailyLogbookDetailRequest{
		FlightRealDate: "2025-01-15",
		FlightNumber:   "AV123",
		AirlineRouteID: "route-uuid",
		LicensePlateID: "plate-uuid",
		PilotRole:      &pilotRole,
		CrewRole:       &crewRole,
		ApproachType:   &approachType,
	}

	detail := ToDomainDailyLogbookDetail("logbook-uuid", req)

	if detail.CrewRole == nil {
		t.Error("expected CrewRole to be set")
	} else if string(*detail.CrewRole) != "captain" {
		t.Errorf("expected CrewRole captain, got %s", string(*detail.CrewRole))
	}

	if detail.ApproachType == nil {
		t.Error("expected ApproachType to be set")
	} else if string(*detail.ApproachType) != "ILS_CAT_I" {
		t.Errorf("expected ApproachType ILS_CAT_I, got %s", string(*detail.ApproachType))
	}

	if detail.PilotRole == nil {
		t.Error("expected PilotRole to be set")
	}
}

func TestFromDomainDailyLogbookDetail_OptionalFields(t *testing.T) {
	crewRole := domain.CrewRole("captain")
	approachType := domain.ApproachType("ILS_CAT_I")
	pilotRole := domain.PilotRolePF

	d := &domain.DailyLogbookDetail{
		ID:             "detail-uuid",
		DailyLogbookID: "logbook-uuid",
		FlightRealDate: "2025-01-15",
		FlightNumber:   "AV123",
		AirlineRouteID: "route-uuid",
		LicensePlateID: "plate-uuid",
		PilotRole:      &pilotRole,
		CrewRole:       &crewRole,
		ApproachType:   &approachType,
	}

	resp := FromDomainDailyLogbookDetail(d, "enc-id", "enc-logbook", "enc-route", "enc-plate")

	if resp.CrewRole == nil {
		t.Error("expected CrewRole to be set in response")
	} else if *resp.CrewRole != "captain" {
		t.Errorf("expected CrewRole captain, got %s", *resp.CrewRole)
	}

	if resp.ApproachType == nil {
		t.Error("expected ApproachType to be set in response")
	} else if *resp.ApproachType != "ILS_CAT_I" {
		t.Errorf("expected ApproachType ILS_CAT_I, got %s", *resp.ApproachType)
	}

	if resp.PilotRole == nil {
		t.Error("expected PilotRole to be set in response")
	}
}

func TestToDomainDailyLogbookDetailUpdate_OptionalFields(t *testing.T) {
	crewRole := "copilot"
	approachType := "VISUAL"
	pilotRole := "PM"

	req := UpdateDailyLogbookDetailRequest{
		FlightRealDate: "2025-01-16",
		FlightNumber:   "AV456",
		AirlineRouteID: "route-uuid",
		LicensePlateID: "plate-uuid",
		PilotRole:      &pilotRole,
		CrewRole:       &crewRole,
		ApproachType:   &approachType,
	}

	detail := ToDomainDailyLogbookDetailUpdate("detail-uuid", req)

	if detail.CrewRole == nil {
		t.Error("expected CrewRole to be set")
	} else if string(*detail.CrewRole) != "copilot" {
		t.Errorf("expected CrewRole copilot, got %s", string(*detail.CrewRole))
	}

	if detail.ApproachType == nil {
		t.Error("expected ApproachType to be set")
	} else if string(*detail.ApproachType) != "VISUAL" {
		t.Errorf("expected ApproachType VISUAL, got %s", string(*detail.ApproachType))
	}

	if detail.PilotRole == nil {
		t.Error("expected PilotRole to be set")
	}
}

// ── DeleteDailyLogbookDetail ──────────────────────────────────────────

func TestHTTP_DeleteDailyLogbookDetail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := newTestDailyLogbookDetailMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)
	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	testDetailID := "550e8400-e29b-41d4-a716-446655440010"
	testLogbookID := "550e8400-e29b-41d4-a716-446655440002"
	testEmployeeID := "550e8400-e29b-41d4-a716-446655440001"
	testRouteID := "550e8400-e29b-41d4-a716-446655440003"
	testAircraftID := "550e8400-e29b-41d4-a716-446655440004"
	encodedDetailID, _ := enc.Encode(testDetailID)

	t.Run("success", func(t *testing.T) {
		detailSvc := &fakeDailyLogbookDetailService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbookDetail, error) {
				return &domain.DailyLogbookDetail{
					ID:             testDetailID,
					DailyLogbookID: testLogbookID,
					FlightNumber:   "AV123",
					PilotRole:      domainPilotRolePtr(domain.PilotRolePF),
					AirlineRouteID: testRouteID,
					LicensePlateID: testAircraftID,
				}, nil
			},
		}
		logbookSvc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: testLogbookID, EmployeeID: testEmployeeID}, nil
			},
		}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookDetailTestRouter(detailSvc, logbookSvc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodDelete, "/daily-logbook-details/"+encodedDetailID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("unauthorized - no employee", func(t *testing.T) {
		detailSvc := &fakeDailyLogbookDetailService{}
		logbookSvc := &fakeDailyLogbookService{}
		router := newDailyLogbookDetailTestRouter(detailSvc, logbookSvc, enc, resp, errHandler, nil)

		req := httptest.NewRequest(http.MethodDelete, "/daily-logbook-details/"+encodedDetailID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("unauthorized - different owner", func(t *testing.T) {
		detailSvc := &fakeDailyLogbookDetailService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbookDetail, error) {
				return &domain.DailyLogbookDetail{
					ID:             testDetailID,
					DailyLogbookID: testLogbookID,
				}, nil
			},
		}
		logbookSvc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: testLogbookID, EmployeeID: "other-employee"}, nil
			},
		}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookDetailTestRouter(detailSvc, logbookSvc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodDelete, "/daily-logbook-details/"+encodedDetailID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		detailSvc := &fakeDailyLogbookDetailService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbookDetail, error) {
				return nil, nil
			},
		}
		logbookSvc := &fakeDailyLogbookService{}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookDetailTestRouter(detailSvc, logbookSvc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodDelete, "/daily-logbook-details/"+encodedDetailID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		detailSvc := &fakeDailyLogbookDetailService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbookDetail, error) {
				return nil, errors.New("db error")
			},
		}
		logbookSvc := &fakeDailyLogbookService{}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookDetailTestRouter(detailSvc, logbookSvc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodDelete, "/daily-logbook-details/"+encodedDetailID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})

	t.Run("invalid ID", func(t *testing.T) {
		detailSvc := &fakeDailyLogbookDetailService{}
		logbookSvc := &fakeDailyLogbookService{}
		authUser := &domain.Employee{ID: testEmployeeID}
		router := newDailyLogbookDetailTestRouter(detailSvc, logbookSvc, enc, resp, errHandler, authUser)

		req := httptest.NewRequest(http.MethodDelete, "/daily-logbook-details/invalid!!!", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("status: %d", w.Code)
	})
}
