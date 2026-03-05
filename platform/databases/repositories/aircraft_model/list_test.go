package aircraftmodel

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestListAircraftModels_All(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtAll, _ := db.Prepare(QueryGetAll)

	r := &repository{stmtGetAll: stmtAll}

	cols := []string{"id", "model_name", "aircraft_type_name", "engine_type_name", "family", "manufacturer", "status"}
	rows := sqlmock.NewRows(cols).
		AddRow("am1", "Boeing 737", "Narrow Body", "Turbofan", "737", "Boeing", true)
	prep.ExpectQuery().WillReturnRows(rows)

	result, err := r.ListAircraftModels(context.Background(), map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
}

func TestListAircraftModels_ByEngineType(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtByEngine, _ := db.Prepare(QueryGetByEngineType)

	r := &repository{stmtGetByEngineType: stmtByEngine}

	cols := []string{"id", "model_name", "aircraft_type_name", "engine_type_name", "family", "manufacturer", "status"}
	rows := sqlmock.NewRows(cols).
		AddRow("am1", "ATR 72", "Turboprop", "Turboprop", "ATR", "ATR", true)
	prep.ExpectQuery().WithArgs("Turboprop").WillReturnRows(rows)

	result, err := r.ListAircraftModels(context.Background(), map[string]interface{}{"engine_type": "Turboprop"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
}

func TestListAircraftModels_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtAll, _ := db.Prepare(QueryGetAll)

	r := &repository{stmtGetAll: stmtAll}

	prep.ExpectQuery().WillReturnError(errors.New("db error"))

	_, err = r.ListAircraftModels(context.Background(), map[string]interface{}{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestListAircraftModels_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtAll, _ := db.Prepare(QueryGetAll)

	r := &repository{stmtGetAll: stmtAll}

	rows := sqlmock.NewRows([]string{"id"}).AddRow("only-one")
	prep.ExpectQuery().WillReturnRows(rows)

	_, err = r.ListAircraftModels(context.Background(), map[string]interface{}{})
	if err == nil {
		t.Fatal("expected scan error")
	}
}
