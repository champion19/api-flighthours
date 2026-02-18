package interactor

import (
	"context"
	"errors"
	"testing"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/input"
	"github.com/champion19/api-flighthours/core/ports/output"
)

type airportFakeTx struct{}

func (t *airportFakeTx) Commit() error   { return nil }
func (t *airportFakeTx) Rollback() error { return nil }

type fakeAirportServiceForInteractor struct {
	getByIDFn      func(ctx context.Context, id string) (*domain.Airport, error)
	updateStatusFn func(ctx context.Context, id string, status bool) error
	activateFn     func(ctx context.Context, id string) error
	deactivateFn   func(ctx context.Context, id string) error
	beginTxFn      func(ctx context.Context) (output.Tx, error)
	listAirportsFn func(ctx context.Context, filters map[string]interface{}) ([]domain.Airport, error)
	getByTypeFn    func(ctx context.Context, airportType string) ([]domain.Airport, error)
	activateTxFn   func(ctx context.Context, tx output.Tx, id string) error
	deactivateTxFn func(ctx context.Context, tx output.Tx, id string) error
}

var _ input.AirportService = (*fakeAirportServiceForInteractor)(nil)

func (f *fakeAirportServiceForInteractor) BeginTx(ctx context.Context) (output.Tx, error) {
	if f.beginTxFn != nil {
		return f.beginTxFn(ctx)
	}
	return &airportFakeTx{}, nil
}

