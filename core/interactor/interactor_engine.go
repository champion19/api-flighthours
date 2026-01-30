package interactor

import (
	"context"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/input"
	"github.com/champion19/api-flighthours/middleware"
	"github.com/champion19/api-flighthours/platform/logger"
)


type EngineInteractor struct {
	service input.EngineService
}


func NewEngineInteractor(service input.EngineService) *EngineInteractor {
	return &EngineInteractor{
		service: service,
	}
}

func (i *EngineInteractor) GetEngineByID(ctx context.Context, id string) (*domain.Engine, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogEngineGet, "engine_id", id)

	engine, err := i.service.GetEngineByID(ctx, id)
	if err != nil {
		if err == domain.ErrEngineNotFound {
			log.Warn(logger.LogEngineNotFound, "engine_id", id)
			return nil, domain.ErrEngineNotFound
		}
		log.Error(logger.LogEngineGetError, "engine_id", id, "error", err)
		return nil, err
	}

	log.Success(logger.LogEngineGetOK, "engine_id", id, "engine_name", engine.Name)
	return engine, nil
}

func (i *EngineInteractor) ListEngines(ctx context.Context) ([]domain.Engine, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogEngineList)

	engines, err := i.service.ListEngines(ctx)
	if err != nil {
		log.Error(logger.LogEngineListError, "error", err)
		return nil, err
	}

	log.Success(logger.LogEngineListOK, "count", len(engines))
	return engines, nil
}
