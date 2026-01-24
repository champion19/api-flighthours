package interactor

import (
	"context"
	"errors"
	"testing"

	"github.com/Nerzal/gocloak/v13"
	"github.com/champion19/flighthours-api/core/interactor/dto"
	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/input"
	"github.com/champion19/flighthours-api/core/ports/output"
	platformLogger "github.com/champion19/flighthours-api/platform/logger"
)

type noopLogger struct{}

func (noopLogger) Info(string, ...any)    {}
func (noopLogger) Error(string, ...any)   {}
func (noopLogger) Debug(string, ...any)   {}
func (noopLogger) Warn(string, ...any)    {}
func (noopLogger) Success(string, ...any) {}
func (noopLogger) Fatal(string, ...any)   {}
func (noopLogger) Panic(string, ...any)   {}
func (noopLogger) WithTraceID(string) platformLogger.Logger {
	return noopLogger{}
}

type fakeTx struct {
	commitFn   func() error
	rollbackFn func() error

	committed  bool
	rolledBack bool
}

func (t *fakeTx) Commit() error {
	t.committed = true
	if t.commitFn != nil {
		return t.commitFn()
	}
	return nil
}
func (t *fakeTx) Rollback() error {
	t.rolledBack = true
	if t.rollbackFn != nil {
		return t.rollbackFn()
	}
	return nil
}

type fakeService struct {
	registerEmployeeFn func(ctx context.Context, employee domain.Employee) (*dto.RegisterEmployee, error)
	checkAndCleanFn    func(ctx context.Context, email string) error
	beginTxFn          func(ctx context.Context) (output.Tx, error)
	saveEmployeeFn     func(ctx context.Context, tx output.Tx, employee domain.Employee) error
	createUserFn       func(ctx context.Context, employee *domain.Employee) (string, error)
	setPasswordFn      func(ctx context.Context, userID string, password string) error
	assignRoleFn       func(ctx context.Context, userID string, role string) error
	updateKcIDFn       func(ctx context.Context, tx output.Tx, employeeID string, keycloakUserID string) error
	rollbackKcFn       func(ctx context.Context, kcID string) error
}

var _ input.Service = (*fakeService)(nil)

func (f *fakeService) BeginTx(ctx context.Context) (output.Tx, error) {
	return f.beginTxFn(ctx)
}
func (f *fakeService) RegisterEmployee(ctx context.Context, employee domain.Employee) (*dto.RegisterEmployee, error) {
	return f.registerEmployeeFn(ctx, employee)
}
func (f *fakeService) GetEmployeeByEmail(context.Context, string) (*domain.Employee, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeService) GetEmployeeByID(context.Context, string) (*domain.Employee, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeService) LocateEmployee(context.Context, string) (*dto.RegisterEmployee, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeService) CheckAndCleanInconsistentState(ctx context.Context, email string) error {
	return f.checkAndCleanFn(ctx, email)
}
func (f *fakeService) SaveEmployeeToDB(ctx context.Context, tx output.Tx, employee domain.Employee) error {
	return f.saveEmployeeFn(ctx, tx, employee)
}
func (f *fakeService) UpdateEmployeeKeycloakID(ctx context.Context, tx output.Tx, employeeID string, keycloakUserID string) error {
	return f.updateKcIDFn(ctx, tx, employeeID, keycloakUserID)
}
func (f *fakeService) CreateUserInKeycloak(ctx context.Context, employee *domain.Employee) (string, error) {
	return f.createUserFn(ctx, employee)
}
func (f *fakeService) SetUserPassword(ctx context.Context, userID string, password string) error {
	return f.setPasswordFn(ctx, userID, password)
}
func (f *fakeService) AssignUserRole(ctx context.Context, userID string, role string) error {
	return f.assignRoleFn(ctx, userID, role)
}
func (f *fakeService) SendVerificationEmail(context.Context, string) error {
	return nil
}
func (f *fakeService) SendPasswordResetEmail(context.Context, string) error {
	return nil
}
func (f *fakeService) Login(context.Context, string, string) (*gocloak.JWT, error) {
	return &gocloak.JWT{}, nil
}
func (f *fakeService) VerifyEmailByToken(context.Context, string) (string, error) {
	return "", nil
}
func (f *fakeService) GetUserByEmail(context.Context, string) (*gocloak.User, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeService) RollbackEmployee(context.Context, string) error {
	return errors.New("not implemented")
}
func (f *fakeService) RollbackKeycloakUser(ctx context.Context, kcID string) error {
	return f.rollbackKcFn(ctx, kcID)
}
func (f *fakeService) UpdatePassword(context.Context, string, string) (string, error) {
	return "", nil
}
func (f *fakeService) UpdateEmployee(context.Context, domain.Employee, bool, string) error {
	return nil
}
func (f *fakeService) ChangePassword(context.Context, string, string, string) (string, error) {
	return "", nil
}
func (f *fakeService) DeleteEmployee(context.Context, string, string) error {
	return nil
}
func (f *fakeService) GetEmployeeByKeycloakID(context.Context, string) (*domain.Employee, error) {
	return nil, errors.New("not implemented")
}

// GetEmployeesByRole mock for HU47 - Virtual Entity pattern
func (f *fakeService) GetEmployeesByRole(context.Context, string) ([]domain.Employee, error) {
	return nil, errors.New("not implemented")
}

