package interactor

import (
	"context"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/input"
	"github.com/champion19/api-flighthours/middleware"
	"github.com/champion19/api-flighthours/platform/logger"
)

type AirlineEmployeeInteractor struct {
	service input.AirlineEmployeeService
}

func NewAirlineEmployeeInteractor(service input.AirlineEmployeeService) *AirlineEmployeeInteractor {
	return &AirlineEmployeeInteractor{
		service: service,
	}
}

func (i *AirlineEmployeeInteractor) GetAirlineEmployeeByID(ctx context.Context, id string) (*domain.AirlineEmployee, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogAirlineEmployeeGet, "operation", "get_airline_employee", "employee_id", id)

	employee, err := i.service.GetAirlineEmployeeByID(ctx, id)
	if err != nil {
		log.Error(logger.LogAirlineEmployeeGetError, "operation", "get_airline_employee", "employee_id", id, "error", err)
		return nil, err
	}

	log.Success(logger.LogAirlineEmployeeGetOK, "employee_id", id)
	return employee, nil
}

func (i *AirlineEmployeeInteractor) AddAirlineEmployee(ctx context.Context, employeeID string, airlineInfo domain.AirlineEmployee) (err error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogAirlineEmployeeCreate, "operation", "add_airline_employee", "employee_id", employeeID)

	airlineInfo.ID = employeeID

	tx, err := i.service.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogAirlineEmployeeCreateError, "operation", "add_airline_employee", "employee_id", employeeID, "error", err)
		return err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(logger.LogAirlineEmployeeCreateError, "employee_id", employeeID, "rollback_error", rbErr, "original_error", err)
			} else {
				log.Warn(logger.LogAirlineEmployeeCreateError, "employee_id", employeeID, "rollback", "ok")
			}
		}
	}()

	if err = i.service.AddAirlineEmployeeTx(ctx, tx, airlineInfo); err != nil {
		log.Error(logger.LogAirlineEmployeeCreateError, "operation", "add_airline_employee", "employee_id", employeeID, "error", err)
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Error(logger.LogAirlineEmployeeCreateError, "operation", "add_airline_employee", "employee_id", employeeID, "commit_error", err)
		return err
	}

	log.Success(logger.LogAirlineEmployeeUpdateOK, "employee_id", employeeID, "airline_id", airlineInfo.AirlineID)
	return nil
}

func (i *AirlineEmployeeInteractor) UpdateAirlineEmployee(ctx context.Context, employeeID string, airlineInfo domain.AirlineEmployee) (err error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogAirlineEmployeeUpdate, "operation", "update_airline_employee", "employee_id", employeeID)

	airlineInfo.ID = employeeID

	tx, err := i.service.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogAirlineEmployeeUpdateError, "operation", "update_airline_employee", "employee_id", employeeID, "error", err)
		return err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(logger.LogAirlineEmployeeUpdateError, "employee_id", employeeID, "rollback_error", rbErr, "original_error", err)
			} else {
				log.Warn(logger.LogAirlineEmployeeUpdateError, "employee_id", employeeID, "rollback", "ok")
			}
		}
	}()

	if err = i.service.UpdateAirlineEmployeeTx(ctx, tx, airlineInfo); err != nil {
		log.Error(logger.LogAirlineEmployeeUpdateError, "operation", "update_airline_employee", "employee_id", employeeID, "error", err)
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Error(logger.LogAirlineEmployeeUpdateError, "operation", "update_airline_employee", "employee_id", employeeID, "commit_error", err)
		return err
	}

	log.Success(logger.LogAirlineEmployeeUpdateOK, "employee_id", employeeID, "airline_id", airlineInfo.AirlineID)
	return nil
}

func (i *AirlineEmployeeInteractor) ActivateAirlineEmployee(ctx context.Context, employeeID string) (err error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogAirlineEmployeeActivate, "operation", "activate_airline_employee", "employee_id", employeeID)

	existingInfo, err := i.service.GetAirlineEmployeeByID(ctx, employeeID)
	if err != nil || existingInfo == nil || existingInfo.AirlineID == "" {
		log.Error(logger.LogAirlineEmployeeNotFound, "operation", "activate_airline_employee", "employee_id", employeeID, "error", "no airline info")
		return domain.ErrAirlineEmployeeNotFound
	}

	tx, err := i.service.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogAirlineEmployeeActivateError, "operation", "activate_airline_employee", "employee_id", employeeID, "error", err)
		return err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(logger.LogAirlineEmployeeActivateError, "employee_id", employeeID, "rollback_error", rbErr, "original_error", err)
			} else {
				log.Warn(logger.LogAirlineEmployeeActivateError, "employee_id", employeeID, "rollback", "ok")
			}
		}
	}()

	if err = i.service.ActivateAirlineEmployeeTx(ctx, tx, employeeID); err != nil {
		log.Error(logger.LogAirlineEmployeeActivateError, "operation", "activate_airline_employee", "employee_id", employeeID, "error", err)
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Error(logger.LogAirlineEmployeeActivateError, "operation", "activate_airline_employee", "employee_id", employeeID, "commit_error", err)
		return err
	}

	log.Success(logger.LogAirlineEmployeeActivateOK, "employee_id", employeeID)
	return nil
}

func (i *AirlineEmployeeInteractor) DeactivateAirlineEmployee(ctx context.Context, employeeID string) (err error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogAirlineEmployeeDeactivate, "operation", "deactivate_airline_employee", "employee_id", employeeID)

	existingInfo, err := i.service.GetAirlineEmployeeByID(ctx, employeeID)
	if err != nil || existingInfo == nil || existingInfo.AirlineID == "" {
		log.Error(logger.LogAirlineEmployeeNotFound, "operation", "deactivate_airline_employee", "employee_id", employeeID, "error", "no airline info")
		return domain.ErrAirlineEmployeeNotFound
	}

	tx, err := i.service.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogAirlineEmployeeDeactivateError, "operation", "deactivate_airline_employee", "employee_id", employeeID, "error", err)
		return err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(logger.LogAirlineEmployeeDeactivateError, "employee_id", employeeID, "rollback_error", rbErr, "original_error", err)
			} else {
				log.Warn(logger.LogAirlineEmployeeDeactivateError, "employee_id", employeeID, "rollback", "ok")
			}
		}
	}()

	if err = i.service.DeactivateAirlineEmployeeTx(ctx, tx, employeeID); err != nil {
		log.Error(logger.LogAirlineEmployeeDeactivateError, "operation", "deactivate_airline_employee", "employee_id", employeeID, "error", err)
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Error(logger.LogAirlineEmployeeDeactivateError, "operation", "deactivate_airline_employee", "employee_id", employeeID, "commit_error", err)
		return err
	}

	log.Success(logger.LogAirlineEmployeeDeactivateOK, "employee_id", employeeID)
	return nil
}
