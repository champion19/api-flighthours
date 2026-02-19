package daily_logbook_detail

import (
	"context"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/databases/common"
	"github.com/champion19/api-flighthours/platform/logger"
)

// DeleteDailyLogbookDetail deletes a daily logbook detail by its ID
func (r *repository) DeleteDailyLogbookDetail(ctx context.Context, tx output.Tx, id string) error {
	log.Info(logger.LogDailyLogbookDetailDelete, "id", id)

	sqlTx, err := common.CastTx(tx)
	if err != nil {
		log.Error(logger.LogDailyLogbookDetailDeleteError, "error", "invalid transaction type")
		return err
	}

	stmt := sqlTx.Tx.StmtContext(ctx, r.stmtDelete)

	result, err := stmt.ExecContext(ctx, id)
	if err != nil {
		log.Error(logger.LogDailyLogbookDetailDeleteError, "id", id, "error", err)
		return domain.ErrFlightCannotDelete
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error(logger.LogDailyLogbookDetailDeleteError, "id", id, "error", err)
		return domain.ErrFlightCannotDelete
	}

	if rowsAffected == 0 {
		log.Warn(logger.LogDailyLogbookDetailNotFound, "id", id)
		return domain.ErrFlightNotFound
	}

	log.Info(logger.LogDailyLogbookDetailDeleteOK, "id", id)
	return nil
}
