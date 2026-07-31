package handlers

import (
	"github.com/champion19/api-flighthours/core/interactor"
	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/input"
	"github.com/champion19/api-flighthours/middleware"
	"github.com/champion19/api-flighthours/platform/cache/messaging"
	"github.com/champion19/api-flighthours/platform/cookie"
	"github.com/champion19/api-flighthours/platform/logger"
	"github.com/champion19/api-flighthours/tools/idencoder"
	"github.com/gin-gonic/gin"
)

type handler struct {
	EmployeeService              input.Service
	Interactor                   input.EmployeeInteractor
	IDEncoder                    *idencoder.HashidsEncoder
	Response                     *middleware.ResponseHandler
	CookieManager                *cookie.Manager
	MessageInteractor            *interactor.MessageInteractor
	MessagingCache               *messaging.MessageCache
	AirlineInteractor            *interactor.AirlineInteractor
	AirportInteractor            *interactor.AirportInteractor
	AirlineEmployeeInteractor    *interactor.AirlineEmployeeInteractor
	EngineInteractor             *interactor.EngineInteractor
	ManufacturerInteractor       *interactor.ManufacturerInteractor
	AircraftModelInteractor      *interactor.AircraftModelInteractor
	TailNumberInteractor         *interactor.TailNumberInteractor
	DailyLogbookDetailInteractor *interactor.DailyLogbookDetailInteractor
	DailyLogbookInteractor       *interactor.DailyLogbookInteractor
	FlightSummaryInteractor      *interactor.FlightSummaryInteractor
	CrewMemberInteractor         *interactor.CrewMemberInteractor
}

type HandlerDeps struct {
	Service                      input.Service
	EmployeeInteractor           input.EmployeeInteractor
	IDEncoder                    *idencoder.HashidsEncoder
	Response                     *middleware.ResponseHandler
	CookieManager                *cookie.Manager
	MessageInteractor            *interactor.MessageInteractor
	MessagingCache               *messaging.MessageCache
	AirlineInteractor            *interactor.AirlineInteractor
	AirlineEmployeeInteractor    *interactor.AirlineEmployeeInteractor
	EngineInteractor             *interactor.EngineInteractor
	DailyLogbookDetailInteractor *interactor.DailyLogbookDetailInteractor
	DailyLogbookInteractor       *interactor.DailyLogbookInteractor
	ManufacturerInteractor       *interactor.ManufacturerInteractor
	AirportInteractor            *interactor.AirportInteractor
	AircraftModelInteractor      *interactor.AircraftModelInteractor
	TailNumberInteractor         *interactor.TailNumberInteractor
	FlightSummaryInteractor      *interactor.FlightSummaryInteractor
	CrewMemberInteractor         *interactor.CrewMemberInteractor
}

func New(deps HandlerDeps) *handler {
	return &handler{
		EmployeeService:              deps.Service,
		Interactor:                   deps.EmployeeInteractor,
		IDEncoder:                    deps.IDEncoder,
		Response:                     deps.Response,
		CookieManager:                deps.CookieManager,
		MessageInteractor:            deps.MessageInteractor,
		MessagingCache:               deps.MessagingCache,
		AirlineInteractor:            deps.AirlineInteractor,
		AirlineEmployeeInteractor:    deps.AirlineEmployeeInteractor,
		EngineInteractor:             deps.EngineInteractor,
		DailyLogbookDetailInteractor: deps.DailyLogbookDetailInteractor,
		DailyLogbookInteractor:       deps.DailyLogbookInteractor,
		ManufacturerInteractor:       deps.ManufacturerInteractor,
		AirportInteractor:            deps.AirportInteractor,
		AircraftModelInteractor:      deps.AircraftModelInteractor,
		TailNumberInteractor:         deps.TailNumberInteractor,
		FlightSummaryInteractor:      deps.FlightSummaryInteractor,
		CrewMemberInteractor:         deps.CrewMemberInteractor,
	}
}

var log logger.Logger = logger.NewSlogLogger()

func (h *handler) EncodeID(uuid string) (string, error) {
	encodedID, err := h.IDEncoder.Encode(uuid)
	if err != nil {
		log.Error(logger.LogMessageIDEncodeError,
			"uuid", uuid,
			"error", err)
		return "", err
	}
	return encodedID, nil
}

func (h *handler) DecodeID(encodedID string) (string, error) {
	uuid, err := h.IDEncoder.Decode(encodedID)
	if err != nil {
		log.Error(logger.LogMessageIDDecodeError,
			"encoded_id", encodedID,
			"error", err)
		return "", err
	}
	return uuid, nil
}

func (h *handler) HandleIDEncodingError(c *gin.Context, uuid string, err error) {
	log.Error(logger.LogMessageIDEncodeError,
		"uuid", uuid,
		"error", err,
		"client_ip", c.ClientIP())
	c.Error(domain.ErrInternalServer)
}

func (h *handler) HandleIDDecodingError(c *gin.Context, encodedID string, err error) {
	log.Error(logger.LogMessageIDDecodeError,
		"encoded_id", encodedID,
		"error", err,
		"client_ip", c.ClientIP())
	c.Error(domain.ErrInvalidID)
}

// encodeTailNumberID re-encodes a DailyLogbookResponse's raw tail_number_id
// (a UUID at this point, resolved from the request) into an obfuscated ID,
// keeping it consistent with every other ID in the API response.
func (h *handler) encodeTailNumberID(response *DailyLogbookResponse) {
	if response.TailNumberID == nil || *response.TailNumberID == "" {
		return
	}
	if encoded, err := h.EncodeID(*response.TailNumberID); err == nil {
		response.TailNumberID = &encoded
	}
}

func (h *handler) resolveID(inputID string) (string, string) {
	if inputID == "" {
		return "", ""
	}

	uuid, err := h.DecodeID(inputID)
	if err != nil {
		log.Warn(logger.LogMessageIDDecodeError, "encoded_id", inputID, "error", err)
		return "", ""
	}
	return uuid, inputID
}
