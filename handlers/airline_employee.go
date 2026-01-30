package handlers

import (
	"time"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
)

// AirlineEmployeeResponse - Response DTO for airline employee data (HU26)
type AirlineEmployeeResponse struct {
	ID                   string    `json:"id"`
	Name                 string    `json:"name"`
	AirlineID            string    `json:"airline_id"`
	AirlineName          string    `json:"airline_name"`
	AirlineCode          string    `json:"airline_code"`
	Email                string    `json:"email"`
	IdentificationNumber string    `json:"identification_number"`
	Bp                   string    `json:"bp,omitempty"`
	StartDate            time.Time `json:"start_date"`
	EndDate              time.Time `json:"end_date"`
	Active               bool      `json:"active"`
	Role                 string    `json:"role"`
	Links                []Link    `json:"_links,omitempty"`
}

// FromDomainAirlineEmployee converts domain.AirlineEmployee to AirlineEmployeeResponse with encoded ID
func FromDomainAirlineEmployee(employee *domain.AirlineEmployee, encodedID, encodedAirlineID string) AirlineEmployeeResponse {
	return AirlineEmployeeResponse{
		ID:        encodedID,
		AirlineID: encodedAirlineID,
		Bp:        employee.Bp,
		StartDate: employee.StartDate,
		EndDate:   employee.EndDate,
		Active:    employee.Active,
	}
}

// AirlineEmployeeRequest - Request DTO for creating/updating airline employee
type AirlineEmployeeRequest struct {
	AirlineID string `json:"airline_id" binding:"required"`
	Bp        string `json:"bp"`
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date"`
	Active    bool   `json:"active"`
}

// Sanitize trims whitespace from AirlineEmployeeRequest fields
func (r *AirlineEmployeeRequest) Sanitize() {
	r.AirlineID = TrimString(r.AirlineID)
	r.Bp = TrimString(r.Bp)
	r.StartDate = TrimString(r.StartDate)
	r.EndDate = TrimString(r.EndDate)
}

// ToDomain converts AirlineEmployeeRequest to domain.AirlineEmployee
func (r *AirlineEmployeeRequest) ToDomain() (domain.AirlineEmployee, error) {
	startDate, err := time.Parse("2006-01-02", r.StartDate)
	if err != nil {
		return domain.AirlineEmployee{}, domain.ErrInvalidDateFormat
	}

	var endDate time.Time
	if r.EndDate != "" {
		endDate, err = time.Parse("2006-01-02", r.EndDate)
		if err != nil {
			return domain.AirlineEmployee{}, domain.ErrInvalidDateFormat
		}
	}

	// Validate start date is before end date
	if !endDate.IsZero() && startDate.After(endDate) {
		return domain.AirlineEmployee{}, domain.ErrStartDateAfterEndDate
	}

	return domain.AirlineEmployee{
		AirlineID: r.AirlineID,
		Bp:        r.Bp,
		StartDate: startDate,
		EndDate:   endDate,
		Active:    r.Active,
	}, nil
}

// AirlineEmployeeCreateResponse - Response DTO for created airline employee (HU28)
type AirlineEmployeeCreateResponse struct {
	ID    string `json:"id"`
	Links []Link `json:"_links,omitempty"`
}

// EmployeeAirlineInfoResponse returns airline info for the authenticated employee (HU24)
type EmployeeAirlineInfoResponse struct {
	AirlineID   string `json:"airline_id,omitempty"`
	AirlineName string `json:"airline_name,omitempty"`
	AirlineCode string `json:"airline_code,omitempty"`
	Bp          string `json:"bp,omitempty"`
	Links       []Link `json:"_links,omitempty"`
}

// AddEmployeeAirlineRequest to add airline info for the authenticated employee (HU26)
// Note: 'active' field is not included - it will be managed by activate/deactivate HUs
type AddEmployeeAirlineRequest struct {
	AirlineID string `json:"airline_id" binding:"required"`
	Bp        string `json:"bp"`
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date"`
}

func (u *AddEmployeeAirlineRequest) Sanitize() {
	u.AirlineID = TrimString(u.AirlineID)
	u.Bp = TrimString(u.Bp)
	u.StartDate = TrimString(u.StartDate)
	u.EndDate = TrimString(u.EndDate)
}

// AddEmployeeAirlineResponse returns the result of adding airline info (HU26)
// Note: 'active' is not returned - new employees default to active=true
type AddEmployeeAirlineResponse struct {
	AirlineID   string `json:"airline_id"`
	AirlineName string `json:"airline_name,omitempty"`
	Bp          string `json:"bp,omitempty"`
	StartDate   string `json:"start_date,omitempty"`
	EndDate     string `json:"end_date,omitempty"`
	Links       []Link `json:"_links,omitempty"`
}

// UpdateEmployeeAirlineRequest to update airline info for the authenticated employee (HU25)
type UpdateEmployeeAirlineRequest struct {
	AirlineID string `json:"airline_id" binding:"required"`
	Bp        string `json:"bp"`
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date"`
}

func (u *UpdateEmployeeAirlineRequest) Sanitize() {
	u.AirlineID = TrimString(u.AirlineID)
	u.Bp = TrimString(u.Bp)
	u.StartDate = TrimString(u.StartDate)
	u.EndDate = TrimString(u.EndDate)
}

// UpdateEmployeeAirlineResponse returns the result of updating airline info (HU25)
type UpdateEmployeeAirlineResponse struct {
	Updated     bool   `json:"updated"`
	AirlineID   string `json:"airline_id"`
	AirlineName string `json:"airline_name,omitempty"`
	Bp          string `json:"bp,omitempty"`
	StartDate   string `json:"start_date,omitempty"`
	EndDate     string `json:"end_date,omitempty"`
	Links       []Link `json:"_links,omitempty"`
}

// StatusChangeResponse returns the result of activating/deactivating airline info (HU27/HU28)
type StatusChangeResponse struct {
	Success bool   `json:"success"`
	Active  bool   `json:"active"`
	Links   []Link `json:"_links,omitempty"`
}
