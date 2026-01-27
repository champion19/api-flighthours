package interactor

import (
	"context"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/input"
	"github.com/champion19/api-flighthours/middleware"
	"github.com/champion19/api-flighthours/platform/logger"
)

// AirlineInteractor orchestrates airline business operations
type AirlineInteractor struct {
	service input.AirlineService
	logger  logger.Logger
}

// NewAirlineInteractor creates a new airline interactor instance
func NewAirlineInteractor(service input.AirlineService, log logger.Logger) *AirlineInteractor {
	return &AirlineInteractor{
		service: service,
		logger:  log,
	}
}

// GetAirlineByID retrieves an airline by its ID
func (i *AirlineInteractor) GetAirlineByID(ctx context.Context, id string) (*domain.Airline, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

	log.Info(logger.LogAirlineGet, "airline_id", id)

	airline, err := i.service.GetAirlineByID(ctx, id)
	if err != nil {
		if err == domain.ErrAirlineNotFound {
			log.Warn(logger.LogAirlineNotFound, "airline_id", id)
			return nil, domain.ErrAirlineNotFound
		}
		log.Error(logger.LogAirlineGetError, "airline_id", id, "error", err)
		return nil, err
	}

	log.Success(logger.LogAirlineGetOK, "airline_id", id, "airline_name", airline.AirlineName)
	return airline, nil
}

func (i *AirlineInteractor) ListAirlines(ctx context.Context, filters map[string]interface{}) ([]domain.Airline, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

	log.Info(logger.LogAirlineList, "filters", filters)

	airlines, err := i.service.ListAirlines(ctx, filters)
	if err != nil {
		log.Error(logger.LogAirlineListError, "error", err)
		return nil, err
	}

	log.Success(logger.LogAirlineListOK, "count", len(airlines))
	return airlines, nil
}
