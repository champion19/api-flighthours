package daily_logbook

import (
	"context"

	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/databases/common"
)

// UpdateDailyLogbookStatus updates only the status of a daily logbook
// Idempotent: if the status is already the desired value, MySQL reports 0 rows affected
// which is NOT an error. Existence is validated beforehand in the interactor.
func (r *repository) UpdateDailyLogbookStatus(ctx context.Context, tx output.Tx, id string, status bool) error {
	sqlTx := tx.(*common.SQLTX)

	_, err := sqlTx.ExecContext(ctx, QueryUpdateStatus, status, id)
	return err
}
