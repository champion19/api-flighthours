package domain

import "github.com/google/uuid"

// PilotRole represents the role of the pilot in a flight segment
type PilotRole string

const (
	PilotRolePF   PilotRole = "PF"   // Pilot Flying
	PilotRolePM   PilotRole = "PM"   // Pilot Monitoring
	PilotRolePFTO PilotRole = "PFTO" // Pilot Flying Take-Off
	PilotRolePFL  PilotRole = "PFL"  // Pilot Flying Landing
)

// CrewRole represents the crew position (captain or first officer) in a flight segment
type CrewRole string

const (
	CrewRoleCaptain      CrewRole = "captain"       // Capitán
	CrewRoleFirstOfficer CrewRole = "first officer" // Primer Oficial
)

// ApproachType represents the type of approach during landing
type ApproachType string

const (
	ApproachTypeAPV    ApproachType = "APV"    // Approach with Vertical Guidance
	ApproachTypeVisual ApproachType = "VISUAL" // Visual Approach
)

// ValidPilotRoles contains all valid pilot roles
var ValidPilotRoles = []PilotRole{PilotRolePF, PilotRolePM, PilotRolePFTO, PilotRolePFL}

// ValidCrewRoles contains all valid crew roles
var ValidCrewRoles = []CrewRole{CrewRoleCaptain, CrewRoleFirstOfficer}

// ValidApproachTypes contains all valid approach types
var ValidApproachTypes = []ApproachType{ApproachTypeAPV, ApproachTypeVisual}

// IsValidPilotRole checks if a string is a valid pilot role
func IsValidPilotRole(role string) bool {
	for _, r := range ValidPilotRoles {
		if string(r) == role {
			return true
		}
	}
	return false
}

// IsValidCrewRole checks if a string is a valid crew role
func IsValidCrewRole(role string) bool {
	if role == "" {
		return true // NULL is valid
	}
	for _, r := range ValidCrewRoles {
		if string(r) == role {
			return true
		}
	}
	return false
}

// IsValidApproachType checks if a string is a valid approach type
func IsValidApproachType(approachType string) bool {
	if approachType == "" {
		return true // NULL is valid
	}
	for _, at := range ValidApproachTypes {
		if string(at) == approachType {
			return true
		}
	}
	return false
}

// DailyLogbookDetail represents a flight segment within a daily logbook
// This is the CORE entity of the flight hours tracking system
type DailyLogbookDetail struct {
	ID                  string        `json:"id"`
	DailyLogbookID      string        `json:"daily_logbook_id"`
	FlightRealDate      string        `json:"flight_real_date"` // DATE format: YYYY-MM-DD
	FlightNumber        string        `json:"flight_number"`
	AirlineRouteID      string        `json:"airline_route_id"`
	TailNumberID        string        `json:"tail_number_id"`
	Passengers          *int          `json:"passengers,omitempty"`
	OutTime             *string       `json:"out_time,omitempty"`     // Hora salida de bloque (OUT)
	TakeoffTime         *string       `json:"takeoff_time,omitempty"` // Hora despegue (OFF)
	LandingTime         *string       `json:"landing_time,omitempty"` // Hora aterrizaje (ON)
	InTime              *string       `json:"in_time,omitempty"`      // Hora llegada a bloque (IN)
	PilotRole           *PilotRole    `json:"pilot_role,omitempty"`
	CrewRole            *CrewRole     `json:"crew_role,omitempty"`
	CompanionName       *string       `json:"companion_name,omitempty"`
	AirTime             *string       `json:"air_time,omitempty"`   // Tiempo de vuelo (ON - OFF)
	BlockTime           *string       `json:"block_time,omitempty"` // Tiempo de bloque (IN - OUT)
	ApproachType        *ApproachType `json:"approach_type,omitempty"`
	FlightType          *string       `json:"flight_type,omitempty"` // 'COMMERCIAL', 'TRAINING', 'FERRY', 'CHECK', 'POSITIONING'
	EmployeeLogbookID   *string       `json:"employee_logbook_id,omitempty"`
	LogDate             string        `json:"log_date,omitempty"`
	RouteCode           string        `json:"route_code,omitempty"`            // e.g., "BOG-CLO"
	OriginIataCode      string        `json:"origin_iata_code,omitempty"`      // Origin airport IATA
	DestinationIataCode string        `json:"destination_iata_code,omitempty"` // Destination airport IATA
	AirlineCode         string        `json:"airline_code,omitempty"`          // Airline IATA code
	TailNumber          string        `json:"tail_number,omitempty"`           // Aircraft registration
	ModelName           string        `json:"model_name,omitempty"`            // Aircraft model name
}

// SetID generates a new UUID for the detail
func (d *DailyLogbookDetail) SetID() {
	d.ID = uuid.New().String()
}

// ToLogger returns a slice of strings for logging purposes
func (d *DailyLogbookDetail) ToLogger() []string {
	pilotRole := ""
	if d.PilotRole != nil {
		pilotRole = string(*d.PilotRole)
	}
	crewRole := ""
	if d.CrewRole != nil {
		crewRole = string(*d.CrewRole)
	}
	return []string{
		"id:" + d.ID,
		"daily_logbook_id:" + d.DailyLogbookID,
		"flight_number:" + d.FlightNumber,
		"flight_real_date:" + d.FlightRealDate,
		"route_code:" + d.RouteCode,
		"tail_number:" + d.TailNumber,
		"pilot_role:" + pilotRole,
		"crew_role:" + crewRole,
	}
}
