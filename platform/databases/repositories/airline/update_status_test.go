package airline

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/platform/databases/common"
)

type fakeTxAL struct{}

func (f *fakeTxAL) Commit() error   { return nil }
func (f *fakeTxAL) Rollback() error { return nil }

func TestUpdateAirlineStatus_InvalidTx(t *testing.T) {
	r := &repository{}
	err := r.UpdateAirlineStatus(context.Background(), &fakeTxAL{}, "id-1", true)
	if !errors.Is(err, domain.ErrInvalidTransaction) {
		t.Fatalf("expected ErrInvalidTransaction, got %v", err)
	}
}

func TestUpdateAirlineStatus_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs("id-1").
		WillReturnError(errors.New("query error"))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.UpdateAirlineStatus(context.Background(), sqlTx, "id-1", true)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUpdateAirlineStatus_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs("id-1").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.UpdateAirlineStatus(context.Background(), sqlTx, "id-1", true)
	if !errors.Is(err, domain.ErrAirlineNotFound) {
		t.Fatalf("expected ErrAirlineNotFound, got %v", err)
	}
}

func TestUpdateAirlineStatus_ExecError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs("id-1").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
	mock.ExpectExec("UPDATE airline SET status").WillReturnError(errors.New("exec fail"))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.UpdateAirlineStatus(context.Background(), sqlTx, "id-1", true)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUpdateAirlineStatus_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs("id-1").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
	mock.ExpectExec("UPDATE airline SET status").WillReturnResult(sqlmock.NewResult(0, 1))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.UpdateAirlineStatus(context.Background(), sqlTx, "id-1", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
