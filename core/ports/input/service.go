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
	RefreshToken(ctx context.Context, refreshToken string) (*gocloak.JWT, error)
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
	ActivateAirlineTx(ctx context.Context, tx output.Tx, id string) error
	DeactivateAirlineTx(ctx context.Context, tx output.Tx, id string) error
}

type AirportService interface {
	BeginTx(ctx context.Context) (output.Tx, error)
	GetAirportByID(ctx context.Context, id string) (*domain.Airport, error)
	ListAirports(ctx context.Context, filters map[string]interface{}) ([]domain.Airport, error)
	GetAirportsByType(ctx context.Context, airportType string) ([]domain.Airport, error)
	ActivateAirportTx(ctx context.Context, tx output.Tx, id string) error
	DeactivateAirportTx(ctx context.Context, tx output.Tx, id string) error
}

type AirlineEmployeeService interface {
	BeginTx(ctx context.Context) (output.Tx, error)
	GetAirlineEmployeeByID(ctx context.Context, id string) (*domain.AirlineEmployee, error)
	AddAirlineEmployeeTx(ctx context.Context, tx output.Tx, employee domain.AirlineEmployee) error
	UpdateAirlineEmployeeTx(ctx context.Context, tx output.Tx, employee domain.AirlineEmployee) error
	ActivateAirlineEmployeeTx(ctx context.Context, tx output.Tx, id string) error
	DeactivateAirlineEmployeeTx(ctx context.Context, tx output.Tx, id string) error
}

type TailNumberService interface {
	BeginTx(ctx context.Context) (output.Tx, error)
	GetTailNumberByID(ctx context.Context, id string) (*domain.TailNumber, error)
	GetTailNumberByPlate(ctx context.Context, plate string) (*domain.TailNumber, error)
	ListTailNumbers(ctx context.Context, filters map[string]interface{}) ([]domain.TailNumber, error)
	CreateTailNumberTx(ctx context.Context, tx output.Tx, registration domain.TailNumber) error
	UpdateTailNumberTx(ctx context.Context, tx output.Tx, registration domain.TailNumber) error
}

type AircraftModelService interface {
	BeginTx(ctx context.Context) (output.Tx, error)
	GetAircraftModelByID(ctx context.Context, id string) (*domain.AircraftModel, error)
	ListAircraftModels(ctx context.Context, filters map[string]interface{}) ([]domain.AircraftModel, error)
	GetAircraftModelsByFamily(ctx context.Context, family string) ([]domain.AircraftModel, error)
	ActivateAircraftModelTx(ctx context.Context, tx output.Tx, id string) error
	DeactivateAircraftModelTx(ctx context.Context, tx output.Tx, id string) error
}

type RouteService interface {
	GetRouteByID(ctx context.Context, id string) (*domain.Route, error)
	ListRoutes(ctx context.Context, filters map[string]interface{}) ([]domain.Route, error)
	GetRouteByAirports(ctx context.Context, originAirportID, destinationAirportID string) (*domain.Route, error)
}

type AirlineRouteService interface {
	BeginTx(ctx context.Context) (output.Tx, error)
	GetAirlineRouteByID(ctx context.Context, id string) (*domain.AirlineRoute, error)
	ListAirlineRoutes(ctx context.Context, filters map[string]interface{}) ([]domain.AirlineRoute, error)
	ListAirlineRoutesByAirlineID(ctx context.Context, airlineID string) ([]domain.AirlineRoute, error)
	ActivateAirlineRouteTx(ctx context.Context, tx output.Tx, id string) error
	DeactivateAirlineRouteTx(ctx context.Context, tx output.Tx, id string) error
	GetAirlineRouteByRouteAndAirline(ctx context.Context, routeID, airlineID string) (*domain.AirlineRoute, error)
	SaveAirlineRouteTx(ctx context.Context, tx output.Tx, airlineRoute domain.AirlineRoute) error
}

type DailyLogbookDetailService interface {
	BeginTx(ctx context.Context) (output.Tx, error)
	GetDailyLogbookDetailByID(ctx context.Context, id string) (*domain.DailyLogbookDetail, error)
	ListDailyLogbookDetailsByLogbook(ctx context.Context, logbookID string) ([]domain.DailyLogbookDetail, error)
	ListDailyLogbookDetailsByEmployee(ctx context.Context, employeeID string) ([]domain.DailyLogbookDetail, error)
	CreateDailyLogbookDetailTx(ctx context.Context, tx output.Tx, detail domain.DailyLogbookDetail) error
	UpdateDailyLogbookDetailTx(ctx context.Context, tx output.Tx, detail domain.DailyLogbookDetail) error
	ValidateTimeSequence(outTime, takeoffTime, landingTime, inTime string) error
	ExistsByUniqueKey(ctx context.Context, employeeLogbookID, flightRealDate, flightNumber, tailNumberID string) (bool, error)
	DeleteDailyLogbookDetailTx(ctx context.Context, tx output.Tx, id string) error
}

type DailyLogbookService interface {
	BeginTx(ctx context.Context) (output.Tx, error)
	GetDailyLogbookByID(ctx context.Context, id string) (*domain.DailyLogbook, error)
	ListDailyLogbooksByEmployee(ctx context.Context, employeeID string, filters map[string]interface{}) ([]domain.DailyLogbook, error)
	CreateDailyLogbookTx(ctx context.Context, tx output.Tx, logbook domain.DailyLogbook) error
	UpdateDailyLogbookTx(ctx context.Context, tx output.Tx, logbook domain.DailyLogbook) error
	ActivateDailyLogbookTx(ctx context.Context, tx output.Tx, id string) error
	DeactivateDailyLogbookTx(ctx context.Context, tx output.Tx, id string) error
	DeleteDailyLogbookTx(ctx context.Context, tx output.Tx, id string) error
}

type EngineService interface {
	GetEngineByID(ctx context.Context, id string) (*domain.Engine, error)
	ListEngines(ctx context.Context) ([]domain.Engine, error)
}

type ManufacturerService interface {
	GetManufacturerByID(ctx context.Context, id string) (*domain.Manufacturer, error)
	ListManufacturers(ctx context.Context) ([]domain.Manufacturer, error)
}

type FlightSummaryService interface {
	GetFlightHoursSummary(ctx context.Context, employeeID, startDate, endDate string) (*domain.FlightHoursSummary, error)
	GetRecentFlights(ctx context.Context, employeeID string, limit int) ([]domain.DailyLogbookDetail, error)
	GetLandingCount(ctx context.Context, employeeID, startDate, endDate string) (int, error)
	GetDailyFlightSeconds(ctx context.Context, employeeID, date string) (int, error)
	CalculatePeriodDates(period, referenceDate string) (startDate, endDate string, err error)
	BuildFlightAlerts(ctx context.Context, employeeID string) ([]domain.FlightAlert, error)
}

type EmployeeFlightSummaryService interface {
	AccumulateFlightHours(ctx context.Context, tx output.Tx, employeeID string, detail domain.DailyLogbookDetail, isDeletion bool) error
	GetSummariesByEmployee(ctx context.Context, employeeID, periodType string) ([]domain.EmployeeFlightSummary, error)
	GetCurrentPeriodSummary(ctx context.Context, employeeID, periodType string, year, number int) (*domain.EmployeeFlightSummary, error)
	GetAllSummaries(ctx context.Context, employeeID string) ([]domain.EmployeeFlightSummary, error)
}
