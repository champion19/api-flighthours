package airport

import domain "github.com/champion19/api-flighthours/core/interactor/services/domain"

// Airport is the database entity for airport table
type Airport struct {
	ID          string  `db:"id"`
	Name        string  `db:"name"`
	City        *string `db:"city"`
	Country     *string `db:"country"`
	IATACode    *string `db:"iata_code"`
	Status      bool    `db:"status"`
	AirportType *string `db:"airport_type"`
}

// ToDomain converts the database entity to domain model
func (a *Airport) ToDomain() *domain.Airport {
	airport := &domain.Airport{
		ID:     a.ID,
		Name:   a.Name,
		Status: a.Status,
	}

	if a.City != nil {
		airport.City = *a.City
	}
	if a.Country != nil {
		airport.Country = *a.Country
	}
	if a.IATACode != nil {
		airport.IATACode = *a.IATACode
	}
	if a.AirportType != nil {
		airport.AirportType = *a.AirportType
	}

	return airport
}

// FromDomain converts a domain model to database entity
func FromDomain(domainAirport *domain.Airport) *Airport {
	airport := &Airport{
		ID:     domainAirport.ID,
		Name:   domainAirport.Name,
		Status: domainAirport.Status,
	}

	if domainAirport.City != "" {
		airport.City = &domainAirport.City
	}
	if domainAirport.Country != "" {
		airport.Country = &domainAirport.Country
	}
	if domainAirport.IATACode != "" {
		airport.IATACode = &domainAirport.IATACode
	}
	if domainAirport.AirportType != "" {
		airport.AirportType = &domainAirport.AirportType
	}

	return airport
}
