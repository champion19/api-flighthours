package services

import (
	"context"
	"errors"
	"testing"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
)

// airportFakeTx is a test transaction with tracking for commit/rollback (for airport tests)
type airportFakeTx struct {
	committed  bool
	rolledBack bool
}

func (t *airportFakeTx) Commit() error {
	t.committed = true
	return nil
}

func (t *airportFakeTx) Rollback() error {
	t.rolledBack = true
	return nil
}

type fakeAirportRepo struct {
	getByIDFn      func(ctx context.Context, id string) (*domain.Airport, error)
	updateStatusFn func(ctx context.Context, tx output.Tx, id string, status bool) error
	beginTxFn      func(ctx context.Context) (output.Tx, error)
	listAirportsFn func(ctx context.Context, filters map[string]interface{}) ([]domain.Airport, error)
	getByTypeFn    func(ctx context.Context, airportType string) ([]domain.Airport, error)
}

func (f fakeAirportRepo) BeginTx(ctx context.Context) (output.Tx, error) {
	if f.beginTxFn != nil {
		return f.beginTxFn(ctx)
	}
	return &airportFakeTx{}, nil
}

func (f fakeAirportRepo) GetAirportByID(ctx context.Context, id string) (*domain.Airport, error) {
	if f.getByIDFn != nil {
		return f.getByIDFn(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (f fakeAirportRepo) UpdateAirportStatus(ctx context.Context, tx output.Tx, id string, status bool) error {
	if f.updateStatusFn != nil {
		return f.updateStatusFn(ctx, tx, id, status)
	}
	return errors.New("not implemented")
}

func (f fakeAirportRepo) ListAirports(ctx context.Context, filters map[string]interface{}) ([]domain.Airport, error) {
	if f.listAirportsFn != nil {
		return f.listAirportsFn(ctx, filters)
	}
	return nil, errors.New("not implemented")
}

func (f fakeAirportRepo) GetAirportsByType(ctx context.Context, airportType string) ([]domain.Airport, error) {
	if f.getByTypeFn != nil {
		return f.getByTypeFn(ctx, airportType)
	}
	return nil, errors.New("not implemented")
}

func TestAirportService_BeginTx(t *testing.T) {
	ctx := context.Background()

	t.Run("returns tx when repo succeeds", func(t *testing.T) {
		expectedTx := &airportFakeTx{}
		svc := NewAirportService(fakeAirportRepo{
			beginTxFn: func(context.Context) (output.Tx, error) { return expectedTx, nil },
		})

		tx, err := svc.BeginTx(ctx)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if tx != expectedTx {
			t.Error("expected tx to match expected")
		}
	})

	t.Run("returns error when repo fails", func(t *testing.T) {
		repoErr := errors.New("db connection failed")
		svc := NewAirportService(fakeAirportRepo{
			beginTxFn: func(context.Context) (output.Tx, error) { return nil, repoErr },
		})

		_, err := svc.BeginTx(ctx)
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected %v, got %v", repoErr, err)
		}
	})
}

func TestAirportService_GetAirportByID(t *testing.T) {
	ctx := context.Background()

	t.Run("returns airport when found", func(t *testing.T) {
		expectedAirport := &domain.Airport{
			ID:       "airport-123",
			Name:     "El Dorado International",
			IATACode: "BOG",
			Status:   true,
		}
		svc := NewAirportService(fakeAirportRepo{
			getByIDFn: func(context.Context, string) (*domain.Airport, error) {
				return expectedAirport, nil
			},
		})

		result, err := svc.GetAirportByID(ctx, "airport-123")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result.ID != expectedAirport.ID {
			t.Fatalf("expected ID %s, got %s", expectedAirport.ID, result.ID)
		}
		if result.Name != expectedAirport.Name {
			t.Fatalf("expected name %s, got %s", expectedAirport.Name, result.Name)
		}
	})

	t.Run("returns error when not found", func(t *testing.T) {
		svc := NewAirportService(fakeAirportRepo{
			getByIDFn: func(context.Context, string) (*domain.Airport, error) {
				return nil, domain.ErrAirportNotFound
			},
		})

		_, err := svc.GetAirportByID(ctx, "non-existent")
		if !errors.Is(err, domain.ErrAirportNotFound) {
			t.Fatalf("expected %v, got %v", domain.ErrAirportNotFound, err)
		}
	})

	t.Run("propagates repository error", func(t *testing.T) {
		repoErr := errors.New("database connection error")
		svc := NewAirportService(fakeAirportRepo{
			getByIDFn: func(context.Context, string) (*domain.Airport, error) {
				return nil, repoErr
			},
		})

		_, err := svc.GetAirportByID(ctx, "airport-123")
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected %v, got %v", repoErr, err)
		}
	})
}

func TestAirportService_ListAirports(t *testing.T) {
	ctx := context.Background()

	t.Run("success => returns list of airports", func(t *testing.T) {
		expectedAirports := []domain.Airport{
			{ID: "airport-1", Name: "El Dorado", IATACode: "BOG"},
			{ID: "airport-2", Name: "Jose Maria Cordova", IATACode: "MDE"},
		}
		svc := NewAirportService(fakeAirportRepo{
			listAirportsFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.Airport, error) {
				return expectedAirports, nil
			},
		})

		result, err := svc.ListAirports(ctx, nil)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(result) != 2 {
			t.Errorf("expected 2 airports, got %d", len(result))
		}
	})

	t.Run("empty list => returns empty slice", func(t *testing.T) {
		svc := NewAirportService(fakeAirportRepo{
			listAirportsFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.Airport, error) {
				return []domain.Airport{}, nil
			},
		})

		result, err := svc.ListAirports(ctx, nil)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(result) != 0 {
			t.Errorf("expected 0 airports, got %d", len(result))
		}
	})

	t.Run("propagates repository error", func(t *testing.T) {
		repoErr := errors.New("database error")
		svc := NewAirportService(fakeAirportRepo{
			listAirportsFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.Airport, error) {
				return nil, repoErr
			},
		})

		_, err := svc.ListAirports(ctx, nil)
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected %v, got %v", repoErr, err)
		}
	})
}

