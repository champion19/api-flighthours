package daily_logbook_detail

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/platform/databases/common"
)

type fakeTxDLSave struct{}

func (f *fakeTxDLSave) Commit() error   { return nil }
func (f *fakeTxDLSave) Rollback() error { return nil }

func TestSaveDailyLogbookDetail_InvalidTx(t *testing.T) {
	r := &repository{}
	err := r.SaveDailyLogbookDetail(context.Background(), &fakeTxDLSave{}, domain.DailyLogbookDetail{})
	if !errors.Is(err, domain.ErrInvalidTransaction) {
		t.Fatalf("expected ErrInvalidTransaction, got %v", err)
	}
}

func TestSaveDailyLogbookDetail_ExecError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO daily_logbook_detail")
	stmt, err := db.Prepare("INSERT INTO daily_logbook_detail")
	if err != nil {
		t.Fatal(err)
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO daily_logbook_detail").WillReturnError(errors.New("insert fail"))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{stmtInsert: stmt}

	err = r.SaveDailyLogbookDetail(context.Background(), sqlTx, domain.DailyLogbookDetail{ID: "d-1"})
	if !errors.Is(err, domain.ErrFlightCannotSave) {
		t.Fatalf("expected ErrFlightCannotSave, got %v", err)
	}
}

func TestSaveDailyLogbookDetail_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO daily_logbook_detail")
	stmt, err := db.Prepare("INSERT INTO daily_logbook_detail")
	if err != nil {
		t.Fatal(err)
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO daily_logbook_detail").WillReturnResult(sqlmock.NewResult(1, 1))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{stmtInsert: stmt}

	err = r.SaveDailyLogbookDetail(context.Background(), sqlTx, domain.DailyLogbookDetail{
		ID:             "d-1",
		DailyLogbookID: "lb-1",
		FlightNumber:   "AV-123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
