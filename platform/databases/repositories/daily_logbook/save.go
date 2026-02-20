package daily_logbook

import (
	"context"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/databases/common"
)

// SaveDailyLogbook creates a new daily logbook entry
func (r *repository) SaveDailyLogbook(ctx context.Context, tx output.Tx, logbook domain.DailyLogbook) error {
	sqlTx, err := common.CastTx(tx)
	if err != nil {
		return err
	}

	_, err = sqlTx.ExecContext(ctx, QueryInsert,
		logbook.ID,
		logbook.LogDate,
		logbook.EmployeeID,
		logbook.BookPage,
		logbook.Status,
	)
	if err != nil {
		return domain.ErrDailyLogbookCannotSave
	}

	return nil
}
