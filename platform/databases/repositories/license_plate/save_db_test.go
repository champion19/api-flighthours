package licenseplate

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/platform/databases/common"
)

type fakeTxLP struct{}

func (f *fakeTxLP) Commit() error   { return nil }
func (f *fakeTxLP) Rollback() error { return nil }

func TestSaveLicensePlate_InvalidTx(t *testing.T) {
	r := &repository{}
	err := r.SaveLicensePlate(context.Background(), &fakeTxLP{}, domain.LicensePlate{})
	if !errors.Is(err, domain.ErrInvalidTransaction) {
		t.Fatalf("expected ErrInvalidTransaction, got %v", err)
	}
}

func TestSaveLicensePlate_DuplicateEntry(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO license_plate").WillReturnError(errors.New("Duplicate entry 'HK-5432' for key"))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.SaveLicensePlate(context.Background(), sqlTx, domain.LicensePlate{ID: "lp-1", LicensePlate: "HK-5432"})
	if !errors.Is(err, domain.ErrLicensePlateDuplicatePlate) {
		t.Fatalf("expected ErrLicensePlateDuplicatePlate, got %v", err)
	}
}

func TestSaveLicensePlate_FKAircraftModel(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO license_plate").WillReturnError(errors.New("foreign key constraint fails aircraft_model"))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.SaveLicensePlate(context.Background(), sqlTx, domain.LicensePlate{ID: "lp-1"})
	if !errors.Is(err, domain.ErrLicensePlateInvalidModel) {
		t.Fatalf("expected ErrLicensePlateInvalidModel, got %v", err)
	}
}

func TestSaveLicensePlate_FKAirline(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO license_plate").WillReturnError(errors.New("foreign key constraint fails airline"))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.SaveLicensePlate(context.Background(), sqlTx, domain.LicensePlate{ID: "lp-1"})
	if !errors.Is(err, domain.ErrLicensePlateInvalidAirline) {
		t.Fatalf("expected ErrLicensePlateInvalidAirline, got %v", err)
	}
}

func TestSaveLicensePlate_GenericError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO license_plate").WillReturnError(errors.New("unknown error"))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.SaveLicensePlate(context.Background(), sqlTx, domain.LicensePlate{ID: "lp-1"})
	if !errors.Is(err, domain.ErrLicensePlateCannotSave) {
		t.Fatalf("expected ErrLicensePlateCannotSave, got %v", err)
	}
}

func TestSaveLicensePlate_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO license_plate").WillReturnResult(sqlmock.NewResult(1, 1))

	tx, _ := db.Begin()
	sqlTx := common.NewSQLTx(tx)
	r := &repository{}

	err = r.SaveLicensePlate(context.Background(), sqlTx, domain.LicensePlate{
		ID:              "lp-1",
		LicensePlate:    "HK-1234",
		AircraftModelID: "am-1",
		AirlineID:       "al-1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
