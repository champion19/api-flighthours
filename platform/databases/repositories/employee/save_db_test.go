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

type fakeTxES struct{}

func (f *fakeTxES) Commit() error   { return nil }
func (f *fakeTxES) Rollback() error { return nil }

func TestSave_InvalidTx(t *testing.T) {
	r := &repository{}
	err := r.Save(context.Background(), &fakeTxES{}, domain.Employee{})
	if !errors.Is(err, domain.ErrInvalidTransaction) {
		t.Fatalf("expected ErrInvalidTransaction, got %v", err)
	}
}

func TestSave_DuplicateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO employee").WillReturnError(&mysql.MySQLError{Number: 1062, Message: "Duplicate entry"})

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.Save(context.Background(), sqlTx, domain.Employee{ID: "emp-1"})
	if !errors.Is(err, domain.ErrDuplicateUser) {
		t.Fatalf("expected ErrDuplicateUser, got %v", err)
	}
}

func TestSave_GenericError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO employee").WillReturnError(errors.New("unexpected error"))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.Save(context.Background(), sqlTx, domain.Employee{ID: "emp-1"})
	if !errors.Is(err, domain.ErrUserCannotSave) {
		t.Fatalf("expected ErrUserCannotSave, got %v", err)
	}
}

func TestSave_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO employee").WillReturnResult(sqlmock.NewResult(1, 1))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.Save(context.Background(), sqlTx, domain.Employee{ID: "emp-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
