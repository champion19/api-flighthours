package message

import (
	"context"
	"database/sql"

	"github.com/champion19/flighthours-api/core/ports/output"
	cachetypes "github.com/champion19/flighthours-api/platform/cache/types"
	"github.com/champion19/flighthours-api/platform/databases/common"
	"github.com/champion19/flighthours-api/platform/logger"
)

const (
	queryGetAllActive = `
		SELECT id, message_code, type, category, module, message_title, message_content, is_active, created_at, updated_at
		FROM system_messages
		WHERE is_active = true`

	queryGetByCode = `
		SELECT id, message_code, type, category, module, message_title, message_content, is_active, created_at, updated_at
		FROM system_messages
		WHERE message_code = ?
		LIMIT 1`

	queryGetByCodeForCache = `
		SELECT id, message_code, type, category, module, message_title, message_content, is_active, created_at, updated_at
		FROM system_messages
		WHERE message_code = ? AND is_active = true
		LIMIT 1`

	queryGetByCodeWithStatus = `
		SELECT id, message_code, type, category, module, message_title, message_content, is_active, created_at, updated_at
		FROM system_messages
		WHERE message_code = ?
		LIMIT 1`

	queryGetByID = `
		SELECT id, message_code, type, category, module, message_title, message_content, is_active, created_at, updated_at
		FROM system_messages
		WHERE id = ?
		LIMIT 1`

	queryGetByType = `
		SELECT id, message_code, type, category, module, message_title, message_content, is_active, created_at, updated_at
		FROM system_messages
		WHERE type = ? AND is_active = true`

	queryGetByModule = `
		SELECT id, message_code, type, category, module, message_title, message_content, is_active, created_at, updated_at
		FROM system_messages
		WHERE module = ? AND is_active = true`

	queryMessageSave = `INSERT INTO system_messages
		(id, message_code, type, category, module, message_title, message_content, is_active)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	queryMessageDelete = `DELETE FROM system_messages WHERE id = ?`
)

// repository implements output.MessageRepository

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

// NewMessageRepository creates a new message repository
func NewMessageRepository(db *sql.DB) (MessageRepository, error) {
	if db == nil {
		return nil, sql.ErrConnDone
	}

	stmtGetAllActive, err := db.Prepare(queryGetAllActive)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing statement", err)
		return nil, err
	}

	stmtGetByCode, err := db.Prepare(queryGetByCode)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing statement", err)
		return nil, err
	}

	stmtGetByCodeForCache, err := db.Prepare(queryGetByCodeForCache)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing statement", err)
		return nil, err
	}

	stmtGetByCodeWithStatus, err := db.Prepare(queryGetByCodeWithStatus)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing statement", err)
		return nil, err
	}

	stmtGetByID, err := db.Prepare(queryGetByID)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing statement", err)
		return nil, err
	}

	stmtGetByType, err := db.Prepare(queryGetByType)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing statement", err)
		return nil, err
	}

	stmtGetByModule, err := db.Prepare(queryGetByModule)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing statement", err)
		return nil, err
	}

	stmtMessageSave, err := db.Prepare(queryMessageSave)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing statement", err)
		return nil, err
	}

	stmtMessageDelete, err := db.Prepare(queryMessageDelete)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing statement", err)
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
	return common.NewSQLTx(tx), nil
}
