package daily_logbook_detail

import (
	"context"
	"database/sql"

	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/databases/common"
	"github.com/champion19/api-flighthours/platform/logger"
)

const (
	// Shared SELECT columns for daily_logbook_detail queries (used by QueryByID, QueryByLogbook, QueryByEmployee)
	selectDetailColumns = `
		dld.id,
		dld.daily_logbook_id,
		dld.flight_real_date,
		dld.flight_number,
		dld.origin_airport_id,
		dld.destination_airport_id,
		dld.actual_tail_number_id,
		dld.passengers,
		dld.out_time,
		dld.takeoff_time,
		dld.landing_time,
		dld.in_time,
		dld.pilot_role,
		dld.crew_role,
		dld.air_time,
		dld.block_time,
		dld.approach_category,
		dld.approach_subtype,
		dld.autoland,
		dld.flight_type,
		dld.employee_logbook_id,
		dl.log_date,
		tn.tail_number,
		am.model_name,
		CONCAT(orig.iata_code, '-', dest.iata_code) as route_code,
		orig.iata_code as origin_iata_code,
		dest.iata_code as destination_iata_code,
		airl.airline_code`

	// Shared JOIN clause for daily_logbook_detail queries
	detailJoins = `
		FROM daily_logbook_detail dld
		INNER JOIN daily_logbook dl ON dld.daily_logbook_id = dl.id
		INNER JOIN tail_number tn ON dld.actual_tail_number_id = tn.id
		INNER JOIN aircraft_model am ON tn.aircraft_model_id = am.id
		INNER JOIN airport orig ON dld.origin_airport_id = orig.id
		INNER JOIN airport dest ON dld.destination_airport_id = dest.id
		INNER JOIN employee emp ON dl.employee_id = emp.id
		INNER JOIN airline airl ON emp.airline = airl.id`

	// Query for getting a detail by ID with JOINs for denormalized data
	QueryByID = `
		SELECT
			dld.id,
			dld.daily_logbook_id,
			dld.flight_real_date,
			dld.flight_number,
			dld.origin_airport_id,
			dld.destination_airport_id,
			dld.actual_tail_number_id,
			dld.passengers,
			dld.out_time,
			dld.takeoff_time,
			dld.landing_time,
			dld.in_time,
			dld.pilot_role,
			dld.crew_role,
			dld.air_time,
			dld.block_time,
			dld.approach_category,
			dld.approach_subtype,
			dld.autoland,
			dld.flight_type,
			dld.employee_logbook_id,
			dl.log_date,
			tn.tail_number,
			am.model_name,
			CONCAT(orig.iata_code, '-', dest.iata_code) as route_code,
			orig.iata_code as origin_iata_code,
			dest.iata_code as destination_iata_code,
			airl.airline_code
		FROM daily_logbook_detail dld
		INNER JOIN daily_logbook dl ON dld.daily_logbook_id = dl.id
		INNER JOIN tail_number tn ON dld.actual_tail_number_id = tn.id
		INNER JOIN aircraft_model am ON tn.aircraft_model_id = am.id
		INNER JOIN airport orig ON dld.origin_airport_id = orig.id
		INNER JOIN airport dest ON dld.destination_airport_id = dest.id
		INNER JOIN employee emp ON dl.employee_id = emp.id
		INNER JOIN airline airl ON emp.airline = airl.id
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
			dld.origin_airport_id,
			dld.destination_airport_id,
			dld.actual_tail_number_id,
			dld.passengers,
			dld.out_time,
			dld.takeoff_time,
			dld.landing_time,
			dld.in_time,
			dld.pilot_role,
			dld.crew_role,
			dld.air_time,
			dld.block_time,
			dld.approach_category,
			dld.approach_subtype,
			dld.autoland,
			dld.flight_type,
			dld.employee_logbook_id,
			dl.log_date,
			tn.tail_number,
			am.model_name,
			CONCAT(orig.iata_code, '-', dest.iata_code) as route_code,
			orig.iata_code as origin_iata_code,
			dest.iata_code as destination_iata_code,
			airl.airline_code
		FROM daily_logbook_detail dld
		INNER JOIN daily_logbook dl ON dld.daily_logbook_id = dl.id
		INNER JOIN tail_number tn ON dld.actual_tail_number_id = tn.id
		INNER JOIN aircraft_model am ON tn.aircraft_model_id = am.id
		INNER JOIN airport orig ON dld.origin_airport_id = orig.id
		INNER JOIN airport dest ON dld.destination_airport_id = dest.id
		INNER JOIN employee emp ON dl.employee_id = emp.id
		INNER JOIN airline airl ON emp.airline = airl.id
		WHERE dld.daily_logbook_id = ?
		ORDER BY dld.out_time ASC
	`

	// Query for listing all flight details by employee ID (joins through daily_logbook)
	QueryByEmployee = `SELECT` + selectDetailColumns + detailJoins + `
		WHERE dl.employee_id = ?
		ORDER BY dld.flight_real_date DESC, dld.out_time ASC`

	// Insert query
	QueryInsert = `
		INSERT INTO daily_logbook_detail (
			id, daily_logbook_id, flight_real_date, flight_number,
			origin_airport_id, destination_airport_id, actual_tail_number_id, passengers,
			out_time, takeoff_time, landing_time, in_time,
			pilot_role, crew_role,
			air_time, block_time,
			approach_category, approach_subtype, autoland, flight_type, employee_logbook_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	// Update query
	QueryUpdate = `
		UPDATE daily_logbook_detail SET
			flight_real_date = ?,
			flight_number = ?,
			origin_airport_id = ?,
			destination_airport_id = ?,
			actual_tail_number_id = ?,
			passengers = ?,
			out_time = ?,
			takeoff_time = ?,
			landing_time = ?,
			in_time = ?,
			pilot_role = ?,
			crew_role = ?,
			air_time = ?,
			block_time = ?,
			approach_category = ?,
			approach_subtype = ?,
			autoland = ?,
			flight_type = ?
		WHERE id = ?`

	// Delete query
	QueryDelete = `DELETE FROM daily_logbook_detail WHERE id = ?`

	// Check duplicate by unique key
	QueryExistsByUniqueKey = `
		SELECT COUNT(*) FROM daily_logbook_detail
		WHERE employee_logbook_id = ?
		  AND flight_real_date = ?
		  AND flight_number = ?
		  AND actual_tail_number_id = ?`
)

var log logger.Logger = logger.NewSlogLogger()

type repository struct {
	stmtGetByID           *sql.Stmt
	stmtGetByLogbook      *sql.Stmt
	stmtGetByEmployee     *sql.Stmt
	stmtInsert            *sql.Stmt
	stmtUpdate            *sql.Stmt
	stmtExistsByUniqueKey *sql.Stmt
	stmtDelete            *sql.Stmt
	db                    *sql.DB
}

func NewDailyLogbookDetailRepository(db *sql.DB) (*repository, error) {
	if db == nil {
		return nil, sql.ErrConnDone
	}

	prepare := func(query string) (*sql.Stmt, error) {
		stmt, err := db.Prepare(query)
		if err != nil {
			log.Error(logger.LogDailyLogbookDetailRepoInitError, "error preparing statement", err)
		}
		return stmt, err
	}

	stmtGetByID, err := prepare(QueryByID)
	if err != nil {
		return nil, err
	}
	stmtGetByLogbook, err := prepare(QueryByLogbook)
	if err != nil {
		return nil, err
	}
	stmtGetByEmployee, err := prepare(QueryByEmployee)
	if err != nil {
		return nil, err
	}
	stmtInsert, err := prepare(QueryInsert)
	if err != nil {
		return nil, err
	}
	stmtUpdate, err := prepare(QueryUpdate)
	if err != nil {
		return nil, err
	}
	stmtExistsByUniqueKey, err := prepare(QueryExistsByUniqueKey)
	if err != nil {
		return nil, err
	}
	stmtDelete, err := prepare(QueryDelete)
	if err != nil {
		return nil, err
	}

	log.Info(logger.LogDailyLogbookDetailRepoInitOK)

	return &repository{
		db:                    db,
		stmtGetByID:           stmtGetByID,
		stmtGetByLogbook:      stmtGetByLogbook,
		stmtGetByEmployee:     stmtGetByEmployee,
		stmtInsert:            stmtInsert,
		stmtUpdate:            stmtUpdate,
		stmtExistsByUniqueKey: stmtExistsByUniqueKey,
		stmtDelete:            stmtDelete,
	}, nil
}

func (r *repository) BeginTx(ctx context.Context) (output.Tx, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return common.NewSQLTx(tx), nil
}
