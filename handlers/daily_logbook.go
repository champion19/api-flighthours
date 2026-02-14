package handlers

import (
	"time"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
)

// DailyLogbookResponse - Response DTO for daily logbook data
type DailyLogbookResponse struct {
	ID         string `json:"id"`
	LogDate    string `json:"log_date"`
	EmployeeID string `json:"employee_id"`
	BookPage   *int   `json:"book_page,omitempty"`
	Status     string `json:"status"`
	Links      []Link `json:"_links,omitempty"`
}

// FromDomainDailyLogbook converts domain.DailyLogbook to DailyLogbookResponse with encoded IDs
func FromDomainDailyLogbook(logbook *domain.DailyLogbook, encodedID, encodedEmployeeID string) DailyLogbookResponse {
	status := "inactive"
	if logbook.Status {
		status = "active"
	}
	return DailyLogbookResponse{
		ID:         encodedID,
		LogDate:    logbook.LogDate.Format("2006-01-02"),
		EmployeeID: encodedEmployeeID,
		BookPage:   logbook.BookPage,
		Status:     status,
	}
}

// CreateDailyLogbookRequest - Request DTO for creating a daily logbook
type CreateDailyLogbookRequest struct {
	LogDate  string `json:"log_date" binding:"required"`
	BookPage *int   `json:"book_page,omitempty"`
}

// Sanitize trims whitespace from CreateDailyLogbookRequest fields
func (r *CreateDailyLogbookRequest) Sanitize() {
	r.LogDate = TrimString(r.LogDate)
}

// ToDomain converts the request to a domain model
func (r *CreateDailyLogbookRequest) ToDomain(employeeID string) (*domain.DailyLogbook, error) {
	logDate, err := time.ParseInLocation("2006-01-02", r.LogDate, time.Local)
	if err != nil {
		return nil, domain.ErrInvalidDateFormat
	}

	logbook := &domain.DailyLogbook{
		LogDate:    logDate,
		EmployeeID: employeeID,
		BookPage:   r.BookPage,
		Status:     true, // New logbooks are active by default
	}
	logbook.SetID()
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
