package interactor

import (
	"context"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/input"
	"github.com/champion19/api-flighthours/middleware"
	"github.com/champion19/api-flighthours/platform/logger"
)


type AirlineInteractor struct {
	service input.AirlineService
}


func NewAirlineInteractor(service input.AirlineService) *AirlineInteractor {
	return &AirlineInteractor{
		service: service,
	}
}


func (i *AirlineInteractor) GetAirlineByID(ctx context.Context, id string) (*domain.Airline, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

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
	log := log.WithTraceID(traceID)

	log.Info(logger.LogAirlineList, "filters", filters)

	airlines, err := i.service.ListAirlines(ctx, filters)
	if err != nil {
		log.Error(logger.LogAirlineListError, "error", err)
		return nil, err
	}

	log.Success(logger.LogAirlineListOK, "count", len(airlines))
	return airlines, nil
}
