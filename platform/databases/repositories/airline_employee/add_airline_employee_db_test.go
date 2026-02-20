package airline_employee

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/platform/databases/common"
)

type fakeTxAE struct{}

func (f *fakeTxAE) Commit() error   { return nil }
func (f *fakeTxAE) Rollback() error { return nil }

func TestAddAirlineEmployee_InvalidTx(t *testing.T) {
	r := &repository{}
	err := r.AddAirlineEmployee(context.Background(), &fakeTxAE{}, domain.AirlineEmployee{})
	if !errors.Is(err, domain.ErrInvalidTransaction) {
		t.Fatalf("expected ErrInvalidTransaction, got %v", err)
	}
}

func TestAddAirlineEmployee_ExecError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE employee SET").WillReturnError(errors.New("db error"))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.AddAirlineEmployee(context.Background(), sqlTx, domain.AirlineEmployee{ID: "emp-1"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAddAirlineEmployee_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE employee SET").WillReturnResult(sqlmock.NewResult(0, 1))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.AddAirlineEmployee(context.Background(), sqlTx, domain.AirlineEmployee{ID: "emp-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
