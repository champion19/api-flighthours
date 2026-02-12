package services

import (
	"context"
	"errors"
	"testing"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
)

func pilotRolePtr(r domain.PilotRole) *domain.PilotRole { return &r }

// mock daily logbook detail repository
type mockDailyLogbookDetailRepo struct {
	getByIDFn        func(ctx context.Context, id string) (*domain.DailyLogbookDetail, error)
	listByLogbookFn  func(ctx context.Context, logbookID string) ([]domain.DailyLogbookDetail, error)
	listByEmployeeFn func(ctx context.Context, employeeID string) ([]domain.DailyLogbookDetail, error)
	saveFn           func(ctx context.Context, tx output.Tx, detail domain.DailyLogbookDetail) error
	updateFn         func(ctx context.Context, tx output.Tx, detail domain.DailyLogbookDetail) error
	beginTxFn        func(ctx context.Context) (output.Tx, error)
}

func (m *mockDailyLogbookDetailRepo) GetDailyLogbookDetailByID(ctx context.Context, id string) (*domain.DailyLogbookDetail, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *mockDailyLogbookDetailRepo) ListDailyLogbookDetailsByLogbook(ctx context.Context, logbookID string) ([]domain.DailyLogbookDetail, error) {
	if m.listByLogbookFn != nil {
		return m.listByLogbookFn(ctx, logbookID)
	}
	return nil, nil
}

func (m *mockDailyLogbookDetailRepo) ListDailyLogbookDetailsByEmployee(ctx context.Context, employeeID string) ([]domain.DailyLogbookDetail, error) {
	if m.listByEmployeeFn != nil {
		return m.listByEmployeeFn(ctx, employeeID)
	}
	return nil, nil
}

func (m *mockDailyLogbookDetailRepo) SaveDailyLogbookDetail(ctx context.Context, tx output.Tx, detail domain.DailyLogbookDetail) error {
	if m.saveFn != nil {
		return m.saveFn(ctx, tx, detail)
	}
	return nil
}

func (m *mockDailyLogbookDetailRepo) UpdateDailyLogbookDetail(ctx context.Context, tx output.Tx, detail domain.DailyLogbookDetail) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, tx, detail)
	}
	return nil
}

func (m *mockDailyLogbookDetailRepo) BeginTx(ctx context.Context) (output.Tx, error) {
	if m.beginTxFn != nil {
		return m.beginTxFn(ctx)
	}
	return &mockTx{}, nil
}

func TestNewDailyLogbookDetailService(t *testing.T) {
	repo := &mockDailyLogbookDetailRepo{}
	svc := NewDailyLogbookDetailService(repo)
	if svc == nil {
		t.Error("expected non-nil service")
	}
}

func TestDailyLogbookDetailService_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expected := &domain.DailyLogbookDetail{ID: "d-1", FlightNumber: "AV123"}
		repo := &mockDailyLogbookDetailRepo{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbookDetail, error) {
				return expected, nil
			},
		}
		svc := NewDailyLogbookDetailService(repo)
		result, err := svc.GetDailyLogbookDetailByID(context.Background(), "d-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.ID != "d-1" {
			t.Errorf("expected ID 'd-1', got %q", result.ID)
		}
	})

	t.Run("error", func(t *testing.T) {
		repo := &mockDailyLogbookDetailRepo{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbookDetail, error) {
				return nil, errors.New("not found")
			},
		}
		svc := NewDailyLogbookDetailService(repo)
		_, err := svc.GetDailyLogbookDetailByID(context.Background(), "d-1")
		if err == nil {
			t.Error("expected error")
		}
	})
}

func TestDailyLogbookDetailService_ListByLogbook(t *testing.T) {
	expected := []domain.DailyLogbookDetail{
		{ID: "d-1", DailyLogbookID: "lb-1"},
		{ID: "d-2", DailyLogbookID: "lb-1"},
	}
	repo := &mockDailyLogbookDetailRepo{
		listByLogbookFn: func(ctx context.Context, logbookID string) ([]domain.DailyLogbookDetail, error) {
			return expected, nil
		},
	}
	svc := NewDailyLogbookDetailService(repo)
	result, err := svc.ListDailyLogbookDetailsByLogbook(context.Background(), "lb-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 details, got %d", len(result))
	}
}

