package services

import (
	"context"
	"errors"
	"testing"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
)

// ═══════════════════════════════════════════
// Mock for EmployeeFlightSummaryRepository
// ═══════════════════════════════════════════

type mockEmpFlightSummaryRepo struct {
	beginTxFn          func(ctx context.Context) (output.Tx, error)
	upsertFn           func(ctx context.Context, tx output.Tx, employeeID string, period domain.PeriodInfo, delta domain.SummaryDelta) error
	getByEmployeeFn    func(ctx context.Context, employeeID, periodType string) ([]domain.EmployeeFlightSummary, error)
	getCurrentPeriodFn func(ctx context.Context, employeeID, periodType string, year, number int) (*domain.EmployeeFlightSummary, error)
	getAllFn           func(ctx context.Context, employeeID string) ([]domain.EmployeeFlightSummary, error)
}

func (m *mockEmpFlightSummaryRepo) BeginTx(ctx context.Context) (output.Tx, error) {
	if m.beginTxFn != nil {
		return m.beginTxFn(ctx)
	}
	return &mockTx{}, nil
}

func (m *mockEmpFlightSummaryRepo) UpsertSummary(ctx context.Context, tx output.Tx, employeeID string, period domain.PeriodInfo, delta domain.SummaryDelta) error {
	if m.upsertFn != nil {
		return m.upsertFn(ctx, tx, employeeID, period, delta)
	}
	return nil
}

func (m *mockEmpFlightSummaryRepo) GetSummariesByEmployee(ctx context.Context, employeeID, periodType string) ([]domain.EmployeeFlightSummary, error) {
	if m.getByEmployeeFn != nil {
		return m.getByEmployeeFn(ctx, employeeID, periodType)
	}
	return nil, nil
}

func (m *mockEmpFlightSummaryRepo) GetCurrentPeriodSummary(ctx context.Context, employeeID, periodType string, year, number int) (*domain.EmployeeFlightSummary, error) {
	if m.getCurrentPeriodFn != nil {
		return m.getCurrentPeriodFn(ctx, employeeID, periodType, year, number)
	}
	return nil, nil
}

func (m *mockEmpFlightSummaryRepo) GetAllSummaries(ctx context.Context, employeeID string) ([]domain.EmployeeFlightSummary, error) {
	if m.getAllFn != nil {
		return m.getAllFn(ctx, employeeID)
	}
	return nil, nil
}

// ═══════════════════════════════════════════
// Tests
// ═══════════════════════════════════════════

func TestNewEmployeeFlightSummaryService(t *testing.T) {
	svc := NewEmployeeFlightSummaryService(&mockEmpFlightSummaryRepo{})
	if svc == nil {
		t.Error("expected non-nil service")
	}
}

