package licenseplate

import(
"context"
"strings"

 domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
 "github.com/champion19/api-flighthours/core/ports/output"
 "github.com/champion19/api-flighthours/platform/databases/common"
)


func (r *repository) SaveLicensePlate(ctx context.Context, tx output.Tx, registration domain.LicensePlate) error {
	sqlTx := tx.(*common.SQLTX)

	_, err := sqlTx.ExecContext(ctx, QueryInsert,
		registration.ID,
		registration.LicensePlate,
		registration.AircraftModelID,
		registration.AirlineID,
	)
	if err != nil {
		log.Error("SaveLicensePlate failed",
			"id", registration.ID,
			"license_plate", registration.LicensePlate,
			"aircraft_model_id", registration.AircraftModelID,
			"airline_id", registration.AirlineID,
			"error", err.Error())


		if strings.Contains(err.Error(), "Duplicate entry") {
			return domain.ErrLicensePlateDuplicatePlate
		}
		if strings.Contains(err.Error(), "foreign key constraint") || strings.Contains(err.Error(), "FOREIGN KEY") {

			if strings.Contains(err.Error(), "aircraft_model") {
				return domain.ErrLicensePlateInvalidModel
			}
			if strings.Contains(err.Error(), "airline") {
				return domain.ErrLicensePlateInvalidAirline
			}
		}
		return domain.ErrLicensePlateCannotSave
	}
	return nil
}