func TestDailyLogbookDetailService_ListByEmployee(t *testing.T) {
	expected := []domain.DailyLogbookDetail{
		{ID: "d-1"},
	}
	repo := &mockDailyLogbookDetailRepo{
		listByEmployeeFn: func(ctx context.Context, employeeID string) ([]domain.DailyLogbookDetail, error) {
			return expected, nil
		},
	}
	svc := NewDailyLogbookDetailService(repo)
	result, err := svc.ListDailyLogbookDetailsByEmployee(context.Background(), "emp-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 detail, got %d", len(result))
	}
}

func TestDailyLogbookDetailService_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &mockDailyLogbookDetailRepo{
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return &mockTx{}, nil
			},
			saveFn: func(ctx context.Context, tx output.Tx, detail domain.DailyLogbookDetail) error {
				return nil
			},
		}
		svc := NewDailyLogbookDetailService(repo)
		err := svc.CreateDailyLogbookDetail(context.Background(), domain.DailyLogbookDetail{
			ID: "d-1", DailyLogbookID: "lb-1", FlightNumber: "AV123",
			OutTime: strPtr("08:00"), TakeoffTime: strPtr("08:15"), LandingTime: strPtr("09:30"), InTime: strPtr("09:45"),
			PilotRole: pilotRolePtr(domain.PilotRolePF),
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("begin tx error", func(t *testing.T) {
		repo := &mockDailyLogbookDetailRepo{
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return nil, errors.New("tx error")
			},
		}
		svc := NewDailyLogbookDetailService(repo)
		err := svc.CreateDailyLogbookDetail(context.Background(), domain.DailyLogbookDetail{ID: "d-1"})
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("save error rolls back", func(t *testing.T) {
		rollbackCalled := false
		repo := &mockDailyLogbookDetailRepo{
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return &mockTx{
					rollbackFn: func() error { rollbackCalled = true; return nil },
				}, nil
			},
			saveFn: func(ctx context.Context, tx output.Tx, detail domain.DailyLogbookDetail) error {
				return errors.New("save error")
			},
		}
		svc := NewDailyLogbookDetailService(repo)
		err := svc.CreateDailyLogbookDetail(context.Background(), domain.DailyLogbookDetail{ID: "d-1"})
		if err == nil {
			t.Error("expected error")
		}
		if !rollbackCalled {
			t.Error("expected rollback")
		}
	})
}

func TestDailyLogbookDetailService_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &mockDailyLogbookDetailRepo{
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return &mockTx{}, nil
			},
			updateFn: func(ctx context.Context, tx output.Tx, detail domain.DailyLogbookDetail) error {
				return nil
			},
		}
		svc := NewDailyLogbookDetailService(repo)
		err := svc.UpdateDailyLogbookDetail(context.Background(), domain.DailyLogbookDetail{ID: "d-1"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("begin tx error", func(t *testing.T) {
		repo := &mockDailyLogbookDetailRepo{
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return nil, errors.New("tx error")
			},
		}
		svc := NewDailyLogbookDetailService(repo)
		err := svc.UpdateDailyLogbookDetail(context.Background(), domain.DailyLogbookDetail{ID: "d-1"})
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("update error rolls back", func(t *testing.T) {
		rollbackCalled := false
		repo := &mockDailyLogbookDetailRepo{
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return &mockTx{
					rollbackFn: func() error { rollbackCalled = true; return nil },
				}, nil
			},
			updateFn: func(ctx context.Context, tx output.Tx, detail domain.DailyLogbookDetail) error {
				return errors.New("update error")
			},
		}
		svc := NewDailyLogbookDetailService(repo)
		err := svc.UpdateDailyLogbookDetail(context.Background(), domain.DailyLogbookDetail{ID: "d-1"})
		if err == nil {
			t.Error("expected error")
		}
		if !rollbackCalled {
			t.Error("expected rollback")
		}
	})
}

func TestDailyLogbookDetailService_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &mockDailyLogbookDetailRepo{
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return &mockTx{}, nil
			},
		}
		svc := NewDailyLogbookDetailService(repo)
		err := svc.DeleteDailyLogbookDetail(context.Background(), "d-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("begin tx error", func(t *testing.T) {
		repo := &mockDailyLogbookDetailRepo{
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return nil, errors.New("tx error")
			},
		}
		svc := NewDailyLogbookDetailService(repo)
		err := svc.DeleteDailyLogbookDetail(context.Background(), "d-1")
		if err == nil {
			t.Error("expected error")
		}
	})
}

