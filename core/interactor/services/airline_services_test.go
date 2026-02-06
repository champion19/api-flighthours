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

func TestAirlineService_ListAirlines(t *testing.T) {
	ctx := context.Background()

	t.Run("returns airlines list", func(t *testing.T) {
		airlines := []domain.Airline{
			{ID: "air-1", AirlineName: "Avianca", AirlineCode: "AV"},
			{ID: "air-2", AirlineName: "LATAM", AirlineCode: "LA"},
		}
		svc := NewAirlineService(fakeAirlineRepo{
			listAirlinesFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.Airline, error) {
				return airlines, nil
			},
		})

		result, err := svc.ListAirlines(ctx, nil)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(result) != 2 {
			t.Fatalf("expected 2 airlines, got %d", len(result))
		}
	})

	t.Run("passes filters to repository", func(t *testing.T) {
		var capturedFilters map[string]interface{}
		svc := NewAirlineService(fakeAirlineRepo{
			listAirlinesFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.Airline, error) {
				capturedFilters = filters
				return []domain.Airline{}, nil
			},
		})

		filters := map[string]interface{}{"status": true}
		svc.ListAirlines(ctx, filters)

		if capturedFilters["status"] != true {
			t.Error("expected filters to be passed to repository")
		}
	})

	t.Run("propagates repository error", func(t *testing.T) {
		repoErr := errors.New("database error")
		svc := NewAirlineService(fakeAirlineRepo{
			listAirlinesFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.Airline, error) {
				return nil, repoErr
			},
		})

		_, err := svc.ListAirlines(ctx, nil)
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected %v, got %v", repoErr, err)
		}
	})
}

func TestAirlineService_UpdateAirlineStatus(t *testing.T) {
	ctx := context.Background()

	t.Run("updates status successfully", func(t *testing.T) {
		tx := &airlineFakeTx{}
		updateCalled := false
		svc := NewAirlineService(fakeAirlineRepo{
			beginTxFn: func(context.Context) (output.Tx, error) { return tx, nil },
			updateStatusFn: func(ctx context.Context, tx output.Tx, id string, active bool) error {
				updateCalled = true
				if id != "airline-123" {
					t.Errorf("expected id 'airline-123', got %s", id)
				}
				if active != true {
					t.Error("expected active to be true")
				}
				return nil
			},
		})

		err := svc.UpdateAirlineStatus(ctx, "airline-123", true)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !updateCalled {
			t.Error("expected update to be called")
		}
		if !tx.committed {
			t.Error("expected transaction to be committed")
		}
	})

	t.Run("propagates repository error", func(t *testing.T) {
		tx := &airlineFakeTx{}
		repoErr := errors.New("update failed")
		svc := NewAirlineService(fakeAirlineRepo{
			beginTxFn: func(context.Context) (output.Tx, error) { return tx, nil },
			updateStatusFn: func(ctx context.Context, tx output.Tx, id string, active bool) error {
				return repoErr
			},
		})

		err := svc.UpdateAirlineStatus(ctx, "airline-123", false)
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected %v, got %v", repoErr, err)
		}
		if !tx.rolledBack {
			t.Error("expected transaction to be rolled back")
		}
	})

	t.Run("returns error if BeginTx fails", func(t *testing.T) {
		beginErr := errors.New("cannot begin transaction")
		svc := NewAirlineService(fakeAirlineRepo{
			beginTxFn: func(context.Context) (output.Tx, error) { return nil, beginErr },
		})

		err := svc.UpdateAirlineStatus(ctx, "airline-123", true)
		if !errors.Is(err, beginErr) {
			t.Fatalf("expected %v, got %v", beginErr, err)
		}
	})
}

func TestAirlineService_ActivateAirline(t *testing.T) {
	ctx := context.Background()

	t.Run("activates airline successfully", func(t *testing.T) {
		tx := &airlineFakeTx{}
		svc := NewAirlineService(fakeAirlineRepo{
			beginTxFn: func(context.Context) (output.Tx, error) { return tx, nil },
			updateStatusFn: func(ctx context.Context, tx output.Tx, id string, active bool) error {
				if !active {
					t.Error("expected active to be true for activation")
				}
				return nil
			},
		})

		err := svc.ActivateAirline(ctx, "airline-123")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !tx.committed {
			t.Error("expected transaction to be committed")
		}
	})

	t.Run("rolls back on error", func(t *testing.T) {
		tx := &airlineFakeTx{}
		repoErr := errors.New("update failed")
		svc := NewAirlineService(fakeAirlineRepo{
			beginTxFn: func(context.Context) (output.Tx, error) { return tx, nil },
			updateStatusFn: func(ctx context.Context, tx output.Tx, id string, active bool) error {
				return repoErr
			},
		})

		err := svc.ActivateAirline(ctx, "airline-123")
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected %v, got %v", repoErr, err)
		}
		if !tx.rolledBack {
			t.Error("expected transaction to be rolled back")
		}
	})

	t.Run("returns error if BeginTx fails", func(t *testing.T) {
		beginErr := errors.New("cannot begin transaction")
		svc := NewAirlineService(fakeAirlineRepo{
			beginTxFn: func(context.Context) (output.Tx, error) { return nil, beginErr },
		})

		err := svc.ActivateAirline(ctx, "airline-123")
		if !errors.Is(err, beginErr) {
			t.Fatalf("expected %v, got %v", beginErr, err)
		}
	})
}

func TestAirlineService_DeactivateAirline(t *testing.T) {
	ctx := context.Background()

	t.Run("deactivates airline successfully", func(t *testing.T) {
		tx := &airlineFakeTx{}
		svc := NewAirlineService(fakeAirlineRepo{
			beginTxFn: func(context.Context) (output.Tx, error) { return tx, nil },
			updateStatusFn: func(ctx context.Context, tx output.Tx, id string, active bool) error {
				if active {
					t.Error("expected active to be false for deactivation")
				}
				return nil
			},
		})

		err := svc.DeactivateAirline(ctx, "airline-123")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !tx.committed {
			t.Error("expected transaction to be committed")
		}
	})

	t.Run("rolls back on error", func(t *testing.T) {
		tx := &airlineFakeTx{}
		repoErr := errors.New("update failed")
		svc := NewAirlineService(fakeAirlineRepo{
			beginTxFn: func(context.Context) (output.Tx, error) { return tx, nil },
			updateStatusFn: func(ctx context.Context, tx output.Tx, id string, active bool) error {
				return repoErr
			},
		})

		err := svc.DeactivateAirline(ctx, "airline-123")
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected %v, got %v", repoErr, err)
		}
		if !tx.rolledBack {
			t.Error("expected transaction to be rolled back")
		}
	})

	t.Run("returns error if BeginTx fails", func(t *testing.T) {
		beginErr := errors.New("cannot begin transaction")
		svc := NewAirlineService(fakeAirlineRepo{
			beginTxFn: func(context.Context) (output.Tx, error) { return nil, beginErr },
		})

		err := svc.DeactivateAirline(ctx, "airline-123")
		if !errors.Is(err, beginErr) {
			t.Fatalf("expected %v, got %v", beginErr, err)
		}
	})
}
