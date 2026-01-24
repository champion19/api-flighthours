package employee

import (
	"context"
	"database/sql"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/output"
	"github.com/champion19/flighthours-api/platform/databases/common"
	"github.com/champion19/flighthours-api/platform/logger"
)

const (
	QuerySave         = "INSERT INTO employee(id,name,airline,email,identification_number,bp,start_date,end_date,active,role,keycloak_user_id) VALUES(?,?,?,?,?,?,?,?,?,?,?)"
	QueryByEmail      = "SELECT id,name,airline,email,identification_number,bp,start_date,end_date,active,role,keycloak_user_id FROM employee WHERE email=? LIMIT 1"
	QueryByID         = "SELECT id,name,airline,email,identification_number,bp,start_date,end_date,active,role,keycloak_user_id FROM employee WHERE id=? LIMIT 1"
	QueryByKeycloakID = "SELECT id,name,airline,email,identification_number,bp,start_date,end_date,active,role,keycloak_user_id FROM employee WHERE keycloak_user_id=? LIMIT 1"
	QueryUpdate       = "UPDATE employee SET name=?,airline=?,email=?,identification_number=?,bp=?,start_date=?,end_date=?,active=?,role=?,keycloak_user_id=? WHERE id=?"
	QueryDelete       = "DELETE FROM employee WHERE id=?"
	QueryPatch        = "UPDATE employee SET keycloak_user_id=? WHERE id=?"
	// HU47 - Get employees by role (Virtual Entity pattern - no new table needed)
	QueryByRole = "SELECT id,name,airline,email,identification_number,bp,start_date,end_date,active,role,keycloak_user_id FROM employee WHERE role=? ORDER BY name"
)

var log logger.Logger = logger.NewSlogLogger()

type repository struct {
	stmtSave       *sql.Stmt
	stmtGetByEmail *sql.Stmt
	stmtGetByID    *sql.Stmt
	stmtUpdate     *sql.Stmt
	stmtDelete     *sql.Stmt
	stmtPatch      *sql.Stmt
	stmtGetByRole  *sql.Stmt
	db             *sql.DB
}

func NewClientRepository(db *sql.DB) (*repository, error) {
	if db == nil {
		return nil, sql.ErrConnDone
	}
	stmtSave, err := db.Prepare(QuerySave)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing statement", err)
		return nil, err
	}
	stmtGetByEmail, err := db.Prepare(QueryByEmail)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing statement", err)
		return nil, err
	}
	stmtGetByID, err := db.Prepare(QueryByID)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing statement", err)
		return nil, err
	}
	stmtUpdate, err := db.Prepare(QueryUpdate)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing statement", err)
		return nil, err
	}
	stmtDelete, err := db.Prepare(QueryDelete)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing statement", err)
		return nil, err
	}

	stmtPatch, err := db.Prepare(QueryPatch)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing statement", err)
		return nil, err
	}

	// HU47 - Prepare statement for role lookup
	stmtGetByRole, err := db.Prepare(QueryByRole)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing role statement", err)
		return nil, err
	}

	return &repository{
		db:             db,
		stmtSave:       stmtSave,
		stmtGetByEmail: stmtGetByEmail,
		stmtGetByID:    stmtGetByID,
		stmtUpdate:     stmtUpdate,
		stmtDelete:     stmtDelete,
		stmtPatch:      stmtPatch,
		stmtGetByRole:  stmtGetByRole,
	}, nil
}

func (r *repository) BeginTx(ctx context.Context) (output.Tx, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return common.NewSQLTx(tx), nil
}

// GetEmployeesByRole retrieves all employees for a specific role (HU47 - Virtual Entity pattern)
// This implements the "Derived Values" approach - no new table needed, we query employees by role field
func (r *repository) GetEmployeesByRole(ctx context.Context, role string) ([]domain.Employee, error) {
	rows, err := r.stmtGetByRole.QueryContext(ctx, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []domain.Employee
	for rows.Next() {
		var e Employee
		if err := rows.Scan(
			&e.ID,
			&e.Name,
			&e.Airline,
			&e.Email,
			&e.IdentificationNumber,
			&e.Bp,
			&e.StartDate,
			&e.EndDate,
			&e.Active,
			&e.Role,
			&e.KeycloakUserID,
		); err != nil {
			return nil, err
		}
		employees = append(employees, e.ToDomain())
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Return sql.ErrNoRows if no employees found (caller handles not found logic)
	if len(employees) == 0 {
		return nil, sql.ErrNoRows
	}

	return employees, nil
}
