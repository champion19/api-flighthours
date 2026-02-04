package route

import (
	"database/sql"

	"github.com/champion19/api-flighthours/platform/logger"
)

const (
	// Query with JOINs to get denormalized data from airport table
	// estimated_flight_time is returned as TIME format "HH:MM:SS"
	QueryByID = `
		SELECT
			r.id,
			r.origin_airport_id,
			ao.iata_code AS origin_iata_code,
			ao.name AS origin_airport_name,
			r.destination_airport_id,
			ad.iata_code AS destination_iata_code,
			ad.name AS destination_airport_name,
			r.origin_country,
			r.destination_country,
			r.airport_type,
			r.estimated_flight_time,
			CONCAT(ao.iata_code, '-', ad.iata_code) AS route_code
		FROM route r
		JOIN airport ao ON r.origin_airport_id = ao.id
		JOIN airport ad ON r.destination_airport_id = ad.id
		WHERE r.id = ?
		LIMIT 1
	`

	QueryGetAll = `
		SELECT
			r.id,
			r.origin_airport_id,
			ao.iata_code AS origin_iata_code,
			ao.name AS origin_airport_name,
			r.destination_airport_id,
			ad.iata_code AS destination_iata_code,
			ad.name AS destination_airport_name,
			r.origin_country,
			r.destination_country,
			r.airport_type,
			r.estimated_flight_time,
			CONCAT(ao.iata_code, '-', ad.iata_code) AS route_code
		FROM route r
		JOIN airport ao ON r.origin_airport_id = ao.id
		JOIN airport ad ON r.destination_airport_id = ad.id
		ORDER BY ao.iata_code, ad.iata_code
	`

	QueryGetByAirportType = `
		SELECT
			r.id,
			r.origin_airport_id,
			ao.iata_code AS origin_iata_code,
			ao.name AS origin_airport_name,
			r.destination_airport_id,
			ad.iata_code AS destination_iata_code,
			ad.name AS destination_airport_name,
			r.origin_country,
			r.destination_country,
			r.airport_type,
			r.estimated_flight_time,
			CONCAT(ao.iata_code, '-', ad.iata_code) AS route_code
		FROM route r
		JOIN airport ao ON r.origin_airport_id = ao.id
		JOIN airport ad ON r.destination_airport_id = ad.id
		WHERE r.airport_type = ?
		ORDER BY ao.iata_code, ad.iata_code
	`

	QueryGetByOriginCountry = `
		SELECT
			r.id,
			r.origin_airport_id,
			ao.iata_code AS origin_iata_code,
			ao.name AS origin_airport_name,
			r.destination_airport_id,
			ad.iata_code AS destination_iata_code,
			ad.name AS destination_airport_name,
			r.origin_country,
			r.destination_country,
			r.airport_type,
			r.estimated_flight_time,
			CONCAT(ao.iata_code, '-', ad.iata_code) AS route_code
		FROM route r
		JOIN airport ao ON r.origin_airport_id = ao.id
		JOIN airport ad ON r.destination_airport_id = ad.id
		WHERE r.origin_country = ?
		ORDER BY ao.iata_code, ad.iata_code
	`
)

var log logger.Logger = logger.NewSlogLogger()

type repository struct {
	stmtGetByID            *sql.Stmt
	stmtGetAll             *sql.Stmt
	stmtGetByAirportType   *sql.Stmt
	stmtGetByOriginCountry *sql.Stmt
	db                     *sql.DB
}

// NewRouteRepository creates a new route repository with prepared statements
func NewRouteRepository(db *sql.DB) (*repository, error) {
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

	stmtGetByAirportType, err := db.Prepare(QueryGetByAirportType)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing statement", err)
		return nil, err
	}

	stmtGetByOriginCountry, err := db.Prepare(QueryGetByOriginCountry)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing statement", err)
		return nil, err
	}

	return &repository{
		db:                     db,
		stmtGetByID:            stmtGetByID,
		stmtGetAll:             stmtGetAll,
		stmtGetByAirportType:   stmtGetByAirportType,
		stmtGetByOriginCountry: stmtGetByOriginCountry,
	}, nil
}
