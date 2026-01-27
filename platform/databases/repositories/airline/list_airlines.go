package airline

import (
	"context"
	"database/sql"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
)

// ListAirlines retrieves all airlines, optionally filtered by status
func (r *repository) ListAirlines(ctx context.Context, filters map[string]interface{}) ([]domain.Airline, error) {
	var rows *sql.Rows
	var err error

	// Check if status filter is provided
	if status, ok := filters["status"]; ok {
		statusBool, isBool := status.(bool)
		if isBool {
			rows, err = r.stmtGetByStatus.QueryContext(ctx, statusBool)
		} else {
			rows, err = r.stmtGetAll.QueryContext(ctx)
		}
	} else {
		rows, err = r.stmtGetAll.QueryContext(ctx)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var airlines []domain.Airline
	for rows.Next() {
		var a Airline
		if err := rows.Scan(&a.ID, &a.AirlineName, &a.AirlineCode, &a.Status); err != nil {
			return nil, err
		}
		airlines = append(airlines, *a.ToDomain())
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return airlines, nil
}
