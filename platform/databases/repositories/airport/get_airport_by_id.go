package airport

import (
	"context"
	"database/sql"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
)


func (r *repository) GetAirportByID(ctx context.Context, id string) (*domain.Airport, error) {
	var a Airport
	err := r.stmtGetByID.QueryRowContext(ctx, id).Scan(
		&a.ID,
		&a.Name,
		&a.City,
		&a.Country,
		&a.IATACode,
		&a.OACICode,
		&a.Status,
		&a.AirportType,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrAirportNotFound
		}
		return nil, err
	}
	return a.ToDomain(), nil
}
