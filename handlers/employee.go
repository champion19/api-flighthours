package handlers

import (
	"time"

	domain "github.com/champion19/flighthours-api/core/interactor/services/domain"
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

// FromDomain convierte un domain.Employee a EmployeeResponse
// Esta función excluye explícitamente el campo Password para no exponerlo en la API
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

// ResendVerificationEmailRequest - DTO para reenviar email de verificación
type ResendVerificationEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResendVerificationEmailResponse - Respuesta de reenvío de email de verificación
type ResendVerificationEmailResponse struct {
	Sent  bool   `json:"sent"`
	Email string `json:"email,omitempty"`
}

// PasswordResetRequest - DTO para solicitar recuperación de contraseña
type PasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// PasswordResetResponse - Respuesta de solicitud de recuperación de contraseña
type PasswordResetResponse struct {
	Sent bool `json:"sent"`
}

// VerifyEmailRequest - DTO para verificar email mediante token proxy
// Este token es un JWT que contiene el email del usuario
type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}

// VerifyEmailResponse - Respuesta de verificación de email
type VerifyEmailResponse struct {
	Verified bool   `json:"verified"`
	Email    string `json:"email,omitempty"`
}

// UpdatePasswordRequest - DTO para actualizar contraseña con token de Keycloak
type UpdatePasswordRequest struct {
	Token           string `json:"token" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=8"`
}

// UpdatePasswordResponse - Respuesta de actualización de contraseña
type UpdatePasswordResponse struct {
	Updated bool   `json:"updated"`
	Email   string `json:"email,omitempty"`
}

// ChangePasswordRequest - DTO para cambiar contraseña cuando el usuario conoce su contraseña actual
// Este flujo no requiere salir de la API, el usuario está autenticado y cambia su contraseña directamente
type ChangePasswordRequest struct {
	Email           string `json:"email" binding:"required,email"`
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=8"`
}

// ChangePasswordResponse - Respuesta de cambio de contraseña
type ChangePasswordResponse struct {
	Changed bool   `json:"changed"`
	Email   string `json:"email,omitempty"`
}

// UpdateEmployeeRequest - DTO para actualizar información general del empleado
// Excluye email y password ya que se manejan en endpoints separados
type UpdateEmployeeRequest struct {
	Name                 string `json:"name"`
	Airline              string `json:"airline"`
	IdentificationNumber string `json:"identificationNumber"`
	Bp                   string `json:"bp"`
	StartDate            string `json:"start_date"`
	EndDate              string `json:"end_date"`
	Active               bool   `json:"active"`
	Role                 string `json:"role"`
}

// ===== SANITIZE METHODS =====
// These methods implement the Sanitizable interface and clean input data
// NOTE: Passwords are NOT trimmed - whitespace could be intentional in passwords

// Sanitize trims whitespace from all string fields in EmployeeRequest
func (e *EmployeeRequest) Sanitize() {
	e.Name = TrimString(e.Name)
	e.Airline = TrimString(e.Airline)
	e.Email = TrimString(e.Email)
	// Password NOT trimmed - whitespace could be intentional
	e.IdentificationNumber = TrimString(e.IdentificationNumber)
	e.Bp = TrimString(e.Bp)
	e.StartDate = TrimString(e.StartDate)
	e.EndDate = TrimString(e.EndDate)
	e.Role = TrimString(e.Role)
}

// Sanitize trims whitespace from LoginRequest fields
func (l *LoginRequest) Sanitize() {
	l.Email = TrimString(l.Email)
	// Password NOT trimmed
}

// Sanitize trims whitespace from ResendVerificationEmailRequest fields
func (r *ResendVerificationEmailRequest) Sanitize() {
	r.Email = TrimString(r.Email)
}

// Sanitize trims whitespace from PasswordResetRequest fields
func (p *PasswordResetRequest) Sanitize() {
	p.Email = TrimString(p.Email)
}

// Sanitize trims whitespace from ChangePasswordRequest fields
func (c *ChangePasswordRequest) Sanitize() {
	c.Email = TrimString(c.Email)
	// Passwords NOT trimmed
}

// Sanitize trims whitespace from UpdateEmployeeRequest fields
func (u *UpdateEmployeeRequest) Sanitize() {
	u.Name = TrimString(u.Name)
	u.Airline = TrimString(u.Airline)
	u.IdentificationNumber = TrimString(u.IdentificationNumber)
	u.Bp = TrimString(u.Bp)
	u.StartDate = TrimString(u.StartDate)
	u.EndDate = TrimString(u.EndDate)
	u.Role = TrimString(u.Role)
}

// ToUpdateData convierte el request a un mapa de campos actualizables
// Se usa con un empleado existente para preservar email, password y keycloak_user_id
// Returns error if date validation fails (invalid format or start_date > end_date)
func (u UpdateEmployeeRequest) ToUpdateData(existing *domain.Employee) (domain.Employee, error) {
	layout := "2006-01-02"

	startDate := existing.StartDate
	if u.StartDate != "" {
		parsed, err := time.Parse(layout, u.StartDate)
		if err != nil {
			return domain.Employee{}, domain.ErrInvalidDateFormat
		}
		startDate = parsed
	}

	endDate := existing.EndDate
	if u.EndDate != "" {
		parsed, err := time.Parse(layout, u.EndDate)
		if err != nil {
			return domain.Employee{}, domain.ErrInvalidDateFormat
		}
		endDate = parsed
	}

	// Validate: start_date cannot be after end_date
	// Only validate if both dates are set (endDate is not zero)
	if !endDate.IsZero() && startDate.After(endDate) {
		return domain.Employee{}, domain.ErrStartDateAfterEndDate
	}

	return domain.Employee{
		ID:                   existing.ID,
		Name:                 u.Name,
		Airline:              u.Airline,
		Email:                existing.Email,    // Preservar email
		Password:             existing.Password, // Preservar password (aunque no está en BD)
		IdentificationNumber: u.IdentificationNumber,
		Bp:                   u.Bp,
		StartDate:            startDate,
		EndDate:              endDate,
		Active:               u.Active,
		Role:                 u.Role,
		KeycloakUserID:       existing.KeycloakUserID, // Preservar keycloak_user_id
	}, nil
}

// UpdateEmployeeResponse - Respuesta de actualización de empleado
type UpdateEmployeeResponse struct {
	ID      string `json:"id"`
	Updated bool   `json:"updated"`
	Links   []Link `json:"_links,omitempty"`
}

func (e EmployeeRequest) ToDomain() (domain.Employee, error) {
	layout := "2006-01-02"

	startDate, err := time.Parse(layout, e.StartDate)
	if err != nil {
		return domain.Employee{}, domain.ErrInvalidDateFormat
	}

	var endDate time.Time
	if e.EndDate != "" {
		endDate, err = time.Parse(layout, e.EndDate)
		if err != nil {
			return domain.Employee{}, domain.ErrInvalidDateFormat
		}

		// Validate: start_date cannot be after end_date
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
