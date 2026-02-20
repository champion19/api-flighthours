package daily_logbook_detail

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/platform/databases/common"
)

type fakeTxDLDUpdate struct{}

func (f *fakeTxDLDUpdate) Commit() error   { return nil }
func (f *fakeTxDLDUpdate) Rollback() error { return nil }

func TestUpdateDailyLogbookDetail_InvalidTx(t *testing.T) {
	r := &repository{}
	err := r.UpdateDailyLogbookDetail(context.Background(), &fakeTxDLDUpdate{}, domain.DailyLogbookDetail{})
	if !errors.Is(err, domain.ErrInvalidTransaction) {
		t.Fatalf("expected ErrInvalidTransaction, got %v", err)
	}
}

func TestUpdateDailyLogbookDetail_ExecError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectPrepare("UPDATE daily_logbook_detail")
	stmt, err := db.Prepare("UPDATE daily_logbook_detail SET")
	if err != nil {
		t.Fatal(err)
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE daily_logbook_detail").WillReturnError(errors.New("update fail"))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{stmtUpdate: stmt}

	err = r.UpdateDailyLogbookDetail(context.Background(), sqlTx, domain.DailyLogbookDetail{ID: "d-1"})
	if !errors.Is(err, domain.ErrFlightCannotUpdate) {
		t.Fatalf("expected ErrFlightCannotUpdate, got %v", err)
	}
}

func TestUpdateDailyLogbookDetail_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectPrepare("UPDATE daily_logbook_detail")
	stmt, err := db.Prepare("UPDATE daily_logbook_detail SET")
	if err != nil {
		t.Fatal(err)
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE daily_logbook_detail").WillReturnResult(sqlmock.NewResult(0, 1))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{stmtUpdate: stmt}

	err = r.UpdateDailyLogbookDetail(context.Background(), sqlTx, domain.DailyLogbookDetail{
		ID:             "d-1",
		DailyLogbookID: "lb-1",
		FlightNumber:   "AV-100",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
