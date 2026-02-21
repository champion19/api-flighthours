package daily_logbook_detail

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestListDailyLogbookDetailsByEmployee_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectPrepare("SELECT")
	stmtByEmp, _ := db.Prepare("SELECT 1")

	r := &repository{stmtGetByEmployee: stmtByEmp}

	mock.ExpectQuery("SELECT").WithArgs("emp1").WillReturnError(errors.New("db error"))

	_, err = r.ListDailyLogbookDetailsByEmployee(context.Background(), "emp1")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestListDailyLogbookDetailsByEmployee_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectPrepare("SELECT")
	stmtByEmp, _ := db.Prepare("SELECT 1")

	r := &repository{stmtGetByEmployee: stmtByEmp}

	rows := sqlmock.NewRows([]string{"id"}).AddRow("one")
	mock.ExpectQuery("SELECT").WithArgs("emp1").WillReturnRows(rows)

	_, err = r.ListDailyLogbookDetailsByEmployee(context.Background(), "emp1")
	if err == nil {
		t.Fatal("expected scan error")
	}
}

func TestListDailyLogbookDetailsByEmployee_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectPrepare("SELECT")
	stmtByEmp, _ := db.Prepare("SELECT 1")

	r := &repository{stmtGetByEmployee: stmtByEmp}

	rows := sqlmock.NewRows(detailCols)
	mock.ExpectQuery("SELECT").WithArgs("emp1").WillReturnRows(rows)

	result, err := r.ListDailyLogbookDetailsByEmployee(context.Background(), "emp1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil && len(result) != 0 {
		t.Fatalf("expected empty, got %d", len(result))
	}
}

// Shared column definitions for daily_logbook_detail queries
var detailCols = []string{
	"id", "daily_logbook_id", "flight_real_date", "flight_number",
	"airline_route_id", "license_plate_id", "passengers",
	"out_time", "takeoff_time", "landing_time", "in_time",
	"pilot_role", "companion_name", "crew_role",
	"air_time", "block_time", "duty_time", "approach_type", "flight_type",
	"employee_logbook_id", "log_date", "license_plate", "model_name",
	"route_code", "origin_iata_code", "destination_iata_code", "airline_code",
}
