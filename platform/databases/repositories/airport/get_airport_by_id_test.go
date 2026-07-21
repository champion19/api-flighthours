package airport

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

func TestGetAirportByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtByID, _ := db.Prepare(QueryByID)
	r := &repository{stmtGetByID: stmtByID}

	cols := []string{"id", "name", "city", "country", "iata_code", "oaci_code", "status", "airport_type"}
	prep.ExpectQuery().WithArgs("ap1").WillReturnRows(
		sqlmock.NewRows(cols).AddRow("ap1", "El Dorado", "Bogotá", "Colombia", "BOG", "SKBO", true, "International"),
	)

	result, err := r.GetAirportByID(context.Background(), "ap1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.IATACode != "BOG" {
		t.Errorf("expected BOG, got %s", result.IATACode)
	}
}

func TestGetAirportByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtByID, _ := db.Prepare(QueryByID)
	r := &repository{stmtGetByID: stmtByID}

	prep.ExpectQuery().WithArgs("unknown").WillReturnError(sql.ErrNoRows)

	_, err = r.GetAirportByID(context.Background(), "unknown")
	if !errors.Is(err, domain.ErrAirportNotFound) {
		t.Fatalf("expected ErrAirportNotFound, got %v", err)
	}
}

func TestGetAirportByID_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtByID, _ := db.Prepare(QueryByID)
	r := &repository{stmtGetByID: stmtByID}

	prep.ExpectQuery().WithArgs("ap1").WillReturnError(errors.New("db error"))

	_, err = r.GetAirportByID(context.Background(), "ap1")
	if err == nil {
		t.Fatal("expected error")
	}
}
