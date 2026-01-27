package airline

import (
	"context"
	"database/sql"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
)

// GetAirlineByID retrieves an airline by its UUID
func (r *repository) GetAirlineByID(ctx context.Context, id string) (*domain.Airline, error) {
	var a Airline
	err := r.stmtGetByID.QueryRowContext(ctx, id).Scan(
		&a.ID,
		&a.AirlineName,
		&a.AirlineCode,
		&a.Status,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrAirlineNotFound
		}
		return nil, err
	}
	return a.ToDomain(), nil
}
