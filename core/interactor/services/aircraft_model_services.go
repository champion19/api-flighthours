package services

import (
	"context"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/logger"
)

// AircraftModelService implements the business logic for aircraft model operations
type AircraftModelService struct {
	repo   output.AircraftModelRepository
	logger logger.Logger
}

// NewAircraftModelService creates a new aircraft model service
func NewAircraftModelService(repo output.AircraftModelRepository, log logger.Logger) *AircraftModelService {
	return &AircraftModelService{
		repo:   repo,
		logger: log,
	}
}

// BeginTx starts a new database transaction
func (s *AircraftModelService) BeginTx(ctx context.Context) (output.Tx, error) {
	return s.repo.BeginTx(ctx)
}

// GetAircraftModelByID retrieves an aircraft model by its ID
func (s *AircraftModelService) GetAircraftModelByID(ctx context.Context, id string) (*domain.AircraftModel, error) {
	return s.repo.GetAircraftModelByID(ctx, id)
}

// ListAircraftModels retrieves all aircraft models with optional filters
func (s *AircraftModelService) ListAircraftModels(ctx context.Context, filters map[string]interface{}) ([]domain.AircraftModel, error) {
	return s.repo.ListAircraftModels(ctx, filters)
}

// GetAircraftModelsByFamily retrieves all aircraft models for a specific family (HU30)
func (s *AircraftModelService) GetAircraftModelsByFamily(ctx context.Context, family string) ([]domain.AircraftModel, error) {
	return s.repo.GetAircraftModelsByFamily(ctx, family)
}

// UpdateAircraftModelStatus updates the status of an aircraft model with transaction handling
func (s *AircraftModelService) UpdateAircraftModelStatus(ctx context.Context, id string, status bool) error {
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err = s.repo.UpdateAircraftModelStatus(ctx, tx, id, status); err != nil {
		return err
	}

	return tx.Commit()
}

// ActivateAircraftModel sets the aircraft model status to true (active) - HU42
func (s *AircraftModelService) ActivateAircraftModel(ctx context.Context, id string) error {
	return s.UpdateAircraftModelStatus(ctx, id, true)
}

// DeactivateAircraftModel sets the aircraft model status to false (inactive) - HU41
func (s *AircraftModelService) DeactivateAircraftModel(ctx context.Context, id string) error {
	return s.UpdateAircraftModelStatus(ctx, id, false)
}

// ActivateAircraftModelTx sets the aircraft model status to true using an external transaction
func (s *AircraftModelService) ActivateAircraftModelTx(ctx context.Context, tx output.Tx, id string) error {
	return s.repo.UpdateAircraftModelStatus(ctx, tx, id, true)
}

// DeactivateAircraftModelTx sets the aircraft model status to false using an external transaction
func (s *AircraftModelService) DeactivateAircraftModelTx(ctx context.Context, tx output.Tx, id string) error {
	return s.repo.UpdateAircraftModelStatus(ctx, tx, id, false)
}
