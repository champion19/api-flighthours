package daily_logbook_detail

import (
	"context"
	"database/sql"

	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/databases/common"
	"github.com/champion19/api-flighthours/platform/logger"
)

const (
	// Query for getting a detail by ID with JOINs for denormalized data
	QueryByID = `
		SELECT
			dld.id,
			dld.daily_logbook_id,
			dld.flight_real_date,
			dld.flight_number,
			dld.airline_route_id,
			dld.actual_license_plate_id,
			dld.passengers,
			dld.out_time,
			dld.takeoff_time,
			dld.landing_time,
			dld.in_time,
			dld.pilot_role,
			dld.companion_name,
			dld.air_time,
			dld.block_time,
			dld.duty_time,
			dld.approach_type,
			dld.flight_type,
			dld.employee_logbook_id,
			dl.log_date,
			lp.license_plate,
			am.model_name,
			CONCAT(orig.iata_code, '-', dest.iata_code) as route_code,
			orig.iata_code as origin_iata_code,
			dest.iata_code as destination_iata_code,
			airl.airline_code
		FROM daily_logbook_detail dld
		INNER JOIN daily_logbook dl ON dld.daily_logbook_id = dl.id
		INNER JOIN license_plate lp ON dld.actual_license_plate_id = lp.id
		INNER JOIN aircraft_model am ON lp.aircraft_model_id = am.id
		INNER JOIN airline_route alr ON dld.airline_route_id = alr.id
		INNER JOIN route r ON alr.route_id = r.id
		INNER JOIN airport orig ON r.origin_airport_id = orig.id
		INNER JOIN airport dest ON r.destination_airport_id = dest.id
		INNER JOIN airline airl ON alr.airline_id = airl.id
		WHERE dld.id = ?
		LIMIT 1
	`

	// Query for listing details by logbook ID
	QueryByLogbook = `
		SELECT
			dld.id,
			dld.daily_logbook_id,
			dld.flight_real_date,
			dld.flight_number,
			dld.airline_route_id,
			dld.actual_license_plate_id,
			dld.passengers,
			dld.out_time,
			dld.takeoff_time,
			dld.landing_time,
			dld.in_time,
			dld.pilot_role,
			dld.companion_name,
			dld.air_time,
			dld.block_time,
			dld.duty_time,
			dld.approach_type,
			dld.flight_type,
			dld.employee_logbook_id,
			dl.log_date,
			lp.license_plate,
			am.model_name,
			CONCAT(orig.iata_code, '-', dest.iata_code) as route_code,
			orig.iata_code as origin_iata_code,
			dest.iata_code as destination_iata_code,
			airl.airline_code
		FROM daily_logbook_detail dld
		INNER JOIN daily_logbook dl ON dld.daily_logbook_id = dl.id
		INNER JOIN license_plate lp ON dld.actual_license_plate_id = lp.id
		INNER JOIN aircraft_model am ON lp.aircraft_model_id = am.id
		INNER JOIN airline_route alr ON dld.airline_route_id = alr.id
		INNER JOIN route r ON alr.route_id = r.id
		INNER JOIN airport orig ON r.origin_airport_id = orig.id
		INNER JOIN airport dest ON r.destination_airport_id = dest.id
		INNER JOIN airline airl ON alr.airline_id = airl.id
		WHERE dld.daily_logbook_id = ?
		ORDER BY dld.out_time ASC
	`

	// Query for listing all flight details by employee ID (joins through daily_logbook)
	QueryByEmployee = `
		SELECT
			dld.id,
			dld.daily_logbook_id,
			dld.flight_real_date,
			dld.flight_number,
			dld.airline_route_id,
			dld.actual_license_plate_id,
			dld.passengers,
			dld.out_time,
			dld.takeoff_time,
			dld.landing_time,
			dld.in_time,
			dld.pilot_role,
			dld.companion_name,
			dld.air_time,
			dld.block_time,
			dld.duty_time,
			dld.approach_type,
			dld.flight_type,
			dld.employee_logbook_id,
			dl.log_date,
			lp.license_plate,
			am.model_name,
			CONCAT(orig.iata_code, '-', dest.iata_code) as route_code,
			orig.iata_code as origin_iata_code,
			dest.iata_code as destination_iata_code,
			airl.airline_code
		FROM daily_logbook_detail dld
		INNER JOIN daily_logbook dl ON dld.daily_logbook_id = dl.id
		INNER JOIN license_plate lp ON dld.actual_license_plate_id = lp.id
		INNER JOIN aircraft_model am ON lp.aircraft_model_id = am.id
		INNER JOIN airline_route alr ON dld.airline_route_id = alr.id
		INNER JOIN route r ON alr.route_id = r.id
		INNER JOIN airport orig ON r.origin_airport_id = orig.id
		INNER JOIN airport dest ON r.destination_airport_id = dest.id
		INNER JOIN airline airl ON alr.airline_id = airl.id
		WHERE dl.employee_id = ?
		ORDER BY dld.flight_real_date DESC, dld.out_time ASC
	`

	// Insert query
	QueryInsert = `
		INSERT INTO daily_logbook_detail (
			id, daily_logbook_id, flight_real_date, flight_number,
			airline_route_id, actual_license_plate_id, passengers,
			out_time, takeoff_time, landing_time, in_time,
			pilot_role, companion_name,
			air_time, block_time, duty_time,
			approach_type, flight_type, employee_logbook_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	// Update query
	QueryUpdate = `
		UPDATE daily_logbook_detail SET
			flight_real_date = ?,
			flight_number = ?,
			airline_route_id = ?,
			actual_license_plate_id = ?,
			passengers = ?,
			out_time = ?,
			takeoff_time = ?,
			landing_time = ?,
			in_time = ?,
			pilot_role = ?,
			companion_name = ?,
			air_time = ?,
			block_time = ?,
			duty_time = ?,
			approach_type = ?,
			flight_type = ?
		WHERE id = ?
	`

	// Delete query
	QueryDelete = `DELETE FROM daily_logbook_detail WHERE id = ?`
)

var log logger.Logger = logger.NewSlogLogger()

type repository struct {
	stmtGetByID       *sql.Stmt
	stmtGetByLogbook  *sql.Stmt
	stmtGetByEmployee *sql.Stmt
	stmtInsert        *sql.Stmt
	stmtUpdate        *sql.Stmt
	db                *sql.DB
}

func NewDailyLogbookDetailRepository(db *sql.DB) (*repository, error) {
	if db == nil {
		return nil, sql.ErrConnDone
	}

	stmtGetByID, err := db.Prepare(QueryByID)
	if err != nil {
		log.Error(logger.LogDailyLogbookDetailRepoInitError, "error preparing statement", err)
		return nil, err
	}

	stmtGetByLogbook, err := db.Prepare(QueryByLogbook)
	if err != nil {
		log.Error(logger.LogDailyLogbookDetailRepoInitError, "error preparing statement", err)
		return nil, err
	}

	stmtGetByEmployee, err := db.Prepare(QueryByEmployee)
	if err != nil {
		log.Error(logger.LogDailyLogbookDetailRepoInitError, "error preparing statement", err)
		return nil, err
	}

	stmtInsert, err := db.Prepare(QueryInsert)
	if err != nil {
		log.Error(logger.LogDailyLogbookDetailRepoInitError, "error preparing statement", err)
		return nil, err
	}

	stmtUpdate, err := db.Prepare(QueryUpdate)
	if err != nil {
		log.Error(logger.LogDailyLogbookDetailRepoInitError, "error preparing statement", err)
		return nil, err
	}

	log.Info(logger.LogDailyLogbookDetailRepoInitOK)

	return &repository{
		db:                db,
		stmtGetByID:       stmtGetByID,
		stmtGetByLogbook:  stmtGetByLogbook,
		stmtGetByEmployee: stmtGetByEmployee,
		stmtInsert:        stmtInsert,
		stmtUpdate:        stmtUpdate,
	}, nil
}

func (r *repository) BeginTx(ctx context.Context) (output.Tx, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return common.NewSQLTx(tx), nil
}
