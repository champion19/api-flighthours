package handlers

import(
domain  "github.com/champion19/api-flighthours/core/interactor/services/domain"
"github.com/google/uuid"
)


type TailNumberResponse struct {
	ID              string `json:"id"`
	TailNumber    string `json:"tail_number"` // Numero de Matrícula
	ModelName       string `json:"model_name"`    // Modelo (denormalized)
	AirlineName     string `json:"airline_name"`  // Aerolínea (denormalized)
	AircraftModelID string `json:"aircraft_model_id"`
	AirlineID       string `json:"airline_id"`
	Links           []Link `json:"_links,omitempty"`
}


func FromDomainTailNumber(reg *domain.TailNumber, encodedID, encodedModelID, encodedAirlineID string) TailNumberResponse {
	return TailNumberResponse{
		ID:              encodedID,
		TailNumber:    reg.TailNumber,
		ModelName:       reg.ModelName,
		AirlineName:     reg.AirlineName,
		AircraftModelID: encodedModelID,
		AirlineID:       encodedAirlineID,
	}
}


type CreateTailNumberRequest struct {
	TailNumber    string `json:"tail_number" binding:"required"`
	AircraftModelID string `json:"aircraft_model_id" binding:"required"`
	AirlineID       string `json:"airline_id" binding:"required"`
}


func (r *CreateTailNumberRequest) Sanitize() {
	r.TailNumber = TrimString(r.TailNumber)
	r.AircraftModelID = TrimString(r.AircraftModelID)
	r.AirlineID = TrimString(r.AirlineID)
}


func (r *CreateTailNumberRequest) ToDomain() domain.TailNumber {
	return domain.TailNumber{
		ID:              uuid.New().String(),
		TailNumber:    r.TailNumber,
		AircraftModelID: r.AircraftModelID,
		AirlineID:       r.AirlineID,
	}
}


type UpdateTailNumberRequest struct {
	TailNumber    string `json:"tail_number" binding:"required"`
	AircraftModelID string `json:"aircraft_model_id" binding:"required"`
	AirlineID       string `json:"airline_id" binding:"required"`
}


func (r *UpdateTailNumberRequest) Sanitize() {
	r.TailNumber = TrimString(r.TailNumber)
	r.AircraftModelID = TrimString(r.AircraftModelID)
	r.AirlineID = TrimString(r.AirlineID)
}


func (r *UpdateTailNumberRequest) ToDomain(id string) domain.TailNumber {
	return domain.TailNumber{
		ID:              id,
		TailNumber:    r.TailNumber,
		AircraftModelID: r.AircraftModelID,
		AirlineID:       r.AirlineID,
	}
}


type TailNumberListResponse struct {
	Registrations []TailNumberResponse `json:"registrations"`
	Total         int                            `json:"total"`
	Links         []Link                         `json:"_links,omitempty"`
}



func ToTailNumberListResponse(registrations []domain.TailNumber, encodeFunc func(string) (string, error), baseURL string) TailNumberListResponse {
	response := TailNumberListResponse{
		Registrations: make([]TailNumberResponse, 0, len(registrations)),
		Total:         len(registrations),
	}

	for _, reg := range registrations {
		encodedID, err := encodeFunc(reg.ID)
		if err != nil {
			encodedID = reg.ID
		}
		encodedModelID, err := encodeFunc(reg.AircraftModelID)
		if err != nil {
			encodedModelID = reg.AircraftModelID
		}
		encodedAirlineID, err := encodeFunc(reg.AirlineID)
		if err != nil {
			encodedAirlineID = reg.AirlineID
		}
		regResp := TailNumberResponse{
			ID:              encodedID,
			TailNumber:    reg.TailNumber,
			ModelName:       reg.ModelName,
			AirlineName:     reg.AirlineName,
			AircraftModelID: encodedModelID,
			AirlineID:       encodedAirlineID,
		}

		if baseURL != "" {
			regResp.Links = BuildTailNumberLinks(baseURL, encodedID)
		}
		response.Registrations = append(response.Registrations, regResp)
	}


	if baseURL != "" {
		response.Links = BuildTailNumberListLinks(baseURL)
	}

	return response
}
