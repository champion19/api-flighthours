package daily_logbook_detail

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/platform/databases/common"
)

type fakeTxDLDel struct{}

func (f *fakeTxDLDel) Commit() error   { return nil }
func (f *fakeTxDLDel) Rollback() error { return nil }

func TestDeleteDailyLogbookDetail_InvalidTx(t *testing.T) {
	r := &repository{}
	err := r.DeleteDailyLogbookDetail(context.Background(), &fakeTxDLDel{}, "id-1")
	if !errors.Is(err, domain.ErrInvalidTransaction) {
		t.Fatalf("expected ErrInvalidTransaction, got %v", err)
	}
}

func TestDeleteDailyLogbookDetail_ExecError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Prepare the statement that the repository will use
	mock.ExpectPrepare("DELETE FROM daily_logbook_detail")
	stmt, err := db.Prepare("DELETE FROM daily_logbook_detail WHERE id = ?")
	if err != nil {
		t.Fatal(err)
	}

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM daily_logbook_detail").WillReturnError(errors.New("exec fail"))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{stmtDelete: stmt}

	err = r.DeleteDailyLogbookDetail(context.Background(), sqlTx, "id-1")
	if !errors.Is(err, domain.ErrFlightCannotDelete) {
		t.Fatalf("expected ErrFlightCannotDelete, got %v", err)
	}
}

func TestDeleteDailyLogbookDetail_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectPrepare("DELETE FROM daily_logbook_detail")
	stmt, err := db.Prepare("DELETE FROM daily_logbook_detail WHERE id = ?")
	if err != nil {
		t.Fatal(err)
	}

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM daily_logbook_detail").WillReturnResult(sqlmock.NewResult(0, 0))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{stmtDelete: stmt}

	err = r.DeleteDailyLogbookDetail(context.Background(), sqlTx, "id-1")
	if !errors.Is(err, domain.ErrFlightNotFound) {
		t.Fatalf("expected ErrFlightNotFound, got %v", err)
	}
}

func TestDeleteDailyLogbookDetail_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectPrepare("DELETE FROM daily_logbook_detail")
	stmt, err := db.Prepare("DELETE FROM daily_logbook_detail WHERE id = ?")
	if err != nil {
		t.Fatal(err)
	}

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM daily_logbook_detail").WillReturnResult(sqlmock.NewResult(0, 1))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{stmtDelete: stmt}

	err = r.DeleteDailyLogbookDetail(context.Background(), sqlTx, "id-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewDailyLogbookDetailRepository_NilDB(t *testing.T) {
	_, err := NewDailyLogbookDetailRepository(nil)
	if err != sql.ErrConnDone {
		t.Fatalf("expected sql.ErrConnDone, got %v", err)
	}
}

func TestBeginTx_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()

	r := &repository{db: db}
	tx, err := r.BeginTx(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tx == nil {
		t.Fatal("expected non-nil tx")
	}
}

func TestBeginTx_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin().WillReturnError(errors.New("begin fail"))

	r := &repository{db: db}
	_, err = r.BeginTx(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}
