package handlers

import domain "github.com/champion19/api-flighthours/core/interactor/services/domain"

// EngineResponse - Response DTO for engine data
type EngineResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Links []Link `json:"_links,omitempty"`
}

// FromDomainEngine converts domain.Engine to EngineResponse with encoded ID
func FromDomainEngine(engine *domain.Engine, encodedID string) EngineResponse {
	return EngineResponse{
		ID:   encodedID,
		Name: engine.Name,
	}
}

// EngineListResponse - Response DTO for listing engines
type EngineListResponse struct {
	Engines []EngineResponse `json:"engines"`
	Total   int              `json:"total"`
	Links   []Link           `json:"_links,omitempty"`
}

// ToEngineListResponse converts a slice of domain.Engine to EngineListResponse
// baseURL is used to build HATEOAS links for each engine
func ToEngineListResponse(engines []domain.Engine, encodeFunc func(string) (string, error), baseURL string) EngineListResponse {
	response := EngineListResponse{
		Engines: make([]EngineResponse, 0, len(engines)),
		Total:   len(engines),
	}

	for _, engine := range engines {
		encodedID, err := encodeFunc(engine.ID)
		if err != nil {
			// If encoding fails, use the original UUID
			encodedID = engine.ID
		}
		engineResp := EngineResponse{
			ID:   encodedID,
			Name: engine.Name,
		}
		// Add HATEOAS links to each engine
		if baseURL != "" {
			engineResp.Links = BuildEngineLinks(baseURL, encodedID)
		}
		response.Engines = append(response.Engines, engineResp)
	}

	// Add collection-level links
	if baseURL != "" {
		response.Links = BuildEngineListLinks(baseURL)
	}

	return response
}
