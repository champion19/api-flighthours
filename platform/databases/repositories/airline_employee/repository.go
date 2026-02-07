package airline_employee

import (
	"context"
	"database/sql"

	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/databases/common"
	"github.com/champion19/api-flighthours/platform/logger"
)

const (
	// Query to get airline-specific info for an employee
	// Uses LEFT JOIN to handle employees without airline, and filters by airline status when assigned
	QueryByID = `
		SELECT
			e.id,
			e.airline,
			e.bp,
			e.start_date,
			e.end_date,
			e.active
		FROM employee e
		LEFT JOIN airline a ON e.airline = a.id
		WHERE e.id = ? AND (e.airline IS NULL OR a.status = TRUE)
		LIMIT 1
	`

	// Query to update only airline-specific fields for an employee
	QueryUpdateAirlineInfo = `
		UPDATE employee SET
			airline = ?,
			bp = ?,
			start_date = ?,
			end_date = ?,
			active = ?
		WHERE id = ?
	`

	QueryUpdateStatus = `
		UPDATE employee SET active = ?
		WHERE id = ? AND airline IS NOT NULL
	`
)

var log logger.Logger = logger.NewSlogLogger()

type repository struct {
	stmtGetByID      *sql.Stmt
	stmtUpdateStatus *sql.Stmt
	db               *sql.DB
}

func NewAirlineEmployeeRepository(db *sql.DB) (*repository, error) {
	if db == nil {
		return nil, sql.ErrConnDone
	}

	stmtGetByID, err := db.Prepare(QueryByID)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing statement", err)
		return nil, err
	}

	stmtUpdateStatus, err := db.Prepare(QueryUpdateStatus)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing statement", err)
		return nil, err
	}

	log.Info(logger.LogDatabaseAvailable, "repository", "airline_employee")

	return &repository{
		db:               db,
		stmtGetByID:      stmtGetByID,
		stmtUpdateStatus: stmtUpdateStatus,
	}, nil
}

func (r *repository) BeginTx(ctx context.Context) (output.Tx, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return common.NewSQLTx(tx), nil
}
