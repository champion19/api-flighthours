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
		log.Debug(logger.LogAirlineEmployeeServiceGetError, "id", id, "error", err)
		return nil, err
	}
	return employee, nil
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
