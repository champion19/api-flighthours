package services

import (
	"context"
	"errors"
	"testing"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

// ═══════════════════════════════════════════
// Mock for FlightSummaryRepository
// ═══════════════════════════════════════════

type mockFlightSummaryRepo struct {
	getSummaryFn   func(ctx context.Context, employeeID, startDate, endDate string) ([]domain.PilotRoleBreakdown, error)
	getRecentFn    func(ctx context.Context, employeeID string, limit int) ([]domain.DailyLogbookDetail, error)
	getLandingsFn  func(ctx context.Context, employeeID, startDate, endDate string) (int, error)
	getDailySecsFn func(ctx context.Context, employeeID, date string) (int, error)
}

func (m *mockFlightSummaryRepo) GetFlightHoursSummary(ctx context.Context, employeeID, startDate, endDate string) ([]domain.PilotRoleBreakdown, error) {
	if m.getSummaryFn != nil {
		return m.getSummaryFn(ctx, employeeID, startDate, endDate)
	}
	return nil, nil
}

func (m *mockFlightSummaryRepo) GetRecentFlights(ctx context.Context, employeeID string, limit int) ([]domain.DailyLogbookDetail, error) {
	if m.getRecentFn != nil {
		return m.getRecentFn(ctx, employeeID, limit)
	}
	return nil, nil
}

func (m *mockFlightSummaryRepo) GetLandingCount(ctx context.Context, employeeID, startDate, endDate string) (int, error) {
	if m.getLandingsFn != nil {
		return m.getLandingsFn(ctx, employeeID, startDate, endDate)
	}
	return 0, nil
}

func (m *mockFlightSummaryRepo) GetDailyFlightSeconds(ctx context.Context, employeeID, date string) (int, error) {
	if m.getDailySecsFn != nil {
		return m.getDailySecsFn(ctx, employeeID, date)
	}
	return 0, nil
}

// ═══════════════════════════════════════════
// Tests: NewFlightSummaryService
// ═══════════════════════════════════════════

func TestNewFlightSummaryService(t *testing.T) {
	svc := NewFlightSummaryService(&mockFlightSummaryRepo{})
	if svc == nil {
		t.Error("expected non-nil service")
	}
}

// ═══════════════════════════════════════════
// Tests: GetFlightHoursSummary
// ═══════════════════════════════════════════

func TestFlightSummaryService_GetFlightHoursSummary(t *testing.T) {
	t.Run("success with breakdown", func(t *testing.T) {
		repo := &mockFlightSummaryRepo{
			getSummaryFn: func(ctx context.Context, employeeID, startDate, endDate string) ([]domain.PilotRoleBreakdown, error) {
				return []domain.PilotRoleBreakdown{
					{PilotRole: "PF", TotalSeconds: 7200, FlightCount: 3},
					{PilotRole: "PM", TotalSeconds: 3600, FlightCount: 2},
				}, nil
			},
			getLandingsFn: func(ctx context.Context, employeeID, startDate, endDate string) (int, error) {
				return 5, nil
			},
		}
		svc := NewFlightSummaryService(repo)
		result, err := svc.GetFlightHoursSummary(context.Background(), "emp-1", "2026-01-01", "2026-01-31")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.TotalMinutes != 180 { // 7200/60 + 3600/60 = 120 + 60
			t.Errorf("expected 180 total minutes, got %d", result.TotalMinutes)
		}
		if result.TotalFlights != 5 {
			t.Errorf("expected 5 total flights, got %d", result.TotalFlights)
		}
		if result.TotalLandings != 5 {
			t.Errorf("expected 5 landings, got %d", result.TotalLandings)
		}
		if result.Breakdown["PF"] != 120 {
			t.Errorf("expected PF=120 min, got %d", result.Breakdown["PF"])
		}
	})

	t.Run("repo error on summary", func(t *testing.T) {
		repo := &mockFlightSummaryRepo{
			getSummaryFn: func(ctx context.Context, employeeID, startDate, endDate string) ([]domain.PilotRoleBreakdown, error) {
				return nil, errors.New("db error")
			},
		}
		svc := NewFlightSummaryService(repo)
		_, err := svc.GetFlightHoursSummary(context.Background(), "emp-1", "2026-01-01", "2026-01-31")
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("repo error on landing count", func(t *testing.T) {
		repo := &mockFlightSummaryRepo{
			getSummaryFn: func(ctx context.Context, employeeID, startDate, endDate string) ([]domain.PilotRoleBreakdown, error) {
				return []domain.PilotRoleBreakdown{}, nil
			},
			getLandingsFn: func(ctx context.Context, employeeID, startDate, endDate string) (int, error) {
				return 0, errors.New("landing error")
			},
		}
		svc := NewFlightSummaryService(repo)
		_, err := svc.GetFlightHoursSummary(context.Background(), "emp-1", "2026-01-01", "2026-01-31")
		if err == nil {
			t.Error("expected error")
		}
	})
}

// ═══════════════════════════════════════════
// Tests: GetRecentFlights, GetLandingCount, GetDailyFlightSeconds
// ═══════════════════════════════════════════

func TestFlightSummaryService_GetRecentFlights(t *testing.T) {
	repo := &mockFlightSummaryRepo{
		getRecentFn: func(ctx context.Context, employeeID string, limit int) ([]domain.DailyLogbookDetail, error) {
			return []domain.DailyLogbookDetail{{ID: "d-1"}}, nil
		},
	}
	svc := NewFlightSummaryService(repo)
	result, err := svc.GetRecentFlights(context.Background(), "emp-1", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 flight, got %d", len(result))
	}
}

func TestFlightSummaryService_GetLandingCount(t *testing.T) {
	repo := &mockFlightSummaryRepo{
		getLandingsFn: func(ctx context.Context, employeeID, startDate, endDate string) (int, error) {
			return 7, nil
		},
	}
	svc := NewFlightSummaryService(repo)
	count, err := svc.GetLandingCount(context.Background(), "emp-1", "2026-01-01", "2026-03-31")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 7 {
		t.Errorf("expected 7, got %d", count)
	}
}

func TestFlightSummaryService_GetDailyFlightSeconds(t *testing.T) {
	repo := &mockFlightSummaryRepo{
		getDailySecsFn: func(ctx context.Context, employeeID, date string) (int, error) {
			return 3600, nil
		},
	}
	svc := NewFlightSummaryService(repo)
	secs, err := svc.GetDailyFlightSeconds(context.Background(), "emp-1", "2026-01-15")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secs != 3600 {
		t.Errorf("expected 3600, got %d", secs)
	}
}

// ═══════════════════════════════════════════
// Tests: CalculatePeriodDates
// ═══════════════════════════════════════════

func TestFlightSummaryService_CalculatePeriodDates(t *testing.T) {
	svc := NewFlightSummaryService(&mockFlightSummaryRepo{})

	t.Run("monthly", func(t *testing.T) {
		start, end, err := svc.CalculatePeriodDates(domain.PeriodMonthly, "2026-03-10")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if start != "2026-03-01" {
			t.Errorf("expected start 2026-03-01, got %s", start)
		}
		if end != "2026-03-31" {
			t.Errorf("expected end 2026-03-31, got %s", end)
		}
	})

	t.Run("bimonthly jan", func(t *testing.T) {
		start, end, err := svc.CalculatePeriodDates(domain.PeriodBimonthly, "2026-01-15")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if start != "2026-01-01" {
			t.Errorf("expected start 2026-01-01, got %s", start)
		}
		if end != "2026-02-28" {
			t.Errorf("expected end 2026-02-28, got %s", end)
		}
	})

	t.Run("quarterly Q2", func(t *testing.T) {
		start, end, err := svc.CalculatePeriodDates(domain.PeriodQuarterly, "2026-05-20")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if start != "2026-04-01" {
			t.Errorf("expected start 2026-04-01, got %s", start)
		}
		if end != "2026-06-30" {
			t.Errorf("expected end 2026-06-30, got %s", end)
		}
	})

	t.Run("semiannual first half", func(t *testing.T) {
		start, end, err := svc.CalculatePeriodDates(domain.PeriodSemiannual, "2026-03-01")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if start != "2026-01-01" || end != "2026-06-30" {
			t.Errorf("expected 2026-01-01 to 2026-06-30, got %s to %s", start, end)
		}
	})

	t.Run("semiannual second half", func(t *testing.T) {
		start, end, err := svc.CalculatePeriodDates(domain.PeriodSemiannual, "2026-09-01")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if start != "2026-07-01" || end != "2026-12-31" {
			t.Errorf("expected 2026-07-01 to 2026-12-31, got %s to %s", start, end)
		}
	})

	t.Run("annual", func(t *testing.T) {
		start, end, err := svc.CalculatePeriodDates(domain.PeriodAnnual, "2026-06-15")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if start != "2026-01-01" || end != "2026-12-31" {
			t.Errorf("expected full year, got %s to %s", start, end)
		}
	})

	t.Run("custom returns error", func(t *testing.T) {
		_, _, err := svc.CalculatePeriodDates(domain.PeriodCustom, "2026-01-01")
		if err == nil {
			t.Error("expected error for custom period")
		}
	})

	t.Run("invalid period", func(t *testing.T) {
		_, _, err := svc.CalculatePeriodDates("weekly", "2026-01-01")
		if err == nil {
			t.Error("expected error for invalid period")
		}
	})

	t.Run("invalid date format", func(t *testing.T) {
		_, _, err := svc.CalculatePeriodDates(domain.PeriodMonthly, "not-a-date")
		if err == nil {
			t.Error("expected error for invalid date")
		}
	})

	t.Run("empty reference date uses now", func(t *testing.T) {
		start, end, err := svc.CalculatePeriodDates(domain.PeriodAnnual, "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if start == "" || end == "" {
			t.Error("expected non-empty dates")
		}
	})
}

// ═══════════════════════════════════════════
// Tests: BuildFlightAlerts
// ═══════════════════════════════════════════

func TestFlightSummaryService_BuildFlightAlerts(t *testing.T) {
	t.Run("no alerts when under thresholds", func(t *testing.T) {
		repo := &mockFlightSummaryRepo{
			getSummaryFn: func(ctx context.Context, _, _, _ string) ([]domain.PilotRoleBreakdown, error) {
				return []domain.PilotRoleBreakdown{
					{PilotRole: "PF", TotalSeconds: 60 * 60, FlightCount: 1}, // 60 min = 1h, well under 50h limit
				}, nil
			},
			getLandingsFn: func(ctx context.Context, _, _, _ string) (int, error) {
				return 5, nil // Above 3 minimum
			},
		}
		svc := NewFlightSummaryService(repo)
		alerts, err := svc.BuildFlightAlerts(context.Background(), "emp-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(alerts) != 0 {
			t.Errorf("expected 0 alerts, got %d", len(alerts))
		}
	})

	t.Run("warning when at limit", func(t *testing.T) {
		repo := &mockFlightSummaryRepo{
			getSummaryFn: func(ctx context.Context, _, _, _ string) ([]domain.PilotRoleBreakdown, error) {
				return []domain.PilotRoleBreakdown{
					{PilotRole: "PF", TotalSeconds: 3000 * 60, FlightCount: 10}, // 3000 min = 50h = at the 15-day limit
				}, nil
			},
			getLandingsFn: func(ctx context.Context, _, _, _ string) (int, error) {
				return 5, nil
			},
		}
		svc := NewFlightSummaryService(repo)
		alerts, err := svc.BuildFlightAlerts(context.Background(), "emp-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// Should have WARNING alerts for exceeding hour limits
		hasWarning := false
		for _, a := range alerts {
			if a.Severity == domain.AlertSeverityWarning {
				hasWarning = true
			}
		}
		if !hasWarning {
			t.Error("expected at least one WARNING alert")
		}
	})

	t.Run("info when at 80% threshold", func(t *testing.T) {
		repo := &mockFlightSummaryRepo{
			getSummaryFn: func(ctx context.Context, _, _, _ string) ([]domain.PilotRoleBreakdown, error) {
				// 80% of 3000 = 2400 min
				return []domain.PilotRoleBreakdown{
					{PilotRole: "PF", TotalSeconds: 2400 * 60, FlightCount: 5},
				}, nil
			},
			getLandingsFn: func(ctx context.Context, _, _, _ string) (int, error) {
				return 5, nil
			},
		}
		svc := NewFlightSummaryService(repo)
		alerts, err := svc.BuildFlightAlerts(context.Background(), "emp-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		hasInfo := false
		for _, a := range alerts {
			if a.Severity == domain.AlertSeverityInfo && a.Type == domain.AlertTypeHourLimit15Days {
				hasInfo = true
			}
		}
		if !hasInfo {
			t.Error("expected INFO alert for 80% threshold")
		}
	})

	t.Run("landing alert when 90-day count < 3", func(t *testing.T) {
		repo := &mockFlightSummaryRepo{
			getSummaryFn: func(ctx context.Context, _, _, _ string) ([]domain.PilotRoleBreakdown, error) {
				return nil, nil
			},
			getLandingsFn: func(ctx context.Context, _, _, _ string) (int, error) {
				return 1, nil // Below minimum
			},
		}
		svc := NewFlightSummaryService(repo)
		alerts, err := svc.BuildFlightAlerts(context.Background(), "emp-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		hasLandingAlert := false
		for _, a := range alerts {
			if a.Type == domain.AlertTypeMinLandings90D {
				hasLandingAlert = true
			}
		}
		if !hasLandingAlert {
			t.Error("expected landing alert")
		}
	})

	t.Run("notice alert when 50-day count < 3 but others OK", func(t *testing.T) {
		callCount := 0
		repo := &mockFlightSummaryRepo{
			getSummaryFn: func(ctx context.Context, _, _, _ string) ([]domain.PilotRoleBreakdown, error) {
				return nil, nil
			},
			getLandingsFn: func(ctx context.Context, _, startDate, _ string) (int, error) {
				callCount++
				// Windows: 50, 60, 70, 80, 90
				// Return: 50d=1, 60d=3, 70d=3, 80d=3, 90d=3
				if callCount <= 5 {
					// 50-day window
					if callCount == 1 {
						return 1, nil // 50d: below 3
					}
					return 3, nil // others: OK
				}
				return 3, nil
			},
		}
		svc := NewFlightSummaryService(repo)
		alerts, err := svc.BuildFlightAlerts(context.Background(), "emp-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		hasNotice := false
		for _, a := range alerts {
			if a.Type == domain.AlertTypeMinLandings90D && a.Severity == domain.AlertSeverityNotice {
				hasNotice = true
			}
		}
		if !hasNotice {
			t.Error("expected NOTICE alert for 50-day window")
		}
	})

	t.Run("summary repo error continues gracefully", func(t *testing.T) {
		repo := &mockFlightSummaryRepo{
			getSummaryFn: func(ctx context.Context, _, _, _ string) ([]domain.PilotRoleBreakdown, error) {
				return nil, errors.New("db error")
			},
			getLandingsFn: func(ctx context.Context, _, _, _ string) (int, error) {
				return 5, nil
			},
		}
		svc := NewFlightSummaryService(repo)
		alerts, err := svc.BuildFlightAlerts(context.Background(), "emp-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// Should not crash, just skip hour alerts
		_ = alerts
	})

	t.Run("landing repo error continues gracefully", func(t *testing.T) {
		repo := &mockFlightSummaryRepo{
			getSummaryFn: func(ctx context.Context, _, _, _ string) ([]domain.PilotRoleBreakdown, error) {
				return nil, nil
			},
			getLandingsFn: func(ctx context.Context, _, _, _ string) (int, error) {
				return 0, errors.New("landing error")
			},
		}
		svc := NewFlightSummaryService(repo)
		alerts, err := svc.BuildFlightAlerts(context.Background(), "emp-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		_ = alerts
	})
}
