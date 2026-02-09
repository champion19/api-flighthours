package licenseplate

import (
	"context"
	"strings"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/databases/common"
)

func (r *repository) UpdateLicensePlate(ctx context.Context, tx output.Tx, registration domain.LicensePlate) error {
	sqlTx := tx.(*common.SQLTX)

	result, err := sqlTx.ExecContext(ctx, QueryUpdate,
		registration.LicensePlate,
		registration.AircraftModelID,
		registration.AirlineID,
		registration.ID,
	)
	if err != nil {
		log.Error("UpdateLicensePlate failed",
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
		return domain.ErrLicensePlateCannotUpdate
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return domain.ErrLicensePlateCannotUpdate
	}

	if rowsAffected == 0 {
		return domain.ErrLicensePlateNotFound
	}

	return nil
}
