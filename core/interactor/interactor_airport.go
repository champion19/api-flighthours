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

// ActivateAirport sets an airport's status to active
func (i *AirportInteractor) ActivateAirport(ctx context.Context, id string) (err error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogAirportActivate, "airport_id", id)

	_, err = i.service.GetAirportByID(ctx, id)
	if err != nil {
		log.Error(logger.LogAirportNotFound, "airport_id", id)
		return err
	}

	tx, err := i.service.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogAirportActivateError, "airport_id", id, "error", err)
		return err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(logger.LogAirportActivateError, "airport_id", id, "rollback_error", rbErr, "original_error", err)
			} else {
				log.Warn(logger.LogAirportActivateError, "airport_id", id, "rollback", "ok")
			}
		}
	}()

	if err = i.service.ActivateAirportTx(ctx, tx, id); err != nil {
		log.Error(logger.LogAirportActivateError, "airport_id", id, "error", err)
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Error(logger.LogAirportActivateError, "airport_id", id, "commit_error", err)
		return err
	}

	log.Success(logger.LogAirportActivateOK, "airport_id", id)
	return nil
}

func (i *AirportInteractor) DeactivateAirport(ctx context.Context, id string) (err error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogAirportDeactivate, "airport_id", id)

	_, err = i.service.GetAirportByID(ctx, id)
	if err != nil {
		log.Error(logger.LogAirportNotFound, "airport_id", id)
		return err
	}

	tx, err := i.service.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogAirportDeactivateError, "airport_id", id, "error", err)
		return err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(logger.LogAirportDeactivateError, "airport_id", id, "rollback_error", rbErr, "original_error", err)
			} else {
				log.Warn(logger.LogAirportDeactivateError, "airport_id", id, "rollback", "ok")
			}
		}
	}()

	if err = i.service.DeactivateAirportTx(ctx, tx, id); err != nil {
		log.Error(logger.LogAirportDeactivateError, "airport_id", id, "error", err)
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Error(logger.LogAirportDeactivateError, "airport_id", id, "commit_error", err)
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
