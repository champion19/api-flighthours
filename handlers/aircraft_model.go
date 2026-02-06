package handlers

import domain "github.com/champion19/api-flighthours/core/interactor/services/domain"

// AircraftModelResponse - Response DTO for aircraft model data
type AircraftModelResponse struct {
	ID               string `json:"id"`
	ModelName        string `json:"model_name"`
	AircraftTypeName string `json:"aircraft_type_name"`
	EngineTypeName   string `json:"engine_type_name,omitempty"`
	Family           string `json:"family"`
	Manufacturer     string `json:"manufacturer,omitempty"`
	Status           bool   `json:"status"`
	Links            []Link `json:"_links,omitempty"`
}

// FromDomainAircraftModel converts domain.AircraftModel to AircraftModelResponse with encoded ID
func FromDomainAircraftModel(model *domain.AircraftModel, encodedID string) AircraftModelResponse {
	return AircraftModelResponse{
		ID:               encodedID,
		ModelName:        model.ModelName,
		AircraftTypeName: model.AircraftTypeName,
		EngineTypeName:   model.EngineTypeName,
		Family:           model.Family,
		Manufacturer:     model.Manufacturer,
		Status:           model.Status,
	}
}

// AircraftModelStatusResponse - Response DTO for aircraft model status change
type AircraftModelStatusResponse struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Updated bool   `json:"updated"`
	Links   []Link `json:"_links,omitempty"`
}

// AircraftModelListResponse - Response DTO for listing aircraft models
type AircraftModelListResponse struct {
	AircraftModels []AircraftModelResponse `json:"aircraft_models"`
	Total          int                     `json:"total"`
	Links          []Link                  `json:"_links,omitempty"`
}

// ToAircraftModelListResponse converts a slice of domain.AircraftModel to AircraftModelListResponse
// baseURL is used to build HATEOAS links for each model
func ToAircraftModelListResponse(models []domain.AircraftModel, encodeFunc func(string) (string, error), baseURL string) AircraftModelListResponse {
	response := AircraftModelListResponse{
		AircraftModels: make([]AircraftModelResponse, 0, len(models)),
		Total:          len(models),
	}

	for _, model := range models {
		encodedID, err := encodeFunc(model.ID)
		if err != nil {
			// If encoding fails, use the original UUID
			encodedID = model.ID
		}
		modelResp := AircraftModelResponse{
			ID:               encodedID,
			ModelName:        model.ModelName,
			AircraftTypeName: model.AircraftTypeName,
			EngineTypeName:   model.EngineTypeName,
			Family:           model.Family,
			Manufacturer:     model.Manufacturer,
			Status:           model.Status,
		}
		// Add HATEOAS links to each model
		if baseURL != "" {
			modelResp.Links = BuildAircraftModelLinks(baseURL, encodedID)
		}
		response.AircraftModels = append(response.AircraftModels, modelResp)
	}

	// Add collection-level links
	if baseURL != "" {
		response.Links = BuildAircraftModelListLinks(baseURL)
	}

	return response
}
