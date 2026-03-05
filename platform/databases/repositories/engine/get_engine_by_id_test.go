package engine

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

func TestGetEngineByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtByID, _ := db.Prepare(QueryByID)

	r := &repository{stmtGetByID: stmtByID}

	prep.ExpectQuery().WithArgs("e1").WillReturnRows(
		sqlmock.NewRows([]string{"id", "name"}).AddRow("e1", "Turbofan"),
	)

	result, err := r.GetEngineByID(context.Background(), "e1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Name != "Turbofan" {
		t.Errorf("expected Turbofan, got %s", result.Name)
	}
}

func TestGetEngineByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtByID, _ := db.Prepare(QueryByID)

	r := &repository{stmtGetByID: stmtByID}

	prep.ExpectQuery().WithArgs("unknown").WillReturnError(sql.ErrNoRows)

	_, err = r.GetEngineByID(context.Background(), "unknown")
	if !errors.Is(err, domain.ErrEngineNotFound) {
		t.Fatalf("expected ErrEngineNotFound, got %v", err)
	}
}

func TestGetEngineByID_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtByID, _ := db.Prepare(QueryByID)

	r := &repository{stmtGetByID: stmtByID}

	prep.ExpectQuery().WithArgs("e1").WillReturnError(errors.New("db error"))

	_, err = r.GetEngineByID(context.Background(), "e1")
	if err == nil {
		t.Fatal("expected error")
	}
}
