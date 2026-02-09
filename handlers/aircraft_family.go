package handlers

import domain "github.com/champion19/api-flighthours/core/interactor/services/domain"

// AircraftFamilyResponse - Response DTO for aircraft family data
type AircraftFamilyResponse struct {
	Family       string `json:"family"`
	Manufacturer string `json:"manufacturer"`
}

// AircraftFamilyListResponse - Response DTO for listing aircraft families
type AircraftFamilyListResponse struct {
	AircraftFamilies []AircraftFamilyResponse `json:"aircraft_families"`
	Total            int                      `json:"total"`
	Links            []Link                   `json:"_links,omitempty"`
}

// ToAircraftFamilyListResponse converts a slice of domain.AircraftModel to AircraftFamilyListResponse
func ToAircraftFamilyListResponse(models []domain.AircraftModel, baseURL string) AircraftFamilyListResponse {
	// Deduplicate by family name
	seen := make(map[string]bool)
	families := make([]AircraftFamilyResponse, 0)

	for _, model := range models {
		if !seen[model.Family] {
			seen[model.Family] = true
			families = append(families, AircraftFamilyResponse{
				Family:       model.Family,
				Manufacturer: model.Manufacturer,
			})
		}
	}

	response := AircraftFamilyListResponse{
		AircraftFamilies: families,
		Total:            len(families),
	}

	if baseURL != "" {
		response.Links = BuildAircraftFamilyListLinks(baseURL)
	}

	return response
}
