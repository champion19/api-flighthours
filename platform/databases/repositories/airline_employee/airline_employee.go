package airline_employee

import (
	"time"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
)

// AirlineEmployee repository struct - airline-specific fields only
// This maps to the airline-related columns in the employee table
type AirlineEmployee struct {
	ID        string    `db:"id"`
	AirlineID *string   `db:"airline"`
	Bp        *string   `db:"bp"`
	StartDate time.Time `db:"start_date"`
	EndDate   time.Time `db:"end_date"`
	Active    bool      `db:"active"`
}

func (e *AirlineEmployee) ToDomain() *domain.AirlineEmployee {
	airlineID := ""
	if e.AirlineID != nil {
		airlineID = *e.AirlineID
	}
	bp := ""
	if e.Bp != nil {
		bp = *e.Bp
	}

	return &domain.AirlineEmployee{
		ID:        e.ID,
		AirlineID: airlineID,
		Bp:        bp,
		StartDate: e.StartDate,
		EndDate:   e.EndDate,
		Active:    e.Active,
	}
}

func stringPtrOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func FromDomain(domainEmployee *domain.AirlineEmployee) *AirlineEmployee {
	return &AirlineEmployee{
		ID:        domainEmployee.ID,
		AirlineID: stringPtrOrNil(domainEmployee.AirlineID),
		Bp:        stringPtrOrNil(domainEmployee.Bp),
		StartDate: domainEmployee.StartDate,
		EndDate:   domainEmployee.EndDate,
		Active:    domainEmployee.Active,
	}
}
