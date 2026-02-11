package interactor

import (
	"context"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/input"
	"github.com/champion19/api-flighthours/platform/logger"
)

// DailyLogbookDetailInteractor orchestrates daily logbook detail operations
// This is the CORE interactor for flight segment tracking
type DailyLogbookDetailInteractor struct {
	service        input.DailyLogbookDetailService
	logbookService input.DailyLogbookService // For ownership verification
}

// NewDailyLogbookDetailInteractor creates a new DailyLogbookDetailInteractor
func NewDailyLogbookDetailInteractor(
	service input.DailyLogbookDetailService,
	logbookService input.DailyLogbookService,
) *DailyLogbookDetailInteractor {
	return &DailyLogbookDetailInteractor{
		service:        service,
		logbookService: logbookService,
	}
}

// GetDailyLogbookDetailByID retrieves a detail by ID
func (i *DailyLogbookDetailInteractor) GetDailyLogbookDetailByID(ctx context.Context, traceID string, id string) (*domain.DailyLogbookDetail, error) {
	log.Info(logger.LogDailyLogbookDetailGet, "trace_id", traceID, "id", id)

	detail, err := i.service.GetDailyLogbookDetailByID(ctx, id)
	if err != nil {
		log.Error(logger.LogDailyLogbookDetailGetError, "trace_id", traceID, "error", err)
		return nil, err
	}

	if detail == nil {
		log.Warn(logger.LogDailyLogbookDetailNotFound, "trace_id", traceID, "id", id)
		return nil, domain.ErrFlightNotFound
	}

	log.Info(logger.LogDailyLogbookDetailGetOK, "trace_id", traceID, "id", id)
	return detail, nil
}

// ListDailyLogbookDetailsByLogbook lists all details for a logbook
func (i *DailyLogbookDetailInteractor) ListDailyLogbookDetailsByLogbook(ctx context.Context, traceID string, logbookID string) ([]domain.DailyLogbookDetail, error) {
	log.Info(logger.LogDailyLogbookDetailList, "trace_id", traceID, "logbook_id", logbookID)

	// Verify logbook exists
	logbook, err := i.logbookService.GetDailyLogbookByID(ctx, logbookID)
	if err != nil {
		log.Error(logger.LogDailyLogbookDetailListError, "trace_id", traceID, "error", err)
		return nil, err
	}
	if logbook == nil {
		log.Warn(logger.LogDailyLogbookNotFound, "trace_id", traceID, "logbook_id", logbookID)
		return nil, domain.ErrFlightInvalidLogbook
	}

	details, err := i.service.ListDailyLogbookDetailsByLogbook(ctx, logbookID)
	if err != nil {
		log.Error(logger.LogDailyLogbookDetailListError, "trace_id", traceID, "error", err)
		return nil, err
	}

	log.Info(logger.LogDailyLogbookDetailListOK, "trace_id", traceID, "count", len(details))
	return details, nil
}

// ListDailyLogbookDetailsByEmployee lists all flight details for an employee
func (i *DailyLogbookDetailInteractor) ListDailyLogbookDetailsByEmployee(ctx context.Context, traceID string, employeeID string) ([]domain.DailyLogbookDetail, error) {
	log.Info(logger.LogDailyLogbookDetailList, "trace_id", traceID, "employee_id", employeeID)

	details, err := i.service.ListDailyLogbookDetailsByEmployee(ctx, employeeID)
	if err != nil {
		log.Error(logger.LogDailyLogbookDetailListError, "trace_id", traceID, "error", err)
		return nil, err
	}

	log.Info(logger.LogDailyLogbookDetailListOK, "trace_id", traceID, "employee_id", employeeID, "count", len(details))
	return details, nil
}

