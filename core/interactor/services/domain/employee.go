package domain

import (
	"time"

	"github.com/google/uuid"

)

type Employee struct {
	ID                   string    `json:"id"`
	Name                 string    `json:"name"`
	Airline              string    `json:"airline"`
	Email                string    `json:"email"`
	Password             string    `json:"password"`
	IdentificationNumber string    `json:"identification_number"`
	Bp                   string    `json:"bp"`
	StartDate            time.Time `json:"start_date"`
	EndDate              time.Time `json:"end_date"`
	Active               bool      `json:"active"`
	Role                 string    `json:"role,omitempty"`
	KeycloakUserID       string    `json:"keycloak_user_id,omitempty"`
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



