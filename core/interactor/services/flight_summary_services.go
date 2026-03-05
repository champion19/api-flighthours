package services

import (
	"context"
	"fmt"
	"time"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/logger"
)

type flightSummaryService struct {
	repo output.FlightSummaryRepository
}

const dateFormat = "2006-01-02"

func NewFlightSummaryService(repo output.FlightSummaryRepository) *flightSummaryService {
	return &flightSummaryService{repo: repo}
}

// GetFlightHoursSummary builds the full summary by calling the repo and aggregating results
func (s *flightSummaryService) GetFlightHoursSummary(ctx context.Context, employeeID, startDate, endDate string) (*domain.FlightHoursSummary, error) {
	log.Info(logger.LogFlightSummaryGet, "employee_id", employeeID, "start", startDate, "end", endDate)

	breakdown, err := s.repo.GetFlightHoursSummary(ctx, employeeID, startDate, endDate)
	if err != nil {
		log.Error(logger.LogFlightSummaryGetError, "error", err)
		return nil, err
	}

	summary := &domain.FlightHoursSummary{
		StartDate: startDate,
		EndDate:   endDate,
		Breakdown: make(map[string]int),
	}

	for _, b := range breakdown {
		minutes := b.TotalSeconds / 60
		summary.Breakdown[b.PilotRole] = minutes
		summary.TotalMinutes += minutes
		summary.TotalFlights += b.FlightCount
	}

	// Get landing count for the same period
	landingCount, err := s.repo.GetLandingCount(ctx, employeeID, startDate, endDate)
	if err != nil {
		log.Error(logger.LogFlightSummaryGetError, "action", "landing_count", "error", err)
		return nil, err
	}
	summary.TotalLandings = landingCount

	log.Info(logger.LogFlightSummaryGetOK, "total_minutes", summary.TotalMinutes, "total_flights", summary.TotalFlights)
	return summary, nil
}

// GetRecentFlights delegates to the repository
func (s *flightSummaryService) GetRecentFlights(ctx context.Context, employeeID string, limit int) ([]domain.DailyLogbookDetail, error) {
	return s.repo.GetRecentFlights(ctx, employeeID, limit)
}

// GetLandingCount delegates to the repository
func (s *flightSummaryService) GetLandingCount(ctx context.Context, employeeID, startDate, endDate string) (int, error) {
	return s.repo.GetLandingCount(ctx, employeeID, startDate, endDate)
}

// GetDailyFlightSeconds delegates to the repository
func (s *flightSummaryService) GetDailyFlightSeconds(ctx context.Context, employeeID, date string) (int, error) {
	return s.repo.GetDailyFlightSeconds(ctx, employeeID, date)
}

// CalculatePeriodDates computes start and end dates for a given period type.
// This is business logic and therefore lives in the service layer.
func (s *flightSummaryService) CalculatePeriodDates(period, referenceDate string) (string, string, error) {
	log.Info(logger.LogFlightSummaryPeriodCalc, "period", period, "reference_date", referenceDate)

	if !domain.IsValidPeriod(period) {
		return "", "", fmt.Errorf("invalid period: %s", period)
	}

	var ref time.Time
	var err error

	if referenceDate != "" {
		ref, err = time.Parse(dateFormat, referenceDate)
		if err != nil {
			return "", "", domain.ErrInvalidDateFormat
		}
	} else {
		ref = time.Now()
	}

	var start, end time.Time

	switch period {
	case domain.PeriodMonthly:
		start = time.Date(ref.Year(), ref.Month(), 1, 0, 0, 0, 0, ref.Location())
		end = start.AddDate(0, 1, -1)

	case domain.PeriodBimonthly:
		// Current bimonth: Jan-Feb, Mar-Apr, May-Jun, Jul-Aug, Sep-Oct, Nov-Dec
		bimonthStart := ((ref.Month()-1)/2)*2 + 1
		start = time.Date(ref.Year(), time.Month(bimonthStart), 1, 0, 0, 0, 0, ref.Location())
		end = start.AddDate(0, 2, -1)

	case domain.PeriodQuarterly:
		quarterStart := ((ref.Month()-1)/3)*3 + 1
		start = time.Date(ref.Year(), time.Month(quarterStart), 1, 0, 0, 0, 0, ref.Location())
		end = start.AddDate(0, 3, -1)

	case domain.PeriodSemiannual:
		if ref.Month() <= 6 {
			start = time.Date(ref.Year(), 1, 1, 0, 0, 0, 0, ref.Location())
			end = time.Date(ref.Year(), 6, 30, 0, 0, 0, 0, ref.Location())
		} else {
			start = time.Date(ref.Year(), 7, 1, 0, 0, 0, 0, ref.Location())
			end = time.Date(ref.Year(), 12, 31, 0, 0, 0, 0, ref.Location())
		}

	case domain.PeriodAnnual:
		start = time.Date(ref.Year(), 1, 1, 0, 0, 0, 0, ref.Location())
		end = time.Date(ref.Year(), 12, 31, 0, 0, 0, 0, ref.Location())

	case domain.PeriodCustom:
		return "", "", fmt.Errorf("custom period requires explicit start_date and end_date")
	}

	return start.Format(dateFormat), end.Format(dateFormat), nil
}

