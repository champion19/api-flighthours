package services

import (
	"context"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
)

// RouteService implements the business logic for route operations
type RouteService struct {
	repo output.RouteRepository
}

// NewRouteService creates a new route service
func NewRouteService(repo output.RouteRepository) *RouteService {
	return &RouteService{
		repo: repo,
	}
}

// GetRouteByID retrieves a route by its ID
func (s *RouteService) GetRouteByID(ctx context.Context, id string) (*domain.Route, error) {
	return s.repo.GetRouteByID(ctx, id)
}

// ListRoutes retrieves all routes with optional filters
func (s *RouteService) ListRoutes(ctx context.Context, filters map[string]interface{}) ([]domain.Route, error) {
	return s.repo.ListRoutes(ctx, filters)
}

// GetRouteByAirports retrieves the physical route (if any) between two airports
func (s *RouteService) GetRouteByAirports(ctx context.Context, originAirportID, destinationAirportID string) (*domain.Route, error) {
	return s.repo.GetRouteByAirports(ctx, originAirportID, destinationAirportID)
}
