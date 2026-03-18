package airport

import (
	"context"
	"database/sql"

	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/databases/common"
	"github.com/champion19/api-flighthours/platform/logger"
)

const (
	QueryByID         = "SELECT id, name, iata_code, status, airport_type FROM airport WHERE id = ? LIMIT 1"
	QueryUpdateStatus = "UPDATE airport SET status = ? WHERE id = ?"
	QueryGetAll       = "SELECT id, name, iata_code, status, airport_type FROM airport ORDER BY name"
	QueryGetByStatus  = "SELECT id, name, iata_code, status, airport_type FROM airport WHERE status = ? ORDER BY name"
	QueryGetByType    = "SELECT id, name, iata_code, status, airport_type FROM airport WHERE airport_type = ? ORDER BY name"
)

var log logger.Logger = logger.NewSlogLogger()

const errPreparingStatement = "error preparing statement"

type repository struct {
	stmtGetByID      *sql.Stmt
	stmtUpdateStatus *sql.Stmt
	stmtGetAll       *sql.Stmt
	stmtGetByStatus  *sql.Stmt
	stmtGetByType    *sql.Stmt
	db               *sql.DB
}

func NewAirportRepository(db *sql.DB) (*repository, error) {
	if db == nil {
		return nil, sql.ErrConnDone
	}

	stmtGetByID, err := db.Prepare(QueryByID)
	if err != nil {
		log.Error(logger.LogAirportRepoInitError, errPreparingStatement, err)
		return nil, err
	}

	stmtUpdateStatus, err := db.Prepare(QueryUpdateStatus)
	if err != nil {
		log.Error(logger.LogAirportRepoInitError, errPreparingStatement, err)
		return nil, err
	}

	stmtGetAll, err := db.Prepare(QueryGetAll)
	if err != nil {
		log.Error(logger.LogAirportRepoInitError, errPreparingStatement, err)
		return nil, err
	}

	stmtGetByStatus, err := db.Prepare(QueryGetByStatus)
	if err != nil {
		log.Error(logger.LogAirportRepoInitError, errPreparingStatement, err)
		return nil, err
	}

	stmtGetByType, err := db.Prepare(QueryGetByType)
	if err != nil {
		log.Error(logger.LogAirportRepoInitError, "error preparing airport type statement", err)
		return nil, err
	}

	return &repository{
		db:               db,
		stmtGetByID:      stmtGetByID,
		stmtUpdateStatus: stmtUpdateStatus,
		stmtGetAll:       stmtGetAll,
		stmtGetByStatus:  stmtGetByStatus,
		stmtGetByType:    stmtGetByType,
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
