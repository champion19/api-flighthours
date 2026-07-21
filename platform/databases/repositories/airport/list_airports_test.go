package airport

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestListAirports_All(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectPrepare("SELECT")
	stmtAll, _ := db.Prepare(QueryGetAll)

	r := &repository{stmtGetAll: stmtAll}

	rows := sqlmock.NewRows([]string{"id", "name", "city", "country", "iata_code", "oaci_code", "status", "airport_type"}).
		AddRow("ap1", "El Dorado", "Bogotá", "Colombia", "BOG", "SKBO", true, "International")
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	result, err := r.ListAirports(context.Background(), map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
}

func TestListAirports_ByStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectPrepare("SELECT")
	stmtByStatus, _ := db.Prepare(QueryGetByStatus)

	r := &repository{stmtGetByStatus: stmtByStatus}

	rows := sqlmock.NewRows([]string{"id", "name", "city", "country", "iata_code", "oaci_code", "status", "airport_type"}).
		AddRow("ap1", "El Dorado", "Bogotá", "Colombia", "BOG", "SKBO", true, "International")
	mock.ExpectQuery("SELECT").WithArgs(true).WillReturnRows(rows)

	result, err := r.ListAirports(context.Background(), map[string]interface{}{"status": true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
}

func TestListAirports_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectPrepare("SELECT")
	stmtAll, _ := db.Prepare(QueryGetAll)

	r := &repository{stmtGetAll: stmtAll}

	mock.ExpectQuery("SELECT").WillReturnError(errors.New("db error"))

	_, err = r.ListAirports(context.Background(), map[string]interface{}{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestListAirports_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectPrepare("SELECT")
	stmtAll, _ := db.Prepare(QueryGetAll)

	r := &repository{stmtGetAll: stmtAll}

	rows := sqlmock.NewRows([]string{"id"}).AddRow("one")
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	_, err = r.ListAirports(context.Background(), map[string]interface{}{})
	if err == nil {
		t.Fatal("expected scan error")
	}
}
