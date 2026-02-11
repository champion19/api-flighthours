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
	OutTime        string  `json:"out_time"`     // TIME format HH:MM
	TakeoffTime    string  `json:"takeoff_time"` // TIME format HH:MM
	LandingTime    string  `json:"landing_time"` // TIME format HH:MM
	InTime         string  `json:"in_time"`      // TIME format HH:MM
	PilotRole      string  `json:"pilot_role"`
	CompanionName  *string `json:"companion_name,omitempty"`
	AirTime        string  `json:"air_time"`            // TIME format HH:MM
	BlockTime      string  `json:"block_time"`          // TIME format HH:MM
	DutyTime       *string `json:"duty_time,omitempty"` // TIME format HH:MM (nullable)
	ApproachType   *string `json:"approach_type,omitempty"`
	FlightType     *string `json:"flight_type,omitempty"`
}

func (r *CreateDailyLogbookDetailRequest) Sanitize() {
	r.FlightRealDate = TrimString(r.FlightRealDate)
	r.FlightNumber = TrimString(r.FlightNumber)
	r.AirlineRouteID = TrimString(r.AirlineRouteID)
	r.LicensePlateID = TrimString(r.LicensePlateID)
	r.OutTime = TrimString(r.OutTime)
	r.TakeoffTime = TrimString(r.TakeoffTime)
	r.LandingTime = TrimString(r.LandingTime)
	r.InTime = TrimString(r.InTime)
	r.PilotRole = TrimString(r.PilotRole)
	r.CompanionName = TrimStringPtr(r.CompanionName)
	r.AirTime = TrimString(r.AirTime)
	r.BlockTime = TrimString(r.BlockTime)
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
	OutTime        string  `json:"out_time"`     // TIME format HH:MM
	TakeoffTime    string  `json:"takeoff_time"` // TIME format HH:MM
	LandingTime    string  `json:"landing_time"` // TIME format HH:MM
	InTime         string  `json:"in_time"`      // TIME format HH:MM
	PilotRole      string  `json:"pilot_role"`
	CompanionName  *string `json:"companion_name,omitempty"`
	AirTime        string  `json:"air_time"`            // TIME format HH:MM
	BlockTime      string  `json:"block_time"`          // TIME format HH:MM
	DutyTime       *string `json:"duty_time,omitempty"` // TIME format HH:MM (nullable)
	ApproachType   *string `json:"approach_type,omitempty"`
	FlightType     *string `json:"flight_type,omitempty"`
}

func (r *UpdateDailyLogbookDetailRequest) Sanitize() {
	r.FlightRealDate = TrimString(r.FlightRealDate)
	r.FlightNumber = TrimString(r.FlightNumber)
	r.AirlineRouteID = TrimString(r.AirlineRouteID)
	r.LicensePlateID = TrimString(r.LicensePlateID)
	r.OutTime = TrimString(r.OutTime)
	r.TakeoffTime = TrimString(r.TakeoffTime)
	r.LandingTime = TrimString(r.LandingTime)
	r.InTime = TrimString(r.InTime)
	r.PilotRole = TrimString(r.PilotRole)
	r.CompanionName = TrimStringPtr(r.CompanionName)
	r.AirTime = TrimString(r.AirTime)
	r.BlockTime = TrimString(r.BlockTime)
	r.DutyTime = TrimStringPtr(r.DutyTime)
	r.ApproachType = TrimStringPtr(r.ApproachType)
	r.FlightType = TrimStringPtr(r.FlightType)
}

// DailyLogbookDetailResponse represents the response for a detail
type DailyLogbookDetailResponse struct {
	ID                  string            `json:"id"`
	DailyLogbookID      string            `json:"daily_logbook_id"`
	FlightRealDate      string            `json:"flight_real_date"`
	FlightNumber        string            `json:"flight_number"`
	AirlineRouteID      string            `json:"airline_route_id"`
	LicensePlateID      string            `json:"license_plate_id"`
	Passengers          *int              `json:"passengers,omitempty"`
	OutTime             string            `json:"out_time"`
	TakeoffTime         string            `json:"takeoff_time"`
	LandingTime         string            `json:"landing_time"`
	InTime              string            `json:"in_time"`
	PilotRole           string            `json:"pilot_role"`
	CompanionName       *string           `json:"companion_name,omitempty"`
	AirTime             string            `json:"air_time"`
	BlockTime           string            `json:"block_time"`
	DutyTime            *string           `json:"duty_time,omitempty"`
	ApproachType        *string           `json:"approach_type,omitempty"`
	FlightType          *string           `json:"flight_type,omitempty"`
	LogDate             string            `json:"log_date,omitempty"`
	RouteCode           string            `json:"route_code,omitempty"`
	OriginIataCode      string            `json:"origin_iata_code,omitempty"`
	DestinationIataCode string            `json:"destination_iata_code,omitempty"`
	AirlineCode         string            `json:"airline_code,omitempty"`
	LicensePlate        string            `json:"license_plate,omitempty"`
	ModelName           string            `json:"model_name,omitempty"`
	Links               map[string]string `json:"_links,omitempty"`
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
		PilotRole:      domain.PilotRole(req.PilotRole),
		CompanionName:  req.CompanionName,
		AirTime:        req.AirTime,
		BlockTime:      req.BlockTime,
		DutyTime:       req.DutyTime,
		FlightType:     req.FlightType,
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
		PilotRole:      domain.PilotRole(req.PilotRole),
		CompanionName:  req.CompanionName,
		AirTime:        req.AirTime,
		BlockTime:      req.BlockTime,
		DutyTime:       req.DutyTime,
		FlightType:     req.FlightType,
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
		PilotRole:           string(d.PilotRole),
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

	if d.ApproachType != nil {
		approachTypeStr := string(*d.ApproachType)
		response.ApproachType = &approachTypeStr
	}

	response.FlightType = d.FlightType

	return response
}
