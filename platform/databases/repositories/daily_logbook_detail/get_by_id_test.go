package daily_logbook_detail

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetDailyLogbookDetailByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtByID, _ := db.Prepare("SELECT 1")
	r := &repository{stmtGetByID: stmtByID}

	prep.ExpectQuery().WithArgs("unknown").WillReturnError(sql.ErrNoRows)

	result, err := r.GetDailyLogbookDetailByID(context.Background(), "unknown")
	if err != nil {
		t.Fatalf("expected nil error for not found, got %v", err)
	}
	if result != nil {
		t.Fatal("expected nil result")
	}
}

func TestGetDailyLogbookDetailByID_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtByID, _ := db.Prepare("SELECT 1")
	r := &repository{stmtGetByID: stmtByID}

	prep.ExpectQuery().WithArgs("id1").WillReturnError(errors.New("db error"))

	_, err = r.GetDailyLogbookDetailByID(context.Background(), "id1")
	if err == nil {
		t.Fatal("expected error")
	}
}
