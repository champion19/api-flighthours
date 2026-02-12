package interactor

import (
	"context"
	"errors"
	"testing"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
)

func strPtr(s string) *string { return &s }

// fakeDailyLogbookDetailService implements input.DailyLogbookDetailService
type fakeDailyLogbookDetailService struct {
	getByIDFn        func(ctx context.Context, id string) (*domain.DailyLogbookDetail, error)
	listByLogbookFn  func(ctx context.Context, logbookID string) ([]domain.DailyLogbookDetail, error)
	listByEmployeeFn func(ctx context.Context, employeeID string) ([]domain.DailyLogbookDetail, error)
	createFn         func(ctx context.Context, detail domain.DailyLogbookDetail) error
	updateFn         func(ctx context.Context, detail domain.DailyLogbookDetail) error
	validateTimeFn   func(outTime, takeoffTime, landingTime, inTime string) error
	createTxFn       func(ctx context.Context, tx output.Tx, detail domain.DailyLogbookDetail) error
	updateTxFn       func(ctx context.Context, tx output.Tx, detail domain.DailyLogbookDetail) error
}

func (f *fakeDailyLogbookDetailService) BeginTx(ctx context.Context) (output.Tx, error) {
	return &fakeTx{}, nil
}

func (f *fakeDailyLogbookDetailService) GetDailyLogbookDetailByID(ctx context.Context, id string) (*domain.DailyLogbookDetail, error) {
	if f.getByIDFn != nil {
		return f.getByIDFn(ctx, id)
	}
	return nil, nil
}

func (f *fakeDailyLogbookDetailService) ListDailyLogbookDetailsByLogbook(ctx context.Context, logbookID string) ([]domain.DailyLogbookDetail, error) {
	if f.listByLogbookFn != nil {
		return f.listByLogbookFn(ctx, logbookID)
	}
	return nil, nil
}

func (f *fakeDailyLogbookDetailService) ListDailyLogbookDetailsByEmployee(ctx context.Context, employeeID string) ([]domain.DailyLogbookDetail, error) {
	if f.listByEmployeeFn != nil {
		return f.listByEmployeeFn(ctx, employeeID)
	}
	return nil, nil
}

func (f *fakeDailyLogbookDetailService) CreateDailyLogbookDetail(ctx context.Context, detail domain.DailyLogbookDetail) error {
	if f.createFn != nil {
		return f.createFn(ctx, detail)
	}
	return nil
}

func (f *fakeDailyLogbookDetailService) UpdateDailyLogbookDetail(ctx context.Context, detail domain.DailyLogbookDetail) error {
	if f.updateFn != nil {
		return f.updateFn(ctx, detail)
	}
	return nil
}

func (f *fakeDailyLogbookDetailService) ValidateTimeSequence(outTime, takeoffTime, landingTime, inTime string) error {
	if f.validateTimeFn != nil {
		return f.validateTimeFn(outTime, takeoffTime, landingTime, inTime)
	}
	return nil
}

func (f *fakeDailyLogbookDetailService) CreateDailyLogbookDetailTx(ctx context.Context, tx output.Tx, detail domain.DailyLogbookDetail) error {
	if f.createTxFn != nil {
		return f.createTxFn(ctx, tx, detail)
	}
	return nil
}

func (f *fakeDailyLogbookDetailService) UpdateDailyLogbookDetailTx(ctx context.Context, tx output.Tx, detail domain.DailyLogbookDetail) error {
	if f.updateTxFn != nil {
		return f.updateTxFn(ctx, tx, detail)
	}
	return nil
}

func TestNewDailyLogbookDetailInteractor(t *testing.T) {
	svc := &fakeDailyLogbookDetailService{}
	logbookSvc := &fakeDailyLogbookService{}
	inter := NewDailyLogbookDetailInteractor(svc, logbookSvc)
	if inter == nil {
		t.Error("expected non-nil interactor")
	}
}

func TestDailyLogbookDetailInteractor_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expected := &domain.DailyLogbookDetail{ID: "d-1", FlightNumber: "AV123"}
		svc := &fakeDailyLogbookDetailService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbookDetail, error) {
				return expected, nil
			},
		}
		inter := NewDailyLogbookDetailInteractor(svc, &fakeDailyLogbookService{})
		result, err := inter.GetDailyLogbookDetailByID(context.Background(), "trace-1", "d-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.ID != "d-1" {
			t.Errorf("expected ID 'd-1', got %q", result.ID)
		}
	})

	t.Run("not found returns nil", func(t *testing.T) {
		svc := &fakeDailyLogbookDetailService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbookDetail, error) {
				return nil, nil
			},
		}
		inter := NewDailyLogbookDetailInteractor(svc, &fakeDailyLogbookService{})
		_, err := inter.GetDailyLogbookDetailByID(context.Background(), "trace-1", "d-1")
		if err == nil {
			t.Error("expected ErrFlightNotFound error")
		}
	})

	t.Run("service error", func(t *testing.T) {
		svc := &fakeDailyLogbookDetailService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbookDetail, error) {
				return nil, errors.New("db error")
			},
		}
		inter := NewDailyLogbookDetailInteractor(svc, &fakeDailyLogbookService{})
		_, err := inter.GetDailyLogbookDetailByID(context.Background(), "trace-1", "d-1")
		if err == nil {
			t.Error("expected error")
		}
	})
}

