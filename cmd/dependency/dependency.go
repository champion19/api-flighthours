package dependency

import (
	"context"
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
	repo "github.com/champion19/api-flighthours/platform/databases/repositories/employee"
	engineRepo "github.com/champion19/api-flighthours/platform/databases/repositories/engine"
	licensePlateRepo "github.com/champion19/api-flighthours/platform/databases/repositories/license_plate"
	manufacturerRepo "github.com/champion19/api-flighthours/platform/databases/repositories/manufacturer"
	messageRepo "github.com/champion19/api-flighthours/platform/databases/repositories/message"
	routeRepo "github.com/champion19/api-flighthours/platform/databases/repositories/route"
	"github.com/champion19/api-flighthours/platform/identity_provider/keycloak"
	"github.com/champion19/api-flighthours/platform/jwt"
	"github.com/champion19/api-flighthours/platform/logger"
	"github.com/champion19/api-flighthours/tools/idencoder"
)

var log logger.Logger = logger.NewSlogLogger()

type Dependencies struct {
	EmployeeService           input.Service
	EmployeeRepo              output.Repository
	Interactor                *interactor.Interactor
	KeycloakClient            output.AuthClient
	Config                    *config.Config
	IDEncoder                 *idencoder.HashidsEncoder
	ResponseHandler           *middleware.ResponseHandler
	MessagingCache            *messagingCache.MessageCache
	MessageInteractor         *interactor.MessageInteractor
	JWTValidator              *jwt.JWKSValidator
	AirlineInteractor         *interactor.AirlineInteractor
	AirlineEmployeeInteractor *interactor.AirlineEmployeeInteractor
	EngineInteractor          *interactor.EngineInteractor
	RouteInteractor           *interactor.RouteInteractor
	AirlineRouteInteractor    *interactor.AirlineRouteInteractor
	AirportInteractor         *interactor.AirportInteractor
	ManufacturerInteractor    *interactor.ManufacturerInteractor
	AircraftModelInteractor   *interactor.AircraftModelInteractor
	LicensePlateInteractor    *interactor.LicensePlateInteractor
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

	airlineRepository, err := airlineRepo.NewAirlineRepository(db)
	if err != nil {
		log.Error(logger.LogAirlineRepoInitError, "error", err)
		return nil, err
	}
	log.Success(logger.LogAirlineRepoInitOK)

	airlineService := services.NewAirlineService(airlineRepository)
	airlineInteractor := interactor.NewAirlineInteractor(airlineService)

	airportRepository, err := airportRepo.NewAirportRepository(db)
	if err != nil {
		log.Error(logger.LogAirportRepoInitError, "error", err)
		return nil, err
	}
	log.Success(logger.LogAirportRepoInitOK)

	airportService := services.NewAirportService(airportRepository)
	airportInteractor := interactor.NewAirportInteractor(airportService)

	licensePlateRepository, err := licensePlateRepo.NewLicensePlateRepository(db)
	if err != nil {
		log.Error(logger.LogLicensePlateRepoInitError, "error", err)
		return nil, err
	}
	log.Success(logger.LogLicensePlateRepoInitOK)

	licensePlateService := services.NewLicensePlateService(licensePlateRepository)
	licensePlateInteractor := interactor.NewLicensePlateInteractor(licensePlateService, log)

	aircraftModelRepository, err := aircraftModelRepo.NewAircraftModelRepository(db)
	if err != nil {
		log.Error(logger.LogAircraftModelRepoInitError, "error", err)
		return nil, err
	}
	log.Success(logger.LogAircraftModelRepoInitOK)

	aircraftModelService := services.NewAircraftModelService(aircraftModelRepository, log)
	aircraftModelInteractor := interactor.NewAircraftModelInteractor(aircraftModelService)

	routeRepository, err := routeRepo.NewRouteRepository(db)
	if err != nil {
		log.Error(logger.LogRouteRepoInitError, "error", err)
		return nil, err
	}
	log.Success(logger.LogRouteRepoInitOK)

	routeService := services.NewRouteService(routeRepository, log)
	routeInteractor := interactor.NewRouteInteractor(routeService, log)

	airlineRouteRepository, err := airlineRouteRepo.NewAirlineRouteRepository(db)
	if err != nil {
		log.Error(logger.LogAirlineRouteRepoInitError, "error", err)
		return nil, err
	}
	log.Success(logger.LogAirlineRouteRepoInitOK)

	airlineRouteService := services.NewAirlineRouteService(airlineRouteRepository)
	airlineRouteInteractor := interactor.NewAirlineRouteInteractor(airlineRouteService)

	engineRepository, err := engineRepo.NewEngineRepository(db)
	if err != nil {
		log.Error(logger.LogEngineRepoInitError, "error", err, "repository", "engine")
		return nil, err
	}
	log.Success(logger.LogEngineRepoInitOK, "repository", "engine")

	engineService := services.NewEngineService(engineRepository)
	engineInteractor := interactor.NewEngineInteractor(engineService)

	manufacturerRepository, err := manufacturerRepo.NewManufacturerRepository(db)
	if err != nil {
		log.Error(logger.LogManufacturerRepoInitError, "error", err)
		return nil, err
	}
	log.Success(logger.LogManufacturerRepoInitOK)

	manufacturerService := services.NewManufacturerService(manufacturerRepository)
	manufacturerInteractor := interactor.NewManufacturerInteractor(manufacturerService)

	airlineEmployeeRepository, err := airlineEmployeeRepo.NewAirlineEmployeeRepository(db)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error", err, "repository", "airline_employee")
		return nil, err
	}
	log.Success(logger.LogDatabaseAvailable, "repository", "airline_employee")

	airlineEmployeeService := services.NewAirlineEmployeeService(airlineEmployeeRepository)
	airlineEmployeeInteractor := interactor.NewAirlineEmployeeInteractor(airlineEmployeeService)

	var jwtValidator *jwt.JWKSValidator
	jwtConfig := jwt.JWKSConfig{
		JWKSURL:         cfg.GetKeycloakJWKSURL(),
		Issuer:          cfg.GetKeycloakIssuerURL(),
		RefreshInterval: 15 * time.Minute,
	}
	jwtValidator, err = jwt.NewJWKSValidator(context.Background(), jwtConfig)
	if err != nil {
		log.Warn(logger.LogJWKSValidatorInitFailed, "error", err)
		jwtValidator = nil
	} else {
		log.Success(logger.LogJWKSValidatorInitOK, "jwks_url", jwtConfig.JWKSURL)
	}

	return &Dependencies{
		EmployeeService:           employeeService,
		EmployeeRepo:              employeeRepo,
		Interactor:                interactorFacade,
		KeycloakClient:            keycloakClient,
		Config:                    cfg,
		IDEncoder:                 encoder,
		ResponseHandler:           responseHandler,
		MessagingCache:            messagingCache,
		MessageInteractor:         messageInteractor,
		JWTValidator:              jwtValidator,
		AirlineInteractor:         airlineInteractor,
		AirlineEmployeeInteractor: airlineEmployeeInteractor,
		EngineInteractor:          engineInteractor,
		RouteInteractor:           routeInteractor,
		AirlineRouteInteractor:    airlineRouteInteractor,
		AirportInteractor:         airportInteractor,
		ManufacturerInteractor:    manufacturerInteractor,
		AircraftModelInteractor:   aircraftModelInteractor,
		LicensePlateInteractor:    licensePlateInteractor,
	}, nil
}
