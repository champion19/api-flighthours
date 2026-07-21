package route

import (
	"database/sql"

	"github.com/champion19/api-flighthours/platform/logger"
)

const (
	// Query with JOINs to get denormalized data from airport table
	QueryByID = `
		SELECT
			r.id,
			r.origin_airport_id,
			ao.iata_code AS origin_iata_code,
			ao.name AS origin_airport_name,
			ao.country AS origin_country,
			r.destination_airport_id,
			ad.iata_code AS destination_iata_code,
			ad.name AS destination_airport_name,
			ad.country AS destination_country,
			r.airport_type,
			r.estimated_flight_time,
			CONCAT(ao.iata_code, '-', ad.iata_code) AS route_code
		FROM route r
		JOIN airport ao ON r.origin_airport_id = ao.id
		JOIN airport ad ON r.destination_airport_id = ad.id
		WHERE r.id = ? AND ao.status = TRUE AND ad.status = TRUE
		LIMIT 1
	`

	QueryGetAll = `
		SELECT
			r.id,
			r.origin_airport_id,
			ao.iata_code AS origin_iata_code,
			ao.name AS origin_airport_name,
			ao.country AS origin_country,
			r.destination_airport_id,
			ad.iata_code AS destination_iata_code,
			ad.name AS destination_airport_name,
			ad.country AS destination_country,
			r.airport_type,
			r.estimated_flight_time,
			CONCAT(ao.iata_code, '-', ad.iata_code) AS route_code
		FROM route r
		JOIN airport ao ON r.origin_airport_id = ao.id
		JOIN airport ad ON r.destination_airport_id = ad.id
		WHERE ao.status = TRUE AND ad.status = TRUE
		ORDER BY ao.iata_code, ad.iata_code
	`

	QueryGetByAirportType = `
		SELECT
			r.id,
			r.origin_airport_id,
			ao.iata_code AS origin_iata_code,
			ao.name AS origin_airport_name,
			ao.country AS origin_country,
			r.destination_airport_id,
			ad.iata_code AS destination_iata_code,
			ad.name AS destination_airport_name,
			ad.country AS destination_country,
			r.airport_type,
			r.estimated_flight_time,
			CONCAT(ao.iata_code, '-', ad.iata_code) AS route_code
		FROM route r
		JOIN airport ao ON r.origin_airport_id = ao.id
		JOIN airport ad ON r.destination_airport_id = ad.id
		WHERE r.airport_type = ? AND ao.status = TRUE AND ad.status = TRUE
		ORDER BY ao.iata_code, ad.iata_code
	`
)

var log logger.Logger = logger.NewSlogLogger()

const errPreparingStatement = "error preparing statement"

type repository struct {
	stmtGetByID          *sql.Stmt
	stmtGetAll           *sql.Stmt
	stmtGetByAirportType *sql.Stmt
	db                   *sql.DB
}

// NewRouteRepository creates a new route repository with prepared statements
func NewRouteRepository(db *sql.DB) (*repository, error) {
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

	stmtGetByAirportType, err := db.Prepare(QueryGetByAirportType)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, errPreparingStatement, err)
		return nil, err
	}

	return &repository{
		db:                   db,
		stmtGetByID:          stmtGetByID,
		stmtGetAll:           stmtGetAll,
		stmtGetByAirportType: stmtGetByAirportType,
	}, nil
}