func TestDailyLogbookDetailInteractor_ListByLogbook(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		logbookSvc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: "lb-1", EmployeeID: "emp-1"}, nil
			},
		}
		svc := &fakeDailyLogbookDetailService{
			listByLogbookFn: func(ctx context.Context, logbookID string) ([]domain.DailyLogbookDetail, error) {
				return []domain.DailyLogbookDetail{{ID: "d-1"}, {ID: "d-2"}}, nil
			},
		}
		inter := NewDailyLogbookDetailInteractor(svc, logbookSvc)
		result, err := inter.ListDailyLogbookDetailsByLogbook(context.Background(), "trace-1", "lb-1", "emp-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 2 {
			t.Errorf("expected 2 details, got %d", len(result))
		}
	})

	t.Run("logbook not found", func(t *testing.T) {
		logbookSvc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return nil, nil
			},
		}
		svc := &fakeDailyLogbookDetailService{}
		inter := NewDailyLogbookDetailInteractor(svc, logbookSvc)
		_, err := inter.ListDailyLogbookDetailsByLogbook(context.Background(), "trace-1", "lb-1", "emp-1")
		if err == nil {
			t.Error("expected error for missing logbook")
		}
	})

	t.Run("logbook service error", func(t *testing.T) {
		logbookSvc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return nil, errors.New("db error")
			},
		}
		svc := &fakeDailyLogbookDetailService{}
		inter := NewDailyLogbookDetailInteractor(svc, logbookSvc)
		_, err := inter.ListDailyLogbookDetailsByLogbook(context.Background(), "trace-1", "lb-1", "emp-1")
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("unauthorized", func(t *testing.T) {
		logbookSvc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: "lb-1", EmployeeID: "emp-other"}, nil
			},
		}
		svc := &fakeDailyLogbookDetailService{}
		inter := NewDailyLogbookDetailInteractor(svc, logbookSvc)
		_, err := inter.ListDailyLogbookDetailsByLogbook(context.Background(), "trace-1", "lb-1", "emp-1")
		if err == nil {
			t.Error("expected unauthorized error")
		}
	})

	t.Run("list service error", func(t *testing.T) {
		logbookSvc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: "lb-1", EmployeeID: "emp-1"}, nil
			},
		}
		svc := &fakeDailyLogbookDetailService{
			listByLogbookFn: func(ctx context.Context, logbookID string) ([]domain.DailyLogbookDetail, error) {
				return nil, errors.New("list error")
			},
		}
		inter := NewDailyLogbookDetailInteractor(svc, logbookSvc)
		_, err := inter.ListDailyLogbookDetailsByLogbook(context.Background(), "trace-1", "lb-1", "emp-1")
		if err == nil {
			t.Error("expected list error")
		}
	})
}

func TestDailyLogbookDetailInteractor_ListByEmployee(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &fakeDailyLogbookDetailService{
			listByEmployeeFn: func(ctx context.Context, employeeID string) ([]domain.DailyLogbookDetail, error) {
				return []domain.DailyLogbookDetail{{ID: "d-1"}}, nil
			},
		}
		inter := NewDailyLogbookDetailInteractor(svc, &fakeDailyLogbookService{})
		result, err := inter.ListDailyLogbookDetailsByEmployee(context.Background(), "trace-1", "emp-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 1 {
			t.Errorf("expected 1, got %d", len(result))
		}
	})

	t.Run("error", func(t *testing.T) {
		svc := &fakeDailyLogbookDetailService{
			listByEmployeeFn: func(ctx context.Context, employeeID string) ([]domain.DailyLogbookDetail, error) {
				return nil, errors.New("db error")
			},
		}
		inter := NewDailyLogbookDetailInteractor(svc, &fakeDailyLogbookService{})
		_, err := inter.ListDailyLogbookDetailsByEmployee(context.Background(), "trace-1", "emp-1")
		if err == nil {
			t.Error("expected error")
		}
	})
}

