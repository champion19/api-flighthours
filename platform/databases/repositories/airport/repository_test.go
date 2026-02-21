package airport

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestNewAirportRepository_NilDB(t *testing.T) {
	repo, err := NewAirportRepository(nil)
	if err != sql.ErrConnDone {
		t.Fatalf("expected sql.ErrConnDone, got %v", err)
	}
	if repo != nil {
		t.Fatal("expected nil repo")
	}
}

func TestNewAirportRepository_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	for i := 0; i < 5; i++ {
		mock.ExpectPrepare("SELECT|UPDATE")
	}

	repo, err := NewAirportRepository(db)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo == nil {
		t.Fatal("expected non-nil repo")
	}
}

func TestNewAirportRepository_PrepareError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectPrepare("SELECT").WillReturnError(sql.ErrConnDone)

	repo, err := NewAirportRepository(db)
	if err == nil {
		t.Fatal("expected error")
	}
	if repo != nil {
		t.Fatal("expected nil repo")
	}
}

func TestAirportRepository_BeginTx_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	for i := 0; i < 5; i++ {
		mock.ExpectPrepare("SELECT|UPDATE")
	}

	repo, err := NewAirportRepository(db)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mock.ExpectBegin()
	tx, err := repo.BeginTx(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tx == nil {
		t.Fatal("expected non-nil tx")
	}
}

func TestAirportRepository_BeginTx_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	for i := 0; i < 5; i++ {
		mock.ExpectPrepare("SELECT|UPDATE")
	}

	repo, err := NewAirportRepository(db)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mock.ExpectBegin().WillReturnError(sql.ErrConnDone)
	tx, err := repo.BeginTx(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	if tx != nil {
		t.Fatal("expected nil tx")
	}
}
