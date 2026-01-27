package dependency

import (
	"context"
	"time"

	"github.com/champion19/flighthours-api/config"
	"github.com/champion19/flighthours-api/core/interactor"
	"github.com/champion19/flighthours-api/core/interactor/services"
	"github.com/champion19/flighthours-api/core/ports/input"
	"github.com/champion19/flighthours-api/core/ports/output"
	"github.com/champion19/flighthours-api/middleware"
	messagingCache "github.com/champion19/flighthours-api/platform/cache/messaging"
	mysql "github.com/champion19/flighthours-api/platform/databases/mysql"
	repo "github.com/champion19/flighthours-api/platform/databases/repositories/employee"
	messageRepo "github.com/champion19/flighthours-api/platform/databases/repositories/message"
	"github.com/champion19/flighthours-api/platform/identity_provider/keycloak"
	"github.com/champion19/flighthours-api/platform/jwt"
	"github.com/champion19/flighthours-api/platform/logger"
	"github.com/champion19/flighthours-api/tools/idencoder"
)

type Dependencies struct {
	EmployeeService   input.Service
	EmployeeRepo      output.Repository
	Interactor        *interactor.Interactor
	KeycloakClient    output.AuthClient
	Config            *config.Config
	Logger            logger.Logger
	IDEncoder         *idencoder.HashidsEncoder
	ResponseHandler   *middleware.ResponseHandler
	MessagingCache    *messagingCache.MessageCache
	MessageInteractor *interactor.MessageInteractor
	JWTValidator      *jwt.JWKSValidator
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

	db, err := mysql.GetDB(cfg.Database, log)
	if err != nil {
		log.Error(logger.LogAppDatabaseError, "error", err)
		return nil, err
	}
	log.Success(logger.LogAppDatabaseConnected)

	keycloakClient, err := keycloak.NewClient(&cfg.Keycloak, log)
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

	employeeService := services.NewService(employeeRepo, keycloakClient, log)
	interactorFacade := interactor.NewInteractor(employeeService, log)

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

	messageService := services.NewMessageService(msgRepo, log)
	messageInteractor := interactor.NewMessageInteractor(messageService, log)
	log.Success(logger.LogDependencyMessageIntInit)

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
		EmployeeService:   employeeService,
		EmployeeRepo:      employeeRepo,
		Interactor:        interactorFacade,
		KeycloakClient:    keycloakClient,
		Config:            cfg,
		Logger:            log,
		IDEncoder:         encoder,
		ResponseHandler:   responseHandler,
		MessagingCache:    messagingCache,
		MessageInteractor: messageInteractor,
		JWTValidator:      jwtValidator,
	}, nil
}
