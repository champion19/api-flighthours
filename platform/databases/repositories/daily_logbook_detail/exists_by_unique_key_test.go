package daily_logbook_detail

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestExistsByUniqueKey_Exists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtExists, _ := db.Prepare("SELECT 1")
	r := &repository{stmtExistsByUniqueKey: stmtExists}

	prep.ExpectQuery().WithArgs("lb1", "2026-01-15", "FL100", "lp1").WillReturnRows(
		sqlmock.NewRows([]string{"count"}).AddRow(1),
	)

	exists, err := r.ExistsByUniqueKey(context.Background(), "lb1", "2026-01-15", "FL100", "lp1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !exists {
		t.Fatal("expected exists=true")
	}
}

func TestExistsByUniqueKey_NotExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtExists, _ := db.Prepare("SELECT 1")
	r := &repository{stmtExistsByUniqueKey: stmtExists}

	prep.ExpectQuery().WithArgs("lb1", "2026-01-15", "FL100", "lp1").WillReturnRows(
		sqlmock.NewRows([]string{"count"}).AddRow(0),
	)

	exists, err := r.ExistsByUniqueKey(context.Background(), "lb1", "2026-01-15", "FL100", "lp1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exists {
		t.Fatal("expected exists=false")
	}
}

func TestExistsByUniqueKey_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtExists, _ := db.Prepare("SELECT 1")
	r := &repository{stmtExistsByUniqueKey: stmtExists}

	prep.ExpectQuery().WithArgs("lb1", "2026-01-15", "FL100", "lp1").WillReturnError(errors.New("db error"))

	_, err = r.ExistsByUniqueKey(context.Background(), "lb1", "2026-01-15", "FL100", "lp1")
	if err == nil {
		t.Fatal("expected error")
	}
}
