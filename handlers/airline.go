package handlers

import (
	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
)

// AirlineResponse is the DTO for airline data in API responses
type AirlineResponse struct {
	ID          string `json:"id"`
	AirlineName string `json:"airline_name"`
	AirlineCode string `json:"airline_code"`
	Status      string `json:"status"`
	Links       []Link `json:"_links,omitempty"`
}

// FromDomainAirline converts a domain Airline to an AirlineResponse DTO
func FromDomainAirline(airline *domain.Airline, encodedID string) AirlineResponse {
	return AirlineResponse{
		ID:          encodedID,
		AirlineName: airline.AirlineName,
		AirlineCode: airline.AirlineCode,
		Status:      airline.Status,
	}
}

// AirlineListResponse - Response DTO for listing airlines
type AirlineListResponse struct {
	Airlines []AirlineResponse `json:"airlines"`
	Total    int               `json:"total"`
	Links    []Link            `json:"_links,omitempty"`
}

// ToAirlineListResponse converts a slice of domain.Airline to AirlineListResponse
// baseURL is used to build HATEOAS links for each airline
func ToAirlineListResponse(airlines []domain.Airline, encodeFunc func(string) (string, error), baseURL string) AirlineListResponse {
	response := AirlineListResponse{
		Airlines: make([]AirlineResponse, 0, len(airlines)),
		Total:    len(airlines),
	}

	for _, airline := range airlines {
		encodedID, err := encodeFunc(airline.ID)
		if err != nil {
			// If encoding fails, use the original UUID
			encodedID = airline.ID
		}
		airlineResp := AirlineResponse{
			ID:          encodedID,
			AirlineName: airline.AirlineName,
			AirlineCode: airline.AirlineCode,
			Status:      airline.Status,
		}
		// Add HATEOAS links to each airline
		if baseURL != "" {
			airlineResp.Links = BuildAirlineLinks(baseURL, encodedID)
		}
		response.Airlines = append(response.Airlines, airlineResp)
	}

	// Add collection-level links
	if baseURL != "" {
		response.Links = BuildAirlineListLinks(baseURL)
	}

	return response
}
