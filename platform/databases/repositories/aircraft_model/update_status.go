package aircraftmodel

import (
	"context"

	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/databases/common"
)

// BeginTx starts a new database transaction
func (r *repository) BeginTx(ctx context.Context) (output.Tx, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return common.NewSQLTx(tx), nil
}

// UpdateAircraftModelStatus updates the status of an aircraft model (active/inactive)
// HU41: Inactivar la información del Tipo Aeronave
// HU42: Activar la información del Tipo Aeronave
// Idempotent: if the status is already the desired value, MySQL reports 0 rows affected
// which is NOT an error. Existence is validated beforehand in the interactor.
func (r *repository) UpdateAircraftModelStatus(ctx context.Context, tx output.Tx, id string, status bool) error {
	sqlTx := tx.(*common.SQLTX)

	_, err := sqlTx.ExecContext(ctx, QueryUpdateStatus, status, id)
	return err
}
