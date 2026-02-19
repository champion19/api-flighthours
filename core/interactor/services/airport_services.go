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

// ActivateAirportTx sets the airport status to true using an external transaction
func (s *AirportService) ActivateAirportTx(ctx context.Context, tx output.Tx, id string) error {
	return s.repo.UpdateAirportStatus(ctx, tx, id, true)
}

// DeactivateAirportTx sets the airport status to false using an external transaction
func (s *AirportService) DeactivateAirportTx(ctx context.Context, tx output.Tx, id string) error {
	return s.repo.UpdateAirportStatus(ctx, tx, id, false)
}
