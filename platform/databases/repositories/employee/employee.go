package employee

import (
"time"

"github.com/champion19/flighthours-api/core/interactor/services/domain"
)

type Employee struct {
	ID                   string    `db:"id"`
	Name                 string    `db:"name"`
	Airline              *string    `db:"airline"`
	Email                string    `db:"email"`
	IdentificationNumber string    `db:"identification_number"`
	Bp                   *string    `db:"bp"`
	StartDate            time.Time `db:"start_date"`
	EndDate              time.Time `db:"end_date"`
	Active               bool      `db:"active"`
	Role                 string    `db:"role"`
	KeycloakUserID       *string    `db:"keycloak_user_id"`
}

func (e Employee) ToDomain() domain.Employee {
	airline := ""
	if e.Airline != nil {
		airline = *e.Airline
	}
	bp := ""
	if e.Bp != nil {
		bp = *e.Bp
	}
	keycloakUserID := ""
	if e.KeycloakUserID != nil {
		keycloakUserID = *e.KeycloakUserID
	}

	return domain.Employee{
		ID:                   e.ID,
		Name:                 e.Name,
		Email:                e.Email,
		Airline:              airline,
		IdentificationNumber: e.IdentificationNumber,
		Bp:                   bp,
		StartDate:            e.StartDate,
		EndDate:              e.EndDate,
		Active:               e.Active,
		Role:                 e.Role,
		KeycloakUserID:       keycloakUserID,
	}
}

// stringPtrOrNil returns nil for empty strings, otherwise a pointer to the string
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
		Email:                domainEmployee.Email,
		Airline:              stringPtrOrNil(domainEmployee.Airline),
		IdentificationNumber: domainEmployee.IdentificationNumber,
		Bp:                   stringPtrOrNil(domainEmployee.Bp),
		StartDate:            domainEmployee.StartDate,
		EndDate:              domainEmployee.EndDate,
		Active:               domainEmployee.Active,
		Role:                 domainEmployee.Role,
		KeycloakUserID:       stringPtrOrNil(domainEmployee.KeycloakUserID),
	}
}
