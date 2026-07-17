package flight_summary

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// ═══════════════════════════════════════════
// Constants to avoid SonarQube duplication warnings
// ═══════════════════════════════════════════

const (
	testEmployeeID    = "emp-1"
	testStartDate     = "2026-01-01"
	testEndDate       = "2026-01-31"
	testDate          = "2026-01-15"
	testStartQ1       = "2026-01-01"
	testEndQ1         = "2026-03-31"
	stmtSelect        = "SELECT"
	stmtSelectCount   = "SELECT COUNT"
	stmtSelectCoal    = "SELECT COALESCE"
	errMsgUnexpected  = "unexpected error: %v"
	errMsgExpected    = "expected error"
	errMsgExpectedNil = "expected nil repo"
)

// ═══════════════════════════════════════════
// Helper: create a valid repository with sqlmock
// ═══════════════════════════════════════════

func newMockRepo(t *testing.T) (*repository, sqlmock.Sqlmock, *sql.DB) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	mock.ExpectPrepare(stmtSelect)
	mock.ExpectPrepare(stmtSelectCount)
	mock.ExpectPrepare(stmtSelectCoal)
	mock.ExpectPrepare(stmtSelect)

	repo, err := NewFlightSummaryRepository(db)
	if err != nil {
		t.Fatalf(errMsgUnexpected, err)
	}
	return repo, mock, db
}

// ═══════════════════════════════════════════
// NewFlightSummaryRepository
// ═══════════════════════════════════════════

func TestNewFlightSummaryRepository(t *testing.T) {
	t.Run("nil db", func(t *testing.T) {
		repo, err := NewFlightSummaryRepository(nil)
		if err != sql.ErrConnDone {
			t.Fatalf("expected sql.ErrConnDone, got %v", err)
		}
		if repo != nil {
			t.Fatal(errMsgExpectedNil)
		}
	})

	t.Run("success", func(t *testing.T) {
		repo, _, db := newMockRepo(t)
		defer db.Close()
		if repo == nil {
			t.Fatal("expected non-nil repo")
		}
	})

	t.Run("prepare error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}
		defer db.Close()

		mock.ExpectPrepare(stmtSelect).WillReturnError(sql.ErrConnDone)

		repo, err := NewFlightSummaryRepository(db)
		if err == nil {
			t.Fatal(errMsgExpected)
		}
		if repo != nil {
			t.Fatal(errMsgExpectedNil)
		}
	})
}

// ═══════════════════════════════════════════
// GetFlightHoursSummary
// ═══════════════════════════════════════════

func TestGetFlightHoursSummary(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo, mock, db := newMockRepo(t)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"pilot_role", "total_seconds", "flight_count"}).
			AddRow("PF", 7200, 3).
			AddRow("PM", 3600, 2)
		mock.ExpectQuery(stmtSelect).WithArgs(testEmployeeID, testStartDate, testEndDate).WillReturnRows(rows)

		result, err := repo.GetFlightHoursSummary(context.Background(), testEmployeeID, testStartDate, testEndDate)
		if err != nil {
			t.Fatalf(errMsgUnexpected, err)
		}
		if len(result) != 2 {
			t.Errorf("expected 2 breakdowns, got %d", len(result))
		}
		if result[0].PilotRole != "PF" {
			t.Errorf("expected PF, got %s", result[0].PilotRole)
		}
	})

	t.Run("query error", func(t *testing.T) {
		repo, mock, db := newMockRepo(t)
		defer db.Close()

		mock.ExpectQuery(stmtSelect).WithArgs(testEmployeeID, testStartDate, testEndDate).WillReturnError(sql.ErrConnDone)

		_, err := repo.GetFlightHoursSummary(context.Background(), testEmployeeID, testStartDate, testEndDate)
		if err == nil {
			t.Fatal(errMsgExpected)
		}
	})

	t.Run("empty result", func(t *testing.T) {
		repo, mock, db := newMockRepo(t)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"pilot_role", "total_seconds", "flight_count"})
		mock.ExpectQuery(stmtSelect).WithArgs(testEmployeeID, testStartDate, testEndDate).WillReturnRows(rows)

		result, err := repo.GetFlightHoursSummary(context.Background(), testEmployeeID, testStartDate, testEndDate)
		if err != nil {
			t.Fatalf(errMsgUnexpected, err)
		}
		if len(result) != 0 {
			t.Errorf("expected 0, got %d", len(result))
		}
	})
}

// ═══════════════════════════════════════════
// GetLandingCount
// ═══════════════════════════════════════════

func TestGetLandingCount(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo, mock, db := newMockRepo(t)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"count"}).AddRow(5)
		mock.ExpectQuery(stmtSelectCount).WithArgs(testEmployeeID, testStartQ1, testEndQ1).WillReturnRows(rows)

		count, err := repo.GetLandingCount(context.Background(), testEmployeeID, testStartQ1, testEndQ1)
		if err != nil {
			t.Fatalf(errMsgUnexpected, err)
		}
		if count != 5 {
			t.Errorf("expected 5, got %d", count)
		}
	})

	t.Run("error", func(t *testing.T) {
		repo, mock, db := newMockRepo(t)
		defer db.Close()

		mock.ExpectQuery(stmtSelectCount).WillReturnError(sql.ErrConnDone)

		_, err := repo.GetLandingCount(context.Background(), testEmployeeID, testStartQ1, testEndQ1)
		if err == nil {
			t.Fatal(errMsgExpected)
		}
	})
}

