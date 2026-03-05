package domain

import (
	"time"

	"github.com/google/uuid"
)

// EmployeeFlightSummary stores accumulated flight data per period
type EmployeeFlightSummary struct {
	ID              string `json:"id"`
	EmployeeID      string `json:"employee_id"`
	PeriodType      string `json:"period_type"`
	PeriodYear      int    `json:"period_year"`
	PeriodNumber    int    `json:"period_number"`
	PeriodStart     string `json:"period_start"`
	PeriodEnd       string `json:"period_end"`
	TotalAirTime    int    `json:"total_air_time"`   // minutes
	TotalBlockTime  int    `json:"total_block_time"` // minutes
	TotalFlights    int    `json:"total_flights"`
	TotalLandings   int    `json:"total_landings"`
	CatApproaches   int    `json:"cat_approaches"`
	LastCatApproach string `json:"last_cat_approach_date,omitempty"`
	LastUpdated     string `json:"last_updated,omitempty"`
}

// SetID generates a new UUID
func (s *EmployeeFlightSummary) SetID() {
	s.ID = uuid.New().String()
}

// Summary period types — match DB ENUM values
const (
	SummaryPeriodMonthly    = "MONTHLY"
	SummaryPeriodQuarterly  = "QUARTERLY"
	SummaryPeriodAnnual     = "ANNUAL"
	SummaryPeriodFirst15    = "PERIOD_1_15"
	SummaryPeriodSecondHalf = "PERIOD_16_31"
)

// AllSummaryPeriods lists every period type we track
var AllSummaryPeriods = []string{
	SummaryPeriodFirst15,
	SummaryPeriodSecondHalf,
	SummaryPeriodMonthly,
	SummaryPeriodQuarterly,
	SummaryPeriodAnnual,
}

// PeriodInfo describes a single period row to upsert
type PeriodInfo struct {
	PeriodType   string
	PeriodYear   int
	PeriodNumber int
	PeriodStart  string // YYYY-MM-DD
	PeriodEnd    string // YYYY-MM-DD
}

// SummaryDelta holds the delta values for a flight summary upsert
type SummaryDelta struct {
	AirTime   int
	BlockTime int
	Flights   int
	Landings  int
}

// dateFormatISO is the standard ISO date format used for period boundaries
const dateFormatISO = "2006-01-02"

// GetAffectedPeriods returns all period rows that contain a given flight date.
// For example, 2026-01-10 affects: PERIOD_1_15 (Jan #1), MONTHLY (Jan), QUARTERLY (Q1), ANNUAL (2026).
func GetAffectedPeriods(flightDate time.Time) []PeriodInfo {
	year := flightDate.Year()
	month := flightDate.Month()
	day := flightDate.Day()

	var periods []PeriodInfo

	// 1. PERIOD_1_15 or PERIOD_16_31
	if day <= 15 {
		start := time.Date(year, month, 1, 0, 0, 0, 0, flightDate.Location())
		end := time.Date(year, month, 15, 0, 0, 0, 0, flightDate.Location())
		periods = append(periods, PeriodInfo{
			PeriodType:   SummaryPeriodFirst15,
			PeriodYear:   year,
			PeriodNumber: int(month)*2 - 1, // Jan→1, Feb→3, Mar→5 ...
			PeriodStart:  start.Format(dateFormatISO),
			PeriodEnd:    end.Format(dateFormatISO),
		})
	} else {
		start := time.Date(year, month, 16, 0, 0, 0, 0, flightDate.Location())
		end := time.Date(year, month+1, 0, 0, 0, 0, 0, flightDate.Location()) // last day of month
		periods = append(periods, PeriodInfo{
			PeriodType:   SummaryPeriodSecondHalf,
			PeriodYear:   year,
			PeriodNumber: int(month) * 2, // Jan→2, Feb→4, Mar→6 ...
			PeriodStart:  start.Format(dateFormatISO),
			PeriodEnd:    end.Format(dateFormatISO),
		})
	}

	// 2. MONTHLY
	monthStart := time.Date(year, month, 1, 0, 0, 0, 0, flightDate.Location())
	monthEnd := time.Date(year, month+1, 0, 0, 0, 0, 0, flightDate.Location())
	periods = append(periods, PeriodInfo{
		PeriodType:   SummaryPeriodMonthly,
		PeriodYear:   year,
		PeriodNumber: int(month),
		PeriodStart:  monthStart.Format(dateFormatISO),
		PeriodEnd:    monthEnd.Format(dateFormatISO),
	})

	// 3. QUARTERLY (Q1=Jan-Mar, Q2=Apr-Jun, Q3=Jul-Sep, Q4=Oct-Dec)
	quarter := (int(month) - 1) / 3
	quarterNumber := quarter + 1
	qStart := time.Date(year, time.Month(quarter*3+1), 1, 0, 0, 0, 0, flightDate.Location())
	qEnd := time.Date(year, time.Month(quarter*3+4), 0, 0, 0, 0, 0, flightDate.Location())
	periods = append(periods, PeriodInfo{
		PeriodType:   SummaryPeriodQuarterly,
		PeriodYear:   year,
		PeriodNumber: quarterNumber,
		PeriodStart:  qStart.Format(dateFormatISO),
		PeriodEnd:    qEnd.Format(dateFormatISO),
	})

	// 4. ANNUAL
	periods = append(periods, PeriodInfo{
		PeriodType:   SummaryPeriodAnnual,
		PeriodYear:   year,
		PeriodNumber: 1,
		PeriodStart:  time.Date(year, 1, 1, 0, 0, 0, 0, flightDate.Location()).Format(dateFormatISO),
		PeriodEnd:    time.Date(year, 12, 31, 0, 0, 0, 0, flightDate.Location()).Format(dateFormatISO),
	})

	return periods
}

// IsLandingRole returns true if the pilot_role counts as a landing
// PF (Pilot Flying) and PFL (Pilot Flying Landing) count as landings
func IsLandingRole(role *PilotRole) bool {
	if role == nil {
		return false
	}
	return *role == PilotRolePF || *role == PilotRolePFL
}
