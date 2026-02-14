package interactor

import (
	"context"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/input"
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

	tx, err := i.service.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogDailyLogbookCreateError, "error", err)
		return err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(logger.LogDailyLogbookCreateError, "rollback_error", rbErr, "original_error", err)
			} else {
				log.Warn(logger.LogDailyLogbookCreateError, "rollback", "ok")
			}
		}
	}()

	if err = i.service.CreateDailyLogbookTx(ctx, tx, logbook); err != nil {
		log.Error(logger.LogDailyLogbookCreateError, "error", err)
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Error(logger.LogDailyLogbookCreateError, "commit_error", err)
		return err
	}

	log.Success(logger.LogDailyLogbookCreateOK, logbook.ToLogger())
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

	tx, err := i.service.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogDailyLogbookActivateError, "error", err)
		return err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(logger.LogDailyLogbookActivateError, "rollback_error", rbErr, "original_error", err)
			} else {
				log.Warn(logger.LogDailyLogbookActivateError, "rollback", "ok")
			}
		}
	}()

	if err = i.service.ActivateDailyLogbookTx(ctx, tx, id); err != nil {
		log.Error(logger.LogDailyLogbookActivateError, "logbook_id", id, "error", err)
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Error(logger.LogDailyLogbookActivateError, "commit_error", err)
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

	tx, err := i.service.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogDailyLogbookDeactivateError, "error", err)
		return err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(logger.LogDailyLogbookDeactivateError, "rollback_error", rbErr, "original_error", err)
			} else {
				log.Warn(logger.LogDailyLogbookDeactivateError, "rollback", "ok")
			}
		}
	}()

	if err = i.service.DeactivateDailyLogbookTx(ctx, tx, id); err != nil {
		log.Error(logger.LogDailyLogbookDeactivateError, "logbook_id", id, "error", err)
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Error(logger.LogDailyLogbookDeactivateError, "commit_error", err)
		return err
	}

	log.Success(logger.LogDailyLogbookDeactivateOK, "logbook_id", id)
	return nil
}
