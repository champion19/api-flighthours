package dependency

import (
	"context"
	"database/sql"
	"time"

	"github.com/champion19/api-flighthours/config"
	"github.com/champion19/api-flighthours/core/interactor"
	"github.com/champion19/api-flighthours/core/interactor/services"
	"github.com/champion19/api-flighthours/core/ports/input"
	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/middleware"
	messagingCache "github.com/champion19/api-flighthours/platform/cache/messaging"
	mysql "github.com/champion19/api-flighthours/platform/databases/mysql"
	aircraftModelRepo "github.com/champion19/api-flighthours/platform/databases/repositories/aircraft_model"
	airlineRepo "github.com/champion19/api-flighthours/platform/databases/repositories/airline"
	airlineEmployeeRepo "github.com/champion19/api-flighthours/platform/databases/repositories/airline_employee"
	airlineRouteRepo "github.com/champion19/api-flighthours/platform/databases/repositories/airline_route"
	airportRepo "github.com/champion19/api-flighthours/platform/databases/repositories/airport"
	dailyLogbookRepo "github.com/champion19/api-flighthours/platform/databases/repositories/daily_logbook"
	dailyLogbookDetailRepo "github.com/champion19/api-flighthours/platform/databases/repositories/daily_logbook_detail"
	repo "github.com/champion19/api-flighthours/platform/databases/repositories/employee"
	employeeFlightSummaryRepo "github.com/champion19/api-flighthours/platform/databases/repositories/employee_flight_summary"
	engineRepo "github.com/champion19/api-flighthours/platform/databases/repositories/engine"
	flightSummaryRepo "github.com/champion19/api-flighthours/platform/databases/repositories/flight_summary"
	manufacturerRepo "github.com/champion19/api-flighthours/platform/databases/repositories/manufacturer"
	messageRepo "github.com/champion19/api-flighthours/platform/databases/repositories/message"
	routeRepo "github.com/champion19/api-flighthours/platform/databases/repositories/route"
	tailNumberRepo "github.com/champion19/api-flighthours/platform/databases/repositories/tail_number"
	"github.com/champion19/api-flighthours/platform/identity_provider/keycloak"
	"github.com/champion19/api-flighthours/platform/jwt"
	"github.com/champion19/api-flighthours/platform/logger"
	"github.com/champion19/api-flighthours/tools/idencoder"
)

var log logger.Logger = logger.NewSlogLogger()

type Dependencies struct {
	EmployeeService              input.Service
	EmployeeRepo                 output.Repository
	Interactor                   *interactor.Interactor
	KeycloakClient               output.AuthClient
	Config                       *config.Config
	IDEncoder                    *idencoder.HashidsEncoder
	ResponseHandler              *middleware.ResponseHandler
	MessagingCache               *messagingCache.MessageCache
	MessageInteractor            *interactor.MessageInteractor
	JWTValidator                 output.TokenValidator
	AirlineInteractor            *interactor.AirlineInteractor
	AirlineEmployeeInteractor    *interactor.AirlineEmployeeInteractor
	EngineInteractor             *interactor.EngineInteractor
	RouteInteractor              *interactor.RouteInteractor
	AirlineRouteInteractor       *interactor.AirlineRouteInteractor
	AirportInteractor            *interactor.AirportInteractor
	ManufacturerInteractor       *interactor.ManufacturerInteractor
	AircraftModelInteractor      *interactor.AircraftModelInteractor
	TailNumberInteractor         *interactor.TailNumberInteractor
	DailyLogbookDetailInteractor *interactor.DailyLogbookDetailInteractor
	DailyLogbookInteractor       *interactor.DailyLogbookInteractor
	FlightSummaryInteractor      *interactor.FlightSummaryInteractor
}

