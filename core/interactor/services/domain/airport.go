package domain

// Airport represents the airport domain model
type Airport struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	City        string `json:"city"`
	Country     string `json:"country"`
	IATACode    string `json:"iata_code"`
	Status      bool   `json:"status"`
	AirportType string `json:"airport_type"`
}

// ToLogger returns a slice of strings for logging airport information
func (a *Airport) ToLogger() []string {
	status := "inactive"
	if a.Status {
		status = "active"
	}
	return []string{
		"id:" + a.ID,
		"name:" + a.Name,
		"iata_code:" + a.IATACode,
		"city:" + a.City,
		"country:" + a.Country,
		"status:" + status,
	}
}

// IsActive returns true if the airport status is active
func (a *Airport) IsActive() bool {
	return a.Status
}
