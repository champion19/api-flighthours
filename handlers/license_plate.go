package handlers

import(
domain  "github.com/champion19/api-flighthours/core/interactor/services/domain"
"github.com/google/uuid"
)


type LicensePlateResponse struct {
	ID              string `json:"id"`
	LicensePlate    string `json:"license_plate"` // Numero de Matrícula
	ModelName       string `json:"model_name"`    // Modelo (denormalized)
	AirlineName     string `json:"airline_name"`  // Aerolínea (denormalized)
	AircraftModelID string `json:"aircraft_model_id"`
	AirlineID       string `json:"airline_id"`
	Links           []Link `json:"_links,omitempty"`
}


func FromDomainLicensePlate(reg *domain.LicensePlate, encodedID, encodedModelID, encodedAirlineID string) LicensePlateResponse {
	return LicensePlateResponse{
		ID:              encodedID,
		LicensePlate:    reg.LicensePlate,
		ModelName:       reg.ModelName,
		AirlineName:     reg.AirlineName,
		AircraftModelID: encodedModelID,
		AirlineID:       encodedAirlineID,
	}
}


type CreateLicensePlateRequest struct {
	LicensePlate    string `json:"license_plate" binding:"required"`
	AircraftModelID string `json:"aircraft_model_id" binding:"required"`
	AirlineID       string `json:"airline_id" binding:"required"`
}


func (r *CreateLicensePlateRequest) Sanitize() {
	r.LicensePlate = TrimString(r.LicensePlate)
	r.AircraftModelID = TrimString(r.AircraftModelID)
	r.AirlineID = TrimString(r.AirlineID)
}


func (r *CreateLicensePlateRequest) ToDomain() domain.LicensePlate {
	return domain.LicensePlate{
		ID:              uuid.New().String(),
		LicensePlate:    r.LicensePlate,
		AircraftModelID: r.AircraftModelID,
		AirlineID:       r.AirlineID,
	}
}


type UpdateLicensePlateRequest struct {
	LicensePlate    string `json:"license_plate" binding:"required"`
	AircraftModelID string `json:"aircraft_model_id" binding:"required"`
	AirlineID       string `json:"airline_id" binding:"required"`
}


func (r *UpdateLicensePlateRequest) Sanitize() {
	r.LicensePlate = TrimString(r.LicensePlate)
	r.AircraftModelID = TrimString(r.AircraftModelID)
	r.AirlineID = TrimString(r.AirlineID)
}


func (r *UpdateLicensePlateRequest) ToDomain(id string) domain.LicensePlate {
	return domain.LicensePlate{
		ID:              id,
		LicensePlate:    r.LicensePlate,
		AircraftModelID: r.AircraftModelID,
		AirlineID:       r.AirlineID,
	}
}


type LicensePlateListResponse struct {
	Registrations []LicensePlateResponse `json:"registrations"`
	Total         int                            `json:"total"`
	Links         []Link                         `json:"_links,omitempty"`
}



func ToLicensePlateListResponse(registrations []domain.LicensePlate, encodeFunc func(string) (string, error), baseURL string) LicensePlateListResponse {
	response := LicensePlateListResponse{
		Registrations: make([]LicensePlateResponse, 0, len(registrations)),
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
		regResp := LicensePlateResponse{
			ID:              encodedID,
			LicensePlate:    reg.LicensePlate,
			ModelName:       reg.ModelName,
			AirlineName:     reg.AirlineName,
			AircraftModelID: encodedModelID,
			AirlineID:       encodedAirlineID,
		}

		if baseURL != "" {
			regResp.Links = BuildLicensePlateLinks(baseURL, encodedID)
		}
		response.Registrations = append(response.Registrations, regResp)
	}


	if baseURL != "" {
		response.Links = BuildLicensePlateListLinks(baseURL)
	}

	return response
}
