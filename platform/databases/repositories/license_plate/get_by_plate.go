package licenseplate

import (
	"context"
	"database/sql"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
)

func (r *repository) GetLicensePlateByPlate(ctx context.Context, plate string) (*domain.LicensePlate, error) {
	var ar LicensePlate
	err := r.stmtGetByLicensePlate.QueryRowContext(ctx, plate).Scan(
		&ar.ID,
		&ar.LicensePlate,
		&ar.AircraftModelID,
		&ar.AirlineID,
		&ar.ModelName,
		&ar.AirlineName,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrLicensePlateNotFound
		}
		return nil, err
	}
	return ar.ToDomain(), nil
}
