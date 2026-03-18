package employee

import (
	"context"
	"database/sql"

	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/databases/common"
	"github.com/champion19/api-flighthours/platform/logger"
)

const (
	QuerySave         = "INSERT INTO employee(id,name,airline,email,identification_number,bp,start_date,end_date,active,role,keycloak_user_id) VALUES(?,?,?,?,?,?,?,?,?,?,?)"
	QueryByEmail      = "SELECT id,name,airline,email,identification_number,bp,start_date,end_date,active,role,keycloak_user_id FROM employee WHERE email=? LIMIT 1"
	QueryByID         = "SELECT id,name,airline,email,identification_number,bp,start_date,end_date,active,role,keycloak_user_id FROM employee WHERE id=? LIMIT 1"
	QueryByKeycloakID = "SELECT id,name,airline,email,identification_number,bp,start_date,end_date,active,role,keycloak_user_id FROM employee WHERE keycloak_user_id=? LIMIT 1"
	QueryUpdate       = "UPDATE employee SET name=?,airline=?,email=?,identification_number=?,bp=?,start_date=?,end_date=?,active=?,role=?,keycloak_user_id=? WHERE id=?"
	QueryDelete       = "DELETE FROM employee WHERE id=?"
	QueryPatch        = "UPDATE employee SET keycloak_user_id=? WHERE id=?"
	QueryList         = "SELECT id,name,airline,email,identification_number,bp,start_date,end_date,active,role,keycloak_user_id FROM employee ORDER BY name"
)

var log logger.Logger = logger.NewSlogLogger()

const errPreparingStatement = "error preparing statement"

type repository struct {
	stmtSave       *sql.Stmt
	stmtGetByEmail *sql.Stmt
	stmtGetByID    *sql.Stmt
	stmtUpdate     *sql.Stmt
	stmtDelete     *sql.Stmt
	stmtPatch      *sql.Stmt
	db             *sql.DB
}

func NewClientRepository(db *sql.DB) (*repository, error) {
	if db == nil {
		return nil, sql.ErrConnDone
	}
	stmtSave, err := db.Prepare(QuerySave)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, errPreparingStatement, err)
		return nil, err
	}
	stmtGetByEmail, err := db.Prepare(QueryByEmail)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, errPreparingStatement, err)
		return nil, err
	}
	stmtGetByID, err := db.Prepare(QueryByID)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, errPreparingStatement, err)
		return nil, err
	}
	stmtUpdate, err := db.Prepare(QueryUpdate)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, errPreparingStatement, err)
		return nil, err
	}
	stmtDelete, err := db.Prepare(QueryDelete)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, errPreparingStatement, err)
		return nil, err
	}

	stmtPatch, err := db.Prepare(QueryPatch)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, errPreparingStatement, err)
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
	}, nil
}

func (r *repository) BeginTx(ctx context.Context) (output.Tx, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	return common.NewSQLTx(tx), nil
}
