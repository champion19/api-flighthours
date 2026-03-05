package interactor

import (
	"context"

	"github.com/champion19/api-flighthours/core/interactor/helpers"
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/input"
	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/logger"
)

// DailyLogbookDetailInteractor orchestrates daily logbook detail operations
// This is the CORE interactor for flight segment tracking
type DailyLogbookDetailInteractor struct {
	service        input.DailyLogbookDetailService
	logbookService input.DailyLogbookService          // For ownership verification
	summaryService input.EmployeeFlightSummaryService // For hours accumulation
}

// NewDailyLogbookDetailInteractor creates a new DailyLogbookDetailInteractor
func NewDailyLogbookDetailInteractor(
	service input.DailyLogbookDetailService,
	logbookService input.DailyLogbookService,
	summaryService input.EmployeeFlightSummaryService,
) *DailyLogbookDetailInteractor {
	return &DailyLogbookDetailInteractor{
		service:        service,
		logbookService: logbookService,
		summaryService: summaryService,
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

	if err = i.validateTimeFields(detail); err != nil {
		log.Error(logger.LogDailyLogbookDetailCreateError, "trace_id", traceID, "error", "invalid time sequence")
		return err
	}

	// Generate ID if not set
	if detail.ID == "" {
		detail.SetID()
	}

	if err = i.checkDuplicateFlight(ctx, traceID, detail); err != nil {
		return err
	}

	err = helpers.RunWithTx(ctx, i.service, log, logger.LogDailyLogbookDetailCreateError,
		func(ctx context.Context, tx output.Tx) error {
			if err := i.service.CreateDailyLogbookDetailTx(ctx, tx, detail); err != nil {
				return err
			}
			i.accumulateHours(ctx, tx, traceID, employeeID, detail, false, "accumulate_on_create")
			return nil
		})
	if err != nil {
		return err
	}

	log.Info(logger.LogDailyLogbookDetailCreateOK, "trace_id", traceID, "id", detail.ID)
	return nil
}

// UpdateDailyLogbookDetail updates an existing detail, verifying ownership
func (i *DailyLogbookDetailInteractor) UpdateDailyLogbookDetail(ctx context.Context, traceID string, detail domain.DailyLogbookDetail, employeeID string) (err error) {
	log.Info(logger.LogDailyLogbookDetailUpdate, "trace_id", traceID, "data", detail.ToLogger())

	existing, err := i.fetchExistingDetail(ctx, traceID, detail.ID)
	if err != nil {
		return err
	}

	// Verify ownership and active status via detail's logbook
	if err = i.VerifyLogbookActiveAndOwned(ctx, existing.DailyLogbookID, employeeID); err != nil {
		log.Warn(logger.LogDailyLogbookDetailUpdateError, "trace_id", traceID, "error", err)
		return err
	}

	if err = i.validateTimeFields(detail); err != nil {
		log.Error(logger.LogDailyLogbookDetailUpdateError, "trace_id", traceID, "error", "invalid time sequence")
		return err
	}

	// Preserve the daily_logbook_id from existing record (cannot change parent)
	detail.DailyLogbookID = existing.DailyLogbookID

	err = helpers.RunWithTx(ctx, i.service, log, logger.LogDailyLogbookDetailUpdateError,
		func(ctx context.Context, tx output.Tx) error {
			// First, reverse the old detail's accumulation
			i.accumulateHours(ctx, tx, traceID, employeeID, *existing, true, "reverse_on_update")
			if err := i.service.UpdateDailyLogbookDetailTx(ctx, tx, detail); err != nil {
				return err
			}
			// Then, accumulate the new detail
			i.accumulateHours(ctx, tx, traceID, employeeID, detail, false, "accumulate_on_update")
			return nil
		})
	if err != nil {
		return err
	}

	log.Info(logger.LogDailyLogbookDetailUpdateOK, "trace_id", traceID, "id", detail.ID)
	return nil
}

// DeleteDailyLogbookDetail deletes a detail
func (i *DailyLogbookDetailInteractor) DeleteDailyLogbookDetail(ctx context.Context, traceID string, id string) error {
	log.Info(logger.LogDailyLogbookDetailDelete, "trace_id", traceID, "id", id)

	// Verify detail exists
	existing, err := i.service.GetDailyLogbookDetailByID(ctx, id)
	if err != nil {
		log.Error(logger.LogDailyLogbookDetailDeleteError, "trace_id", traceID, "error", err)
		return err
	}
	if existing == nil {
		log.Warn(logger.LogDailyLogbookDetailNotFound, "trace_id", traceID, "id", id)
		return domain.ErrFlightNotFound
	}

	err = helpers.RunWithTx(ctx, i.service, log, logger.LogDailyLogbookDetailDeleteError,
		func(ctx context.Context, tx output.Tx) error {
			if err := i.service.DeleteDailyLogbookDetailTx(ctx, tx, id); err != nil {
				return err
			}
			// Reverse the deleted detail's accumulation
			if i.summaryService != nil {
				if err := i.summaryService.AccumulateFlightHours(ctx, tx, "", *existing, true); err != nil {
					log.Warn(logger.LogFlightSummaryGetError, "trace_id", traceID, "action", "reverse_on_delete", "error", err)
				}
			}
			return nil
		})
	if err != nil {
		return err
	}

	log.Info(logger.LogDailyLogbookDetailDeleteOK, "trace_id", traceID, "id", id)
	return nil
}

// ────────────────────────────────────────────────────
// Private helpers (extracted to reduce cognitive complexity)
// ────────────────────────────────────────────────────

// validateTimeFields validates the time sequence if all 4 time fields are provided.
func (i *DailyLogbookDetailInteractor) validateTimeFields(detail domain.DailyLogbookDetail) error {
	if detail.OutTime == nil || detail.TakeoffTime == nil || detail.LandingTime == nil || detail.InTime == nil {
		return nil
	}
	return i.service.ValidateTimeSequence(*detail.OutTime, *detail.TakeoffTime, *detail.LandingTime, *detail.InTime)
}

// checkDuplicateFlight checks for duplicate flights when EmployeeLogbookID is set.
func (i *DailyLogbookDetailInteractor) checkDuplicateFlight(ctx context.Context, traceID string, detail domain.DailyLogbookDetail) error {
	if detail.EmployeeLogbookID == nil {
		return nil
	}
	exists, err := i.service.ExistsByUniqueKey(ctx, *detail.EmployeeLogbookID, detail.FlightRealDate, detail.FlightNumber, detail.TailNumberID)
	if err != nil {
		log.Error(logger.LogDailyLogbookDetailCreateError, "trace_id", traceID, "error", err)
		return err
	}
	if exists {
		log.Warn(logger.LogDailyLogbookDetailDuplicate, "trace_id", traceID,
			"employee_logbook_id", *detail.EmployeeLogbookID,
			"flight_real_date", detail.FlightRealDate,
			"flight_number", detail.FlightNumber,
			"tail_number_id", detail.TailNumberID)
		return domain.ErrFlightDuplicate
	}
	return nil
}

// accumulateHours safely accumulates or reverses flight hours via summaryService.
func (i *DailyLogbookDetailInteractor) accumulateHours(ctx context.Context, tx output.Tx, traceID, employeeID string, detail domain.DailyLogbookDetail, isDelete bool, action string) {
	if i.summaryService == nil {
		return
	}
	if err := i.summaryService.AccumulateFlightHours(ctx, tx, employeeID, detail, isDelete); err != nil {
		log.Warn(logger.LogFlightSummaryGetError, "trace_id", traceID, "action", action, "error", err)
	}
}

// fetchExistingDetail fetches and validates the existence of a detail record.
func (i *DailyLogbookDetailInteractor) fetchExistingDetail(ctx context.Context, traceID, detailID string) (*domain.DailyLogbookDetail, error) {
	existing, err := i.service.GetDailyLogbookDetailByID(ctx, detailID)
	if err != nil {
		log.Error(logger.LogDailyLogbookDetailUpdateError, "trace_id", traceID, "error", err)
		return nil, err
	}
	if existing == nil {
		log.Warn(logger.LogDailyLogbookDetailNotFound, "trace_id", traceID, "id", detailID)
		return nil, domain.ErrFlightNotFound
	}
	return existing, nil
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
