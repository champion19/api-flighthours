package services

import (
	"context"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
)

// DailyLogbookService implements the business logic for daily logbook operations
type DailyLogbookService struct {
	repo output.DailyLogbookRepository
}

// NewDailyLogbookService creates a new daily logbook service
func NewDailyLogbookService(repo output.DailyLogbookRepository) *DailyLogbookService {
	return &DailyLogbookService{
		repo: repo,
	}
}

// BeginTx starts a new database transaction
func (s *DailyLogbookService) BeginTx(ctx context.Context) (output.Tx, error) {
	return s.repo.BeginTx(ctx)
}

// GetDailyLogbookByID retrieves a daily logbook by its ID
func (s *DailyLogbookService) GetDailyLogbookByID(ctx context.Context, id string) (*domain.DailyLogbook, error) {
	return s.repo.GetDailyLogbookByID(ctx, id)
}

// ListDailyLogbooksByEmployee retrieves all daily logbooks for a specific employee
func (s *DailyLogbookService) ListDailyLogbooksByEmployee(ctx context.Context, employeeID string, filters map[string]interface{}) ([]domain.DailyLogbook, error) {
	return s.repo.ListDailyLogbooksByEmployee(ctx, employeeID, filters)
}

// UpdateDailyLogbookTx updates an existing daily logbook using an external transaction
func (s *DailyLogbookService) UpdateDailyLogbookTx(ctx context.Context, tx output.Tx, logbook domain.DailyLogbook) error {
	return s.repo.UpdateDailyLogbook(ctx, tx, logbook)
}

// DeleteDailyLogbookTx deletes a daily logbook using an external transaction
func (s *DailyLogbookService) DeleteDailyLogbookTx(ctx context.Context, tx output.Tx, id string) error {
	return s.repo.DeleteDailyLogbook(ctx, tx, id)
}

// CreateDailyLogbookTx creates a new daily logbook using an external transaction
func (s *DailyLogbookService) CreateDailyLogbookTx(ctx context.Context, tx output.Tx, logbook domain.DailyLogbook) error {
	// Generate ID if not set
	if logbook.ID == "" {
		logbook.SetID()
	}

	return s.repo.SaveDailyLogbook(ctx, tx, logbook)
}

// ActivateDailyLogbookTx sets the daily logbook status to true using an external transaction
func (s *DailyLogbookService) ActivateDailyLogbookTx(ctx context.Context, tx output.Tx, id string) error {
	return s.repo.UpdateDailyLogbookStatus(ctx, tx, id, true)
}

// DeactivateDailyLogbookTx sets the daily logbook status to false using an external transaction
func (s *DailyLogbookService) DeactivateDailyLogbookTx(ctx context.Context, tx output.Tx, id string) error {
	return s.repo.UpdateDailyLogbookStatus(ctx, tx, id, false)
}
