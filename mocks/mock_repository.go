package mocks

import (
	"context"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/output"
	"github.com/stretchr/testify/mock"
)

// MockTx is a mock implementation of output.Tx
type MockTx struct {
	mock.Mock
}

func (m *MockTx) Commit() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockTx) Rollback() error {
	args := m.Called()
	return args.Error(0)
}

// MockRepository is a mock implementation of output.Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) BeginTx(ctx context.Context) (output.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(output.Tx), args.Error(1)
}

func (m *MockRepository) Save(ctx context.Context, tx output.Tx, employee domain.Employee) error {
	args := m.Called(ctx, tx, employee)
	return args.Error(0)
}

func (m *MockRepository) UpdateEmployee(ctx context.Context, tx output.Tx, employee domain.Employee) error {
	args := m.Called(ctx, tx, employee)
	return args.Error(0)
}

func (m *MockRepository) PatchEmployee(ctx context.Context, tx output.Tx, id string, keycloakUserID string) error {
	args := m.Called(ctx, tx, id, keycloakUserID)
	return args.Error(0)
}

func (m *MockRepository) DeleteEmployee(ctx context.Context, tx output.Tx, id string) error {
	args := m.Called(ctx, tx, id)
	return args.Error(0)
}

func (m *MockRepository) GetEmployeeByEmail(ctx context.Context, email string) (*domain.Employee, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Employee), args.Error(1)
}

func (m *MockRepository) GetEmployeeByID(ctx context.Context, id string) (*domain.Employee, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Employee), args.Error(1)
}

func (m *MockRepository) GetEmployeeByKeycloakID(ctx context.Context, keycloakUserID string) (*domain.Employee, error) {
	args := m.Called(ctx, keycloakUserID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Employee), args.Error(1)
}

// GetEmployeesByRole mock for HU47 - Virtual Entity pattern (Crew Member Types)
func (m *MockRepository) GetEmployeesByRole(ctx context.Context, role string) ([]domain.Employee, error) {
	args := m.Called(ctx, role)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Employee), args.Error(1)
}
