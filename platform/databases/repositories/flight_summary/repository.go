package flight_summary

import (
	"context"
	"database/sql"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/databases/common"
	"github.com/champion19/api-flighthours/platform/logger"
)

const (
	// Aggregates air_time by pilot_role for a given employee and date range
	QueryFlightHoursSummary = `
		SELECT
			COALESCE(dld.pilot_role, 'UNSET') as pilot_role,
			COALESCE(SUM(TIME_TO_SEC(dld.air_time)), 0) as total_seconds,
			COUNT(*) as flight_count
		FROM daily_logbook_detail dld
		INNER JOIN daily_logbook dl ON dld.daily_logbook_id = dl.id
		WHERE dl.employee_id = ?
		  AND dld.flight_real_date >= ?
		  AND dld.flight_real_date <= ?
		GROUP BY dld.pilot_role`

	// Counts landings (only PF and PFL roles count as landings) for a date range
	QueryLandingCount = `
		SELECT COUNT(*)
		FROM daily_logbook_detail dld
		INNER JOIN daily_logbook dl ON dld.daily_logbook_id = dl.id
		WHERE dl.employee_id = ?
		  AND dld.flight_real_date >= ?
		  AND dld.flight_real_date <= ?
		  AND dld.pilot_role IN ('PF', 'PFL')`

	// Gets the sum of air_time (in seconds) for flights on a specific date to estimate consecutive hours
	QueryDailyFlightSeconds = `
		SELECT COALESCE(SUM(TIME_TO_SEC(dld.air_time)), 0) as total_seconds
		FROM daily_logbook_detail dld
		INNER JOIN daily_logbook dl ON dld.daily_logbook_id = dl.id
		WHERE dl.employee_id = ?
		  AND dld.flight_real_date = ?`

	// Gets the last N flights ordered by date descending, with full join data
	QueryRecentFlights = `SELECT` + selectDetailColumns + detailJoins + `
		WHERE dl.employee_id = ?
		ORDER BY dld.flight_real_date DESC, dld.out_time DESC
		LIMIT ?`

	// Shared SELECT columns (reuses daily_logbook_detail pattern)
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
		dld.companion_name,
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

	detailJoins = `
		FROM daily_logbook_detail dld
		INNER JOIN daily_logbook dl ON dld.daily_logbook_id = dl.id
		INNER JOIN tail_number tn ON dld.actual_tail_number_id = tn.id
		INNER JOIN aircraft_model am ON tn.aircraft_model_id = am.id
		INNER JOIN airport orig ON dld.origin_airport_id = orig.id
		INNER JOIN airport dest ON dld.destination_airport_id = dest.id
		INNER JOIN employee emp ON dl.employee_id = emp.id
		INNER JOIN airline airl ON emp.airline = airl.id`
)

var log logger.Logger = logger.NewSlogLogger()

type repository struct {
	db                     *sql.DB
	stmtFlightHoursSummary *sql.Stmt
	stmtLandingCount       *sql.Stmt
	stmtDailyFlightSeconds *sql.Stmt
	stmtRecentFlights      *sql.Stmt
}

func NewFlightSummaryRepository(db *sql.DB) (*repository, error) {
	if db == nil {
		return nil, sql.ErrConnDone
	}

	prepare := func(query string) (*sql.Stmt, error) {
		stmt, err := db.Prepare(query)
		if err != nil {
			log.Error(logger.LogFlightSummaryRepoInitError, "error preparing statement", err)
		}
		return stmt, err
	}

	stmtFlightHoursSummary, err := prepare(QueryFlightHoursSummary)
	if err != nil {
		return nil, err
	}
	stmtLandingCount, err := prepare(QueryLandingCount)
	if err != nil {
		return nil, err
	}
	stmtDailyFlightSeconds, err := prepare(QueryDailyFlightSeconds)
	if err != nil {
		return nil, err
	}
	stmtRecentFlights, err := prepare(QueryRecentFlights)
	if err != nil {
		return nil, err
	}

	log.Info(logger.LogFlightSummaryRepoInitOK)

	return &repository{
		db:                     db,
		stmtFlightHoursSummary: stmtFlightHoursSummary,
		stmtLandingCount:       stmtLandingCount,
		stmtDailyFlightSeconds: stmtDailyFlightSeconds,
		stmtRecentFlights:      stmtRecentFlights,
	}, nil
}

