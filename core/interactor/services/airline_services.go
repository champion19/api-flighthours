package services

import (
	"context"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/logger"
)

// AirlineService implements the business logic for airline operations
type AirlineService struct {
	repo   output.AirlineRepository
	logger logger.Logger
}

// NewAirlineService creates a new airline service
func NewAirlineService(repo output.AirlineRepository, log logger.Logger) *AirlineService {
	return &AirlineService{
		repo:   repo,
		logger: log,
	}
}

// BeginTx starts a new database transaction
func (s *AirlineService) BeginTx(ctx context.Context) (output.Tx, error) {
	return s.repo.BeginTx(ctx)
}

// GetAirlineByID retrieves an airline by its ID
func (s *AirlineService) GetAirlineByID(ctx context.Context, id string) (*domain.Airline, error) {
	return s.repo.GetAirlineByID(ctx, id)
}

// ListAirlines retrieves all airlines with optional filters
func (s *AirlineService) ListAirlines(ctx context.Context, filters map[string]interface{}) ([]domain.Airline, error) {
	return s.repo.ListAirlines(ctx, filters)
}
