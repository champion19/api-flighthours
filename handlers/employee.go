package handlers

import (
	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
)

// EmployeeRequest - Request DTO for employee registration (core data only)
type EmployeeRequest struct {
	Name                 string `json:"name"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	IdentificationNumber string `json:"identificationNumber"`
	Role                 string `json:"role"`
}

// EmployeeResponse - Response DTO for employee data (core data only)
type EmployeeResponse struct {
	ID                   string `json:"id"`
	Name                 string `json:"name"`
	Email                string `json:"email"`
	IdentificationNumber string `json:"identification_number"`
	Role                 string `json:"role"`
	Links                []Link `json:"_links,omitempty"`
}

// FromDomain converts domain.Employee to EmployeeResponse with encoded ID
func FromDomain(employee *domain.Employee, encodedID string) EmployeeResponse {
	return EmployeeResponse{
		ID:                   encodedID,
		Name:                 employee.Name,
		Email:                employee.Email,
		IdentificationNumber: employee.IdentificationNumber,
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

// UpdateEmployeeRequest - Request DTO for updating employee core data (role is immutable)
type UpdateEmployeeRequest struct {
	Name                 string `json:"name"`
	IdentificationNumber string `json:"identificationNumber"`
}

func (e *EmployeeRequest) Sanitize() {
	e.Name = TrimString(e.Name)
	e.Email = TrimString(e.Email)
	e.IdentificationNumber = TrimString(e.IdentificationNumber)
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
}

// ToUpdateData merges UpdateEmployeeRequest with existing employee data (preserves role)
func (u UpdateEmployeeRequest) ToUpdateData(existing *domain.Employee) domain.Employee {
	return domain.Employee{
		ID:                   existing.ID,
		Name:                 u.Name,
		Email:                existing.Email,
		Password:             existing.Password,
		IdentificationNumber: u.IdentificationNumber,
		Role:                 existing.Role, // Role is immutable - preserve from existing
		KeycloakUserID:       existing.KeycloakUserID,
	}
}

type UpdateEmployeeResponse struct {
	ID      string `json:"id"`
	Updated bool   `json:"updated"`
	Links   []Link `json:"_links,omitempty"`
}

// ToDomain converts EmployeeRequest to domain.Employee
func (e EmployeeRequest) ToDomain() domain.Employee {
	return domain.Employee{
		Name:                 e.Name,
		Email:                e.Email,
		Password:             e.Password,
		IdentificationNumber: e.IdentificationNumber,
		Role:                 e.Role,
	}
}
