package interactor

import (
	"context"
	"errors"
	"testing"

	"github.com/Nerzal/gocloak/v13"
	"github.com/champion19/api-flighthours/core/interactor/dto"
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/input"
	"github.com/champion19/api-flighthours/core/ports/output"
	platformLogger "github.com/champion19/api-flighthours/platform/logger"
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
	registerEmployeeFn      func(ctx context.Context, employee domain.Employee) (*dto.RegisterEmployee, error)
	checkAndCleanFn         func(ctx context.Context, email string) error
	beginTxFn               func(ctx context.Context) (output.Tx, error)
	saveEmployeeFn          func(ctx context.Context, tx output.Tx, employee domain.Employee) error
	createUserFn            func(ctx context.Context, employee *domain.Employee) (string, error)
	setPasswordFn           func(ctx context.Context, userID string, password string) error
	assignRoleFn            func(ctx context.Context, userID string, role string) error
	updateKcIDFn            func(ctx context.Context, tx output.Tx, employeeID string, keycloakUserID string) error
	rollbackKcFn            func(ctx context.Context, kcID string) error
	loginFn                 func(ctx context.Context, email, password string) (*gocloak.JWT, error)
	getEmployeeByIDFn       func(ctx context.Context, id string) (*domain.Employee, error)
	deleteEmployeeFn        func(ctx context.Context, id, kcID string) error
	updateEmployeeFn        func(ctx context.Context, tx output.Tx, e domain.Employee) error
	sendPasswordResetFn     func(ctx context.Context, email string) error
	updatePasswordFn        func(ctx context.Context, token, newPass string) (string, error)
	changePasswordFn        func(ctx context.Context, email, current, newPass string) (string, error)
	verifyEmailByTokenFn    func(ctx context.Context, token string) (string, error)
	getUserByEmailFn        func(ctx context.Context, email string) (*gocloak.User, error)
	sendVerificationEmailFn func(ctx context.Context, userID string) error
	refreshTokenFn          func(ctx context.Context, refreshToken string) (*gocloak.JWT, error)
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
func (f *fakeService) GetEmployeeByID(ctx context.Context, id string) (*domain.Employee, error) {
	if f.getEmployeeByIDFn != nil {
		return f.getEmployeeByIDFn(ctx, id)
	}
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
func (f *fakeService) SendVerificationEmail(ctx context.Context, userID string) error {
	if f.sendVerificationEmailFn != nil {
		return f.sendVerificationEmailFn(ctx, userID)
	}
	return nil
}
func (f *fakeService) SendPasswordResetEmail(ctx context.Context, email string) error {
	if f.sendPasswordResetFn != nil {
		return f.sendPasswordResetFn(ctx, email)
	}
	return nil
}
func (f *fakeService) Login(ctx context.Context, email, password string) (*gocloak.JWT, error) {
	if f.loginFn != nil {
		return f.loginFn(ctx, email, password)
	}
	return &gocloak.JWT{}, nil
}
func (f *fakeService) RefreshToken(ctx context.Context, refreshToken string) (*gocloak.JWT, error) {
	if f.refreshTokenFn != nil {
		return f.refreshTokenFn(ctx, refreshToken)
	}
	return &gocloak.JWT{}, nil
}
func (f *fakeService) VerifyEmailByToken(ctx context.Context, token string) (string, error) {
	if f.verifyEmailByTokenFn != nil {
		return f.verifyEmailByTokenFn(ctx, token)
	}
	return "", nil
}
func (f *fakeService) GetUserByEmail(ctx context.Context, email string) (*gocloak.User, error) {
	if f.getUserByEmailFn != nil {
		return f.getUserByEmailFn(ctx, email)
	}
	return nil, errors.New("not implemented")
}
func (f *fakeService) RollbackEmployee(context.Context, string) error {
	return errors.New("not implemented")
}
func (f *fakeService) RollbackKeycloakUser(ctx context.Context, kcID string) error {
	return f.rollbackKcFn(ctx, kcID)
}
func (f *fakeService) UpdatePassword(ctx context.Context, token, newPass string) (string, error) {
	if f.updatePasswordFn != nil {
		return f.updatePasswordFn(ctx, token, newPass)
	}
	return "", nil
}
func (f *fakeService) ChangePassword(ctx context.Context, email, current, newPass string) (string, error) {
	if f.changePasswordFn != nil {
		return f.changePasswordFn(ctx, email, current, newPass)
	}
	return "", nil
}
func (f *fakeService) DeleteEmployee(ctx context.Context, id, kcID string) error {
	if f.deleteEmployeeFn != nil {
		return f.deleteEmployeeFn(ctx, id, kcID)
	}
	return nil
}
func (f *fakeService) GetEmployeeByKeycloakID(context.Context, string) (*domain.Employee, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeService) UpdateEmployee(ctx context.Context, tx output.Tx, e domain.Employee) error {
	if f.updateEmployeeFn != nil {
		return f.updateEmployeeFn(ctx, tx, e)
	}
	return nil
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
		i := NewInteractor(svc)

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
		i := NewInteractor(svc)

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

		i := NewInteractor(svc)
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

		i := NewInteractor(svc)
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

	t.Run("error causes tx rollback to fail => logs error but still returns original error", func(t *testing.T) {
		tx := &fakeTx{rollbackFn: func() error { return errors.New("rollback failed") }}

		svc := &fakeService{
			registerEmployeeFn: func(ctx context.Context, employee domain.Employee) (*dto.RegisterEmployee, error) {
				return &dto.RegisterEmployee{Employee: employee, Message: "ok"}, nil
			},
			checkAndCleanFn: func(context.Context, string) error { return nil },
			beginTxFn:       func(context.Context) (output.Tx, error) { return tx, nil },
			saveEmployeeFn: func(context.Context, output.Tx, domain.Employee) error {
				return errors.New("save failed")
			},
		}

		i := NewInteractor(svc)
		_, err := i.RegisterEmployee(ctx, mkEmployee())
		if err == nil {
			t.Fatalf("expected error")
		}
		// tx.Rollback was called (even though it failed)
		if !tx.rolledBack {
			t.Fatalf("expected tx.Rollback to be called")
		}
	})

	t.Run("keycloak rollback fails => logs error but still returns original error", func(t *testing.T) {
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
			rollbackKcFn: func(context.Context, string) error {
				calledRollbackKC++
				return errors.New("keycloak rollback failed")
			},
		}

		i := NewInteractor(svc)
		_, err := i.RegisterEmployee(ctx, mkEmployee())
		if err == nil {
			t.Fatalf("expected error")
		}
		if calledRollbackKC != 1 {
			t.Fatalf("expected keycloak rollback called once, got %d", calledRollbackKC)
		}
	})

	t.Run("generic service error => returns error directly", func(t *testing.T) {
		genericErr := errors.New("unexpected service failure")
		svc := &fakeService{
			registerEmployeeFn: func(context.Context, domain.Employee) (*dto.RegisterEmployee, error) {
				return nil, genericErr
			},
		}
		i := NewInteractor(svc)

		res, err := i.RegisterEmployee(ctx, mkEmployee())
		if res != nil {
			t.Fatalf("expected nil result")
		}
		if !errors.Is(err, genericErr) {
			t.Fatalf("expected %v, got %v", genericErr, err)
		}
	})

	t.Run("ErrIncompleteRegistration + CheckAndClean succeeds => returns ErrIncompleteRegistration", func(t *testing.T) {
		svc := &fakeService{
			registerEmployeeFn: func(context.Context, domain.Employee) (*dto.RegisterEmployee, error) {
				return nil, domain.ErrIncompleteRegistration
			},
			checkAndCleanFn: func(context.Context, string) error { return nil },
		}
		i := NewInteractor(svc)

		_, err := i.RegisterEmployee(ctx, mkEmployee())
		if !errors.Is(err, domain.ErrIncompleteRegistration) {
			t.Fatalf("expected %v, got %v", domain.ErrIncompleteRegistration, err)
		}
	})

	t.Run("BeginTx error => returns error", func(t *testing.T) {
		svc := &fakeService{
			registerEmployeeFn: func(ctx context.Context, employee domain.Employee) (*dto.RegisterEmployee, error) {
				return &dto.RegisterEmployee{Employee: employee, Message: "ok"}, nil
			},
			checkAndCleanFn: func(context.Context, string) error { return nil },
			beginTxFn: func(context.Context) (output.Tx, error) {
				return nil, errors.New("begin tx failed")
			},
		}
		i := NewInteractor(svc)

		_, err := i.RegisterEmployee(ctx, mkEmployee())
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("CreateUserInKeycloak error => rolls back tx", func(t *testing.T) {
		tx := &fakeTx{}
		svc := &fakeService{
			registerEmployeeFn: func(ctx context.Context, employee domain.Employee) (*dto.RegisterEmployee, error) {
				return &dto.RegisterEmployee{Employee: employee, Message: "ok"}, nil
			},
			checkAndCleanFn: func(context.Context, string) error { return nil },
			beginTxFn:       func(context.Context) (output.Tx, error) { return tx, nil },
			saveEmployeeFn:  func(context.Context, output.Tx, domain.Employee) error { return nil },
			createUserFn: func(context.Context, *domain.Employee) (string, error) {
				return "", errors.New("keycloak create failed")
			},
		}
		i := NewInteractor(svc)

		_, err := i.RegisterEmployee(ctx, mkEmployee())
		if err == nil {
			t.Fatal("expected error")
		}
		if !tx.rolledBack {
			t.Error("expected tx.Rollback to be called")
		}
	})

	t.Run("AssignUserRole error => rolls back keycloak and tx", func(t *testing.T) {
		tx := &fakeTx{}
		rollbackKCCalled := false
		svc := &fakeService{
			registerEmployeeFn: func(ctx context.Context, employee domain.Employee) (*dto.RegisterEmployee, error) {
				return &dto.RegisterEmployee{Employee: employee, Message: "ok"}, nil
			},
			checkAndCleanFn: func(context.Context, string) error { return nil },
			beginTxFn:       func(context.Context) (output.Tx, error) { return tx, nil },
			saveEmployeeFn:  func(context.Context, output.Tx, domain.Employee) error { return nil },
			createUserFn:    func(context.Context, *domain.Employee) (string, error) { return "kc1", nil },
			setPasswordFn:   func(context.Context, string, string) error { return nil },
			assignRoleFn: func(context.Context, string, string) error {
				return errors.New("assign role failed")
			},
			rollbackKcFn: func(context.Context, string) error {
				rollbackKCCalled = true
				return nil
			},
		}
		i := NewInteractor(svc)

		_, err := i.RegisterEmployee(ctx, mkEmployee())
		if err == nil {
			t.Fatal("expected error")
		}
		if !tx.rolledBack {
			t.Error("expected tx.Rollback")
		}
		if !rollbackKCCalled {
			t.Error("expected keycloak rollback")
		}
	})

	t.Run("UpdateEmployeeKeycloakID error => rolls back keycloak and tx", func(t *testing.T) {
		tx := &fakeTx{}
		svc := &fakeService{
			registerEmployeeFn: func(ctx context.Context, employee domain.Employee) (*dto.RegisterEmployee, error) {
				return &dto.RegisterEmployee{Employee: employee, Message: "ok"}, nil
			},
			checkAndCleanFn: func(context.Context, string) error { return nil },
			beginTxFn:       func(context.Context) (output.Tx, error) { return tx, nil },
			saveEmployeeFn:  func(context.Context, output.Tx, domain.Employee) error { return nil },
			createUserFn:    func(context.Context, *domain.Employee) (string, error) { return "kc1", nil },
			setPasswordFn:   func(context.Context, string, string) error { return nil },
			assignRoleFn:    func(context.Context, string, string) error { return nil },
			updateKcIDFn: func(context.Context, output.Tx, string, string) error {
				return errors.New("update kc id failed")
			},
			rollbackKcFn: func(context.Context, string) error { return nil },
		}
		i := NewInteractor(svc)

		_, err := i.RegisterEmployee(ctx, mkEmployee())
		if err == nil {
			t.Fatal("expected error")
		}
		if !tx.rolledBack {
			t.Error("expected tx.Rollback")
		}
	})

	t.Run("Commit error => returns error", func(t *testing.T) {
		tx := &fakeTx{commitFn: func() error { return errors.New("commit failed") }}
		svc := &fakeService{
			registerEmployeeFn: func(ctx context.Context, employee domain.Employee) (*dto.RegisterEmployee, error) {
				return &dto.RegisterEmployee{Employee: employee, Message: "ok"}, nil
			},
			checkAndCleanFn: func(context.Context, string) error { return nil },
			beginTxFn:       func(context.Context) (output.Tx, error) { return tx, nil },
			saveEmployeeFn:  func(context.Context, output.Tx, domain.Employee) error { return nil },
			createUserFn:    func(context.Context, *domain.Employee) (string, error) { return "kc1", nil },
			setPasswordFn:   func(context.Context, string, string) error { return nil },
			assignRoleFn:    func(context.Context, string, string) error { return nil },
			updateKcIDFn:    func(context.Context, output.Tx, string, string) error { return nil },
			rollbackKcFn:    func(context.Context, string) error { return nil },
		}
		i := NewInteractor(svc)

		_, err := i.RegisterEmployee(ctx, mkEmployee())
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("SendVerificationEmail error => still succeeds", func(t *testing.T) {
		sendEmailCalled := false
		tx := &fakeTx{}
		svc := &fakeService{
			registerEmployeeFn: func(ctx context.Context, employee domain.Employee) (*dto.RegisterEmployee, error) {
				return &dto.RegisterEmployee{Employee: employee, Message: "ok"}, nil
			},
			checkAndCleanFn: func(context.Context, string) error { return nil },
			beginTxFn:       func(context.Context) (output.Tx, error) { return tx, nil },
			saveEmployeeFn:  func(context.Context, output.Tx, domain.Employee) error { return nil },
			createUserFn:    func(context.Context, *domain.Employee) (string, error) { return "kc1", nil },
			setPasswordFn:   func(context.Context, string, string) error { return nil },
			assignRoleFn:    func(context.Context, string, string) error { return nil },
			updateKcIDFn:    func(context.Context, output.Tx, string, string) error { return nil },
			rollbackKcFn:    func(context.Context, string) error { return nil },
			sendVerificationEmailFn: func(context.Context, string) error {
				sendEmailCalled = true
				return errors.New("email send failed")
			},
		}
		i := NewInteractor(svc)

		res, err := i.RegisterEmployee(ctx, mkEmployee())
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if res == nil {
			t.Fatal("expected non-nil result")
		}
		// SendVerificationEmail runs in goroutine; give it a moment
		_ = sendEmailCalled
	})
}

func TestInteractor_Login(t *testing.T) {
	ctx := context.Background()

	t.Run("success => returns token response", func(t *testing.T) {
		svc := &fakeService{
			loginFn: func(ctx context.Context, email, password string) (*gocloak.JWT, error) {
				return &gocloak.JWT{
					AccessToken:  "access-token-123",
					RefreshToken: "refresh-token-456",
					ExpiresIn:    3600,
					TokenType:    "Bearer",
				}, nil
			},
		}
		i := NewInteractor(svc)

		result, err := i.Login(ctx, "test@example.com", "password123")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result.AccessToken != "access-token-123" {
			t.Errorf("expected access-token-123, got %s", result.AccessToken)
		}
	})

	t.Run("login fails => returns error", func(t *testing.T) {
		loginErr := errors.New("invalid credentials")
		svc := &fakeService{
			loginFn: func(ctx context.Context, email, password string) (*gocloak.JWT, error) {
				return nil, loginErr
			},
		}
		i := NewInteractor(svc)

		_, err := i.Login(ctx, "test@example.com", "wrongpassword")
		if !errors.Is(err, loginErr) {
			t.Fatalf("expected %v, got %v", loginErr, err)
		}
	})
}

func TestInteractor_DeleteEmployee(t *testing.T) {
	ctx := context.Background()

	t.Run("success => deletes employee", func(t *testing.T) {
		deleteCalled := false
		svc := &fakeService{
			getEmployeeByIDFn: func(ctx context.Context, id string) (*domain.Employee, error) {
				return &domain.Employee{ID: id, KeycloakUserID: "kc-123"}, nil
			},
			deleteEmployeeFn: func(ctx context.Context, id, kcID string) error {
				deleteCalled = true
				return nil
			},
		}
		i := NewInteractor(svc)

		err := i.DeleteEmployee(ctx, "employee-123")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !deleteCalled {
			t.Fatal("expected delete to be called")
		}
	})

	t.Run("employee not found => returns error", func(t *testing.T) {
		svc := &fakeService{
			getEmployeeByIDFn: func(ctx context.Context, id string) (*domain.Employee, error) {
				return nil, domain.ErrPersonNotFound
			},
		}
		i := NewInteractor(svc)

		err := i.DeleteEmployee(ctx, "non-existent")
		if !errors.Is(err, domain.ErrPersonNotFound) {
			t.Fatalf("expected %v, got %v", domain.ErrPersonNotFound, err)
		}
	})

	t.Run("delete service fails => returns error", func(t *testing.T) {
		deleteErr := errors.New("delete service error")
		svc := &fakeService{
			getEmployeeByIDFn: func(ctx context.Context, id string) (*domain.Employee, error) {
				return &domain.Employee{ID: "em1", KeycloakUserID: "kc1"}, nil
			},
			deleteEmployeeFn: func(ctx context.Context, id, kcID string) error {
				return deleteErr
			},
		}
		i := NewInteractor(svc)

		err := i.DeleteEmployee(ctx, "em1")
		if !errors.Is(err, deleteErr) {
			t.Fatalf("expected %v, got %v", deleteErr, err)
		}
	})
}

func TestInteractor_UpdateEmployee(t *testing.T) {
	ctx := context.Background()

	t.Run("success => returns update result", func(t *testing.T) {
		tx := &fakeTx{}
		svc := &fakeService{
			beginTxFn: func(ctx context.Context) (output.Tx, error) { return tx, nil },
			updateEmployeeFn: func(ctx context.Context, tx output.Tx, e domain.Employee) error {
				return nil
			},
		}
		i := NewInteractor(svc)

		result, err := i.UpdateEmployee(ctx, domain.Employee{ID: "e1", Name: "Updated Name"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !result.Updated {
			t.Error("expected Updated to be true")
		}
		if !tx.committed {
			t.Error("expected tx.Commit to be called")
		}
	})

	t.Run("begin tx fails => returns error", func(t *testing.T) {
		svc := &fakeService{
			beginTxFn: func(ctx context.Context) (output.Tx, error) {
				return nil, errors.New("tx error")
			},
		}
		i := NewInteractor(svc)

		_, err := i.UpdateEmployee(ctx, domain.Employee{ID: "e1"})
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("update fails => returns error", func(t *testing.T) {
		tx := &fakeTx{}
		updateErr := errors.New("update failed")
		svc := &fakeService{
			beginTxFn: func(ctx context.Context) (output.Tx, error) { return tx, nil },
			updateEmployeeFn: func(ctx context.Context, txArg output.Tx, e domain.Employee) error {
				return updateErr
			},
		}
		i := NewInteractor(svc)

		_, err := i.UpdateEmployee(ctx, domain.Employee{ID: "e1"})
		if !errors.Is(err, updateErr) {
			t.Fatalf("expected %v, got %v", updateErr, err)
		}
		if !tx.rolledBack {
			t.Error("expected tx.Rollback to be called")
		}
	})

	t.Run("commit fails => returns error", func(t *testing.T) {
		tx := &fakeTx{commitFn: func() error { return errors.New("commit failed") }}
		svc := &fakeService{
			beginTxFn: func(ctx context.Context) (output.Tx, error) { return tx, nil },
			updateEmployeeFn: func(ctx context.Context, txArg output.Tx, e domain.Employee) error {
				return nil
			},
		}
		i := NewInteractor(svc)

		_, err := i.UpdateEmployee(ctx, domain.Employee{ID: "e1"})
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("update fails and rollback fails => returns original error", func(t *testing.T) {
		tx := &fakeTx{rollbackFn: func() error { return errors.New("rollback failed") }}
		updateErr := errors.New("update failed")
		svc := &fakeService{
			beginTxFn: func(ctx context.Context) (output.Tx, error) { return tx, nil },
			updateEmployeeFn: func(ctx context.Context, txArg output.Tx, e domain.Employee) error {
				return updateErr
			},
		}
		i := NewInteractor(svc)

		_, err := i.UpdateEmployee(ctx, domain.Employee{ID: "e1"})
		if !errors.Is(err, updateErr) {
			t.Fatalf("expected %v, got %v", updateErr, err)
		}
	})
}

func TestInteractor_RequestPasswordReset(t *testing.T) {
	ctx := context.Background()

	t.Run("always returns nil even on error", func(t *testing.T) {
		svc := &fakeService{
			sendPasswordResetFn: func(ctx context.Context, email string) error {
				return errors.New("email sending failed")
			},
		}
		i := NewInteractor(svc)

		err := i.RequestPasswordReset(ctx, "test@example.com")
		if err != nil {
			t.Fatalf("expected no error (masked), got %v", err)
		}
	})

	t.Run("success => sends password reset email", func(t *testing.T) {
		sendCalled := false
		svc := &fakeService{
			sendPasswordResetFn: func(ctx context.Context, email string) error {
				sendCalled = true
				return nil
			},
		}
		i := NewInteractor(svc)

		err := i.RequestPasswordReset(ctx, "test@example.com")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !sendCalled {
			t.Fatal("expected SendPasswordResetEmail to be called")
		}
	})
}

func TestInteractor_UpdatePassword(t *testing.T) {
	ctx := context.Background()

	t.Run("password mismatch => returns error", func(t *testing.T) {
		svc := &fakeService{}
		i := NewInteractor(svc)

		_, err := i.UpdatePassword(ctx, "token", "newpass", "different")
		if !errors.Is(err, domain.ErrPasswordMismatch) {
			t.Fatalf("expected %v, got %v", domain.ErrPasswordMismatch, err)
		}
	})

	t.Run("success => returns email", func(t *testing.T) {
		svc := &fakeService{
			updatePasswordFn: func(ctx context.Context, token, newPass string) (string, error) {
				return "user@example.com", nil
			},
		}
		i := NewInteractor(svc)

		email, err := i.UpdatePassword(ctx, "valid-token", "newpass", "newpass")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if email != "user@example.com" {
			t.Errorf("expected user@example.com, got %s", email)
		}
	})

	t.Run("invalid token => returns ErrInvalidToken", func(t *testing.T) {
		svc := &fakeService{
			updatePasswordFn: func(ctx context.Context, token, newPass string) (string, error) {
				return "", domain.ErrInvalidToken
			},
		}
		i := NewInteractor(svc)

		_, err := i.UpdatePassword(ctx, "invalid-token", "newpass", "newpass")
		if !errors.Is(err, domain.ErrInvalidToken) {
			t.Fatalf("expected %v, got %v", domain.ErrInvalidToken, err)
		}
	})

	t.Run("password update failed => returns ErrPasswordUpdateFailed", func(t *testing.T) {
		svc := &fakeService{
			updatePasswordFn: func(ctx context.Context, token, newPass string) (string, error) {
				return "", domain.ErrPasswordUpdateFailed
			},
		}
		i := NewInteractor(svc)

		_, err := i.UpdatePassword(ctx, "valid-token", "newpass", "newpass")
		if !errors.Is(err, domain.ErrPasswordUpdateFailed) {
			t.Fatalf("expected %v, got %v", domain.ErrPasswordUpdateFailed, err)
		}
	})
}

func TestInteractor_ChangePassword(t *testing.T) {
	ctx := context.Background()

	t.Run("password mismatch => returns error", func(t *testing.T) {
		svc := &fakeService{}
		i := NewInteractor(svc)

		_, err := i.ChangePassword(ctx, "email@test.com", "current", "new1", "new2")
		if !errors.Is(err, domain.ErrPasswordMismatch) {
			t.Fatalf("expected %v, got %v", domain.ErrPasswordMismatch, err)
		}
	})

	t.Run("success => returns email", func(t *testing.T) {
		svc := &fakeService{
			changePasswordFn: func(ctx context.Context, email, current, newPass string) (string, error) {
				return email, nil
			},
		}
		i := NewInteractor(svc)

		result, err := i.ChangePassword(ctx, "email@test.com", "current", "newpass", "newpass")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result != "email@test.com" {
			t.Errorf("expected email@test.com, got %s", result)
		}
	})

	t.Run("invalid current password => returns error", func(t *testing.T) {
		svc := &fakeService{
			changePasswordFn: func(ctx context.Context, email, current, newPass string) (string, error) {
				return "", domain.ErrInvalidCurrentPassword
			},
		}
		i := NewInteractor(svc)

		_, err := i.ChangePassword(ctx, "email@test.com", "wrong", "newpass", "newpass")
		if !errors.Is(err, domain.ErrInvalidCurrentPassword) {
			t.Fatalf("expected %v, got %v", domain.ErrInvalidCurrentPassword, err)
		}
	})

	t.Run("user not found => returns error", func(t *testing.T) {
		svc := &fakeService{
			changePasswordFn: func(ctx context.Context, email, current, newPass string) (string, error) {
				return "", domain.ErrUserNotFound
			},
		}
		i := NewInteractor(svc)

		_, err := i.ChangePassword(ctx, "unknown@test.com", "current", "newpass", "newpass")
		if !errors.Is(err, domain.ErrUserNotFound) {
			t.Fatalf("expected %v, got %v", domain.ErrUserNotFound, err)
		}
	})

	t.Run("password update failed => returns error", func(t *testing.T) {
		svc := &fakeService{
			changePasswordFn: func(ctx context.Context, email, current, newPass string) (string, error) {
				return "", domain.ErrPasswordUpdateFailed
			},
		}
		i := NewInteractor(svc)

		_, err := i.ChangePassword(ctx, "email@test.com", "current", "newpass", "newpass")
		if !errors.Is(err, domain.ErrPasswordUpdateFailed) {
			t.Fatalf("expected %v, got %v", domain.ErrPasswordUpdateFailed, err)
		}
	})
}

func TestInteractor_VerifyEmailByToken(t *testing.T) {
	ctx := context.Background()

	t.Run("success => returns email", func(t *testing.T) {
		svc := &fakeService{
			verifyEmailByTokenFn: func(ctx context.Context, token string) (string, error) {
				return "verified@example.com", nil
			},
		}
		i := NewInteractor(svc)

		email, err := i.VerifyEmailByToken(ctx, "valid-token")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if email != "verified@example.com" {
			t.Errorf("expected verified@example.com, got %s", email)
		}
	})

	t.Run("invalid token => returns error", func(t *testing.T) {
		svc := &fakeService{
			verifyEmailByTokenFn: func(ctx context.Context, token string) (string, error) {
				return "", domain.ErrInvalidToken
			},
		}
		i := NewInteractor(svc)

		_, err := i.VerifyEmailByToken(ctx, "invalid-token")
		if !errors.Is(err, domain.ErrInvalidToken) {
			t.Fatalf("expected %v, got %v", domain.ErrInvalidToken, err)
		}
	})

	t.Run("user not found => returns error", func(t *testing.T) {
		svc := &fakeService{
			verifyEmailByTokenFn: func(ctx context.Context, token string) (string, error) {
				return "unknown@example.com", domain.ErrUserNotFound
			},
		}
		i := NewInteractor(svc)

		_, err := i.VerifyEmailByToken(ctx, "token-for-unknown-user")
		if !errors.Is(err, domain.ErrUserNotFound) {
			t.Fatalf("expected %v, got %v", domain.ErrUserNotFound, err)
		}
	})

	t.Run("email already verified => returns error", func(t *testing.T) {
		svc := &fakeService{
			verifyEmailByTokenFn: func(ctx context.Context, token string) (string, error) {
				return "verified@example.com", domain.ErrEmailAlreadyVerified
			},
		}
		i := NewInteractor(svc)

		_, err := i.VerifyEmailByToken(ctx, "token-for-verified-email")
		if !errors.Is(err, domain.ErrEmailAlreadyVerified) {
			t.Fatalf("expected %v, got %v", domain.ErrEmailAlreadyVerified, err)
		}
	})
}

func TestInteractor_ResendVerificationEmail(t *testing.T) {
	ctx := context.Background()

	t.Run("user not found => returns error", func(t *testing.T) {
		svc := &fakeService{
			getUserByEmailFn: func(ctx context.Context, email string) (*gocloak.User, error) {
				return nil, errors.New("not found")
			},
		}
		i := NewInteractor(svc)

		err := i.ResendVerificationEmail(ctx, "nonexistent@example.com")
		if !errors.Is(err, domain.ErrUserNotFound) {
			t.Fatalf("expected %v, got %v", domain.ErrUserNotFound, err)
		}
	})

	t.Run("email already verified => returns error", func(t *testing.T) {
		verified := true
		svc := &fakeService{
			getUserByEmailFn: func(ctx context.Context, email string) (*gocloak.User, error) {
				id := "user-123"
				return &gocloak.User{ID: &id, EmailVerified: &verified}, nil
			},
		}
		i := NewInteractor(svc)

		err := i.ResendVerificationEmail(ctx, "verified@example.com")
		if !errors.Is(err, domain.ErrEmailAlreadyVerified) {
			t.Fatalf("expected %v, got %v", domain.ErrEmailAlreadyVerified, err)
		}
	})

	t.Run("success => sends verification email", func(t *testing.T) {
		notVerified := false
		sendCalled := false
		svc := &fakeService{
			getUserByEmailFn: func(ctx context.Context, email string) (*gocloak.User, error) {
				id := "user-123"
				return &gocloak.User{ID: &id, EmailVerified: &notVerified}, nil
			},
			sendVerificationEmailFn: func(ctx context.Context, userID string) error {
				sendCalled = true
				return nil
			},
		}
		i := NewInteractor(svc)

		err := i.ResendVerificationEmail(ctx, "test@example.com")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !sendCalled {
			t.Fatal("expected SendVerificationEmail to be called")
		}
	})

	t.Run("send email fails => returns error", func(t *testing.T) {
		notVerified := false
		sendErr := errors.New("keycloak error")
		svc := &fakeService{
			getUserByEmailFn: func(ctx context.Context, email string) (*gocloak.User, error) {
				id := "user-123"
				return &gocloak.User{ID: &id, EmailVerified: &notVerified}, nil
			},
			sendVerificationEmailFn: func(ctx context.Context, userID string) error {
				return sendErr
			},
		}
		i := NewInteractor(svc)

		err := i.ResendVerificationEmail(ctx, "test@example.com")
		if !errors.Is(err, sendErr) {
			t.Fatalf("expected %v, got %v", sendErr, err)
		}
	})
}

func TestInteractor_RefreshToken(t *testing.T) {
	ctx := context.Background()

	t.Run("success - returns TokenResponse", func(t *testing.T) {
		svc := &fakeService{
			refreshTokenFn: func(ctx context.Context, rt string) (*gocloak.JWT, error) {
				return &gocloak.JWT{
					AccessToken:  "new-access",
					RefreshToken: "new-refresh",
					ExpiresIn:    3600,
					TokenType:    "Bearer",
				}, nil
			},
		}

		i := NewInteractor(svc)
		result, err := i.RefreshToken(ctx, "valid-token")
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if result.AccessToken != "new-access" {
			t.Errorf("expected new-access, got %s", result.AccessToken)
		}
	})

	t.Run("service error - propagates", func(t *testing.T) {
		errorMsg := errors.New("refresh failed")
		svc := &fakeService{
			refreshTokenFn: func(ctx context.Context, rt string) (*gocloak.JWT, error) {
				return nil, errorMsg
			},
		}

		i := NewInteractor(svc)
		_, err := i.RefreshToken(ctx, "expired-token")
		if !errors.Is(err, errorMsg) {
			t.Fatalf("expected %v, got %v", errorMsg, err)
		}
	})
}
