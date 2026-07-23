package airline_route

import (
	"context"
	"database/sql"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
)

// GetAirlineRouteByRouteAndAirline looks up the airline_route link (if any)
// between a physical route and an airline, regardless of status. Used to
// find-or-create the pending link when an employee picks an origin/
// destination pair whose route exists but isn't yet linked to their airline.
func (r *repository) GetAirlineRouteByRouteAndAirline(ctx context.Context, routeID, airlineID string) (*domain.AirlineRoute, error) {
	var airlineRoute AirlineRoute
	var estimatedFlightTime sql.NullString
	var originIataCode, originOaciCode sql.NullString
	var destinationIataCode, destinationOaciCode sql.NullString
	var routeCode sql.NullString

	err := r.stmtGetByRouteAndAirline.QueryRowContext(ctx, routeID, airlineID).Scan(
		&airlineRoute.ID,
		&airlineRoute.RouteID,
		&airlineRoute.AirlineID,
		&airlineRoute.Status,
		&airlineRoute.AirlineCode,
		&airlineRoute.AirlineName,
		&airlineRoute.OriginAirportID,
		&originIataCode,
		&originOaciCode,
		&airlineRoute.DestinationAirportID,
		&destinationIataCode,
		&destinationOaciCode,
		&routeCode,
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
	if originIataCode.Valid {
		airlineRoute.OriginIataCode = originIataCode.String
	}
	if originOaciCode.Valid {
		airlineRoute.OriginOaciCode = originOaciCode.String
	}
	if destinationIataCode.Valid {
		airlineRoute.DestinationIataCode = destinationIataCode.String
	}
	if destinationOaciCode.Valid {
		airlineRoute.DestinationOaciCode = destinationOaciCode.String
	}
	if routeCode.Valid {
		airlineRoute.RouteCode = routeCode.String
	}

	return airlineRoute.ToDomain(), nil
}
