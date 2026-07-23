package airline_route

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

type fakeTxAR struct{}

func (f *fakeTxAR) Commit() error   { return nil }
func (f *fakeTxAR) Rollback() error { return nil }

func TestUpdateAirlineRouteStatus_InvalidTx(t *testing.T) {
	r := &repository{}
	err := r.UpdateAirlineRouteStatus(context.Background(), &fakeTxAR{}, "id-1", domain.AirlineRouteStatusActive)
	if !errors.Is(err, domain.ErrInvalidTransaction) {
		t.Fatalf("expected ErrInvalidTransaction, got %v", err)
	}
}

func TestUpdateAirlineRouteStatus_ExecError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Prepare the stmtUpdateStatus
	stmtUpdate := mock.ExpectPrepare("UPDATE airline_route")

	mock.ExpectBegin()
	stmtUpdate.ExpectExec().
		WithArgs(domain.AirlineRouteStatusActive, "id-1", domain.AirlineRouteStatusActive).
		WillReturnError(errors.New("exec error"))

	stmt, _ := db.Prepare(QueryUpdateStatus)
	r := &repository{stmtUpdateStatus: stmt, db: db}

	tx, _ := db.Begin()
	err = r.UpdateAirlineRouteStatus(context.Background(), tx, "id-1", domain.AirlineRouteStatusActive)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUpdateAirlineRouteStatus_RowsAffected_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	stmtUpdate := mock.ExpectPrepare("UPDATE airline_route")

	mock.ExpectBegin()
	stmtUpdate.ExpectExec().
		WithArgs(domain.AirlineRouteStatusActive, "id-1", domain.AirlineRouteStatusActive).
		WillReturnResult(sqlmock.NewResult(0, 1))

	stmt, _ := db.Prepare(QueryUpdateStatus)
	r := &repository{stmtUpdateStatus: stmt, db: db}

	tx, _ := db.Begin()
	err = r.UpdateAirlineRouteStatus(context.Background(), tx, "id-1", domain.AirlineRouteStatusActive)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
