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
		log.Debug("GetAirlineEmployeeByID: error", "id", id, "error", err)
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
