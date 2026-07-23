package route

import (
	"context"
	"database/sql"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
)

// ListRoutes retrieves all routes with optional filters
func (r *repository) ListRoutes(ctx context.Context, filters map[string]interface{}) ([]domain.Route, error) {
	var rows *sql.Rows
	var err error

	// Check if filtering by airport type
	if airportType, ok := filters["airport_type"]; ok {
		rows, err = r.stmtGetByAirportType.QueryContext(ctx, airportType)
	} else {
		rows, err = r.stmtGetAll.QueryContext(ctx)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var routes []domain.Route
	for rows.Next() {
		var route Route
		var estimatedFlightTime sql.NullString
		var originCountry, destinationCountry sql.NullString
		// iata_code/oaci_code/route_code can be NULL when an airport was
		// registered with only one of the two codes.
		var originIataCode, originOaciCode sql.NullString
		var destinationIataCode, destinationOaciCode sql.NullString
		var routeCode sql.NullString

		if err := rows.Scan(
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
		); err != nil {
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

		routes = append(routes, *route.ToDomain())
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return routes, nil
}
