package services

import (
	"context"
	"errors"
	"testing"

	"github.com/Nerzal/gocloak/v13"
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
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

type fakeRepo struct {
	getByEmailFn func(ctx context.Context, email string) (*domain.Employee, error)
}

func (f fakeRepo) BeginTx(ctx context.Context) (output.Tx, error) {
	return nil, errors.New("not implemented")
}
func (f fakeRepo) Save(context.Context, output.Tx, domain.Employee) error {
	return errors.New("not implemented")
}
func (f fakeRepo) UpdateEmployee(context.Context, output.Tx, domain.Employee) error {
	return errors.New("not implemented")
}
func (f fakeRepo) PatchEmployee(context.Context, output.Tx, string, string) error {
	return errors.New("not implemented")
}
func (f fakeRepo) DeleteEmployee(context.Context, output.Tx, string) error {
	return errors.New("not implemented")
}
func (f fakeRepo) GetEmployeeByEmail(ctx context.Context, email string) (*domain.Employee, error) {
	return f.getByEmailFn(ctx, email)
}
func (f fakeRepo) GetEmployeeByID(context.Context, string) (*domain.Employee, error) {
	return nil, errors.New("not implemented")
}
func (f fakeRepo) GetEmployeeByKeycloakID(context.Context, string) (*domain.Employee, error) {
	return nil, errors.New("not implemented")
}

// GetEmployeesByRole mock
func (f fakeRepo) GetEmployeesByRole(context.Context, string) ([]domain.Employee, error) {
	return nil, errors.New("not implemented")
}

type fakeKeycloak struct {
	getUserByEmailFn func(ctx context.Context, email string) (*gocloak.User, error)
	createUserFn     func(ctx context.Context, employee *domain.Employee) (string, error)
}

func (f fakeKeycloak) LoginUser(context.Context, string, string) (*gocloak.JWT, error) {
	return nil, errors.New("not implemented")
}
func (f fakeKeycloak) CreateUser(ctx context.Context, employee *domain.Employee) (string, error) {
	return f.createUserFn(ctx, employee)
}
func (f fakeKeycloak) GetUserByEmail(ctx context.Context, email string) (*gocloak.User, error) {
	return f.getUserByEmailFn(ctx, email)
}
func (fakeKeycloak) GetUserByID(context.Context, string) (*gocloak.User, error) {
	return nil, errors.New("not implemented")
}
func (fakeKeycloak) UpdateUser(context.Context, *gocloak.User) error {
	return errors.New("not implemented")
}
func (fakeKeycloak) DeleteUser(context.Context, string) error { return errors.New("not implemented") }
func (fakeKeycloak) SetPassword(context.Context, string, string, bool) error {
	return errors.New("not implemented")
}
func (fakeKeycloak) AssignRole(context.Context, string, string) error {
	return errors.New("not implemented")
}
func (fakeKeycloak) RemoveRole(context.Context, string, string) error {
	return errors.New("not implemented")
}
func (fakeKeycloak) GetUserRoles(context.Context, string) ([]*gocloak.Role, error) {
	return nil, errors.New("not implemented")
}
func (fakeKeycloak) SendVerificationEmail(context.Context, string) error {
	return errors.New("not implemented")
}
func (fakeKeycloak) SendPasswordResetEmail(context.Context, string) error {
	return errors.New("not implemented")
}
func (fakeKeycloak) VerifyEmail(context.Context, string) error { return errors.New("not implemented") }
func (fakeKeycloak) Logout(context.Context, string) error      { return errors.New("not implemented") }
func (fakeKeycloak) RefreshToken(context.Context, string) (*gocloak.JWT, error) {
	return nil, errors.New("not implemented")
}
func (fakeKeycloak) ValidateActionToken(context.Context, string) (string, string, error) {
	return "", "", errors.New("not implemented")
}

func TestService_RegisterEmployee(t *testing.T) {
	ctx := context.Background()

	mkEmployee := func() domain.Employee {
		return domain.Employee{Email: "a@b.com"}
	}

	t.Run("db connection error => ErrDatabaseUnavailable", func(t *testing.T) {
		svc := NewService(
			fakeRepo{getByEmailFn: func(context.Context, string) (*domain.Employee, error) {
				return nil, errors.New("connection refused")
			}},
			fakeKeycloak{getUserByEmailFn: func(context.Context, string) (*gocloak.User, error) {
				return nil, errors.New("should not be called")
			}},
			noopLogger{},
		)

		_, err := svc.RegisterEmployee(ctx, mkEmployee())
		if !errors.Is(err, domain.ErrDatabaseUnavailable) {
			t.Fatalf("expected %v, got %v", domain.ErrDatabaseUnavailable, err)
		}
	})

	t.Run("keycloak connection error => ErrKeycloakUnavailable", func(t *testing.T) {
		svc := NewService(
			fakeRepo{getByEmailFn: func(context.Context, string) (*domain.Employee, error) {
				return nil, errors.New("record not found")
			}},
			fakeKeycloak{getUserByEmailFn: func(context.Context, string) (*gocloak.User, error) {
				return nil, errors.New("context deadline exceeded")
			}},
			noopLogger{},
		)

		_, err := svc.RegisterEmployee(ctx, mkEmployee())
		if !errors.Is(err, domain.ErrKeycloakUnavailable) {
			t.Fatalf("expected %v, got %v", domain.ErrKeycloakUnavailable, err)
		}
	})

	t.Run("exists in db and keycloak => ErrDuplicateUser", func(t *testing.T) {
		kcID := "kc1"
		svc := NewService(
			fakeRepo{getByEmailFn: func(context.Context, string) (*domain.Employee, error) {
				return &domain.Employee{ID: "db1", Email: "a@b.com"}, nil
			}},
			fakeKeycloak{getUserByEmailFn: func(context.Context, string) (*gocloak.User, error) {
				return &gocloak.User{ID: &kcID}, nil
			}},
			noopLogger{},
		)

		_, err := svc.RegisterEmployee(ctx, mkEmployee())
		if !errors.Is(err, domain.ErrDuplicateUser) {
			t.Fatalf("expected %v, got %v", domain.ErrDuplicateUser, err)
		}
	})

	t.Run("exists only in db => ErrIncompleteRegistration", func(t *testing.T) {
		svc := NewService(
			fakeRepo{getByEmailFn: func(context.Context, string) (*domain.Employee, error) {
				return &domain.Employee{ID: "db1", Email: "a@b.com"}, nil
			}},
			fakeKeycloak{getUserByEmailFn: func(context.Context, string) (*gocloak.User, error) {
				return nil, errors.New("404")
			}},
			noopLogger{},
		)

		_, err := svc.RegisterEmployee(ctx, mkEmployee())
		if !errors.Is(err, domain.ErrIncompleteRegistration) {
			t.Fatalf("expected %v, got %v", domain.ErrIncompleteRegistration, err)
		}
	})

	t.Run("exists only in keycloak => ErrIncompleteRegistration", func(t *testing.T) {
		kcID := "kc1"
		svc := NewService(
			fakeRepo{getByEmailFn: func(context.Context, string) (*domain.Employee, error) {
				return nil, errors.New("not found")
			}},
			fakeKeycloak{getUserByEmailFn: func(context.Context, string) (*gocloak.User, error) {
				return &gocloak.User{ID: &kcID}, nil
			}},
			noopLogger{},
		)

		_, err := svc.RegisterEmployee(ctx, mkEmployee())
		if !errors.Is(err, domain.ErrIncompleteRegistration) {
			t.Fatalf("expected %v, got %v", domain.ErrIncompleteRegistration, err)
		}
	})

	t.Run("not in db nor keycloak => success", func(t *testing.T) {
		svc := NewService(
			fakeRepo{getByEmailFn: func(context.Context, string) (*domain.Employee, error) {
				return nil, errors.New("not found")
			}},
			fakeKeycloak{getUserByEmailFn: func(context.Context, string) (*gocloak.User, error) {
				return nil, errors.New("not found")
			}},
			noopLogger{},
		)

		res, err := svc.RegisterEmployee(ctx, mkEmployee())
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if res == nil {
			t.Fatalf("expected non-nil result")
		}
		if res.Employee.Email != "a@b.com" {
			t.Fatalf("unexpected employee in result: %+v", res.Employee)
		}
	})
}

func TestService_CreateUserInKeycloak(t *testing.T) {
	ctx := context.Background()
	employee := &domain.Employee{Email: "a@b.com"}

	t.Run("connection error => ErrKeycloakUnavailable", func(t *testing.T) {
		svc := NewService(
			fakeRepo{getByEmailFn: func(context.Context, string) (*domain.Employee, error) {
				return nil, nil
			}},
			fakeKeycloak{createUserFn: func(context.Context, *domain.Employee) (string, error) {
				return "", errors.New("connect: connection refused")
			}},
			noopLogger{},
		)

		_, err := svc.CreateUserInKeycloak(ctx, employee)
		if !errors.Is(err, domain.ErrKeycloakUnavailable) {
			t.Fatalf("expected %v, got %v", domain.ErrKeycloakUnavailable, err)
		}
	})

	t.Run("generic error => ErrKeycloakUserCreationFailed", func(t *testing.T) {
		svc := NewService(
			fakeRepo{getByEmailFn: func(context.Context, string) (*domain.Employee, error) {
				return nil, nil
			}},
			fakeKeycloak{createUserFn: func(context.Context, *domain.Employee) (string, error) {
				return "", errors.New("some keycloak error")
			}},
			noopLogger{},
		)

		_, err := svc.CreateUserInKeycloak(ctx, employee)
		if !errors.Is(err, domain.ErrKeycloakUserCreationFailed) {
			t.Fatalf("expected %v, got %v", domain.ErrKeycloakUserCreationFailed, err)
		}
	})

	t.Run("success", func(t *testing.T) {
		svc := NewService(
			fakeRepo{getByEmailFn: func(context.Context, string) (*domain.Employee, error) {
				return nil, nil
			}},
			fakeKeycloak{createUserFn: func(context.Context, *domain.Employee) (string, error) {
				return "kc1", nil
			}},
			noopLogger{},
		)

		id, err := svc.CreateUserInKeycloak(ctx, employee)
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if id != "kc1" {
			t.Fatalf("expected kc1, got %s", id)
		}
	})
}
