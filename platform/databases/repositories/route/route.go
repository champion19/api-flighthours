package route

import domain "github.com/champion19/api-flighthours/core/interactor/services/domain"

type Route struct {
	ID                     string `db:"id"`
	OriginAirportID        string `db:"origin_airport_id"`
	OriginIataCode         string `db:"origin_iata_code"`
	OriginAirportName      string `db:"origin_airport_name"`
	DestinationAirportID   string `db:"destination_airport_id"`
	DestinationIataCode    string `db:"destination_iata_code"`
	DestinationAirportName string `db:"destination_airport_name"`
	AirportType            string `db:"airport_type"`
	EstimatedFlightTime    string `db:"estimated_flight_time"`
	RouteCode              string `db:"route_code"`
}

func (r *Route) ToDomain() *domain.Route {
	return &domain.Route{
		ID:                     r.ID,
		OriginAirportID:        r.OriginAirportID,
		OriginIataCode:         r.OriginIataCode,
		OriginAirportName:      r.OriginAirportName,
		DestinationAirportID:   r.DestinationAirportID,
		DestinationIataCode:    r.DestinationIataCode,
		DestinationAirportName: r.DestinationAirportName,
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
		OriginAirportName:      domainRoute.OriginAirportName,
		DestinationAirportID:   domainRoute.DestinationAirportID,
		DestinationIataCode:    domainRoute.DestinationIataCode,
		DestinationAirportName: domainRoute.DestinationAirportName,
		AirportType:            domainRoute.AirportType,
		EstimatedFlightTime:    domainRoute.EstimatedFlightTime,
		RouteCode:              domainRoute.RouteCode,
	}
}
