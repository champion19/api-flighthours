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

type AirportRepository interface {
	BeginTx(ctx context.Context) (Tx, error)
	GetAirportByID(ctx context.Context, id string) (*domain.Airport, error)
	ListAirports(ctx context.Context, filters map[string]interface{}) ([]domain.Airport, error)
	GetAirportsByType(ctx context.Context, airportType string) ([]domain.Airport, error)
	UpdateAirportStatus(ctx context.Context, tx Tx, id string, status bool) error
}

type LicensePlateRepository interface {
	BeginTx(ctx context.Context) (Tx, error)
	GetLicensePlateByID(ctx context.Context, id string) (*domain.LicensePlate, error)
	GetLicensePlateByPlate(ctx context.Context, plate string) (*domain.LicensePlate, error)
	ListLicensePlates(ctx context.Context, filters map[string]interface{}) ([]domain.LicensePlate, error)
	SaveLicensePlate(ctx context.Context, tx Tx, registration domain.LicensePlate) error
	UpdateLicensePlate(ctx context.Context, tx Tx, registration domain.LicensePlate) error
}

type AircraftModelRepository interface {
	BeginTx(ctx context.Context) (Tx, error)
	GetAircraftModelByID(ctx context.Context, id string) (*domain.AircraftModel, error)
	ListAircraftModels(ctx context.Context, filters map[string]interface{}) ([]domain.AircraftModel, error)
	GetAircraftModelsByFamily(ctx context.Context, family string) ([]domain.AircraftModel, error)
	UpdateAircraftModelStatus(ctx context.Context, tx Tx, id string, status bool) error
}

type RouteRepository interface {
	GetRouteByID(ctx context.Context, id string) (*domain.Route, error)
	ListRoutes(ctx context.Context, filters map[string]interface{}) ([]domain.Route, error)
}
type AirlineRouteRepository interface {
	BeginTx(ctx context.Context) (Tx, error)
	GetAirlineRouteByID(ctx context.Context, id string) (*domain.AirlineRoute, error)
	ListAirlineRoutes(ctx context.Context, filters map[string]interface{}) ([]domain.AirlineRoute, error)
	ListAirlineRoutesByAirlineID(ctx context.Context, airlineID string) ([]domain.AirlineRoute, error)
	UpdateAirlineRouteStatus(ctx context.Context, tx Tx, id string, status bool) error
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

type ManufacturerRepository interface {
	GetManufacturerByID(ctx context.Context, id string) (*domain.Manufacturer, error)
	ListManufacturers(ctx context.Context) ([]domain.Manufacturer, error)
}

type DailyLogbookDetailRepository interface {
	BeginTx(ctx context.Context) (Tx, error)
	GetDailyLogbookDetailByID(ctx context.Context, id string) (*domain.DailyLogbookDetail, error)
	ListDailyLogbookDetailsByLogbook(ctx context.Context, logbookID string) ([]domain.DailyLogbookDetail, error)
	ListDailyLogbookDetailsByEmployee(ctx context.Context, employeeID string) ([]domain.DailyLogbookDetail, error)
	SaveDailyLogbookDetail(ctx context.Context, tx Tx, detail domain.DailyLogbookDetail) error
	UpdateDailyLogbookDetail(ctx context.Context, tx Tx, detail domain.DailyLogbookDetail) error
	ExistsByUniqueKey(ctx context.Context, employeeID, flightRealDate, flightNumber, licensePlateID string) (bool, error)
}

type DailyLogbookRepository interface {
	BeginTx(ctx context.Context) (Tx, error)
	GetDailyLogbookByID(ctx context.Context, id string) (*domain.DailyLogbook, error)
	ListDailyLogbooksByEmployee(ctx context.Context, employeeID string, filters map[string]interface{}) ([]domain.DailyLogbook, error)
	SaveDailyLogbook(ctx context.Context, tx Tx, logbook domain.DailyLogbook) error
	UpdateDailyLogbook(ctx context.Context, tx Tx, logbook domain.DailyLogbook) error
	UpdateDailyLogbookStatus(ctx context.Context, tx Tx, id string, status bool) error
}
