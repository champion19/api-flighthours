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
	getByCityFn    func(ctx context.Context, city string) ([]domain.Airport, error)
	getByCountryFn func(ctx context.Context, country string) ([]domain.Airport, error)
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

func (f fakeAirportRepo) GetAirportsByCity(ctx context.Context, city string) ([]domain.Airport, error) {
	if f.getByCityFn != nil {
		return f.getByCityFn(ctx, city)
	}
	return nil, errors.New("not implemented")
}

func (f fakeAirportRepo) GetAirportsByCountry(ctx context.Context, country string) ([]domain.Airport, error) {
	if f.getByCountryFn != nil {
		return f.getByCountryFn(ctx, country)
	}
	return nil, errors.New("not implemented")
}

func (f fakeAirportRepo) GetAirportsByType(ctx context.Context, airportType string) ([]domain.Airport, error) {
	if f.getByTypeFn != nil {
		return f.getByTypeFn(ctx, airportType)
	}
	return nil, errors.New("not implemented")
}

func TestAirportService_GetAirportByID(t *testing.T) {
	ctx := context.Background()

	t.Run("returns airport when found", func(t *testing.T) {
		expectedAirport := &domain.Airport{
			ID:       "airport-123",
			Name:     "El Dorado International",
			City:     "Bogota",
			Country:  "Colombia",
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

func TestAirportService_UpdateAirportStatus(t *testing.T) {
	ctx := context.Background()

	t.Run("success => commit", func(t *testing.T) {
		tx := &airportFakeTx{}
		svc := NewAirportService(fakeAirportRepo{
			beginTxFn: func(context.Context) (output.Tx, error) { return tx, nil },
			updateStatusFn: func(context.Context, output.Tx, string, bool) error {
				return nil
			},
		})

		err := svc.UpdateAirportStatus(ctx, "airport-123", true)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !tx.committed {
			t.Fatalf("expected tx.Commit to be called")
		}
	})

	t.Run("begin tx fails => returns error", func(t *testing.T) {
		beginErr := errors.New("cannot begin transaction")
		svc := NewAirportService(fakeAirportRepo{
			beginTxFn: func(context.Context) (output.Tx, error) { return nil, beginErr },
		})

		err := svc.UpdateAirportStatus(ctx, "airport-123", true)
		if !errors.Is(err, beginErr) {
			t.Fatalf("expected %v, got %v", beginErr, err)
		}
	})

	t.Run("update fails => rollback", func(t *testing.T) {
		tx := &airportFakeTx{}
		updateErr := errors.New("update failed")
		svc := NewAirportService(fakeAirportRepo{
			beginTxFn: func(context.Context) (output.Tx, error) { return tx, nil },
			updateStatusFn: func(context.Context, output.Tx, string, bool) error {
				return updateErr
			},
		})

		err := svc.UpdateAirportStatus(ctx, "airport-123", true)
		if !errors.Is(err, updateErr) {
			t.Fatalf("expected %v, got %v", updateErr, err)
		}
	})
}

func TestAirportService_DeactivateAirport(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		tx := &airportFakeTx{}
		var receivedStatus bool
		svc := NewAirportService(fakeAirportRepo{
			beginTxFn: func(context.Context) (output.Tx, error) { return tx, nil },
			updateStatusFn: func(_ context.Context, _ output.Tx, _ string, status bool) error {
				receivedStatus = status
				return nil
			},
		})

		err := svc.DeactivateAirport(ctx, "airport-123")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if receivedStatus {
			t.Fatalf("expected status false, got true")
		}
	})

	t.Run("propagates error", func(t *testing.T) {
		tx := &airportFakeTx{}
		updateErr := errors.New("update failed")
		svc := NewAirportService(fakeAirportRepo{
			beginTxFn: func(context.Context) (output.Tx, error) { return tx, nil },
			updateStatusFn: func(context.Context, output.Tx, string, bool) error {
				return updateErr
			},
		})

		err := svc.DeactivateAirport(ctx, "airport-123")
		if !errors.Is(err, updateErr) {
			t.Fatalf("expected %v, got %v", updateErr, err)
		}
	})
}
