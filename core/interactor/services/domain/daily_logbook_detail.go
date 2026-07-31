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

// CrewRole represents the crew position the logged-in pilot flew as on a flight segment
type CrewRole string

const (
	CrewRoleCaptain          CrewRole = "captain"            // Capitán
	CrewRoleFirstOfficer     CrewRole = "first officer"      // Primer Oficial
	CrewRoleInstructor       CrewRole = "instructor"         // Instructor
	CrewRoleLineCheckCaptain CrewRole = "line check captain" // Line Check Captain
	CrewRoleSafetyPilot      CrewRole = "safety pilot"       // Safety Pilot
)

// CrewMemberRole represents the role of another crew member (not the logged-in pilot)
// assigned to a flight segment: the first officer or a member of the cabin crew.
type CrewMemberRole string

const (
	CrewMemberRoleFirstOfficer    CrewMemberRole = "first_officer"    // Primer Oficial
	CrewMemberRolePurser          CrewMemberRole = "purser"           // Jefe de Cabina
	CrewMemberRoleFlightAttendant CrewMemberRole = "flight_attendant" // Tripulante de Cabina
)

// ValidCrewMemberRoles contains all valid crew member roles for daily_logbook_detail_crew assignments
var ValidCrewMemberRoles = []CrewMemberRole{CrewMemberRoleFirstOfficer, CrewMemberRolePurser, CrewMemberRoleFlightAttendant}

// IsValidCrewMemberRole checks if a string is a valid crew member role
func IsValidCrewMemberRole(role string) bool {
	for _, r := range ValidCrewMemberRoles {
		if string(r) == role {
			return true
		}
	}
	return false
}

// ApproachCategory represents the top-level approach procedure used during landing
type ApproachCategory string

const (
	ApproachCategoryRNP    ApproachCategory = "RNP"    // Required Navigation Performance
	ApproachCategoryILS    ApproachCategory = "ILS"    // Instrument Landing System
	ApproachCategoryVisual ApproachCategory = "VISUAL" // Visual Approach
)

// ApproachSubtype represents the specific procedure within an ApproachCategory.
// Only meaningful for RNP and ILS — VISUAL has no subtype.
type ApproachSubtype string

const (
	ApproachSubtypeLNAV        ApproachSubtype = "LNAV"
	ApproachSubtypeLNAVVNAV    ApproachSubtype = "LNAV/VNAV"
	ApproachSubtypeRNPAR03     ApproachSubtype = "RNP AR = 0.3"
	ApproachSubtypeRNPARLT03   ApproachSubtype = "RNP AR < 0.3"
	ApproachSubtypeCatI        ApproachSubtype = "CAT I"
	ApproachSubtypeCatII       ApproachSubtype = "CAT II"
	ApproachSubtypeCatIIIGT175 ApproachSubtype = "CAT III > 175"
	ApproachSubtypeCatIIILT175 ApproachSubtype = "CAT III < 175"
)

// ValidPilotRoles contains all valid pilot roles
var ValidPilotRoles = []PilotRole{PilotRolePF, PilotRolePM, PilotRolePFTO, PilotRolePFL}

// ValidCrewRoles contains all valid crew roles
var ValidCrewRoles = []CrewRole{
	CrewRoleCaptain,
	CrewRoleFirstOfficer,
	CrewRoleInstructor,
	CrewRoleLineCheckCaptain,
	CrewRoleSafetyPilot,
}

// ValidApproachCategories contains all valid approach categories
var ValidApproachCategories = []ApproachCategory{ApproachCategoryRNP, ApproachCategoryILS, ApproachCategoryVisual}

// approachSubtypesByCategory maps each category that has subtypes to its valid subtypes.
// ApproachCategoryVisual is intentionally absent — it has no subtypes.
var approachSubtypesByCategory = map[ApproachCategory][]ApproachSubtype{
	ApproachCategoryRNP: {
		ApproachSubtypeLNAV,
		ApproachSubtypeLNAVVNAV,
		ApproachSubtypeRNPAR03,
		ApproachSubtypeRNPARLT03,
	},
	ApproachCategoryILS: {
		ApproachSubtypeCatI,
		ApproachSubtypeCatII,
		ApproachSubtypeCatIIIGT175,
		ApproachSubtypeCatIIILT175,
	},
}

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