func TestEmployeeFlightSummaryService_AccumulateFlightHours(t *testing.T) {
	t.Run("success creation", func(t *testing.T) {
		upsertCalls := 0
		repo := &mockEmpFlightSummaryRepo{
			upsertFn: func(ctx context.Context, tx output.Tx, employeeID string, period domain.PeriodInfo, delta domain.SummaryDelta) error {
				upsertCalls++
				if delta.AirTime < 0 {
					t.Error("expected positive air time for creation")
				}
				return nil
			},
		}
		svc := NewEmployeeFlightSummaryService(repo)
		airTime := "02:30:00"
		blockTime := "03:00:00"
		pf := domain.PilotRolePF
		err := svc.AccumulateFlightHours(context.Background(), &mockTx{}, "emp-1", domain.DailyLogbookDetail{
			FlightRealDate: "2026-01-10",
			AirTime:        &airTime,
			BlockTime:      &blockTime,
			PilotRole:      &pf,
		}, false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if upsertCalls != 4 { // PERIOD_1_15, MONTHLY, QUARTERLY, ANNUAL
			t.Errorf("expected 4 upsert calls, got %d", upsertCalls)
		}
	})

	t.Run("success deletion negates deltas", func(t *testing.T) {
		repo := &mockEmpFlightSummaryRepo{
			upsertFn: func(ctx context.Context, tx output.Tx, employeeID string, period domain.PeriodInfo, delta domain.SummaryDelta) error {
				if delta.AirTime > 0 {
					t.Error("expected negative air time for deletion")
				}
				if delta.Flights > 0 {
					t.Error("expected negative flights delta for deletion")
				}
				return nil
			},
		}
		svc := NewEmployeeFlightSummaryService(repo)
		airTime := "01:00:00"
		pm := domain.PilotRolePM
		err := svc.AccumulateFlightHours(context.Background(), &mockTx{}, "emp-1", domain.DailyLogbookDetail{
			FlightRealDate: "2026-01-10",
			AirTime:        &airTime,
			PilotRole:      &pm,
		}, true)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("invalid date", func(t *testing.T) {
		svc := NewEmployeeFlightSummaryService(&mockEmpFlightSummaryRepo{})
		err := svc.AccumulateFlightHours(context.Background(), &mockTx{}, "emp-1", domain.DailyLogbookDetail{
			FlightRealDate: "not-a-date",
		}, false)
		if err == nil {
			t.Error("expected error for invalid date")
		}
	})

	t.Run("upsert error", func(t *testing.T) {
		repo := &mockEmpFlightSummaryRepo{
			upsertFn: func(ctx context.Context, tx output.Tx, employeeID string, period domain.PeriodInfo, delta domain.SummaryDelta) error {
				return errors.New("upsert failed")
			},
		}
		svc := NewEmployeeFlightSummaryService(repo)
		err := svc.AccumulateFlightHours(context.Background(), &mockTx{}, "emp-1", domain.DailyLogbookDetail{
			FlightRealDate: "2026-01-10",
		}, false)
		if err == nil {
			t.Error("expected error from upsert")
		}
	})

	t.Run("nil time fields", func(t *testing.T) {
		repo := &mockEmpFlightSummaryRepo{
			upsertFn: func(ctx context.Context, tx output.Tx, employeeID string, period domain.PeriodInfo, delta domain.SummaryDelta) error {
				if delta.AirTime != 0 {
					t.Errorf("expected 0 air time for nil, got %d", delta.AirTime)
				}
				return nil
			},
		}
		svc := NewEmployeeFlightSummaryService(repo)
		err := svc.AccumulateFlightHours(context.Background(), &mockTx{}, "emp-1", domain.DailyLogbookDetail{
			FlightRealDate: "2026-01-10",
		}, false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("landing role counts landing", func(t *testing.T) {
		repo := &mockEmpFlightSummaryRepo{
			upsertFn: func(ctx context.Context, tx output.Tx, employeeID string, period domain.PeriodInfo, delta domain.SummaryDelta) error {
				if delta.Landings != 1 {
					t.Errorf("expected landing delta 1 for PFL role, got %d", delta.Landings)
				}
				return nil
			},
		}
		svc := NewEmployeeFlightSummaryService(repo)
		pfl := domain.PilotRolePFL
		err := svc.AccumulateFlightHours(context.Background(), &mockTx{}, "emp-1", domain.DailyLogbookDetail{
			FlightRealDate: "2026-01-10",
			PilotRole:      &pfl,
		}, false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestEmployeeFlightSummaryService_Delegates(t *testing.T) {
	t.Run("GetSummariesByEmployee", func(t *testing.T) {
		repo := &mockEmpFlightSummaryRepo{
			getByEmployeeFn: func(ctx context.Context, employeeID, periodType string) ([]domain.EmployeeFlightSummary, error) {
				return []domain.EmployeeFlightSummary{{ID: "s-1"}}, nil
			},
		}
		svc := NewEmployeeFlightSummaryService(repo)
		result, err := svc.GetSummariesByEmployee(context.Background(), "emp-1", "MONTHLY")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 1 {
			t.Errorf("expected 1, got %d", len(result))
		}
	})

	t.Run("GetCurrentPeriodSummary", func(t *testing.T) {
		repo := &mockEmpFlightSummaryRepo{
			getCurrentPeriodFn: func(ctx context.Context, employeeID, periodType string, year, number int) (*domain.EmployeeFlightSummary, error) {
				return &domain.EmployeeFlightSummary{ID: "s-1"}, nil
			},
		}
		svc := NewEmployeeFlightSummaryService(repo)
		result, err := svc.GetCurrentPeriodSummary(context.Background(), "emp-1", "MONTHLY", 2026, 1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.ID != "s-1" {
			t.Errorf("expected ID s-1, got %s", result.ID)
		}
	})

	t.Run("GetAllSummaries", func(t *testing.T) {
		repo := &mockEmpFlightSummaryRepo{
			getAllFn: func(ctx context.Context, employeeID string) ([]domain.EmployeeFlightSummary, error) {
				return []domain.EmployeeFlightSummary{{ID: "s-1"}, {ID: "s-2"}}, nil
			},
		}
		svc := NewEmployeeFlightSummaryService(repo)
		result, err := svc.GetAllSummaries(context.Background(), "emp-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 2 {
			t.Errorf("expected 2, got %d", len(result))
		}
	})
}

// ═══════════════════════════════════════════
// Tests: parseTimeToMinutes
// ═══════════════════════════════════════════

func Test_parseTimeToMinutes(t *testing.T) {
	tests := []struct {
		name string
		time *string
		want int
	}{
		{"nil", nil, 0},
		{"empty", strPtrLocal(""), 0},
		{"HH:MM:SS", strPtrLocal("02:30:00"), 150},
		{"HH:MM", strPtrLocal("01:15"), 75},
		{"invalid", strPtrLocal("not-time"), 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseTimeToMinutes(tt.time); got != tt.want {
				t.Errorf("parseTimeToMinutes() = %d, want %d", got, tt.want)
			}
		})
	}
}

func strPtrLocal(s string) *string { return &s }
