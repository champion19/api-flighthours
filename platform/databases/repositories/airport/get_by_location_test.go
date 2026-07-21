package airport

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetAirportsByType_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtByType, _ := db.Prepare(QueryGetByType)
	r := &repository{stmtGetByType: stmtByType}

	cols := []string{"id", "name", "city", "country", "iata_code", "oaci_code", "status", "airport_type"}
	prep.ExpectQuery().WithArgs("International").WillReturnRows(
		sqlmock.NewRows(cols).AddRow("ap1", "El Dorado", "Bogotá", "Colombia", "BOG", "SKBO", true, "International"),
	)

	result, err := r.GetAirportsByType(context.Background(), "International")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
}

func TestGetAirportsByType_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtByType, _ := db.Prepare(QueryGetByType)
	r := &repository{stmtGetByType: stmtByType}

	cols := []string{"id", "name", "city", "country", "iata_code", "oaci_code", "status", "airport_type"}
	prep.ExpectQuery().WithArgs("Unknown").WillReturnRows(sqlmock.NewRows(cols))

	_, err = r.GetAirportsByType(context.Background(), "Unknown")
	if err == nil {
		t.Fatal("expected ErrNoRows for empty result")
	}
}

func TestGetAirportsByType_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtByType, _ := db.Prepare(QueryGetByType)
	r := &repository{stmtGetByType: stmtByType}

	prep.ExpectQuery().WithArgs("International").WillReturnError(errors.New("db error"))

	_, err = r.GetAirportsByType(context.Background(), "International")
	if err == nil {
		t.Fatal("expected error")
	}
}
