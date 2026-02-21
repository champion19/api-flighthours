package employee

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestNewClientRepository_NilDB(t *testing.T) {
	repo, err := NewClientRepository(nil)
	if err != sql.ErrConnDone {
		t.Fatalf("expected sql.ErrConnDone, got %v", err)
	}
	if repo != nil {
		t.Fatal("expected nil repo")
	}
}

func TestNewClientRepository_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	for i := 0; i < 6; i++ {
		mock.ExpectPrepare("INSERT|SELECT|UPDATE|DELETE")
	}

	repo, err := NewClientRepository(db)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo == nil {
		t.Fatal("expected non-nil repo")
	}
}

func TestNewClientRepository_PrepareError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectPrepare("INSERT").WillReturnError(sql.ErrConnDone)

	repo, err := NewClientRepository(db)
	if err == nil {
		t.Fatal("expected error")
	}
	if repo != nil {
		t.Fatal("expected nil repo")
	}
}

func TestClientRepository_BeginTx_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	for i := 0; i < 6; i++ {
		mock.ExpectPrepare("INSERT|SELECT|UPDATE|DELETE")
	}

	repo, err := NewClientRepository(db)
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

func TestClientRepository_BeginTx_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	for i := 0; i < 6; i++ {
		mock.ExpectPrepare("INSERT|SELECT|UPDATE|DELETE")
	}

	repo, err := NewClientRepository(db)
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
