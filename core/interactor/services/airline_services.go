package services

import (
	"context"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
)

type AirlineService struct {
	repo output.AirlineRepository
}

func NewAirlineService(repo output.AirlineRepository) *AirlineService {
	return &AirlineService{
		repo: repo,
	}
}

func (s *AirlineService) BeginTx(ctx context.Context) (output.Tx, error) {
	return s.repo.BeginTx(ctx)
}
func (s *AirlineService) GetAirlineByID(ctx context.Context, id string) (*domain.Airline, error) {
	return s.repo.GetAirlineByID(ctx, id)
}

func (s *AirlineService) ListAirlines(ctx context.Context, filters map[string]interface{}) ([]domain.Airline, error) {
	return s.repo.ListAirlines(ctx, filters)
}

// ActivateAirlineTx sets the airline status to true using an external transaction
func (s *AirlineService) ActivateAirlineTx(ctx context.Context, tx output.Tx, id string) error {
	return s.repo.UpdateAirlineStatus(ctx, tx, id, true)
}

// DeactivateAirlineTx sets the airline status to false using an external transaction
func (s *AirlineService) DeactivateAirlineTx(ctx context.Context, tx output.Tx, id string) error {
	return s.repo.UpdateAirlineStatus(ctx, tx, id, false)
}
