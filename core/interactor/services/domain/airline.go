package domain

// Airline represents an airline entity in the domain layer
type Airline struct {
	ID          string `json:"id"`
	AirlineName string `json:"airline_name"`
	AirlineCode string `json:"airline_code"`
	Status      string `json:"status"`
}

// ToLogger returns a slice of strings with relevant information for logging
func (a *Airline) ToLogger() []string {
	return []string{
		"id:" + a.ID,
		"name:" + a.AirlineName,
		"code:" + a.AirlineCode,
		"status:" + a.Status,
	}
}
