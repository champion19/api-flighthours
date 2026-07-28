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
	"github.com/champion19/api-flighthours/middleware"
	"github.com/champion19/api-flighthours/platform/cache/messaging"
	cachetypes "github.com/champion19/api-flighthours/platform/cache/types"
	"github.com/champion19/api-flighthours/tools/idencoder"
	"github.com/gin-gonic/gin"
)

// ═══════════════════════════════════════════
// Fake FlightSummaryService for controller tests
// ═══════════════════════════════════════════

type fakeFlightSummaryServiceCtrl struct {
	getSummaryFn  func(ctx context.Context, employeeID, startDate, endDate string) (*domain.FlightHoursSummary, error)
	getRecentFn   func(ctx context.Context, employeeID string, limit int) ([]domain.DailyLogbookDetail, error)
	getLandingsFn func(ctx context.Context, employeeID, startDate, endDate string) (int, error)
	getDailyFn    func(ctx context.Context, employeeID, date string) (int, error)
	calcPeriodFn  func(period, referenceDate string) (string, string, error)
	buildAlertsFn func(ctx context.Context, employeeID string) ([]domain.FlightAlert, error)
}

func (f *fakeFlightSummaryServiceCtrl) GetFlightHoursSummary(ctx context.Context, employeeID, startDate, endDate string) (*domain.FlightHoursSummary, error) {
	if f.getSummaryFn != nil {
		return f.getSummaryFn(ctx, employeeID, startDate, endDate)
	}
	return &domain.FlightHoursSummary{Breakdown: make(map[string]int)}, nil
}

func (f *fakeFlightSummaryServiceCtrl) GetRecentFlights(ctx context.Context, employeeID string, limit int) ([]domain.DailyLogbookDetail, error) {
	if f.getRecentFn != nil {
		return f.getRecentFn(ctx, employeeID, limit)
	}
	return nil, nil
}

func (f *fakeFlightSummaryServiceCtrl) GetLandingCount(ctx context.Context, employeeID, startDate, endDate string) (int, error) {
	if f.getLandingsFn != nil {
		return f.getLandingsFn(ctx, employeeID, startDate, endDate)
	}
	return 0, nil
}

func (f *fakeFlightSummaryServiceCtrl) GetDailyFlightSeconds(ctx context.Context, employeeID, date string) (int, error) {
	if f.getDailyFn != nil {
		return f.getDailyFn(ctx, employeeID, date)
	}
	return 0, nil
}

func (f *fakeFlightSummaryServiceCtrl) CalculatePeriodDates(period, referenceDate string) (string, string, error) {
	if f.calcPeriodFn != nil {
		return f.calcPeriodFn(period, referenceDate)
	}
	return "2026-01-01", "2026-01-31", nil
}

func (f *fakeFlightSummaryServiceCtrl) BuildFlightAlerts(ctx context.Context, employeeID string) ([]domain.FlightAlert, error) {
	if f.buildAlertsFn != nil {
		return f.buildAlertsFn(ctx, employeeID)
	}
	return nil, nil
}

// ═══════════════════════════════════════════
// Helper: message cache & router
// ═══════════════════════════════════════════

func newFlightSummaryMessageCache(t *testing.T) *messaging.MessageCache {
	t.Helper()
	repo := fakeMessageCacheRepo{messages: []cachetypes.CachedMessage{
		{Code: domain.MsgFlightSummaryGetOK, Type: cachetypes.TypeSuccess, Content: "summary ok"},
		{Code: domain.MsgFlightSummaryGetErr, Type: cachetypes.TypeError, Content: "summary error"},
		{Code: domain.MsgFlightSummaryInvalid, Type: cachetypes.TypeError, Content: "invalid period"},
		{Code: domain.MsgFlightAlertsGetOK, Type: cachetypes.TypeSuccess, Content: "alerts ok"},
		{Code: domain.MsgFlightAlertsGetErr, Type: cachetypes.TypeError, Content: "alerts error"},
		{Code: domain.MsgRecentFlightsGetOK, Type: cachetypes.TypeSuccess, Content: "recent ok"},
		{Code: domain.MsgRecentFlightsGetErr, Type: cachetypes.TypeError, Content: "recent error"},
		{Code: domain.MsgServerError, Type: cachetypes.TypeError, Content: "server error"},
	}}
	cache := messaging.NewMessageCache(repo, 0)
	if err := cache.LoadMessages(context.Background()); err != nil {
		t.Fatalf("failed to load message cache: %v", err)
	}
	return cache
}

