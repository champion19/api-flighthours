package airline

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

func TestGetAirlineByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtByID, _ := db.Prepare(QueryByID)
	r := &repository{stmtGetByID: stmtByID}

	cols := []string{"id", "airline_name", "airline_code", "status"}
	prep.ExpectQuery().WithArgs("a1").WillReturnRows(
		sqlmock.NewRows(cols).AddRow("a1", "Test Air", "TST", "active"),
	)

	result, err := r.GetAirlineByID(context.Background(), "a1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "a1" {
		t.Errorf("expected a1, got %s", result.ID)
	}
}

func TestGetAirlineByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtByID, _ := db.Prepare(QueryByID)
	r := &repository{stmtGetByID: stmtByID}

	prep.ExpectQuery().WithArgs("unknown").WillReturnError(sql.ErrNoRows)

	_, err = r.GetAirlineByID(context.Background(), "unknown")
	if !errors.Is(err, domain.ErrAirlineNotFound) {
		t.Fatalf("expected ErrAirlineNotFound, got %v", err)
	}
}

func TestGetAirlineByID_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtByID, _ := db.Prepare(QueryByID)
	r := &repository{stmtGetByID: stmtByID}

	prep.ExpectQuery().WithArgs("a1").WillReturnError(errors.New("db error"))

	_, err = r.GetAirlineByID(context.Background(), "a1")
	if err == nil {
		t.Fatal("expected error")
	}
}
