package domain

import "testing"

func TestIsValidPeriod(t *testing.T) {
	tests := []struct {
		period string
		want   bool
	}{
		{PeriodMonthly, true},
		{PeriodBimonthly, true},
		{PeriodQuarterly, true},
		{PeriodSemiannual, true},
		{PeriodAnnual, true},
		{PeriodCustom, true},
		{"invalid", false},
		{"", false},
		{"weekly", false},
	}
	for _, tt := range tests {
		t.Run(tt.period, func(t *testing.T) {
			if got := IsValidPeriod(tt.period); got != tt.want {
				t.Errorf("IsValidPeriod(%q) = %v, want %v", tt.period, got, tt.want)
			}
		})
	}
}

func TestFormatMinutesToHHMM(t *testing.T) {
	tests := []struct {
		name    string
		minutes int
		want    string
	}{
		{"zero", 0, "0:00"},
		{"one hour", 60, "1:00"},
		{"90 minutes", 90, "1:30"},
		{"8 hours 15 min", 495, "8:15"},
		{"100 hours", 6000, "100:00"},
		{"single digit minutes", 5, "0:05"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatMinutesToHHMM(tt.minutes); got != tt.want {
				t.Errorf("FormatMinutesToHHMM(%d) = %q, want %q", tt.minutes, got, tt.want)
			}
		})
	}
}

func TestValidPeriods(t *testing.T) {
	if len(ValidPeriods) != 6 {
		t.Errorf("expected 6 valid periods, got %d", len(ValidPeriods))
	}
}

func TestAlertConstants(t *testing.T) {
	// Verify regulatory limits are set correctly
	if LimitHours15Days != 3000 {
		t.Errorf("expected 3000 min for 15-day limit, got %d", LimitHours15Days)
	}
	if LimitHoursMonthly != 5400 {
		t.Errorf("expected 5400 min for monthly limit, got %d", LimitHoursMonthly)
	}
	if LimitHoursQuarter != 16200 {
		t.Errorf("expected 16200 min for quarterly limit, got %d", LimitHoursQuarter)
	}
	if LimitHoursAnnual != 60000 {
		t.Errorf("expected 60000 min for annual limit, got %d", LimitHoursAnnual)
	}
	if MinLandings90Days != 3 {
		t.Errorf("expected 3 min landings, got %d", MinLandings90Days)
	}
	if AlertWarningPercent != 80 {
		t.Errorf("expected 80%% warning threshold, got %d", AlertWarningPercent)
	}
}