func newFlightSummaryTestRouter(
	svc input.FlightSummaryService,
	resp *middleware.ResponseHandler,
	errHandler *middleware.ErrorHandler,
	enc *idencoder.HashidsEncoder,
	authUser *domain.Employee,
) *gin.Engine {
	flightInteractor := interactor.NewFlightSummaryInteractor(svc)
	h := New(HandlerDeps{
		EmployeeInteractor:      &fakeEmployeeInteractor{},
		IDEncoder:               enc,
		Response:                resp,
		FlightSummaryInteractor: flightInteractor,
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
	r.GET("/employees/flight-hours-summary", h.GetFlightHoursSummary())
	r.GET("/employees/flight-alerts", h.GetFlightAlerts())
	r.GET("/employees/recent-flights", h.GetRecentFlights())
	return r
}

// ═══════════════════════════════════════════
// EP1: GetFlightHoursSummary
// ═══════════════════════════════════════════

func TestHTTP_GetFlightHoursSummary(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := newFlightSummaryMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)
	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	if err != nil {
		t.Fatalf("encoder: %v", err)
	}

	t.Run("success monthly", func(t *testing.T) {
		svc := &fakeFlightSummaryServiceCtrl{
			getSummaryFn: func(ctx context.Context, _, _, _ string) (*domain.FlightHoursSummary, error) {
				return &domain.FlightHoursSummary{
					TotalMinutes: 150,
					TotalFlights: 5,
					Breakdown:    map[string]int{"PF": 90, "PM": 60},
				}, nil
			},
		}
		router := newFlightSummaryTestRouter(svc, resp, errHandler, enc, &domain.Employee{ID: "emp-1"})
		req := httptest.NewRequest(http.MethodGet, "/employees/flight-hours-summary?period=monthly", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		t.Logf("status: %d, body: %s", w.Code, w.Body.String())
	})

	t.Run("custom period without dates", func(t *testing.T) {
		svc := &fakeFlightSummaryServiceCtrl{}
		router := newFlightSummaryTestRouter(svc, resp, errHandler, enc, &domain.Employee{ID: "emp-1"})
		req := httptest.NewRequest(http.MethodGet, "/employees/flight-hours-summary?period=custom", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		t.Logf("status: %d", w.Code)
	})

	t.Run("custom period with dates", func(t *testing.T) {
		svc := &fakeFlightSummaryServiceCtrl{
			getSummaryFn: func(ctx context.Context, _, _, _ string) (*domain.FlightHoursSummary, error) {
				return &domain.FlightHoursSummary{Breakdown: map[string]int{}}, nil
			},
		}
		router := newFlightSummaryTestRouter(svc, resp, errHandler, enc, &domain.Employee{ID: "emp-1"})
		req := httptest.NewRequest(http.MethodGet, "/employees/flight-hours-summary?period=custom&start_date=2026-01-01&end_date=2026-06-30", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		t.Logf("status: %d", w.Code)
	})

	t.Run("invalid period", func(t *testing.T) {
		svc := &fakeFlightSummaryServiceCtrl{}
		router := newFlightSummaryTestRouter(svc, resp, errHandler, enc, &domain.Employee{ID: "emp-1"})
		req := httptest.NewRequest(http.MethodGet, "/employees/flight-hours-summary?period=invalid", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		t.Logf("status: %d", w.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		svc := &fakeFlightSummaryServiceCtrl{}
		router := newFlightSummaryTestRouter(svc, resp, errHandler, enc, nil)
		req := httptest.NewRequest(http.MethodGet, "/employees/flight-hours-summary?period=monthly", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		t.Logf("status: %d", w.Code)
	})

	t.Run("interactor error", func(t *testing.T) {
		svc := &fakeFlightSummaryServiceCtrl{
			getSummaryFn: func(ctx context.Context, _, _, _ string) (*domain.FlightHoursSummary, error) {
				return nil, errors.New("db error")
			},
		}
		router := newFlightSummaryTestRouter(svc, resp, errHandler, enc, &domain.Employee{ID: "emp-1"})
		req := httptest.NewRequest(http.MethodGet, "/employees/flight-hours-summary?period=monthly", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		t.Logf("status: %d", w.Code)
	})

	t.Run("default period when empty", func(t *testing.T) {
		svc := &fakeFlightSummaryServiceCtrl{
			getSummaryFn: func(ctx context.Context, _, _, _ string) (*domain.FlightHoursSummary, error) {
				return &domain.FlightHoursSummary{Breakdown: map[string]int{}}, nil
			},
		}
		router := newFlightSummaryTestRouter(svc, resp, errHandler, enc, &domain.Employee{ID: "emp-1"})
		req := httptest.NewRequest(http.MethodGet, "/employees/flight-hours-summary", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		t.Logf("status: %d", w.Code)
	})
}

// ═══════════════════════════════════════════
// EP2: GetFlightAlerts
// ═══════════════════════════════════════════

func TestHTTP_GetFlightAlerts(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := newFlightSummaryMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)
	enc, _ := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})

	t.Run("success", func(t *testing.T) {
		svc := &fakeFlightSummaryServiceCtrl{
			buildAlertsFn: func(ctx context.Context, employeeID string) ([]domain.FlightAlert, error) {
				return []domain.FlightAlert{{Type: "TEST", Severity: "INFO"}}, nil
			},
		}
		router := newFlightSummaryTestRouter(svc, resp, errHandler, enc, &domain.Employee{ID: "emp-1"})
		req := httptest.NewRequest(http.MethodGet, "/employees/flight-alerts", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		t.Logf("status: %d, body: %s", w.Code, w.Body.String())
	})

	t.Run("unauthorized", func(t *testing.T) {
		svc := &fakeFlightSummaryServiceCtrl{}
		router := newFlightSummaryTestRouter(svc, resp, errHandler, enc, nil)
		req := httptest.NewRequest(http.MethodGet, "/employees/flight-alerts", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		t.Logf("status: %d", w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &fakeFlightSummaryServiceCtrl{
			buildAlertsFn: func(ctx context.Context, employeeID string) ([]domain.FlightAlert, error) {
				return nil, errors.New("alert error")
			},
		}
		router := newFlightSummaryTestRouter(svc, resp, errHandler, enc, &domain.Employee{ID: "emp-1"})
		req := httptest.NewRequest(http.MethodGet, "/employees/flight-alerts", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		t.Logf("status: %d", w.Code)
	})
}

// ═══════════════════════════════════════════
// EP3: GetRecentFlights
// ═══════════════════════════════════════════

func TestHTTP_GetRecentFlights(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := newFlightSummaryMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)
	enc, _ := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})

	t.Run("success", func(t *testing.T) {
		svc := &fakeFlightSummaryServiceCtrl{
			getRecentFn: func(ctx context.Context, employeeID string, limit int) ([]domain.DailyLogbookDetail, error) {
				pf := domain.PilotRolePF
				return []domain.DailyLogbookDetail{
					{
						ID:                   "550e8400-e29b-41d4-a716-446655440010",
						DailyLogbookID:       "550e8400-e29b-41d4-a716-446655440002",
						FlightNumber:         "AV123",
						OriginAirportID:      "550e8400-e29b-41d4-a716-446655440003",
						DestinationAirportID: "550e8400-e29b-41d4-a716-446655440005",
						TailNumberID:         "550e8400-e29b-41d4-a716-446655440004",
						PilotRole:            &pf,
					},
				}, nil
			},
		}
		router := newFlightSummaryTestRouter(svc, resp, errHandler, enc, &domain.Employee{ID: "emp-1"})
		req := httptest.NewRequest(http.MethodGet, "/employees/recent-flights", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		t.Logf("status: %d, body: %s", w.Code, w.Body.String())
	})

	t.Run("unauthorized", func(t *testing.T) {
		svc := &fakeFlightSummaryServiceCtrl{}
		router := newFlightSummaryTestRouter(svc, resp, errHandler, enc, nil)
		req := httptest.NewRequest(http.MethodGet, "/employees/recent-flights", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		t.Logf("status: %d", w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &fakeFlightSummaryServiceCtrl{
			getRecentFn: func(ctx context.Context, employeeID string, limit int) ([]domain.DailyLogbookDetail, error) {
				return nil, errors.New("db error")
			},
		}
		router := newFlightSummaryTestRouter(svc, resp, errHandler, enc, &domain.Employee{ID: "emp-1"})
		req := httptest.NewRequest(http.MethodGet, "/employees/recent-flights", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		t.Logf("status: %d", w.Code)
	})

	t.Run("empty result", func(t *testing.T) {
		svc := &fakeFlightSummaryServiceCtrl{
			getRecentFn: func(ctx context.Context, employeeID string, limit int) ([]domain.DailyLogbookDetail, error) {
				return []domain.DailyLogbookDetail{}, nil
			},
		}
		router := newFlightSummaryTestRouter(svc, resp, errHandler, enc, &domain.Employee{ID: "emp-1"})
		req := httptest.NewRequest(http.MethodGet, "/employees/recent-flights", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		t.Logf("status: %d", w.Code)
	})
}
