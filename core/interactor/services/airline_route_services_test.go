package services

import (
	"context"
	"errors"
	"testing"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
)

// Mock transaction for testing
type mockTx struct {
	commitFn   func() error
	rollbackFn func() error
}

func (m *mockTx) Commit() error {
	if m.commitFn != nil {
		return m.commitFn()
	}
	return nil
}

func (m *mockTx) Rollback() error {
	if m.rollbackFn != nil {
		return m.rollbackFn()
	}
	return nil
}

func (m *mockTx) ExecContext(ctx context.Context, query string, args ...interface{}) (interface{}, error) {
	return nil, nil
}

func (m *mockTx) QueryRowContext(ctx context.Context, query string, args ...interface{}) interface{} {
	return nil
}

// Mock airline route repository
type mockAirlineRouteRepo struct {
	getByIDFn       func(ctx context.Context, id string) (*domain.AirlineRoute, error)
	listFn          func(ctx context.Context, filters map[string]interface{}) ([]domain.AirlineRoute, error)
	listByAirlineFn func(ctx context.Context, airlineID string) ([]domain.AirlineRoute, error)
	beginTxFn       func(ctx context.Context) (output.Tx, error)
	updateStatusFn  func(ctx context.Context, tx output.Tx, id string, status bool) error
}

func (m *mockAirlineRouteRepo) GetAirlineRouteByID(ctx context.Context, id string) (*domain.AirlineRoute, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *mockAirlineRouteRepo) ListAirlineRoutes(ctx context.Context, filters map[string]interface{}) ([]domain.AirlineRoute, error) {
	if m.listFn != nil {
		return m.listFn(ctx, filters)
	}
	return nil, nil
}

func (m *mockAirlineRouteRepo) ListAirlineRoutesByAirlineID(ctx context.Context, airlineID string) ([]domain.AirlineRoute, error) {
	if m.listByAirlineFn != nil {
		return m.listByAirlineFn(ctx, airlineID)
	}
	return nil, nil
}

func (m *mockAirlineRouteRepo) BeginTx(ctx context.Context) (output.Tx, error) {
	if m.beginTxFn != nil {
		return m.beginTxFn(ctx)
	}
	return &mockTx{}, nil
}

func (m *mockAirlineRouteRepo) UpdateAirlineRouteStatus(ctx context.Context, tx output.Tx, id string, status bool) error {
	if m.updateStatusFn != nil {
		return m.updateStatusFn(ctx, tx, id, status)
	}
	return nil
}

func TestNewAirlineRouteService(t *testing.T) {
	t.Run("creates airline route service", func(t *testing.T) {
		repo := &mockAirlineRouteRepo{}
		service := NewAirlineRouteService(repo)

		if service == nil {
			t.Error("expected non-nil service")
		}
	})
}

func TestAirlineRouteService_GetAirlineRouteByID(t *testing.T) {
	t.Run("returns airline route when found", func(t *testing.T) {
		expected := &domain.AirlineRoute{
			ID:        "ar-123",
			AirlineID: "airline-1",
			RouteID:   "route-1",
		}

		repo := &mockAirlineRouteRepo{
			getByIDFn: func(ctx context.Context, id string) (*domain.AirlineRoute, error) {
				return expected, nil
			},
		}

		service := NewAirlineRouteService(repo)
		result, err := service.GetAirlineRouteByID(context.Background(), "ar-123")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.ID != "ar-123" {
			t.Errorf("expected ID 'ar-123', got %q", result.ID)
		}
	})

	t.Run("returns error when not found", func(t *testing.T) {
		repo := &mockAirlineRouteRepo{
			getByIDFn: func(ctx context.Context, id string) (*domain.AirlineRoute, error) {
				return nil, errors.New("not found")
			},
		}

		service := NewAirlineRouteService(repo)
		_, err := service.GetAirlineRouteByID(context.Background(), "non-existent")

		if err == nil {
			t.Error("expected error for non-existent airline route")
		}
	})
}

