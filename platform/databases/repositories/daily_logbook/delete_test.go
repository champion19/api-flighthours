package daily_logbook

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/platform/databases/common"
)

type fakeTxDLD struct{}

func (f *fakeTxDLD) Commit() error   { return nil }
func (f *fakeTxDLD) Rollback() error { return nil }

func TestDeleteDailyLogbook_InvalidTx(t *testing.T) {
	r := &repository{}
	err := r.DeleteDailyLogbook(context.Background(), &fakeTxDLD{}, "id-1")
	if !errors.Is(err, domain.ErrInvalidTransaction) {
		t.Fatalf("expected ErrInvalidTransaction, got %v", err)
	}
}

func TestDeleteDailyLogbook_CascadeDeleteError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM daily_logbook_detail").WillReturnError(errors.New("cascade fail"))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.DeleteDailyLogbook(context.Background(), sqlTx, "id-1")
	if !errors.Is(err, domain.ErrDailyLogbookCannotDelete) {
		t.Fatalf("expected ErrDailyLogbookCannotDelete, got %v", err)
	}
}

func TestDeleteDailyLogbook_MainDeleteError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM daily_logbook_detail").WillReturnResult(sqlmock.NewResult(0, 2))
	mock.ExpectExec("DELETE FROM daily_logbook WHERE").WillReturnError(errors.New("delete fail"))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.DeleteDailyLogbook(context.Background(), sqlTx, "id-1")
	if !errors.Is(err, domain.ErrDailyLogbookCannotDelete) {
		t.Fatalf("expected ErrDailyLogbookCannotDelete, got %v", err)
	}
}

func TestDeleteDailyLogbook_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM daily_logbook_detail").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("DELETE FROM daily_logbook WHERE").WillReturnResult(sqlmock.NewResult(0, 0))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.DeleteDailyLogbook(context.Background(), sqlTx, "id-1")
	if !errors.Is(err, domain.ErrDailyLogbookNotFound) {
		t.Fatalf("expected ErrDailyLogbookNotFound, got %v", err)
	}
}

func TestDeleteDailyLogbook_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM daily_logbook_detail").WillReturnResult(sqlmock.NewResult(0, 3))
	mock.ExpectExec("DELETE FROM daily_logbook WHERE").WillReturnResult(sqlmock.NewResult(0, 1))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.DeleteDailyLogbook(context.Background(), sqlTx, "id-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
