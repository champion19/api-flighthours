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

type LicensePlateInteractor struct {
	service input.LicensePlateService
	logger  logger.Logger
}

func NewLicensePlateInteractor(service input.LicensePlateService, log logger.Logger) *LicensePlateInteractor {
	return &LicensePlateInteractor{
		service: service,
		logger:  log,
	}
}

func (i *LicensePlateInteractor) GetLicensePlateByID(ctx context.Context, id string) (*domain.LicensePlate, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogLicensePlateGet, "registration_id", id)

	registration, err := i.service.GetLicensePlateByID(ctx, id)
	if err != nil {
		log.Error(logger.LogLicensePlateGetError, "registration_id", id, "error", err)
		return nil, err
	}

	log.Success(logger.LogLicensePlateGetOK, registration.ToLogger())
	return registration, nil
}

func (i *LicensePlateInteractor) GetLicensePlateByPlate(ctx context.Context, plate string) (*domain.LicensePlate, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogLicensePlateGet, "license_plate", plate)

	registration, err := i.service.GetLicensePlateByPlate(ctx, plate)
	if err != nil {
		log.Error(logger.LogLicensePlateGetError, "license_plate", plate, "error", err)
		return nil, err
	}

	log.Success(logger.LogLicensePlateGetOK, registration.ToLogger())
	return registration, nil
}

func (i *LicensePlateInteractor) ListLicensePlates(ctx context.Context, filters map[string]interface{}) ([]domain.LicensePlate, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogLicensePlateList, "filters", filters)

	registrations, err := i.service.ListLicensePlates(ctx, filters)
	if err != nil {
		log.Error(logger.LogLicensePlateListError, "error", err)
		return nil, err
	}

	log.Success(logger.LogLicensePlateListOK, "count", len(registrations))
	return registrations, nil
}

func (i *LicensePlateInteractor) CreateLicensePlate(ctx context.Context, registration domain.LicensePlate) (err error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogLicensePlateCreate, registration.ToLogger())

	err = helpers.RunWithTx(ctx, i.service, log, logger.LogLicensePlateCreateError,
		func(ctx context.Context, tx output.Tx) error {
			return i.service.CreateLicensePlateTx(ctx, tx, registration)
		})
	if err != nil {
		return err
	}

	log.Success(logger.LogLicensePlateCreateOK, registration.ToLogger())
	return nil
}

func (i *LicensePlateInteractor) UpdateLicensePlate(ctx context.Context, registration domain.LicensePlate) (err error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogLicensePlateUpdate, registration.ToLogger())

	err = helpers.RunWithTx(ctx, i.service, log, logger.LogLicensePlateUpdateError,
		func(ctx context.Context, tx output.Tx) error {
			return i.service.UpdateLicensePlateTx(ctx, tx, registration)
		})
	if err != nil {
		return err
	}

	log.Success(logger.LogLicensePlateUpdateOK, registration.ToLogger())
	return nil
}
