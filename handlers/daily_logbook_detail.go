package handlers

import (
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

// CrewAssignmentRequest is one command crew/cabin crew assignment sent as part of a flight leg payload.
// Either CrewMemberID (an existing roster entry, picked from search) or Name (a brand-new
// person, created in the same save transaction) must be set.
type CrewAssignmentRequest struct {
	CrewMemberID string `json:"crew_member_id,omitempty"`
	Name         string `json:"name,omitempty"`
	BP           string `json:"bp,omitempty"`
	Role         string `json:"role"` // captain, first officer, instructor, line check captain, safety pilot (Tripulación de Mando); purser, flight_attendant (Tripulación de Cabina)
}

type CreateDailyLogbookDetailRequest struct {
	FlightRealDate       string                  `json:"flight_real_date"`
	FlightNumber         string                  `json:"flight_number"`
	OriginAirportID      string                  `json:"origin_airport_id"`
	DestinationAirportID string                  `json:"destination_airport_id"`
	TailNumberID         string                  `json:"tail_number_id"`
	Passengers           *int                    `json:"passengers,omitempty"`
	OutTime              *string                 `json:"out_time,omitempty"`     // TIME format HH:MM (nullable)
	TakeoffTime          *string                 `json:"takeoff_time,omitempty"` // TIME format HH:MM (nullable)
	LandingTime          *string                 `json:"landing_time,omitempty"` // TIME format HH:MM (nullable)
	InTime               *string                 `json:"in_time,omitempty"`      // TIME format HH:MM (nullable)
	PilotRole            *string                 `json:"pilot_role,omitempty"`   // nullable
	CrewRole             *string                 `json:"crew_role,omitempty"`    // nullable
	AirTime              *string                 `json:"air_time,omitempty"`   // TIME format HH:MM (nullable)
	BlockTime            *string                 `json:"block_time,omitempty"` // TIME format HH:MM (nullable)
	ApproachCategory     *string                 `json:"approach_category,omitempty"`
	ApproachSubtype      *string                 `json:"approach_subtype,omitempty"`
	Autoland             *bool                   `json:"autoland,omitempty"`
	FlightType           *string                 `json:"flight_type,omitempty"`
	Crew                 []CrewAssignmentRequest `json:"crew,omitempty"` // command crew + cabin crew — optional
}

func (r *CreateDailyLogbookDetailRequest) Sanitize() {
	r.FlightRealDate = TrimString(r.FlightRealDate)
	r.FlightNumber = TrimString(r.FlightNumber)
	r.OriginAirportID = TrimString(r.OriginAirportID)
	r.DestinationAirportID = TrimString(r.DestinationAirportID)
	r.TailNumberID = TrimString(r.TailNumberID)
	r.OutTime = TrimStringPtr(r.OutTime)
	r.TakeoffTime = TrimStringPtr(r.TakeoffTime)
	r.LandingTime = TrimStringPtr(r.LandingTime)
	r.InTime = TrimStringPtr(r.InTime)
	r.PilotRole = TrimStringPtr(r.PilotRole)
	r.CrewRole = TrimStringPtr(r.CrewRole)
	r.AirTime = TrimStringPtr(r.AirTime)
	r.BlockTime = TrimStringPtr(r.BlockTime)
	r.ApproachCategory = TrimStringPtr(r.ApproachCategory)
	r.ApproachSubtype = TrimStringPtr(r.ApproachSubtype)
	r.FlightType = TrimStringPtr(r.FlightType)
}

type UpdateDailyLogbookDetailRequest struct {
	FlightRealDate       string                  `json:"flight_real_date"`
	FlightNumber         string                  `json:"flight_number"`
	OriginAirportID      string                  `json:"origin_airport_id"`
	DestinationAirportID string                  `json:"destination_airport_id"`
	TailNumberID         string                  `json:"tail_number_id"`
	Passengers           *int                    `json:"passengers,omitempty"`
	OutTime              *string                 `json:"out_time,omitempty"`     // TIME format HH:MM (nullable)
	TakeoffTime          *string                 `json:"takeoff_time,omitempty"` // TIME format HH:MM (nullable)
	LandingTime          *string                 `json:"landing_time,omitempty"` // TIME format HH:MM (nullable)
	InTime               *string                 `json:"in_time,omitempty"`      // TIME format HH:MM (nullable)
	PilotRole            *string                 `json:"pilot_role,omitempty"`   // nullable
	CrewRole             *string                 `json:"crew_role,omitempty"`    // nullable
	AirTime              *string                 `json:"air_time,omitempty"`   // TIME format HH:MM (nullable)
	BlockTime            *string                 `json:"block_time,omitempty"` // TIME format HH:MM (nullable)
	ApproachCategory     *string                 `json:"approach_category,omitempty"`
	ApproachSubtype      *string                 `json:"approach_subtype,omitempty"`
	Autoland             *bool                   `json:"autoland,omitempty"`
	FlightType           *string                 `json:"flight_type,omitempty"`
	Crew                 []CrewAssignmentRequest `json:"crew,omitempty"` // command crew + cabin crew — optional; empty array clears
}

func (r *UpdateDailyLogbookDetailRequest) Sanitize() {
	r.FlightRealDate = TrimString(r.FlightRealDate)
	r.FlightNumber = TrimString(r.FlightNumber)
	r.OriginAirportID = TrimString(r.OriginAirportID)
	r.DestinationAirportID = TrimString(r.DestinationAirportID)
	r.TailNumberID = TrimString(r.TailNumberID)
	r.OutTime = TrimStringPtr(r.OutTime)
	r.TakeoffTime = TrimStringPtr(r.TakeoffTime)
	r.LandingTime = TrimStringPtr(r.LandingTime)
	r.InTime = TrimStringPtr(r.InTime)
	r.PilotRole = TrimStringPtr(r.PilotRole)
	r.CrewRole = TrimStringPtr(r.CrewRole)
	r.AirTime = TrimStringPtr(r.AirTime)
	r.BlockTime = TrimStringPtr(r.BlockTime)
	r.ApproachCategory = TrimStringPtr(r.ApproachCategory)
	r.ApproachSubtype = TrimStringPtr(r.ApproachSubtype)
	r.FlightType = TrimStringPtr(r.FlightType)
}

// ============================================
// RESPONSE DTOs
// ============================================

// CrewAssignmentResponse is one command crew/cabin crew assignment on a flight leg response.
type CrewAssignmentResponse struct {
	ID           string  `json:"id"`
	CrewMemberID string  `json:"crew_member_id"`
	Name         string  `json:"name"`
	BP           *string `json:"bp,omitempty"`
	Role         string  `json:"role"`
}

// DailyLogbookDetailResponse represents the response for a detail
type DailyLogbookDetailResponse struct {
	ID                   string                   `json:"id"`
	DailyLogbookID       string                   `json:"daily_logbook_id"`
	FlightRealDate       string                   `json:"flight_real_date"`
	FlightNumber         string                   `json:"flight_number"`
	OriginAirportID      string                   `json:"origin_airport_id"`
	DestinationAirportID string                   `json:"destination_airport_id"`
	TailNumberID         string                   `json:"tail_number_id"`
	Passengers           *int                     `json:"passengers,omitempty"`
	OutTime              *string                  `json:"out_time,omitempty"`
	TakeoffTime          *string                  `json:"takeoff_time,omitempty"`
	LandingTime          *string                  `json:"landing_time,omitempty"`
	InTime               *string                  `json:"in_time,omitempty"`
	PilotRole            *string                  `json:"pilot_role,omitempty"`
	CrewRole             *string                  `json:"crew_role,omitempty"`
	AirTime              *string                  `json:"air_time,omitempty"`
	BlockTime            *string                  `json:"block_time,omitempty"`
	ApproachCategory     *string                  `json:"approach_category,omitempty"`
	ApproachSubtype      *string                  `json:"approach_subtype,omitempty"`
	Autoland             *bool                    `json:"autoland,omitempty"`
	FlightType           *string                  `json:"flight_type,omitempty"`
	LogDate              string                   `json:"log_date,omitempty"`
	RouteCode            string                   `json:"route_code,omitempty"`
	OriginIataCode       string                   `json:"origin_iata_code,omitempty"`
	DestinationIataCode  string                   `json:"destination_iata_code,omitempty"`
	AirlineCode          string                   `json:"airline_code,omitempty"`
	TailNumber           string                   `json:"tail_number,omitempty"`
	ModelName            string                   `json:"model_name,omitempty"`
	Crew                 []CrewAssignmentResponse `json:"crew,omitempty"`
	Links                []Link                   `json:"_links,omitempty"`
}

// ============================================
// MAPPERS
// ============================================

// ToDomainDailyLogbookDetail converts a create request to domain model
func ToDomainDailyLogbookDetail(logbookID string, req CreateDailyLogbookDetailRequest) domain.DailyLogbookDetail {
	detail := domain.DailyLogbookDetail{
		DailyLogbookID:       logbookID,
		FlightRealDate:       req.FlightRealDate,
		FlightNumber:         req.FlightNumber,
		OriginAirportID:      req.OriginAirportID,
		DestinationAirportID: req.DestinationAirportID,
		TailNumberID:         req.TailNumberID,
		Passengers:           req.Passengers,
		OutTime:              req.OutTime,
		TakeoffTime:          req.TakeoffTime,
		LandingTime:          req.LandingTime,
		InTime:               req.InTime,
		AirTime:              req.AirTime,
		BlockTime:            req.BlockTime,
		FlightType:           req.FlightType,
	}

	if req.PilotRole != nil {
		pilotRole := domain.PilotRole(*req.PilotRole)
		detail.PilotRole = &pilotRole
	}

	if req.CrewRole != nil {
		crewRole := domain.CrewRole(*req.CrewRole)
		detail.CrewRole = &crewRole
	}

	if req.ApproachCategory != nil {
		approachCategory := domain.ApproachCategory(*req.ApproachCategory)
		detail.ApproachCategory = &approachCategory
	}
	detail.ApproachSubtype = req.ApproachSubtype
	detail.Autoland = req.Autoland
	detail.Crew = toDomainCrewAssignments(req.Crew)

	return detail
}

// ToDomainDailyLogbookDetailUpdate converts an update request to domain model
func ToDomainDailyLogbookDetailUpdate(id string, req UpdateDailyLogbookDetailRequest) domain.DailyLogbookDetail {
	detail := domain.DailyLogbookDetail{
		ID:                   id,
		FlightRealDate:       req.FlightRealDate,
		FlightNumber:         req.FlightNumber,
		OriginAirportID:      req.OriginAirportID,
		DestinationAirportID: req.DestinationAirportID,
		TailNumberID:         req.TailNumberID,
		Passengers:           req.Passengers,
		OutTime:              req.OutTime,
		TakeoffTime:          req.TakeoffTime,
		LandingTime:          req.LandingTime,
		InTime:               req.InTime,
		AirTime:              req.AirTime,
		BlockTime:            req.BlockTime,
		FlightType:           req.FlightType,
	}

	if req.PilotRole != nil {
		pilotRole := domain.PilotRole(*req.PilotRole)
		detail.PilotRole = &pilotRole
	}

	if req.CrewRole != nil {
		crewRole := domain.CrewRole(*req.CrewRole)
		detail.CrewRole = &crewRole
	}

	if req.ApproachCategory != nil {
		approachCategory := domain.ApproachCategory(*req.ApproachCategory)
		detail.ApproachCategory = &approachCategory
	}
	detail.ApproachSubtype = req.ApproachSubtype
	detail.Autoland = req.Autoland
	detail.Crew = toDomainCrewAssignments(req.Crew)

	return detail
}

// toDomainCrewAssignments converts request crew rows to domain assignments.
// CrewMemberID is expected to already be resolved to a real UUID by the controller
// (same pattern as origin/destination airport IDs) before this is called — rows with
// no CrewMemberID carry Name/BP instead, to be resolved/created by the interactor
// within the same save transaction.
func toDomainCrewAssignments(rows []CrewAssignmentRequest) []domain.CrewAssignment {
	if rows == nil {
		return nil
	}
	assignments := make([]domain.CrewAssignment, 0, len(rows))
	for _, row := range rows {
		assignment := domain.CrewAssignment{
			CrewMemberID: row.CrewMemberID,
			Name:         row.Name,
			Role:         domain.CrewMemberRole(row.Role),
		}
		if row.BP != "" {
			bp := row.BP
			assignment.BP = &bp
		}
		assignments = append(assignments, assignment)
	}
	return assignments
}

// FromDomainDailyLogbookDetail converts domain model to response DTO
func FromDomainDailyLogbookDetail(d *domain.DailyLogbookDetail, encodedID, encodedLogbookID, encodedOriginAirportID, encodedDestinationAirportID, encodedAircraftID string) DailyLogbookDetailResponse {
	response := DailyLogbookDetailResponse{
		ID:                   encodedID,
		DailyLogbookID:       encodedLogbookID,
		FlightRealDate:       d.FlightRealDate,
		FlightNumber:         d.FlightNumber,
		OriginAirportID:      encodedOriginAirportID,
		DestinationAirportID: encodedDestinationAirportID,
		TailNumberID:         encodedAircraftID,
		Passengers:           d.Passengers,
		OutTime:              d.OutTime,
		TakeoffTime:          d.TakeoffTime,
		LandingTime:          d.LandingTime,
		InTime:               d.InTime,
		AirTime:              d.AirTime,
		BlockTime:            d.BlockTime,
		LogDate:              d.LogDate,
		RouteCode:            d.RouteCode,
		OriginIataCode:       d.OriginIataCode,
		DestinationIataCode:  d.DestinationIataCode,
		AirlineCode:          d.AirlineCode,
		TailNumber:           d.TailNumber,
		ModelName:            d.ModelName,
	}

	if d.PilotRole != nil {
		pilotRoleStr := string(*d.PilotRole)
		response.PilotRole = &pilotRoleStr
	}

	if d.CrewRole != nil {
		crewRoleStr := string(*d.CrewRole)
		response.CrewRole = &crewRoleStr
	}

	if d.ApproachCategory != nil {
		approachCategoryStr := string(*d.ApproachCategory)
		response.ApproachCategory = &approachCategoryStr
	}
	response.ApproachSubtype = d.ApproachSubtype
	response.Autoland = d.Autoland

	response.FlightType = d.FlightType

	if d.Crew != nil {
		crew := make([]CrewAssignmentResponse, 0, len(d.Crew))
		for _, a := range d.Crew {
			crew = append(crew, CrewAssignmentResponse{
				ID:           a.ID,
				CrewMemberID: a.CrewMemberID, // encoded by the caller (controller) before returning to the client
				Name:         a.Name,
				BP:           a.BP,
				Role:         string(a.Role),
			})
		}
		response.Crew = crew
	}

	return response
}
