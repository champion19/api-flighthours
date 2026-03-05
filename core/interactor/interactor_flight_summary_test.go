package interactor

import (
	"context"
	"errors"
	"testing"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
)

// ═══════════════════════════════════════════
// Fake FlightSummaryService
// ═══════════════════════════════════════════

type fakeFlightSummaryService struct {
	getSummaryFn   func(ctx context.Context, employeeID, startDate, endDate string) (*domain.FlightHoursSummary, error)
	getRecentFn    func(ctx context.Context, employeeID string, limit int) ([]domain.DailyLogbookDetail, error)
	getLandingsFn  func(ctx context.Context, employeeID, startDate, endDate string) (int, error)
	getDailySecsFn func(ctx context.Context, employeeID, date string) (int, error)
	calcPeriodFn   func(period, referenceDate string) (string, string, error)
	buildAlertsFn  func(ctx context.Context, employeeID string) ([]domain.FlightAlert, error)
}

func (f *fakeFlightSummaryService) GetFlightHoursSummary(ctx context.Context, employeeID, startDate, endDate string) (*domain.FlightHoursSummary, error) {
	if f.getSummaryFn != nil {
		return f.getSummaryFn(ctx, employeeID, startDate, endDate)
	}
	return &domain.FlightHoursSummary{}, nil
}

func (f *fakeFlightSummaryService) GetRecentFlights(ctx context.Context, employeeID string, limit int) ([]domain.DailyLogbookDetail, error) {
	if f.getRecentFn != nil {
		return f.getRecentFn(ctx, employeeID, limit)
	}
	return nil, nil
}

func (f *fakeFlightSummaryService) GetLandingCount(ctx context.Context, employeeID, startDate, endDate string) (int, error) {
	if f.getLandingsFn != nil {
		return f.getLandingsFn(ctx, employeeID, startDate, endDate)
	}
	return 0, nil
}

func (f *fakeFlightSummaryService) GetDailyFlightSeconds(ctx context.Context, employeeID, date string) (int, error) {
	if f.getDailySecsFn != nil {
		return f.getDailySecsFn(ctx, employeeID, date)
	}
	return 0, nil
}

func (f *fakeFlightSummaryService) CalculatePeriodDates(period, referenceDate string) (string, string, error) {
	if f.calcPeriodFn != nil {
		return f.calcPeriodFn(period, referenceDate)
	}
	return "2026-01-01", "2026-01-31", nil
}

func (f *fakeFlightSummaryService) BuildFlightAlerts(ctx context.Context, employeeID string) ([]domain.FlightAlert, error) {
	if f.buildAlertsFn != nil {
		return f.buildAlertsFn(ctx, employeeID)
	}
	return nil, nil
}

// fakeTx for unused Tx parameter — needed to satisfy the interface if not already defined in this package
// If fakeTx is already defined (from other test files), comment this out
// type fakeTx struct{}
// func (f *fakeTx) Commit() error   { return nil }
// func (f *fakeTx) Rollback() error { return nil }

// ═══════════════════════════════════════════
// Tests: FlightSummaryInteractor
// ═══════════════════════════════════════════

func TestNewFlightSummaryInteractor(t *testing.T) {
	inter := NewFlightSummaryInteractor(&fakeFlightSummaryService{})
	if inter == nil {
		t.Error("expected non-nil interactor")
	}
}

