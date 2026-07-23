package route

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

func TestGetRouteByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtByID, _ := db.Prepare(QueryByID)
	r := &repository{stmtGetByID: stmtByID}

	cols := []string{"id", "origin_airport_id", "origin_iata_code", "origin_oaci_code", "origin_airport_name", "origin_country",
		"destination_airport_id", "destination_iata_code", "destination_oaci_code", "destination_airport_name", "destination_country",
		"airport_type", "estimated_flight_time", "route_code"}
	prep.ExpectQuery().WithArgs("r1").WillReturnRows(
		sqlmock.NewRows(cols).AddRow("r1", "ap1", "BOG", "SKBO", "El Dorado", "Colombia", "ap2", "JFK", "KJFK", "Kennedy", "Estados Unidos", "International", sql.NullString{String: "5h", Valid: true}, "BOG-JFK"),
	)

	result, err := r.GetRouteByID(context.Background(), "r1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.RouteCode != "BOG-JFK" {
		t.Errorf("expected BOG-JFK, got %s", result.RouteCode)
	}
}

func TestGetRouteByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtByID, _ := db.Prepare(QueryByID)
	r := &repository{stmtGetByID: stmtByID}

	prep.ExpectQuery().WithArgs("unknown").WillReturnError(sql.ErrNoRows)

	_, err = r.GetRouteByID(context.Background(), "unknown")
	if !errors.Is(err, domain.ErrRouteNotFound) {
		t.Fatalf("expected ErrRouteNotFound, got %v", err)
	}
}

func TestGetRouteByID_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtByID, _ := db.Prepare(QueryByID)
	r := &repository{stmtGetByID: stmtByID}

	prep.ExpectQuery().WithArgs("r1").WillReturnError(errors.New("db error"))

	_, err = r.GetRouteByID(context.Background(), "r1")
	if err == nil {
		t.Fatal("expected error")
	}
}
