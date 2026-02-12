package interactor

import (
	"context"
	"errors"
	"testing"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/input"
	"github.com/champion19/api-flighthours/core/ports/output"
)

// fakeAircraftModelService implements input.AircraftModelService for testing
type fakeAircraftModelService struct {
	beginTxFn      func(ctx context.Context) (output.Tx, error)
	getByIDFn      func(ctx context.Context, id string) (*domain.AircraftModel, error)
	listFn         func(ctx context.Context, filters map[string]interface{}) ([]domain.AircraftModel, error)
	getByFamilyFn  func(ctx context.Context, family string) ([]domain.AircraftModel, error)
	activateFn     func(ctx context.Context, id string) error
	deactivateFn   func(ctx context.Context, id string) error
	activateTxFn   func(ctx context.Context, tx output.Tx, id string) error
	deactivateTxFn func(ctx context.Context, tx output.Tx, id string) error
}

var _ input.AircraftModelService = (*fakeAircraftModelService)(nil)

func (f *fakeAircraftModelService) BeginTx(ctx context.Context) (output.Tx, error) {
	if f.beginTxFn != nil {
		return f.beginTxFn(ctx)
	}
	return &fakeTx{}, nil
}

func (f *fakeAircraftModelService) GetAircraftModelByID(ctx context.Context, id string) (*domain.AircraftModel, error) {
	if f.getByIDFn != nil {
		return f.getByIDFn(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (f *fakeAircraftModelService) ListAircraftModels(ctx context.Context, filters map[string]interface{}) ([]domain.AircraftModel, error) {
	if f.listFn != nil {
		return f.listFn(ctx, filters)
	}
	return nil, errors.New("not implemented")
}

func (f *fakeAircraftModelService) GetAircraftModelsByFamily(ctx context.Context, family string) ([]domain.AircraftModel, error) {
	if f.getByFamilyFn != nil {
		return f.getByFamilyFn(ctx, family)
	}
	return nil, errors.New("not implemented")
}

func (f *fakeAircraftModelService) ActivateAircraftModel(ctx context.Context, id string) error {
	if f.activateFn != nil {
		return f.activateFn(ctx, id)
	}
	return nil
}

func (f *fakeAircraftModelService) DeactivateAircraftModel(ctx context.Context, id string) error {
	if f.deactivateFn != nil {
		return f.deactivateFn(ctx, id)
	}
	return nil
}

func (f *fakeAircraftModelService) ActivateAircraftModelTx(ctx context.Context, tx output.Tx, id string) error {
	if f.activateTxFn != nil {
		return f.activateTxFn(ctx, tx, id)
	}
	return nil
}

func (f *fakeAircraftModelService) DeactivateAircraftModelTx(ctx context.Context, tx output.Tx, id string) error {
	if f.deactivateTxFn != nil {
		return f.deactivateTxFn(ctx, tx, id)
	}
	return nil
}

func TestNewAircraftModelInteractor(t *testing.T) {
	svc := &fakeAircraftModelService{}
	interactor := NewAircraftModelInteractor(svc)
	if interactor == nil {
		t.Error("expected non-nil AircraftModelInteractor")
	}
}

func TestAircraftModelInteractor_GetAircraftModelByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expected := &domain.AircraftModel{
			ID:        "model-123",
			ModelName: "Boeing 737-800",
			Family:    "737",
			Status:    true,
		}
		svc := &fakeAircraftModelService{
			getByIDFn: func(ctx context.Context, id string) (*domain.AircraftModel, error) {
				return expected, nil
			},
		}
		i := NewAircraftModelInteractor(svc)

		result, err := i.GetAircraftModelByID(context.Background(), "model-123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.ID != expected.ID {
			t.Errorf("expected ID %q, got %q", expected.ID, result.ID)
		}
		if result.ModelName != expected.ModelName {
			t.Errorf("expected ModelName %q, got %q", expected.ModelName, result.ModelName)
		}
	})

	t.Run("not found", func(t *testing.T) {
		svc := &fakeAircraftModelService{
			getByIDFn: func(ctx context.Context, id string) (*domain.AircraftModel, error) {
				return nil, domain.ErrAircraftModelNotFound
			},
		}
		i := NewAircraftModelInteractor(svc)

		_, err := i.GetAircraftModelByID(context.Background(), "nonexistent")
		if err != domain.ErrAircraftModelNotFound {
			t.Errorf("expected ErrAircraftModelNotFound, got %v", err)
		}
	})

	t.Run("service error", func(t *testing.T) {
		svc := &fakeAircraftModelService{
			getByIDFn: func(ctx context.Context, id string) (*domain.AircraftModel, error) {
				return nil, errors.New("database error")
			},
		}
		i := NewAircraftModelInteractor(svc)

		_, err := i.GetAircraftModelByID(context.Background(), "model-123")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestAircraftModelInteractor_ListAircraftModels(t *testing.T) {
	t.Run("success without filters", func(t *testing.T) {
		expected := []domain.AircraftModel{
			{ID: "model-1", ModelName: "Boeing 737-800", Family: "737"},
			{ID: "model-2", ModelName: "Airbus A320", Family: "A320"},
		}
		svc := &fakeAircraftModelService{
			listFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.AircraftModel, error) {
				return expected, nil
			},
		}
		i := NewAircraftModelInteractor(svc)

		result, err := i.ListAircraftModels(context.Background(), nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 2 {
			t.Errorf("expected 2 models, got %d", len(result))
		}
	})

	t.Run("success with filters", func(t *testing.T) {
		svc := &fakeAircraftModelService{
			listFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.AircraftModel, error) {
				if _, ok := filters["engine_type"]; !ok {
					t.Error("expected 'engine_type' filter to be present")
				}
				return []domain.AircraftModel{{ID: "model-1"}}, nil
			},
		}
		i := NewAircraftModelInteractor(svc)

		filters := map[string]interface{}{"engine_type": "JET"}
		result, err := i.ListAircraftModels(context.Background(), filters)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 1 {
			t.Errorf("expected 1 model, got %d", len(result))
		}
	})

	t.Run("service error", func(t *testing.T) {
		svc := &fakeAircraftModelService{
			listFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.AircraftModel, error) {
				return nil, errors.New("database error")
			},
		}
		i := NewAircraftModelInteractor(svc)

		_, err := i.ListAircraftModels(context.Background(), nil)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestAircraftModelInteractor_GetAircraftModelsByFamily(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expected := []domain.AircraftModel{
			{ID: "model-1", ModelName: "737-700", Family: "737"},
			{ID: "model-2", ModelName: "737-800", Family: "737"},
			{ID: "model-3", ModelName: "737-900", Family: "737"},
		}
		svc := &fakeAircraftModelService{
			getByFamilyFn: func(ctx context.Context, family string) ([]domain.AircraftModel, error) {
				if family != "737" {
					t.Errorf("expected family '737', got %q", family)
				}
				return expected, nil
			},
		}
		i := NewAircraftModelInteractor(svc)

		result, err := i.GetAircraftModelsByFamily(context.Background(), "737")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 3 {
			t.Errorf("expected 3 models, got %d", len(result))
		}
	})

	t.Run("empty family", func(t *testing.T) {
		svc := &fakeAircraftModelService{
			getByFamilyFn: func(ctx context.Context, family string) ([]domain.AircraftModel, error) {
				return []domain.AircraftModel{}, nil
			},
		}
		i := NewAircraftModelInteractor(svc)

		result, err := i.GetAircraftModelsByFamily(context.Background(), "nonexistent")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 0 {
			t.Errorf("expected 0 models, got %d", len(result))
		}
	})

	t.Run("service error", func(t *testing.T) {
		svc := &fakeAircraftModelService{
			getByFamilyFn: func(ctx context.Context, family string) ([]domain.AircraftModel, error) {
				return nil, errors.New("database error")
			},
		}
		i := NewAircraftModelInteractor(svc)

		_, err := i.GetAircraftModelsByFamily(context.Background(), "737")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestAircraftModelInteractor_ActivateAircraftModel(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &fakeAircraftModelService{
			getByIDFn: func(ctx context.Context, id string) (*domain.AircraftModel, error) {
				return &domain.AircraftModel{ID: id, ModelName: "737-800"}, nil
			},
			activateTxFn: func(ctx context.Context, tx output.Tx, id string) error {
				return nil
			},
		}
		i := NewAircraftModelInteractor(svc)

		err := i.ActivateAircraftModel(context.Background(), "model-123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		svc := &fakeAircraftModelService{
			getByIDFn: func(ctx context.Context, id string) (*domain.AircraftModel, error) {
				return nil, domain.ErrAircraftModelNotFound
			},
		}
		i := NewAircraftModelInteractor(svc)

		err := i.ActivateAircraftModel(context.Background(), "nonexistent")
		if err != domain.ErrAircraftModelNotFound {
			t.Errorf("expected ErrAircraftModelNotFound, got %v", err)
		}
	})

	t.Run("activation service error", func(t *testing.T) {
		svc := &fakeAircraftModelService{
			getByIDFn: func(ctx context.Context, id string) (*domain.AircraftModel, error) {
				return &domain.AircraftModel{ID: id}, nil
			},
			activateTxFn: func(ctx context.Context, tx output.Tx, id string) error {
				return errors.New("activation failed")
			},
		}
		i := NewAircraftModelInteractor(svc)

		err := i.ActivateAircraftModel(context.Background(), "model-123")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("begin tx error", func(t *testing.T) {
		svc := &fakeAircraftModelService{
			getByIDFn: func(ctx context.Context, id string) (*domain.AircraftModel, error) {
				return &domain.AircraftModel{ID: id}, nil
			},
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return nil, errors.New("tx error")
			},
		}
		i := NewAircraftModelInteractor(svc)

		err := i.ActivateAircraftModel(context.Background(), "model-123")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("commit error", func(t *testing.T) {
		svc := &fakeAircraftModelService{
			getByIDFn: func(ctx context.Context, id string) (*domain.AircraftModel, error) {
				return &domain.AircraftModel{ID: id}, nil
			},
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return &fakeTx{commitFn: func() error { return errors.New("commit error") }}, nil
			},
			activateTxFn: func(ctx context.Context, tx output.Tx, id string) error {
				return nil
			},
		}
		i := NewAircraftModelInteractor(svc)

		err := i.ActivateAircraftModel(context.Background(), "model-123")
		if err == nil {
			t.Error("expected commit error, got nil")
		}
	})
}

func TestAircraftModelInteractor_DeactivateAircraftModel(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &fakeAircraftModelService{
			getByIDFn: func(ctx context.Context, id string) (*domain.AircraftModel, error) {
				return &domain.AircraftModel{ID: id, ModelName: "737-800"}, nil
			},
			deactivateTxFn: func(ctx context.Context, tx output.Tx, id string) error {
				return nil
			},
		}
		i := NewAircraftModelInteractor(svc)

		err := i.DeactivateAircraftModel(context.Background(), "model-123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		svc := &fakeAircraftModelService{
			getByIDFn: func(ctx context.Context, id string) (*domain.AircraftModel, error) {
				return nil, domain.ErrAircraftModelNotFound
			},
		}
		i := NewAircraftModelInteractor(svc)

		err := i.DeactivateAircraftModel(context.Background(), "nonexistent")
		if err != domain.ErrAircraftModelNotFound {
			t.Errorf("expected ErrAircraftModelNotFound, got %v", err)
		}
	})

	t.Run("deactivation service error", func(t *testing.T) {
		svc := &fakeAircraftModelService{
			getByIDFn: func(ctx context.Context, id string) (*domain.AircraftModel, error) {
				return &domain.AircraftModel{ID: id}, nil
			},
			deactivateTxFn: func(ctx context.Context, tx output.Tx, id string) error {
				return errors.New("deactivation failed")
			},
		}
		i := NewAircraftModelInteractor(svc)

		err := i.DeactivateAircraftModel(context.Background(), "model-123")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("begin tx error", func(t *testing.T) {
		svc := &fakeAircraftModelService{
			getByIDFn: func(ctx context.Context, id string) (*domain.AircraftModel, error) {
				return &domain.AircraftModel{ID: id}, nil
			},
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return nil, errors.New("tx error")
			},
		}
		i := NewAircraftModelInteractor(svc)

		err := i.DeactivateAircraftModel(context.Background(), "model-123")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("commit error", func(t *testing.T) {
		svc := &fakeAircraftModelService{
			getByIDFn: func(ctx context.Context, id string) (*domain.AircraftModel, error) {
				return &domain.AircraftModel{ID: id}, nil
			},
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return &fakeTx{commitFn: func() error { return errors.New("commit error") }}, nil
			},
			deactivateTxFn: func(ctx context.Context, tx output.Tx, id string) error {
				return nil
			},
		}
		i := NewAircraftModelInteractor(svc)

		err := i.DeactivateAircraftModel(context.Background(), "model-123")
		if err == nil {
			t.Error("expected commit error, got nil")
		}
	})
}
