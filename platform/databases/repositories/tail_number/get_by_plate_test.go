package tailnumber

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

func TestGetTailNumberByPlate_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtByPlate, _ := db.Prepare(QueryGetByTailNumber)

	r := &repository{stmtGetByTailNumber: stmtByPlate}

	prep.ExpectQuery().WithArgs("HK-1234").WillReturnRows(
		sqlmock.NewRows([]string{"id", "tail_number", "aircraft_model_id", "airline_id", "model_name", "airline_name"}).
			AddRow("lp1", "HK-1234", "am1", "a1", "Boeing 737", "Test Air"),
	)

	result, err := r.GetTailNumberByPlate(context.Background(), "HK-1234")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.TailNumber != "HK-1234" {
		t.Errorf("expected HK-1234, got %s", result.TailNumber)
	}
}

func TestGetTailNumberByPlate_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtByPlate, _ := db.Prepare(QueryGetByTailNumber)

	r := &repository{stmtGetByTailNumber: stmtByPlate}

	prep.ExpectQuery().WithArgs("UNKNOWN").WillReturnError(sql.ErrNoRows)

	_, err = r.GetTailNumberByPlate(context.Background(), "UNKNOWN")
	if !errors.Is(err, domain.ErrTailNumberNotFound) {
		t.Fatalf("expected ErrTailNumberNotFound, got %v", err)
	}
}

func TestGetTailNumberByPlate_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtByPlate, _ := db.Prepare(QueryGetByTailNumber)

	r := &repository{stmtGetByTailNumber: stmtByPlate}

	prep.ExpectQuery().WithArgs("HK-1234").WillReturnError(errors.New("db error"))

	_, err = r.GetTailNumberByPlate(context.Background(), "HK-1234")
	if err == nil {
		t.Fatal("expected error")
	}
}
