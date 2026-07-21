package domain

type Airport struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	City        string `json:"city"`
	Country     string `json:"country"`
	IATACode    string `json:"iata_code"`
	OACICode    string `json:"oaci_code"`
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
		"city:" + a.City,
		"country:" + a.Country,
		"iata_code:" + a.IATACode,
		"oaci_code:" + a.OACICode,
		"airport_type:" + a.AirportType,
		"status:" + status,
	}
}

func (a *Airport) IsActive() bool {
	return a.Status
}