func Init() (*Dependencies, error) {
	log := logger.NewSlogLogger()
	log.Info(logger.LogAppStarting)

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error(logger.LogAppConfigError, "error", err)
		return nil, err
	}
	log.Info(logger.LogAppConfigLoaded)

	middleware.PrometheusInit()
	log.Success(logger.LogPrometheusInitOK)

	db, err := mysql.GetDB(cfg.Database)
	if err != nil {
		log.Error(logger.LogAppDatabaseError, "error", err)
		return nil, err
	}
	log.Success(logger.LogAppDatabaseConnected)

	keycloakClient, err := keycloak.NewClient(&cfg.Keycloak)
	if err != nil {
		log.Error(logger.LogKeycloakClientError, "error", err)
		return nil, err
	}
	log.Success(logger.LogKeycloakClientOK)

	employeeRepo, err := repo.NewClientRepository(db)
	if err != nil {
		log.Error(logger.LogEmployeeRepoInitError, "error", err)
		return nil, err
	}

	employeeService := services.NewService(employeeRepo, keycloakClient)
	interactorFacade := interactor.NewInteractor(employeeService)

	encoder, err := idencoder.NewHashidsEncoder(idencoder.Config{
		Secret:    cfg.IDEncoder.Secret,
		MinLength: cfg.IDEncoder.MinLength,
	}, log)
	if err != nil {
		log.Error(logger.LogIDEncoderInitError, "error", err)
		return nil, err
	}
	log.Success(logger.LogIDEncodeOK)

	msgRepo, err := messageRepo.NewMessageRepository(db)
	if err != nil {
		log.Error(logger.LogRepoMsgInitError, "error", err)
		return nil, err
	}
	log.Success(logger.LogDependencyMessageRepoInit)

	refreshInterval := 5 * time.Minute
	messagingCache := messagingCache.NewMessageCache(msgRepo, refreshInterval)

	if err := messagingCache.LoadMessages(context.Background()); err != nil {
		log.Warn(logger.LogMsgCacheLoadError, "error", err)
	}
	log.Success(logger.LogMsgCacheInit, "messages_loaded", messagingCache.MessageCount())

	messagingCache.StartAutoRefresh(context.Background())

	responseHandler := middleware.NewResponseHandler(messagingCache)

	messageService := services.NewMessageService(msgRepo)
	messageInteractor := interactor.NewMessageInteractor(messageService)
	log.Success(logger.LogDependencyMessageIntInit)

	deps, err := initDomainDependencies(db, log)
	if err != nil {
		return nil, err
	}

	var jwtValidator output.TokenValidator
	jwtConfig := jwt.JWKSConfig{
		JWKSURL:         cfg.GetKeycloakJWKSURL(),
		Issuer:          cfg.GetKeycloakIssuerURL(),
		RefreshInterval: 15 * time.Minute,
	}
	concreteValidator, err := jwt.NewJWKSValidator(context.Background(), jwtConfig)
	if err != nil {
		log.Warn(logger.LogJWKSValidatorInitFailed, "error", err)
		jwtValidator = nil
	} else {
		log.Success(logger.LogJWKSValidatorInitOK, "jwks_url", jwtConfig.JWKSURL)
		jwtValidator = concreteValidator
	}

	return &Dependencies{
		EmployeeService:              employeeService,
		EmployeeRepo:                 employeeRepo,
		Interactor:                   interactorFacade,
		KeycloakClient:               keycloakClient,
		Config:                       cfg,
		IDEncoder:                    encoder,
		ResponseHandler:              responseHandler,
		MessagingCache:               messagingCache,
		MessageInteractor:            messageInteractor,
		JWTValidator:                 jwtValidator,
		AirlineInteractor:            deps.airlineInteractor,
		AirlineEmployeeInteractor:    deps.airlineEmployeeInteractor,
		EngineInteractor:             deps.engineInteractor,
		RouteInteractor:              deps.routeInteractor,
		AirlineRouteInteractor:       deps.airlineRouteInteractor,
		AirportInteractor:            deps.airportInteractor,
		ManufacturerInteractor:       deps.manufacturerInteractor,
		AircraftModelInteractor:      deps.aircraftModelInteractor,
		TailNumberInteractor:         deps.tailNumberInteractor,
		DailyLogbookDetailInteractor: deps.dailyLogbookDetailInteractor,
		DailyLogbookInteractor:       deps.dailyLogbookInteractor,
		FlightSummaryInteractor:      deps.flightSummaryInteractor,
	}, nil
}

type domainDeps struct {
	airlineInteractor            *interactor.AirlineInteractor
	airlineEmployeeInteractor    *interactor.AirlineEmployeeInteractor
	engineInteractor             *interactor.EngineInteractor
	routeInteractor              *interactor.RouteInteractor
	airlineRouteInteractor       *interactor.AirlineRouteInteractor
	airportInteractor            *interactor.AirportInteractor
	manufacturerInteractor       *interactor.ManufacturerInteractor
	aircraftModelInteractor      *interactor.AircraftModelInteractor
	tailNumberInteractor         *interactor.TailNumberInteractor
	dailyLogbookDetailInteractor *interactor.DailyLogbookDetailInteractor
	dailyLogbookInteractor       *interactor.DailyLogbookInteractor
	flightSummaryInteractor      *interactor.FlightSummaryInteractor
}

