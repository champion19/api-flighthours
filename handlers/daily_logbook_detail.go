package handlers

import (
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

type CreateDailyLogbookDetailRequest struct {
	FlightRealDate string  `json:"flight_real_date"`
	FlightNumber   string  `json:"flight_number"`
	AirlineRouteID string  `json:"airline_route_id"`
	TailNumberID   string  `json:"tail_number_id"`
	Passengers     *int    `json:"passengers,omitempty"`
	OutTime        *string `json:"out_time,omitempty"`     // TIME format HH:MM (nullable)
	TakeoffTime    *string `json:"takeoff_time,omitempty"` // TIME format HH:MM (nullable)
	LandingTime    *string `json:"landing_time,omitempty"` // TIME format HH:MM (nullable)
	InTime         *string `json:"in_time,omitempty"`      // TIME format HH:MM (nullable)
	PilotRole      *string `json:"pilot_role,omitempty"`   // nullable
	CrewRole       *string `json:"crew_role,omitempty"`    // nullable
	CompanionName  *string `json:"companion_name,omitempty"`
	AirTime        *string `json:"air_time,omitempty"`   // TIME format HH:MM (nullable)
	BlockTime      *string `json:"block_time,omitempty"` // TIME format HH:MM (nullable)
	ApproachType   *string `json:"approach_type,omitempty"`
	FlightType     *string `json:"flight_type,omitempty"`
}

func (r *CreateDailyLogbookDetailRequest) Sanitize() {
	r.FlightRealDate = TrimString(r.FlightRealDate)
	r.FlightNumber = TrimString(r.FlightNumber)
	r.AirlineRouteID = TrimString(r.AirlineRouteID)
	r.TailNumberID = TrimString(r.TailNumberID)
	r.OutTime = TrimStringPtr(r.OutTime)
	r.TakeoffTime = TrimStringPtr(r.TakeoffTime)
	r.LandingTime = TrimStringPtr(r.LandingTime)
	r.InTime = TrimStringPtr(r.InTime)
	r.PilotRole = TrimStringPtr(r.PilotRole)
	r.CrewRole = TrimStringPtr(r.CrewRole)
	r.CompanionName = TrimStringPtr(r.CompanionName)
	r.AirTime = TrimStringPtr(r.AirTime)
	r.BlockTime = TrimStringPtr(r.BlockTime)
	r.ApproachType = TrimStringPtr(r.ApproachType)
	r.FlightType = TrimStringPtr(r.FlightType)
}

type UpdateDailyLogbookDetailRequest struct {
	FlightRealDate string  `json:"flight_real_date"`
	FlightNumber   string  `json:"flight_number"`
	AirlineRouteID string  `json:"airline_route_id"`
	TailNumberID   string  `json:"tail_number_id"`
	Passengers     *int    `json:"passengers,omitempty"`
	OutTime        *string `json:"out_time,omitempty"`     // TIME format HH:MM (nullable)
	TakeoffTime    *string `json:"takeoff_time,omitempty"` // TIME format HH:MM (nullable)
	LandingTime    *string `json:"landing_time,omitempty"` // TIME format HH:MM (nullable)
	InTime         *string `json:"in_time,omitempty"`      // TIME format HH:MM (nullable)
	PilotRole      *string `json:"pilot_role,omitempty"`   // nullable
	CrewRole       *string `json:"crew_role,omitempty"`    // nullable
	CompanionName  *string `json:"companion_name,omitempty"`
	AirTime        *string `json:"air_time,omitempty"`   // TIME format HH:MM (nullable)
	BlockTime      *string `json:"block_time,omitempty"` // TIME format HH:MM (nullable)
	ApproachType   *string `json:"approach_type,omitempty"`
	FlightType     *string `json:"flight_type,omitempty"`
}

func (r *UpdateDailyLogbookDetailRequest) Sanitize() {
	r.FlightRealDate = TrimString(r.FlightRealDate)
	r.FlightNumber = TrimString(r.FlightNumber)
	r.AirlineRouteID = TrimString(r.AirlineRouteID)
	r.TailNumberID = TrimString(r.TailNumberID)
	r.OutTime = TrimStringPtr(r.OutTime)
	r.TakeoffTime = TrimStringPtr(r.TakeoffTime)
	r.LandingTime = TrimStringPtr(r.LandingTime)
	r.InTime = TrimStringPtr(r.InTime)
	r.PilotRole = TrimStringPtr(r.PilotRole)
	r.CrewRole = TrimStringPtr(r.CrewRole)
	r.CompanionName = TrimStringPtr(r.CompanionName)
	r.AirTime = TrimStringPtr(r.AirTime)
	r.BlockTime = TrimStringPtr(r.BlockTime)
	r.ApproachType = TrimStringPtr(r.ApproachType)
	r.FlightType = TrimStringPtr(r.FlightType)
}

// ============================================
// RESPONSE DTOs
// ============================================

// DailyLogbookDetailResponse represents the response for a detail
type DailyLogbookDetailResponse struct {
	ID                  string  `json:"id"`
	DailyLogbookID      string  `json:"daily_logbook_id"`
	FlightRealDate      string  `json:"flight_real_date"`
	FlightNumber        string  `json:"flight_number"`
	AirlineRouteID      string  `json:"airline_route_id"`
	TailNumberID        string  `json:"tail_number_id"`
	Passengers          *int    `json:"passengers,omitempty"`
	OutTime             *string `json:"out_time,omitempty"`
	TakeoffTime         *string `json:"takeoff_time,omitempty"`
	LandingTime         *string `json:"landing_time,omitempty"`
	InTime              *string `json:"in_time,omitempty"`
	PilotRole           *string `json:"pilot_role,omitempty"`
	CrewRole            *string `json:"crew_role,omitempty"`
	CompanionName       *string `json:"companion_name,omitempty"`
	AirTime             *string `json:"air_time,omitempty"`
	BlockTime           *string `json:"block_time,omitempty"`
	ApproachType        *string `json:"approach_type,omitempty"`
	FlightType          *string `json:"flight_type,omitempty"`
	LogDate             string  `json:"log_date,omitempty"`
	RouteCode           string  `json:"route_code,omitempty"`
	OriginIataCode      string  `json:"origin_iata_code,omitempty"`
	DestinationIataCode string  `json:"destination_iata_code,omitempty"`
	AirlineCode         string  `json:"airline_code,omitempty"`
	TailNumber          string  `json:"tail_number,omitempty"`
	ModelName           string  `json:"model_name,omitempty"`
	Links               []Link  `json:"_links,omitempty"`
}

// ============================================
// MAPPERS
// ============================================

// ToDomainDailyLogbookDetail converts a create request to domain model
func ToDomainDailyLogbookDetail(logbookID string, req CreateDailyLogbookDetailRequest) domain.DailyLogbookDetail {
	detail := domain.DailyLogbookDetail{
		DailyLogbookID: logbookID,
		FlightRealDate: req.FlightRealDate,
		FlightNumber:   req.FlightNumber,
		AirlineRouteID: req.AirlineRouteID,
		TailNumberID:   req.TailNumberID,
		Passengers:     req.Passengers,
		OutTime:        req.OutTime,
		TakeoffTime:    req.TakeoffTime,
		LandingTime:    req.LandingTime,
		InTime:         req.InTime,
		CompanionName:  req.CompanionName,
		AirTime:        req.AirTime,
		BlockTime:      req.BlockTime,
		FlightType:     req.FlightType,
	}

	if req.PilotRole != nil {
		pilotRole := domain.PilotRole(*req.PilotRole)
		detail.PilotRole = &pilotRole
	}

	if req.CrewRole != nil {
		crewRole := domain.CrewRole(*req.CrewRole)
		detail.CrewRole = &crewRole
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
		TailNumberID:   req.TailNumberID,
		Passengers:     req.Passengers,
		OutTime:        req.OutTime,
		TakeoffTime:    req.TakeoffTime,
		LandingTime:    req.LandingTime,
		InTime:         req.InTime,
		CompanionName:  req.CompanionName,
		AirTime:        req.AirTime,
		BlockTime:      req.BlockTime,
		FlightType:     req.FlightType,
	}

	if req.PilotRole != nil {
		pilotRole := domain.PilotRole(*req.PilotRole)
		detail.PilotRole = &pilotRole
	}

	if req.CrewRole != nil {
		crewRole := domain.CrewRole(*req.CrewRole)
		detail.CrewRole = &crewRole
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
		TailNumberID:        encodedAircraftID,
		Passengers:          d.Passengers,
		OutTime:             d.OutTime,
		TakeoffTime:         d.TakeoffTime,
		LandingTime:         d.LandingTime,
		InTime:              d.InTime,
		CompanionName:       d.CompanionName,
		AirTime:             d.AirTime,
		BlockTime:           d.BlockTime,
		LogDate:             d.LogDate,
		RouteCode:           d.RouteCode,
		OriginIataCode:      d.OriginIataCode,
		DestinationIataCode: d.DestinationIataCode,
		AirlineCode:         d.AirlineCode,
		TailNumber:          d.TailNumber,
		ModelName:           d.ModelName,
	}

	if d.PilotRole != nil {
		pilotRoleStr := string(*d.PilotRole)
		response.PilotRole = &pilotRoleStr
	}

	if d.CrewRole != nil {
		crewRoleStr := string(*d.CrewRole)
		response.CrewRole = &crewRoleStr
	}

	if d.ApproachType != nil {
		approachTypeStr := string(*d.ApproachType)
		response.ApproachType = &approachTypeStr
	}

	response.FlightType = d.FlightType

	return response
}
