package daily_logbook

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/platform/databases/common"
)

type fakeTxDLS struct{}

func (f *fakeTxDLS) Commit() error   { return nil }
func (f *fakeTxDLS) Rollback() error { return nil }

func TestSaveDailyLogbook_InvalidTx(t *testing.T) {
	r := &repository{}
	err := r.SaveDailyLogbook(context.Background(), &fakeTxDLS{}, domain.DailyLogbook{})
	if !errors.Is(err, domain.ErrInvalidTransaction) {
		t.Fatalf("expected ErrInvalidTransaction, got %v", err)
	}
}

func TestSaveDailyLogbook_ExecError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO daily_logbook").WillReturnError(errors.New("insert fail"))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.SaveDailyLogbook(context.Background(), sqlTx, domain.DailyLogbook{ID: "lb-1"})
	if !errors.Is(err, domain.ErrDailyLogbookCannotSave) {
		t.Fatalf("expected ErrDailyLogbookCannotSave, got %v", err)
	}
}

func TestSaveDailyLogbook_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO daily_logbook").WillReturnResult(sqlmock.NewResult(1, 1))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.SaveDailyLogbook(context.Background(), sqlTx, domain.DailyLogbook{ID: "lb-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
