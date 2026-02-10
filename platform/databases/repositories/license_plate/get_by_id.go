package licenseplate

import(
	"context"
	"database/sql"
	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
)


func (r *repository) GetLicensePlateByID(ctx context.Context, id string) (*domain.LicensePlate, error) {
	var ar LicensePlate
	err := r.stmtGetByID.QueryRowContext(ctx, id).Scan(
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
