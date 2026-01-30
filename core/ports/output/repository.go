package output

import (
	"context"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

type Tx interface {
	Commit() error
	Rollback() error
}

type Repository interface {
	BeginTx(ctx context.Context) (Tx, error)

	Save(ctx context.Context, tx Tx, employee domain.Employee) error
	UpdateEmployee(ctx context.Context, tx Tx, employee domain.Employee) error
	PatchEmployee(ctx context.Context, tx Tx, id string, keycloakUserID string) error
	DeleteEmployee(ctx context.Context, tx Tx, id string) error
	GetEmployeeByEmail(ctx context.Context, email string) (*domain.Employee, error)
	GetEmployeeByID(ctx context.Context, id string) (*domain.Employee, error)
	GetEmployeeByKeycloakID(ctx context.Context, keycloakUserID string) (*domain.Employee, error)
	GetEmployeesByRole(ctx context.Context, role string) ([]domain.Employee, error)
}

type MessageRepository interface {
	BeginTx(ctx context.Context) (Tx, error)

	SaveMessage(ctx context.Context, tx Tx, message domain.Message) error
	UpdateMessage(ctx context.Context, tx Tx, message domain.Message) error
	DeleteMessage(ctx context.Context, tx Tx, id string) error
	GetAllActive(ctx context.Context) ([]domain.Message, error)
	GetByID(ctx context.Context, id string) (*domain.Message, error)
	GetByCode(ctx context.Context, code string) (*domain.Message, error)
	GetByType(ctx context.Context, msgType string) ([]domain.Message, error)
	GetByModule(ctx context.Context, module string) ([]domain.Message, error)
}

type AirlineRepository interface {
	BeginTx(ctx context.Context) (Tx, error)
	GetAirlineByID(ctx context.Context, id string) (*domain.Airline, error)
	ListAirlines(ctx context.Context, filters map[string]interface{}) ([]domain.Airline, error)
	UpdateAirlineStatus(ctx context.Context, tx Tx, id string, active bool) error
}

type AirlineEmployeeRepository interface {
	BeginTx(ctx context.Context) (Tx, error)

	GetAirlineEmployeeByID(ctx context.Context, id string) (*domain.AirlineEmployee, error)
	AddAirlineEmployee(ctx context.Context, tx Tx, employee domain.AirlineEmployee) error
	UpdateAirlineEmployee(ctx context.Context, tx Tx, employee domain.AirlineEmployee) error
	UpdateAirlineEmployeeStatus(ctx context.Context, tx Tx, id string, active bool) error
}

type EngineRepository interface {
	GetEngineByID(ctx context.Context, id string) (*domain.Engine, error)
	ListEngines(ctx context.Context) ([]domain.Engine, error)
}
