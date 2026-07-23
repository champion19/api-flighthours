package route

import (
	"context"
	"database/sql"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
)

// GetRouteByAirports looks up the physical route (if any) between two
// airports by their ID — used to resolve/auto-request an airline_route link
// without depending on IATA/OACI codes, which can be missing.
func (r *repository) GetRouteByAirports(ctx context.Context, originAirportID, destinationAirportID string) (*domain.Route, error) {
	var route Route
	var estimatedFlightTime sql.NullString
	var originCountry, destinationCountry sql.NullString
	var originIataCode, originOaciCode sql.NullString
	var destinationIataCode, destinationOaciCode sql.NullString
	var routeCode sql.NullString

	err := r.stmtGetByAirports.QueryRowContext(ctx, originAirportID, destinationAirportID).Scan(
		&route.ID,
		&route.OriginAirportID,
		&originIataCode,
		&originOaciCode,
		&route.OriginAirportName,
		&originCountry,
		&route.DestinationAirportID,
		&destinationIataCode,
		&destinationOaciCode,
		&route.DestinationAirportName,
		&destinationCountry,
		&route.AirportType,
		&estimatedFlightTime,
		&routeCode,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrRouteNotFound
		}
		return nil, err
	}

	if estimatedFlightTime.Valid {
		route.EstimatedFlightTime = estimatedFlightTime.String
	}
	if originCountry.Valid {
		route.OriginCountry = originCountry.String
	}
	if destinationCountry.Valid {
		route.DestinationCountry = destinationCountry.String
	}
	if originIataCode.Valid {
		route.OriginIataCode = originIataCode.String
	}
	if originOaciCode.Valid {
		route.OriginOaciCode = originOaciCode.String
	}
	if destinationIataCode.Valid {
		route.DestinationIataCode = destinationIataCode.String
	}
	if destinationOaciCode.Valid {
		route.DestinationOaciCode = destinationOaciCode.String
	}
	if routeCode.Valid {
		route.RouteCode = routeCode.String
	}

	return route.ToDomain(), nil
}