func TestAirlineRouteService_ListAirlineRoutes(t *testing.T) {
	t.Run("returns list of airline routes", func(t *testing.T) {
		expected := []domain.AirlineRoute{
			{ID: "ar-1", AirlineID: "airline-1", RouteID: "route-1"},
			{ID: "ar-2", AirlineID: "airline-1", RouteID: "route-2"},
		}

		repo := &mockAirlineRouteRepo{
			listFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.AirlineRoute, error) {
				return expected, nil
			},
		}

		service := NewAirlineRouteService(repo)
		result, err := service.ListAirlineRoutes(context.Background(), nil)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 2 {
			t.Errorf("expected 2 airline routes, got %d", len(result))
		}
	})
}

func TestAirlineRouteService_ListAirlineRoutesByAirlineID(t *testing.T) {
	t.Run("returns routes for specific airline", func(t *testing.T) {
		expected := []domain.AirlineRoute{
			{ID: "ar-1", AirlineID: "airline-1", RouteID: "route-1"},
		}

		repo := &mockAirlineRouteRepo{
			listByAirlineFn: func(ctx context.Context, airlineID string) ([]domain.AirlineRoute, error) {
				return expected, nil
			},
		}

		service := NewAirlineRouteService(repo)
		result, err := service.ListAirlineRoutesByAirlineID(context.Background(), "airline-1")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 1 {
			t.Errorf("expected 1 airline route, got %d", len(result))
		}
	})
}

func TestAirlineRouteService_ActivateAirlineRoute(t *testing.T) {
	t.Run("activates airline route successfully", func(t *testing.T) {
		repo := &mockAirlineRouteRepo{
			getByIDFn: func(ctx context.Context, id string) (*domain.AirlineRoute, error) {
				return &domain.AirlineRoute{ID: "ar-123"}, nil
			},
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return &mockTx{}, nil
			},
			updateStatusFn: func(ctx context.Context, tx output.Tx, id string, status bool) error {
				if !status {
					t.Error("expected status to be true for activation")
				}
				return nil
			},
		}

		service := NewAirlineRouteService(repo)
		err := service.ActivateAirlineRoute(context.Background(), "ar-123")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("returns error when airline route not found", func(t *testing.T) {
		repo := &mockAirlineRouteRepo{
			getByIDFn: func(ctx context.Context, id string) (*domain.AirlineRoute, error) {
				return nil, errors.New("not found")
			},
		}

		service := NewAirlineRouteService(repo)
		err := service.ActivateAirlineRoute(context.Background(), "non-existent")

		if err == nil {
			t.Error("expected error for non-existent airline route")
		}
	})

	t.Run("returns error when BeginTx fails", func(t *testing.T) {
		repo := &mockAirlineRouteRepo{
			getByIDFn: func(ctx context.Context, id string) (*domain.AirlineRoute, error) {
				return &domain.AirlineRoute{ID: "ar-123"}, nil
			},
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return nil, errors.New("tx failed")
			},
		}

		service := NewAirlineRouteService(repo)
		err := service.ActivateAirlineRoute(context.Background(), "ar-123")

		if err == nil {
			t.Error("expected error when BeginTx fails")
		}
	})

	t.Run("returns error and rolls back when UpdateStatus fails", func(t *testing.T) {
		rollbackCalled := false
		repo := &mockAirlineRouteRepo{
			getByIDFn: func(ctx context.Context, id string) (*domain.AirlineRoute, error) {
				return &domain.AirlineRoute{ID: "ar-123"}, nil
			},
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return &mockTx{
					rollbackFn: func() error {
						rollbackCalled = true
						return nil
					},
				}, nil
			},
			updateStatusFn: func(ctx context.Context, tx output.Tx, id string, status bool) error {
				return errors.New("update failed")
			},
		}

		service := NewAirlineRouteService(repo)
		err := service.ActivateAirlineRoute(context.Background(), "ar-123")

		if err == nil {
			t.Error("expected error when UpdateStatus fails")
		}
		if !rollbackCalled {
			t.Error("expected rollback to be called")
		}
	})

	t.Run("returns error when Commit fails", func(t *testing.T) {
		repo := &mockAirlineRouteRepo{
			getByIDFn: func(ctx context.Context, id string) (*domain.AirlineRoute, error) {
				return &domain.AirlineRoute{ID: "ar-123"}, nil
			},
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return &mockTx{
					commitFn: func() error { return errors.New("commit failed") },
				}, nil
			},
			updateStatusFn: func(ctx context.Context, tx output.Tx, id string, status bool) error {
				return nil
			},
		}

		service := NewAirlineRouteService(repo)
		err := service.ActivateAirlineRoute(context.Background(), "ar-123")

		if err == nil {
			t.Error("expected error when Commit fails")
		}
	})
}

