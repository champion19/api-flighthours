package employee_flight_summary

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
	testEmployeeID   = "emp-1"
	stmtInsert       = "INSERT INTO"
	stmtSelect       = "SELECT id, employee_id"
	errMsgUnexpected = "unexpected error: %v"
	errMsgExpected   = "expected error"
	testPeriodStart  = "2026-01-01"
	testPeriodEnd    = "2026-01-31"
	testLastUpdated  = "2026-01-31 23:59:59"
)

var summaryCols = []string{
	"id", "employee_id", "period_type", "period_year", "period_number",
	"period_start", "period_end", "total_air_time", "total_block_time",
	"total_flights", "total_landings", "last_updated",
}

// ═══════════════════════════════════════════
// Helper
// ═══════════════════════════════════════════

func newMockRepo(t *testing.T) (*repository, sqlmock.Sqlmock, *sql.DB) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	mock.ExpectPrepare(stmtInsert)
	mock.ExpectPrepare(stmtSelect)
	mock.ExpectPrepare(stmtSelect)
	mock.ExpectPrepare(stmtSelect)

	repo, err := NewEmployeeFlightSummaryRepository(db)
	if err != nil {
		t.Fatalf(errMsgUnexpected, err)
	}
	return repo, mock, db
}

// ═══════════════════════════════════════════
// NewEmployeeFlightSummaryRepository
// ═══════════════════════════════════════════

func TestNewRepository(t *testing.T) {
	t.Run("nil db", func(t *testing.T) {
		repo, err := NewEmployeeFlightSummaryRepository(nil)
		if err != sql.ErrConnDone {
			t.Fatalf("expected sql.ErrConnDone, got %v", err)
		}
		if repo != nil {
			t.Fatal("expected nil repo")
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

		mock.ExpectPrepare(stmtInsert).WillReturnError(sql.ErrConnDone)

		repo, err := NewEmployeeFlightSummaryRepository(db)
		if err == nil {
			t.Fatal(errMsgExpected)
		}
		if repo != nil {
			t.Fatal("expected nil repo")
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

// ═══════════════════════════════════════════
// GetSummariesByEmployee
// ═══════════════════════════════════════════

func TestGetSummariesByEmployee(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo, mock, db := newMockRepo(t)
		defer db.Close()

		rows := sqlmock.NewRows(summaryCols).
			AddRow("s-1", testEmployeeID, "MONTHLY", 2026, 1, testPeriodStart, testPeriodEnd, 9000, 10800, 5, 3, testLastUpdated).
			AddRow("s-2", testEmployeeID, "MONTHLY", 2025, 12, "2025-12-01", "2025-12-31", 7200, 8400, 4, 2, nil)
		mock.ExpectQuery(stmtSelect).WithArgs(testEmployeeID, "MONTHLY").WillReturnRows(rows)

		result, err := repo.GetSummariesByEmployee(context.Background(), testEmployeeID, "MONTHLY")
		if err != nil {
			t.Fatalf(errMsgUnexpected, err)
		}
		if len(result) != 2 {
			t.Fatalf("expected 2, got %d", len(result))
		}
		if result[0].TotalAirTime != 9000 {
			t.Errorf("expected 9000, got %d", result[0].TotalAirTime)
		}
		if result[0].LastUpdated != testLastUpdated {
			t.Errorf("expected last_updated set, got %q", result[0].LastUpdated)
		}
		if result[1].LastUpdated != "" {
			t.Errorf("expected empty last_updated for nil, got %q", result[1].LastUpdated)
		}
	})

	t.Run("error", func(t *testing.T) {
		repo, mock, db := newMockRepo(t)
		defer db.Close()

		mock.ExpectQuery(stmtSelect).WillReturnError(sql.ErrConnDone)

		_, err := repo.GetSummariesByEmployee(context.Background(), testEmployeeID, "MONTHLY")
		if err == nil {
			t.Fatal(errMsgExpected)
		}
	})
}

// ═══════════════════════════════════════════
// GetCurrentPeriodSummary
// ═══════════════════════════════════════════

func TestGetCurrentPeriodSummary(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo, mock, db := newMockRepo(t)
		defer db.Close()

		rows := sqlmock.NewRows(summaryCols).
			AddRow("s-1", testEmployeeID, "MONTHLY", 2026, 1, testPeriodStart, testPeriodEnd, 9000, 10800, 5, 3, testLastUpdated)
		mock.ExpectQuery(stmtSelect).WithArgs(testEmployeeID, "MONTHLY", 2026, 1).WillReturnRows(rows)

		result, err := repo.GetCurrentPeriodSummary(context.Background(), testEmployeeID, "MONTHLY", 2026, 1)
		if err != nil {
			t.Fatalf(errMsgUnexpected, err)
		}
		if result == nil {
			t.Fatal("expected non-nil result")
		}
		if result.ID != "s-1" {
			t.Errorf("expected s-1, got %s", result.ID)
		}
	})

	t.Run("not found", func(t *testing.T) {
		repo, mock, db := newMockRepo(t)
		defer db.Close()

		rows := sqlmock.NewRows(summaryCols)
		mock.ExpectQuery(stmtSelect).WithArgs(testEmployeeID, "MONTHLY", 2026, 1).WillReturnRows(rows)

		result, err := repo.GetCurrentPeriodSummary(context.Background(), testEmployeeID, "MONTHLY", 2026, 1)
		if err != nil {
			t.Fatalf(errMsgUnexpected, err)
		}
		if result != nil {
			t.Errorf("expected nil for not found, got %+v", result)
		}
	})
}

// ═══════════════════════════════════════════
// GetAllSummaries
// ═══════════════════════════════════════════

func TestGetAllSummaries(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo, mock, db := newMockRepo(t)
		defer db.Close()

		rows := sqlmock.NewRows(summaryCols).
			AddRow("s-1", testEmployeeID, "MONTHLY", 2026, 1, testPeriodStart, testPeriodEnd, 9000, 10800, 5, 3, nil)
		mock.ExpectQuery(stmtSelect).WithArgs(testEmployeeID).WillReturnRows(rows)

		result, err := repo.GetAllSummaries(context.Background(), testEmployeeID)
		if err != nil {
			t.Fatalf(errMsgUnexpected, err)
		}
		if len(result) != 1 {
			t.Errorf("expected 1, got %d", len(result))
		}
	})

	t.Run("error", func(t *testing.T) {
		repo, mock, db := newMockRepo(t)
		defer db.Close()

		mock.ExpectQuery(stmtSelect).WillReturnError(sql.ErrConnDone)

		_, err := repo.GetAllSummaries(context.Background(), testEmployeeID)
		if err == nil {
			t.Fatal(errMsgExpected)
		}
	})
}
