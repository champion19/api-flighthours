package airline_route

import (
	"context"
	"strings"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/databases/common"
	"github.com/champion19/api-flighthours/platform/logger"
)

// SaveAirlineRoute inserts a new airline_route link (typically with status
// Pending, auto-created when an employee's airline has no link yet for an
// existing route). Relies on the uq_airline_route_route_airline unique
// constraint to make concurrent creation attempts safe — a duplicate insert
// is reported back as ErrAirlineRouteAlreadyExists so the caller can re-fetch
// the row another request just created instead of failing outright.
func (r *repository) SaveAirlineRoute(ctx context.Context, tx output.Tx, airlineRoute domain.AirlineRoute) error {
	sqlTx, err := common.CastTx(tx)
	if err != nil {
		return err
	}

	_, err = sqlTx.ExecContext(ctx, QueryInsert,
		airlineRoute.ID,
		airlineRoute.RouteID,
		airlineRoute.AirlineID,
		airlineRoute.Status,
	)
	if err != nil {
		log.Error(logger.LogAirlineRouteRepoInitError,
			"id", airlineRoute.ID,
			"route_id", airlineRoute.RouteID,
			"airline_id", airlineRoute.AirlineID,
			"error", err.Error())

		if strings.Contains(err.Error(), "Duplicate entry") {
			return domain.ErrAirlineRouteAlreadyExists
		}
		return domain.ErrAirlineRouteCannotSave
	}

	return nil
}
