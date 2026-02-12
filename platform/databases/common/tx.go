package common

import (
	"context"
	"database/sql"
)
type SQLTX struct {
	Tx     *sql.Tx
	closed bool
}

func NewSQLTx(tx *sql.Tx) *SQLTX {
	return &SQLTX{
		Tx:     tx,
		closed: false,
	}
}

func (t *SQLTX) Commit() error {
	if t.closed {
		panic("sqlTx: commit on closed transaction")
	}
	t.closed = true
	return t.Tx.Commit()
}

func (t *SQLTX) Rollback() error {
	if t.closed {
		panic("sqlTx: rollback on closed transaction")
	}
	t.closed = true
	return t.Tx.Rollback()
}

func (t *SQLTX) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return t.Tx.ExecContext(ctx, query, args...)
}
func (t *SQLTX) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return t.Tx.QueryRowContext(ctx, query, args...)
}
