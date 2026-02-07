package input

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
	"github.com/champion19/api-flighthours/core/interactor/dto"
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
)

type Service interface {
	BeginTx(ctx context.Context) (output.Tx, error)
	RegisterEmployee(ctx context.Context, employee domain.Employee) (*dto.RegisterEmployee, error)
	GetEmployeeByEmail(ctx context.Context, email string) (*domain.Employee, error)
	GetEmployeeByID(ctx context.Context, id string) (*domain.Employee, error)
	GetEmployeeByKeycloakID(ctx context.Context, keycloakUserID string) (*domain.Employee, error)
	LocateEmployee(ctx context.Context, id string) (*dto.RegisterEmployee, error)
	CheckAndCleanInconsistentState(ctx context.Context, email string) error
	GetEmployeesByRole(ctx context.Context, role string) ([]domain.Employee, error)
	SaveEmployeeToDB(ctx context.Context, tx output.Tx, employee domain.Employee) error
	UpdateEmployee(ctx context.Context, tx output.Tx, employee domain.Employee) error
	UpdateEmployeeKeycloakID(ctx context.Context, tx output.Tx, employeeID string, keycloakUserID string) error
	DeleteEmployee(ctx context.Context, employeeID string, keycloakUserID string) error
	CreateUserInKeycloak(ctx context.Context, employee *domain.Employee) (string, error)
	SetUserPassword(ctx context.Context, userID string, password string) error
	AssignUserRole(ctx context.Context, userID string, role string) error
	GetUserByEmail(ctx context.Context, email string) (*gocloak.User, error)
	SendVerificationEmail(ctx context.Context, userID string) error
	SendPasswordResetEmail(ctx context.Context, email string) error
	Login(ctx context.Context, email, password string) (*gocloak.JWT, error)
	VerifyEmailByToken(ctx context.Context, token string) (string, error)
	UpdatePassword(ctx context.Context, token, newPassword string) (string, error)
	ChangePassword(ctx context.Context, email, currentPassword, newPassword string) (string, error)
	RollbackEmployee(ctx context.Context, employeeID string) error
	RollbackKeycloakUser(ctx context.Context, KeycloakUserID string) error
}
type MessageService interface {
	BeginTx(ctx context.Context) (output.Tx, error)
	ValidateMessage(ctx context.Context, message domain.Message) error
	GetMessageByID(ctx context.Context, id string) (*domain.Message, error)
	GetMessageByCode(ctx context.Context, code string) (*domain.Message, error)
	ListMessages(ctx context.Context, filters map[string]interface{}) ([]domain.Message, error)
	ListActiveMessages(ctx context.Context) ([]domain.Message, error)
	SaveMessageToDB(ctx context.Context, tx output.Tx, message domain.Message) error
	UpdateMessageInDB(ctx context.Context, tx output.Tx, message domain.Message) error
	DeleteMessageFromDB(ctx context.Context, tx output.Tx, id string) error
}

type AirlineService interface {
	BeginTx(ctx context.Context) (output.Tx, error)
	GetAirlineByID(ctx context.Context, id string) (*domain.Airline, error)
	ListAirlines(ctx context.Context, filters map[string]interface{}) ([]domain.Airline, error)
	UpdateAirlineStatus(ctx context.Context, id string, status bool) error
	ActivateAirline(ctx context.Context, id string) error
	DeactivateAirline(ctx context.Context, id string) error
}

type AirportService interface {
	BeginTx(ctx context.Context) (output.Tx, error)
	GetAirportByID(ctx context.Context, id string) (*domain.Airport, error)
	ListAirports(ctx context.Context, filters map[string]interface{}) ([]domain.Airport, error)
	GetAirportsByType(ctx context.Context, airportType string) ([]domain.Airport, error)
	UpdateAirportStatus(ctx context.Context, id string, status bool) error
	ActivateAirport(ctx context.Context, id string) error
	DeactivateAirport(ctx context.Context, id string) error
}

type AirlineEmployeeService interface {
	BeginTx(ctx context.Context) (output.Tx, error)
	GetAirlineEmployeeByID(ctx context.Context, id string) (*domain.AirlineEmployee, error)
	AddAirlineEmployee(ctx context.Context, employee domain.AirlineEmployee) error
	UpdateAirlineEmployee(ctx context.Context, employee domain.AirlineEmployee) error
	ActivateAirlineEmployee(ctx context.Context, id string) error
	DeactivateAirlineEmployee(ctx context.Context, id string) error
}

type AircraftModelService interface {
	BeginTx(ctx context.Context) (output.Tx, error)
	GetAircraftModelByID(ctx context.Context, id string) (*domain.AircraftModel, error)
	ListAircraftModels(ctx context.Context, filters map[string]interface{}) ([]domain.AircraftModel, error)
	GetAircraftModelsByFamily(ctx context.Context, family string) ([]domain.AircraftModel, error)
	ActivateAircraftModel(ctx context.Context, id string) error
	DeactivateAircraftModel(ctx context.Context, id string) error
}


type RouteService interface {
	GetRouteByID(ctx context.Context, id string) (*domain.Route, error)
	ListRoutes(ctx context.Context, filters map[string]interface{}) ([]domain.Route, error)
}

type AirlineRouteService interface {
	BeginTx(ctx context.Context) (output.Tx, error)
	GetAirlineRouteByID(ctx context.Context, id string) (*domain.AirlineRoute, error)
	ListAirlineRoutes(ctx context.Context, filters map[string]interface{}) ([]domain.AirlineRoute, error)
	ListAirlineRoutesByAirlineID(ctx context.Context, airlineID string) ([]domain.AirlineRoute, error)
	ActivateAirlineRoute(ctx context.Context, id string) error
	DeactivateAirlineRoute(ctx context.Context, id string) error
}

type EngineService interface {
	GetEngineByID(ctx context.Context, id string) (*domain.Engine, error)
	ListEngines(ctx context.Context) ([]domain.Engine, error)
}

type ManufacturerService interface {
	GetManufacturerByID(ctx context.Context, id string) (*domain.Manufacturer, error)
	ListManufacturers(ctx context.Context) ([]domain.Manufacturer, error)
}
