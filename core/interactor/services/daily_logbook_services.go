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

// CreateDailyLogbook creates a new daily logbook entry with transaction handling
func (s *DailyLogbookService) CreateDailyLogbook(ctx context.Context, logbook domain.DailyLogbook) error {
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Generate ID if not set
	if logbook.ID == "" {
		logbook.SetID()
	}

	if err = s.repo.SaveDailyLogbook(ctx, tx, logbook); err != nil {
		return err
	}

	return tx.Commit()
}

// CreateDailyLogbookTx creates a new daily logbook using an external transaction
func (s *DailyLogbookService) CreateDailyLogbookTx(ctx context.Context, tx output.Tx, logbook domain.DailyLogbook) error {
	// Generate ID if not set
	if logbook.ID == "" {
		logbook.SetID()
	}

	return s.repo.SaveDailyLogbook(ctx, tx, logbook)
}
