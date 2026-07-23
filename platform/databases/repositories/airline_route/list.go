package airline_route

import (
	"context"
	"database/sql"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
)

func (r *repository) ListAirlineRoutes(ctx context.Context, filters map[string]interface{}) ([]domain.AirlineRoute, error) {
	var rows *sql.Rows
	var err error

	// Check filters
	if airlineID, ok := filters["airline_id"]; ok {
		rows, err = r.stmtGetByAirlineID.QueryContext(ctx, airlineID)
	} else if airlineCode, ok := filters["airline_code"]; ok {
		rows, err = r.stmtGetByAirlineCode.QueryContext(ctx, airlineCode)
	} else if status, ok := filters["status"]; ok {
		rows, err = r.stmtGetByStatus.QueryContext(ctx, status)
	} else {
		rows, err = r.stmtGetAll.QueryContext(ctx)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var airlineRoutes []domain.AirlineRoute
	for rows.Next() {
		var ar AirlineRoute
		var estimatedFlightTime sql.NullString
		// iata_code/oaci_code/route_code can be NULL when an airport was
		// registered with only one of the two codes.
		var originIataCode, originOaciCode sql.NullString
		var destinationIataCode, destinationOaciCode sql.NullString
		var routeCode sql.NullString

		if err := rows.Scan(
			&ar.ID,
			&ar.RouteID,
			&ar.AirlineID,
			&ar.Status,
			&ar.AirlineCode,
			&ar.AirlineName,
			&ar.OriginAirportID,
			&originIataCode,
			&originOaciCode,
			&ar.DestinationAirportID,
			&destinationIataCode,
			&destinationOaciCode,
			&routeCode,
			&ar.OriginAirportName,
			&ar.DestinationAirportName,
			&ar.AirportType,
			&estimatedFlightTime,
		); err != nil {
			return nil, err
		}

		if estimatedFlightTime.Valid {
			ar.EstimatedFlightTime = estimatedFlightTime.String
		}
		if originIataCode.Valid {
			ar.OriginIataCode = originIataCode.String
		}
		if originOaciCode.Valid {
			ar.OriginOaciCode = originOaciCode.String
		}
		if destinationIataCode.Valid {
			ar.DestinationIataCode = destinationIataCode.String
		}
		if destinationOaciCode.Valid {
			ar.DestinationOaciCode = destinationOaciCode.String
		}
		if routeCode.Valid {
			ar.RouteCode = routeCode.String
		}

		airlineRoutes = append(airlineRoutes, *ar.ToDomain())
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return airlineRoutes, nil
}
