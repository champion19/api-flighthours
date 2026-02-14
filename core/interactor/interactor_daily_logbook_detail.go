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

// ListDailyLogbookDetailsByLogbook lists all details for a logbook, verifying ownership
func (i *DailyLogbookDetailInteractor) ListDailyLogbookDetailsByLogbook(ctx context.Context, traceID string, logbookID string, employeeID string) ([]domain.DailyLogbookDetail, error) {
	log.Info(logger.LogDailyLogbookDetailList, "trace_id", traceID, "logbook_id", logbookID)

	// Verify logbook exists and ownership
	logbook, err := i.logbookService.GetDailyLogbookByID(ctx, logbookID)
	if err != nil {
		log.Error(logger.LogDailyLogbookDetailListError, "trace_id", traceID, "error", err)
		return nil, err
	}
	if logbook == nil {
		log.Warn(logger.LogDailyLogbookNotFound, "trace_id", traceID, "logbook_id", logbookID)
		return nil, domain.ErrFlightInvalidLogbook
	}
	if logbook.EmployeeID != employeeID {
		log.Warn(logger.LogDailyLogbookDetailListError, "trace_id", traceID, "error", "unauthorized")
		return nil, domain.ErrFlightUnauthorized
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

// CreateDailyLogbookDetail creates a new detail, verifying logbook ownership
func (i *DailyLogbookDetailInteractor) CreateDailyLogbookDetail(ctx context.Context, traceID string, detail domain.DailyLogbookDetail, employeeID string) (err error) {
	log.Info(logger.LogDailyLogbookDetailCreate, "trace_id", traceID, "data", detail.ToLogger())

	// Verify logbook ownership and active status
	if err = i.VerifyLogbookActiveAndOwned(ctx, detail.DailyLogbookID, employeeID); err != nil {
		log.Warn(logger.LogDailyLogbookDetailCreateError, "trace_id", traceID, "error", err)
		return err
	}

	// Validate time sequence only if all 4 time fields are provided
	if detail.OutTime != nil && detail.TakeoffTime != nil && detail.LandingTime != nil && detail.InTime != nil {
		if err = i.service.ValidateTimeSequence(*detail.OutTime, *detail.TakeoffTime, *detail.LandingTime, *detail.InTime); err != nil {
			log.Error(logger.LogDailyLogbookDetailCreateError, "trace_id", traceID, "error", "invalid time sequence")
			return err
		}
	}

	// Generate ID if not set
	if detail.ID == "" {
		detail.SetID()
	}

	// Check for duplicate flight
	if detail.EmployeeLogbookID != nil {
		exists, err := i.service.ExistsByUniqueKey(ctx, *detail.EmployeeLogbookID, detail.FlightRealDate, detail.FlightNumber, detail.LicensePlateID)
		if err != nil {
			log.Error(logger.LogDailyLogbookDetailCreateError, "trace_id", traceID, "error", err)
			return err
		}
		if exists {
			log.Warn(logger.LogDailyLogbookDetailDuplicate, "trace_id", traceID,
				"employee_logbook_id", *detail.EmployeeLogbookID,
				"flight_real_date", detail.FlightRealDate,
				"flight_number", detail.FlightNumber,
				"license_plate_id", detail.LicensePlateID)
			return domain.ErrFlightDuplicate
		}
	}

	tx, err := i.service.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogDailyLogbookDetailCreateError, "trace_id", traceID, "error", err)
		return err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(logger.LogDailyLogbookDetailCreateError, "trace_id", traceID, "rollback_error", rbErr)
			} else {
				log.Warn(logger.LogDailyLogbookDetailCreateError, "trace_id", traceID, "rollback", "ok")
			}
		}
	}()

	if err = i.service.CreateDailyLogbookDetailTx(ctx, tx, detail); err != nil {
		log.Error(logger.LogDailyLogbookDetailCreateError, "trace_id", traceID, "error", err)
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Error(logger.LogDailyLogbookDetailCreateError, "trace_id", traceID, "commit_error", err)
		return err
	}

	log.Info(logger.LogDailyLogbookDetailCreateOK, "trace_id", traceID, "id", detail.ID)
	return nil
}

// UpdateDailyLogbookDetail updates an existing detail, verifying ownership
func (i *DailyLogbookDetailInteractor) UpdateDailyLogbookDetail(ctx context.Context, traceID string, detail domain.DailyLogbookDetail, employeeID string) (err error) {
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

	// Verify ownership and active status via detail's logbook
	if err = i.VerifyLogbookActiveAndOwned(ctx, existing.DailyLogbookID, employeeID); err != nil {
		log.Warn(logger.LogDailyLogbookDetailUpdateError, "trace_id", traceID, "error", err)
		return err
	}

	// Validate time sequence only if all 4 time fields are provided
	if detail.OutTime != nil && detail.TakeoffTime != nil && detail.LandingTime != nil && detail.InTime != nil {
		if err = i.service.ValidateTimeSequence(*detail.OutTime, *detail.TakeoffTime, *detail.LandingTime, *detail.InTime); err != nil {
			log.Error(logger.LogDailyLogbookDetailUpdateError, "trace_id", traceID, "error", "invalid time sequence")
			return err
		}
	}

	// Preserve the daily_logbook_id from existing record (cannot change parent)
	detail.DailyLogbookID = existing.DailyLogbookID

	tx, err := i.service.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogDailyLogbookDetailUpdateError, "trace_id", traceID, "error", err)
		return err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(logger.LogDailyLogbookDetailUpdateError, "trace_id", traceID, "rollback_error", rbErr)
			} else {
				log.Warn(logger.LogDailyLogbookDetailUpdateError, "trace_id", traceID, "rollback", "ok")
			}
		}
	}()

	if err = i.service.UpdateDailyLogbookDetailTx(ctx, tx, detail); err != nil {
		log.Error(logger.LogDailyLogbookDetailUpdateError, "trace_id", traceID, "error", err)
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Error(logger.LogDailyLogbookDetailUpdateError, "trace_id", traceID, "commit_error", err)
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

// VerifyLogbookActiveAndOwned verifies ownership AND that the logbook is active (open for modifications)
func (i *DailyLogbookDetailInteractor) VerifyLogbookActiveAndOwned(ctx context.Context, logbookID string, employeeID string) error {
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
	if !logbook.Status {
		return domain.ErrDailyLogbookInactive
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
