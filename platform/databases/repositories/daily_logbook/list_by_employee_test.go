package daily_logbook

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestListDailyLogbooksByEmployee_All(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtByEmp, _ := db.Prepare(QueryByEmployee)

	r := &repository{stmtGetByEmployee: stmtByEmp}

	rows := sqlmock.NewRows([]string{"id", "log_date", "employee_id", "book_page", "status", "tail_number_id", "tail_number", "crew_role"}).
		AddRow("dl1", time.Now(), "emp1", 1, true, nil, nil, nil)
	prep.ExpectQuery().WithArgs("emp1").WillReturnRows(rows)

	result, err := r.ListDailyLogbooksByEmployee(context.Background(), "emp1", map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
}

func TestListDailyLogbooksByEmployee_ByStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtByEmpStatus, _ := db.Prepare(QueryByEmployeeAndStatus)

	r := &repository{stmtGetByEmployeeAndStatus: stmtByEmpStatus}

	rows := sqlmock.NewRows([]string{"id", "log_date", "employee_id", "book_page", "status", "tail_number_id", "tail_number", "crew_role"}).
		AddRow("dl1", time.Now(), "emp1", 1, true, nil, nil, nil)
	prep.ExpectQuery().WithArgs("emp1", true).WillReturnRows(rows)

	result, err := r.ListDailyLogbooksByEmployee(context.Background(), "emp1", map[string]interface{}{"status": true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
}

func TestListDailyLogbooksByEmployee_StatusNonBool(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtByEmp, _ := db.Prepare(QueryByEmployee)

	r := &repository{stmtGetByEmployee: stmtByEmp}

	rows := sqlmock.NewRows([]string{"id", "log_date", "employee_id", "book_page", "status"})
	prep.ExpectQuery().WithArgs("emp1").WillReturnRows(rows)

	_, err = r.ListDailyLogbooksByEmployee(context.Background(), "emp1", map[string]interface{}{"status": "invalid"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListDailyLogbooksByEmployee_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtByEmp, _ := db.Prepare(QueryByEmployee)

	r := &repository{stmtGetByEmployee: stmtByEmp}

	prep.ExpectQuery().WithArgs("emp1").WillReturnError(errors.New("db error"))

	_, err = r.ListDailyLogbooksByEmployee(context.Background(), "emp1", map[string]interface{}{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestListDailyLogbooksByEmployee_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("SELECT")
	stmtByEmp, _ := db.Prepare(QueryByEmployee)

	r := &repository{stmtGetByEmployee: stmtByEmp}

	rows := sqlmock.NewRows([]string{"id"}).AddRow("one")
	prep.ExpectQuery().WithArgs("emp1").WillReturnRows(rows)

	_, err = r.ListDailyLogbooksByEmployee(context.Background(), "emp1", map[string]interface{}{})
	if err == nil {
		t.Fatal("expected scan error")
	}
}
