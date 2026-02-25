package employee_flight_summary

import (
	"context"
	"database/sql"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/databases/common"
	"github.com/champion19/api-flighthours/platform/logger"
	"github.com/google/uuid"
)

const (
	// INSERT ON DUPLICATE KEY UPDATE — atomically upserts a period row
	QueryUpsertSummary = `
		INSERT INTO employee_flight_summary
			(id, employee_id, period_type, period_year, period_number, period_start, period_end,
			 total_air_time, total_block_time, total_flights, total_landings)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			total_air_time   = total_air_time + VALUES(total_air_time),
			total_block_time = total_block_time + VALUES(total_block_time),
			total_flights    = total_flights + VALUES(total_flights),
			total_landings   = total_landings + VALUES(total_landings)`

	// Get all summary rows for an employee by period type
	QueryGetSummaries = `
		SELECT id, employee_id, period_type, period_year, period_number,
		       period_start, period_end, total_air_time, total_block_time,
		       total_duty_time, total_flights, total_landings, last_updated
		FROM employee_flight_summary
		WHERE employee_id = ? AND period_type = ?
		ORDER BY period_year DESC, period_number DESC`

	// Get a specific period row
	QueryGetCurrentPeriod = `
		SELECT id, employee_id, period_type, period_year, period_number,
		       period_start, period_end, total_air_time, total_block_time,
		       total_duty_time, total_flights, total_landings, last_updated
		FROM employee_flight_summary
		WHERE employee_id = ? AND period_type = ? AND period_year = ? AND period_number = ?`

	// Get all current period summaries for an employee (for alerts dashboard)
	QueryGetAllCurrentSummaries = `
		SELECT id, employee_id, period_type, period_year, period_number,
		       period_start, period_end, total_air_time, total_block_time,
		       total_duty_time, total_flights, total_landings, last_updated
		FROM employee_flight_summary
		WHERE employee_id = ?
		ORDER BY period_type, period_year DESC, period_number DESC`
)

var log logger.Logger = logger.NewSlogLogger()

type repository struct {
	db                    *sql.DB
	stmtUpsert            *sql.Stmt
	stmtGetSummaries      *sql.Stmt
	stmtGetCurrentPeriod  *sql.Stmt
	stmtGetAllCurrentSumm *sql.Stmt
}

func NewEmployeeFlightSummaryRepository(db *sql.DB) (*repository, error) {
	if db == nil {
		return nil, sql.ErrConnDone
	}

	prepare := func(query string) (*sql.Stmt, error) {
		stmt, err := db.Prepare(query)
		if err != nil {
			log.Error(logger.LogFlightSummaryRepoInitError, "error", err)
		}
		return stmt, err
	}

	stmtUpsert, err := prepare(QueryUpsertSummary)
	if err != nil {
		return nil, err
	}
	stmtGetSummaries, err := prepare(QueryGetSummaries)
	if err != nil {
		return nil, err
	}
	stmtGetCurrentPeriod, err := prepare(QueryGetCurrentPeriod)
	if err != nil {
		return nil, err
	}
	stmtGetAllCurrentSumm, err := prepare(QueryGetAllCurrentSummaries)
	if err != nil {
		return nil, err
	}

	log.Info(logger.LogFlightSummaryRepoInitOK, "repository", "employee_flight_summary")

	return &repository{
		db:                    db,
		stmtUpsert:            stmtUpsert,
		stmtGetSummaries:      stmtGetSummaries,
		stmtGetCurrentPeriod:  stmtGetCurrentPeriod,
		stmtGetAllCurrentSumm: stmtGetAllCurrentSumm,
	}, nil
}

func (r *repository) BeginTx(ctx context.Context) (output.Tx, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return common.NewSQLTx(tx), nil
}

