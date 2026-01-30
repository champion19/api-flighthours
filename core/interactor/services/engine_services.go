package services

import (
	"context"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
)

type EngineService struct {
	repo output.EngineRepository
}
func NewEngineService(repo output.EngineRepository) *EngineService {
	return &EngineService{
		repo: repo,
	}
}

func (s *EngineService) GetEngineByID(ctx context.Context, id string) (*domain.Engine, error) {
	return s.repo.GetEngineByID(ctx, id)
}
func (s *EngineService) ListEngines(ctx context.Context) ([]domain.Engine, error) {
	return s.repo.ListEngines(ctx)
}
