package daily_logbook_detail

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestListDailyLogbookDetailsByLogbook_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectPrepare("SELECT")
	stmtByLB, _ := db.Prepare("SELECT 1")

	r := &repository{stmtGetByLogbook: stmtByLB}

	mock.ExpectQuery("SELECT").WithArgs("lb1").WillReturnError(errors.New("db error"))

	_, err = r.ListDailyLogbookDetailsByLogbook(context.Background(), "lb1")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestListDailyLogbookDetailsByLogbook_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectPrepare("SELECT")
	stmtByLB, _ := db.Prepare("SELECT 1")

	r := &repository{stmtGetByLogbook: stmtByLB}

	rows := sqlmock.NewRows([]string{"id"}).AddRow("one")
	mock.ExpectQuery("SELECT").WithArgs("lb1").WillReturnRows(rows)

	_, err = r.ListDailyLogbookDetailsByLogbook(context.Background(), "lb1")
	if err == nil {
		t.Fatal("expected scan error")
	}
}

func TestListDailyLogbookDetailsByLogbook_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectPrepare("SELECT")
	stmtByLB, _ := db.Prepare("SELECT 1")

	r := &repository{stmtGetByLogbook: stmtByLB}

	cols := []string{
		"id", "daily_logbook_id", "flight_real_date", "flight_number",
		"origin_airport_id", "destination_airport_id", "tail_number_id", "passengers",
		"out_time", "takeoff_time", "landing_time", "in_time",
		"pilot_role", "companion_name", "crew_role",
		"air_time", "block_time", "approach_category", "approach_subtype", "autoland", "flight_type",
		"employee_logbook_id", "log_date", "tail_number", "model_name",
		"route_code", "origin_iata_code", "destination_iata_code", "airline_code",
	}
	rows := sqlmock.NewRows(cols)
	mock.ExpectQuery("SELECT").WithArgs("lb1").WillReturnRows(rows)

	result, err := r.ListDailyLogbookDetailsByLogbook(context.Background(), "lb1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil && len(result) != 0 {
		t.Fatalf("expected empty, got %d", len(result))
	}
}
