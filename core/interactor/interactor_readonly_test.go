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
)

type fakeServiceForReadOnly struct {
	locateRes *dto.RegisterEmployee
	locateErr error
}

func (f *fakeServiceForReadOnly) BeginTx(context.Context) (output.Tx, error) {
	return &fakeTx{}, nil
}
func (f *fakeServiceForReadOnly) RegisterEmployee(context.Context, domain.Employee) (*dto.RegisterEmployee, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeServiceForReadOnly) GetEmployeeByEmail(context.Context, string) (*domain.Employee, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeServiceForReadOnly) GetEmployeeByID(context.Context, string) (*domain.Employee, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeServiceForReadOnly) LocateEmployee(ctx context.Context, id string) (*dto.RegisterEmployee, error) {
	return f.locateRes, f.locateErr
}
func (f *fakeServiceForReadOnly) CheckAndCleanInconsistentState(context.Context, string) error {
	return nil
}
func (f *fakeServiceForReadOnly) SaveEmployeeToDB(context.Context, output.Tx, domain.Employee) error {
	return nil
}
func (f *fakeServiceForReadOnly) UpdateEmployeeKeycloakID(context.Context, output.Tx, string, string) error {
	return nil
}
func (f *fakeServiceForReadOnly) CreateUserInKeycloak(context.Context, *domain.Employee) (string, error) {
	return "", errors.New("not implemented")
}
func (f *fakeServiceForReadOnly) SetUserPassword(context.Context, string, string) error {
	return nil
}
func (f *fakeServiceForReadOnly) AssignUserRole(context.Context, string, string) error {
	return nil
}
func (f *fakeServiceForReadOnly) SendVerificationEmail(context.Context, string) error {
	return nil
}
func (f *fakeServiceForReadOnly) SendPasswordResetEmail(context.Context, string) error {
	return nil
}
func (f *fakeServiceForReadOnly) Login(context.Context, string, string) (*gocloak.JWT, error) {
	return &gocloak.JWT{}, nil
}
func (f *fakeServiceForReadOnly) VerifyEmailByToken(context.Context, string) (string, error) {
	return "", nil
}
func (f *fakeServiceForReadOnly) GetUserByEmail(context.Context, string) (*gocloak.User, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeServiceForReadOnly) RollbackEmployee(context.Context, string) error {
	return nil
}
func (f *fakeServiceForReadOnly) RollbackKeycloakUser(context.Context, string) error {
	return nil
}
func (f *fakeServiceForReadOnly) UpdatePassword(context.Context, string, string) (string, error) {
	return "", nil
}
func (f *fakeServiceForReadOnly) UpdateEmployee(context.Context, domain.Employee, bool, string) error {
	return nil
}
func (f *fakeServiceForReadOnly) ChangePassword(context.Context, string, string, string) (string, error) {
	return "", nil
}
func (f *fakeServiceForReadOnly) DeleteEmployee(context.Context, string, string) error {
	return nil
}
func (f *fakeServiceForReadOnly) GetEmployeeByKeycloakID(context.Context, string) (*domain.Employee, error) {
	return nil, errors.New("not implemented")
}

// GetEmployeesByRole mock for HU47 - Virtual Entity pattern
func (f *fakeServiceForReadOnly) GetEmployeesByRole(context.Context, string) ([]domain.Employee, error) {
	return nil, errors.New("not implemented")
}

var _ input.Service = (*fakeServiceForReadOnly)(nil)

func TestInteractor_Locate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		employee := domain.Employee{ID: "123", Email: "test@example.com"}
		result := &dto.RegisterEmployee{Employee: employee, Message: "found"}

		svc := &fakeServiceForReadOnly{locateRes: result}
		inter := NewInteractor(svc, noopLogger{})

		res, err := inter.Locate(context.Background(), "123")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if res.Employee.ID != "123" {
			t.Fatalf("expected ID 123, got %s", res.Employee.ID)
		}
	})

	t.Run("not found", func(t *testing.T) {
		svc := &fakeServiceForReadOnly{locateErr: domain.ErrPersonNotFound}
		inter := NewInteractor(svc, noopLogger{})

		_, err := inter.Locate(context.Background(), "999")

		if err != domain.ErrPersonNotFound {
			t.Fatalf("expected ErrPersonNotFound, got %v", err)
		}
	})
}
