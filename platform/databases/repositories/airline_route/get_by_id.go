package airline_route

import (
	"context"
	"database/sql"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
)


func (r *repository) GetAirlineRouteByID(ctx context.Context, id string) (*domain.AirlineRoute, error) {
	var airlineRoute AirlineRoute
	var estimatedFlightTime sql.NullString

	err := r.stmtGetByID.QueryRowContext(ctx, id).Scan(
		&airlineRoute.ID,
		&airlineRoute.RouteID,
		&airlineRoute.AirlineID,
		&airlineRoute.Status,
		&airlineRoute.AirlineCode,
		&airlineRoute.AirlineName,
		&airlineRoute.OriginIataCode,
		&airlineRoute.DestinationIataCode,
		&airlineRoute.RouteCode,
		&airlineRoute.OriginAirportName,
		&airlineRoute.DestinationAirportName,
		&airlineRoute.AirportType,
		&estimatedFlightTime,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrAirlineRouteNotFound
		}
		return nil, err
	}

	if estimatedFlightTime.Valid {
		airlineRoute.EstimatedFlightTime = estimatedFlightTime.String
	}

	return airlineRoute.ToDomain(), nil
}
