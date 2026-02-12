package interactor

import (
	"context"
	"errors"
	"testing"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/input"
	"github.com/champion19/api-flighthours/core/ports/output"
)

// fakeAirlineService implements input.AirlineService for testing
type fakeAirlineService struct {
	beginTxFn      func(ctx context.Context) (output.Tx, error)
	getByIDFn      func(ctx context.Context, id string) (*domain.Airline, error)
	listFn         func(ctx context.Context, filters map[string]interface{}) ([]domain.Airline, error)
	updateStatusFn func(ctx context.Context, id string, status bool) error
	activateFn     func(ctx context.Context, id string) error
	deactivateFn   func(ctx context.Context, id string) error
	activateTxFn   func(ctx context.Context, tx output.Tx, id string) error
	deactivateTxFn func(ctx context.Context, tx output.Tx, id string) error
}

var _ input.AirlineService = (*fakeAirlineService)(nil)

func (f *fakeAirlineService) BeginTx(ctx context.Context) (output.Tx, error) {
	if f.beginTxFn != nil {
		return f.beginTxFn(ctx)
	}
	return &fakeTx{}, nil
}

func (f *fakeAirlineService) GetAirlineByID(ctx context.Context, id string) (*domain.Airline, error) {
	if f.getByIDFn != nil {
		return f.getByIDFn(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (f *fakeAirlineService) ListAirlines(ctx context.Context, filters map[string]interface{}) ([]domain.Airline, error) {
	if f.listFn != nil {
		return f.listFn(ctx, filters)
	}
	return nil, errors.New("not implemented")
}

func (f *fakeAirlineService) UpdateAirlineStatus(ctx context.Context, id string, status bool) error {
	if f.updateStatusFn != nil {
		return f.updateStatusFn(ctx, id, status)
	}
	return nil
}

func (f *fakeAirlineService) ActivateAirline(ctx context.Context, id string) error {
	if f.activateFn != nil {
		return f.activateFn(ctx, id)
	}
	return nil
}

func (f *fakeAirlineService) DeactivateAirline(ctx context.Context, id string) error {
	if f.deactivateFn != nil {
		return f.deactivateFn(ctx, id)
	}
	return nil
}

func (f *fakeAirlineService) ActivateAirlineTx(ctx context.Context, tx output.Tx, id string) error {
	if f.activateTxFn != nil {
		return f.activateTxFn(ctx, tx, id)
	}
	return nil
}

func (f *fakeAirlineService) DeactivateAirlineTx(ctx context.Context, tx output.Tx, id string) error {
	if f.deactivateTxFn != nil {
		return f.deactivateTxFn(ctx, tx, id)
	}
	return nil
}

func TestNewAirlineInteractor(t *testing.T) {
	svc := &fakeAirlineService{}
	interactor := NewAirlineInteractor(svc)
	if interactor == nil {
		t.Error("expected non-nil AirlineInteractor")
	}
}

func TestAirlineInteractor_GetAirlineByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expectedAirline := &domain.Airline{
			ID:          "airline-123",
			AirlineName: "Avianca",
			AirlineCode: "AV",
		}
		svc := &fakeAirlineService{
			getByIDFn: func(ctx context.Context, id string) (*domain.Airline, error) {
				return expectedAirline, nil
			},
		}
		interactor := NewAirlineInteractor(svc)

		result, err := interactor.GetAirlineByID(context.Background(), "airline-123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.ID != expectedAirline.ID {
			t.Errorf("expected ID %q, got %q", expectedAirline.ID, result.ID)
		}
		if result.AirlineName != expectedAirline.AirlineName {
			t.Errorf("expected AirlineName %q, got %q", expectedAirline.AirlineName, result.AirlineName)
		}
	})

	t.Run("not found", func(t *testing.T) {
		svc := &fakeAirlineService{
			getByIDFn: func(ctx context.Context, id string) (*domain.Airline, error) {
				return nil, domain.ErrAirlineNotFound
			},
		}
		interactor := NewAirlineInteractor(svc)

		_, err := interactor.GetAirlineByID(context.Background(), "nonexistent")
		if err != domain.ErrAirlineNotFound {
			t.Errorf("expected ErrAirlineNotFound, got %v", err)
		}
	})

	t.Run("service error", func(t *testing.T) {
		svc := &fakeAirlineService{
			getByIDFn: func(ctx context.Context, id string) (*domain.Airline, error) {
				return nil, errors.New("database error")
			},
		}
		interactor := NewAirlineInteractor(svc)

		_, err := interactor.GetAirlineByID(context.Background(), "airline-123")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestAirlineInteractor_ListAirlines(t *testing.T) {
	t.Run("success without filters", func(t *testing.T) {
		expectedAirlines := []domain.Airline{
			{ID: "airline-1", AirlineName: "Avianca", AirlineCode: "AV"},
			{ID: "airline-2", AirlineName: "LATAM", AirlineCode: "LA"},
		}
		svc := &fakeAirlineService{
			listFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.Airline, error) {
				return expectedAirlines, nil
			},
		}
		interactor := NewAirlineInteractor(svc)

		result, err := interactor.ListAirlines(context.Background(), nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 2 {
			t.Errorf("expected 2 airlines, got %d", len(result))
		}
	})

	t.Run("success with filters", func(t *testing.T) {
		svc := &fakeAirlineService{
			listFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.Airline, error) {
				// Verify filters are passed correctly
				if _, ok := filters["status"]; !ok {
					t.Error("expected 'status' filter to be present")
				}
				return []domain.Airline{{ID: "airline-1"}}, nil
			},
		}
		interactor := NewAirlineInteractor(svc)

		filters := map[string]interface{}{"status": true}
		result, err := interactor.ListAirlines(context.Background(), filters)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 1 {
			t.Errorf("expected 1 airline, got %d", len(result))
		}
	})

	t.Run("service error", func(t *testing.T) {
		svc := &fakeAirlineService{
			listFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.Airline, error) {
				return nil, errors.New("database error")
			},
		}
		interactor := NewAirlineInteractor(svc)

		_, err := interactor.ListAirlines(context.Background(), nil)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestAirlineInteractor_ActivateAirline(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &fakeAirlineService{
			getByIDFn: func(ctx context.Context, id string) (*domain.Airline, error) {
				return &domain.Airline{ID: id, AirlineName: "Test"}, nil
			},
			activateTxFn: func(ctx context.Context, tx output.Tx, id string) error {
				return nil
			},
		}
		interactor := NewAirlineInteractor(svc)

		err := interactor.ActivateAirline(context.Background(), "airline-123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		svc := &fakeAirlineService{
			getByIDFn: func(ctx context.Context, id string) (*domain.Airline, error) {
				return nil, domain.ErrAirlineNotFound
			},
		}
		interactor := NewAirlineInteractor(svc)

		err := interactor.ActivateAirline(context.Background(), "nonexistent")
		if err != domain.ErrAirlineNotFound {
			t.Errorf("expected ErrAirlineNotFound, got %v", err)
		}
	})

	t.Run("activation service error", func(t *testing.T) {
		svc := &fakeAirlineService{
			getByIDFn: func(ctx context.Context, id string) (*domain.Airline, error) {
				return &domain.Airline{ID: id}, nil
			},
			activateTxFn: func(ctx context.Context, tx output.Tx, id string) error {
				return errors.New("activation failed")
			},
		}
		interactor := NewAirlineInteractor(svc)

		err := interactor.ActivateAirline(context.Background(), "airline-123")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestAirlineInteractor_DeactivateAirline(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &fakeAirlineService{
			getByIDFn: func(ctx context.Context, id string) (*domain.Airline, error) {
				return &domain.Airline{ID: id, AirlineName: "Test"}, nil
			},
			deactivateTxFn: func(ctx context.Context, tx output.Tx, id string) error {
				return nil
			},
		}
		interactor := NewAirlineInteractor(svc)

		err := interactor.DeactivateAirline(context.Background(), "airline-123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		svc := &fakeAirlineService{
			getByIDFn: func(ctx context.Context, id string) (*domain.Airline, error) {
				return nil, domain.ErrAirlineNotFound
			},
		}
		interactor := NewAirlineInteractor(svc)

		err := interactor.DeactivateAirline(context.Background(), "nonexistent")
		if err != domain.ErrAirlineNotFound {
			t.Errorf("expected ErrAirlineNotFound, got %v", err)
		}
	})

	t.Run("deactivation service error", func(t *testing.T) {
		svc := &fakeAirlineService{
			getByIDFn: func(ctx context.Context, id string) (*domain.Airline, error) {
				return &domain.Airline{ID: id}, nil
			},
			deactivateTxFn: func(ctx context.Context, tx output.Tx, id string) error {
				return errors.New("deactivation failed")
			},
		}
		interactor := NewAirlineInteractor(svc)

		err := interactor.DeactivateAirline(context.Background(), "airline-123")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
