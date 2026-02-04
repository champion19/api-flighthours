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
	} else if originCountry, ok := filters["origin_country"]; ok {
		rows, err = r.stmtGetByOriginCountry.QueryContext(ctx, originCountry)
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
		var originCountry, destinationCountry, estimatedFlightTime sql.NullString

		if err := rows.Scan(
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
		); err != nil {
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

		routes = append(routes, *route.ToDomain())
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return routes, nil
}
