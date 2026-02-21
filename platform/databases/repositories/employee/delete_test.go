package employee

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestDeleteEmployee_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	r := &repository{db: db}

	mock.ExpectBegin()
	mock.ExpectExec("DELETE").WithArgs("e1").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err = r.DeleteEmployee(context.Background(), nil, "e1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteEmployee_BeginTxError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	r := &repository{db: db}

	mock.ExpectBegin().WillReturnError(errors.New("begin error"))

	err = r.DeleteEmployee(context.Background(), nil, "e1")
	if err == nil {
		t.Fatal("expected error")
	}
}
