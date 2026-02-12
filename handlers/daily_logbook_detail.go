package handlers

import (
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

type CreateDailyLogbookDetailRequest struct {
	FlightRealDate string  `json:"flight_real_date"`
	FlightNumber   string  `json:"flight_number"`
	AirlineRouteID string  `json:"airline_route_id"`
	LicensePlateID string  `json:"license_plate_id"`
	Passengers     *int    `json:"passengers,omitempty"`
	OutTime        *string `json:"out_time,omitempty"`     // TIME format HH:MM (nullable)
	TakeoffTime    *string `json:"takeoff_time,omitempty"` // TIME format HH:MM (nullable)
	LandingTime    *string `json:"landing_time,omitempty"` // TIME format HH:MM (nullable)
	InTime         *string `json:"in_time,omitempty"`      // TIME format HH:MM (nullable)
	PilotRole      *string `json:"pilot_role,omitempty"`   // nullable
	CompanionName  *string `json:"companion_name,omitempty"`
	AirTime        *string `json:"air_time,omitempty"`   // TIME format HH:MM (nullable)
	BlockTime      *string `json:"block_time,omitempty"` // TIME format HH:MM (nullable)
	DutyTime       *string `json:"duty_time,omitempty"`  // TIME format HH:MM (nullable)
	ApproachType   *string `json:"approach_type,omitempty"`
	FlightType     *string `json:"flight_type,omitempty"`
}

func (r *CreateDailyLogbookDetailRequest) Sanitize() {
	r.FlightRealDate = TrimString(r.FlightRealDate)
	r.FlightNumber = TrimString(r.FlightNumber)
	r.AirlineRouteID = TrimString(r.AirlineRouteID)
	r.LicensePlateID = TrimString(r.LicensePlateID)
	r.OutTime = TrimStringPtr(r.OutTime)
	r.TakeoffTime = TrimStringPtr(r.TakeoffTime)
	r.LandingTime = TrimStringPtr(r.LandingTime)
	r.InTime = TrimStringPtr(r.InTime)
	r.PilotRole = TrimStringPtr(r.PilotRole)
	r.CompanionName = TrimStringPtr(r.CompanionName)
	r.AirTime = TrimStringPtr(r.AirTime)
	r.BlockTime = TrimStringPtr(r.BlockTime)
	r.DutyTime = TrimStringPtr(r.DutyTime)
	r.ApproachType = TrimStringPtr(r.ApproachType)
	r.FlightType = TrimStringPtr(r.FlightType)
}

type UpdateDailyLogbookDetailRequest struct {
	FlightRealDate string  `json:"flight_real_date"`
	FlightNumber   string  `json:"flight_number"`
	AirlineRouteID string  `json:"airline_route_id"`
	LicensePlateID string  `json:"license_plate_id"`
	Passengers     *int    `json:"passengers,omitempty"`
	OutTime        *string `json:"out_time,omitempty"`     // TIME format HH:MM (nullable)
	TakeoffTime    *string `json:"takeoff_time,omitempty"` // TIME format HH:MM (nullable)
	LandingTime    *string `json:"landing_time,omitempty"` // TIME format HH:MM (nullable)
	InTime         *string `json:"in_time,omitempty"`      // TIME format HH:MM (nullable)
	PilotRole      *string `json:"pilot_role,omitempty"`   // nullable
	CompanionName  *string `json:"companion_name,omitempty"`
	AirTime        *string `json:"air_time,omitempty"`   // TIME format HH:MM (nullable)
	BlockTime      *string `json:"block_time,omitempty"` // TIME format HH:MM (nullable)
	DutyTime       *string `json:"duty_time,omitempty"`  // TIME format HH:MM (nullable)
	ApproachType   *string `json:"approach_type,omitempty"`
	FlightType     *string `json:"flight_type,omitempty"`
}

func (r *UpdateDailyLogbookDetailRequest) Sanitize() {
	r.FlightRealDate = TrimString(r.FlightRealDate)
	r.FlightNumber = TrimString(r.FlightNumber)
	r.AirlineRouteID = TrimString(r.AirlineRouteID)
	r.LicensePlateID = TrimString(r.LicensePlateID)
	r.OutTime = TrimStringPtr(r.OutTime)
	r.TakeoffTime = TrimStringPtr(r.TakeoffTime)
	r.LandingTime = TrimStringPtr(r.LandingTime)
	r.InTime = TrimStringPtr(r.InTime)
	r.PilotRole = TrimStringPtr(r.PilotRole)
	r.CompanionName = TrimStringPtr(r.CompanionName)
	r.AirTime = TrimStringPtr(r.AirTime)
	r.BlockTime = TrimStringPtr(r.BlockTime)
	r.DutyTime = TrimStringPtr(r.DutyTime)
	r.ApproachType = TrimStringPtr(r.ApproachType)
	r.FlightType = TrimStringPtr(r.FlightType)
}

// DailyLogbookDetailResponse represents the response for a detail
type DailyLogbookDetailResponse struct {
	ID                  string  `json:"id"`
	DailyLogbookID      string  `json:"-"`
	FlightRealDate      string  `json:"flight_real_date"`
	FlightNumber        string  `json:"flight_number"`
	AirlineRouteID      string  `json:"-"`
	LicensePlateID      string  `json:"-"`
	Passengers          *int    `json:"-"`
	OutTime             *string `json:"-"`
	TakeoffTime         *string `json:"-"`
	LandingTime         *string `json:"-"`
	InTime              *string `json:"-"`
	PilotRole           *string `json:"-"`
	CompanionName       *string `json:"-"`
	AirTime             *string `json:"-"`
	BlockTime           *string `json:"-"`
	DutyTime            *string `json:"-"`
	ApproachType        *string `json:"-"`
	FlightType          *string `json:"-"`
	LogDate             string  `json:"-"`
	RouteCode           string  `json:"route_code,omitempty"`
	OriginIataCode      string  `json:"origin_iata_code,omitempty"`
	DestinationIataCode string  `json:"destination_iata_code,omitempty"`
	AirlineCode         string  `json:"airline_code,omitempty"`
	LicensePlate        string  `json:"-"`
	ModelName           string  `json:"-"`
	Links               []Link  `json:"_links,omitempty"`
}

// ToDomainDailyLogbookDetail converts a create request to domain model
func ToDomainDailyLogbookDetail(logbookID string, req CreateDailyLogbookDetailRequest) domain.DailyLogbookDetail {
	detail := domain.DailyLogbookDetail{
		DailyLogbookID: logbookID,
		FlightRealDate: req.FlightRealDate,
		FlightNumber:   req.FlightNumber,
		AirlineRouteID: req.AirlineRouteID,
		LicensePlateID: req.LicensePlateID,
		Passengers:     req.Passengers,
		OutTime:        req.OutTime,
		TakeoffTime:    req.TakeoffTime,
		LandingTime:    req.LandingTime,
		InTime:         req.InTime,
		CompanionName:  req.CompanionName,
		AirTime:        req.AirTime,
		BlockTime:      req.BlockTime,
		DutyTime:       req.DutyTime,
		FlightType:     req.FlightType,
	}

	if req.PilotRole != nil {
		pilotRole := domain.PilotRole(*req.PilotRole)
		detail.PilotRole = &pilotRole
	}

	if req.ApproachType != nil {
		approachType := domain.ApproachType(*req.ApproachType)
		detail.ApproachType = &approachType
	}

	return detail
}

// ToDomainDailyLogbookDetailUpdate converts an update request to domain model
func ToDomainDailyLogbookDetailUpdate(id string, req UpdateDailyLogbookDetailRequest) domain.DailyLogbookDetail {
	detail := domain.DailyLogbookDetail{
		ID:             id,
		FlightRealDate: req.FlightRealDate,
		FlightNumber:   req.FlightNumber,
		AirlineRouteID: req.AirlineRouteID,
		LicensePlateID: req.LicensePlateID,
		Passengers:     req.Passengers,
		OutTime:        req.OutTime,
		TakeoffTime:    req.TakeoffTime,
		LandingTime:    req.LandingTime,
		InTime:         req.InTime,
		CompanionName:  req.CompanionName,
		AirTime:        req.AirTime,
		BlockTime:      req.BlockTime,
		DutyTime:       req.DutyTime,
		FlightType:     req.FlightType,
	}

	if req.PilotRole != nil {
		pilotRole := domain.PilotRole(*req.PilotRole)
		detail.PilotRole = &pilotRole
	}

	if req.ApproachType != nil {
		approachType := domain.ApproachType(*req.ApproachType)
		detail.ApproachType = &approachType
	}

	return detail
}

// FromDomainDailyLogbookDetail converts domain model to response DTO
func FromDomainDailyLogbookDetail(d *domain.DailyLogbookDetail, encodedID, encodedLogbookID, encodedRouteID, encodedAircraftID string) DailyLogbookDetailResponse {
	response := DailyLogbookDetailResponse{
		ID:                  encodedID,
		DailyLogbookID:      encodedLogbookID,
		FlightRealDate:      d.FlightRealDate,
		FlightNumber:        d.FlightNumber,
		AirlineRouteID:      encodedRouteID,
		LicensePlateID:      encodedAircraftID,
		Passengers:          d.Passengers,
		OutTime:             d.OutTime,
		TakeoffTime:         d.TakeoffTime,
		LandingTime:         d.LandingTime,
		InTime:              d.InTime,
		CompanionName:       d.CompanionName,
		AirTime:             d.AirTime,
		BlockTime:           d.BlockTime,
		DutyTime:            d.DutyTime,
		LogDate:             d.LogDate,
		RouteCode:           d.RouteCode,
		OriginIataCode:      d.OriginIataCode,
		DestinationIataCode: d.DestinationIataCode,
		AirlineCode:         d.AirlineCode,
		LicensePlate:        d.LicensePlate,
		ModelName:           d.ModelName,
	}

	if d.PilotRole != nil {
		pilotRoleStr := string(*d.PilotRole)
		response.PilotRole = &pilotRoleStr
	}

	if d.ApproachType != nil {
		approachTypeStr := string(*d.ApproachType)
		response.ApproachType = &approachTypeStr
	}

	response.FlightType = d.FlightType

	return response
}
