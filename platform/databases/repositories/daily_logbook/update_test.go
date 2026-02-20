package daily_logbook

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/platform/databases/common"
)

type fakeTxDLU struct{}

func (f *fakeTxDLU) Commit() error   { return nil }
func (f *fakeTxDLU) Rollback() error { return nil }

// --- UpdateDailyLogbook ---

func TestUpdateDailyLogbook_InvalidTx(t *testing.T) {
	r := &repository{}
	err := r.UpdateDailyLogbook(context.Background(), &fakeTxDLU{}, domain.DailyLogbook{})
	if !errors.Is(err, domain.ErrInvalidTransaction) {
		t.Fatalf("expected ErrInvalidTransaction, got %v", err)
	}
}

func TestUpdateDailyLogbook_ExecError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE daily_logbook SET").WillReturnError(errors.New("update fail"))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.UpdateDailyLogbook(context.Background(), sqlTx, domain.DailyLogbook{ID: "lb-1"})
	if !errors.Is(err, domain.ErrDailyLogbookCannotUpdate) {
		t.Fatalf("expected ErrDailyLogbookCannotUpdate, got %v", err)
	}
}

func TestUpdateDailyLogbook_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE daily_logbook SET").WillReturnResult(sqlmock.NewResult(0, 0))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.UpdateDailyLogbook(context.Background(), sqlTx, domain.DailyLogbook{ID: "lb-1"})
	if !errors.Is(err, domain.ErrDailyLogbookNotFound) {
		t.Fatalf("expected ErrDailyLogbookNotFound, got %v", err)
	}
}

func TestUpdateDailyLogbook_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE daily_logbook SET").WillReturnResult(sqlmock.NewResult(0, 1))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.UpdateDailyLogbook(context.Background(), sqlTx, domain.DailyLogbook{ID: "lb-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// --- UpdateDailyLogbookStatus ---

func TestUpdateDailyLogbookStatus_InvalidTx(t *testing.T) {
	r := &repository{}
	err := r.UpdateDailyLogbookStatus(context.Background(), &fakeTxDLU{}, "id-1", true)
	if !errors.Is(err, domain.ErrInvalidTransaction) {
		t.Fatalf("expected ErrInvalidTransaction, got %v", err)
	}
}

func TestUpdateDailyLogbookStatus_ExecError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE daily_logbook SET status").WillReturnError(errors.New("exec fail"))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.UpdateDailyLogbookStatus(context.Background(), sqlTx, "id-1", true)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUpdateDailyLogbookStatus_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE daily_logbook SET status").WillReturnResult(sqlmock.NewResult(0, 1))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.UpdateDailyLogbookStatus(context.Background(), sqlTx, "id-1", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