func initDomainDependencies(db *sql.DB, log logger.Logger) (*domainDeps, error) {
	airlineRepository, err := airlineRepo.NewAirlineRepository(db)
	if err != nil {
		log.Error(logger.LogAirlineRepoInitError, "error", err)
		return nil, err
	}
	log.Success(logger.LogAirlineRepoInitOK)

	airlineService := services.NewAirlineService(airlineRepository)

	airportRepository, err := airportRepo.NewAirportRepository(db)
	if err != nil {
		log.Error(logger.LogAirportRepoInitError, "error", err)
		return nil, err
	}
	log.Success(logger.LogAirportRepoInitOK)

	airportService := services.NewAirportService(airportRepository)

	dailyLogbookRepository, err := dailyLogbookRepo.NewDailyLogbookRepository(db)
	if err != nil {
		log.Error(logger.LogDailyLogbookRepoInitError, "error", err)
		return nil, err
	}
	log.Success(logger.LogDailyLogbookRepoInitOK)

	dailyLogbookService := services.NewDailyLogbookService(dailyLogbookRepository)

	tailNumberRepository, err := tailNumberRepo.NewTailNumberRepository(db)
	if err != nil {
		log.Error(logger.LogTailNumberRepoInitError, "error", err)
		return nil, err
	}
	log.Success(logger.LogTailNumberRepoInitOK)

	aircraftModelRepository, err := aircraftModelRepo.NewAircraftModelRepository(db)
	if err != nil {
		log.Error(logger.LogAircraftModelRepoInitError, "error", err)
		return nil, err
	}
	log.Success(logger.LogAircraftModelRepoInitOK)

	routeRepository, err := routeRepo.NewRouteRepository(db)
	if err != nil {
		log.Error(logger.LogRouteRepoInitError, "error", err)
		return nil, err
	}
	log.Success(logger.LogRouteRepoInitOK)

	airlineRouteRepository, err := airlineRouteRepo.NewAirlineRouteRepository(db)
	if err != nil {
		log.Error(logger.LogAirlineRouteRepoInitError, "error", err)
		return nil, err
	}
	log.Success(logger.LogAirlineRouteRepoInitOK)

	dailyLogbookDetailRepository, err := dailyLogbookDetailRepo.NewDailyLogbookDetailRepository(db)
	if err != nil {
		log.Error(logger.LogDailyLogbookDetailRepoInitError, "error", err)
		return nil, err
	}
	log.Success(logger.LogDailyLogbookDetailRepoInitOK)

	engineRepository, err := engineRepo.NewEngineRepository(db)
	if err != nil {
		log.Error(logger.LogEngineRepoInitError, "error", err, "repository", "engine")
		return nil, err
	}
	log.Success(logger.LogEngineRepoInitOK, "repository", "engine")

	manufacturerRepository, err := manufacturerRepo.NewManufacturerRepository(db)
	if err != nil {
		log.Error(logger.LogManufacturerRepoInitError, "error", err)
		return nil, err
	}
	log.Success(logger.LogManufacturerRepoInitOK)

	airlineEmployeeRepository, err := airlineEmployeeRepo.NewAirlineEmployeeRepository(db)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error", err, "repository", "airline_employee")
		return nil, err
	}
	log.Success(logger.LogDatabaseAvailable, "repository", "airline_employee")

	return &domainDeps{
		airlineInteractor:            interactor.NewAirlineInteractor(airlineService),
		airlineEmployeeInteractor:    interactor.NewAirlineEmployeeInteractor(services.NewAirlineEmployeeService(airlineEmployeeRepository)),
		engineInteractor:             interactor.NewEngineInteractor(services.NewEngineService(engineRepository)),
		routeInteractor:              interactor.NewRouteInteractor(services.NewRouteService(routeRepository)),
		airlineRouteInteractor:       interactor.NewAirlineRouteInteractor(services.NewAirlineRouteService(airlineRouteRepository)),
		airportInteractor:            interactor.NewAirportInteractor(airportService),
		manufacturerInteractor:       interactor.NewManufacturerInteractor(services.NewManufacturerService(manufacturerRepository)),
		aircraftModelInteractor:      interactor.NewAircraftModelInteractor(services.NewAircraftModelService(aircraftModelRepository, log)),
		tailNumberInteractor:         interactor.NewTailNumberInteractor(services.NewTailNumberService(tailNumberRepository), log),
		dailyLogbookDetailInteractor: interactor.NewDailyLogbookDetailInteractor(services.NewDailyLogbookDetailService(dailyLogbookDetailRepository), dailyLogbookService, initEmployeeFlightSummaryService(db, log)),
		dailyLogbookInteractor:       interactor.NewDailyLogbookInteractor(dailyLogbookService),
		flightSummaryInteractor:      initFlightSummaryInteractor(db, log),
	}, nil
}

func initFlightSummaryInteractor(db *sql.DB, log logger.Logger) *interactor.FlightSummaryInteractor {
	fsRepo, err := flightSummaryRepo.NewFlightSummaryRepository(db)
	if err != nil {
		log.Warn(logger.LogFlightSummaryRepoInitError, "error", err)
		return nil
	}
	log.Success(logger.LogFlightSummaryRepoInitOK)

	fsService := services.NewFlightSummaryService(fsRepo)
	return interactor.NewFlightSummaryInteractor(fsService)
}

func initEmployeeFlightSummaryService(db *sql.DB, log logger.Logger) *services.EmployeeFlightSummaryServiceImpl {
	efsRepo, err := employeeFlightSummaryRepo.NewEmployeeFlightSummaryRepository(db)
	if err != nil {
		log.Warn(logger.LogFlightSummaryRepoInitError, "error", err, "repository", "employee_flight_summary")
		return nil
	}
	log.Success(logger.LogFlightSummaryRepoInitOK, "repository", "employee_flight_summary")

	return services.NewEmployeeFlightSummaryService(efsRepo)
}
