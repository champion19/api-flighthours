package dependency

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/champion19/api-flighthours/platform/logger"
)

func TestInitFlightSummaryInteractor(t *testing.T) {
	log := logger.NewSlogLogger()

	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}
		defer db.Close()

		// 4 prepared statements for flight_summary repo
		mock.ExpectPrepare("SELECT")
		mock.ExpectPrepare("SELECT COUNT")
		mock.ExpectPrepare("SELECT COALESCE")
		mock.ExpectPrepare("SELECT")

		result := initFlightSummaryInteractor(db, log)
		if result == nil {
			t.Error("expected non-nil interactor")
		}
	})

	t.Run("prepare error returns nil", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}
		defer db.Close()

		mock.ExpectPrepare("SELECT").WillReturnError(sql.ErrConnDone)

		result := initFlightSummaryInteractor(db, log)
		if result != nil {
			t.Error("expected nil interactor on error")
		}
	})
}

func TestInitEmployeeFlightSummaryService(t *testing.T) {
	log := logger.NewSlogLogger()

	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}
		defer db.Close()

		// 4 prepared statements for employee_flight_summary repo
		mock.ExpectPrepare("INSERT INTO")
		mock.ExpectPrepare("SELECT id, employee_id")
		mock.ExpectPrepare("SELECT id, employee_id")
		mock.ExpectPrepare("SELECT id, employee_id")

		result := initEmployeeFlightSummaryService(db, log)
		if result == nil {
			t.Error("expected non-nil service")
		}
	})

	t.Run("prepare error returns nil", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}
		defer db.Close()

		mock.ExpectPrepare("INSERT INTO").WillReturnError(sql.ErrConnDone)

		result := initEmployeeFlightSummaryService(db, log)
		if result != nil {
			t.Error("expected nil service on error")
		}
	})
}
