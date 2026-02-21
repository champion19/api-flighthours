package airline_route

import (
	"context"
	"database/sql"

	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/logger"
)

const (
	QueryByID = `
		SELECT
			ar.id,
			ar.route_id,
			ar.airline_id,
			ar.status,
			a.airline_code,
			a.airline_name,
			ao.iata_code AS origin_iata_code,
			ad.iata_code AS destination_iata_code,
			CONCAT(ao.iata_code, '-', ad.iata_code) AS route_code,
			ao.name AS origin_airport_name,
			ad.name AS destination_airport_name,
			r.airport_type,
			r.estimated_flight_time
		FROM airline_route ar
		JOIN airline a ON ar.airline_id = a.id
		JOIN route r ON ar.route_id = r.id
		JOIN airport ao ON r.origin_airport_id = ao.id
		JOIN airport ad ON r.destination_airport_id = ad.id
		WHERE ar.id = ? AND ao.status = TRUE AND ad.status = TRUE AND a.status = TRUE
		LIMIT 1
	`

	QueryGetAll = `
		SELECT
			ar.id,
			ar.route_id,
			ar.airline_id,
			ar.status,
			a.airline_code,
			a.airline_name,
			ao.iata_code AS origin_iata_code,
			ad.iata_code AS destination_iata_code,
			CONCAT(ao.iata_code, '-', ad.iata_code) AS route_code,
			ao.name AS origin_airport_name,
			ad.name AS destination_airport_name,
			r.airport_type,
			r.estimated_flight_time
		FROM airline_route ar
		JOIN airline a ON ar.airline_id = a.id
		JOIN route r ON ar.route_id = r.id
		JOIN airport ao ON r.origin_airport_id = ao.id
		JOIN airport ad ON r.destination_airport_id = ad.id
		WHERE ao.status = TRUE AND ad.status = TRUE AND a.status = TRUE
		ORDER BY a.airline_code, ao.iata_code, ad.iata_code
	`

	QueryGetByAirlineID = `
		SELECT
			ar.id,
			ar.route_id,
			ar.airline_id,
			ar.status,
			a.airline_code,
			a.airline_name,
			ao.iata_code AS origin_iata_code,
			ad.iata_code AS destination_iata_code,
			CONCAT(ao.iata_code, '-', ad.iata_code) AS route_code,
			ao.name AS origin_airport_name,
			ad.name AS destination_airport_name,
			r.airport_type,
			r.estimated_flight_time
		FROM airline_route ar
		JOIN airline a ON ar.airline_id = a.id
		JOIN route r ON ar.route_id = r.id
		JOIN airport ao ON r.origin_airport_id = ao.id
		JOIN airport ad ON r.destination_airport_id = ad.id
		WHERE ar.airline_id = ? AND ar.status = true AND ao.status = TRUE AND ad.status = TRUE AND a.status = TRUE
		ORDER BY ao.iata_code, ad.iata_code
	`

	QueryGetByAirlineCode = `
		SELECT
			ar.id,
			ar.route_id,
			ar.airline_id,
			ar.status,
			a.airline_code,
			a.airline_name,
			ao.iata_code AS origin_iata_code,
			ad.iata_code AS destination_iata_code,
			CONCAT(ao.iata_code, '-', ad.iata_code) AS route_code,
			ao.name AS origin_airport_name,
			ad.name AS destination_airport_name,
			r.airport_type,
			r.estimated_flight_time
		FROM airline_route ar
		JOIN airline a ON ar.airline_id = a.id
		JOIN route r ON ar.route_id = r.id
		JOIN airport ao ON r.origin_airport_id = ao.id
		JOIN airport ad ON r.destination_airport_id = ad.id
		WHERE a.airline_code = ? AND ao.status = TRUE AND ad.status = TRUE AND a.status = TRUE
		ORDER BY ao.iata_code, ad.iata_code
	`

	QueryGetByStatus = `
		SELECT
			ar.id,
			ar.route_id,
			ar.airline_id,
			ar.status,
			a.airline_code,
			a.airline_name,
			ao.iata_code AS origin_iata_code,
			ad.iata_code AS destination_iata_code,
			CONCAT(ao.iata_code, '-', ad.iata_code) AS route_code,
			ao.name AS origin_airport_name,
			ad.name AS destination_airport_name,
			r.airport_type,
			r.estimated_flight_time
		FROM airline_route ar
		JOIN airline a ON ar.airline_id = a.id
		JOIN route r ON ar.route_id = r.id
		JOIN airport ao ON r.origin_airport_id = ao.id
		JOIN airport ad ON r.destination_airport_id = ad.id
		WHERE ar.status = ? AND ao.status = TRUE AND ad.status = TRUE AND a.status = TRUE
		ORDER BY a.airline_code, ao.iata_code, ad.iata_code
	`

	QueryUpdateStatus = `
		UPDATE airline_route
		SET status = ?
		WHERE id = ? AND status != ?
	`
)

var log logger.Logger = logger.NewSlogLogger()

const errPreparingStatement = "error preparing statement"

type repository struct {
	stmtGetByID          *sql.Stmt
	stmtGetAll           *sql.Stmt
	stmtGetByAirlineID   *sql.Stmt
	stmtGetByAirlineCode *sql.Stmt
	stmtGetByStatus      *sql.Stmt
	stmtUpdateStatus     *sql.Stmt
	db                   *sql.DB
}

func NewAirlineRouteRepository(db *sql.DB) (*repository, error) {
	if db == nil {
		return nil, sql.ErrConnDone
	}

	stmtGetByID, err := db.Prepare(QueryByID)
	if err != nil {
		log.Error(logger.LogAirlineRouteRepoInitError, errPreparingStatement, err)
		return nil, err
	}

	stmtGetAll, err := db.Prepare(QueryGetAll)
	if err != nil {
		log.Error(logger.LogAirlineRouteRepoInitError, errPreparingStatement, err)
		return nil, err
	}

	stmtGetByAirlineID, err := db.Prepare(QueryGetByAirlineID)
	if err != nil {
		log.Error(logger.LogAirlineRouteRepoInitError, errPreparingStatement, err)
		return nil, err
	}

	stmtGetByAirlineCode, err := db.Prepare(QueryGetByAirlineCode)
	if err != nil {
		log.Error(logger.LogAirlineRouteRepoInitError, errPreparingStatement, err)
		return nil, err
	}

	stmtGetByStatus, err := db.Prepare(QueryGetByStatus)
	if err != nil {
		log.Error(logger.LogAirlineRouteRepoInitError, errPreparingStatement, err)
		return nil, err
	}

	stmtUpdateStatus, err := db.Prepare(QueryUpdateStatus)
	if err != nil {
		log.Error(logger.LogAirlineRouteRepoInitError, errPreparingStatement, err)
		return nil, err
	}

	log.Info(logger.LogAirlineRouteRepoInitOK)

	return &repository{
		db:                   db,
		stmtGetByID:          stmtGetByID,
		stmtGetAll:           stmtGetAll,
		stmtGetByAirlineID:   stmtGetByAirlineID,
		stmtGetByAirlineCode: stmtGetByAirlineCode,
		stmtGetByStatus:      stmtGetByStatus,
		stmtUpdateStatus:     stmtUpdateStatus,
	}, nil
}

func (r *repository) BeginTx(ctx context.Context) (output.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}
