package licenseplate

import (
	"context"
	"database/sql"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
)

func (r *repository) ListLicensePlates(ctx context.Context, filters map[string]interface{}) ([]domain.LicensePlate, error) {
	var rows *sql.Rows
	var err error

	if licensePlate, ok := filters["license_plate"]; ok {
		plateStr, isString := licensePlate.(string)
		if isString && plateStr != "" {
			rows, err = r.stmtGetByLicensePlate.QueryContext(ctx, plateStr)
		} else {
			rows, err = r.stmtGetAll.QueryContext(ctx)
		}
	} else if airlineID, ok := filters["airline_id"]; ok {
		airlineIDStr, isString := airlineID.(string)
		if isString && airlineIDStr != "" {
			rows, err = r.stmtGetByAirline.QueryContext(ctx, airlineIDStr)
		} else {
			rows, err = r.stmtGetAll.QueryContext(ctx)
		}
	} else {
		rows, err = r.stmtGetAll.QueryContext(ctx)
	}

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
