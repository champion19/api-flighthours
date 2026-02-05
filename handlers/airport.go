package handlers

import (
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

type AirportResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	IATACode    string `json:"iata_code,omitempty"`
	Status      string `json:"status"`
	AirportType string `json:"airport_type,omitempty"`
	Links       []Link `json:"_links,omitempty"`
}

func FromDomainAirport(airport *domain.Airport, encodedID string) AirportResponse {
	status := "inactive"
	if airport.Status {
		status = "active"
	}
	return AirportResponse{
		ID:          encodedID,
		Name:        airport.Name,
		IATACode:    airport.IATACode,
		Status:      status,
		AirportType: airport.AirportType,
	}
}

type AirportStatusResponse struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Updated bool   `json:"updated"`
	Links   []Link `json:"_links,omitempty"`
}

type AirportListResponse struct {
	Airports []AirportResponse `json:"airports"`
	Total    int               `json:"total"`
	Links    []Link            `json:"_links,omitempty"`
}

func ToAirportListResponse(airports []domain.Airport, encodeFunc func(string) (string, error), baseURL string) AirportListResponse {
	response := AirportListResponse{
		Airports: make([]AirportResponse, 0, len(airports)),
		Total:    len(airports),
	}

	for _, airport := range airports {
		encodedID, err := encodeFunc(airport.ID)
		if err != nil {
			encodedID = airport.ID
		}
		status := "inactive"
		if airport.Status {
			status = "active"
		}
		airportResp := AirportResponse{
			ID:          encodedID,
			Name:        airport.Name,
			IATACode:    airport.IATACode,
			Status:      status,
			AirportType: airport.AirportType,
		}
		if baseURL != "" {
			airportResp.Links = BuildAirportLinks(baseURL, encodedID)
		}
		response.Airports = append(response.Airports, airportResp)
	}

	if baseURL != "" {
		response.Links = BuildAirportListLinks(baseURL)
	}

	return response
}
