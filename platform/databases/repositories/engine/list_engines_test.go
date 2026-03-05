package engine

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestListEngines_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtAll, _ := db.Prepare(QueryGetAll)

	r := &repository{stmtGetAll: stmtAll, db: db}

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow("e1", "Turbofan").
		AddRow("e2", "Turboprop")
	prep.ExpectQuery().WillReturnRows(rows)

	result, err := r.ListEngines(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result))
	}
}

func TestListEngines_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtAll, _ := db.Prepare(QueryGetAll)

	r := &repository{stmtGetAll: stmtAll, db: db}

	prep.ExpectQuery().WillReturnError(errors.New("db error"))

	_, err = r.ListEngines(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestListEngines_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtAll, _ := db.Prepare(QueryGetAll)

	r := &repository{stmtGetAll: stmtAll, db: db}

	rows := sqlmock.NewRows([]string{"id"}).AddRow("only-one")
	prep.ExpectQuery().WillReturnRows(rows)

	_, err = r.ListEngines(context.Background())
	if err == nil {
		t.Fatal("expected scan error")
	}
}
