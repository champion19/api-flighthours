package airline_route

import domain "github.com/champion19/api-flighthours/core/interactor/services/domain"

type AirlineRoute struct {
	ID                     string `db:"id"`
	RouteID                string `db:"route_id"`
	AirlineID              string `db:"airline_id"`
	Status                 string `db:"status"`
	AirlineCode            string `db:"airline_code"`
	AirlineName            string `db:"airline_name"`
	OriginAirportID        string `db:"origin_airport_id"`
	OriginIataCode         string `db:"origin_iata_code"`
	OriginOaciCode         string `db:"origin_oaci_code"`
	DestinationAirportID   string `db:"destination_airport_id"`
	DestinationIataCode    string `db:"destination_iata_code"`
	DestinationOaciCode    string `db:"destination_oaci_code"`
	RouteCode              string `db:"route_code"`
	OriginAirportName      string `db:"origin_airport_name"`
	DestinationAirportName string `db:"destination_airport_name"`
	AirportType            string `db:"airport_type"`
	EstimatedFlightTime    string `db:"estimated_flight_time"`
}

func (ar *AirlineRoute) ToDomain() *domain.AirlineRoute {
	return &domain.AirlineRoute{
		ID:                     ar.ID,
		RouteID:                ar.RouteID,
		AirlineID:              ar.AirlineID,
		Status:                 ar.Status,
		AirlineCode:            ar.AirlineCode,
		AirlineName:            ar.AirlineName,
		OriginAirportID:        ar.OriginAirportID,
		OriginIataCode:         ar.OriginIataCode,
		OriginOaciCode:         ar.OriginOaciCode,
		DestinationAirportID:   ar.DestinationAirportID,
		DestinationIataCode:    ar.DestinationIataCode,
		DestinationOaciCode:    ar.DestinationOaciCode,
		RouteCode:              ar.RouteCode,
		OriginAirportName:      ar.OriginAirportName,
		DestinationAirportName: ar.DestinationAirportName,
		AirportType:            ar.AirportType,
		EstimatedFlightTime:    ar.EstimatedFlightTime,
	}
}

func FromDomain(domainAirlineRoute *domain.AirlineRoute) *AirlineRoute {
	return &AirlineRoute{
		ID:                     domainAirlineRoute.ID,
		RouteID:                domainAirlineRoute.RouteID,
		AirlineID:              domainAirlineRoute.AirlineID,
		Status:                 domainAirlineRoute.Status,
		AirlineCode:            domainAirlineRoute.AirlineCode,
		AirlineName:            domainAirlineRoute.AirlineName,
		OriginAirportID:        domainAirlineRoute.OriginAirportID,
		OriginIataCode:         domainAirlineRoute.OriginIataCode,
		OriginOaciCode:         domainAirlineRoute.OriginOaciCode,
		DestinationAirportID:   domainAirlineRoute.DestinationAirportID,
		DestinationIataCode:    domainAirlineRoute.DestinationIataCode,
		DestinationOaciCode:    domainAirlineRoute.DestinationOaciCode,
		RouteCode:              domainAirlineRoute.RouteCode,
		OriginAirportName:      domainAirlineRoute.OriginAirportName,
		DestinationAirportName: domainAirlineRoute.DestinationAirportName,
		AirportType:            domainAirlineRoute.AirportType,
		EstimatedFlightTime:    domainAirlineRoute.EstimatedFlightTime,
	}
}
