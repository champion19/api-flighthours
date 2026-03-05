package interactor

import (
	"context"

	"github.com/champion19/api-flighthours/core/interactor/helpers"
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/input"
	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/middleware"
	"github.com/champion19/api-flighthours/platform/logger"
)

type TailNumberInteractor struct {
	service input.TailNumberService
	logger  logger.Logger
}

func NewTailNumberInteractor(service input.TailNumberService, log logger.Logger) *TailNumberInteractor {
	return &TailNumberInteractor{
		service: service,
		logger:  log,
	}
}

func (i *TailNumberInteractor) GetTailNumberByID(ctx context.Context, id string) (*domain.TailNumber, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogTailNumberGet, "registration_id", id)

	registration, err := i.service.GetTailNumberByID(ctx, id)
	if err != nil {
		log.Error(logger.LogTailNumberGetError, "registration_id", id, "error", err)
		return nil, err
	}

	log.Success(logger.LogTailNumberGetOK, registration.ToLogger())
	return registration, nil
}

func (i *TailNumberInteractor) GetTailNumberByPlate(ctx context.Context, plate string) (*domain.TailNumber, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogTailNumberGet, "tail_number", plate)

	registration, err := i.service.GetTailNumberByPlate(ctx, plate)
	if err != nil {
		log.Error(logger.LogTailNumberGetError, "tail_number", plate, "error", err)
		return nil, err
	}

	log.Success(logger.LogTailNumberGetOK, registration.ToLogger())
	return registration, nil
}

func (i *TailNumberInteractor) ListTailNumbers(ctx context.Context, filters map[string]interface{}) ([]domain.TailNumber, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogTailNumberList, "filters", filters)

	registrations, err := i.service.ListTailNumbers(ctx, filters)
	if err != nil {
		log.Error(logger.LogTailNumberListError, "error", err)
		return nil, err
	}

	log.Success(logger.LogTailNumberListOK, "count", len(registrations))
	return registrations, nil
}

func (i *TailNumberInteractor) CreateTailNumber(ctx context.Context, registration domain.TailNumber) (err error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogTailNumberCreate, registration.ToLogger())

	err = helpers.RunWithTx(ctx, i.service, log, logger.LogTailNumberCreateError,
		func(ctx context.Context, tx output.Tx) error {
			return i.service.CreateTailNumberTx(ctx, tx, registration)
		})
	if err != nil {
		return err
	}

	log.Success(logger.LogTailNumberCreateOK, registration.ToLogger())
	return nil
}

func (i *TailNumberInteractor) UpdateTailNumber(ctx context.Context, registration domain.TailNumber) (err error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogTailNumberUpdate, registration.ToLogger())

	err = helpers.RunWithTx(ctx, i.service, log, logger.LogTailNumberUpdateError,
		func(ctx context.Context, tx output.Tx) error {
			return i.service.UpdateTailNumberTx(ctx, tx, registration)
		})
	if err != nil {
		return err
	}

	log.Success(logger.LogTailNumberUpdateOK, registration.ToLogger())
	return nil
}
