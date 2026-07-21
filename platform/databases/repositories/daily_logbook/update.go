package daily_logbook

import (
	"context"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/databases/common"
)

// UpdateDailyLogbook updates an existing daily logbook entry
func (r *repository) UpdateDailyLogbook(ctx context.Context, tx output.Tx, logbook domain.DailyLogbook) error {
	sqlTx, err := common.CastTx(tx)
	if err != nil {
		return err
	}

	result, err := sqlTx.ExecContext(ctx, QueryUpdate,
		logbook.LogDate,
		logbook.BookPage,
		logbook.Status,
		logbook.TailNumberID,
		logbook.ID,
	)
	if err != nil {
		return domain.ErrDailyLogbookCannotUpdate
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrDailyLogbookNotFound
	}

	return nil
}

// UpdateDailyLogbookStatus updates only the status of a daily logbook
// Idempotent: if the status is already the desired value, MySQL reports 0 rows affected
// which is NOT an error. Existence is validated beforehand in the interactor.
func (r *repository) UpdateDailyLogbookStatus(ctx context.Context, tx output.Tx, id string, status bool) error {
	sqlTx, err := common.CastTx(tx)
	if err != nil {
		return err
	}

	_, err = sqlTx.ExecContext(ctx, QueryUpdateStatus, status, id)
	return err
}