func TestInteractor_RegisterEmployee(t *testing.T) {
	ctx := context.Background()

	mkEmployee := func() domain.Employee {
		return domain.Employee{Email: "a@b.com", Password: "pw", Role: "user"}
	}

	t.Run("service returns ErrIncompleteRegistration => cleanup called and returns same error", func(t *testing.T) {
		cleanupCalled := 0
		svc := &fakeService{
			registerEmployeeFn: func(context.Context, domain.Employee) (*dto.RegisterEmployee, error) {
				return nil, domain.ErrIncompleteRegistration
			},
			checkAndCleanFn: func(context.Context, string) error {
				cleanupCalled++
				return nil
			},
		}
		i := NewInteractor(svc, noopLogger{})

		res, err := i.RegisterEmployee(ctx, mkEmployee())
		if res != nil {
			t.Fatalf("expected nil result")
		}
		if !errors.Is(err, domain.ErrIncompleteRegistration) {
			t.Fatalf("expected %v, got %v", domain.ErrIncompleteRegistration, err)
		}
		if cleanupCalled != 1 {
			t.Fatalf("expected cleanup called 1 time, got %d", cleanupCalled)
		}
	})

	t.Run("service returns ErrIncompleteRegistration and cleanup fails => returns cleanup error", func(t *testing.T) {
		cleanupErr := errors.New("cleanup failed")
		svc := &fakeService{
			registerEmployeeFn: func(context.Context, domain.Employee) (*dto.RegisterEmployee, error) {
				return nil, domain.ErrIncompleteRegistration
			},
			checkAndCleanFn: func(context.Context, string) error {
				return cleanupErr
			},
		}
		i := NewInteractor(svc, noopLogger{})

		res, err := i.RegisterEmployee(ctx, mkEmployee())
		if res != nil {
			t.Fatalf("expected nil result")
		}
		if !errors.Is(err, cleanupErr) {
			t.Fatalf("expected %v, got %v", cleanupErr, err)
		}
	})

	t.Run("happy path => commit called and result populated", func(t *testing.T) {
		tx := &fakeTx{}
		called := struct {
			checkClean int
			save       int
			createKC   int
			setPwd     int
			assignRole int
			patchKC    int
			rollbackKC int
		}{}

		svc := &fakeService{
			registerEmployeeFn: func(ctx context.Context, employee domain.Employee) (*dto.RegisterEmployee, error) {
				return &dto.RegisterEmployee{Employee: employee, Message: "ok"}, nil
			},
			checkAndCleanFn: func(context.Context, string) error {
				called.checkClean++
				return nil
			},
			beginTxFn: func(context.Context) (output.Tx, error) { return tx, nil },
			saveEmployeeFn: func(context.Context, output.Tx, domain.Employee) error {
				called.save++
				return nil
			},
			createUserFn: func(context.Context, *domain.Employee) (string, error) {
				called.createKC++
				return "kc1", nil
			},
			setPasswordFn: func(context.Context, string, string) error {
				called.setPwd++
				return nil
			},
			assignRoleFn: func(context.Context, string, string) error {
				called.assignRole++
				return nil
			},
			updateKcIDFn: func(context.Context, output.Tx, string, string) error {
				called.patchKC++
				return nil
			},
			rollbackKcFn: func(context.Context, string) error {
				called.rollbackKC++
				return nil
			},
		}

		i := NewInteractor(svc, noopLogger{})
		res, err := i.RegisterEmployee(ctx, mkEmployee())
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if res == nil {
			t.Fatalf("expected non-nil result")
		}
		if res.Employee.KeycloakUserID != "kc1" {
			t.Fatalf("expected KeycloakUserID kc1, got %q", res.Employee.KeycloakUserID)
		}
		if res.Employee.ID == "" {
			t.Fatalf("expected employee ID to be set")
		}
		if !tx.committed {
			t.Fatalf("expected tx.Commit to be called")
		}
		if tx.rolledBack {
			t.Fatalf("did not expect rollback")
		}
		if called.checkClean != 1 || called.save != 1 || called.createKC != 1 || called.setPwd != 1 || called.assignRole != 1 || called.patchKC != 1 {
			t.Fatalf("unexpected call counts: %+v", called)
		}
		if called.rollbackKC != 0 {
			t.Fatalf("did not expect keycloak rollback")
		}
	})

	t.Run("error after keycloak created => rollback tx and rollback keycloak", func(t *testing.T) {
		tx := &fakeTx{}
		calledRollbackKC := 0

		svc := &fakeService{
			registerEmployeeFn: func(ctx context.Context, employee domain.Employee) (*dto.RegisterEmployee, error) {
				return &dto.RegisterEmployee{Employee: employee, Message: "ok"}, nil
			},
			checkAndCleanFn: func(context.Context, string) error { return nil },
			beginTxFn:       func(context.Context) (output.Tx, error) { return tx, nil },
			saveEmployeeFn:  func(context.Context, output.Tx, domain.Employee) error { return nil },
			createUserFn: func(context.Context, *domain.Employee) (string, error) {
				return "kc1", nil
			},
			setPasswordFn: func(context.Context, string, string) error {
				return errors.New("set password failed")
			},
			assignRoleFn: func(context.Context, string, string) error { return nil },
			updateKcIDFn: func(context.Context, output.Tx, string, string) error { return nil },
			rollbackKcFn: func(context.Context, string) error {
				calledRollbackKC++
				return nil
			},
		}

		i := NewInteractor(svc, noopLogger{})
		_, err := i.RegisterEmployee(ctx, mkEmployee())
		if err == nil {
			t.Fatalf("expected error")
		}
		if !tx.rolledBack {
			t.Fatalf("expected tx.Rollback to be called")
		}
		if calledRollbackKC != 1 {
			t.Fatalf("expected keycloak rollback called once, got %d", calledRollbackKC)
		}
	})
}