// UpsertSummary atomically inserts or updates a period summary row.
// delta values can be positive (create) or negative (delete).
func (r *repository) UpsertSummary(ctx context.Context, tx output.Tx, employeeID string, period domain.PeriodInfo, airTimeDelta, blockTimeDelta, flightsDelta, landingsDelta int) error {
	log.Info(logger.LogFlightSummaryGet, "action", "upsert_summary",
		"employee_id", employeeID,
		"period_type", period.PeriodType,
		"period_year", period.PeriodYear,
		"period_number", period.PeriodNumber)

	sqlTx, err := common.CastTx(tx)
	if err != nil {
		return err
	}

	newID := uuid.New().String()

	stmt := sqlTx.Tx.StmtContext(ctx, r.stmtUpsert)
	_, err = stmt.ExecContext(ctx,
		newID, employeeID,
		period.PeriodType, period.PeriodYear, period.PeriodNumber,
		period.PeriodStart, period.PeriodEnd,
		airTimeDelta, blockTimeDelta, flightsDelta, landingsDelta,
	)
	if err != nil {
		log.Error(logger.LogFlightSummaryGetError, "action", "upsert_summary", "error", err)
		return err
	}

	log.Info(logger.LogFlightSummaryGetOK, "action", "upsert_summary", "period_type", period.PeriodType)
	return nil
}

// GetSummariesByEmployee returns all summary rows for a period type, ordered by newest first
func (r *repository) GetSummariesByEmployee(ctx context.Context, employeeID, periodType string) ([]domain.EmployeeFlightSummary, error) {
	rows, err := r.stmtGetSummaries.QueryContext(ctx, employeeID, periodType)
	if err != nil {
		log.Error(logger.LogFlightSummaryGetError, "action", "get_summaries", "error", err)
		return nil, err
	}
	defer rows.Close()

	return scanSummaries(rows)
}

// GetCurrentPeriodSummary returns a single period summary
func (r *repository) GetCurrentPeriodSummary(ctx context.Context, employeeID, periodType string, year, number int) (*domain.EmployeeFlightSummary, error) {
	row := r.stmtGetCurrentPeriod.QueryRowContext(ctx, employeeID, periodType, year, number)
	return scanSingleSummary(row)
}

// GetAllSummaries returns all summary rows for an employee (for alerts dashboard)
func (r *repository) GetAllSummaries(ctx context.Context, employeeID string) ([]domain.EmployeeFlightSummary, error) {
	rows, err := r.stmtGetAllCurrentSumm.QueryContext(ctx, employeeID)
	if err != nil {
		log.Error(logger.LogFlightSummaryGetError, "action", "get_all_summaries", "error", err)
		return nil, err
	}
	defer rows.Close()

	return scanSummaries(rows)
}

func scanSummaries(rows *sql.Rows) ([]domain.EmployeeFlightSummary, error) {
	var summaries []domain.EmployeeFlightSummary
	for rows.Next() {
		s, err := scanSummaryRow(rows)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, *s)
	}
	return summaries, rows.Err()
}

func scanSummaryRow(scanner interface {
	Scan(dest ...interface{}) error
}) (*domain.EmployeeFlightSummary, error) {
	var s domain.EmployeeFlightSummary
	var lastUpdated sql.NullString

	err := scanner.Scan(
		&s.ID, &s.EmployeeID, &s.PeriodType,
		&s.PeriodYear, &s.PeriodNumber,
		&s.PeriodStart, &s.PeriodEnd,
		&s.TotalAirTime, &s.TotalBlockTime, &s.TotalDutyTime,
		&s.TotalFlights, &s.TotalLandings,
		&lastUpdated,
	)
	if err != nil {
		return nil, err
	}

	if lastUpdated.Valid {
		s.LastUpdated = lastUpdated.String
	}

	return &s, nil
}

func scanSingleSummary(row *sql.Row) (*domain.EmployeeFlightSummary, error) {
	var s domain.EmployeeFlightSummary
	var lastUpdated sql.NullString

	err := row.Scan(
		&s.ID, &s.EmployeeID, &s.PeriodType,
		&s.PeriodYear, &s.PeriodNumber,
		&s.PeriodStart, &s.PeriodEnd,
		&s.TotalAirTime, &s.TotalBlockTime, &s.TotalDutyTime,
		&s.TotalFlights, &s.TotalLandings,
		&lastUpdated,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if lastUpdated.Valid {
		s.LastUpdated = lastUpdated.String
	}

	return &s, nil
}
