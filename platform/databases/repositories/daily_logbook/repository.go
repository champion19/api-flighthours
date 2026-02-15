package daily_logbook

import (
	"context"
	"database/sql"

	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/databases/common"
	"github.com/champion19/api-flighthours/platform/logger"
)

const (
	QueryByID                = "SELECT id, log_date, employee_id, book_page, status FROM daily_logbook WHERE id = ? LIMIT 1"
	QueryByEmployee          = "SELECT id, log_date, employee_id, book_page, status FROM daily_logbook WHERE employee_id = ? ORDER BY log_date DESC"
	QueryByEmployeeAndStatus = "SELECT id, log_date, employee_id, book_page, status FROM daily_logbook WHERE employee_id = ? AND status = ? ORDER BY log_date DESC"
	QueryInsert              = "INSERT INTO daily_logbook (id, log_date, employee_id, book_page, status) VALUES (?, ?, ?, ?, ?)"
	QueryUpdate              = "UPDATE daily_logbook SET log_date = ?, book_page = ?, status = ? WHERE id = ?"
	QueryDeleteDetails       = "DELETE FROM daily_logbook_detail WHERE daily_logbook_id = ?"
	QueryDelete              = "DELETE FROM daily_logbook WHERE id = ?"
	QueryUpdateStatus        = "UPDATE daily_logbook SET status = ? WHERE id = ?"
)

var log logger.Logger = logger.NewSlogLogger()

type repository struct {
	stmtGetByID                *sql.Stmt
	stmtGetByEmployee          *sql.Stmt
	stmtGetByEmployeeAndStatus *sql.Stmt
	stmtInsert                 *sql.Stmt
	stmtUpdate                 *sql.Stmt
	stmtDelete                 *sql.Stmt
	stmtUpdateStatus           *sql.Stmt
	db                         *sql.DB
}

// NewDailyLogbookRepository creates a new daily logbook repository with prepared statements
func NewDailyLogbookRepository(db *sql.DB) (*repository, error) {
	if db == nil {
		return nil, sql.ErrConnDone
	}

	stmtGetByID, err := db.Prepare(QueryByID)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing statement", err)
		return nil, err
	}

	stmtGetByEmployee, err := db.Prepare(QueryByEmployee)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing statement", err)
		return nil, err
	}

	stmtGetByEmployeeAndStatus, err := db.Prepare(QueryByEmployeeAndStatus)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing statement", err)
		return nil, err
	}

	stmtInsert, err := db.Prepare(QueryInsert)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing statement", err)
		return nil, err
	}
	stmtUpdate, err := db.Prepare(QueryUpdate)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing statement", err)
		return nil, err
	}
	stmtDelete, err := db.Prepare(QueryDelete)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing statement", err)
		return nil, err
	}
	stmtUpdateStatus, err := db.Prepare(QueryUpdateStatus)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing statement", err)
		return nil, err
	}

	return &repository{
		db:                         db,
		stmtGetByID:                stmtGetByID,
		stmtGetByEmployee:          stmtGetByEmployee,
		stmtGetByEmployeeAndStatus: stmtGetByEmployeeAndStatus,
		stmtInsert:                 stmtInsert,
		stmtUpdate:                 stmtUpdate,
		stmtDelete:                 stmtDelete,
		stmtUpdateStatus:           stmtUpdateStatus,
	}, nil
}

// BeginTx starts a new database transaction
func (r *repository) BeginTx(ctx context.Context) (output.Tx, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return common.NewSQLTx(tx), nil
}
