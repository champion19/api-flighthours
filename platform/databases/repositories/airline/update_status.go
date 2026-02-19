package airline

import (
	"context"
	"database/sql"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/databases/common"
	"github.com/champion19/api-flighthours/platform/logger"
)

const (
	QueryUpdateStatus = `UPDATE airline SET status = ? WHERE id = ?`
)

// UpdateAirlineStatus updates only the active status of an airline (HU3/HU4)
// This operation is IDEMPOTENT - calling it multiple times with the same status returns success
func (r *repository) UpdateAirlineStatus(ctx context.Context, tx output.Tx, id string, active bool) error {
	log.Debug(logger.LogAirlineRepoUpdateStatusStart,
		"airline_id", id,
		"active", active)

	// Cast the transaction to the concrete type
	dbTx, ok := tx.(*common.SQLTX)
	if !ok {
		log.Error(logger.LogDatabaseUnavailable, "error", logger.LogErrInvalidTransaction)
		return domain.ErrInvalidTransaction
	}

	// First, check if the airline exists
	var exists bool
	err := dbTx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM airline WHERE id = ?)", id).Scan(&exists)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "airline_id", id, "error", err)
		if err == sql.ErrNoRows {
			return domain.ErrAirlineNotFound
		}
		return err
	}

	if !exists {
		log.Warn(logger.LogAirlineRepoUpdateStatusNotFound, "airline_id", id)
		return domain.ErrAirlineNotFound
	}

	// Update the status (idempotent - OK if already in desired state)
	_, err = dbTx.ExecContext(ctx, QueryUpdateStatus, active, id)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable,
			"airline_id", id,
			"error", err)
		return err
	}

	log.Debug(logger.LogAirlineRepoUpdateStatusOK,
		"airline_id", id,
		"active", active)

	return nil
}
