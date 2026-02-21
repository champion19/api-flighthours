package airline

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestListAirlines_All(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectPrepare("SELECT")
	stmtAll, _ := db.Prepare(QueryGetAll)

	r := &repository{stmtGetAll: stmtAll}

	rows := sqlmock.NewRows([]string{"id", "airline_name", "airline_code", "status"}).
		AddRow("a1", "Test Air", "TST", "active")
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	result, err := r.ListAirlines(context.Background(), map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}
}

func TestListAirlines_ByStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectPrepare("SELECT")
	stmtByStatus, _ := db.Prepare(QueryGetByStatus)

	r := &repository{stmtGetByStatus: stmtByStatus}

	rows := sqlmock.NewRows([]string{"id", "airline_name", "airline_code", "status"}).
		AddRow("a1", "Active Air", "ACT", "active")
	mock.ExpectQuery("SELECT").WithArgs(true).WillReturnRows(rows)

	result, err := r.ListAirlines(context.Background(), map[string]interface{}{"status": true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}
}

func TestListAirlines_StatusNonBool(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectPrepare("SELECT")
	stmtAll, _ := db.Prepare(QueryGetAll)

	r := &repository{stmtGetAll: stmtAll}

	rows := sqlmock.NewRows([]string{"id", "airline_name", "airline_code", "status"})
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	_, err = r.ListAirlines(context.Background(), map[string]interface{}{"status": "invalid"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListAirlines_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectPrepare("SELECT")
	stmtAll, _ := db.Prepare(QueryGetAll)

	r := &repository{stmtGetAll: stmtAll}

	mock.ExpectQuery("SELECT").WillReturnError(errors.New("db error"))

	_, err = r.ListAirlines(context.Background(), map[string]interface{}{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestListAirlines_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectPrepare("SELECT")
	stmtAll, _ := db.Prepare(QueryGetAll)

	r := &repository{stmtGetAll: stmtAll}

	rows := sqlmock.NewRows([]string{"id"}).AddRow("only-one")
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	_, err = r.ListAirlines(context.Background(), map[string]interface{}{})
	if err == nil {
		t.Fatal("expected scan error")
	}
}
