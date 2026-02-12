package common

import (
	"database/sql"
	"testing"
)

func TestNewSQLTx(t *testing.T) {
	t.Run("creates new SQLTX wrapper", func(t *testing.T) {
		// Can't create real sql.Tx without connection, just verify nil handling
		tx := NewSQLTx(nil)
		if tx == nil {
			t.Error("expected non-nil SQLTX")
		}
		if tx.closed != false {
			t.Error("expected closed to be false initially")
		}
	})
}

func TestSQLTX_CommitPanicsWhenClosed(t *testing.T) {
	t.Run("panics when committing closed transaction", func(t *testing.T) {
		tx := &SQLTX{Tx: nil, closed: true}

		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic when committing closed transaction")
			}
		}()

		_ = tx.Commit()
	})
}

func TestSQLTX_RollbackPanicsWhenClosed(t *testing.T) {
	t.Run("panics when rolling back closed transaction", func(t *testing.T) {
		tx := &SQLTX{Tx: nil, closed: true}

		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic when rolling back closed transaction")
			}
		}()

		_ = tx.Rollback()
	})
}

// Note: ExecContext and QueryRowContext require a real *sql.Tx to test
// which requires a database connection. These are better tested via
// integration tests.
func TestSQLTX_ClosedFieldAccessible(t *testing.T) {
	t.Run("closed field is accessible", func(t *testing.T) {
		tx := NewSQLTx(nil)
		if tx.closed {
			t.Error("new transaction should not be closed")
		}
	})
}

// Mock test for ErrConnDone constant
func TestSQLErrConnDone(t *testing.T) {
	t.Run("sql.ErrConnDone is defined", func(t *testing.T) {
		if sql.ErrConnDone == nil {
			t.Error("expected sql.ErrConnDone to be non-nil")
		}
	})
}
