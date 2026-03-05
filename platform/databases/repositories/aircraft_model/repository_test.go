package aircraftmodel

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestNewAircraftModelRepository_NilDB(t *testing.T) {
	repo, err := NewAircraftModelRepository(nil)
	if err != sql.ErrConnDone {
		t.Fatalf("expected sql.ErrConnDone, got %v", err)
	}
	if repo != nil {
		t.Fatal("expected nil repo")
	}
}

func TestNewAircraftModelRepository_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	for i := 0; i < 5; i++ {
		mock.ExpectPrepare("SELECT|UPDATE")
	}

	repo, err := NewAircraftModelRepository(db)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo == nil {
		t.Fatal("expected non-nil repo")
	}
}

func TestNewAircraftModelRepository_PrepareError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectPrepare("SELECT").WillReturnError(sql.ErrConnDone)

	repo, err := NewAircraftModelRepository(db)
	if err == nil {
		t.Fatal("expected error")
	}
	if repo != nil {
		t.Fatal("expected nil repo")
	}
}