// IsValidApproachCategory checks if a string is a valid approach category
func IsValidApproachCategory(category string) bool {
	if category == "" {
		return true // NULL is valid
	}
	for _, c := range ValidApproachCategories {
		if string(c) == category {
			return true
		}
	}
	return false
}

// IsValidApproachSubtype checks if a subtype string is valid for the given category.
// An empty subtype is valid for VISUAL (or when no category is set) but invalid for RNP/ILS.
func IsValidApproachSubtype(category ApproachCategory, subtype string) bool {
	allowed, hasSubtypes := approachSubtypesByCategory[category]
	if !hasSubtypes {
		return subtype == ""
	}
	for _, s := range allowed {
		if string(s) == subtype {
			return true
		}
	}
	return false
}

// ValidateApproachFields enforces the cross-field rules that a flat enum can't express:
//   - subtype must belong to the selected category (or be empty for VISUAL)
//   - autoland can only be set when the category is ILS
func ValidateApproachFields(category *ApproachCategory, subtype *string, autoland *bool) error {
	if category != nil && !IsValidApproachCategory(string(*category)) {
		return ErrFlightInvalidApproach
	}

	var cat ApproachCategory
	if category != nil {
		cat = *category
	}

	subtypeValue := ""
	if subtype != nil {
		subtypeValue = *subtype
	}
	if !IsValidApproachSubtype(cat, subtypeValue) {
		return ErrFlightInvalidApproach
	}

	if autoland != nil && *autoland && cat != ApproachCategoryILS {
		return ErrFlightInvalidApproach
	}

	return nil
}

// DailyLogbookDetail represents a flight segment within a daily logbook
// This is the CORE entity of the flight hours tracking system
type DailyLogbookDetail struct {
	ID                   string            `json:"id"`
	DailyLogbookID       string            `json:"daily_logbook_id"`
	FlightRealDate       string            `json:"flight_real_date"` // DATE format: YYYY-MM-DD
	FlightNumber         string            `json:"flight_number"`
	OriginAirportID      string            `json:"origin_airport_id"`
	DestinationAirportID string            `json:"destination_airport_id"`
	TailNumberID         string            `json:"tail_number_id"`
	Passengers           *int              `json:"passengers,omitempty"`
	OutTime              *string           `json:"out_time,omitempty"`     // Hora salida de bloque (OUT)
	TakeoffTime          *string           `json:"takeoff_time,omitempty"` // Hora despegue (OFF)
	LandingTime          *string           `json:"landing_time,omitempty"` // Hora aterrizaje (ON)
	InTime               *string           `json:"in_time,omitempty"`      // Hora llegada a bloque (IN)
	PilotRole            *PilotRole        `json:"pilot_role,omitempty"`
	CrewRole             *CrewRole         `json:"crew_role,omitempty"`
	AirTime              *string           `json:"air_time,omitempty"`   // Tiempo de vuelo (ON - OFF)
	BlockTime            *string           `json:"block_time,omitempty"` // Tiempo de bloque (IN - OUT)
	ApproachCategory     *ApproachCategory `json:"approach_category,omitempty"`
	ApproachSubtype      *string           `json:"approach_subtype,omitempty"`
	Autoland             *bool             `json:"autoland,omitempty"`
	FlightType           *string           `json:"flight_type,omitempty"` // 'COMMERCIAL', 'TRAINING', 'FERRY', 'CHECK', 'POSITIONING'
	EmployeeLogbookID    *string           `json:"employee_logbook_id,omitempty"`
	LogDate              string            `json:"log_date,omitempty"`
	RouteCode            string            `json:"route_code,omitempty"`            // e.g., "BOG-CLO"
	OriginIataCode       string            `json:"origin_iata_code,omitempty"`      // Origin airport IATA
	DestinationIataCode  string            `json:"destination_iata_code,omitempty"` // Destination airport IATA
	AirlineCode          string            `json:"airline_code,omitempty"`          // Airline IATA code
	TailNumber           string            `json:"tail_number,omitempty"`           // Aircraft registration
	ModelName            string            `json:"model_name,omitempty"`            // Aircraft model name
	Crew                 []CrewAssignment  `json:"crew,omitempty"`                  // First Officer + cabin crew assigned to this leg
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
