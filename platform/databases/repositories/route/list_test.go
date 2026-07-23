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
		"id", "origin_airport_id", "origin_iata_code", "origin_oaci_code", "origin_airport_name", "origin_country",
		"destination_airport_id", "destination_iata_code", "destination_oaci_code", "destination_airport_name", "destination_country",
		"airport_type", "estimated_flight_time", "route_code",
	}).AddRow("r1", "ap1", "BOG", "SKBO", "El Dorado", "Colombia", "ap2", "JFK", "KJFK", "Kennedy", "Estados Unidos", "International", sql.NullString{String: "5h", Valid: true}, "BOG-JFK")
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
		"id", "origin_airport_id", "origin_iata_code", "origin_oaci_code", "origin_airport_name", "origin_country",
		"destination_airport_id", "destination_iata_code", "destination_oaci_code", "destination_airport_name", "destination_country",
		"airport_type", "estimated_flight_time", "route_code",
	}).AddRow("r1", "ap1", "BOG", "SKBO", "El Dorado", "Colombia", "ap2", "MDE", "SKRG", "Cordova", "Colombia", "Domestic", sql.NullString{Valid: false}, "BOG-MDE")
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

// An airport registered with only an OACI code (no IATA) has iata_code =
// NULL in the DB. Before this fix, scanning that NULL into a plain Go
// string crashed the whole query — breaking route resolution for every
// employee, not just this one route. Confirms it no longer does.
func TestListRoutes_OriginHasNoIataCode_OnlyOaci(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectPrepare("SELECT")
	stmtAll, _ := db.Prepare(QueryGetAll)

	r := &repository{stmtGetAll: stmtAll}

	rows := sqlmock.NewRows([]string{
		"id", "origin_airport_id", "origin_iata_code", "origin_oaci_code", "origin_airport_name", "origin_country",
		"destination_airport_id", "destination_iata_code", "destination_oaci_code", "destination_airport_name", "destination_country",
		"airport_type", "estimated_flight_time", "route_code",
	}).AddRow("r1", "ap1", nil, "SKRG", "Olaya Herrera", "Colombia", "ap2", "JFK", "KJFK", "Kennedy", "Estados Unidos", "International", nil, "SKRG-JFK")
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	result, err := r.ListRoutes(context.Background(), map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error scanning a NULL origin_iata_code: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
	if result[0].OriginIataCode != "" {
		t.Errorf("expected empty OriginIataCode, got %q", result[0].OriginIataCode)
	}
	if result[0].OriginOaciCode != "SKRG" {
		t.Errorf("expected OriginOaciCode 'SKRG', got %q", result[0].OriginOaciCode)
	}
	if result[0].RouteCode != "SKRG-JFK" {
		t.Errorf("expected RouteCode 'SKRG-JFK' (OACI fallback), got %q", result[0].RouteCode)
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
