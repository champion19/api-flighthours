package services

import (
	"context"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
)

// AirportService implements the business logic for airport operations
type AirportService struct {
	repo output.AirportRepository
}

// NewAirportService creates a new airport service
func NewAirportService(repo output.AirportRepository) *AirportService {
	return &AirportService{
		repo: repo,
	}
}

// BeginTx starts a new database transaction
func (s *AirportService) BeginTx(ctx context.Context) (output.Tx, error) {
	return s.repo.BeginTx(ctx)
}

// GetAirportByID retrieves an airport by its ID
func (s *AirportService) GetAirportByID(ctx context.Context, id string) (*domain.Airport, error) {
	return s.repo.GetAirportByID(ctx, id)
}

// ListAirports retrieves all airports with optional filters
func (s *AirportService) ListAirports(ctx context.Context, filters map[string]interface{}) ([]domain.Airport, error) {
	return s.repo.ListAirports(ctx, filters)
}

// GetAirportsByType retrieves all airports for a specific airport type
func (s *AirportService) GetAirportsByType(ctx context.Context, airportType string) ([]domain.Airport, error) {
	return s.repo.GetAirportsByType(ctx, airportType)
}

// UpdateAirportStatus updates the status of an airport with transaction handling
func (s *AirportService) UpdateAirportStatus(ctx context.Context, id string, status bool) error {
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err = s.repo.UpdateAirportStatus(ctx, tx, id, status); err != nil {
		return err
	}

	return tx.Commit()
}

// DeactivateAirport sets the airport status to false (inactive)
func (s *AirportService) DeactivateAirport(ctx context.Context, id string) error {
	return s.UpdateAirportStatus(ctx, id, false)
}
