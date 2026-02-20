package employee

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/platform/databases/common"
	"github.com/go-sql-driver/mysql"
)

type fakeTxEU struct{}

func (f *fakeTxEU) Commit() error   { return nil }
func (f *fakeTxEU) Rollback() error { return nil }

func TestUpdateEmployee_InvalidTx(t *testing.T) {
	r := &repository{}
	err := r.UpdateEmployee(context.Background(), &fakeTxEU{}, domain.Employee{})
	if !errors.Is(err, domain.ErrInvalidTransaction) {
		t.Fatalf("expected ErrInvalidTransaction, got %v", err)
	}
}

func TestUpdateEmployee_ForeignKeyError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE employee SET").WillReturnError(&mysql.MySQLError{Number: 1452, Message: "foreign key constraint"})

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.UpdateEmployee(context.Background(), sqlTx, domain.Employee{ID: "emp-1"})
	if !errors.Is(err, domain.ErrInvalidForeignKey) {
		t.Fatalf("expected ErrInvalidForeignKey, got %v", err)
	}
}

func TestUpdateEmployee_DataTooLong(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE employee SET").WillReturnError(&mysql.MySQLError{Number: 1406, Message: "Data too long"})

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.UpdateEmployee(context.Background(), sqlTx, domain.Employee{ID: "emp-1"})
	if !errors.Is(err, domain.ErrDataTooLong) {
		t.Fatalf("expected ErrDataTooLong, got %v", err)
	}
}

func TestUpdateEmployee_DuplicateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE employee SET").WillReturnError(&mysql.MySQLError{Number: 1062, Message: "Duplicate entry"})

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.UpdateEmployee(context.Background(), sqlTx, domain.Employee{ID: "emp-1"})
	if !errors.Is(err, domain.ErrDuplicateUser) {
		t.Fatalf("expected ErrDuplicateUser, got %v", err)
	}
}

func TestUpdateEmployee_GenericError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE employee SET").WillReturnError(errors.New("generic error"))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.UpdateEmployee(context.Background(), sqlTx, domain.Employee{ID: "emp-1"})
	if !errors.Is(err, domain.ErrUserCannotUpdate) {
		t.Fatalf("expected ErrUserCannotUpdate, got %v", err)
	}
}

func TestUpdateEmployee_Success(t *testing.T) {
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

	err = r.UpdateEmployee(context.Background(), sqlTx, domain.Employee{ID: "emp-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
