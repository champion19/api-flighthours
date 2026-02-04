package services

import (
	"context"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
)

// AirlineRouteService implements the business logic for airline routes
type AirlineRouteService struct {
	repo output.AirlineRouteRepository
}

// NewAirlineRouteService creates a new airline route service
func NewAirlineRouteService(repo output.AirlineRouteRepository) *AirlineRouteService {
	return &AirlineRouteService{
		repo: repo,
	}
}

// BeginTx starts a new database transaction
func (s *AirlineRouteService) BeginTx(ctx context.Context) (output.Tx, error) {
	return s.repo.BeginTx(ctx)
}

// GetAirlineRouteByID retrieves an airline route by ID
func (s *AirlineRouteService) GetAirlineRouteByID(ctx context.Context, id string) (*domain.AirlineRoute, error) {
	return s.repo.GetAirlineRouteByID(ctx, id)
}

// ListAirlineRoutes retrieves all airline routes with optional filters
func (s *AirlineRouteService) ListAirlineRoutes(ctx context.Context, filters map[string]interface{}) ([]domain.AirlineRoute, error) {
	return s.repo.ListAirlineRoutes(ctx, filters)
}

func (s *AirlineRouteService) ListAirlineRoutesByAirlineID(ctx context.Context, airlineID string) ([]domain.AirlineRoute, error) {
	return s.repo.ListAirlineRoutesByAirlineID(ctx, airlineID)
}

// ActivateAirlineRoute activates an airline route
func (s *AirlineRouteService) ActivateAirlineRoute(ctx context.Context, id string) error {
	// First check if the airline route exists
	_, err := s.repo.GetAirlineRouteByID(ctx, id)
	if err != nil {
		return err
	}

	// Start transaction
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return err
	}

	// Update status
	err = s.repo.UpdateAirlineRouteStatus(ctx, tx, id, true)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// DeactivateAirlineRoute deactivates an airline route
func (s *AirlineRouteService) DeactivateAirlineRoute(ctx context.Context, id string) error {
	// First check if the airline route exists
	_, err := s.repo.GetAirlineRouteByID(ctx, id)
	if err != nil {
		return err
	}

	// Start transaction
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return err
	}

	// Update status
	err = s.repo.UpdateAirlineRouteStatus(ctx, tx, id, false)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
