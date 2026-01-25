package input

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
	"github.com/champion19/flighthours-api/core/interactor/dto"
	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/output"
)

type Service interface {
	BeginTx(ctx context.Context) (output.Tx, error)

	//employee-validaciones y consultas
	RegisterEmployee(ctx context.Context, employee domain.Employee) (*dto.RegisterEmployee, error)
	GetEmployeeByEmail(ctx context.Context, email string) (*domain.Employee, error)
	GetEmployeeByID(ctx context.Context, id string) (*domain.Employee, error)
	GetEmployeeByKeycloakID(ctx context.Context, keycloakUserID string) (*domain.Employee, error)
	LocateEmployee(ctx context.Context, id string) (*dto.RegisterEmployee, error)
	CheckAndCleanInconsistentState(ctx context.Context, email string) error
	GetEmployeesByRole(ctx context.Context, role string) ([]domain.Employee, error)

	//employee- operaciones transaccionales de BD
	SaveEmployeeToDB(ctx context.Context, tx output.Tx, employee domain.Employee) error
	UpdateEmployeeKeycloakID(ctx context.Context, tx output.Tx, employeeID string, keycloakUserID string) error
	// DeleteEmployee elimina un empleado de la BD y de Keycloak
	// Primero elimina de Keycloak, luego de la BD
	DeleteEmployee(ctx context.Context, employeeID string, keycloakUserID string) error

	//employee-operaciones de keycloak
	CreateUserInKeycloak(ctx context.Context, employee *domain.Employee) (string, error)
	SetUserPassword(ctx context.Context, userID string, password string) error
	AssignUserRole(ctx context.Context, userID string, role string) error
	GetUserByEmail(ctx context.Context, email string) (*gocloak.User, error)
	SendVerificationEmail(ctx context.Context, userID string) error
	SendPasswordResetEmail(ctx context.Context, email string) error
	Login(ctx context.Context, email, password string) (*gocloak.JWT, error)
	VerifyEmailByToken(ctx context.Context, token string) (string, error)
	UpdatePassword(ctx context.Context, token, newPassword string) (string, error)
	// ChangePassword changes the password for an authenticated user who knows their current password
	// Returns the email of the user whose password was changed
	ChangePassword(ctx context.Context, email, currentPassword, newPassword string) (string, error)

	//employee- compensaciones (rollback)
	RollbackEmployee(ctx context.Context, employeeID string) error
	RollbackKeycloakUser(ctx context.Context, KeycloakUserID string) error
}
type MessageService interface {
	// Transacciones
	BeginTx(ctx context.Context) (output.Tx, error)

	// Messages - Validaciones y operaciones
	ValidateMessage(ctx context.Context, message domain.Message) error
	GetMessageByID(ctx context.Context, id string) (*domain.Message, error)
	GetMessageByCode(ctx context.Context, code string) (*domain.Message, error)
	ListMessages(ctx context.Context, filters map[string]interface{}) ([]domain.Message, error)
	ListActiveMessages(ctx context.Context) ([]domain.Message, error)

	// Messages - Operaciones transaccionales de BD
	SaveMessageToDB(ctx context.Context, tx output.Tx, message domain.Message) error
	UpdateMessageInDB(ctx context.Context, tx output.Tx, message domain.Message) error
	DeleteMessageFromDB(ctx context.Context, tx output.Tx, id string) error
}
