package services

import (
	"context"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/logger"
)

var log logger.Logger = logger.NewSlogLogger()

type airlineEmployeeService struct {
	repository output.AirlineEmployeeRepository
}

func NewAirlineEmployeeService(repository output.AirlineEmployeeRepository) *airlineEmployeeService {
	return &airlineEmployeeService{
		repository: repository,
	}
}

func (s *airlineEmployeeService) BeginTx(ctx context.Context) (output.Tx, error) {
	return s.repository.BeginTx(ctx)
}
func (s *airlineEmployeeService) GetAirlineEmployeeByID(ctx context.Context, id string) (*domain.AirlineEmployee, error) {
	employee, err := s.repository.GetAirlineEmployeeByID(ctx, id)
	if err != nil {
		log.Debug(logger.LogAirlineEmployeeGetError, "id", id, "error", err)
		return nil, err
	}
	return employee, nil
}

func (s *airlineEmployeeService) AddAirlineEmployee(ctx context.Context, employee domain.AirlineEmployee) error {
	tx, err := s.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error", err)
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err = s.repository.AddAirlineEmployee(ctx, tx, employee); err != nil {
		log.Error(logger.LogDatabaseUnavailable, "operation", "add", "employee_id", employee.ID, "error", err)
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Error(logger.LogDatabaseUnavailable, "operation", "commit", "employee_id", employee.ID, "error", err)
		return err
	}

	log.Info(logger.LogDatabaseAvailable, "operation", "add_airline_employee", "employee_id", employee.ID)
	return nil
}

func (s *airlineEmployeeService) UpdateAirlineEmployee(ctx context.Context, employee domain.AirlineEmployee) error {
	tx, err := s.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error", err)
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err = s.repository.UpdateAirlineEmployee(ctx, tx, employee); err != nil {
		log.Error(logger.LogDatabaseUnavailable, "operation", "update", "employee_id", employee.ID, "error", err)
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Error(logger.LogDatabaseUnavailable, "operation", "commit", "employee_id", employee.ID, "error", err)
		return err
	}

	log.Info(logger.LogDatabaseAvailable, "operation", "update_airline_employee", "employee_id", employee.ID)
	return nil
}

func (s *airlineEmployeeService) ActivateAirlineEmployee(ctx context.Context, id string) error {
	tx, err := s.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error", err)
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err = s.repository.UpdateAirlineEmployeeStatus(ctx, tx, id, true); err != nil {
		log.Error(logger.LogDatabaseUnavailable, "operation", "activate", "employee_id", id, "error", err)
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Error(logger.LogDatabaseUnavailable, "operation", "commit", "employee_id", id, "error", err)
		return err
	}

	log.Info(logger.LogDatabaseAvailable, "operation", "activate_airline_employee", "employee_id", id)
	return nil
}

func (s *airlineEmployeeService) DeactivateAirlineEmployee(ctx context.Context, id string) error {
	tx, err := s.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error", err)
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err = s.repository.UpdateAirlineEmployeeStatus(ctx, tx, id, false); err != nil {
		log.Error(logger.LogDatabaseUnavailable, "operation", "deactivate", "employee_id", id, "error", err)
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Error(logger.LogDatabaseUnavailable, "operation", "commit", "employee_id", id, "error", err)
		return err
	}

	log.Info(logger.LogDatabaseAvailable, "operation", "deactivate_airline_employee", "employee_id", id)
	return nil
}

// AddAirlineEmployeeTx adds an airline employee using an external transaction
func (s *airlineEmployeeService) AddAirlineEmployeeTx(ctx context.Context, tx output.Tx, employee domain.AirlineEmployee) error {
	return s.repository.AddAirlineEmployee(ctx, tx, employee)
}

// UpdateAirlineEmployeeTx updates an airline employee using an external transaction
func (s *airlineEmployeeService) UpdateAirlineEmployeeTx(ctx context.Context, tx output.Tx, employee domain.AirlineEmployee) error {
	return s.repository.UpdateAirlineEmployee(ctx, tx, employee)
}

// ActivateAirlineEmployeeTx activates an airline employee using an external transaction
func (s *airlineEmployeeService) ActivateAirlineEmployeeTx(ctx context.Context, tx output.Tx, id string) error {
	return s.repository.UpdateAirlineEmployeeStatus(ctx, tx, id, true)
}

// DeactivateAirlineEmployeeTx deactivates an airline employee using an external transaction
func (s *airlineEmployeeService) DeactivateAirlineEmployeeTx(ctx context.Context, tx output.Tx, id string) error {
	return s.repository.UpdateAirlineEmployeeStatus(ctx, tx, id, false)
}
