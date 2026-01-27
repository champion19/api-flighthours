package handlers

import (
	"time"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
)

type EmployeeRequest struct {
	Name                 string `json:"name"`
	Airline              string `json:"airline"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	IdentificationNumber string `json:"identificationNumber"`
	Bp                   string `json:"bp"`
	StartDate            string `json:"start_date"`
	EndDate              string `json:"end_date"`
	Active               bool   `json:"active"`
	Role                 string `json:"role"`
}

type EmployeeResponse struct {
	ID                   string    `json:"id"`
	Name                 string    `json:"name"`
	Airline              string    `json:"airline,omitempty"`
	Email                string    `json:"email"`
	IdentificationNumber string    `json:"identification_number"`
	Bp                   string    `json:"bp,omitempty"`
	StartDate            time.Time `json:"start_date"`
	EndDate              time.Time `json:"end_date"`
	Active               bool      `json:"active"`
	Role                 string    `json:"role"`
	Links                []Link    `json:"_links,omitempty"`
}

func FromDomain(employee *domain.Employee, encodedID string) EmployeeResponse {
	return EmployeeResponse{
		ID:                   encodedID,
		Name:                 employee.Name,
		Airline:              employee.Airline,
		Email:                employee.Email,
		IdentificationNumber: employee.IdentificationNumber,
		Bp:                   employee.Bp,
		StartDate:            employee.StartDate,
		EndDate:              employee.EndDate,
		Active:               employee.Active,
		Role:                 employee.Role,
	}
}

type RegisterEmployeeResponse struct {
	Links []Link `json:"_links"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type ResendVerificationEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResendVerificationEmailResponse struct {
	Sent  bool   `json:"sent"`
	Email string `json:"email,omitempty"`
}

type PasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type PasswordResetResponse struct {
	Sent bool `json:"sent"`
}

type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}

type VerifyEmailResponse struct {
	Verified bool   `json:"verified"`
	Email    string `json:"email,omitempty"`
}

type UpdatePasswordRequest struct {
	Token           string `json:"token" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=8"`
}

type UpdatePasswordResponse struct {
	Updated bool   `json:"updated"`
	Email   string `json:"email,omitempty"`
}

type ChangePasswordRequest struct {
	Email           string `json:"email" binding:"required,email"`
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=8"`
}

type ChangePasswordResponse struct {
	Changed bool   `json:"changed"`
	Email   string `json:"email,omitempty"`
}

type UpdateEmployeeRequest struct {
	Name                 string `json:"name"`
	IdentificationNumber string `json:"identificationNumber"`
	StartDate            string `json:"start_date"`
	EndDate              string `json:"end_date"`
	Active               bool   `json:"active"`
	Role                 string `json:"role"`
}

func (e *EmployeeRequest) Sanitize() {
	e.Name = TrimString(e.Name)
	e.Airline = TrimString(e.Airline)
	e.Email = TrimString(e.Email)
	e.IdentificationNumber = TrimString(e.IdentificationNumber)
	e.Bp = TrimString(e.Bp)
	e.StartDate = TrimString(e.StartDate)
	e.EndDate = TrimString(e.EndDate)
	e.Role = TrimString(e.Role)
}

func (l *LoginRequest) Sanitize() {
	l.Email = TrimString(l.Email)
}

func (r *ResendVerificationEmailRequest) Sanitize() {
	r.Email = TrimString(r.Email)
}

func (p *PasswordResetRequest) Sanitize() {
	p.Email = TrimString(p.Email)
}

func (c *ChangePasswordRequest) Sanitize() {
	c.Email = TrimString(c.Email)
}

func (u *UpdateEmployeeRequest) Sanitize() {
	u.Name = TrimString(u.Name)
	u.IdentificationNumber = TrimString(u.IdentificationNumber)
	u.StartDate = TrimString(u.StartDate)
	u.EndDate = TrimString(u.EndDate)
	u.Role = TrimString(u.Role)
}

func (u UpdateEmployeeRequest) ToUpdateData(existing *domain.Employee) (domain.Employee, error) {
	layout := "2006-01-02"
	loc := time.Local 

	startDate := existing.StartDate
	if u.StartDate != "" {
		parsed, err := time.ParseInLocation(layout, u.StartDate, loc)
		if err != nil {
			return domain.Employee{}, domain.ErrInvalidDateFormat
		}
		startDate = parsed
	}

	endDate := existing.EndDate
	if u.EndDate != "" {
		parsed, err := time.ParseInLocation(layout, u.EndDate, loc)
		if err != nil {
			return domain.Employee{}, domain.ErrInvalidDateFormat
		}
		endDate = parsed
	}

	if !endDate.IsZero() && startDate.After(endDate) {
		return domain.Employee{}, domain.ErrStartDateAfterEndDate
	}

	return domain.Employee{
		ID:                   existing.ID,
		Name:                 u.Name,
		Airline:              existing.Airline,
		Email:                existing.Email,
		Password:             existing.Password,
		IdentificationNumber: u.IdentificationNumber,
		Bp:                   existing.Bp,
		StartDate:            startDate,
		EndDate:              endDate,
		Active:               u.Active,
		Role:                 u.Role,
		KeycloakUserID:       existing.KeycloakUserID,
	}, nil
}

type UpdateEmployeeResponse struct {
	ID      string `json:"id"`
	Updated bool   `json:"updated"`
	Links   []Link `json:"_links,omitempty"`
}

func (e EmployeeRequest) ToDomain() (domain.Employee, error) {
	layout := "2006-01-02"
	loc := time.Local

	startDate, err := time.ParseInLocation(layout, e.StartDate, loc)
	if err != nil {
		return domain.Employee{}, domain.ErrInvalidDateFormat
	}

	var endDate time.Time
	if e.EndDate != "" {
		endDate, err = time.ParseInLocation(layout, e.EndDate, loc)
		if err != nil {
			return domain.Employee{}, domain.ErrInvalidDateFormat
		}

		if startDate.After(endDate) {
			return domain.Employee{}, domain.ErrStartDateAfterEndDate
		}
	}

	return domain.Employee{
		Name:                 e.Name,
		Airline:              e.Airline,
		Email:                e.Email,
		Password:             e.Password,
		IdentificationNumber: e.IdentificationNumber,
		Bp:                   e.Bp,
		StartDate:            startDate,
		EndDate:              endDate,
		Active:               e.Active,
		Role:                 e.Role,
	}, nil
}
