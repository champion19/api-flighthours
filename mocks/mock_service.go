package mocks

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
	"github.com/champion19/flighthours-api/core/interactor/dto"
	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/output"
	"github.com/stretchr/testify/mock"
)

// MockService is a mock implementation of input.Service
type MockService struct {
	mock.Mock
}

// Transacciones

func (m *MockService) BeginTx(ctx context.Context) (output.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(output.Tx), args.Error(1)
}

// Person - Validaciones y consultas

func (m *MockService) RegisterEmployee(ctx context.Context, employee domain.Employee) (*dto.RegisterEmployee, error) {
	args := m.Called(ctx, employee)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.RegisterEmployee), args.Error(1)
}

func (m *MockService) GetEmployeeByEmail(ctx context.Context, email string) (*domain.Employee, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Employee), args.Error(1)
}

func (m *MockService) GetEmployeeByID(ctx context.Context, id string) (*domain.Employee, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Employee), args.Error(1)
}

func (m *MockService) LocateEmployee(ctx context.Context, id string) (*dto.RegisterEmployee, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.RegisterEmployee), args.Error(1)
}

func (m *MockService) CheckAndCleanInconsistentState(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

// Person - Operaciones transaccionales de BD

func (m *MockService) SaveEmployeeToDB(ctx context.Context, tx output.Tx, employee domain.Employee) error {
	args := m.Called(ctx, tx, employee)
	return args.Error(0)
}

func (m *MockService) UpdateEmployeeKeycloakID(ctx context.Context, tx output.Tx, employeeID string, keycloakUserID string) error {
	args := m.Called(ctx, tx, employeeID, keycloakUserID)
	return args.Error(0)
}

// Person - Operaciones de Keycloak

func (m *MockService) CreateUserInKeycloak(ctx context.Context, employee *domain.Employee) (string, error) {
	args := m.Called(ctx, employee)
	return args.String(0), args.Error(1)
}

func (m *MockService) SetUserPassword(ctx context.Context, userID string, password string) error {
	args := m.Called(ctx, userID, password)
	return args.Error(0)
}

func (m *MockService) AssignUserRole(ctx context.Context, userID string, role string) error {
	args := m.Called(ctx, userID, role)
	return args.Error(0)
}

// Person - Compensaciones (rollback)

func (m *MockService) RollbackEmployee(ctx context.Context, employeeID string) error {
	args := m.Called(ctx, employeeID)
	return args.Error(0)
}

func (m *MockService) RollbackKeycloakUser(ctx context.Context, keycloakUserID string) error {
	args := m.Called(ctx, keycloakUserID)
	return args.Error(0)
}

func (m *MockService) GetUserByEmail(ctx context.Context, email string) (*gocloak.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*gocloak.User), args.Error(1)
}

func (m *MockService) SendVerificationEmail(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockService) SendPasswordResetEmail(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

func (m *MockService) Login(ctx context.Context, email, password string) (*gocloak.JWT, error) {
	args := m.Called(ctx, email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*gocloak.JWT), args.Error(1)
}

func (m *MockService) VerifyEmailByToken(ctx context.Context, token string) (string, error) {
	args := m.Called(ctx, token)
	return args.String(0), args.Error(1)
}