func TestDailyLogbookDetailInteractor_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &fakeDailyLogbookDetailService{
			createTxFn: func(ctx context.Context, tx output.Tx, detail domain.DailyLogbookDetail) error {
				return nil
			},
		}
		logbookSvc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: "lb-1", EmployeeID: "emp-1"}, nil
			},
		}
		inter := NewDailyLogbookDetailInteractor(svc, logbookSvc)
		err := inter.CreateDailyLogbookDetail(context.Background(), "trace-1", domain.DailyLogbookDetail{
			ID: "d-1", DailyLogbookID: "lb-1", OutTime: strPtr("08:00"), TakeoffTime: strPtr("08:15"), LandingTime: strPtr("09:30"), InTime: strPtr("09:45"),
		}, "emp-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("time validation error", func(t *testing.T) {
		svc := &fakeDailyLogbookDetailService{
			validateTimeFn: func(out, takeoff, landing, in string) error {
				return domain.ErrFlightInvalidTimeSequence
			},
		}
		logbookSvc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: "lb-1", EmployeeID: "emp-1"}, nil
			},
		}
		inter := NewDailyLogbookDetailInteractor(svc, logbookSvc)
		err := inter.CreateDailyLogbookDetail(context.Background(), "trace-1", domain.DailyLogbookDetail{
			DailyLogbookID: "lb-1", OutTime: strPtr("10:00"), TakeoffTime: strPtr("08:15"), LandingTime: strPtr("09:30"), InTime: strPtr("09:45"),
		}, "emp-1")
		if err == nil {
			t.Error("expected validation error")
		}
	})

	t.Run("create error", func(t *testing.T) {
		svc := &fakeDailyLogbookDetailService{
			createTxFn: func(ctx context.Context, tx output.Tx, detail domain.DailyLogbookDetail) error {
				return errors.New("create error")
			},
		}
		logbookSvc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: "lb-1", EmployeeID: "emp-1"}, nil
			},
		}
		inter := NewDailyLogbookDetailInteractor(svc, logbookSvc)
		err := inter.CreateDailyLogbookDetail(context.Background(), "trace-1", domain.DailyLogbookDetail{
			ID: "d-1", DailyLogbookID: "lb-1", OutTime: strPtr("08:00"), TakeoffTime: strPtr("08:15"), LandingTime: strPtr("09:30"), InTime: strPtr("09:45"),
		}, "emp-1")
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("ownership error", func(t *testing.T) {
		logbookSvc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: "lb-1", EmployeeID: "emp-other"}, nil
			},
		}
		inter := NewDailyLogbookDetailInteractor(&fakeDailyLogbookDetailService{}, logbookSvc)
		err := inter.CreateDailyLogbookDetail(context.Background(), "trace-1", domain.DailyLogbookDetail{
			DailyLogbookID: "lb-1",
		}, "emp-1")
		if err == nil {
			t.Error("expected ownership error")
		}
	})

	t.Run("success without time fields", func(t *testing.T) {
		svc := &fakeDailyLogbookDetailService{
			createTxFn: func(ctx context.Context, tx output.Tx, detail domain.DailyLogbookDetail) error {
				return nil
			},
		}
		logbookSvc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: "lb-1", EmployeeID: "emp-1"}, nil
			},
		}
		inter := NewDailyLogbookDetailInteractor(svc, logbookSvc)
		err := inter.CreateDailyLogbookDetail(context.Background(), "trace-1", domain.DailyLogbookDetail{
			ID: "d-1", DailyLogbookID: "lb-1",
		}, "emp-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("generates ID when empty", func(t *testing.T) {
		svc := &fakeDailyLogbookDetailService{
			createTxFn: func(ctx context.Context, tx output.Tx, detail domain.DailyLogbookDetail) error {
				if detail.ID == "" {
					t.Error("expected non-empty ID")
				}
				return nil
			},
		}
		logbookSvc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: "lb-1", EmployeeID: "emp-1"}, nil
			},
		}
		inter := NewDailyLogbookDetailInteractor(svc, logbookSvc)
		err := inter.CreateDailyLogbookDetail(context.Background(), "trace-1", domain.DailyLogbookDetail{
			DailyLogbookID: "lb-1",
		}, "emp-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestDailyLogbookDetailInteractor_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &fakeDailyLogbookDetailService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbookDetail, error) {
				return &domain.DailyLogbookDetail{ID: "d-1", DailyLogbookID: "lb-1"}, nil
			},
			updateTxFn: func(ctx context.Context, tx output.Tx, detail domain.DailyLogbookDetail) error {
				return nil
			},
		}
		logbookSvc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: "lb-1", EmployeeID: "emp-1"}, nil
			},
		}
		inter := NewDailyLogbookDetailInteractor(svc, logbookSvc)
		err := inter.UpdateDailyLogbookDetail(context.Background(), "trace-1", domain.DailyLogbookDetail{
			ID: "d-1", OutTime: strPtr("08:00"), TakeoffTime: strPtr("08:15"), LandingTime: strPtr("09:30"), InTime: strPtr("09:45"),
		}, "emp-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		svc := &fakeDailyLogbookDetailService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbookDetail, error) {
				return nil, nil
			},
		}
		inter := NewDailyLogbookDetailInteractor(svc, &fakeDailyLogbookService{})
		err := inter.UpdateDailyLogbookDetail(context.Background(), "trace-1", domain.DailyLogbookDetail{
			ID: "d-1", OutTime: strPtr("08:00"), TakeoffTime: strPtr("08:15"), LandingTime: strPtr("09:30"), InTime: strPtr("09:45"),
		}, "emp-1")
		if err == nil {
			t.Error("expected error for not found")
		}
	})

	t.Run("time validation error", func(t *testing.T) {
		svc := &fakeDailyLogbookDetailService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbookDetail, error) {
				return &domain.DailyLogbookDetail{ID: "d-1", DailyLogbookID: "lb-1"}, nil
			},
			validateTimeFn: func(out, takeoff, landing, in string) error {
				return domain.ErrFlightInvalidTimeSequence
			},
		}
		logbookSvc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: "lb-1", EmployeeID: "emp-1"}, nil
			},
		}
		inter := NewDailyLogbookDetailInteractor(svc, logbookSvc)
		err := inter.UpdateDailyLogbookDetail(context.Background(), "trace-1", domain.DailyLogbookDetail{
			ID: "d-1", OutTime: strPtr("10:00"), TakeoffTime: strPtr("08:15"), LandingTime: strPtr("09:30"), InTime: strPtr("09:45"),
		}, "emp-1")
		if err == nil {
			t.Error("expected validation error")
		}
	})

	t.Run("ownership error", func(t *testing.T) {
		svc := &fakeDailyLogbookDetailService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbookDetail, error) {
				return &domain.DailyLogbookDetail{ID: "d-1", DailyLogbookID: "lb-1"}, nil
			},
		}
		logbookSvc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: "lb-1", EmployeeID: "emp-other"}, nil
			},
		}
		inter := NewDailyLogbookDetailInteractor(svc, logbookSvc)
		err := inter.UpdateDailyLogbookDetail(context.Background(), "trace-1", domain.DailyLogbookDetail{
			ID: "d-1",
		}, "emp-1")
		if err == nil {
			t.Error("expected ownership error")
		}
	})

	t.Run("update tx error", func(t *testing.T) {
		svc := &fakeDailyLogbookDetailService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbookDetail, error) {
				return &domain.DailyLogbookDetail{ID: "d-1", DailyLogbookID: "lb-1"}, nil
			},
			updateTxFn: func(ctx context.Context, tx output.Tx, detail domain.DailyLogbookDetail) error {
				return errors.New("update error")
			},
		}
		logbookSvc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: "lb-1", EmployeeID: "emp-1"}, nil
			},
		}
		inter := NewDailyLogbookDetailInteractor(svc, logbookSvc)
		err := inter.UpdateDailyLogbookDetail(context.Background(), "trace-1", domain.DailyLogbookDetail{
			ID: "d-1",
		}, "emp-1")
		if err == nil {
			t.Error("expected update error")
		}
	})

	t.Run("success without time fields", func(t *testing.T) {
		svc := &fakeDailyLogbookDetailService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbookDetail, error) {
				return &domain.DailyLogbookDetail{ID: "d-1", DailyLogbookID: "lb-1"}, nil
			},
			updateTxFn: func(ctx context.Context, tx output.Tx, detail domain.DailyLogbookDetail) error {
				return nil
			},
		}
		logbookSvc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: "lb-1", EmployeeID: "emp-1"}, nil
			},
		}
		inter := NewDailyLogbookDetailInteractor(svc, logbookSvc)
		err := inter.UpdateDailyLogbookDetail(context.Background(), "trace-1", domain.DailyLogbookDetail{
			ID: "d-1",
		}, "emp-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("get by ID error", func(t *testing.T) {
		svc := &fakeDailyLogbookDetailService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbookDetail, error) {
				return nil, errors.New("db error")
			},
		}
		inter := NewDailyLogbookDetailInteractor(svc, &fakeDailyLogbookService{})
		err := inter.UpdateDailyLogbookDetail(context.Background(), "trace-1", domain.DailyLogbookDetail{
			ID: "d-1",
		}, "emp-1")
		if err == nil {
			t.Error("expected error")
		}
	})
}

