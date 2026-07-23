package handlers

import domain "github.com/champion19/api-flighthours/core/interactor/services/domain"

type AirlineRouteResponse struct {
	ID                     string `json:"id"`
	RouteID                string `json:"route_id"`
	AirlineID              string `json:"airline_id"`
	Status                 string `json:"status"`
	AirlineCode            string `json:"airline_code"`
	AirlineName            string `json:"airline_name"`
	OriginAirportID        string `json:"origin_airport_id,omitempty"`
	OriginIataCode         string `json:"origin_iata_code"`
	OriginOaciCode         string `json:"origin_oaci_code,omitempty"`
	DestinationAirportID   string `json:"destination_airport_id,omitempty"`
	DestinationIataCode    string `json:"destination_iata_code"`
	DestinationOaciCode    string `json:"destination_oaci_code,omitempty"`
	RouteCode              string `json:"route_code"`
	OriginAirportName      string `json:"origin_airport_name,omitempty"`
	DestinationAirportName string `json:"destination_airport_name,omitempty"`
	AirportType            string `json:"airport_type,omitempty"`
	EstimatedFlightTime    string `json:"estimated_flight_time,omitempty"`
	Links                  []Link `json:"_links,omitempty"`
}

func FromDomainAirlineRoute(ar *domain.AirlineRoute, encodedID, encodedRouteID, encodedAirlineID string) AirlineRouteResponse {
	return AirlineRouteResponse{
		ID:                     encodedID,
		RouteID:                encodedRouteID,
		AirlineID:              encodedAirlineID,
		Status:                 ar.Status,
		AirlineCode:            ar.AirlineCode,
		AirlineName:            ar.AirlineName,
		OriginIataCode:         ar.OriginIataCode,
		DestinationIataCode:    ar.DestinationIataCode,
		RouteCode:              ar.RouteCode,
		OriginAirportName:      ar.OriginAirportName,
		DestinationAirportName: ar.DestinationAirportName,
		AirportType:            ar.AirportType,
		EstimatedFlightTime:    ar.EstimatedFlightTime,
	}
}

type AirlineRouteListResponse struct {
	AirlineRoutes []AirlineRouteResponse `json:"airline_routes"`
	Total         int                    `json:"total"`
	Links         []Link                 `json:"_links,omitempty"`
}

func ToAirlineRouteListResponse(airlineRoutes []domain.AirlineRoute, encodeFunc func(string) (string, error), baseURL string) AirlineRouteListResponse {
	response := AirlineRouteListResponse{
		AirlineRoutes: make([]AirlineRouteResponse, 0, len(airlineRoutes)),
		Total:         len(airlineRoutes),
	}

	for _, ar := range airlineRoutes {
		encodedID, err := encodeFunc(ar.ID)
		if err != nil {
			encodedID = ar.ID
		}
		encodedRouteID, err := encodeFunc(ar.RouteID)
		if err != nil {
			encodedRouteID = ar.RouteID
		}
		encodedAirlineID, err := encodeFunc(ar.AirlineID)
		if err != nil {
			encodedAirlineID = ar.AirlineID
		}
		// Encoded the same way as AirportResponse.ID so the client can match
		// a selected airport (by IATA or OACI) back to this route by ID,
		// instead of by code — codes can be missing for OACI-only airports.
		encodedOriginAirportID, err := encodeFunc(ar.OriginAirportID)
		if err != nil {
			encodedOriginAirportID = ar.OriginAirportID
		}
		encodedDestinationAirportID, err := encodeFunc(ar.DestinationAirportID)
		if err != nil {
			encodedDestinationAirportID = ar.DestinationAirportID
		}
		arResp := AirlineRouteResponse{
			ID:                     encodedID,
			RouteID:                encodedRouteID,
			AirlineID:              encodedAirlineID,
			Status:                 ar.Status,
			AirlineCode:            ar.AirlineCode,
			AirlineName:            ar.AirlineName,
			OriginAirportID:        encodedOriginAirportID,
			OriginIataCode:         ar.OriginIataCode,
			OriginOaciCode:         ar.OriginOaciCode,
			DestinationAirportID:   encodedDestinationAirportID,
			DestinationIataCode:    ar.DestinationIataCode,
			DestinationOaciCode:    ar.DestinationOaciCode,
			RouteCode:              ar.RouteCode,
			OriginAirportName:      ar.OriginAirportName,
			DestinationAirportName: ar.DestinationAirportName,
			AirportType:            ar.AirportType,
			EstimatedFlightTime:    ar.EstimatedFlightTime,
		}

		if baseURL != "" {
			arResp.Links = BuildAirlineRouteLinks(baseURL, encodedID, ar.Status == domain.AirlineRouteStatusActive)
		}
		response.AirlineRoutes = append(response.AirlineRoutes, arResp)
	}

	if baseURL != "" {
		response.Links = BuildAirlineRouteListLinks(baseURL)
	}

	return response
}
