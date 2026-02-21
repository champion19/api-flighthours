package airline_employee

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

func TestGetAirlineEmployeeByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtByID, _ := db.Prepare(QueryByID)
	r := &repository{stmtGetByID: stmtByID}

	cols := []string{"id", "airline", "bp", "start_date", "end_date", "active"}
	now := time.Now()
	prep.ExpectQuery().WithArgs("e1").WillReturnRows(
		sqlmock.NewRows(cols).AddRow("e1", sql.NullString{String: "a1", Valid: true}, sql.NullString{}, now, now, true),
	)

	result, err := r.GetAirlineEmployeeByID(context.Background(), "e1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "e1" {
		t.Errorf("expected e1, got %s", result.ID)
	}
}

func TestGetAirlineEmployeeByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtByID, _ := db.Prepare(QueryByID)
	r := &repository{stmtGetByID: stmtByID}

	prep.ExpectQuery().WithArgs("unknown").WillReturnError(sql.ErrNoRows)

	_, err = r.GetAirlineEmployeeByID(context.Background(), "unknown")
	if !errors.Is(err, domain.ErrAirlineEmployeeNotFound) {
		t.Fatalf("expected ErrAirlineEmployeeNotFound, got %v", err)
	}
}

func TestGetAirlineEmployeeByID_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtByID, _ := db.Prepare(QueryByID)
	r := &repository{stmtGetByID: stmtByID}

	prep.ExpectQuery().WithArgs("e1").WillReturnError(errors.New("db error"))

	_, err = r.GetAirlineEmployeeByID(context.Background(), "e1")
	if err == nil {
		t.Fatal("expected error")
	}
}