// ═══════════════════════════════════════════
// GetDailyFlightSeconds
// ═══════════════════════════════════════════

func TestGetDailyFlightSeconds(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo, mock, db := newMockRepo(t)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"total_seconds"}).AddRow(3600)
		mock.ExpectQuery(stmtSelectCoal).WithArgs(testEmployeeID, testDate).WillReturnRows(rows)

		secs, err := repo.GetDailyFlightSeconds(context.Background(), testEmployeeID, testDate)
		if err != nil {
			t.Fatalf(errMsgUnexpected, err)
		}
		if secs != 3600 {
			t.Errorf("expected 3600, got %d", secs)
		}
	})

	t.Run("error", func(t *testing.T) {
		repo, mock, db := newMockRepo(t)
		defer db.Close()

		mock.ExpectQuery(stmtSelectCoal).WillReturnError(sql.ErrConnDone)

		_, err := repo.GetDailyFlightSeconds(context.Background(), testEmployeeID, testDate)
		if err == nil {
			t.Fatal(errMsgExpected)
		}
	})
}

// ═══════════════════════════════════════════
// GetRecentFlights
// ═══════════════════════════════════════════

func TestGetRecentFlights(t *testing.T) {
	detailCols := []string{
		"id", "daily_logbook_id", "flight_real_date", "flight_number",
		"airline_route_id", "actual_tail_number_id", "passengers",
		"out_time", "takeoff_time", "landing_time", "in_time",
		"pilot_role", "companion_name", "crew_role",
		"air_time", "block_time", "approach_category", "approach_subtype", "autoland", "flight_type",
		"employee_logbook_id", "log_date", "tail_number", "model_name",
		"route_code", "origin_iata_code", "destination_iata_code", "airline_code",
	}

	t.Run("success", func(t *testing.T) {
		repo, mock, db := newMockRepo(t)
		defer db.Close()

		rows := sqlmock.NewRows(detailCols).
			AddRow("d-1", "lb-1", testDate, "AV123",
				"route-1", "tn-1", 120,
				"08:00", "08:15", "09:30", "09:45",
				"PF", nil, nil,
				"01:15", "01:45", nil, nil, nil, nil,
				nil, testDate, "HK-1234", "A320",
				"BOG-CLO", "BOG", "CLO", "AV")
		mock.ExpectQuery(stmtSelect).WithArgs(testEmployeeID, 5).WillReturnRows(rows)

		result, err := repo.GetRecentFlights(context.Background(), testEmployeeID, 5)
		if err != nil {
			t.Fatalf(errMsgUnexpected, err)
		}
		if len(result) != 1 {
			t.Fatalf("expected 1 flight, got %d", len(result))
		}
		if result[0].FlightNumber != "AV123" {
			t.Errorf("expected AV123, got %s", result[0].FlightNumber)
		}
		if result[0].TailNumber != "HK-1234" {
			t.Errorf("expected HK-1234, got %s", result[0].TailNumber)
		}
	})

	t.Run("query error", func(t *testing.T) {
		repo, mock, db := newMockRepo(t)
		defer db.Close()

		mock.ExpectQuery(stmtSelect).WillReturnError(sql.ErrConnDone)

		_, err := repo.GetRecentFlights(context.Background(), testEmployeeID, 5)
		if err == nil {
			t.Fatal(errMsgExpected)
		}
	})

	t.Run("empty result", func(t *testing.T) {
		repo, mock, db := newMockRepo(t)
		defer db.Close()

		rows := sqlmock.NewRows(detailCols)
		mock.ExpectQuery(stmtSelect).WithArgs(testEmployeeID, 5).WillReturnRows(rows)

		result, err := repo.GetRecentFlights(context.Background(), testEmployeeID, 5)
		if err != nil {
			t.Fatalf(errMsgUnexpected, err)
		}
		if len(result) != 0 {
			t.Errorf("expected 0, got %d", len(result))
		}
	})
}

// ═══════════════════════════════════════════
// BeginTx
// ═══════════════════════════════════════════

func TestBeginTx(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo, mock, db := newMockRepo(t)
		defer db.Close()

		mock.ExpectBegin()
		tx, err := repo.BeginTx(context.Background())
		if err != nil {
			t.Fatalf(errMsgUnexpected, err)
		}
		if tx == nil {
			t.Fatal("expected non-nil tx")
		}
	})

	t.Run("error", func(t *testing.T) {
		repo, mock, db := newMockRepo(t)
		defer db.Close()

		mock.ExpectBegin().WillReturnError(sql.ErrConnDone)
		tx, err := repo.BeginTx(context.Background())
		if err == nil {
			t.Fatal(errMsgExpected)
		}
		if tx != nil {
			t.Fatal("expected nil tx")
		}
	})
}
