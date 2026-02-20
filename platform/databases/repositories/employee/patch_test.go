package employee

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/platform/databases/common"
)

type fakeTxEP struct{}

func (f *fakeTxEP) Commit() error   { return nil }
func (f *fakeTxEP) Rollback() error { return nil }

func TestPatchEmployee_InvalidTx(t *testing.T) {
	r := &repository{}
	err := r.PatchEmployee(context.Background(), &fakeTxEP{}, "id-1", "kc-1")
	if !errors.Is(err, domain.ErrInvalidTransaction) {
		t.Fatalf("expected ErrInvalidTransaction, got %v", err)
	}
}

func TestPatchEmployee_ExecError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE employee SET").WillReturnError(errors.New("update fail"))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.PatchEmployee(context.Background(), sqlTx, "id-1", "kc-1")
	if !errors.Is(err, domain.ErrUserCannotSave) {
		t.Fatalf("expected ErrUserCannotSave, got %v", err)
	}
}

func TestPatchEmployee_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE employee SET").WillReturnResult(sqlmock.NewResult(0, 0))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.PatchEmployee(context.Background(), sqlTx, "id-1", "kc-1")
	if !errors.Is(err, domain.ErrPersonNotFound) {
		t.Fatalf("expected ErrPersonNotFound, got %v", err)
	}
}

func TestPatchEmployee_Success(t *testing.T) {
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

	err = r.PatchEmployee(context.Background(), sqlTx, "id-1", "kc-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
