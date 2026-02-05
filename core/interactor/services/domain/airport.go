package domain

type Airport struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	IATACode    string `json:"iata_code"`
	Status      bool   `json:"status"`
	AirportType string `json:"airport_type"`
}

func (a *Airport) ToLogger() []string {
	status := "inactive"
	if a.Status {
		status = "active"
	}
	return []string{
		"id:" + a.ID,
		"name:" + a.Name,
		"iata_code:" + a.IATACode,
		"airport_type:" + a.AirportType,
		"status:" + status,
	}
}

func (a *Airport) IsActive() bool {
	return a.Status
}
