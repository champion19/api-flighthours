package domain

// TailNumber represents the aircraft registration domain model
// It links a tail number (matrícula) with an aircraft model and an airline
type TailNumber struct {
	ID              string `json:"id"`
	TailNumber    string `json:"tail_number"`
	AircraftModelID string `json:"aircraft_model_id"`
	AirlineID       string `json:"airline_id"`
	ModelName   string `json:"model_name,omitempty"`
	AirlineName string `json:"airline_name,omitempty"`
}


func (ar *TailNumber) ToLogger() []string {
	return []string{
		"id:" + ar.ID,
		"tail_number:" + ar.TailNumber,
		"aircraft_model_id:" + ar.AircraftModelID,
		"airline_id:" + ar.AirlineID,
	}
}
