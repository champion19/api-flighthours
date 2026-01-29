package domain

import (
	"github.com/google/uuid"
)

// Employee represents core user data (authentication, identity)
// Airline-specific data is handled by AirlineEmployee
type Employee struct {
	ID                   string `json:"id"`
	Name                 string `json:"name"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	IdentificationNumber string `json:"identification_number"`
	Role                 string `json:"role,omitempty"`
	KeycloakUserID       string `json:"keycloak_user_id,omitempty"`
}

func (e *Employee) SetID() {
	e.ID = uuid.New().String()
}

func (e *Employee) ToLogger() []string {
	return []string{
		"id:" + e.ID,
		"email:" + e.Email,
		"role:" + e.Role,
	}
}
