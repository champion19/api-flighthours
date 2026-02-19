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

// DailyLogbookInteractor orchestrates daily logbook operations
type DailyLogbookInteractor struct {
	service input.DailyLogbookService
}

// NewDailyLogbookInteractor creates a new daily logbook interactor
func NewDailyLogbookInteractor(service input.DailyLogbookService) *DailyLogbookInteractor {
	return &DailyLogbookInteractor{
		service: service,
	}
}

// GetDailyLogbookByID retrieves a daily logbook by its ID
func (i *DailyLogbookInteractor) GetDailyLogbookByID(ctx context.Context, id string, employeeID string) (*domain.DailyLogbook, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogDailyLogbookGet, "logbook_id", id)

	logbook, err := i.service.GetDailyLogbookByID(ctx, id)
	if err != nil {
		log.Error(logger.LogDailyLogbookGetError, "logbook_id", id, "error", err)
		return nil, err
	}

	// Verify ownership
	if logbook.EmployeeID != employeeID {
		log.Warn(logger.LogDailyLogbookGetError, "error", "unauthorized", "logbook_id", id)
		return nil, domain.ErrFlightUnauthorized
	}

	log.Success(logger.LogDailyLogbookGetOK, logbook.ToLogger())
	return logbook, nil
}

// ListDailyLogbooksByEmployee retrieves all daily logbooks for the authenticated employee
func (i *DailyLogbookInteractor) ListDailyLogbooksByEmployee(ctx context.Context, employeeID string, filters map[string]interface{}) ([]domain.DailyLogbook, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogDailyLogbookList, "employee_id", employeeID, "filters", filters)

	logbooks, err := i.service.ListDailyLogbooksByEmployee(ctx, employeeID, filters)
	if err != nil {
		log.Error(logger.LogDailyLogbookListError, "error", err)
		return nil, err
	}

	log.Success(logger.LogDailyLogbookListOK, "count", len(logbooks))
	return logbooks, nil
}

// CreateDailyLogbook creates a new daily logbook for the authenticated employee
func (i *DailyLogbookInteractor) CreateDailyLogbook(ctx context.Context, logbook domain.DailyLogbook) (err error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogDailyLogbookCreate, "employee_id", logbook.EmployeeID)

	err = helpers.RunWithTx(ctx, i.service, log, logger.LogDailyLogbookCreateError,
		func(ctx context.Context, tx output.Tx) error {
			return i.service.CreateDailyLogbookTx(ctx, tx, logbook)
		})
	if err != nil {
		return err
	}

	log.Success(logger.LogDailyLogbookCreateOK, logbook.ToLogger())
	return nil
}

// UpdateDailyLogbook updates an existing daily logbook
func (i *DailyLogbookInteractor) UpdateDailyLogbook(ctx context.Context, logbook domain.DailyLogbook, employeeID string) (err error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogDailyLogbookUpdate, "logbook_id", logbook.ID)

	// Verify logbook exists and belongs to the authenticated employee
	existing, err := i.service.GetDailyLogbookByID(ctx, logbook.ID)
	if err != nil {
		log.Error(logger.LogDailyLogbookNotFound, "logbook_id", logbook.ID)
		return err
	}

	if existing.EmployeeID != employeeID {
		log.Warn(logger.LogDailyLogbookUpdateError, "error", "unauthorized", "logbook_id", logbook.ID)
		return domain.ErrFlightUnauthorized
	}

	err = helpers.RunWithTx(ctx, i.service, log, logger.LogDailyLogbookUpdateError,
		func(ctx context.Context, tx output.Tx) error {
			return i.service.UpdateDailyLogbookTx(ctx, tx, logbook)
		})
	if err != nil {
		return err
	}

	log.Success(logger.LogDailyLogbookUpdateOK, logbook.ToLogger())
	return nil
}

// DeleteDailyLogbook removes a daily logbook
func (i *DailyLogbookInteractor) DeleteDailyLogbook(ctx context.Context, id string) (err error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogDailyLogbookDelete, "logbook_id", id)

	err = helpers.RunWithTx(ctx, i.service, log, logger.LogDailyLogbookDeleteError,
		func(ctx context.Context, tx output.Tx) error {
			return i.service.DeleteDailyLogbookTx(ctx, tx, id)
		})
	if err != nil {
		return err
	}

	log.Success(logger.LogDailyLogbookDeleteOK, "logbook_id", id)
	return nil
}

// ActivateDailyLogbook sets a daily logbook's status to active
func (i *DailyLogbookInteractor) ActivateDailyLogbook(ctx context.Context, id string) (err error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogDailyLogbookActivate, "logbook_id", id)

	// Verify logbook exists
	_, err = i.service.GetDailyLogbookByID(ctx, id)
	if err != nil {
		log.Error(logger.LogDailyLogbookNotFound, "logbook_id", id)
		return err
	}

	err = helpers.RunWithTx(ctx, i.service, log, logger.LogDailyLogbookActivateError,
		func(ctx context.Context, tx output.Tx) error {
			return i.service.ActivateDailyLogbookTx(ctx, tx, id)
		})
	if err != nil {
		return err
	}

	log.Success(logger.LogDailyLogbookActivateOK, "logbook_id", id)
	return nil
}

// DeactivateDailyLogbook sets a daily logbook's status to inactive
func (i *DailyLogbookInteractor) DeactivateDailyLogbook(ctx context.Context, id string) (err error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogDailyLogbookDeactivate, "logbook_id", id)

	// Verify logbook exists
	_, err = i.service.GetDailyLogbookByID(ctx, id)
	if err != nil {
		log.Error(logger.LogDailyLogbookNotFound, "logbook_id", id)
		return err
	}

	err = helpers.RunWithTx(ctx, i.service, log, logger.LogDailyLogbookDeactivateError,
		func(ctx context.Context, tx output.Tx) error {
			return i.service.DeactivateDailyLogbookTx(ctx, tx, id)
		})
	if err != nil {
		return err
	}

	log.Success(logger.LogDailyLogbookDeactivateOK, "logbook_id", id)
	return nil
}