func TestAirportService_GetAirportsByType(t *testing.T) {
	ctx := context.Background()

	t.Run("success => returns airports of type", func(t *testing.T) {
		expectedAirports := []domain.Airport{
			{ID: "airport-1", Name: "El Dorado", AirportType: "INTERNACIONAL"},
			{ID: "airport-2", Name: "Cali Alfonso", AirportType: "INTERNACIONAL"},
		}
		svc := NewAirportService(fakeAirportRepo{
			getByTypeFn: func(ctx context.Context, airportType string) ([]domain.Airport, error) {
				if airportType != "INTERNACIONAL" {
					t.Errorf("expected type INTERNACIONAL, got %s", airportType)
				}
				return expectedAirports, nil
			},
		})

		result, err := svc.GetAirportsByType(ctx, "INTERNACIONAL")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(result) != 2 {
			t.Errorf("expected 2 airports, got %d", len(result))
		}
	})

	t.Run("empty => returns empty slice", func(t *testing.T) {
		svc := NewAirportService(fakeAirportRepo{
			getByTypeFn: func(ctx context.Context, airportType string) ([]domain.Airport, error) {
				return []domain.Airport{}, nil
			},
		})

		result, err := svc.GetAirportsByType(ctx, "NACIONAL")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(result) != 0 {
			t.Errorf("expected 0 airports, got %d", len(result))
		}
	})

	t.Run("propagates repository error", func(t *testing.T) {
		repoErr := errors.New("database error")
		svc := NewAirportService(fakeAirportRepo{
			getByTypeFn: func(ctx context.Context, airportType string) ([]domain.Airport, error) {
				return nil, repoErr
			},
		})

		_, err := svc.GetAirportsByType(ctx, "INTERNACIONAL")
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected %v, got %v", repoErr, err)
		}
	})
}

func TestAirportService_ActivateAirportTx(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &fakeAirportRepo{
			updateStatusFn: func(ctx context.Context, tx output.Tx, id string, status bool) error {
				if !status {
					t.Error("expected status=true")
				}
				return nil
			},
		}
		svc := NewAirportService(repo)
		err := svc.ActivateAirportTx(context.Background(), &airportFakeTx{}, "airport-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		repo := &fakeAirportRepo{
			updateStatusFn: func(ctx context.Context, tx output.Tx, id string, status bool) error {
				return errors.New("db error")
			},
		}
		svc := NewAirportService(repo)
		err := svc.ActivateAirportTx(context.Background(), &airportFakeTx{}, "airport-1")
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestAirportService_DeactivateAirportTx(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &fakeAirportRepo{
			updateStatusFn: func(ctx context.Context, tx output.Tx, id string, status bool) error {
				if status {
					t.Error("expected status=false")
				}
				return nil
			},
		}
		svc := NewAirportService(repo)
		err := svc.DeactivateAirportTx(context.Background(), &airportFakeTx{}, "airport-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		repo := &fakeAirportRepo{
			updateStatusFn: func(ctx context.Context, tx output.Tx, id string, status bool) error {
				return errors.New("db error")
			},
		}
		svc := NewAirportService(repo)
		err := svc.DeactivateAirportTx(context.Background(), &airportFakeTx{}, "airport-1")
		if err == nil {
			t.Fatal("expected error")
		}
	})
}
