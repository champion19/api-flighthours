package message

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/platform/databases/common"
)

type fakeTxMD struct{}

func (f *fakeTxMD) Commit() error   { return nil }
func (f *fakeTxMD) Rollback() error { return nil }

func TestDeleteMessage_InvalidTx(t *testing.T) {
	r := &repository{}
	err := r.DeleteMessage(context.Background(), &fakeTxMD{}, "id-1")
	if !errors.Is(err, domain.ErrInvalidTransaction) {
		t.Fatalf("expected ErrInvalidTransaction, got %v", err)
	}
}

func TestDeleteMessage_ExecError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM system_messages").WillReturnError(errors.New("delete error"))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.DeleteMessage(context.Background(), sqlTx, "id-1")
	if !errors.Is(err, domain.ErrMessageCannotDelete) {
		t.Fatalf("expected ErrMessageCannotDelete, got %v", err)
	}
}

func TestDeleteMessage_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM system_messages").WillReturnResult(sqlmock.NewResult(0, 1))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.DeleteMessage(context.Background(), sqlTx, "id-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
