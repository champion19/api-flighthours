package output

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

type AuthClient interface {
	LoginUser(ctx context.Context, username, password string) (*gocloak.JWT, error)

	CreateUser(ctx context.Context, employee *domain.Employee) (string, error)
	GetUserByEmail(ctx context.Context, email string) (*gocloak.User, error)
	GetUserByID(ctx context.Context, userID string) (*gocloak.User, error)
	UpdateUser(ctx context.Context, user *gocloak.User) error
	DeleteUser(ctx context.Context, userID string) error
	SetPassword(ctx context.Context, userID string, password string, temporary bool) error

	AssignRole(ctx context.Context, userID string, roleName string) error
	RemoveRole(ctx context.Context, userID string, roleName string) error
	GetUserRoles(ctx context.Context, userID string) ([]*gocloak.Role, error)

	SendVerificationEmail(ctx context.Context, userID string) error
	VerifyEmail(ctx context.Context, userID string) error
	SendPasswordResetEmail(ctx context.Context, userID string) error
	ValidateActionToken(ctx context.Context, token string) (userID string, email string, err error)

	Logout(ctx context.Context, refreshToken string) error
	RefreshToken(ctx context.Context, refreshToken string) (*gocloak.JWT, error)
}