func TestDailyLogbookDetailService_ValidateTimeSequence(t *testing.T) {
	svc := NewDailyLogbookDetailService(&mockDailyLogbookDetailRepo{})

	tests := []struct {
		name    string
		out     string
		takeoff string
		landing string
		in      string
		wantErr bool
	}{
		{"valid HH:MM", "08:00", "08:15", "09:30", "09:45", false},
		{"valid HH:MM:SS", "08:00:00", "08:15:00", "09:30:00", "09:45:00", false},
		{"invalid out_time format", "bad", "08:15", "09:30", "09:45", true},
		{"invalid takeoff_time format", "08:00", "bad", "09:30", "09:45", true},
		{"invalid landing_time format", "08:00", "08:15", "bad", "09:45", true},
		{"invalid in_time format", "08:00", "08:15", "09:30", "bad", true},
		{"out not before takeoff", "09:00", "08:15", "09:30", "09:45", true},
		{"takeoff not before landing", "08:00", "10:00", "09:30", "09:45", true},
		{"landing not before in", "08:00", "08:15", "10:00", "09:45", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.ValidateTimeSequence(tt.out, tt.takeoff, tt.landing, tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTimeSequence() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDailyLogbookDetailService_BeginTx(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &mockDailyLogbookDetailRepo{
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return &mockTx{}, nil
			},
		}
		svc := NewDailyLogbookDetailService(repo)
		tx, err := svc.BeginTx(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tx == nil {
			t.Error("expected non-nil tx")
		}
	})

	t.Run("error", func(t *testing.T) {
		repo := &mockDailyLogbookDetailRepo{
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return nil, errors.New("tx error")
			},
		}
		svc := NewDailyLogbookDetailService(repo)
		_, err := svc.BeginTx(context.Background())
		if err == nil {
			t.Error("expected error")
		}
	})
}

func TestDailyLogbookDetailService_CreateTx(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &mockDailyLogbookDetailRepo{
			saveFn: func(ctx context.Context, tx output.Tx, detail domain.DailyLogbookDetail) error {
				return nil
			},
		}
		svc := NewDailyLogbookDetailService(repo)
		err := svc.CreateDailyLogbookDetailTx(context.Background(), &mockTx{}, domain.DailyLogbookDetail{ID: "d-1"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		repo := &mockDailyLogbookDetailRepo{
			saveFn: func(ctx context.Context, tx output.Tx, detail domain.DailyLogbookDetail) error {
				return errors.New("save failed")
			},
		}
		svc := NewDailyLogbookDetailService(repo)
		err := svc.CreateDailyLogbookDetailTx(context.Background(), &mockTx{}, domain.DailyLogbookDetail{})
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestDailyLogbookDetailService_UpdateTx(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &mockDailyLogbookDetailRepo{
			updateFn: func(ctx context.Context, tx output.Tx, detail domain.DailyLogbookDetail) error {
				return nil
			},
		}
		svc := NewDailyLogbookDetailService(repo)
		err := svc.UpdateDailyLogbookDetailTx(context.Background(), &mockTx{}, domain.DailyLogbookDetail{ID: "d-1"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		repo := &mockDailyLogbookDetailRepo{
			updateFn: func(ctx context.Context, tx output.Tx, detail domain.DailyLogbookDetail) error {
				return errors.New("update failed")
			},
		}
		svc := NewDailyLogbookDetailService(repo)
		err := svc.UpdateDailyLogbookDetailTx(context.Background(), &mockTx{}, domain.DailyLogbookDetail{})
		if err == nil {
			t.Fatal("expected error")
		}
	})
}
