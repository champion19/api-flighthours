package airline_employee

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/platform/databases/common"
)

type fakeTxAEU struct{}

func (f *fakeTxAEU) Commit() error   { return nil }
func (f *fakeTxAEU) Rollback() error { return nil }

// --- UpdateAirlineEmployee ---

func TestUpdateAirlineEmployee_InvalidTx(t *testing.T) {
	r := &repository{}
	err := r.UpdateAirlineEmployee(context.Background(), &fakeTxAEU{}, domain.AirlineEmployee{})
	if !errors.Is(err, domain.ErrInvalidTransaction) {
		t.Fatalf("expected ErrInvalidTransaction, got %v", err)
	}
}

func TestUpdateAirlineEmployee_ExecError(t *testing.T) {
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

	err = r.UpdateAirlineEmployee(context.Background(), sqlTx, domain.AirlineEmployee{ID: "emp-1"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUpdateAirlineEmployee_Success(t *testing.T) {
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

	err = r.UpdateAirlineEmployee(context.Background(), sqlTx, domain.AirlineEmployee{ID: "emp-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// --- UpdateAirlineEmployeeStatus ---

func TestUpdateAirlineEmployeeStatus_InvalidTx(t *testing.T) {
	r := &repository{}
	err := r.UpdateAirlineEmployeeStatus(context.Background(), &fakeTxAEU{}, "id-1", true)
	if !errors.Is(err, domain.ErrInvalidTransaction) {
		t.Fatalf("expected ErrInvalidTransaction, got %v", err)
	}
}

func TestUpdateAirlineEmployeeStatus_EmployeeNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT airline IS NOT NULL FROM employee").
		WithArgs("id-1").
		WillReturnError(errors.New("no rows"))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.UpdateAirlineEmployeeStatus(context.Background(), sqlTx, "id-1", true)
	if !errors.Is(err, domain.ErrAirlineEmployeeNotFound) {
		t.Fatalf("expected ErrAirlineEmployeeNotFound, got %v", err)
	}
}

func TestUpdateAirlineEmployeeStatus_NoAirline(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT airline IS NOT NULL FROM employee").
		WithArgs("id-1").
		WillReturnRows(sqlmock.NewRows([]string{"has_airline"}).AddRow(false))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.UpdateAirlineEmployeeStatus(context.Background(), sqlTx, "id-1", true)
	if !errors.Is(err, domain.ErrAirlineEmployeeNotFound) {
		t.Fatalf("expected ErrAirlineEmployeeNotFound, got %v", err)
	}
}

func TestUpdateAirlineEmployeeStatus_ExecError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT airline IS NOT NULL FROM employee").
		WithArgs("id-1").
		WillReturnRows(sqlmock.NewRows([]string{"has_airline"}).AddRow(true))
	mock.ExpectExec("UPDATE employee SET active").WillReturnError(errors.New("exec fail"))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.UpdateAirlineEmployeeStatus(context.Background(), sqlTx, "id-1", true)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUpdateAirlineEmployeeStatus_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT airline IS NOT NULL FROM employee").
		WithArgs("id-1").
		WillReturnRows(sqlmock.NewRows([]string{"has_airline"}).AddRow(true))
	mock.ExpectExec("UPDATE employee SET active").WillReturnResult(sqlmock.NewResult(0, 1))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.UpdateAirlineEmployeeStatus(context.Background(), sqlTx, "id-1", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