func TestDailyLogbookDetailInteractor_VerifyLogbookOwnership(t *testing.T) {
	t.Run("owner matches", func(t *testing.T) {
		logbookSvc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: "lb-1", EmployeeID: "emp-1"}, nil
			},
		}
		inter := NewDailyLogbookDetailInteractor(&fakeDailyLogbookDetailService{}, logbookSvc)
		err := inter.VerifyLogbookOwnership(context.Background(), "lb-1", "emp-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("owner does not match", func(t *testing.T) {
		logbookSvc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: "lb-1", EmployeeID: "emp-1"}, nil
			},
		}
		inter := NewDailyLogbookDetailInteractor(&fakeDailyLogbookDetailService{}, logbookSvc)
		err := inter.VerifyLogbookOwnership(context.Background(), "lb-1", "emp-other")
		if err == nil {
			t.Error("expected unauthorized error")
		}
	})

	t.Run("logbook not found", func(t *testing.T) {
		logbookSvc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return nil, nil
			},
		}
		inter := NewDailyLogbookDetailInteractor(&fakeDailyLogbookDetailService{}, logbookSvc)
		err := inter.VerifyLogbookOwnership(context.Background(), "lb-1", "emp-1")
		if err == nil {
			t.Error("expected error for missing logbook")
		}
	})

	t.Run("service error", func(t *testing.T) {
		logbookSvc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return nil, errors.New("db error")
			},
		}
		inter := NewDailyLogbookDetailInteractor(&fakeDailyLogbookDetailService{}, logbookSvc)
		err := inter.VerifyLogbookOwnership(context.Background(), "lb-1", "emp-1")
		if err == nil {
			t.Error("expected error")
		}
	})
}

