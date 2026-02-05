package airport

import (
	"context"
	"database/sql"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
)


// GetAirportsByCity retrieves all airports for a specific city (HU13 - Virtual Entity pattern)
// This implements the "Derived Values" approach - no new table needed, we query airports by city field
func (r *repository) GetAirportsByCity(ctx context.Context, city string) ([]domain.Airport, error) {
	rows, err := r.stmtGetByCity.QueryContext(ctx, city)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var airports []domain.Airport
	for rows.Next() {
		var a Airport
		if err := rows.Scan(&a.ID, &a.Name, &a.City, &a.Country, &a.IATACode, &a.Status, &a.AirportType); err != nil {
			return nil, err
		}
		airports = append(airports, *a.ToDomain())
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Return empty slice with no error if no airports found (caller handles not found logic)
	if len(airports) == 0 {
		return nil, sql.ErrNoRows
	}

	return airports, nil
}

// GetAirportsByCountry retrieves all airports for a specific country (HU38 - Virtual Entity pattern)
// This implements the "Derived Values" approach - no new table needed, we query airports by country field
func (r *repository) GetAirportsByCountry(ctx context.Context, country string) ([]domain.Airport, error) {
	rows, err := r.stmtGetByCountry.QueryContext(ctx, country)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var airports []domain.Airport
	for rows.Next() {
		var a Airport
		if err := rows.Scan(&a.ID, &a.Name, &a.City, &a.Country, &a.IATACode, &a.Status, &a.AirportType); err != nil {
			return nil, err
		}
		airports = append(airports, *a.ToDomain())
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Return empty slice with no error if no airports found (caller handles not found logic)
	if len(airports) == 0 {
		return nil, sql.ErrNoRows
	}

	return airports, nil
}

// GetAirportsByType retrieves all airports for a specific airport type (HU46 - Virtual Entity pattern)
// This implements the "Derived Values" approach - no new table needed, we query airports by airport_type field
func (r *repository) GetAirportsByType(ctx context.Context, airportType string) ([]domain.Airport, error) {
	rows, err := r.stmtGetByType.QueryContext(ctx, airportType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var airports []domain.Airport
	for rows.Next() {
		var a Airport
		if err := rows.Scan(&a.ID, &a.Name, &a.City, &a.Country, &a.IATACode, &a.Status, &a.AirportType); err != nil {
			return nil, err
		}
		airports = append(airports, *a.ToDomain())
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Return empty slice with no error if no airports found (caller handles not found logic)
	if len(airports) == 0 {
		return nil, sql.ErrNoRows
	}

	return airports, nil
}
