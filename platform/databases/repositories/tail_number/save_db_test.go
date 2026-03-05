package tailnumber

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/platform/databases/common"
)

type fakeTxLP struct{}

func (f *fakeTxLP) Commit() error   { return nil }
func (f *fakeTxLP) Rollback() error { return nil }

func TestSaveTailNumber_InvalidTx(t *testing.T) {
	r := &repository{}
	err := r.SaveTailNumber(context.Background(), &fakeTxLP{}, domain.TailNumber{})
	if !errors.Is(err, domain.ErrInvalidTransaction) {
		t.Fatalf("expected ErrInvalidTransaction, got %v", err)
	}
}

func TestSaveTailNumber_DuplicateEntry(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO tail_number").WillReturnError(errors.New("Duplicate entry 'HK-5432' for key"))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.SaveTailNumber(context.Background(), sqlTx, domain.TailNumber{ID: "lp-1", TailNumber: "HK-5432"})
	if !errors.Is(err, domain.ErrTailNumberDuplicatePlate) {
		t.Fatalf("expected ErrTailNumberDuplicatePlate, got %v", err)
	}
}

func TestSaveTailNumber_FKAircraftModel(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO tail_number").WillReturnError(errors.New("foreign key constraint fails aircraft_model"))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.SaveTailNumber(context.Background(), sqlTx, domain.TailNumber{ID: "lp-1"})
	if !errors.Is(err, domain.ErrTailNumberInvalidModel) {
		t.Fatalf("expected ErrTailNumberInvalidModel, got %v", err)
	}
}

func TestSaveTailNumber_FKAirline(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO tail_number").WillReturnError(errors.New("foreign key constraint fails airline"))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.SaveTailNumber(context.Background(), sqlTx, domain.TailNumber{ID: "lp-1"})
	if !errors.Is(err, domain.ErrTailNumberInvalidAirline) {
		t.Fatalf("expected ErrTailNumberInvalidAirline, got %v", err)
	}
}

func TestSaveTailNumber_GenericError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO tail_number").WillReturnError(errors.New("unknown error"))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.SaveTailNumber(context.Background(), sqlTx, domain.TailNumber{ID: "lp-1"})
	if !errors.Is(err, domain.ErrTailNumberCannotSave) {
		t.Fatalf("expected ErrTailNumberCannotSave, got %v", err)
	}
}

func TestSaveTailNumber_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO tail_number").WillReturnResult(sqlmock.NewResult(1, 1))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.SaveTailNumber(context.Background(), sqlTx, domain.TailNumber{
		ID:              "lp-1",
		TailNumber:    "HK-1234",
		AircraftModelID: "am-1",
		AirlineID:       "al-1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
