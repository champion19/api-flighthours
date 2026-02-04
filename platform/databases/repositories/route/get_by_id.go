package route

import (
	"context"
	"database/sql"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
)

// GetRouteByID retrieves a route by ID with denormalized airport data
func (r *repository) GetRouteByID(ctx context.Context, id string) (*domain.Route, error) {
	var route Route
	var originCountry, destinationCountry, estimatedFlightTime sql.NullString

	err := r.stmtGetByID.QueryRowContext(ctx, id).Scan(
		&route.ID,
		&route.OriginAirportID,
		&route.OriginIataCode,
		&route.OriginAirportName,
		&route.DestinationAirportID,
		&route.DestinationIataCode,
		&route.DestinationAirportName,
		&originCountry,
		&destinationCountry,
		&route.AirportType,
		&estimatedFlightTime,
		&route.RouteCode,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrRouteNotFound
		}
		return nil, err
	}

	if originCountry.Valid {
		route.OriginCountry = originCountry.String
	}
	if destinationCountry.Valid {
		route.DestinationCountry = destinationCountry.String
	}
	if estimatedFlightTime.Valid {
		route.EstimatedFlightTime = estimatedFlightTime.String
	}

	return route.ToDomain(), nil
}
