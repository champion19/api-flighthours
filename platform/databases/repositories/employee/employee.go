package employee

import (
	"time"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

type Employee struct {
	ID                   string     `db:"id"`
	Name                 string     `db:"name"`
	Airline              *string    `db:"airline"`
	Email                string     `db:"email"`
	IdentificationNumber string     `db:"identification_number"`
	Bp                   *string    `db:"bp"`
	StartDate            *time.Time `db:"start_date"`
	EndDate              *time.Time `db:"end_date"`
	Active               bool       `db:"active"`
	Role                 string     `db:"role"`
	KeycloakUserID       *string    `db:"keycloak_user_id"`
}


func (e Employee) ToDomain() domain.Employee {
	keycloakUserID := ""
	if e.KeycloakUserID != nil {
		keycloakUserID = *e.KeycloakUserID
	}

	return domain.Employee{
		ID:                   e.ID,
		Name:                 e.Name,
		Email:                e.Email,
		IdentificationNumber: e.IdentificationNumber,
		Role:                 e.Role,
		KeycloakUserID:       keycloakUserID,
	}
}

func stringPtrOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}


func FromDomain(domainEmployee domain.Employee) Employee {
	return Employee{
		ID:                   domainEmployee.ID,
		Name:                 domainEmployee.Name,
		Airline:              nil,
		Email:                domainEmployee.Email,
		IdentificationNumber: domainEmployee.IdentificationNumber,
		Bp:                   nil,
		StartDate:            nil, // NULL - dates are optional at registration
		EndDate:              nil, // NULL - dates are optional at registration
		Active:               false,
		Role:                 domainEmployee.Role,
		KeycloakUserID:       stringPtrOrNil(domainEmployee.KeycloakUserID),
	}
}