func (f *fakeAirportServiceForInteractor) GetAirportByID(ctx context.Context, id string) (*domain.Airport, error) {
	if f.getByIDFn != nil {
		return f.getByIDFn(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (f *fakeAirportServiceForInteractor) UpdateAirportStatus(ctx context.Context, id string, status bool) error {
	if f.updateStatusFn != nil {
		return f.updateStatusFn(ctx, id, status)
	}
	return errors.New("not implemented")
}

func (f *fakeAirportServiceForInteractor) ActivateAirport(ctx context.Context, id string) error {
	if f.activateFn != nil {
		return f.activateFn(ctx, id)
	}
	return errors.New("not implemented")
}

func (f *fakeAirportServiceForInteractor) DeactivateAirport(ctx context.Context, id string) error {
	if f.deactivateFn != nil {
		return f.deactivateFn(ctx, id)
	}
	return errors.New("not implemented")
}

func (f *fakeAirportServiceForInteractor) ListAirports(ctx context.Context, filters map[string]interface{}) ([]domain.Airport, error) {
	if f.listAirportsFn != nil {
		return f.listAirportsFn(ctx, filters)
	}
	return nil, errors.New("not implemented")
}

func (f *fakeAirportServiceForInteractor) GetAirportsByType(ctx context.Context, airportType string) ([]domain.Airport, error) {
	if f.getByTypeFn != nil {
		return f.getByTypeFn(ctx, airportType)
	}
	return nil, errors.New("not implemented")
}

func (f *fakeAirportServiceForInteractor) ActivateAirportTx(ctx context.Context, tx output.Tx, id string) error {
	if f.activateTxFn != nil {
		return f.activateTxFn(ctx, tx, id)
	}
	return nil
}

func (f *fakeAirportServiceForInteractor) DeactivateAirportTx(ctx context.Context, tx output.Tx, id string) error {
	if f.deactivateTxFn != nil {
		return f.deactivateTxFn(ctx, tx, id)
	}
	return nil
}

func TestAirportInteractor_GetAirportByID(t *testing.T) {
	ctx := context.Background()

	t.Run("success => returns airport", func(t *testing.T) {
		expectedAirport := &domain.Airport{
			ID:       "airport-123",
			Name:     "El Dorado International",
			IATACode: "BOG",
			Status:   true,
		}
		svc := &fakeAirportServiceForInteractor{
			getByIDFn: func(context.Context, string) (*domain.Airport, error) {
				return expectedAirport, nil
			},
		}
		interactor := NewAirportInteractor(svc)

		result, err := interactor.GetAirportByID(ctx, "airport-123")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result == nil {
			t.Fatal("expected airport, got nil")
		}
		if result.ID != expectedAirport.ID {
			t.Errorf("expected ID %s, got %s", expectedAirport.ID, result.ID)
		}
		if result.Name != expectedAirport.Name {
			t.Errorf("expected name %s, got %s", expectedAirport.Name, result.Name)
		}
	})

	t.Run("not found => returns ErrAirportNotFound", func(t *testing.T) {
		svc := &fakeAirportServiceForInteractor{
			getByIDFn: func(context.Context, string) (*domain.Airport, error) {
				return nil, domain.ErrAirportNotFound
			},
		}
		interactor := NewAirportInteractor(svc)

		_, err := interactor.GetAirportByID(ctx, "non-existent")
		if !errors.Is(err, domain.ErrAirportNotFound) {
			t.Fatalf("expected %v, got %v", domain.ErrAirportNotFound, err)
		}
	})

	t.Run("service error => propagate error", func(t *testing.T) {
		serviceErr := errors.New("service unavailable")
		svc := &fakeAirportServiceForInteractor{
			getByIDFn: func(context.Context, string) (*domain.Airport, error) {
				return nil, serviceErr
			},
		}
		interactor := NewAirportInteractor(svc)

		_, err := interactor.GetAirportByID(ctx, "airport-123")
		if !errors.Is(err, serviceErr) {
			t.Fatalf("expected %v, got %v", serviceErr, err)
		}
	})
}

func TestAirportInteractor_DeactivateAirport(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		deactivateCalled := false
		svc := &fakeAirportServiceForInteractor{
			getByIDFn: func(context.Context, string) (*domain.Airport, error) {
				return &domain.Airport{ID: "airport-123", Status: true}, nil
			},
			deactivateTxFn: func(context.Context, output.Tx, string) error {
				deactivateCalled = true
				return nil
			},
		}
		interactor := NewAirportInteractor(svc)

		err := interactor.DeactivateAirport(ctx, "airport-123")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !deactivateCalled {
			t.Fatal("expected DeactivateAirportTx to be called")
		}
	})

	t.Run("airport not found => returns error", func(t *testing.T) {
		svc := &fakeAirportServiceForInteractor{
			getByIDFn: func(context.Context, string) (*domain.Airport, error) {
				return nil, domain.ErrAirportNotFound
			},
		}
		interactor := NewAirportInteractor(svc)

		err := interactor.DeactivateAirport(ctx, "non-existent")
		if !errors.Is(err, domain.ErrAirportNotFound) {
			t.Fatalf("expected %v, got %v", domain.ErrAirportNotFound, err)
		}
	})

	t.Run("deactivate fails => returns error", func(t *testing.T) {
		deactivateErr := errors.New("failed to deactivate")
		svc := &fakeAirportServiceForInteractor{
			getByIDFn: func(context.Context, string) (*domain.Airport, error) {
				return &domain.Airport{ID: "airport-123"}, nil
			},
			deactivateTxFn: func(context.Context, output.Tx, string) error {
				return deactivateErr
			},
		}
		interactor := NewAirportInteractor(svc)

		err := interactor.DeactivateAirport(ctx, "airport-123")
		if !errors.Is(err, deactivateErr) {
			t.Fatalf("expected %v, got %v", deactivateErr, err)
		}
	})

	t.Run("begin tx error", func(t *testing.T) {
		svc := &fakeAirportServiceForInteractor{
			getByIDFn: func(context.Context, string) (*domain.Airport, error) {
				return &domain.Airport{ID: "airport-123"}, nil
			},
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return nil, errors.New("tx start failed")
			},
		}
		interactor := NewAirportInteractor(svc)

		err := interactor.DeactivateAirport(ctx, "airport-123")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("commit error", func(t *testing.T) {
		svc := &fakeAirportServiceForInteractor{
			getByIDFn: func(context.Context, string) (*domain.Airport, error) {
				return &domain.Airport{ID: "airport-123"}, nil
			},
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return &fakeTx{commitFn: func() error { return errors.New("commit failed") }}, nil
			},
			deactivateTxFn: func(context.Context, output.Tx, string) error {
				return nil
			},
		}
		interactor := NewAirportInteractor(svc)

		err := interactor.DeactivateAirport(ctx, "airport-123")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestAirportInteractor_ListAirports(t *testing.T) {
	ctx := context.Background()

	t.Run("success => returns list of airports", func(t *testing.T) {
		expectedAirports := []domain.Airport{
			{ID: "airport-1", Name: "El Dorado", IATACode: "BOG"},
			{ID: "airport-2", Name: "Jose Maria Cordova", IATACode: "MDE"},
		}
		svc := &fakeAirportServiceForInteractor{
			listAirportsFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.Airport, error) {
				return expectedAirports, nil
			},
		}
		interactor := NewAirportInteractor(svc)

		result, err := interactor.ListAirports(ctx, nil)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(result) != 2 {
			t.Errorf("expected 2 airports, got %d", len(result))
		}
	})

	t.Run("success with filters", func(t *testing.T) {
		expectedFilters := map[string]interface{}{"status": true}
		svc := &fakeAirportServiceForInteractor{
			listAirportsFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.Airport, error) {
				if filters["status"] != true {
					t.Errorf("expected status filter true, got %v", filters["status"])
				}
				return []domain.Airport{}, nil
			},
		}
		interactor := NewAirportInteractor(svc)

		_, err := interactor.ListAirports(ctx, expectedFilters)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("empty list => returns empty slice", func(t *testing.T) {
		svc := &fakeAirportServiceForInteractor{
			listAirportsFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.Airport, error) {
				return []domain.Airport{}, nil
			},
		}
		interactor := NewAirportInteractor(svc)

		result, err := interactor.ListAirports(ctx, nil)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(result) != 0 {
			t.Errorf("expected 0 airports, got %d", len(result))
		}
	})

	t.Run("service error => propagate error", func(t *testing.T) {
		serviceErr := errors.New("database error")
		svc := &fakeAirportServiceForInteractor{
			listAirportsFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.Airport, error) {
				return nil, serviceErr
			},
		}
		interactor := NewAirportInteractor(svc)

		_, err := interactor.ListAirports(ctx, nil)
		if !errors.Is(err, serviceErr) {
			t.Fatalf("expected %v, got %v", serviceErr, err)
		}
	})
}

func TestAirportInteractor_GetAirportsByType(t *testing.T) {
	ctx := context.Background()

	t.Run("success => returns airports of type", func(t *testing.T) {
		expectedAirports := []domain.Airport{
			{ID: "airport-1", Name: "El Dorado", AirportType: "INTERNACIONAL"},
			{ID: "airport-2", Name: "Cali Alfonso", AirportType: "INTERNACIONAL"},
		}
		svc := &fakeAirportServiceForInteractor{
			getByTypeFn: func(ctx context.Context, airportType string) ([]domain.Airport, error) {
				if airportType != "INTERNACIONAL" {
					t.Errorf("expected type INTERNACIONAL, got %s", airportType)
				}
				return expectedAirports, nil
			},
		}
		interactor := NewAirportInteractor(svc)

		result, err := interactor.GetAirportsByType(ctx, "INTERNACIONAL")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(result) != 2 {
			t.Errorf("expected 2 airports, got %d", len(result))
		}
	})

	t.Run("empty => returns empty slice", func(t *testing.T) {
		svc := &fakeAirportServiceForInteractor{
			getByTypeFn: func(ctx context.Context, airportType string) ([]domain.Airport, error) {
				return []domain.Airport{}, nil
			},
		}
		interactor := NewAirportInteractor(svc)

		result, err := interactor.GetAirportsByType(ctx, "NACIONAL")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(result) != 0 {
			t.Errorf("expected 0 airports, got %d", len(result))
		}
	})

	t.Run("service error => propagate error", func(t *testing.T) {
		serviceErr := errors.New("database error")
		svc := &fakeAirportServiceForInteractor{
			getByTypeFn: func(ctx context.Context, airportType string) ([]domain.Airport, error) {
				return nil, serviceErr
			},
		}
		interactor := NewAirportInteractor(svc)

		_, err := interactor.GetAirportsByType(ctx, "INTERNACIONAL")
		if !errors.Is(err, serviceErr) {
			t.Fatalf("expected %v, got %v", serviceErr, err)
		}
	})
}

func TestAirportInteractor_ActivateAirport(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		activateCalled := false
		svc := &fakeAirportServiceForInteractor{
			getByIDFn: func(context.Context, string) (*domain.Airport, error) {
				return &domain.Airport{ID: "airport-123", Status: false}, nil
			},
			activateTxFn: func(context.Context, output.Tx, string) error {
				activateCalled = true
				return nil
			},
		}
		interactor := NewAirportInteractor(svc)

		err := interactor.ActivateAirport(ctx, "airport-123")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !activateCalled {
			t.Fatal("expected ActivateAirportTx to be called")
		}
	})

	t.Run("airport not found => returns error", func(t *testing.T) {
		svc := &fakeAirportServiceForInteractor{
			getByIDFn: func(context.Context, string) (*domain.Airport, error) {
				return nil, domain.ErrAirportNotFound
			},
		}
		interactor := NewAirportInteractor(svc)

		err := interactor.ActivateAirport(ctx, "non-existent")
		if !errors.Is(err, domain.ErrAirportNotFound) {
			t.Fatalf("expected %v, got %v", domain.ErrAirportNotFound, err)
		}
	})

	t.Run("activate fails => returns error", func(t *testing.T) {
		activateErr := errors.New("failed to activate")
		svc := &fakeAirportServiceForInteractor{
			getByIDFn: func(context.Context, string) (*domain.Airport, error) {
				return &domain.Airport{ID: "airport-123"}, nil
			},
			activateTxFn: func(context.Context, output.Tx, string) error {
				return activateErr
			},
		}
		interactor := NewAirportInteractor(svc)

		err := interactor.ActivateAirport(ctx, "airport-123")
		if !errors.Is(err, activateErr) {
			t.Fatalf("expected %v, got %v", activateErr, err)
		}
	})

	t.Run("begin tx error", func(t *testing.T) {
		svc := &fakeAirportServiceForInteractor{
			getByIDFn: func(context.Context, string) (*domain.Airport, error) {
				return &domain.Airport{ID: "airport-123"}, nil
			},
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return nil, errors.New("tx start failed")
			},
		}
		interactor := NewAirportInteractor(svc)

		err := interactor.ActivateAirport(ctx, "airport-123")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("commit error", func(t *testing.T) {
		svc := &fakeAirportServiceForInteractor{
			getByIDFn: func(context.Context, string) (*domain.Airport, error) {
				return &domain.Airport{ID: "airport-123"}, nil
			},
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return &fakeTx{commitFn: func() error { return errors.New("commit failed") }}, nil
			},
			activateTxFn: func(context.Context, output.Tx, string) error {
				return nil
			},
		}
		interactor := NewAirportInteractor(svc)

		err := interactor.ActivateAirport(ctx, "airport-123")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
