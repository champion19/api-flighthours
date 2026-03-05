package domain

import "fmt"

// FlightHoursSummary contains the aggregated flight hours for a given period
type FlightHoursSummary struct {
	Period        string         `json:"period"`
	StartDate     string         `json:"start_date"`
	EndDate       string         `json:"end_date"`
	TotalMinutes  int            `json:"total_minutes"`
	TotalFlights  int            `json:"total_flights"`
	TotalLandings int            `json:"total_landings"`
	Breakdown     map[string]int `json:"breakdown"` // pilot_role → total minutes
}

// PilotRoleBreakdown is a single row from the aggregation query
type PilotRoleBreakdown struct {
	PilotRole    string `json:"pilot_role"`
	TotalSeconds int    `json:"total_seconds"`
	FlightCount  int    `json:"flight_count"`
}

// FlightAlert represents a single alert for the employee dashboard
type FlightAlert struct {
	Type         string `json:"type"`
	Severity     string `json:"severity"`
	Message      string `json:"message"`
	CurrentValue int    `json:"current_value"`
	Threshold    int    `json:"threshold"`
}

// Alert types — real regulatory rules
const (
	AlertTypeHourLimit15Days  = "HOUR_LIMIT_15_DAYS"
	AlertTypeHourLimitMonthly = "HOUR_LIMIT_MONTHLY"
	AlertTypeHourLimitQuarter = "HOUR_LIMIT_QUARTERLY"
	AlertTypeHourLimitAnnual  = "HOUR_LIMIT_ANNUAL"
	AlertTypeMinLandings90D   = "MIN_LANDINGS_90_DAYS"
)

// Alert severities
const (
	AlertSeverityWarning = "WARNING"
	AlertSeverityInfo    = "INFO"
	AlertSeverityNotice  = "NOTICE" // Neutral/gray — informational, no stress
)

// Regulatory flight hour limits (in minutes)
const (
	LimitHours15Days  = 50 * 60   // 3000 min = 50h max in a 15-day calendar period
	LimitHoursMonthly = 90 * 60   // 5400 min = 90h max in a calendar month
	LimitHoursQuarter = 270 * 60  // 16200 min = 270h max in a calendar quarter
	LimitHoursAnnual  = 1000 * 60 // 60000 min = 1000h max in a calendar year
	MinLandings90Days = 3         // 3 landings in 90 rolling days for PF/PFL
)

// AlertWarningPercent is the threshold percentage at which INFO alerts are triggered
const AlertWarningPercent = 80

// Valid period types for flight hours summary
const (
	PeriodMonthly    = "monthly"
	PeriodBimonthly  = "bimonthly"
	PeriodQuarterly  = "quarterly"
	PeriodSemiannual = "semiannual"
	PeriodAnnual     = "annual"
	PeriodCustom     = "custom"
)

// ValidPeriods contains all valid period types
var ValidPeriods = []string{
	PeriodMonthly,
	PeriodBimonthly,
	PeriodQuarterly,
	PeriodSemiannual,
	PeriodAnnual,
	PeriodCustom,
}

// IsValidPeriod checks if a string is a valid period type
func IsValidPeriod(period string) bool {
	for _, p := range ValidPeriods {
		if p == period {
			return true
		}
	}
	return false
}

// FormatMinutesToHHMM converts minutes to HH:MM format
func FormatMinutesToHHMM(totalMinutes int) string {
	hours := totalMinutes / 60
	minutes := totalMinutes % 60
	return fmt.Sprintf("%d:%02d", hours, minutes)
}
