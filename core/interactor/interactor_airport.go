package interactor

import (
	"context"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/input"
	"github.com/champion19/api-flighthours/middleware"
	"github.com/champion19/api-flighthours/platform/logger"
)

type AirportInteractor struct {
	service input.AirportService
}

func NewAirportInteractor(service input.AirportService) *AirportInteractor {
	return &AirportInteractor{
		service: service,
	}
}

func (i *AirportInteractor) GetAirportByID(ctx context.Context, id string) (*domain.Airport, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogAirportGet, "airport_id", id)

	airport, err := i.service.GetAirportByID(ctx, id)
	if err != nil {
		log.Error(logger.LogAirportGetError, "airport_id", id, "error", err)
		return nil, err
	}

	log.Success(logger.LogAirportGetOK, airport.ToLogger())
	return airport, nil
}

func (i *AirportInteractor) DeactivateAirport(ctx context.Context, id string) error {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogAirportDeactivate, "airport_id", id)

	_, err := i.service.GetAirportByID(ctx, id)
	if err != nil {
		log.Error(logger.LogAirportNotFound, "airport_id", id)
		return err
	}

	if err = i.service.DeactivateAirport(ctx, id); err != nil {
		log.Error(logger.LogAirportDeactivateError, "airport_id", id, "error", err)
		return err
	}

	log.Success(logger.LogAirportDeactivateOK, "airport_id", id)
	return nil
}

func (i *AirportInteractor) ListAirports(ctx context.Context, filters map[string]interface{}) ([]domain.Airport, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogAirportList, "filters", filters)

	airports, err := i.service.ListAirports(ctx, filters)
	if err != nil {
		log.Error(logger.LogAirportListError, "error", err)
		return nil, err
	}

	log.Success(logger.LogAirportListOK, "count", len(airports))
	return airports, nil
}

func (i *AirportInteractor) GetAirportsByType(ctx context.Context, airportType string) ([]domain.Airport, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogAirportTypeGet, "airport_type", airportType)

	airports, err := i.service.GetAirportsByType(ctx, airportType)
	if err != nil {
		log.Error(logger.LogAirportTypeGetError, "airport_type", airportType, "error", err)
		return nil, err
	}

	log.Success(logger.LogAirportTypeGetOK, "airport_type", airportType, "count", len(airports))
	return airports, nil
}
