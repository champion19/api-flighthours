package services

import (
	"context"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
)

type TailNumberService struct {
	repo output.TailNumberRepository
}

func NewTailNumberService(repo output.TailNumberRepository) *TailNumberService {
	return &TailNumberService{
		repo: repo,
	}
}

func (s *TailNumberService) BeginTx(ctx context.Context) (output.Tx, error) {
	return s.repo.BeginTx(ctx)
}

func (s *TailNumberService) GetTailNumberByID(ctx context.Context, id string) (*domain.TailNumber, error) {
	return s.repo.GetTailNumberByID(ctx, id)
}

func (s *TailNumberService) GetTailNumberByPlate(ctx context.Context, plate string) (*domain.TailNumber, error) {
	return s.repo.GetTailNumberByPlate(ctx, plate)
}

func (s *TailNumberService) ListTailNumbers(ctx context.Context, filters map[string]interface{}) ([]domain.TailNumber, error) {
	return s.repo.ListTailNumbers(ctx, filters)
}

// CreateTailNumberTx creates a license plate using an external transaction
func (s *TailNumberService) CreateTailNumberTx(ctx context.Context, tx output.Tx, registration domain.TailNumber) error {
	return s.repo.SaveTailNumber(ctx, tx, registration)
}

// UpdateTailNumberTx updates a license plate using an external transaction
func (s *TailNumberService) UpdateTailNumberTx(ctx context.Context, tx output.Tx, registration domain.TailNumber) error {
	return s.repo.UpdateTailNumber(ctx, tx, registration)
}
