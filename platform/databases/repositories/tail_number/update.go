package tailnumber

import (
	"context"
	"strings"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/databases/common"
)

func (r *repository) UpdateTailNumber(ctx context.Context, tx output.Tx, registration domain.TailNumber) error {
	sqlTx, err := common.CastTx(tx)
	if err != nil {
		return err
	}

	result, err := sqlTx.ExecContext(ctx, QueryUpdate,
		registration.TailNumber,
		registration.AircraftModelID,
		registration.AirlineID,
		registration.ID,
	)
	if err != nil {
		log.Error("UpdateTailNumber failed",
			"id", registration.ID,
			"tail_number", registration.TailNumber,
			"aircraft_model_id", registration.AircraftModelID,
			"airline_id", registration.AirlineID,
			"error", err.Error())

		if strings.Contains(err.Error(), "Duplicate entry") {
			return domain.ErrTailNumberDuplicatePlate
		}
		if strings.Contains(err.Error(), "foreign key constraint") || strings.Contains(err.Error(), "FOREIGN KEY") {
			if strings.Contains(err.Error(), "aircraft_model") {
				return domain.ErrTailNumberInvalidModel
			}
			if strings.Contains(err.Error(), "airline") {
				return domain.ErrTailNumberInvalidAirline
			}
		}
		return domain.ErrTailNumberCannotUpdate
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return domain.ErrTailNumberCannotUpdate
	}

	if rowsAffected == 0 {
		return domain.ErrTailNumberNotFound
	}

	return nil
}
