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

// Alert types
const (
	AlertTypeMaxConsecutiveHours = "MAX_CONSECUTIVE_HOURS"
	AlertTypeMinMonthlyLandings  = "MIN_MONTHLY_LANDINGS"
)

// Alert severities
const (
	AlertSeverityWarning = "WARNING"
	AlertSeverityInfo    = "INFO"
)

// Default thresholds (until flight_limitation table is populated)
const (
	DefaultMaxConsecutiveMinutes = 600 // 10 hours
	DefaultMinMonthlyLandings    = 3
)

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
