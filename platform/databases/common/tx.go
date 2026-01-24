package common

import (
	"context"
	"database/sql"
)
type SQLTX struct {
	tx     *sql.Tx
	closed bool
}

func NewSQLTx(tx *sql.Tx) *SQLTX {
	return &SQLTX{
		tx:     tx,
		closed: false,
	}
}

// Commit commits the transaction
func (t *SQLTX) Commit() error {
	if t.closed {
		panic("sqlTx: commit on closed transaction")
	}
	t.closed = true
	return t.tx.Commit()
}

// Rollback rolls back the transaction
func (t *SQLTX) Rollback() error {
	if t.closed {
		panic("sqlTx: rollback on closed transaction")
	}
	t.closed = true
	return t.tx.Rollback()
}

// ExecContext executes a query within the transaction
func (t *SQLTX) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return t.tx	.ExecContext(ctx, query, args...)
}

// QueryRowContext queries a single row within the transaction
func (t *SQLTX) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return t.tx.QueryRowContext(ctx, query, args...)
}
