package domain

import (
	"testing"
	"time"
)

// ═══════════════════════════════════════════
// Tests for employee_flight_summary.go
// ═══════════════════════════════════════════

func TestEmployeeFlightSummary_SetID(t *testing.T) {
	s := &EmployeeFlightSummary{}
	s.SetID()
	if s.ID == "" {
		t.Error("expected non-empty ID after SetID()")
	}
}

func TestGetAffectedPeriods_FirstHalf(t *testing.T) {
	// Jan 10, 2026 → should produce: PERIOD_1_15, MONTHLY, QUARTERLY, ANNUAL
	date := time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)
	periods := GetAffectedPeriods(date)

	if len(periods) != 4 {
		t.Fatalf("expected 4 periods, got %d", len(periods))
	}

	// Check PERIOD_1_15
	p := periods[0]
	if p.PeriodType != SummaryPeriodFirst15 {
		t.Errorf("expected %s, got %s", SummaryPeriodFirst15, p.PeriodType)
	}
	if p.PeriodYear != 2026 {
		t.Errorf("expected year 2026, got %d", p.PeriodYear)
	}
	if p.PeriodNumber != 1 {
		t.Errorf("expected period number 1, got %d", p.PeriodNumber)
	}
	if p.PeriodStart != "2026-01-01" {
		t.Errorf("expected start 2026-01-01, got %s", p.PeriodStart)
	}
	if p.PeriodEnd != "2026-01-15" {
		t.Errorf("expected end 2026-01-15, got %s", p.PeriodEnd)
	}

	// Check MONTHLY
	if periods[1].PeriodType != SummaryPeriodMonthly {
		t.Errorf("expected %s, got %s", SummaryPeriodMonthly, periods[1].PeriodType)
	}
	if periods[1].PeriodNumber != 1 {
		t.Errorf("expected monthly period number 1, got %d", periods[1].PeriodNumber)
	}

	// Check QUARTERLY
	if periods[2].PeriodType != SummaryPeriodQuarterly {
		t.Errorf("expected %s, got %s", SummaryPeriodQuarterly, periods[2].PeriodType)
	}
	if periods[2].PeriodNumber != 1 {
		t.Errorf("expected quarterly number 1 (Q1), got %d", periods[2].PeriodNumber)
	}

	// Check ANNUAL
	if periods[3].PeriodType != SummaryPeriodAnnual {
		t.Errorf("expected %s, got %s", SummaryPeriodAnnual, periods[3].PeriodType)
	}
	if periods[3].PeriodStart != "2026-01-01" || periods[3].PeriodEnd != "2026-12-31" {
		t.Errorf("expected annual range 2026-01-01 to 2026-12-31, got %s to %s", periods[3].PeriodStart, periods[3].PeriodEnd)
	}
}

func TestGetAffectedPeriods_SecondHalf(t *testing.T) {
	// Jan 20, 2026 → should produce: PERIOD_16_31, MONTHLY, QUARTERLY, ANNUAL
	date := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)
	periods := GetAffectedPeriods(date)

	if len(periods) != 4 {
		t.Fatalf("expected 4 periods, got %d", len(periods))
	}

	p := periods[0]
	if p.PeriodType != SummaryPeriodSecondHalf {
		t.Errorf("expected %s, got %s", SummaryPeriodSecondHalf, p.PeriodType)
	}
	if p.PeriodNumber != 2 {
		t.Errorf("expected period number 2, got %d", p.PeriodNumber)
	}
	if p.PeriodStart != "2026-01-16" {
		t.Errorf("expected start 2026-01-16, got %s", p.PeriodStart)
	}
}

func TestGetAffectedPeriods_Q2(t *testing.T) {
	// Apr 5, 2026 → Q2
	date := time.Date(2026, 4, 5, 0, 0, 0, 0, time.UTC)
	periods := GetAffectedPeriods(date)

	// Should have PERIOD_1_15, MONTHLY(4), QUARTERLY(2), ANNUAL
	if periods[2].PeriodNumber != 2 {
		t.Errorf("expected quarterly number 2 (Q2), got %d", periods[2].PeriodNumber)
	}
}

func TestGetAffectedPeriods_Q3(t *testing.T) {
	date := time.Date(2026, 8, 18, 0, 0, 0, 0, time.UTC)
	periods := GetAffectedPeriods(date)

	// PERIOD_16_31 + MONTHLY(8) + QUARTERLY(3) + ANNUAL
	if periods[0].PeriodType != SummaryPeriodSecondHalf {
		t.Errorf("expected PERIOD_16_31, got %s", periods[0].PeriodType)
	}
	if periods[2].PeriodNumber != 3 {
		t.Errorf("expected Q3, got %d", periods[2].PeriodNumber)
	}
}

func TestGetAffectedPeriods_Q4(t *testing.T) {
	date := time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC)
	periods := GetAffectedPeriods(date)
	if periods[2].PeriodNumber != 4 {
		t.Errorf("expected Q4, got %d", periods[2].PeriodNumber)
	}
}

func TestGetAffectedPeriods_Day15Boundary(t *testing.T) {
	// Day 15 should be in PERIOD_1_15
	date := time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC)
	periods := GetAffectedPeriods(date)
	if periods[0].PeriodType != SummaryPeriodFirst15 {
		t.Errorf("day 15 should be PERIOD_1_15, got %s", periods[0].PeriodType)
	}
}

func TestGetAffectedPeriods_Day16Boundary(t *testing.T) {
	// Day 16 should be in PERIOD_16_31
	date := time.Date(2026, 3, 16, 0, 0, 0, 0, time.UTC)
	periods := GetAffectedPeriods(date)
	if periods[0].PeriodType != SummaryPeriodSecondHalf {
		t.Errorf("day 16 should be PERIOD_16_31, got %s", periods[0].PeriodType)
	}
}

func TestIsLandingRole(t *testing.T) {
	tests := []struct {
		name string
		role *PilotRole
		want bool
	}{
		{"nil role", nil, false},
		{"PF", func() *PilotRole { r := PilotRolePF; return &r }(), true},
		{"PFL", func() *PilotRole { r := PilotRolePFL; return &r }(), true},
		{"PM", func() *PilotRole { r := PilotRolePM; return &r }(), false},
		{"PFTO", func() *PilotRole { r := PilotRolePFTO; return &r }(), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsLandingRole(tt.role); got != tt.want {
				t.Errorf("IsLandingRole() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAllSummaryPeriods(t *testing.T) {
	if len(AllSummaryPeriods) != 5 {
		t.Errorf("expected 5 summary periods, got %d", len(AllSummaryPeriods))
	}
}
