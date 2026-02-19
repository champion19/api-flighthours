package common

import (
	"testing"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

func TestCastTx_ValidSQLTX(t *testing.T) {
	tx := NewSQLTx(nil)
	result, err := CastTx(tx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != tx {
		t.Error("expected same SQLTX instance")
	}
}

type fakeTx struct{}

func (f *fakeTx) Commit() error   { return nil }
func (f *fakeTx) Rollback() error { return nil }

func TestCastTx_InvalidType(t *testing.T) {
	_, err := CastTx(&fakeTx{})
	if err == nil {
		t.Fatal("expected error for invalid tx type")
	}
	if err != domain.ErrInvalidTransaction {
		t.Errorf("expected ErrInvalidTransaction, got %v", err)
	}
}

func TestCastTx_NilTx(t *testing.T) {
	_, err := CastTx(nil)
	if err == nil {
		t.Fatal("expected error for nil tx")
	}
}
