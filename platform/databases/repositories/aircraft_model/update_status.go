package aircraftmodel

import(
	"context"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
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
func (r *repository) UpdateAircraftModelStatus(ctx context.Context, tx output.Tx, id string, status bool) error {
	sqlTx := tx.(*common.SQLTX)

	result, err := sqlTx.ExecContext(ctx, QueryUpdateStatus, status, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrAircraftModelNotFound
	}

	return nil
}