func TestDailyLogbookDetailInteractor_GetLogbookOwner(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		logbookSvc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: "lb-1", EmployeeID: "emp-1"}, nil
			},
		}
		inter := NewDailyLogbookDetailInteractor(&fakeDailyLogbookDetailService{}, logbookSvc)
		owner, err := inter.GetLogbookOwner(context.Background(), "lb-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if owner != "emp-1" {
			t.Errorf("expected 'emp-1', got %q", owner)
		}
	})

	t.Run("not found", func(t *testing.T) {
		logbookSvc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return nil, nil
			},
		}
		inter := NewDailyLogbookDetailInteractor(&fakeDailyLogbookDetailService{}, logbookSvc)
		_, err := inter.GetLogbookOwner(context.Background(), "lb-1")
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("service error", func(t *testing.T) {
		logbookSvc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return nil, errors.New("db error")
			},
		}
		inter := NewDailyLogbookDetailInteractor(&fakeDailyLogbookDetailService{}, logbookSvc)
		_, err := inter.GetLogbookOwner(context.Background(), "lb-1")
		if err == nil {
			t.Error("expected error")
		}
	})
}

func TestDailyLogbookDetailInteractor_GetDetailLogbookOwner(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		detailSvc := &fakeDailyLogbookDetailService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbookDetail, error) {
				return &domain.DailyLogbookDetail{ID: "d-1", DailyLogbookID: "lb-1"}, nil
			},
		}
		logbookSvc := &fakeDailyLogbookService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbook, error) {
				return &domain.DailyLogbook{ID: "lb-1", EmployeeID: "emp-1"}, nil
			},
		}
		inter := NewDailyLogbookDetailInteractor(detailSvc, logbookSvc)
		owner, err := inter.GetDetailLogbookOwner(context.Background(), "d-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if owner != "emp-1" {
			t.Errorf("expected 'emp-1', got %q", owner)
		}
	})

	t.Run("detail not found", func(t *testing.T) {
		detailSvc := &fakeDailyLogbookDetailService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbookDetail, error) {
				return nil, nil
			},
		}
		inter := NewDailyLogbookDetailInteractor(detailSvc, &fakeDailyLogbookService{})
		_, err := inter.GetDetailLogbookOwner(context.Background(), "d-1")
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("service error", func(t *testing.T) {
		detailSvc := &fakeDailyLogbookDetailService{
			getByIDFn: func(ctx context.Context, id string) (*domain.DailyLogbookDetail, error) {
				return nil, errors.New("db error")
			},
		}
		inter := NewDailyLogbookDetailInteractor(detailSvc, &fakeDailyLogbookService{})
		_, err := inter.GetDetailLogbookOwner(context.Background(), "d-1")
		if err == nil {
			t.Error("expected error")
		}
	})
}
