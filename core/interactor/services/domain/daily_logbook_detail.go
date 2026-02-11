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

// ApproachType represents the type of approach during landing
type ApproachType string

const (
	ApproachTypeNPA    ApproachType = "NPA"    // Non-Precision Approach
	ApproachTypePA     ApproachType = "PA"     // Precision Approach
	ApproachTypeAPV    ApproachType = "APV"    // Approach with Vertical Guidance
	ApproachTypeVisual ApproachType = "VISUAL" // Visual Approach
)

// ValidPilotRoles contains all valid pilot roles
var ValidPilotRoles = []PilotRole{PilotRolePF, PilotRolePM, PilotRolePFTO, PilotRolePFL}

// ValidApproachTypes contains all valid approach types
var ValidApproachTypes = []ApproachType{ApproachTypeNPA, ApproachTypePA, ApproachTypeAPV, ApproachTypeVisual}

// IsValidPilotRole checks if a string is a valid pilot role
func IsValidPilotRole(role string) bool {
	for _, r := range ValidPilotRoles {
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
	ID                           string `json:"id"`
	DailyLogbookID               string `json:"daily_logbook_id"`
	FlightRealDate               string `json:"flight_real_date"` // DATE format: YYYY-MM-DD
	FlightNumber                 string `json:"flight_number"`
	AirlineRouteID               string `json:"airline_route_id"`
	LicensePlateID               string `json:"license_plate_id"`
	Passengers                   *int   `json:"passengers,omitempty"`

	// Flight times (stored as TIME format HH:MM)
	OutTime     string `json:"out_time"`     // Hora salida de bloque (OUT)
	TakeoffTime string `json:"takeoff_time"` // Hora despegue (OFF)
	LandingTime string `json:"landing_time"` // Hora aterrizaje (ON)
	InTime      string `json:"in_time"`      // Hora llegada a bloque (IN)

	// Pilot role and companion
	PilotRole     PilotRole `json:"pilot_role"`
	CompanionName *string   `json:"companion_name,omitempty"`

	// Calculated times (stored as TIME format HH:MM)
	AirTime   string  `json:"air_time"`            // Tiempo de vuelo (ON - OFF)
	BlockTime string  `json:"block_time"`          // Tiempo de bloque (IN - OUT)
	DutyTime  *string `json:"duty_time,omitempty"` // Tiempo de servicio (DUTY)

	// Approach and flight type
	ApproachType      *ApproachType `json:"approach_type,omitempty"`
	FlightType        *string       `json:"flight_type,omitempty"` // 'COMMERCIAL', 'TRAINING', 'FERRY', 'CHECK', 'POSITIONING'
	EmployeeLogbookID *string       `json:"employee_logbook_id,omitempty"`

	// Denormalized fields for display (populated via JOINs)
	// From daily_logbook
	LogDate string `json:"log_date,omitempty"`

	// From airline_route -> route -> airports
	RouteCode           string `json:"route_code,omitempty"`            // e.g., "BOG-CLO"
	OriginIataCode      string `json:"origin_iata_code,omitempty"`      // Origin airport IATA
	DestinationIataCode string `json:"destination_iata_code,omitempty"` // Destination airport IATA
	AirlineCode         string `json:"airline_code,omitempty"`          // Airline IATA code

	// From aircraft_registration -> aircraft_model
	LicensePlate string `json:"license_plate,omitempty"` // Aircraft registration
	ModelName    string `json:"model_name,omitempty"`    // Aircraft model name
}

// SetID generates a new UUID for the detail
func (d *DailyLogbookDetail) SetID() {
	d.ID = uuid.New().String()
}

// ToLogger returns a slice of strings for logging purposes
func (d *DailyLogbookDetail) ToLogger() []string {
	pilotRole := string(d.PilotRole)
	return []string{
		"id:" + d.ID,
		"daily_logbook_id:" + d.DailyLogbookID,
		"flight_number:" + d.FlightNumber,
		"flight_real_date:" + d.FlightRealDate,
		"route_code:" + d.RouteCode,
		"license_plate:" + d.LicensePlate,
		"pilot_role:" + pilotRole,
	}
}
