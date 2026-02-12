package interactor

import (
	"context"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/input"
	"github.com/champion19/api-flighthours/middleware"
	"github.com/champion19/api-flighthours/platform/logger"
)

// AircraftModelInteractor orchestrates aircraft model operations
type AircraftModelInteractor struct {
	service input.AircraftModelService
}

// NewAircraftModelInteractor creates a new aircraft model interactor
func NewAircraftModelInteractor(service input.AircraftModelService) *AircraftModelInteractor {
	return &AircraftModelInteractor{
		service: service,
	}
}

// GetAircraftModelByID retrieves an aircraft model by its ID
func (i *AircraftModelInteractor) GetAircraftModelByID(ctx context.Context, id string) (*domain.AircraftModel, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogAircraftModelGet, "aircraft_model_id", id)

	model, err := i.service.GetAircraftModelByID(ctx, id)
	if err != nil {
		log.Error(logger.LogAircraftModelGetError, "aircraft_model_id", id, "error", err)
		return nil, err
	}

	log.Success(logger.LogAircraftModelGetOK, model.ToLogger())
	return model, nil
}

// ListAircraftModels retrieves all aircraft models with optional filters
func (i *AircraftModelInteractor) ListAircraftModels(ctx context.Context, filters map[string]interface{}) ([]domain.AircraftModel, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogAircraftModelList, "filters", filters)

	models, err := i.service.ListAircraftModels(ctx, filters)
	if err != nil {
		log.Error(logger.LogAircraftModelListError, "error", err)
		return nil, err
	}

	log.Success(logger.LogAircraftModelListOK, "count", len(models))
	return models, nil
}

// GetAircraftModelsByFamily retrieves all aircraft models for a specific family (HU30)
func (i *AircraftModelInteractor) GetAircraftModelsByFamily(ctx context.Context, family string) ([]domain.AircraftModel, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogAircraftModelList, "family", family)

	models, err := i.service.GetAircraftModelsByFamily(ctx, family)
	if err != nil {
		log.Error(logger.LogAircraftModelListError, "family", family, "error", err)
		return nil, err
	}

	log.Success(logger.LogAircraftModelListOK, "family", family, "count", len(models))
	return models, nil
}

// ActivateAircraftModel sets an aircraft model's status to active (HU42)
func (i *AircraftModelInteractor) ActivateAircraftModel(ctx context.Context, id string) (err error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogAircraftModelActivate, "aircraft_model_id", id)

	// Verify aircraft model exists
	_, err = i.service.GetAircraftModelByID(ctx, id)
	if err != nil {
		log.Error(logger.LogAircraftModelNotFound, "aircraft_model_id", id)
		return err
	}

	tx, err := i.service.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogAircraftModelActivateError, "aircraft_model_id", id, "error", err)
		return err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(logger.LogAircraftModelActivateError, "aircraft_model_id", id, "rollback_error", rbErr, "original_error", err)
			} else {
				log.Warn(logger.LogAircraftModelActivateError, "aircraft_model_id", id, "rollback", "ok")
			}
		}
	}()

	if err = i.service.ActivateAircraftModelTx(ctx, tx, id); err != nil {
		log.Error(logger.LogAircraftModelActivateError, "aircraft_model_id", id, "error", err)
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Error(logger.LogAircraftModelActivateError, "aircraft_model_id", id, "commit_error", err)
		return err
	}

	log.Success(logger.LogAircraftModelActivateOK, "aircraft_model_id", id)
	return nil
}

// DeactivateAircraftModel sets an aircraft model's status to inactive (HU41)
func (i *AircraftModelInteractor) DeactivateAircraftModel(ctx context.Context, id string) (err error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogAircraftModelDeactivate, "aircraft_model_id", id)

	// Verify aircraft model exists
	_, err = i.service.GetAircraftModelByID(ctx, id)
	if err != nil {
		log.Error(logger.LogAircraftModelNotFound, "aircraft_model_id", id)
		return err
	}

	tx, err := i.service.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogAircraftModelDeactivateError, "aircraft_model_id", id, "error", err)
		return err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(logger.LogAircraftModelDeactivateError, "aircraft_model_id", id, "rollback_error", rbErr, "original_error", err)
			} else {
				log.Warn(logger.LogAircraftModelDeactivateError, "aircraft_model_id", id, "rollback", "ok")
			}
		}
	}()

	if err = i.service.DeactivateAircraftModelTx(ctx, tx, id); err != nil {
		log.Error(logger.LogAircraftModelDeactivateError, "aircraft_model_id", id, "error", err)
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Error(logger.LogAircraftModelDeactivateError, "aircraft_model_id", id, "commit_error", err)
		return err
	}

	log.Success(logger.LogAircraftModelDeactivateOK, "aircraft_model_id", id)
	return nil
}
