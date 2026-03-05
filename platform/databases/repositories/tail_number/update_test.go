package tailnumber

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/platform/databases/common"
)

type fakeTxLPU struct{}

func (f *fakeTxLPU) Commit() error   { return nil }
func (f *fakeTxLPU) Rollback() error { return nil }

func TestUpdateTailNumber_InvalidTx(t *testing.T) {
	r := &repository{}
	err := r.UpdateTailNumber(context.Background(), &fakeTxLPU{}, domain.TailNumber{})
	if !errors.Is(err, domain.ErrInvalidTransaction) {
		t.Fatalf("expected ErrInvalidTransaction, got %v", err)
	}
}

func TestUpdateTailNumber_DuplicateEntry(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE tail_number SET").WillReturnError(errors.New("Duplicate entry 'HK-5432' for key"))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.UpdateTailNumber(context.Background(), sqlTx, domain.TailNumber{ID: "lp-1", TailNumber: "HK-5432"})
	if !errors.Is(err, domain.ErrTailNumberDuplicatePlate) {
		t.Fatalf("expected ErrTailNumberDuplicatePlate, got %v", err)
	}
}

func TestUpdateTailNumber_FKAircraftModel(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE tail_number SET").WillReturnError(errors.New("foreign key constraint fails aircraft_model"))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.UpdateTailNumber(context.Background(), sqlTx, domain.TailNumber{ID: "lp-1"})
	if !errors.Is(err, domain.ErrTailNumberInvalidModel) {
		t.Fatalf("expected ErrTailNumberInvalidModel, got %v", err)
	}
}

func TestUpdateTailNumber_FKAirline(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE tail_number SET").WillReturnError(errors.New("foreign key constraint fails airline"))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.UpdateTailNumber(context.Background(), sqlTx, domain.TailNumber{ID: "lp-1"})
	if !errors.Is(err, domain.ErrTailNumberInvalidAirline) {
		t.Fatalf("expected ErrTailNumberInvalidAirline, got %v", err)
	}
}

func TestUpdateTailNumber_GenericError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE tail_number SET").WillReturnError(errors.New("unknown error"))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.UpdateTailNumber(context.Background(), sqlTx, domain.TailNumber{ID: "lp-1"})
	if !errors.Is(err, domain.ErrTailNumberCannotUpdate) {
		t.Fatalf("expected ErrTailNumberCannotUpdate, got %v", err)
	}
}

func TestUpdateTailNumber_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE tail_number SET").WillReturnResult(sqlmock.NewResult(0, 0))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.UpdateTailNumber(context.Background(), sqlTx, domain.TailNumber{ID: "lp-1"})
	if !errors.Is(err, domain.ErrTailNumberNotFound) {
		t.Fatalf("expected ErrTailNumberNotFound, got %v", err)
	}
}

func TestUpdateTailNumber_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE tail_number SET").WillReturnResult(sqlmock.NewResult(0, 1))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.UpdateTailNumber(context.Background(), sqlTx, domain.TailNumber{
		ID:              "lp-1",
		TailNumber:    "HK-1234",
		AircraftModelID: "am-1",
		AirlineID:       "al-1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
