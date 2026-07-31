package crew_member

import (
	"context"
	"database/sql"

	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/databases/common"
	"github.com/champion19/api-flighthours/platform/logger"
)

const (
	// QuerySearch finds crew members in a pilot's own roster whose name or bp
	// matches (case-insensitive, via LIKE) — bp is the more reliable identifier
	// (badge/carné number) so it's matched too, not just the free-text name.
	QuerySearch = `
		SELECT id, employee_id, name, bp
		FROM crew_member
		WHERE employee_id = ? AND (name LIKE ? OR bp LIKE ?)
		ORDER BY name ASC`

	// QueryGetByID fetches a single crew member by ID.
	QueryGetByID = `SELECT id, employee_id, name, bp FROM crew_member WHERE id = ? LIMIT 1`

	// QueryGetByEmployeeAndName fetches a crew member by its unique (employee_id, name) key.
	QueryGetByEmployeeAndName = `SELECT id, employee_id, name, bp FROM crew_member WHERE employee_id = ? AND name = ? LIMIT 1`

	// QueryGetByEmployeeAndBP fetches a crew member by (employee_id, bp) — bp is
	// the reliable identifier, so this is checked before falling back to name.
	QueryGetByEmployeeAndBP = `SELECT id, employee_id, name, bp FROM crew_member WHERE employee_id = ? AND bp = ? LIMIT 1`

	// QueryInsertIgnoreDuplicate inserts a new crew member, no-op if (employee_id, name) already exists.
	// If it already exists and a non-empty bp is now provided, fills it in without overwriting an existing value.
	QueryInsertIgnoreDuplicate = `
		INSERT INTO crew_member (id, employee_id, name, bp)
		VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE bp = COALESCE(bp, VALUES(bp))`
)

var log logger.Logger = logger.NewSlogLogger()

type repository struct {
	stmtSearch               *sql.Stmt
	stmtGetByID              *sql.Stmt
	stmtGetByEmployeeAndName *sql.Stmt
	stmtGetByEmployeeAndBP   *sql.Stmt
	stmtInsertIgnore         *sql.Stmt
	db                       *sql.DB
}

func NewCrewMemberRepository(db *sql.DB) (*repository, error) {
	if db == nil {
		return nil, sql.ErrConnDone
	}

	prepare := func(query string) (*sql.Stmt, error) {
		stmt, err := db.Prepare(query)
		if err != nil {
			log.Error(logger.LogCrewMemberRepoInitError, "error preparing statement", err)
		}
		return stmt, err
	}

	stmtSearch, err := prepare(QuerySearch)
	if err != nil {
		return nil, err
	}
	stmtGetByID, err := prepare(QueryGetByID)
	if err != nil {
		return nil, err
	}
	stmtGetByEmployeeAndName, err := prepare(QueryGetByEmployeeAndName)
	if err != nil {
		return nil, err
	}
	stmtGetByEmployeeAndBP, err := prepare(QueryGetByEmployeeAndBP)
	if err != nil {
		return nil, err
	}
	stmtInsertIgnore, err := prepare(QueryInsertIgnoreDuplicate)
	if err != nil {
		return nil, err
	}

	log.Info(logger.LogCrewMemberRepoInitOK)

	return &repository{
		db:                       db,
		stmtSearch:               stmtSearch,
		stmtGetByID:              stmtGetByID,
		stmtGetByEmployeeAndName: stmtGetByEmployeeAndName,
		stmtGetByEmployeeAndBP:   stmtGetByEmployeeAndBP,
		stmtInsertIgnore:         stmtInsertIgnore,
	}, nil
}

func (r *repository) BeginTx(ctx context.Context) (output.Tx, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return common.NewSQLTx(tx), nil
}