func (r *repository) BeginTx(ctx context.Context) (output.Tx, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return common.NewSQLTx(tx), nil
}

// GetFlightHoursSummary returns pilot_role breakdown for a date range
func (r *repository) GetFlightHoursSummary(ctx context.Context, employeeID, startDate, endDate string) ([]domain.PilotRoleBreakdown, error) {
	log.Info(logger.LogFlightSummaryGet, "employee_id", employeeID, "start", startDate, "end", endDate)

	rows, err := r.stmtFlightHoursSummary.QueryContext(ctx, employeeID, startDate, endDate)
	if err != nil {
		log.Error(logger.LogFlightSummaryGetError, "error", err)
		return nil, err
	}
	defer rows.Close()

	var breakdown []domain.PilotRoleBreakdown
	for rows.Next() {
		var b domain.PilotRoleBreakdown
		if err := rows.Scan(&b.PilotRole, &b.TotalSeconds, &b.FlightCount); err != nil {
			log.Error(logger.LogFlightSummaryGetError, "scan_error", err)
			return nil, err
		}
		breakdown = append(breakdown, b)
	}

	if err = rows.Err(); err != nil {
		log.Error(logger.LogFlightSummaryGetError, "rows_error", err)
		return nil, err
	}

	log.Info(logger.LogFlightSummaryGetOK, "employee_id", employeeID, "roles_count", len(breakdown))
	return breakdown, nil
}

// GetLandingCount returns the number of landings in a date range
func (r *repository) GetLandingCount(ctx context.Context, employeeID, startDate, endDate string) (int, error) {
	log.Info(logger.LogFlightSummaryGet, "action", "landing_count", "employee_id", employeeID)

	var count int
	err := r.stmtLandingCount.QueryRowContext(ctx, employeeID, startDate, endDate).Scan(&count)
	if err != nil {
		log.Error(logger.LogFlightSummaryGetError, "action", "landing_count", "error", err)
		return 0, err
	}

	return count, nil
}

// GetDailyFlightSeconds returns the total seconds of air_time for a specific date
func (r *repository) GetDailyFlightSeconds(ctx context.Context, employeeID, date string) (int, error) {
	log.Info(logger.LogFlightSummaryGet, "action", "daily_flight_seconds", "employee_id", employeeID, "date", date)

	var totalSeconds int
	err := r.stmtDailyFlightSeconds.QueryRowContext(ctx, employeeID, date).Scan(&totalSeconds)
	if err != nil {
		log.Error(logger.LogFlightSummaryGetError, "action", "daily_flight_seconds", "error", err)
		return 0, err
	}

	return totalSeconds, nil
}

// GetRecentFlights returns the last N flights for an employee
func (r *repository) GetRecentFlights(ctx context.Context, employeeID string, limit int) ([]domain.DailyLogbookDetail, error) {
	log.Info(logger.LogFlightSummaryGet, "action", "recent_flights", "employee_id", employeeID, "limit", limit)

	rows, err := r.stmtRecentFlights.QueryContext(ctx, employeeID, limit)
	if err != nil {
		log.Error(logger.LogFlightSummaryGetError, "action", "recent_flights", "error", err)
		return nil, err
	}
	defer rows.Close()

	var details []domain.DailyLogbookDetail
	for rows.Next() {
		entity, err := scanDetail(rows)
		if err != nil {
			log.Error(logger.LogFlightSummaryGetError, "action", "recent_flights", "scan_error", err)
			return nil, err
		}
		details = append(details, *entity)
	}

	if err = rows.Err(); err != nil {
		log.Error(logger.LogFlightSummaryGetError, "action", "recent_flights", "rows_error", err)
		return nil, err
	}

	log.Info(logger.LogFlightSummaryGetOK, "action", "recent_flights", "count", len(details))
	return details, nil
}
