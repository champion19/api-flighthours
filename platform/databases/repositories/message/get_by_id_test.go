package message

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	r := &repository{db: db}

	cols := []string{"id", "message_code", "type", "category", "module", "message_title", "message_content", "is_active", "created_at", "updated_at"}
	mock.ExpectQuery("SELECT").WithArgs("m1").WillReturnRows(
		sqlmock.NewRows(cols).AddRow("m1", "MSG_001", "error", "validation", "auth", "Title", "Content", true, time.Now(), time.Now()),
	)

	result, err := r.GetByID(context.Background(), "m1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result.Code != "MSG_001" {
		t.Errorf("expected MSG_001, got %s", result.Code)
	}
}

func TestGetByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	r := &repository{db: db}

	mock.ExpectQuery("SELECT").WithArgs("unknown").WillReturnError(sql.ErrNoRows)

	result, err := r.GetByID(context.Background(), "unknown")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if result != nil {
		t.Fatal("expected nil result")
	}
}

func TestGetByID_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	r := &repository{db: db}

	mock.ExpectQuery("SELECT").WithArgs("m1").WillReturnError(errors.New("db error"))

	_, err = r.GetByID(context.Background(), "m1")
	if err == nil {
		t.Fatal("expected error")
	}
}
