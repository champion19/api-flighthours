package services

import (
	"context"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
)

// ManufacturerService implements the business logic for manufacturer operations
type ManufacturerService struct {
	repo output.ManufacturerRepository
}

// NewManufacturerService creates a new manufacturer service
func NewManufacturerService(repo output.ManufacturerRepository) *ManufacturerService {
	return &ManufacturerService{
		repo: repo,
	}
}

// GetManufacturerByID retrieves a manufacturer by its ID
func (s *ManufacturerService) GetManufacturerByID(ctx context.Context, id string) (*domain.Manufacturer, error) {
	return s.repo.GetManufacturerByID(ctx, id)
}

// ListManufacturers retrieves all manufacturers
func (s *ManufacturerService) ListManufacturers(ctx context.Context) ([]domain.Manufacturer, error) {
	return s.repo.ListManufacturers(ctx)
}
