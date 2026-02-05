package interactor

import (
	"context"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/input"
	"github.com/champion19/api-flighthours/middleware"
	"github.com/champion19/api-flighthours/platform/logger"
)


// AirportInteractor orchestrates airport operations
type AirportInteractor struct {
	service input.AirportService
}

// NewAirportInteractor creates a new airport interactor
func NewAirportInteractor(service input.AirportService) *AirportInteractor {
	return &AirportInteractor{
		service: service,
	}
}

// GetAirportByID retrieves an airport by its ID
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

// DeactivateAirport sets an airport's status to inactive
func (i *AirportInteractor) DeactivateAirport(ctx context.Context, id string) error {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogAirportDeactivate, "airport_id", id)

	// Verify airport exists
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

// ListAirports retrieves all airports with optional filters
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

// GetAirportsByCity retrieves all airports for a specific city (HU13 - Virtual Entity pattern)
func (i *AirportInteractor) GetAirportsByCity(ctx context.Context, city string) ([]domain.Airport, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogAirportList, "city", city)

	airports, err := i.service.GetAirportsByCity(ctx, city)
	if err != nil {
		log.Error(logger.LogAirportListError, "city", city, "error", err)
		return nil, err
	}

	log.Success(logger.LogAirportListOK, "city", city, "count", len(airports))
	return airports, nil
}

// GetAirportsByCountry retrieves all airports for a specific country (HU38 - Virtual Entity pattern)
func (i *AirportInteractor) GetAirportsByCountry(ctx context.Context, country string) ([]domain.Airport, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogAirportList, "country", country)

	airports, err := i.service.GetAirportsByCountry(ctx, country)
	if err != nil {
		log.Error(logger.LogAirportListError, "country", country, "error", err)
		return nil, err
	}

	log.Success(logger.LogAirportListOK, "country", country, "count", len(airports))
	return airports, nil
}

// GetAirportsByType retrieves all airports for a specific airport type (HU46 - Virtual Entity pattern)
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
