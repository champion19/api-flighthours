package airline_route

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

var arCols = []string{
	"id", "route_id", "airline_id", "status", "airline_code", "airline_name",
	"origin_airport_id", "origin_iata_code", "origin_oaci_code",
	"destination_airport_id", "destination_iata_code", "destination_oaci_code", "route_code",
	"origin_airport_name", "destination_airport_name", "airport_type", "estimated_flight_time",
}

func TestListAirlineRoutes_All(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectPrepare("SELECT")
	stmtAll, _ := db.Prepare(QueryGetAll)

	r := &repository{stmtGetAll: stmtAll}

	rows := sqlmock.NewRows(arCols).
		AddRow("ar1", "r1", "a1", true, "TST", "Test Air", "ap1", "BOG", "SKBO", "ap2", "JFK", "KJFK", "BOG-JFK", "El Dorado", "Kennedy", "International", sql.NullString{String: "5h", Valid: true})
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	result, err := r.ListAirlineRoutes(context.Background(), map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
}

func TestListAirlineRoutes_ByAirlineID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectPrepare("SELECT")
	stmtByID, _ := db.Prepare(QueryGetByAirlineID)

	r := &repository{stmtGetByAirlineID: stmtByID}

	rows := sqlmock.NewRows(arCols).
		AddRow("ar1", "r1", "a1", true, "TST", "Test Air", "ap1", "BOG", "SKBO", "ap2", "MDE", "SKRG", "BOG-MDE", "El Dorado", "Cordova", "Domestic", sql.NullString{Valid: false})
	mock.ExpectQuery("SELECT").WithArgs("a1").WillReturnRows(rows)

	result, err := r.ListAirlineRoutes(context.Background(), map[string]interface{}{"airline_id": "a1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
}

func TestListAirlineRoutes_ByAirlineCode(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectPrepare("SELECT")
	stmtByCode, _ := db.Prepare(QueryGetByAirlineCode)

	r := &repository{stmtGetByAirlineCode: stmtByCode}

	rows := sqlmock.NewRows(arCols).
		AddRow("ar1", "r1", "a1", true, "TST", "Test Air", "ap1", "BOG", "SKBO", "ap2", "JFK", "KJFK", "BOG-JFK", "El Dorado", "Kennedy", "International", sql.NullString{Valid: false})
	mock.ExpectQuery("SELECT").WithArgs("TST").WillReturnRows(rows)

	result, err := r.ListAirlineRoutes(context.Background(), map[string]interface{}{"airline_code": "TST"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
}

func TestListAirlineRoutes_ByStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectPrepare("SELECT")
	stmtByStatus, _ := db.Prepare(QueryGetByStatus)

	r := &repository{stmtGetByStatus: stmtByStatus}

	rows := sqlmock.NewRows(arCols)
	mock.ExpectQuery("SELECT").WithArgs(true).WillReturnRows(rows)

	_, err = r.ListAirlineRoutes(context.Background(), map[string]interface{}{"status": true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// An airport registered with only an OACI code (no IATA) has iata_code =
// NULL in the DB. Before this fix, scanning that NULL into a plain Go
// string crashed the whole query — breaking route resolution for every
// employee, not just this one route. Confirms it no longer does.
func TestListAirlineRoutes_DestinationHasNoIataCode_OnlyOaci(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectPrepare("SELECT")
	stmtAll, _ := db.Prepare(QueryGetAll)

	r := &repository{stmtGetAll: stmtAll}

	rows := sqlmock.NewRows(arCols).
		AddRow("ar1", "r1", "a1", true, "TST", "Test Air", "ap1", "BOG", "SKBO", "ap2", nil, "SKRG", "BOG-SKRG", "El Dorado", "Olaya Herrera", "Domestic", sql.NullString{Valid: false})
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	result, err := r.ListAirlineRoutes(context.Background(), map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error scanning a NULL destination_iata_code: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
	if result[0].DestinationIataCode != "" {
		t.Errorf("expected empty DestinationIataCode, got %q", result[0].DestinationIataCode)
	}
	if result[0].DestinationOaciCode != "SKRG" {
		t.Errorf("expected DestinationOaciCode 'SKRG', got %q", result[0].DestinationOaciCode)
	}
	if result[0].RouteCode != "BOG-SKRG" {
		t.Errorf("expected RouteCode 'BOG-SKRG' (OACI fallback), got %q", result[0].RouteCode)
	}
}

func TestListAirlineRoutes_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectPrepare("SELECT")
	stmtAll, _ := db.Prepare(QueryGetAll)

	r := &repository{stmtGetAll: stmtAll}

	mock.ExpectQuery("SELECT").WillReturnError(errors.New("db error"))

	_, err = r.ListAirlineRoutes(context.Background(), map[string]interface{}{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestListAirlineRoutesByAirlineID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectPrepare("SELECT")
	stmtByID, _ := db.Prepare(QueryGetByAirlineID)

	r := &repository{stmtGetByAirlineID: stmtByID}

	rows := sqlmock.NewRows(arCols).
		AddRow("ar1", "r1", "a1", true, "TST", "Test Air", "ap1", "BOG", "SKBO", "ap2", "JFK", "KJFK", "BOG-JFK", "El Dorado", "Kennedy", "International", sql.NullString{Valid: false})
	mock.ExpectQuery("SELECT").WithArgs("a1").WillReturnRows(rows)

	result, err := r.ListAirlineRoutesByAirlineID(context.Background(), "a1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
}
