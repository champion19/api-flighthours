package airport

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/platform/databases/common"
)

type fakeTxAP struct{}

func (f *fakeTxAP) Commit() error   { return nil }
func (f *fakeTxAP) Rollback() error { return nil }

func TestUpdateAirportStatus_InvalidTx(t *testing.T) {
	r := &repository{}
	err := r.UpdateAirportStatus(context.Background(), &fakeTxAP{}, "id-1", true)
	if !errors.Is(err, domain.ErrInvalidTransaction) {
		t.Fatalf("expected ErrInvalidTransaction, got %v", err)
	}
}

func TestUpdateAirportStatus_ExecError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE airport SET status").WillReturnError(errors.New("exec fail"))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.UpdateAirportStatus(context.Background(), sqlTx, "id-1", true)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUpdateAirportStatus_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE airport SET status").WillReturnResult(sqlmock.NewResult(0, 1))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.UpdateAirportStatus(context.Background(), sqlTx, "id-1", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
