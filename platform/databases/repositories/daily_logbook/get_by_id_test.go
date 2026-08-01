package daily_logbook

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

func TestGetDailyLogbookByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtByID, _ := db.Prepare(QueryByID)
	r := &repository{stmtGetByID: stmtByID}

	cols := []string{"id", "log_date", "employee_id", "book_page", "status", "tail_number_id", "tail_number", "crew_role", "created_at"}
	prep.ExpectQuery().WithArgs("dl1").WillReturnRows(
		sqlmock.NewRows(cols).AddRow("dl1", time.Now(), "emp1", 1, true, nil, nil, nil, time.Now()),
	)

	result, err := r.GetDailyLogbookByID(context.Background(), "dl1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "dl1" {
		t.Errorf("expected dl1, got %s", result.ID)
	}
}

func TestGetDailyLogbookByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtByID, _ := db.Prepare(QueryByID)
	r := &repository{stmtGetByID: stmtByID}

	prep.ExpectQuery().WithArgs("unknown").WillReturnError(sql.ErrNoRows)

	_, err = r.GetDailyLogbookByID(context.Background(), "unknown")
	if !errors.Is(err, domain.ErrDailyLogbookNotFound) {
		t.Fatalf("expected ErrDailyLogbookNotFound, got %v", err)
	}
}

func TestGetDailyLogbookByID_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtByID, _ := db.Prepare(QueryByID)
	r := &repository{stmtGetByID: stmtByID}

	prep.ExpectQuery().WithArgs("dl1").WillReturnError(errors.New("db error"))

	_, err = r.GetDailyLogbookByID(context.Background(), "dl1")
	if err == nil {
		t.Fatal("expected error")
	}
}
