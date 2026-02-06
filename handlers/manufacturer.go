package handlers

import("github.com/champion19/api-flighthours/core/interactor/services/domain")

// ManufacturerResponse - Response DTO for manufacturer data
type ManufacturerResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Links []Link `json:"_links,omitempty"`
}

// FromDomainManufacturer converts domain.Manufacturer to ManufacturerResponse with encoded ID
func FromDomainManufacturer(manufacturer *domain.Manufacturer, encodedID string) ManufacturerResponse {
	return ManufacturerResponse{
		ID:   encodedID,
		Name: manufacturer.Name,
	}
}

// ManufacturerListResponse - Response DTO for listing manufacturers
type ManufacturerListResponse struct {
	Manufacturers []ManufacturerResponse `json:"manufacturers"`
	Total         int                    `json:"total"`
	Links         []Link                 `json:"_links,omitempty"`
}

// ToManufacturerListResponse converts a slice of domain.Manufacturer to ManufacturerListResponse
// baseURL is used to build HATEOAS links for each manufacturer
func ToManufacturerListResponse(manufacturers []domain.Manufacturer, encodeFunc func(string) (string, error), baseURL string) ManufacturerListResponse {
	response := ManufacturerListResponse{
		Manufacturers: make([]ManufacturerResponse, 0, len(manufacturers)),
		Total:         len(manufacturers),
	}

	for _, manufacturer := range manufacturers {
		encodedID, err := encodeFunc(manufacturer.ID)
		if err != nil {
			// If encoding fails, use the original UUID
			encodedID = manufacturer.ID
		}
		manufacturerResp := ManufacturerResponse{
			ID:   encodedID,
			Name: manufacturer.Name,
		}
		// Add HATEOAS links to each manufacturer
		if baseURL != "" {
			manufacturerResp.Links = BuildManufacturerLinks(baseURL, encodedID)
		}
		response.Manufacturers = append(response.Manufacturers, manufacturerResp)
	}

	// Add collection-level links
	if baseURL != "" {
		response.Links = BuildManufacturerListLinks(baseURL)
	}

	return response
}
