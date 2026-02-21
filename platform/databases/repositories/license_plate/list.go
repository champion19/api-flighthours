package licenseplate

import (
	"context"
	"database/sql"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
)

func (r *repository) ListLicensePlates(ctx context.Context, filters map[string]interface{}) ([]domain.LicensePlate, error) {
	rows, err := r.resolveListQuery(ctx, filters)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var registrations []domain.LicensePlate
	for rows.Next() {
		var ar LicensePlate
		if err := rows.Scan(
			&ar.ID,
			&ar.LicensePlate,
			&ar.AircraftModelID,
			&ar.AirlineID,
			&ar.ModelName,
			&ar.AirlineName,
		); err != nil {
			return nil, err
		}
		registrations = append(registrations, *ar.ToDomain())
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return registrations, nil
}

func (r *repository) resolveListQuery(ctx context.Context, filters map[string]interface{}) (*sql.Rows, error) {
	if val, ok := filters["license_plate"]; ok {
		if plate, isStr := val.(string); isStr && plate != "" {
			return r.stmtGetByLicensePlate.QueryContext(ctx, plate)
		}
	}

	if val, ok := filters["airline_id"]; ok {
		if airlineID, isStr := val.(string); isStr && airlineID != "" {
			return r.stmtGetByAirline.QueryContext(ctx, airlineID)
		}
	}

	return r.stmtGetAll.QueryContext(ctx)
}
