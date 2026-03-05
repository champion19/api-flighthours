package message

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/platform/databases/common"
)

type fakeTxMS struct{}

func (f *fakeTxMS) Commit() error   { return nil }
func (f *fakeTxMS) Rollback() error { return nil }

func TestSaveMessage_InvalidTx(t *testing.T) {
	r := &repository{}
	err := r.SaveMessage(context.Background(), &fakeTxMS{}, domain.Message{})
	if !errors.Is(err, domain.ErrInvalidTransaction) {
		t.Fatalf("expected ErrInvalidTransaction, got %v", err)
	}
}

func TestSaveMessage_ExecError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO system_messages").WillReturnError(errors.New("insert fail"))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.SaveMessage(context.Background(), sqlTx, domain.Message{ID: "msg-1"})
	if !errors.Is(err, domain.ErrMessageCannotSave) {
		t.Fatalf("expected ErrMessageCannotSave, got %v", err)
	}
}

func TestSaveMessage_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO system_messages").WillReturnResult(sqlmock.NewResult(1, 1))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.SaveMessage(context.Background(), sqlTx, domain.Message{ID: "msg-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
