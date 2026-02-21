package licenseplate

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

var lpCols = []string{"id", "license_plate", "aircraft_model_id", "airline_id", "model_name", "airline_name"}

func TestListLicensePlates_All(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtAll, _ := db.Prepare(QueryGetAll)

	r := &repository{stmtGetAll: stmtAll}

	rows := sqlmock.NewRows(lpCols).
		AddRow("lp1", "HK-1234", "am1", "a1", "Boeing 737", "Test Air")
	prep.ExpectQuery().WillReturnRows(rows)

	result, err := r.ListLicensePlates(context.Background(), map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
}

func TestListLicensePlates_ByAirlineID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtByAirline, _ := db.Prepare(QueryGetByAirline)

	r := &repository{stmtGetByAirline: stmtByAirline}

	rows := sqlmock.NewRows(lpCols).
		AddRow("lp1", "HK-1234", "am1", "a1", "Boeing 737", "Test Air")
	prep.ExpectQuery().WithArgs("a1").WillReturnRows(rows)

	result, err := r.ListLicensePlates(context.Background(), map[string]interface{}{"airline_id": "a1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
}

func TestListLicensePlates_ByPlate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtByPlate, _ := db.Prepare(QueryGetByLicensePlate)

	r := &repository{stmtGetByLicensePlate: stmtByPlate}

	rows := sqlmock.NewRows(lpCols).
		AddRow("lp1", "HK-1234", "am1", "a1", "Boeing 737", "Test Air")
	prep.ExpectQuery().WithArgs("HK-1234").WillReturnRows(rows)

	result, err := r.ListLicensePlates(context.Background(), map[string]interface{}{"license_plate": "HK-1234"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
}

func TestListLicensePlates_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtAll, _ := db.Prepare(QueryGetAll)

	r := &repository{stmtGetAll: stmtAll}

	prep.ExpectQuery().WillReturnError(errors.New("db error"))

	_, err = r.ListLicensePlates(context.Background(), map[string]interface{}{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestListLicensePlates_ScanError(t *testing.T) {
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

	_, err = r.ListLicensePlates(context.Background(), map[string]interface{}{})
	if err == nil {
		t.Fatal("expected scan error")
	}
}
