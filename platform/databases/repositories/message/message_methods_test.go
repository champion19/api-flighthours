package message

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

var msgCols = []string{"id", "message_code", "type", "category", "module", "message_title", "message_content", "is_active", "created_at", "updated_at"}

// --- GetAllActive ---

func TestGetAllActive_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	r := &repository{db: db}

	rows := sqlmock.NewRows(msgCols).
		AddRow("m1", "MSG_001", "error", "validation", "auth", "Title", "Content", true, time.Now(), time.Now())
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	result, err := r.GetAllActive(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
}

func TestGetAllActive_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	r := &repository{db: db}

	mock.ExpectQuery("SELECT").WillReturnError(errors.New("db error"))

	_, err = r.GetAllActive(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

// --- GetAllActiveForCache ---

func TestGetAllActiveForCache_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	r := &repository{db: db}

	rows := sqlmock.NewRows(msgCols).
		AddRow("m1", "MSG_001", "error", "validation", "auth", "Title", "Content", true, time.Now(), time.Now())
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	result, err := r.GetAllActiveForCache(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
	if result[0].Code != "MSG_001" {
		t.Errorf("expected MSG_001, got %s", result[0].Code)
	}
}

func TestGetAllActiveForCache_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	r := &repository{db: db}

	mock.ExpectQuery("SELECT").WillReturnError(errors.New("db error"))

	_, err = r.GetAllActiveForCache(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

// --- GetByCode ---

func TestGetByCode_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	r := &repository{db: db}

	mock.ExpectQuery("SELECT").WithArgs("MSG_001").WillReturnRows(
		sqlmock.NewRows(msgCols).AddRow("m1", "MSG_001", "error", "validation", "auth", "Title", "Content", true, time.Now(), time.Now()),
	)

	result, err := r.GetByCode(context.Background(), "MSG_001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Code != "MSG_001" {
		t.Errorf("expected MSG_001, got %s", result.Code)
	}
}

func TestGetByCode_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	r := &repository{db: db}

	mock.ExpectQuery("SELECT").WithArgs("UNKNOWN").WillReturnError(sql.ErrNoRows)

	result, err := r.GetByCode(context.Background(), "UNKNOWN")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if result != nil {
		t.Fatal("expected nil result")
	}
}

// --- GetByCodeForCache ---

func TestGetByCodeForCache_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	r := &repository{db: db}

	mock.ExpectQuery("SELECT").WithArgs("MSG_001").WillReturnRows(
		sqlmock.NewRows(msgCols).AddRow("m1", "MSG_001", "error", "validation", "auth", "Title", "Content", true, time.Now(), time.Now()),
	)

	result, err := r.GetByCodeForCache(context.Background(), "MSG_001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Code != "MSG_001" {
		t.Errorf("expected MSG_001, got %s", result.Code)
	}
}

func TestGetByCodeForCache_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	r := &repository{db: db}

	mock.ExpectQuery("SELECT").WithArgs("UNKNOWN").WillReturnError(sql.ErrNoRows)

	result, err := r.GetByCodeForCache(context.Background(), "UNKNOWN")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if result != nil {
		t.Fatal("expected nil result")
	}
}

// --- GetByCodeWithStatusForCache ---

func TestGetByCodeWithStatusForCache_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	r := &repository{db: db}

	mock.ExpectQuery("SELECT").WithArgs("MSG_001").WillReturnRows(
		sqlmock.NewRows(msgCols).AddRow("m1", "MSG_001", "error", "validation", "auth", "Title", "Content", true, time.Now(), time.Now()),
	)

	result, err := r.GetByCodeWithStatusForCache(context.Background(), "MSG_001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Code != "MSG_001" {
		t.Errorf("expected MSG_001, got %s", result.Code)
	}
}

func TestGetByCodeWithStatusForCache_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	r := &repository{db: db}

	mock.ExpectQuery("SELECT").WithArgs("UNKNOWN").WillReturnError(sql.ErrNoRows)

	result, err := r.GetByCodeWithStatusForCache(context.Background(), "UNKNOWN")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if result != nil {
		t.Fatal("expected nil result")
	}
}
