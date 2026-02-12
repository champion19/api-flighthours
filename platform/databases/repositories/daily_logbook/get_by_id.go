package daily_logbook

import (
	"context"
	"database/sql"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
)

// GetDailyLogbookByID retrieves a daily logbook by its UUID
func (r *repository) GetDailyLogbookByID(ctx context.Context, id string) (*domain.DailyLogbook, error) {
	var d DailyLogbook
	err := r.stmtGetByID.QueryRowContext(ctx, id).Scan(
		&d.ID,
		&d.LogDate,
		&d.EmployeeID,
		&d.BookPage,
		&d.Status,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrDailyLogbookNotFound
		}
		return nil, err
	}
	return d.ToDomain(), nil
}
