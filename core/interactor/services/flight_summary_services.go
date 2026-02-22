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
		ref, err = time.Parse("2006-01-02", referenceDate)
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

	return start.Format("2006-01-02"), end.Format("2006-01-02"), nil
}

// BuildFlightAlerts evaluates alert conditions and returns any active alerts.
// Business logic: checks consecutive hours today and monthly landing count.
func (s *flightSummaryService) BuildFlightAlerts(ctx context.Context, employeeID string) ([]domain.FlightAlert, error) {
	log.Info(logger.LogFlightSummaryAlerts, "employee_id", employeeID)

	var alerts []domain.FlightAlert
	today := time.Now().Format("2006-01-02")

	// Alert 1: Consecutive hours today
	dailySeconds, err := s.repo.GetDailyFlightSeconds(ctx, employeeID, today)
	if err != nil {
		log.Error(logger.LogFlightSummaryAlertsError, "action", "daily_seconds", "error", err)
		return nil, err
	}

	dailyMinutes := dailySeconds / 60
	if dailyMinutes >= domain.DefaultMaxConsecutiveMinutes*80/100 { // Alert at 80% threshold
		severity := domain.AlertSeverityInfo
		if dailyMinutes >= domain.DefaultMaxConsecutiveMinutes {
			severity = domain.AlertSeverityWarning
		}
		alerts = append(alerts, domain.FlightAlert{
			Type:         domain.AlertTypeMaxConsecutiveHours,
			Severity:     severity,
			Message:      fmt.Sprintf("You have flown %s today (max: %s)", domain.FormatMinutesToHHMM(dailyMinutes), domain.FormatMinutesToHHMM(domain.DefaultMaxConsecutiveMinutes)),
			CurrentValue: dailyMinutes,
			Threshold:    domain.DefaultMaxConsecutiveMinutes,
		})
	}

	// Alert 2: Monthly landing count
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	monthEnd := monthStart.AddDate(0, 1, -1)

	landingCount, err := s.repo.GetLandingCount(ctx, employeeID, monthStart.Format("2006-01-02"), monthEnd.Format("2006-01-02"))
	if err != nil {
		log.Error(logger.LogFlightSummaryAlertsError, "action", "landing_count", "error", err)
		return nil, err
	}

	if landingCount < domain.DefaultMinMonthlyLandings {
		alerts = append(alerts, domain.FlightAlert{
			Type:         domain.AlertTypeMinMonthlyLandings,
			Severity:     domain.AlertSeverityInfo,
			Message:      fmt.Sprintf("You have %d landings this month (minimum: %d)", landingCount, domain.DefaultMinMonthlyLandings),
			CurrentValue: landingCount,
			Threshold:    domain.DefaultMinMonthlyLandings,
		})
	}

	log.Info(logger.LogFlightSummaryAlertsOK, "alert_count", len(alerts))
	return alerts, nil
}
