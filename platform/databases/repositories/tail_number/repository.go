package tailnumber

import (
	"context"
	"database/sql"

	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/databases/common"
	"github.com/champion19/api-flighthours/platform/logger"
)

const (
	// Query with JOINs to get denormalized data (Numero, Modelo, Aerolinea)
	QueryByID = `
		SELECT
			ar.id,
			ar.tail_number,
			ar.aircraft_model_id,
			ar.airline_id,
			COALESCE(am.model_name, '') AS model_name,
			COALESCE(a.airline_name, '') AS airline_name
		FROM tail_number ar
		LEFT JOIN aircraft_model am ON ar.aircraft_model_id = am.id
		LEFT JOIN airline a ON ar.airline_id = a.id
		WHERE ar.id = ?
		LIMIT 1`

	QueryGetAll = `
		SELECT
			ar.id,
			ar.tail_number,
			ar.aircraft_model_id,
			ar.airline_id,
			COALESCE(am.model_name, '') AS model_name,
			COALESCE(a.airline_name, '') AS airline_name
		FROM tail_number ar
		LEFT JOIN aircraft_model am ON ar.aircraft_model_id = am.id
		LEFT JOIN airline a ON ar.airline_id = a.id
		ORDER BY ar.tail_number`

	QueryGetByAirline = `
		SELECT
			ar.id,
			ar.tail_number,
			ar.aircraft_model_id,
			ar.airline_id,
			COALESCE(am.model_name, '') AS model_name,
			COALESCE(a.airline_name, '') AS airline_name
		FROM tail_number ar
		LEFT JOIN aircraft_model am ON ar.aircraft_model_id = am.id
		LEFT JOIN airline a ON ar.airline_id = a.id
		WHERE ar.airline_id = ?
		ORDER BY ar.tail_number`

	QueryGetByTailNumber = `
		SELECT
			ar.id,
			ar.tail_number,
			ar.aircraft_model_id,
			ar.airline_id,
			COALESCE(am.model_name, '') AS model_name,
			COALESCE(a.airline_name, '') AS airline_name
		FROM tail_number ar
		LEFT JOIN aircraft_model am ON ar.aircraft_model_id = am.id
		LEFT JOIN airline a ON ar.airline_id = a.id
		WHERE ar.tail_number = ?
		LIMIT 1`

	QueryInsert = `INSERT INTO tail_number (id, tail_number, aircraft_model_id, airline_id) VALUES (?, ?, ?, ?)`
	QueryUpdate = `UPDATE tail_number SET tail_number = ?, aircraft_model_id = ?, airline_id = ? WHERE id = ?`
)

var log logger.Logger = logger.NewSlogLogger()

const errPreparingStatement = "error preparing statement"

type repository struct {
	stmtGetByID         *sql.Stmt
	stmtGetAll          *sql.Stmt
	stmtGetByAirline    *sql.Stmt
	stmtGetByTailNumber *sql.Stmt
	stmtInsert          *sql.Stmt
	stmtUpdate          *sql.Stmt
	db                  *sql.DB
}

// NewTailNumberRepository creates a new aircraft registration repository with prepared statements
func NewTailNumberRepository(db *sql.DB) (*repository, error) {
	if db == nil {
		return nil, sql.ErrConnDone
	}

	stmtGetByID, err := db.Prepare(QueryByID)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, errPreparingStatement, err)
		return nil, err
	}

	stmtGetAll, err := db.Prepare(QueryGetAll)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, errPreparingStatement, err)
		return nil, err
	}

	stmtGetByAirline, err := db.Prepare(QueryGetByAirline)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, errPreparingStatement, err)
		return nil, err
	}

	stmtGetByTailNumber, err := db.Prepare(QueryGetByTailNumber)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, errPreparingStatement, err)
		return nil, err
	}

	stmtInsert, err := db.Prepare(QueryInsert)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, errPreparingStatement, err)
		return nil, err
	}

	stmtUpdate, err := db.Prepare(QueryUpdate)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, errPreparingStatement, err)
		return nil, err
	}

	return &repository{
		db:                  db,
		stmtGetByID:         stmtGetByID,
		stmtGetAll:          stmtGetAll,
		stmtGetByAirline:    stmtGetByAirline,
		stmtGetByTailNumber: stmtGetByTailNumber,
		stmtInsert:          stmtInsert,
		stmtUpdate:          stmtUpdate,
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
