package services

import (
	"context"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
)


type LicensePlateService struct {
	repo   output.LicensePlateRepository
}


func NewLicensePlateService(repo output.LicensePlateRepository) *LicensePlateService {
	return &LicensePlateService{
		repo:   repo,

	}
}


func (s *LicensePlateService) BeginTx(ctx context.Context) (output.Tx, error) {
	return s.repo.BeginTx(ctx)
}


func (s *LicensePlateService) GetLicensePlateByID(ctx context.Context, id string) (*domain.LicensePlate, error) {
	return s.repo.GetLicensePlateByID(ctx, id)
}


func (s *LicensePlateService) ListLicensePlates(ctx context.Context, filters map[string]interface{}) ([]domain.LicensePlate, error) {
	return s.repo.ListLicensePlates(ctx, filters)
}


func (s *LicensePlateService) CreateLicensePlate(ctx context.Context, registration domain.LicensePlate) error {
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err = s.repo.SaveLicensePlate(ctx, tx, registration); err != nil {
		return err
	}

	return tx.Commit()
}


func (s *LicensePlateService) UpdateLicensePlate(ctx context.Context, registration domain.LicensePlate) error {
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err = s.repo.UpdateLicensePlate(ctx, tx, registration); err != nil {
		return err
	}

	return tx.Commit()
}