// BuildFlightAlerts evaluates all regulatory alert conditions.
// Uses employee_flight_summary rows for hour limits and GetLandingCount for PF/PFL landing currency.
func (s *flightSummaryService) BuildFlightAlerts(ctx context.Context, employeeID string) ([]domain.FlightAlert, error) {
	log.Info(logger.LogFlightSummaryAlerts, "employee_id", employeeID)

	var alerts []domain.FlightAlert
	now := time.Now()

	// Get all affected periods for today to determine current period numbers
	periods := domain.GetAffectedPeriods(now)

	// Build a lookup map of current period info
	periodMap := make(map[string]domain.PeriodInfo)
	for _, p := range periods {
		periodMap[p.PeriodType] = p
	}

	// Define which alerts to check
	type alertCheck struct {
		periodType string
		alertType  string
		limit      int
		label      string
	}

	checks := []alertCheck{
		{domain.SummaryPeriodFirst15, domain.AlertTypeHourLimit15Days, domain.LimitHours15Days, "15-day period"},
		{domain.SummaryPeriodSecondHalf, domain.AlertTypeHourLimit15Days, domain.LimitHours15Days, "15-day period"},
		{domain.SummaryPeriodMonthly, domain.AlertTypeHourLimitMonthly, domain.LimitHoursMonthly, "month"},
		{domain.SummaryPeriodQuarterly, domain.AlertTypeHourLimitQuarter, domain.LimitHoursQuarter, "quarter"},
		{domain.SummaryPeriodAnnual, domain.AlertTypeHourLimitAnnual, domain.LimitHoursAnnual, "year"},
	}

	// Check flight hour limits against employee_flight_summary
	for _, check := range checks {
		pInfo, ok := periodMap[check.periodType]
		if !ok {
			continue // This period type is not active for today's date
		}

		// Get daily flight seconds as a proxy for current accumulated hours
		// Since employee_flight_summary may not be populated yet, fall back to repo query
		currentMinutes := 0

		// Try to read from the aggregation query (real-time data from daily_logbook_detail)
		breakdown, err := s.repo.GetFlightHoursSummary(ctx, employeeID, pInfo.PeriodStart, pInfo.PeriodEnd)
		if err != nil {
			log.Error(logger.LogFlightSummaryAlertsError, "action", "get_summary_for_alert", "period", check.periodType, "error", err)
			continue
		}

		for _, b := range breakdown {
			currentMinutes += b.TotalSeconds / 60
		}

		warningThreshold := check.limit * domain.AlertWarningPercent / 100
		if currentMinutes >= warningThreshold {
			severity := domain.AlertSeverityInfo
			if currentMinutes >= check.limit {
				severity = domain.AlertSeverityWarning
			}
			alerts = append(alerts, domain.FlightAlert{
				Type:         check.alertType,
				Severity:     severity,
				Message:      fmt.Sprintf("You have flown %s this %s (max: %s)", domain.FormatMinutesToHHMM(currentMinutes), check.label, domain.FormatMinutesToHHMM(check.limit)),
				CurrentValue: currentMinutes,
				Threshold:    check.limit,
			})
		}
	}

	// Check landing currency (multiple rolling windows)
	landingAlert := s.buildLandingCurrencyAlert(ctx, employeeID, now)
	if landingAlert != nil {
		alerts = append(alerts, *landingAlert)
	}

	log.Info(logger.LogFlightSummaryAlertsOK, "alert_count", len(alerts))
	return alerts, nil
}

