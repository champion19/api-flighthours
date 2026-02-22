package interactor

import (
	"context"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/input"
	"github.com/champion19/api-flighthours/platform/logger"
)

// FlightSummaryInteractor orchestrates flight summary operations.
// No business logic here — it delegates to the service layer.
type FlightSummaryInteractor struct {
	service input.FlightSummaryService
}

func NewFlightSummaryInteractor(service input.FlightSummaryService) *FlightSummaryInteractor {
	return &FlightSummaryInteractor{service: service}
}

// GetFlightHoursSummary resolves the period dates (via service) and fetches the summary
func (i *FlightSummaryInteractor) GetFlightHoursSummary(ctx context.Context, traceID, employeeID, period, startDate, endDate, referenceDate string) (*domain.FlightHoursSummary, error) {
	log.Info(logger.LogFlightSummaryGet, "trace_id", traceID, "employee_id", employeeID, "period", period)

	// If period is not "custom", compute dates via the service (business logic)
	if period != domain.PeriodCustom {
		var err error
		startDate, endDate, err = i.service.CalculatePeriodDates(period, referenceDate)
		if err != nil {
			log.Error(logger.LogFlightSummaryGetError, "trace_id", traceID, "error", err)
			return nil, err
		}
	}

	summary, err := i.service.GetFlightHoursSummary(ctx, employeeID, startDate, endDate)
	if err != nil {
		log.Error(logger.LogFlightSummaryGetError, "trace_id", traceID, "error", err)
		return nil, err
	}

	summary.Period = period
	return summary, nil
}

// GetFlightAlerts delegates alert evaluation to the service
func (i *FlightSummaryInteractor) GetFlightAlerts(ctx context.Context, traceID, employeeID string) ([]domain.FlightAlert, error) {
	log.Info(logger.LogFlightSummaryAlerts, "trace_id", traceID, "employee_id", employeeID)

	alerts, err := i.service.BuildFlightAlerts(ctx, employeeID)
	if err != nil {
		log.Error(logger.LogFlightSummaryAlertsError, "trace_id", traceID, "error", err)
		return nil, err
	}

	return alerts, nil
}

// GetRecentFlights delegates to the service to get the last N flights
func (i *FlightSummaryInteractor) GetRecentFlights(ctx context.Context, traceID, employeeID string, limit int) ([]domain.DailyLogbookDetail, error) {
	log.Info(logger.LogFlightSummaryRecentFlights, "trace_id", traceID, "employee_id", employeeID, "limit", limit)

	flights, err := i.service.GetRecentFlights(ctx, employeeID, limit)
	if err != nil {
		log.Error(logger.LogFlightSummaryGetError, "trace_id", traceID, "error", err)
		return nil, err
	}

	return flights, nil
}
