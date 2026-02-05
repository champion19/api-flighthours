package airport

import domain "github.com/champion19/api-flighthours/core/interactor/services/domain"

// Airport is the DB model - matches actual database columns
type Airport struct {
	ID          string  `db:"id"`
	Name        string  `db:"name"`
	IATACode    *string `db:"iata_code"`
	Status      bool    `db:"status"`
	AirportType *string `db:"airport_type"`
}

func (a *Airport) ToDomain() *domain.Airport {
	airport := &domain.Airport{
		ID:     a.ID,
		Name:   a.Name,
		Status: a.Status,
	}

	if a.IATACode != nil {
		airport.IATACode = *a.IATACode
	}
	if a.AirportType != nil {
		airport.AirportType = *a.AirportType
	}

	return airport
}

func FromDomain(domainAirport *domain.Airport) *Airport {
	airport := &Airport{
		ID:     domainAirport.ID,
		Name:   domainAirport.Name,
		Status: domainAirport.Status,
	}

	if domainAirport.IATACode != "" {
		airport.IATACode = &domainAirport.IATACode
	}
	if domainAirport.AirportType != "" {
		airport.AirportType = &domainAirport.AirportType
	}

	return airport
}
