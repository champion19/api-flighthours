package licenseplate

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

func TestGetLicensePlateByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtByID, _ := db.Prepare(QueryByID)
	r := &repository{stmtGetByID: stmtByID}

	prep.ExpectQuery().WithArgs("lp1").WillReturnRows(
		sqlmock.NewRows(lpCols).AddRow("lp1", "HK-1234", "am1", "a1", "Boeing 737", "Test Air"),
	)

	result, err := r.GetLicensePlateByID(context.Background(), "lp1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "lp1" {
		t.Errorf("expected lp1, got %s", result.ID)
	}
}

func TestGetLicensePlateByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtByID, _ := db.Prepare(QueryByID)
	r := &repository{stmtGetByID: stmtByID}

	prep.ExpectQuery().WithArgs("unknown").WillReturnError(sql.ErrNoRows)

	_, err = r.GetLicensePlateByID(context.Background(), "unknown")
	if !errors.Is(err, domain.ErrLicensePlateNotFound) {
		t.Fatalf("expected ErrLicensePlateNotFound, got %v", err)
	}
}

func TestGetLicensePlateByID_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtByID, _ := db.Prepare(QueryByID)
	r := &repository{stmtGetByID: stmtByID}

	prep.ExpectQuery().WithArgs("lp1").WillReturnError(errors.New("db error"))

	_, err = r.GetLicensePlateByID(context.Background(), "lp1")
	if err == nil {
		t.Fatal("expected error")
	}
}
