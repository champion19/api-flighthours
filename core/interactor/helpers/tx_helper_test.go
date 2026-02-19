package helpers

import (
	"context"
	"errors"
	"testing"

	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/logger"
)

// fakeTxBeginner implements TxBeginner for testing
type fakeTxBeginner struct {
	beginTxFn func(ctx context.Context) (output.Tx, error)
}

func (f *fakeTxBeginner) BeginTx(ctx context.Context) (output.Tx, error) {
	if f.beginTxFn != nil {
		return f.beginTxFn(ctx)
	}
	return &fakeTxHelper{}, nil
}

type fakeTxHelper struct {
	commitFn   func() error
	rollbackFn func() error
}

func (f *fakeTxHelper) Commit() error {
	if f.commitFn != nil {
		return f.commitFn()
	}
	return nil
}

func (f *fakeTxHelper) Rollback() error {
	if f.rollbackFn != nil {
		return f.rollbackFn()
	}
	return nil
}

func TestRunWithTx(t *testing.T) {
	log := logger.NewSlogLogger()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		tb := &fakeTxBeginner{}
		err := RunWithTx(ctx, tb, log, "test-error",
			func(ctx context.Context, tx output.Tx) error {
				return nil
			})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("begin tx error", func(t *testing.T) {
		tb := &fakeTxBeginner{
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return nil, errors.New("begin failed")
			},
		}
		err := RunWithTx(ctx, tb, log, "test-error",
			func(ctx context.Context, tx output.Tx) error {
				return nil
			})
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("fn error triggers rollback", func(t *testing.T) {
		rolledBack := false
		tb := &fakeTxBeginner{
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return &fakeTxHelper{
					rollbackFn: func() error {
						rolledBack = true
						return nil
					},
				}, nil
			},
		}
		err := RunWithTx(ctx, tb, log, "test-error",
			func(ctx context.Context, tx output.Tx) error {
				return errors.New("fn failed")
			})
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !rolledBack {
			t.Error("expected rollback to be called")
		}
	})

	t.Run("commit error triggers rollback", func(t *testing.T) {
		rolledBack := false
		tb := &fakeTxBeginner{
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return &fakeTxHelper{
					commitFn: func() error {
						return errors.New("commit failed")
					},
					rollbackFn: func() error {
						rolledBack = true
						return nil
					},
				}, nil
			},
		}
		err := RunWithTx(ctx, tb, log, "test-error",
			func(ctx context.Context, tx output.Tx) error {
				return nil
			})
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !rolledBack {
			t.Error("expected rollback to be called on commit error")
		}
	})

	t.Run("rollback error is logged but original error returned", func(t *testing.T) {
		tb := &fakeTxBeginner{
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return &fakeTxHelper{
					rollbackFn: func() error {
						return errors.New("rollback failed")
					},
				}, nil
			},
		}
		fnErr := errors.New("fn failed")
		err := RunWithTx(ctx, tb, log, "test-error",
			func(ctx context.Context, tx output.Tx) error {
				return fnErr
			})
		if err != fnErr {
			t.Errorf("expected original error %v, got %v", fnErr, err)
		}
	})
}
