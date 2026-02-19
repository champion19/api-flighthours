package helpers

import (
	"context"

	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/logger"
)

// TxBeginner is the minimal interface required to start a transaction.
// All services that support transactions already implement this method.
type TxBeginner interface {
	BeginTx(ctx context.Context) (output.Tx, error)
}

// RunWithTx encapsulates the common transactional pattern:
// BeginTx → defer Rollback on error → execute fn → Commit.
// This eliminates the duplicated boilerplate across interactors.
func RunWithTx(ctx context.Context, tb TxBeginner, log logger.Logger, errMsg string,
	fn func(ctx context.Context, tx output.Tx) error) (err error) {

	tx, err := tb.BeginTx(ctx)
	if err != nil {
		log.Error(errMsg, "error", err)
		return err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(errMsg, "rollback_error", rbErr, "original_error", err)
			} else {
				log.Warn(errMsg, "rollback", "ok")
			}
		}
	}()

	if err = fn(ctx, tx); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Error(errMsg, "commit_error", err)
		return err
	}

	return nil
}
