package aircraftmodel

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

var amCols = []string{"id", "model_name", "aircraft_type_name", "engine_type_name", "family", "manufacturer", "status"}

func TestGetAircraftModelByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtByID, _ := db.Prepare(QueryByID)
	r := &repository{stmtGetByID: stmtByID}

	prep.ExpectQuery().WithArgs("am1").WillReturnRows(
		sqlmock.NewRows(amCols).AddRow("am1", "Boeing 737", "Narrow", "Turbofan", "737", "Boeing", true),
	)

	result, err := r.GetAircraftModelByID(context.Background(), "am1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "am1" {
		t.Errorf("expected am1, got %s", result.ID)
	}
}

func TestGetAircraftModelByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtByID, _ := db.Prepare(QueryByID)
	r := &repository{stmtGetByID: stmtByID}

	prep.ExpectQuery().WithArgs("unknown").WillReturnError(sql.ErrNoRows)

	_, err = r.GetAircraftModelByID(context.Background(), "unknown")
	if !errors.Is(err, domain.ErrAircraftModelNotFound) {
		t.Fatalf("expected ErrAircraftModelNotFound, got %v", err)
	}
}

func TestGetAircraftModelsByFamily_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtByFamily, _ := db.Prepare(QueryGetByFamily)
	r := &repository{stmtGetByFamily: stmtByFamily}

	prep.ExpectQuery().WithArgs("737").WillReturnRows(
		sqlmock.NewRows(amCols).AddRow("am1", "Boeing 737-800", "Narrow", "Turbofan", "737", "Boeing", true),
	)

	result, err := r.GetAircraftModelsByFamily(context.Background(), "737")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
}

func TestGetAircraftModelsByFamily_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtByFamily, _ := db.Prepare(QueryGetByFamily)
	r := &repository{stmtGetByFamily: stmtByFamily}

	prep.ExpectQuery().WithArgs("737").WillReturnError(errors.New("db error"))

	_, err = r.GetAircraftModelsByFamily(context.Background(), "737")
	if err == nil {
		t.Fatal("expected error")
	}
}
