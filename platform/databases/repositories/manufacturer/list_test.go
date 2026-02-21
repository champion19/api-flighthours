package manufacturer

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestListManufacturers_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtAll, _ := db.Prepare(QueryGetAll)

	r := &repository{stmtGetAll: stmtAll, db: db}

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow("m1", "Boeing").
		AddRow("m2", "Airbus")
	prep.ExpectQuery().WillReturnRows(rows)

	result, err := r.ListManufacturers(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result))
	}
	if result[0].Name != "Boeing" {
		t.Errorf("expected Boeing, got %s", result[0].Name)
	}
}

func TestListManufacturers_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtAll, _ := db.Prepare(QueryGetAll)

	r := &repository{stmtGetAll: stmtAll, db: db}

	prep.ExpectQuery().WillReturnError(errors.New("db error"))

	_, err = r.ListManufacturers(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestListManufacturers_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtAll, _ := db.Prepare(QueryGetAll)

	r := &repository{stmtGetAll: stmtAll, db: db}

	rows := sqlmock.NewRows([]string{"id"}).AddRow("only-one-col")
	prep.ExpectQuery().WillReturnRows(rows)

	_, err = r.ListManufacturers(context.Background())
	if err == nil {
		t.Fatal("expected scan error")
	}
}

func TestListManufacturers_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtAll, _ := db.Prepare(QueryGetAll)

	r := &repository{stmtGetAll: stmtAll, db: db}

	rows := sqlmock.NewRows([]string{"id", "name"})
	prep.ExpectQuery().WillReturnRows(rows)

	result, err := r.ListManufacturers(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil && len(result) != 0 {
		t.Fatalf("expected empty result, got %d", len(result))
	}
}
