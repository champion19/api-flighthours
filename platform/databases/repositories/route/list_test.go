package route

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestListRoutes_All(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectPrepare("SELECT")
	stmtAll, _ := db.Prepare(QueryGetAll)

	r := &repository{stmtGetAll: stmtAll}

	rows := sqlmock.NewRows([]string{
		"id", "origin_airport_id", "origin_iata_code", "origin_airport_name", "origin_country",
		"destination_airport_id", "destination_iata_code", "destination_airport_name", "destination_country",
		"airport_type", "estimated_flight_time", "route_code",
	}).AddRow("r1", "ap1", "BOG", "El Dorado", "Colombia", "ap2", "JFK", "Kennedy", "Estados Unidos", "International", sql.NullString{String: "5h", Valid: true}, "BOG-JFK")
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	result, err := r.ListRoutes(context.Background(), map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
	if result[0].RouteCode != "BOG-JFK" {
		t.Errorf("expected BOG-JFK, got %s", result[0].RouteCode)
	}
}

func TestListRoutes_ByAirportType(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectPrepare("SELECT")
	stmtByType, _ := db.Prepare(QueryGetByAirportType)

	r := &repository{stmtGetByAirportType: stmtByType}

	rows := sqlmock.NewRows([]string{
		"id", "origin_airport_id", "origin_iata_code", "origin_airport_name", "origin_country",
		"destination_airport_id", "destination_iata_code", "destination_airport_name", "destination_country",
		"airport_type", "estimated_flight_time", "route_code",
	}).AddRow("r1", "ap1", "BOG", "El Dorado", "Colombia", "ap2", "MDE", "Cordova", "Colombia", "Domestic", sql.NullString{Valid: false}, "BOG-MDE")
	mock.ExpectQuery("SELECT").WithArgs("Domestic").WillReturnRows(rows)

	result, err := r.ListRoutes(context.Background(), map[string]interface{}{"airport_type": "Domestic"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
}

func TestListRoutes_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectPrepare("SELECT")
	stmtAll, _ := db.Prepare(QueryGetAll)

	r := &repository{stmtGetAll: stmtAll}

	mock.ExpectQuery("SELECT").WillReturnError(errors.New("db error"))

	_, err = r.ListRoutes(context.Background(), map[string]interface{}{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestListRoutes_ScanError(t *testing.T) {
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

	_, err = r.ListRoutes(context.Background(), map[string]interface{}{})
	if err == nil {
		t.Fatal("expected scan error")
	}
}
