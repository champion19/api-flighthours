package airport

import (
	"context"
	"database/sql"

	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/databases/common"
	"github.com/champion19/api-flighthours/platform/logger"
)

const (
	QueryByID         = "SELECT id, name, city, country, iata_code, status, airport_type FROM airport WHERE id = ? LIMIT 1"
	QueryUpdateStatus = "UPDATE airport SET status = ? WHERE id = ?"
	QueryGetAll       = "SELECT id, name, city, country, iata_code, status, airport_type FROM airport ORDER BY name"
	QueryGetByStatus  = "SELECT id, name, city, country, iata_code, status, airport_type FROM airport WHERE status = ? ORDER BY name"
	// HU13 - Get airports by city (Virtual Entity pattern - no new table needed)
	QueryGetByCity = "SELECT id, name, city, country, iata_code, status, airport_type FROM airport WHERE city = ? ORDER BY name"
	// HU38 - Get airports by country (Virtual Entity pattern - no new table needed)
	QueryGetByCountry = "SELECT id, name, city, country, iata_code, status, airport_type FROM airport WHERE country = ? ORDER BY name"
	// HU46 - Get airports by type (Virtual Entity pattern - no new table needed)
	QueryGetByType = "SELECT id, name, city, country, iata_code, status, airport_type FROM airport WHERE airport_type = ? ORDER BY name"
)

var log logger.Logger = logger.NewSlogLogger()

type repository struct {
	stmtGetByID      *sql.Stmt
	stmtUpdateStatus *sql.Stmt
	stmtGetAll       *sql.Stmt
	stmtGetByStatus  *sql.Stmt
	stmtGetByCity    *sql.Stmt
	stmtGetByCountry *sql.Stmt
	stmtGetByType    *sql.Stmt
	db               *sql.DB
}

// NewAirportRepository creates a new airport repository with prepared statements
func NewAirportRepository(db *sql.DB) (*repository, error) {
	if db == nil {
		return nil, sql.ErrConnDone
	}

	stmtGetByID, err := db.Prepare(QueryByID)
	if err != nil {
		log.Error(logger.LogAirportRepoInitError, "error preparing statement", err)
		return nil, err
	}

	stmtUpdateStatus, err := db.Prepare(QueryUpdateStatus)
	if err != nil {
		log.Error(logger.LogAirportRepoInitError, "error preparing statement", err)
		return nil, err
	}

	stmtGetAll, err := db.Prepare(QueryGetAll)
	if err != nil {
		log.Error(logger.LogAirportRepoInitError, "error preparing statement", err)
		return nil, err
	}

	stmtGetByStatus, err := db.Prepare(QueryGetByStatus)
	if err != nil {
		log.Error(logger.LogAirportRepoInitError, "error preparing statement", err)
		return nil, err
	}

	// HU13 - Prepare statement for city lookup
	stmtGetByCity, err := db.Prepare(QueryGetByCity)
	if err != nil {
		log.Error(logger.LogAirportRepoInitError, "error preparing city statement", err)
		return nil, err
	}

	// HU38 - Prepare statement for country lookup
	stmtGetByCountry, err := db.Prepare(QueryGetByCountry)
	if err != nil {
		log.Error(logger.LogAirportRepoInitError, "error preparing country statement", err)
		return nil, err
	}

	// HU46 - Prepare statement for airport type lookup
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
		stmtGetByCity:    stmtGetByCity,
		stmtGetByCountry: stmtGetByCountry,
		stmtGetByType:    stmtGetByType,
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
