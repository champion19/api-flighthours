package input

import (
	"context"

	"github.com/champion19/api-flighthours/core/interactor/dto"
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

// EmployeeInteractor defines the interface for employee business operations
type EmployeeInteractor interface {
	// RegisterEmployee registers a new employee in the system
	RegisterEmployee(ctx context.Context, employee domain.Employee) (*dto.RegisterEmployee, error)

	// Login authenticates an employee and returns tokens
	Login(ctx context.Context, email, password string) (*dto.TokenResponse, error)

	// RefreshToken refreshes an expired access token using a refresh token
	RefreshToken(ctx context.Context, refreshToken string) (*dto.TokenResponse, error)

	// VerifyEmailByToken verifies an email using a token
	VerifyEmailByToken(ctx context.Context, token string) (string, error)

	// RequestPasswordReset sends a password reset email
	RequestPasswordReset(ctx context.Context, email string) error

	// UpdatePassword updates password using a reset token
	UpdatePassword(ctx context.Context, token, newPassword, confirmPassword string) (string, error)

	// ChangePassword changes password for authenticated user
	ChangePassword(ctx context.Context, email, currentPassword, newPassword, confirmPassword string) (string, error)

	// ResendVerificationEmail resends the verification email
	ResendVerificationEmail(ctx context.Context, email string) error

	// DeleteEmployee deletes an employee
	DeleteEmployee(ctx context.Context, employeeID string) error

	// UpdateEmployee updates employee information
	UpdateEmployee(ctx context.Context, employee domain.Employee) (*dto.UpdateEmployee, error)

	// Locate finds an employee by ID
	Locate(ctx context.Context, id string) (*dto.RegisterEmployee, error)
}
