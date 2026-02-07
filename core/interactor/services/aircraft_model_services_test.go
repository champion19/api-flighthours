package services

import (
	"context"
	"errors"
	"testing"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/logger"
)

// fakeAircraftModelRepo implements output.AircraftModelRepository for testing
type fakeAircraftModelRepo struct {
	beginTxFn      func(ctx context.Context) (output.Tx, error)
	getByIDFn      func(ctx context.Context, id string) (*domain.AircraftModel, error)
	listFn         func(ctx context.Context, filters map[string]interface{}) ([]domain.AircraftModel, error)
	getByFamilyFn  func(ctx context.Context, family string) ([]domain.AircraftModel, error)
	updateStatusFn func(ctx context.Context, tx output.Tx, id string, status bool) error
}

var _ output.AircraftModelRepository = (*fakeAircraftModelRepo)(nil)

func (f *fakeAircraftModelRepo) BeginTx(ctx context.Context) (output.Tx, error) {
	if f.beginTxFn != nil {
		return f.beginTxFn(ctx)
	}
	return &fakeTxAM{}, nil
}

func (f *fakeAircraftModelRepo) GetAircraftModelByID(ctx context.Context, id string) (*domain.AircraftModel, error) {
	if f.getByIDFn != nil {
		return f.getByIDFn(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (f *fakeAircraftModelRepo) ListAircraftModels(ctx context.Context, filters map[string]interface{}) ([]domain.AircraftModel, error) {
	if f.listFn != nil {
		return f.listFn(ctx, filters)
	}
	return nil, errors.New("not implemented")
}

func (f *fakeAircraftModelRepo) GetAircraftModelsByFamily(ctx context.Context, family string) ([]domain.AircraftModel, error) {
	if f.getByFamilyFn != nil {
		return f.getByFamilyFn(ctx, family)
	}
	return nil, errors.New("not implemented")
}

func (f *fakeAircraftModelRepo) UpdateAircraftModelStatus(ctx context.Context, tx output.Tx, id string, status bool) error {
	if f.updateStatusFn != nil {
		return f.updateStatusFn(ctx, tx, id, status)
	}
	return nil
}

// fakeTxAM implements output.Tx for aircraft model tests
type fakeTxAM struct {
	commitErr   error
	rollbackErr error
}

func (f *fakeTxAM) Commit() error   { return f.commitErr }
func (f *fakeTxAM) Rollback() error { return f.rollbackErr }

func newTestAircraftModelLogger() logger.Logger {
	return logger.NewSlogLogger()
}

func TestNewAircraftModelService(t *testing.T) {
	repo := &fakeAircraftModelRepo{}
	log := newTestAircraftModelLogger()
	svc := NewAircraftModelService(repo, log)
	if svc == nil {
		t.Error("expected non-nil AircraftModelService")
	}
}

func TestAircraftModelService_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expected := &domain.AircraftModel{ID: "model-123", ModelName: "737-800"}
		repo := &fakeAircraftModelRepo{
			getByIDFn: func(ctx context.Context, id string) (*domain.AircraftModel, error) {
				return expected, nil
			},
		}
		svc := NewAircraftModelService(repo, newTestAircraftModelLogger())

		result, err := svc.GetAircraftModelByID(context.Background(), "model-123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.ID != expected.ID {
			t.Errorf("expected ID %q, got %q", expected.ID, result.ID)
		}
	})

	t.Run("error", func(t *testing.T) {
		repo := &fakeAircraftModelRepo{
			getByIDFn: func(ctx context.Context, id string) (*domain.AircraftModel, error) {
				return nil, errors.New("db error")
			},
		}
		svc := NewAircraftModelService(repo, newTestAircraftModelLogger())

		_, err := svc.GetAircraftModelByID(context.Background(), "model-123")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestAircraftModelService_ListAircraftModels(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expected := []domain.AircraftModel{
			{ID: "model-1", ModelName: "737-800"},
			{ID: "model-2", ModelName: "A320"},
		}
		repo := &fakeAircraftModelRepo{
			listFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.AircraftModel, error) {
				return expected, nil
			},
		}
		svc := NewAircraftModelService(repo, newTestAircraftModelLogger())

		result, err := svc.ListAircraftModels(context.Background(), nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 2 {
			t.Errorf("expected 2 models, got %d", len(result))
		}
	})

	t.Run("error", func(t *testing.T) {
		repo := &fakeAircraftModelRepo{
			listFn: func(ctx context.Context, filters map[string]interface{}) ([]domain.AircraftModel, error) {
				return nil, errors.New("db error")
			},
		}
		svc := NewAircraftModelService(repo, newTestAircraftModelLogger())

		_, err := svc.ListAircraftModels(context.Background(), nil)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestAircraftModelService_GetByFamily(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expected := []domain.AircraftModel{
			{ID: "m1", ModelName: "737-700", Family: "737"},
			{ID: "m2", ModelName: "737-800", Family: "737"},
		}
		repo := &fakeAircraftModelRepo{
			getByFamilyFn: func(ctx context.Context, family string) ([]domain.AircraftModel, error) {
				return expected, nil
			},
		}
		svc := NewAircraftModelService(repo, newTestAircraftModelLogger())

		result, err := svc.GetAircraftModelsByFamily(context.Background(), "737")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 2 {
			t.Errorf("expected 2 models, got %d", len(result))
		}
	})

	t.Run("error", func(t *testing.T) {
		repo := &fakeAircraftModelRepo{
			getByFamilyFn: func(ctx context.Context, family string) ([]domain.AircraftModel, error) {
				return nil, errors.New("db error")
			},
		}
		svc := NewAircraftModelService(repo, newTestAircraftModelLogger())

		_, err := svc.GetAircraftModelsByFamily(context.Background(), "737")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestAircraftModelService_UpdateStatus(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &fakeAircraftModelRepo{
			updateStatusFn: func(ctx context.Context, tx output.Tx, id string, status bool) error {
				return nil
			},
		}
		svc := NewAircraftModelService(repo, newTestAircraftModelLogger())

		err := svc.UpdateAircraftModelStatus(context.Background(), "model-123", true)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("begin tx error", func(t *testing.T) {
		repo := &fakeAircraftModelRepo{
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return nil, errors.New("tx error")
			},
		}
		svc := NewAircraftModelService(repo, newTestAircraftModelLogger())

		err := svc.UpdateAircraftModelStatus(context.Background(), "model-123", true)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("update error triggers rollback", func(t *testing.T) {
		repo := &fakeAircraftModelRepo{
			updateStatusFn: func(ctx context.Context, tx output.Tx, id string, status bool) error {
				return errors.New("update failed")
			},
		}
		svc := NewAircraftModelService(repo, newTestAircraftModelLogger())

		err := svc.UpdateAircraftModelStatus(context.Background(), "model-123", true)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestAircraftModelService_ActivateDeactivate(t *testing.T) {
	t.Run("activate calls UpdateAircraftModelStatus with true", func(t *testing.T) {
		var receivedStatus bool
		repo := &fakeAircraftModelRepo{
			updateStatusFn: func(ctx context.Context, tx output.Tx, id string, status bool) error {
				receivedStatus = status
				return nil
			},
		}
		svc := NewAircraftModelService(repo, newTestAircraftModelLogger())

		err := svc.ActivateAircraftModel(context.Background(), "model-123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !receivedStatus {
			t.Error("expected status=true for activation")
		}
	})

	t.Run("deactivate calls UpdateAircraftModelStatus with false", func(t *testing.T) {
		var receivedStatus bool
		receivedStatus = true // start with true
		repo := &fakeAircraftModelRepo{
			updateStatusFn: func(ctx context.Context, tx output.Tx, id string, status bool) error {
				receivedStatus = status
				return nil
			},
		}
		svc := NewAircraftModelService(repo, newTestAircraftModelLogger())

		err := svc.DeactivateAircraftModel(context.Background(), "model-123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if receivedStatus {
			t.Error("expected status=false for deactivation")
		}
	})
}
