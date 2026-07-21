package airport

import (
	"context"
	"database/sql"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
)

func (r *repository) GetAirportsByType(ctx context.Context, airportType string) ([]domain.Airport, error) {
	rows, err := r.stmtGetByType.QueryContext(ctx, airportType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var airports []domain.Airport
	for rows.Next() {
		var a Airport
		if err := rows.Scan(&a.ID, &a.Name, &a.City, &a.Country, &a.IATACode, &a.OACICode, &a.Status, &a.AirportType); err != nil {
			return nil, err
		}
		airports = append(airports, *a.ToDomain())
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(airports) == 0 {
		return nil, sql.ErrNoRows
	}

	return airports, nil
}
