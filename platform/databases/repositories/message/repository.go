package message

import (
	"context"
	"database/sql"

	"github.com/champion19/api-flighthours/core/ports/output"
	cachetypes "github.com/champion19/api-flighthours/platform/cache/types"
	"github.com/champion19/api-flighthours/platform/databases/common"
	"github.com/champion19/api-flighthours/platform/logger"
)

const (
	// Shared SELECT columns for message queries
	selectMessageColumns = `id, message_code, type, category, module, message_title, message_content, is_active, created_at, updated_at`

	queryGetAllActive = `SELECT ` + selectMessageColumns + ` FROM system_messages WHERE is_active = true`

	queryGetByCode = `SELECT ` + selectMessageColumns + ` FROM system_messages WHERE message_code = ? LIMIT 1`

	queryGetByCodeForCache = `SELECT ` + selectMessageColumns + ` FROM system_messages WHERE message_code = ? AND is_active = true LIMIT 1`

	queryGetByCodeWithStatus = `SELECT ` + selectMessageColumns + ` FROM system_messages WHERE message_code = ? LIMIT 1`

	queryGetByID = `SELECT ` + selectMessageColumns + ` FROM system_messages WHERE id = ? LIMIT 1`

	queryGetByType = `SELECT ` + selectMessageColumns + ` FROM system_messages WHERE type = ? AND is_active = true`

	queryGetByModule = `SELECT ` + selectMessageColumns + ` FROM system_messages WHERE module = ? AND is_active = true`

	queryMessageSave = `INSERT INTO system_messages
		(id, message_code, type, category, module, message_title, message_content, is_active)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	queryMessageDelete = `DELETE FROM system_messages WHERE id = ?`
)

var log logger.Logger = logger.NewSlogLogger()

type repository struct {
	stmtGetAllActive        *sql.Stmt
	stmtGetByCode           *sql.Stmt
	stmtGetByCodeForCache   *sql.Stmt
	stmtGetByCodeWithStatus *sql.Stmt
	stmtGetByID             *sql.Stmt
	stmtGetByType           *sql.Stmt
	stmtGetByModule         *sql.Stmt
	stmtMessageSave         *sql.Stmt
	stmtMessageDelete       *sql.Stmt
	db                      *sql.DB
}

type MessageRepository interface {
	output.MessageRepository
	cachetypes.MessageCacheRepository
}

func NewMessageRepository(db *sql.DB) (MessageRepository, error) {
	if db == nil {
		return nil, sql.ErrConnDone
	}

	prepare := func(query string) (*sql.Stmt, error) {
		stmt, err := db.Prepare(query)
		if err != nil {
			log.Error(logger.LogDatabaseUnavailable, "error preparing statement", err)
		}
		return stmt, err
	}

	stmtGetAllActive, err := prepare(queryGetAllActive)
	if err != nil {
		return nil, err
	}
	stmtGetByCode, err := prepare(queryGetByCode)
	if err != nil {
		return nil, err
	}
	stmtGetByCodeForCache, err := prepare(queryGetByCodeForCache)
	if err != nil {
		return nil, err
	}
	stmtGetByCodeWithStatus, err := prepare(queryGetByCodeWithStatus)
	if err != nil {
		return nil, err
	}
	stmtGetByID, err := prepare(queryGetByID)
	if err != nil {
		return nil, err
	}
	stmtGetByType, err := prepare(queryGetByType)
	if err != nil {
		return nil, err
	}
	stmtGetByModule, err := prepare(queryGetByModule)
	if err != nil {
		return nil, err
	}
	stmtMessageSave, err := prepare(queryMessageSave)
	if err != nil {
		return nil, err
	}
	stmtMessageDelete, err := prepare(queryMessageDelete)
	if err != nil {
		return nil, err
	}

	return &repository{
		db:                      db,
		stmtGetAllActive:        stmtGetAllActive,
		stmtGetByCode:           stmtGetByCode,
		stmtGetByCodeForCache:   stmtGetByCodeForCache,
		stmtGetByCodeWithStatus: stmtGetByCodeWithStatus,
		stmtGetByID:             stmtGetByID,
		stmtGetByType:           stmtGetByType,
		stmtGetByModule:         stmtGetByModule,
		stmtMessageSave:         stmtMessageSave,
		stmtMessageDelete:       stmtMessageDelete,
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
