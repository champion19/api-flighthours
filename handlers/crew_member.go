package handlers

import (
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

// CreateCrewMemberRequest is the payload to add a person to the pilot's own crew roster.
type CreateCrewMemberRequest struct {
	Name string `json:"name"`
	BP   string `json:"bp,omitempty"`
}

func (r *CreateCrewMemberRequest) Sanitize() {
	r.Name = TrimString(r.Name)
	r.BP = TrimString(r.BP)
}

// CrewMemberResponse represents a single person in a pilot's crew roster.
type CrewMemberResponse struct {
	ID   string  `json:"id"`
	Name string  `json:"name"`
	BP   *string `json:"bp,omitempty"`
}

// FromDomainCrewMember converts a domain crew member to its response DTO.
func FromDomainCrewMember(m *domain.CrewMember, encodedID string) CrewMemberResponse {
	return CrewMemberResponse{
		ID:   encodedID,
		Name: m.Name,
		BP:   m.BP,
	}
}

// resolveCrewAssignments converts obfuscated crew_member_id values in a flight leg
// request's crew list to real UUIDs — same resolution already applied to origin/
// destination airport IDs. Preserves nil (field omitted) vs non-nil-empty (field
// sent as []), since that distinction decides whether ReplaceCrewForDetail runs.
// Rows with no crew_member_id (a brand-new person, identified only by name) are
// passed through untouched — they get created inside the same save transaction.
func (h *handler) resolveCrewAssignments(rows []CrewAssignmentRequest) []CrewAssignmentRequest {
	if rows == nil {
		return nil
	}
	resolved := make([]CrewAssignmentRequest, 0, len(rows))
	for _, row := range rows {
		if row.CrewMemberID == "" {
			resolved = append(resolved, row)
			continue
		}
		crewMemberUUID, _ := h.resolveID(row.CrewMemberID)
		resolved = append(resolved, CrewAssignmentRequest{CrewMemberID: crewMemberUUID, Role: row.Role})
	}
	return resolved
}

// encodeCrewAssignments obfuscates crew_member_id values in a flight leg response's crew list.
func (h *handler) encodeCrewAssignments(rows []CrewAssignmentResponse) []CrewAssignmentResponse {
	if rows == nil {
		return nil
	}
	encoded := make([]CrewAssignmentResponse, 0, len(rows))
	for _, row := range rows {
		encodedID, _ := h.EncodeID(row.CrewMemberID)
		row.CrewMemberID = encodedID
		encoded = append(encoded, row)
	}
	return encoded
}
