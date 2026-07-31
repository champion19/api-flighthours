package domain

import "github.com/google/uuid"

// CrewMember is a person in a pilot's private crew roster (First Officer or
// cabin crew someone has flown with before). Scoped per employee — each
// pilot only sees/searches their own roster.
type CrewMember struct {
	ID         string  `json:"id"`
	EmployeeID string  `json:"employee_id"`
	Name       string  `json:"name"`
	BP         *string `json:"bp,omitempty"` // airline badge/ID number, optional
}

// SetID generates a new UUID for the crew member
func (c *CrewMember) SetID() {
	c.ID = uuid.New().String()
}

// CrewAssignment represents one crew member assigned to a specific flight
// segment (daily_logbook_detail), with the role they performed on that
// flight. The captain is NOT represented here — that's the authenticated
// pilot themselves, already covered by DailyLogbookDetail.CrewRole.
//
// Either CrewMemberID (an existing roster entry) or Name (a brand-new person,
// resolved/created in the same transaction as the flight save) must be set.
type CrewAssignment struct {
	ID           string         `json:"id"`
	CrewMemberID string         `json:"crew_member_id"`
	Name         string         `json:"name"` // set by caller for new people; denormalized from crew_member for display otherwise
	BP           *string        `json:"bp,omitempty"`
	Role         CrewMemberRole `json:"role"`
}
