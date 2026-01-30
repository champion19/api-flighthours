package services

import (
	"context"
	"errors"
	"testing"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
)

type airlineFakeTx struct {
	committed  bool
	rolledBack bool
}

func (t *airlineFakeTx) Commit() error {
	t.committed = true
	return nil
}

func (t *airlineFakeTx) Rollback() error {
	t.rolledBack = true
	return nil
}

type fakeAirlineRepo struct {
	getByIDFn      func(ctx context.Context, id string) (*domain.Airline, error)
	beginTxFn      func(ctx context.Context) (output.Tx, error)
	listAirlinesFn func(ctx context.Context, filters map[string]interface{}) ([]domain.Airline, error)
	updateStatusFn func(ctx context.Context, tx output.Tx, id string, active bool) error
}

func (f fakeAirlineRepo) BeginTx(ctx context.Context) (output.Tx, error) {
	if f.beginTxFn != nil {
		return f.beginTxFn(ctx)
	}
	return &airlineFakeTx{}, nil
}

func (f fakeAirlineRepo) GetAirlineByID(ctx context.Context, id string) (*domain.Airline, error) {
	if f.getByIDFn != nil {
		return f.getByIDFn(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (f fakeAirlineRepo) ListAirlines(ctx context.Context, filters map[string]interface{}) ([]domain.Airline, error) {
	if f.listAirlinesFn != nil {
		return f.listAirlinesFn(ctx, filters)
	}
	return nil, errors.New("not implemented")
}

func (f fakeAirlineRepo) UpdateAirlineStatus(ctx context.Context, tx output.Tx, id string, active bool) error {
	if f.updateStatusFn != nil {
		return f.updateStatusFn(ctx, tx, id, active)
	}
	return nil
}

func TestAirlineService_GetAirlineByID(t *testing.T) {
	ctx := context.Background()

	t.Run("returns airline when found", func(t *testing.T) {
		expectedAirline := &domain.Airline{
			ID:          "airline-123",
			AirlineName: "Test Airlines",
			AirlineCode: "TST",
			Status:      "active",
		}
		svc := NewAirlineService(fakeAirlineRepo{
			getByIDFn: func(context.Context, string) (*domain.Airline, error) {
				return expectedAirline, nil
			},
		})

		result, err := svc.GetAirlineByID(ctx, "airline-123")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result.ID != expectedAirline.ID {
			t.Fatalf("expected ID %s, got %s", expectedAirline.ID, result.ID)
		}
		if result.AirlineName != expectedAirline.AirlineName {
			t.Fatalf("expected name %s, got %s", expectedAirline.AirlineName, result.AirlineName)
		}
	})

	t.Run("returns error when not found", func(t *testing.T) {
		svc := NewAirlineService(fakeAirlineRepo{
			getByIDFn: func(context.Context, string) (*domain.Airline, error) {
				return nil, domain.ErrAirlineNotFound
			},
		})

		_, err := svc.GetAirlineByID(ctx, "non-existent")
		if !errors.Is(err, domain.ErrAirlineNotFound) {
			t.Fatalf("expected %v, got %v", domain.ErrAirlineNotFound, err)
		}
	})

	t.Run("propagates repository error", func(t *testing.T) {
		repoErr := errors.New("database connection error")
		svc := NewAirlineService(fakeAirlineRepo{
			getByIDFn: func(context.Context, string) (*domain.Airline, error) {
				return nil, repoErr
			},
		})

		_, err := svc.GetAirlineByID(ctx, "airline-123")
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected %v, got %v", repoErr, err)
		}
	})
}

func TestAirlineService_BeginTx(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		expectedTx := &airlineFakeTx{}
		svc := NewAirlineService(fakeAirlineRepo{
			beginTxFn: func(context.Context) (output.Tx, error) { return expectedTx, nil },
		})

		tx, err := svc.BeginTx(ctx)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if tx == nil {
			t.Fatal("expected tx, got nil")
		}
	})

	t.Run("propagates error", func(t *testing.T) {
		beginErr := errors.New("cannot begin transaction")
		svc := NewAirlineService(fakeAirlineRepo{
			beginTxFn: func(context.Context) (output.Tx, error) { return nil, beginErr },
		})

		_, err := svc.BeginTx(ctx)
		if !errors.Is(err, beginErr) {
			t.Fatalf("expected %v, got %v", beginErr, err)
		}
	})
}
