package handlers

import (
	"time"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
)

const dateFormatISO = "2006-01-02"

// DailyLogbookResponse - Response DTO for daily logbook data
type DailyLogbookResponse struct {
	ID           string  `json:"id"`
	LogDate      string  `json:"log_date"`
	EmployeeID   string  `json:"employee_id"`
	BookPage     *int    `json:"book_page,omitempty"`
	Status       string  `json:"status"`
	TailNumberID *string `json:"tail_number_id,omitempty"`
	TailNumber   *string `json:"tail_number,omitempty"`
	CrewRole     *string `json:"crew_role,omitempty"`  // default crew role for every flight logged under this book page
	CreatedAt    string  `json:"created_at,omitempty"` // DB-assigned creation timestamp, not editable by the pilot
	Links        []Link  `json:"_links,omitempty"`
}

// FromDomainDailyLogbook converts domain.DailyLogbook to DailyLogbookResponse with encoded IDs
func FromDomainDailyLogbook(logbook *domain.DailyLogbook, encodedID, encodedEmployeeID string) DailyLogbookResponse {
	status := "inactive"
	if logbook.Status {
		status = "active"
	}
	response := DailyLogbookResponse{
		ID:           encodedID,
		LogDate:      logbook.LogDate.Format(dateFormatISO),
		EmployeeID:   encodedEmployeeID,
		BookPage:     logbook.BookPage,
		Status:       status,
		TailNumberID: logbook.TailNumberID,
		TailNumber:   logbook.TailNumber,
	}
	if !logbook.CreatedAt.IsZero() {
		response.CreatedAt = logbook.CreatedAt.Format(time.RFC3339)
	}
	if logbook.CrewRole != nil {
		crewRole := string(*logbook.CrewRole)
		response.CrewRole = &crewRole
	}
	return response
}

// CreateDailyLogbookRequest - Request DTO for creating a daily logbook
type CreateDailyLogbookRequest struct {
	LogDate      string  `json:"log_date" binding:"required"`
	BookPage     *int    `json:"book_page,omitempty"`
	TailNumberID *string `json:"tail_number_id,omitempty"`
	CrewRole     *string `json:"crew_role,omitempty"` // default crew role for every flight logged under this book page
}

// Sanitize trims whitespace from CreateDailyLogbookRequest fields
func (r *CreateDailyLogbookRequest) Sanitize() {
	r.LogDate = TrimString(r.LogDate)
}

// ToDomain converts the request to a domain model
func (r *CreateDailyLogbookRequest) ToDomain(employeeID string) (*domain.DailyLogbook, error) {
	logDate, err := time.ParseInLocation(dateFormatISO, r.LogDate, time.Local)
	if err != nil {
		return nil, domain.ErrInvalidDateFormat
	}

	logbook := &domain.DailyLogbook{
		LogDate:      logDate,
		EmployeeID:   employeeID,
		BookPage:     r.BookPage,
		Status:       true, // New logbooks are active by default
		TailNumberID: r.TailNumberID,
	}
	if r.CrewRole != nil && *r.CrewRole != "" {
		if !domain.IsValidCrewRole(*r.CrewRole) {
			return nil, domain.ErrInvalidRequest
		}
		crewRole := domain.CrewRole(*r.CrewRole)
		logbook.CrewRole = &crewRole
	}
	logbook.SetID()
	return logbook, nil
}

// UpdateDailyLogbookRequest - Request DTO for updating a daily logbook
type UpdateDailyLogbookRequest struct {
	LogDate      string  `json:"log_date" binding:"required"`
	BookPage     *int    `json:"book_page,omitempty"`
	Status       *bool   `json:"status,omitempty"`
	TailNumberID *string `json:"tail_number_id,omitempty"`
	CrewRole     *string `json:"crew_role,omitempty"`
}

// Sanitize trims whitespace from UpdateDailyLogbookRequest fields
func (r *UpdateDailyLogbookRequest) Sanitize() {
	r.LogDate = TrimString(r.LogDate)
}

// ToDomain converts the request to a domain model for update
func (r *UpdateDailyLogbookRequest) ToDomain(id, employeeID string) (*domain.DailyLogbook, error) {
	logDate, err := time.Parse(dateFormatISO, r.LogDate)
	if err != nil {
		return nil, domain.ErrInvalidDateFormat
	}

	status := true
	if r.Status != nil {
		status = *r.Status
	}

	logbook := &domain.DailyLogbook{
		ID:           id,
		LogDate:      logDate,
		EmployeeID:   employeeID,
		BookPage:     r.BookPage,
		Status:       status,
		TailNumberID: r.TailNumberID,
	}
	if r.CrewRole != nil && *r.CrewRole != "" {
		if !domain.IsValidCrewRole(*r.CrewRole) {
			return nil, domain.ErrInvalidRequest
		}
		crewRole := domain.CrewRole(*r.CrewRole)
		logbook.CrewRole = &crewRole
	}
	return logbook, nil
}

// DailyLogbookListResponse - Response DTO for listing daily logbooks
type DailyLogbookListResponse struct {
	Logbooks []DailyLogbookResponse `json:"daily_logbooks"`
	Total    int                    `json:"total"`
	Links    []Link                 `json:"_links,omitempty"`
}

// ToDailyLogbookListResponse converts a slice of domain.DailyLogbook to DailyLogbookListResponse
func ToDailyLogbookListResponse(logbooks []domain.DailyLogbook, encodeFunc func(string) (string, error), baseURL string) DailyLogbookListResponse {
	response := DailyLogbookListResponse{
		Logbooks: make([]DailyLogbookResponse, 0, len(logbooks)),
		Total:    len(logbooks),
	}

	for _, logbook := range logbooks {
		encodedID, err := encodeFunc(logbook.ID)
		if err != nil {
			// If encoding fails, use the original UUID
			encodedID = logbook.ID
		}
		encodedEmployeeID, err := encodeFunc(logbook.EmployeeID)
		if err != nil {
			encodedEmployeeID = logbook.EmployeeID
		}
		logbookResp := FromDomainDailyLogbook(&logbook, encodedID, encodedEmployeeID)
		if logbookResp.TailNumberID != nil && *logbookResp.TailNumberID != "" {
			if encodedTailNumberID, err := encodeFunc(*logbookResp.TailNumberID); err == nil {
				logbookResp.TailNumberID = &encodedTailNumberID
			}
		}
		// Add HATEOAS links to each logbook
		if baseURL != "" {
			logbookResp.Links = BuildDailyLogbookLinks(baseURL, encodedID)
		}
		response.Logbooks = append(response.Logbooks, logbookResp)
	}

	// Add collection-level links
	if baseURL != "" {
		response.Links = BuildDailyLogbookListLinks(baseURL)
	}

	return response
}

// DailyLogbookStatusResponse - Response DTO for daily logbook status changes
type DailyLogbookStatusResponse struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Updated bool   `json:"updated"`
	Links   []Link `json:"_links,omitempty"`
}

// DailyLogbookDeleteResponse - Response DTO for daily logbook deletion
type DailyLogbookDeleteResponse struct {
	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
	Links   []Link `json:"_links,omitempty"`
}
