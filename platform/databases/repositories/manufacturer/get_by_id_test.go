package manufacturer

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

func TestGetManufacturerByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtByID, _ := db.Prepare(QueryByID)
	r := &repository{stmtGetByID: stmtByID}

	prep.ExpectQuery().WithArgs("m1").WillReturnRows(
		sqlmock.NewRows([]string{"id", "name"}).AddRow("m1", "Boeing"),
	)

	result, err := r.GetManufacturerByID(context.Background(), "m1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Name != "Boeing" {
		t.Errorf("expected Boeing, got %s", result.Name)
	}
}

func TestGetManufacturerByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtByID, _ := db.Prepare(QueryByID)
	r := &repository{stmtGetByID: stmtByID}

	prep.ExpectQuery().WithArgs("unknown").WillReturnError(sql.ErrNoRows)

	_, err = r.GetManufacturerByID(context.Background(), "unknown")
	if !errors.Is(err, domain.ErrManufacturerNotFound) {
		t.Fatalf("expected ErrManufacturerNotFound, got %v", err)
	}
}

func TestGetManufacturerByID_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtByID, _ := db.Prepare(QueryByID)
	r := &repository{stmtGetByID: stmtByID}

	prep.ExpectQuery().WithArgs("m1").WillReturnError(errors.New("db error"))

	_, err = r.GetManufacturerByID(context.Background(), "m1")
	if err == nil {
		t.Fatal("expected error")
	}
}
