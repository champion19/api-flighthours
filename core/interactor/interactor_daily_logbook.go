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
	logger  logger.Logger
}

// NewDailyLogbookInteractor creates a new daily logbook interactor
func NewDailyLogbookInteractor(service input.DailyLogbookService, log logger.Logger) *DailyLogbookInteractor {
	return &DailyLogbookInteractor{
		service: service,
		logger:  log,
	}
}

// GetDailyLogbookByID retrieves a daily logbook by its ID
func (i *DailyLogbookInteractor) GetDailyLogbookByID(ctx context.Context, id string) (*domain.DailyLogbook, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

	log.Info(logger.LogDailyLogbookGet, "logbook_id", id)

	logbook, err := i.service.GetDailyLogbookByID(ctx, id)
	if err != nil {
		log.Error(logger.LogDailyLogbookGetError, "logbook_id", id, "error", err)
		return nil, err
	}

	log.Success(logger.LogDailyLogbookGetOK, logbook.ToLogger())
	return logbook, nil
}

// ListDailyLogbooksByEmployee retrieves all daily logbooks for the authenticated employee
func (i *DailyLogbookInteractor) ListDailyLogbooksByEmployee(ctx context.Context, employeeID string, filters map[string]interface{}) ([]domain.DailyLogbook, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

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
func (i *DailyLogbookInteractor) CreateDailyLogbook(ctx context.Context, logbook domain.DailyLogbook) error {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

	log.Info(logger.LogDailyLogbookCreate, "employee_id", logbook.EmployeeID)

	if err := i.service.CreateDailyLogbook(ctx, logbook); err != nil {
		log.Error(logger.LogDailyLogbookCreateError, "error", err)
		return err
	}

	log.Success(logger.LogDailyLogbookCreateOK, logbook.ToLogger())
	return nil
}