func TestAirlineRouteService_DeactivateAirlineRoute(t *testing.T) {
	t.Run("deactivates airline route successfully", func(t *testing.T) {
		repo := &mockAirlineRouteRepo{
			getByIDFn: func(ctx context.Context, id string) (*domain.AirlineRoute, error) {
				return &domain.AirlineRoute{ID: "ar-123"}, nil
			},
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return &mockTx{}, nil
			},
			updateStatusFn: func(ctx context.Context, tx output.Tx, id string, status bool) error {
				if status {
					t.Error("expected status to be false for deactivation")
				}
				return nil
			},
		}

		service := NewAirlineRouteService(repo)
		err := service.DeactivateAirlineRoute(context.Background(), "ar-123")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("returns error when airline route not found", func(t *testing.T) {
		repo := &mockAirlineRouteRepo{
			getByIDFn: func(ctx context.Context, id string) (*domain.AirlineRoute, error) {
				return nil, errors.New("not found")
			},
		}

		service := NewAirlineRouteService(repo)
		err := service.DeactivateAirlineRoute(context.Background(), "non-existent")

		if err == nil {
			t.Error("expected error for non-existent airline route")
		}
	})

	t.Run("returns error when BeginTx fails", func(t *testing.T) {
		repo := &mockAirlineRouteRepo{
			getByIDFn: func(ctx context.Context, id string) (*domain.AirlineRoute, error) {
				return &domain.AirlineRoute{ID: "ar-123"}, nil
			},
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return nil, errors.New("tx failed")
			},
		}

		service := NewAirlineRouteService(repo)
		err := service.DeactivateAirlineRoute(context.Background(), "ar-123")

		if err == nil {
			t.Error("expected error when BeginTx fails")
		}
	})

	t.Run("returns error and rolls back when UpdateStatus fails", func(t *testing.T) {
		rollbackCalled := false
		repo := &mockAirlineRouteRepo{
			getByIDFn: func(ctx context.Context, id string) (*domain.AirlineRoute, error) {
				return &domain.AirlineRoute{ID: "ar-123"}, nil
			},
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return &mockTx{
					rollbackFn: func() error {
						rollbackCalled = true
						return nil
					},
				}, nil
			},
			updateStatusFn: func(ctx context.Context, tx output.Tx, id string, status bool) error {
				return errors.New("update failed")
			},
		}

		service := NewAirlineRouteService(repo)
		err := service.DeactivateAirlineRoute(context.Background(), "ar-123")

		if err == nil {
			t.Error("expected error when UpdateStatus fails")
		}
		if !rollbackCalled {
			t.Error("expected rollback to be called")
		}
	})

	t.Run("returns error when Commit fails", func(t *testing.T) {
		repo := &mockAirlineRouteRepo{
			getByIDFn: func(ctx context.Context, id string) (*domain.AirlineRoute, error) {
				return &domain.AirlineRoute{ID: "ar-123"}, nil
			},
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return &mockTx{
					commitFn: func() error { return errors.New("commit failed") },
				}, nil
			},
			updateStatusFn: func(ctx context.Context, tx output.Tx, id string, status bool) error {
				return nil
			},
		}

		service := NewAirlineRouteService(repo)
		err := service.DeactivateAirlineRoute(context.Background(), "ar-123")

		if err == nil {
			t.Error("expected error when Commit fails")
		}
	})
}

func TestAirlineRouteService_BeginTx(t *testing.T) {
	t.Run("begins transaction successfully", func(t *testing.T) {
		repo := &mockAirlineRouteRepo{
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return &mockTx{}, nil
			},
		}

		service := NewAirlineRouteService(repo)
		tx, err := service.BeginTx(context.Background())

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tx == nil {
			t.Error("expected non-nil transaction")
		}
	})
}
