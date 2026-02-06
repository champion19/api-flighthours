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

func (s *AirlineService) UpdateAirlineStatus(ctx context.Context, id string, status bool) error {
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err = s.repo.UpdateAirlineStatus(ctx, tx, id, status); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *AirlineService) ActivateAirline(ctx context.Context, id string) error {
	return s.UpdateAirlineStatus(ctx, id, true)
}

func (s *AirlineService) DeactivateAirline(ctx context.Context, id string) error {
	return s.UpdateAirlineStatus(ctx, id, false)
}
