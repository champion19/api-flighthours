package engine

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestNewEngineRepository_NilDB(t *testing.T) {
	repo, err := NewEngineRepository(nil)
	if err != sql.ErrConnDone {
		t.Fatalf("expected sql.ErrConnDone, got %v", err)
	}
	if repo != nil {
		t.Fatal("expected nil repo")
	}
}

func TestNewEngineRepository_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectPrepare("SELECT")
	mock.ExpectPrepare("SELECT")

	repo, err := NewEngineRepository(db)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo == nil {
		t.Fatal("expected non-nil repo")
	}
}

func TestNewEngineRepository_PrepareError_First(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectPrepare("SELECT").WillReturnError(sql.ErrConnDone)

	repo, err := NewEngineRepository(db)
	if err == nil {
		t.Fatal("expected error")
	}
	if repo != nil {
		t.Fatal("expected nil repo")
	}
}

func TestNewEngineRepository_PrepareError_Second(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectPrepare("SELECT")
	mock.ExpectPrepare("SELECT").WillReturnError(sql.ErrConnDone)

	repo, err := NewEngineRepository(db)
	if err == nil {
		t.Fatal("expected error")
	}
	if repo != nil {
		t.Fatal("expected nil repo")
	}
}
