package airline

import (
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

// Airline represents the database entity for airline table
type Airline struct {
	ID          string `db:"id"`
	AirlineName string `db:"airline_name"`
	AirlineCode string `db:"airline_code"`
	Status      string `db:"status"`
}

// ToDomain converts the database entity to domain model
func (a *Airline) ToDomain() *domain.Airline {
	return &domain.Airline{
		ID:          a.ID,
		AirlineName: a.AirlineName,
		AirlineCode: a.AirlineCode,
		Status:      a.Status,
	}
}

// FromDomain converts a domain model to database entity
func FromDomain(domainAirline *domain.Airline) *Airline {
	return &Airline{
		ID:          domainAirline.ID,
		AirlineName: domainAirline.AirlineName,
		AirlineCode: domainAirline.AirlineCode,
		Status:      domainAirline.Status,
	}
}
