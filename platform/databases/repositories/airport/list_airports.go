package airport

import(
	"context"
	"database/sql"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
)



// ListAirports retrieves all airports, optionally filtered by status
func (r *repository) ListAirports(ctx context.Context, filters map[string]interface{}) ([]domain.Airport, error) {
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

	var airports []domain.Airport
	for rows.Next() {
		var a Airport
		if err := rows.Scan(&a.ID, &a.Name, &a.City, &a.Country, &a.IATACode, &a.Status, &a.AirportType); err != nil {
			return nil, err
		}
		airports = append(airports, *a.ToDomain())
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return airports, nil
}
