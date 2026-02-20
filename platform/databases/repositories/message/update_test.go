package message

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/platform/databases/common"
)

type fakeTxMU struct{}

func (f *fakeTxMU) Commit() error   { return nil }
func (f *fakeTxMU) Rollback() error { return nil }

func TestUpdateMessage_InvalidTx(t *testing.T) {
	r := &repository{}
	err := r.UpdateMessage(context.Background(), &fakeTxMU{}, domain.Message{})
	if !errors.Is(err, domain.ErrInvalidTransaction) {
		t.Fatalf("expected ErrInvalidTransaction, got %v", err)
	}
}

func TestUpdateMessage_ExecError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE system_messages SET").WillReturnError(errors.New("update fail"))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.UpdateMessage(context.Background(), sqlTx, domain.Message{
		ID:      "msg-1",
		Code:    "CODE-1",
		Title:   "Test Title",
		Content: "Test Content",
		Active:  true,
	})
	if !errors.Is(err, domain.ErrMessageCannotUpdate) {
		t.Fatalf("expected ErrMessageCannotUpdate, got %v", err)
	}
}

func TestUpdateMessage_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE system_messages SET").WillReturnResult(sqlmock.NewResult(0, 1))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.UpdateMessage(context.Background(), sqlTx, domain.Message{
		ID:      "msg-1",
		Code:    "CODE-1",
		Title:   "Test Title",
		Content: "Test Content",
		Active:  true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
