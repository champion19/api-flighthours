package airport

import (
	"context"

	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/databases/common"
)

// UpdateAirportStatus updates the status of an airport (active/inactive)
// Idempotent: if the status is already the desired value, MySQL reports 0 rows affected
// which is NOT an error. Existence is validated beforehand in the interactor.
func (r *repository) UpdateAirportStatus(ctx context.Context, tx output.Tx, id string, status bool) error {
	sqlTx := tx.(*common.SQLTX)

	_, err := sqlTx.ExecContext(ctx, QueryUpdateStatus, status, id)
	return err
}
