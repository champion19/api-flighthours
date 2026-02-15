package daily_logbook

import (
	"context"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/databases/common"
)

// DeleteDailyLogbook removes a daily logbook and its associated details (cascade)
func (r *repository) DeleteDailyLogbook(ctx context.Context, tx output.Tx, id string) error {
	sqlTx := tx.(*common.SQLTX)

	// 1. Cascade: delete associated daily_logbook_detail records first
	if _, err := sqlTx.ExecContext(ctx, QueryDeleteDetails, id); err != nil {
		return domain.ErrDailyLogbookCannotDelete
	}

	// 2. Delete the daily logbook itself
	result, err := sqlTx.ExecContext(ctx, QueryDelete, id)
	if err != nil {
		return domain.ErrDailyLogbookCannotDelete
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
