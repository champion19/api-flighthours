package airline_route

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

func TestGetAirlineRouteByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtByID, _ := db.Prepare(QueryByID)
	r := &repository{stmtGetByID: stmtByID}

	prep.ExpectQuery().WithArgs("ar1").WillReturnRows(
		sqlmock.NewRows(arCols).AddRow("ar1", "r1", "a1", true, "TST", "Test Air", "BOG", "JFK", "BOG-JFK", "El Dorado", "Kennedy", "International", sql.NullString{String: "5h", Valid: true}),
	)

	result, err := r.GetAirlineRouteByID(context.Background(), "ar1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "ar1" {
		t.Errorf("expected ar1, got %s", result.ID)
	}
	if result.EstimatedFlightTime != "5h" {
		t.Errorf("expected 5h, got %s", result.EstimatedFlightTime)
	}
}

func TestGetAirlineRouteByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtByID, _ := db.Prepare(QueryByID)
	r := &repository{stmtGetByID: stmtByID}

	prep.ExpectQuery().WithArgs("unknown").WillReturnError(sql.ErrNoRows)

	_, err = r.GetAirlineRouteByID(context.Background(), "unknown")
	if !errors.Is(err, domain.ErrAirlineRouteNotFound) {
		t.Fatalf("expected ErrAirlineRouteNotFound, got %v", err)
	}
}

func TestGetAirlineRouteByID_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtByID, _ := db.Prepare(QueryByID)
	r := &repository{stmtGetByID: stmtByID}

	prep.ExpectQuery().WithArgs("ar1").WillReturnError(errors.New("db error"))

	_, err = r.GetAirlineRouteByID(context.Background(), "ar1")
	if err == nil {
		t.Fatal("expected error")
	}
}
