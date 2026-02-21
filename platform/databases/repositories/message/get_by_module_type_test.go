package message

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetByModule_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	r := &repository{db: db}

	mock.ExpectQuery("SELECT").WithArgs("auth").WillReturnError(errors.New("db error"))

	_, err = r.GetByModule(context.Background(), "auth")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetByModule_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	r := &repository{db: db}

	cols := []string{"id", "message_code", "type", "category", "module", "message_title", "message_content", "is_active", "created_at", "updated_at"}
	rows := sqlmock.NewRows(cols)
	mock.ExpectQuery("SELECT").WithArgs("auth").WillReturnRows(rows)

	result, err := r.GetByModule(context.Background(), "auth")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Fatalf("expected empty, got %d", len(result))
	}
}

func TestGetByType_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	r := &repository{db: db}

	mock.ExpectQuery("SELECT").WithArgs("error").WillReturnError(errors.New("db error"))

	_, err = r.GetByType(context.Background(), "error")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetByType_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	r := &repository{db: db}

	cols := []string{"id", "message_code", "type", "category", "module", "message_title", "message_content", "is_active", "created_at", "updated_at"}
	rows := sqlmock.NewRows(cols)
	mock.ExpectQuery("SELECT").WithArgs("error").WillReturnRows(rows)

	result, err := r.GetByType(context.Background(), "error")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Fatalf("expected empty, got %d", len(result))
	}
}
