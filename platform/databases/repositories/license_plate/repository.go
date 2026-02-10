package licenseplate

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
			ar.license_plate,
			ar.aircraft_model_id,
			ar.airline_id,
			COALESCE(am.model_name, '') AS model_name,
			COALESCE(a.airline_name, '') AS airline_name
		FROM license_plate ar
		LEFT JOIN aircraft_model am ON ar.aircraft_model_id = am.id
		LEFT JOIN airline a ON ar.airline_id = a.id
		WHERE ar.id = ?
		LIMIT 1`

	QueryGetAll = `
		SELECT
			ar.id,
			ar.license_plate,
			ar.aircraft_model_id,
			ar.airline_id,
			COALESCE(am.model_name, '') AS model_name,
			COALESCE(a.airline_name, '') AS airline_name
		FROM license_plate ar
		LEFT JOIN aircraft_model am ON ar.aircraft_model_id = am.id
		LEFT JOIN airline a ON ar.airline_id = a.id
		ORDER BY ar.license_plate`

	QueryGetByAirline = `
		SELECT
			ar.id,
			ar.license_plate,
			ar.aircraft_model_id,
			ar.airline_id,
			COALESCE(am.model_name, '') AS model_name,
			COALESCE(a.airline_name, '') AS airline_name
		FROM license_plate ar
		LEFT JOIN aircraft_model am ON ar.aircraft_model_id = am.id
		LEFT JOIN airline a ON ar.airline_id = a.id
		WHERE ar.airline_id = ?
		ORDER BY ar.license_plate`

	QueryGetByLicensePlate = `
		SELECT
			ar.id,
			ar.license_plate,
			ar.aircraft_model_id,
			ar.airline_id,
			COALESCE(am.model_name, '') AS model_name,
			COALESCE(a.airline_name, '') AS airline_name
		FROM license_plate ar
		LEFT JOIN aircraft_model am ON ar.aircraft_model_id = am.id
		LEFT JOIN airline a ON ar.airline_id = a.id
		WHERE ar.license_plate = ?
		LIMIT 1`

	QueryInsert = `INSERT INTO license_plate (id, license_plate, aircraft_model_id, airline_id) VALUES (?, ?, ?, ?)`
	QueryUpdate = `UPDATE license_plate SET license_plate = ?, aircraft_model_id = ?, airline_id = ? WHERE id = ?`
)

var log logger.Logger = logger.NewSlogLogger()

type repository struct {
	stmtGetByID           *sql.Stmt
	stmtGetAll            *sql.Stmt
	stmtGetByAirline      *sql.Stmt
	stmtGetByLicensePlate *sql.Stmt
	stmtInsert            *sql.Stmt
	stmtUpdate            *sql.Stmt
	db                    *sql.DB
}

// NewLicensePlateRepository creates a new aircraft registration repository with prepared statements
func NewLicensePlateRepository(db *sql.DB) (*repository, error) {
	if db == nil {
		return nil, sql.ErrConnDone
	}

	stmtGetByID, err := db.Prepare(QueryByID)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing statement", err)
		return nil, err
	}

	stmtGetAll, err := db.Prepare(QueryGetAll)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing statement", err)
		return nil, err
	}

	stmtGetByAirline, err := db.Prepare(QueryGetByAirline)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing statement", err)
		return nil, err
	}

	stmtGetByLicensePlate, err := db.Prepare(QueryGetByLicensePlate)
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

	return &repository{
		db:                    db,
		stmtGetByID:           stmtGetByID,
		stmtGetAll:            stmtGetAll,
		stmtGetByAirline:      stmtGetByAirline,
		stmtGetByLicensePlate: stmtGetByLicensePlate,
		stmtInsert:            stmtInsert,
		stmtUpdate:            stmtUpdate,
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
