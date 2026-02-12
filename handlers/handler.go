package handlers

import (
	"github.com/champion19/api-flighthours/core/interactor"
	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/input"
	"github.com/champion19/api-flighthours/middleware"
	"github.com/champion19/api-flighthours/platform/cache/messaging"
	"github.com/champion19/api-flighthours/platform/logger"
	"github.com/champion19/api-flighthours/tools/idencoder"
	"github.com/gin-gonic/gin"
)

type handler struct {
	EmployeeService              input.Service
	Interactor                   input.EmployeeInteractor
	IDEncoder                    *idencoder.HashidsEncoder
	Response                     *middleware.ResponseHandler
	MessageInteractor            *interactor.MessageInteractor
	MessagingCache               *messaging.MessageCache
	AirlineInteractor            *interactor.AirlineInteractor
	AirportInteractor            *interactor.AirportInteractor
	AirlineEmployeeInteractor    *interactor.AirlineEmployeeInteractor
	EngineInteractor             *interactor.EngineInteractor
	RouteInteractor              *interactor.RouteInteractor
	AirlineRouteInteractor       *interactor.AirlineRouteInteractor
	ManufacturerInteractor       *interactor.ManufacturerInteractor
	AircraftModelInteractor      *interactor.AircraftModelInteractor
	LicensePlateInteractor       *interactor.LicensePlateInteractor
	DailyLogbookDetailInteractor *interactor.DailyLogbookDetailInteractor
	DailyLogbookInteractor       *interactor.DailyLogbookInteractor
}

func New(
	service input.Service,
	employeeInteractor input.EmployeeInteractor,
	idEncoder *idencoder.HashidsEncoder,
	response *middleware.ResponseHandler,
	messageInteractor *interactor.MessageInteractor,
	messagingCache *messaging.MessageCache,
	airlineInteractor *interactor.AirlineInteractor,
	airlineEmployeeInteractor *interactor.AirlineEmployeeInteractor,
	engineInteractor *interactor.EngineInteractor,
	routeInteractor *interactor.RouteInteractor,
	manufacturerInteractor *interactor.ManufacturerInteractor,
	airportInteractor *interactor.AirportInteractor,
	airlineRouteInteractor *interactor.AirlineRouteInteractor,
	dailyLogbookDetailInteractor *interactor.DailyLogbookDetailInteractor,
	dailyLogbookInteractor *interactor.DailyLogbookInteractor,
	aircraftModelInteractor *interactor.AircraftModelInteractor,
	aircraftRegistrationInteractor *interactor.LicensePlateInteractor) *handler {
	return &handler{
		EmployeeService:           service,
		Interactor:                employeeInteractor,
		IDEncoder:                 idEncoder,
		Response:                  response,
		MessageInteractor:         messageInteractor,
		MessagingCache:            messagingCache,
		AirlineInteractor:         airlineInteractor,
		AirlineEmployeeInteractor: airlineEmployeeInteractor,
		EngineInteractor:          engineInteractor,
		RouteInteractor:           routeInteractor,
		AirlineRouteInteractor:    airlineRouteInteractor,
		DailyLogbookDetailInteractor:    dailyLogbookDetailInteractor,
		DailyLogbookInteractor:    dailyLogbookInteractor,
		ManufacturerInteractor:    manufacturerInteractor,
		AirportInteractor:         airportInteractor,
		AircraftModelInteractor:   aircraftModelInteractor,
		LicensePlateInteractor:    aircraftRegistrationInteractor,
	}
}

var Logger = logger.NewSlogLogger()

func (h *handler) EncodeID(uuid string) (string, error) {
	encodedID, err := h.IDEncoder.Encode(uuid)
	if err != nil {
		Logger.Error(logger.LogMessageIDEncodeError,
			"uuid", uuid,
			"error", err)
		return "", err
	}
	return encodedID, nil
}

func (h *handler) DecodeID(encodedID string) (string, error) {
	uuid, err := h.IDEncoder.Decode(encodedID)
	if err != nil {
		Logger.Error(logger.LogMessageIDDecodeError,
			"encoded_id", encodedID,
			"error", err)
		return "", err
	}
	return uuid, nil
}

func (h *handler) HandleIDEncodingError(c *gin.Context, uuid string, err error) {
	Logger.Error(logger.LogMessageIDEncodeError,
		"uuid", uuid,
		"error", err,
		"client_ip", c.ClientIP())
	c.Error(domain.ErrInternalServer)
}

func (h *handler) HandleIDDecodingError(c *gin.Context, encodedID string, err error) {
	Logger.Error(logger.LogMessageIDDecodeError,
		"encoded_id", encodedID,
		"error", err,
		"client_ip", c.ClientIP())
	c.Error(domain.ErrInvalidID)
}

func (h *handler) resolveID(inputID string) (string, string) {
	if inputID == "" {
		return "", ""
	}

	uuid, err := h.DecodeID(inputID)
	if err != nil {
		Logger.Warn(logger.LogMessageIDDecodeError, "encoded_id", inputID, "error", err)
		return "", ""
	}
	return uuid, inputID
}
