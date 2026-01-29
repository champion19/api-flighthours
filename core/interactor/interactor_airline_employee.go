package interactor

import (
	"context"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/input"
	"github.com/champion19/api-flighthours/middleware"
	"github.com/champion19/api-flighthours/platform/logger"
)

// AirlineEmployeeInteractor orchestrates airline employee operations (Release 15)
type AirlineEmployeeInteractor struct {
	service input.AirlineEmployeeService
}

// NewAirlineEmployeeInteractor creates a new airline employee interactor
func NewAirlineEmployeeInteractor(service input.AirlineEmployeeService) *AirlineEmployeeInteractor {
	return &AirlineEmployeeInteractor{
		service: service,
	}
}

// GetAirlineEmployeeByID retrieves airline info for an employee by their ID (HU24)
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

// AddAirlineEmployee adds or updates airline info for an existing employee (HU26)
// The employee must already exist in the system
func (i *AirlineEmployeeInteractor) AddAirlineEmployee(ctx context.Context, employeeID string, airlineInfo domain.AirlineEmployee) error {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogAirlineEmployeeCreate, "operation", "add_airline_employee", "employee_id", employeeID)

	// Set the ID to match the employee
	airlineInfo.ID = employeeID

	if err := i.service.AddAirlineEmployee(ctx, airlineInfo); err != nil {
		log.Error(logger.LogAirlineEmployeeCreateError, "operation", "add_airline_employee", "employee_id", employeeID, "error", err)
		return err
	}

	log.Success(logger.LogAirlineEmployeeUpdateOK, "employee_id", employeeID, "airline_id", airlineInfo.AirlineID)
	return nil
}

// UpdateAirlineEmployee updates airline info for an existing employee (HU25)
// The employee must already have airline info assigned
func (i *AirlineEmployeeInteractor) UpdateAirlineEmployee(ctx context.Context, employeeID string, airlineInfo domain.AirlineEmployee) error {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogAirlineEmployeeUpdate, "operation", "update_airline_employee", "employee_id", employeeID)

	// Set the ID to match the employee
	airlineInfo.ID = employeeID

	if err := i.service.UpdateAirlineEmployee(ctx, airlineInfo); err != nil {
		log.Error(logger.LogAirlineEmployeeUpdateError, "operation", "update_airline_employee", "employee_id", employeeID, "error", err)
		return err
	}

	log.Success(logger.LogAirlineEmployeeUpdateOK, "employee_id", employeeID, "airline_id", airlineInfo.AirlineID)
	return nil
}
