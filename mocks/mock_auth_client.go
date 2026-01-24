package mocks

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/stretchr/testify/mock"
)

// MockAuthClient is a mock implementation of output.AuthClient (Keycloak)
type MockAuthClient struct {
	mock.Mock
}

// Autenticación

func (m *MockAuthClient) LoginUser(ctx context.Context, username, password string) (*gocloak.JWT, error) {
	args := m.Called(ctx, username, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*gocloak.JWT), args.Error(1)
}

// Gestión de usuarios

func (m *MockAuthClient) CreateUser(ctx context.Context, employee *domain.Employee) (string, error) {
	args := m.Called(ctx, employee)
	return args.String(0), args.Error(1)
}

func (m *MockAuthClient) GetUserByEmail(ctx context.Context, email string) (*gocloak.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*gocloak.User), args.Error(1)
}

func (m *MockAuthClient) GetUserByID(ctx context.Context, userID string) (*gocloak.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*gocloak.User), args.Error(1)
}

func (m *MockAuthClient) UpdateUser(ctx context.Context, user *gocloak.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockAuthClient) DeleteUser(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockAuthClient) SetPassword(ctx context.Context, userID string, password string, temporary bool) error {
	args := m.Called(ctx, userID, password, temporary)
	return args.Error(0)
}

// Roles

func (m *MockAuthClient) AssignRole(ctx context.Context, userID string, roleName string) error {
	args := m.Called(ctx, userID, roleName)
	return args.Error(0)
}

func (m *MockAuthClient) RemoveRole(ctx context.Context, userID string, roleName string) error {
	args := m.Called(ctx, userID, roleName)
	return args.Error(0)
}

func (m *MockAuthClient) GetUserRoles(ctx context.Context, userID string) ([]*gocloak.Role, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*gocloak.Role), args.Error(1)
}

// Verificación

func (m *MockAuthClient) SendVerificationEmail(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockAuthClient) SendPasswordResetEmail(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockAuthClient) VerifyEmail(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// Sesiones

func (m *MockAuthClient) Logout(ctx context.Context, refreshToken string) error {
	args := m.Called(ctx, refreshToken)
	return args.Error(0)
}

func (m *MockAuthClient) RefreshToken(ctx context.Context, refreshToken string) (*gocloak.JWT, error) {
	args := m.Called(ctx, refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*gocloak.JWT), args.Error(1)
}

// Token Validation

func (m *MockAuthClient) ValidateActionToken(ctx context.Context, token string) (string, string, error) {
	args := m.Called(ctx, token)
	return args.String(0), args.String(1), args.Error(2)
}
