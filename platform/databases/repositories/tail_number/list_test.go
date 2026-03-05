package tailnumber

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

var lpCols = []string{"id", "tail_number", "aircraft_model_id", "airline_id", "model_name", "airline_name"}

func TestListTailNumbers_All(t *testing.T) {
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

	result, err := r.ListTailNumbers(context.Background(), map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
}

func TestListTailNumbers_ByAirlineID(t *testing.T) {
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

	result, err := r.ListTailNumbers(context.Background(), map[string]interface{}{"airline_id": "a1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
}

func TestListTailNumbers_ByPlate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtByPlate, _ := db.Prepare(QueryGetByTailNumber)

	r := &repository{stmtGetByTailNumber: stmtByPlate}

	rows := sqlmock.NewRows(lpCols).
		AddRow("lp1", "HK-1234", "am1", "a1", "Boeing 737", "Test Air")
	prep.ExpectQuery().WithArgs("HK-1234").WillReturnRows(rows)

	result, err := r.ListTailNumbers(context.Background(), map[string]interface{}{"tail_number": "HK-1234"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
}

func TestListTailNumbers_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtAll, _ := db.Prepare(QueryGetAll)

	r := &repository{stmtGetAll: stmtAll}

	prep.ExpectQuery().WillReturnError(errors.New("db error"))

	_, err = r.ListTailNumbers(context.Background(), map[string]interface{}{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestListTailNumbers_ScanError(t *testing.T) {
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

	_, err = r.ListTailNumbers(context.Background(), map[string]interface{}{})
	if err == nil {
		t.Fatal("expected scan error")
	}
}
