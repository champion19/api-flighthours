package airline

import (
	"context"
	"database/sql"

	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/databases/common"
	"github.com/champion19/api-flighthours/platform/logger"
)

const (
	QueryByID        = "SELECT id, airline_name, airline_code, status FROM airline WHERE id = ? LIMIT 1"
	QueryGetAll      = "SELECT id, airline_name, airline_code, status FROM airline ORDER BY airline_name"
	QueryGetByStatus = "SELECT id, airline_name, airline_code, status FROM airline WHERE status = ? ORDER BY airline_name"
)

var log logger.Logger = logger.NewSlogLogger()

type repository struct {
	db              *sql.DB
	stmtGetByID     *sql.Stmt
	stmtGetAll      *sql.Stmt
	stmtGetByStatus *sql.Stmt
}

// NewAirlineRepository creates a new airline repository instance
func NewAirlineRepository(db *sql.DB) (*repository, error) {
	if db == nil {
		return nil, sql.ErrConnDone
	}

	stmtGetByID, err := db.Prepare(QueryByID)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing airline statement", err)
		return nil, err
	}

	stmtGetAll, err := db.Prepare(QueryGetAll)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing airline list statement", err)
		return nil, err
	}

	stmtGetByStatus, err := db.Prepare(QueryGetByStatus)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing airline status statement", err)
		return nil, err
	}

	return &repository{
		db:              db,
		stmtGetByID:     stmtGetByID,
		stmtGetAll:      stmtGetAll,
		stmtGetByStatus: stmtGetByStatus,
	}, nil
}

func (r *repository) BeginTx(ctx context.Context) (output.Tx, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	return common.NewSQLTx(tx), nil
}