// buildLandingCurrencyAlert evaluates the 90-day rolling landing currency with
// 3-phase progressive alerts. Returns nil if the pilot has enough landings.
func (s *flightSummaryService) buildLandingCurrencyAlert(ctx context.Context, employeeID string, now time.Time) *domain.FlightAlert {
	todayStr := now.Format(dateFormat)

	windows := []int{50, 60, 70, 80, 90}
	landingCounts := make(map[int]int)
	for _, w := range windows {
		start := now.AddDate(0, 0, -w).Format(dateFormat)
		count, err := s.repo.GetLandingCount(ctx, employeeID, start, todayStr)
		if err != nil {
			log.Error(logger.LogFlightSummaryAlertsError, "action", fmt.Sprintf("landing_count_%dd", w), "error", err)
			continue
		}
		landingCounts[w] = count
	}

	count90 := landingCounts[90]
	count80 := landingCounts[80]
	count70 := landingCounts[70]
	count60 := landingCounts[60]
	count50 := landingCounts[50]

	// All windows pass — no alert
	if count90 >= domain.MinLandings90Days &&
		count80 >= domain.MinLandings90Days &&
		count70 >= domain.MinLandings90Days &&
		count60 >= domain.MinLandings90Days &&
		count50 >= domain.MinLandings90Days {
		return nil
	}

	// Phase 3: 90-day window failed → WARNING (simulator required)
	if count90 < domain.MinLandings90Days {
		return &domain.FlightAlert{
			Type:         domain.AlertTypeMinLandings90D,
			Severity:     domain.AlertSeverityWarning,
			Message:      fmt.Sprintf("You have %d PF/PFL landings in the last 90 days (minimum: %d). Simulator recurrency training may be required.", count90, domain.MinLandings90Days),
			CurrentValue: count90,
			Threshold:    domain.MinLandings90Days,
		}
	}

	// Phase 2: 60-80 day windows failing → INFO (approaching expiry)
	if count80 < domain.MinLandings90Days || count70 < domain.MinLandings90Days || count60 < domain.MinLandings90Days {
		windowDays, daysRemaining, countUsed := 60, 30, count60
		if count70 < domain.MinLandings90Days {
			windowDays, daysRemaining, countUsed = 70, 20, count70
		}
		if count80 < domain.MinLandings90Days {
			windowDays, daysRemaining, countUsed = 80, 10, count80
		}
		return &domain.FlightAlert{
			Type:         domain.AlertTypeMinLandings90D,
			Severity:     domain.AlertSeverityInfo,
			Message:      fmt.Sprintf("You have %d PF/PFL landings in the last %d days. Some landings will expire in ~%d days — complete %d landings to stay current.", countUsed, windowDays, daysRemaining, domain.MinLandings90Days),
			CurrentValue: countUsed,
			Threshold:    domain.MinLandings90Days,
		}
	}

	// Phase 1: 50-day window failing → NOTICE (gentle, no stress)
	if count50 < domain.MinLandings90Days {
		return &domain.FlightAlert{
			Type:         domain.AlertTypeMinLandings90D,
			Severity:     domain.AlertSeverityNotice,
			Message:      fmt.Sprintf("You have %d PF/PFL landings in the last 50 days. Aim for %d to maintain landing currency.", count50, domain.MinLandings90Days),
			CurrentValue: count50,
			Threshold:    domain.MinLandings90Days,
		}
	}

	return nil
}