// CreateDailyLogbookDetail creates a new detail
func (i *DailyLogbookDetailInteractor) CreateDailyLogbookDetail(ctx context.Context, traceID string, detail domain.DailyLogbookDetail) error {
	log.Info(logger.LogDailyLogbookDetailCreate, "trace_id", traceID, "data", detail.ToLogger())

	// Validate time sequence
	if err := i.service.ValidateTimeSequence(detail.OutTime, detail.TakeoffTime, detail.LandingTime, detail.InTime); err != nil {
		log.Error(logger.LogDailyLogbookDetailCreateError, "trace_id", traceID, "error", "invalid time sequence")
		return err
	}

	// Generate ID if not set
	if detail.ID == "" {
		detail.SetID()
	}

	err := i.service.CreateDailyLogbookDetail(ctx, detail)
	if err != nil {
		log.Error(logger.LogDailyLogbookDetailCreateError, "trace_id", traceID, "error", err)
		return err
	}

	log.Info(logger.LogDailyLogbookDetailCreateOK, "trace_id", traceID, "id", detail.ID)
	return nil
}

// UpdateDailyLogbookDetail updates an existing detail
func (i *DailyLogbookDetailInteractor) UpdateDailyLogbookDetail(ctx context.Context, traceID string, detail domain.DailyLogbookDetail) error {
	log.Info(logger.LogDailyLogbookDetailUpdate, "trace_id", traceID, "data", detail.ToLogger())

	// Verify detail exists
	existing, err := i.service.GetDailyLogbookDetailByID(ctx, detail.ID)
	if err != nil {
		log.Error(logger.LogDailyLogbookDetailUpdateError, "trace_id", traceID, "error", err)
		return err
	}
	if existing == nil {
		log.Warn(logger.LogDailyLogbookDetailNotFound, "trace_id", traceID, "id", detail.ID)
		return domain.ErrFlightNotFound
	}

	// Validate time sequence
	if err := i.service.ValidateTimeSequence(detail.OutTime, detail.TakeoffTime, detail.LandingTime, detail.InTime); err != nil {
		log.Error(logger.LogDailyLogbookDetailUpdateError, "trace_id", traceID, "error", "invalid time sequence")
		return err
	}

	// Preserve the daily_logbook_id from existing record (cannot change parent)
	detail.DailyLogbookID = existing.DailyLogbookID

	err = i.service.UpdateDailyLogbookDetail(ctx, detail)
	if err != nil {
		log.Error(logger.LogDailyLogbookDetailUpdateError, "trace_id", traceID, "error", err)
		return err
	}

	log.Info(logger.LogDailyLogbookDetailUpdateOK, "trace_id", traceID, "id", detail.ID)
	return nil
}

// VerifyLogbookOwnership verifies that a logbook belongs to the specified employee
func (i *DailyLogbookDetailInteractor) VerifyLogbookOwnership(ctx context.Context, logbookID string, employeeID string) error {
	logbook, err := i.logbookService.GetDailyLogbookByID(ctx, logbookID)
	if err != nil {
		return err
	}
	if logbook == nil {
		return domain.ErrFlightInvalidLogbook
	}
	if logbook.EmployeeID != employeeID {
		return domain.ErrFlightUnauthorized
	}
	return nil
}

// GetLogbookOwner returns the employee ID that owns a logbook
func (i *DailyLogbookDetailInteractor) GetLogbookOwner(ctx context.Context, logbookID string) (string, error) {
	logbook, err := i.logbookService.GetDailyLogbookByID(ctx, logbookID)
	if err != nil {
		return "", err
	}
	if logbook == nil {
		return "", domain.ErrFlightInvalidLogbook
	}
	return logbook.EmployeeID, nil
}

// GetDetailLogbookOwner returns the employee ID that owns the logbook of a detail
func (i *DailyLogbookDetailInteractor) GetDetailLogbookOwner(ctx context.Context, detailID string) (string, error) {
	detail, err := i.service.GetDailyLogbookDetailByID(ctx, detailID)
	if err != nil {
		return "", err
	}
	if detail == nil {
		return "", domain.ErrFlightNotFound
	}
	return i.GetLogbookOwner(ctx, detail.DailyLogbookID)
}
