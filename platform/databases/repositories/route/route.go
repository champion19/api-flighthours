package route

import domain "github.com/champion19/api-flighthours/core/interactor/services/domain"

type Route struct {
	ID                     string `db:"id"`
	OriginAirportID        string `db:"origin_airport_id"`
	OriginIataCode         string `db:"origin_iata_code"`
	OriginOaciCode         string `db:"origin_oaci_code"`
	OriginAirportName      string `db:"origin_airport_name"`
	OriginCountry          string `db:"origin_country"`
	DestinationAirportID   string `db:"destination_airport_id"`
	DestinationIataCode    string `db:"destination_iata_code"`
	DestinationOaciCode    string `db:"destination_oaci_code"`
	DestinationAirportName string `db:"destination_airport_name"`
	DestinationCountry     string `db:"destination_country"`
	AirportType            string `db:"airport_type"`
	EstimatedFlightTime    string `db:"estimated_flight_time"`
	RouteCode              string `db:"route_code"`
}

func (r *Route) ToDomain() *domain.Route {
	return &domain.Route{
		ID:                     r.ID,
		OriginAirportID:        r.OriginAirportID,
		OriginIataCode:         r.OriginIataCode,
		OriginOaciCode:         r.OriginOaciCode,
		OriginAirportName:      r.OriginAirportName,
		OriginCountry:          r.OriginCountry,
		DestinationAirportID:   r.DestinationAirportID,
		DestinationIataCode:    r.DestinationIataCode,
		DestinationOaciCode:    r.DestinationOaciCode,
		DestinationAirportName: r.DestinationAirportName,
		DestinationCountry:     r.DestinationCountry,
		AirportType:            r.AirportType,
		EstimatedFlightTime:    r.EstimatedFlightTime,
		RouteCode:              r.RouteCode,
	}
}

func FromDomain(domainRoute *domain.Route) *Route {
	return &Route{
		ID:                     domainRoute.ID,
		OriginAirportID:        domainRoute.OriginAirportID,
		OriginIataCode:         domainRoute.OriginIataCode,
		OriginOaciCode:         domainRoute.OriginOaciCode,
		OriginAirportName:      domainRoute.OriginAirportName,
		OriginCountry:          domainRoute.OriginCountry,
		DestinationAirportID:   domainRoute.DestinationAirportID,
		DestinationIataCode:    domainRoute.DestinationIataCode,
		DestinationOaciCode:    domainRoute.DestinationOaciCode,
		DestinationAirportName: domainRoute.DestinationAirportName,
		DestinationCountry:     domainRoute.DestinationCountry,
		AirportType:            domainRoute.AirportType,
		EstimatedFlightTime:    domainRoute.EstimatedFlightTime,
		RouteCode:              domainRoute.RouteCode,
	}
}
