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

type Airline struct {
	ID          string `json:"id"`
	AirlineName string `json:"airline_name"`
	AirlineCode string `json:"airline_code"`
	Status      string `json:"status"`
}

func (a *Airline) ToLogger() []string {
	return []string{
		"id:" + a.ID,
		"name:" + a.AirlineName,
		"code:" + a.AirlineCode,
		"status:" + a.Status,
	}
}

func (a *Airline) IsActive() bool {
	return a.Status == "active"
}