func TestFlightSummaryInteractor_GetFlightHoursSummary(t *testing.T) {
	t.Run("named period calculates dates", func(t *testing.T) {
		svc := &fakeFlightSummaryService{
			calcPeriodFn: func(period, referenceDate string) (string, string, error) {
				return "2026-03-01", "2026-03-31", nil
			},
			getSummaryFn: func(ctx context.Context, employeeID, startDate, endDate string) (*domain.FlightHoursSummary, error) {
				if startDate != "2026-03-01" || endDate != "2026-03-31" {
					t.Errorf("expected computed dates, got %s to %s", startDate, endDate)
				}
				return &domain.FlightHoursSummary{TotalMinutes: 100}, nil
			},
		}
		inter := NewFlightSummaryInteractor(svc)
		result, err := inter.GetFlightHoursSummary(context.Background(), "t-1", "emp-1", "monthly", "", "", "2026-03-15")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Period != "monthly" {
			t.Errorf("expected period 'monthly', got %q", result.Period)
		}
	})

	t.Run("custom period uses provided dates", func(t *testing.T) {
		svc := &fakeFlightSummaryService{
			getSummaryFn: func(ctx context.Context, employeeID, startDate, endDate string) (*domain.FlightHoursSummary, error) {
				if startDate != "2026-01-01" || endDate != "2026-06-30" {
					t.Errorf("expected provided dates, got %s to %s", startDate, endDate)
				}
				return &domain.FlightHoursSummary{}, nil
			},
		}
		inter := NewFlightSummaryInteractor(svc)
		_, err := inter.GetFlightHoursSummary(context.Background(), "t-1", "emp-1", "custom", "2026-01-01", "2026-06-30", "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("period calculation error", func(t *testing.T) {
		svc := &fakeFlightSummaryService{
			calcPeriodFn: func(period, referenceDate string) (string, string, error) {
				return "", "", errors.New("invalid period")
			},
		}
		inter := NewFlightSummaryInteractor(svc)
		_, err := inter.GetFlightHoursSummary(context.Background(), "t-1", "emp-1", "invalid", "", "", "")
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("service error", func(t *testing.T) {
		svc := &fakeFlightSummaryService{
			getSummaryFn: func(ctx context.Context, _, _, _ string) (*domain.FlightHoursSummary, error) {
				return nil, errors.New("service error")
			},
		}
		inter := NewFlightSummaryInteractor(svc)
		_, err := inter.GetFlightHoursSummary(context.Background(), "t-1", "emp-1", "custom", "2026-01-01", "2026-01-31", "")
		if err == nil {
			t.Error("expected error")
		}
	})
}

func TestFlightSummaryInteractor_GetFlightAlerts(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &fakeFlightSummaryService{
			buildAlertsFn: func(ctx context.Context, employeeID string) ([]domain.FlightAlert, error) {
				return []domain.FlightAlert{{Type: "TEST"}}, nil
			},
		}
		inter := NewFlightSummaryInteractor(svc)
		alerts, err := inter.GetFlightAlerts(context.Background(), "t-1", "emp-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(alerts) != 1 {
			t.Errorf("expected 1 alert, got %d", len(alerts))
		}
	})

	t.Run("error", func(t *testing.T) {
		svc := &fakeFlightSummaryService{
			buildAlertsFn: func(ctx context.Context, employeeID string) ([]domain.FlightAlert, error) {
				return nil, errors.New("alert error")
			},
		}
		inter := NewFlightSummaryInteractor(svc)
		_, err := inter.GetFlightAlerts(context.Background(), "t-1", "emp-1")
		if err == nil {
			t.Error("expected error")
		}
	})
}

func TestFlightSummaryInteractor_GetRecentFlights(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &fakeFlightSummaryService{
			getRecentFn: func(ctx context.Context, employeeID string, limit int) ([]domain.DailyLogbookDetail, error) {
				return []domain.DailyLogbookDetail{{ID: "d-1"}, {ID: "d-2"}}, nil
			},
		}
		inter := NewFlightSummaryInteractor(svc)
		flights, err := inter.GetRecentFlights(context.Background(), "t-1", "emp-1", 5)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(flights) != 2 {
			t.Errorf("expected 2 flights, got %d", len(flights))
		}
	})

	t.Run("error", func(t *testing.T) {
		svc := &fakeFlightSummaryService{
			getRecentFn: func(ctx context.Context, employeeID string, limit int) ([]domain.DailyLogbookDetail, error) {
				return nil, errors.New("db error")
			},
		}
		inter := NewFlightSummaryInteractor(svc)
		_, err := inter.GetRecentFlights(context.Background(), "t-1", "emp-1", 5)
		if err == nil {
			t.Error("expected error")
		}
	})
}

// Ensure fakeFlightSummaryService properly satisfies unused methods
var _ = func() {
	var svc fakeFlightSummaryService
	var _ output.Tx
	_, _ = svc.GetLandingCount(context.Background(), "", "", "")
	_, _ = svc.GetDailyFlightSeconds(context.Background(), "", "")
}
