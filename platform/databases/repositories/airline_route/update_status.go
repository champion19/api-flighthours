package airline_route

import (
	"context"
	"database/sql"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
)

func (r *repository) UpdateAirlineRouteStatus(ctx context.Context, tx output.Tx, id string, status string) error {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return domain.ErrInvalidTransaction
	}

	stmt := sqlTx.StmtContext(ctx, r.stmtUpdateStatus)
	result, err := stmt.ExecContext(ctx, status, id, status)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		existingRoute, err := r.GetAirlineRouteByID(ctx, id)
		if err != nil {
			return domain.ErrAirlineRouteNotFound
		}

		if existingRoute.Status == status {
			if status == domain.AirlineRouteStatusActive {
				return domain.ErrAirlineRouteAlreadyActive
			}
			return domain.ErrAirlineRouteAlreadyInactive
		}

		return domain.ErrAirlineRouteNotFound
	}

	return nil
}
