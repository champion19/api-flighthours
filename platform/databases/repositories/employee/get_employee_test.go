package employee

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

// Employee getters use r.db.QueryRowContext (not prepared stmt), 11 columns

func TestGetEmployeeByEmail_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	r := &repository{db: db}

	cols := []string{"id", "name", "airline", "email", "identification_number", "bp", "start_date", "end_date", "active", "role", "keycloak_user_id"}
	mock.ExpectQuery("SELECT").WithArgs("test@test.com").WillReturnRows(
		sqlmock.NewRows(cols).AddRow("e1", "John", sql.NullString{String: "a1", Valid: true}, "test@test.com", "123", sql.NullString{}, sql.NullString{}, sql.NullString{}, true, "pilot", sql.NullString{String: "kc1", Valid: true}),
	)

	result, err := r.GetEmployeeByEmail(context.Background(), "test@test.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Email != "test@test.com" {
		t.Errorf("expected test@test.com, got %s", result.Email)
	}
}

func TestGetEmployeeByEmail_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	r := &repository{db: db}

	mock.ExpectQuery("SELECT").WithArgs("unknown@test.com").WillReturnError(sql.ErrNoRows)

	_, err = r.GetEmployeeByEmail(context.Background(), "unknown@test.com")
	if !errors.Is(err, domain.ErrPersonNotFound) {
		t.Fatalf("expected ErrPersonNotFound, got %v", err)
	}
}

func TestGetEmployeeByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	r := &repository{db: db}

	cols := []string{"id", "name", "airline", "email", "identification_number", "bp", "start_date", "end_date", "active", "role", "keycloak_user_id"}
	mock.ExpectQuery("SELECT").WithArgs("e1").WillReturnRows(
		sqlmock.NewRows(cols).AddRow("e1", "John", sql.NullString{String: "a1", Valid: true}, "test@test.com", "123", sql.NullString{}, sql.NullString{}, sql.NullString{}, true, "pilot", sql.NullString{String: "kc1", Valid: true}),
	)

	result, err := r.GetEmployeeByID(context.Background(), "e1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "e1" {
		t.Errorf("expected e1, got %s", result.ID)
	}
}

func TestGetEmployeeByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	r := &repository{db: db}

	mock.ExpectQuery("SELECT").WithArgs("unknown").WillReturnError(sql.ErrNoRows)

	_, err = r.GetEmployeeByID(context.Background(), "unknown")
	if !errors.Is(err, domain.ErrPersonNotFound) {
		t.Fatalf("expected ErrPersonNotFound, got %v", err)
	}
}

func TestGetEmployeeByKeycloakID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	r := &repository{db: db}

	cols := []string{"id", "name", "airline", "email", "identification_number", "bp", "start_date", "end_date", "active", "role", "keycloak_user_id"}
	mock.ExpectQuery("SELECT").WithArgs("kc1").WillReturnRows(
		sqlmock.NewRows(cols).AddRow("e1", "John", sql.NullString{String: "a1", Valid: true}, "test@test.com", "123", sql.NullString{}, sql.NullString{}, sql.NullString{}, true, "pilot", sql.NullString{String: "kc1", Valid: true}),
	)

	result, err := r.GetEmployeeByKeycloakID(context.Background(), "kc1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "e1" {
		t.Errorf("expected e1, got %s", result.ID)
	}
}

func TestGetEmployeeByKeycloakID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	r := &repository{db: db}

	mock.ExpectQuery("SELECT").WithArgs("unknown-kc").WillReturnError(sql.ErrNoRows)

	_, err = r.GetEmployeeByKeycloakID(context.Background(), "unknown-kc")
	if !errors.Is(err, domain.ErrPersonNotFound) {
		t.Fatalf("expected ErrPersonNotFound, got %v", err)
	}
}
