package daily_logbook

import (
	"context"
	"database/sql"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
)

// ListDailyLogbooksByEmployee retrieves all daily logbooks for a specific employee
func (r *repository) ListDailyLogbooksByEmployee(ctx context.Context, employeeID string, filters map[string]interface{}) ([]domain.DailyLogbook, error) {
	var rows *sql.Rows
	var err error

	// Check if status filter is provided
	if status, ok := filters["status"]; ok {
		statusBool, isBool := status.(bool)
		if isBool {
			rows, err = r.stmtGetByEmployeeAndStatus.QueryContext(ctx, employeeID, statusBool)
		} else {
			rows, err = r.stmtGetByEmployee.QueryContext(ctx, employeeID)
		}
	} else {
		rows, err = r.stmtGetByEmployee.QueryContext(ctx, employeeID)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logbooks []domain.DailyLogbook
	for rows.Next() {
		var d DailyLogbook
		if err := rows.Scan(&d.ID, &d.LogDate, &d.EmployeeID, &d.BookPage, &d.Status); err != nil {
			return nil, err
		}
		logbooks = append(logbooks, *d.ToDomain())
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return logbooks, nil
}
