package aircraftmodel

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/platform/databases/common"
)

type fakeTxAM struct{}

func (f *fakeTxAM) Commit() error   { return nil }
func (f *fakeTxAM) Rollback() error { return nil }

// --- BeginTx ---

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

// --- UpdateAircraftModelStatus ---

func TestUpdateAircraftModelStatus_InvalidTx(t *testing.T) {
	r := &repository{}
	err := r.UpdateAircraftModelStatus(context.Background(), &fakeTxAM{}, "id-1", true)
	if !errors.Is(err, domain.ErrInvalidTransaction) {
		t.Fatalf("expected ErrInvalidTransaction, got %v", err)
	}
}

func TestUpdateAircraftModelStatus_ExecError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE aircraft_model SET status").WillReturnError(errors.New("exec fail"))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.UpdateAircraftModelStatus(context.Background(), sqlTx, "id-1", true)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUpdateAircraftModelStatus_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE aircraft_model SET status").WillReturnResult(sqlmock.NewResult(0, 1))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.UpdateAircraftModelStatus(context.Background(), sqlTx, "id-1", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
