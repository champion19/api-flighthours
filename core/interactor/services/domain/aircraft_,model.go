package domain


type AircraftModel struct {
	ID               string `json:"id"`
	ModelName        string `json:"model_name"`
	AircraftTypeName string `json:"aircraft_type_name"`
	EngineTypeName   string `json:"engine_type_name"`
	Family           string `json:"family"`
	Manufacturer     string `json:"manufacturer"`
	Status           bool   `json:"status"`
}


func (am *AircraftModel) ToLogger() []string {
	statusStr := "inactive"
	if am.Status {
		statusStr = "active"
	}
	return []string{
		"id:" + am.ID,
		"model_name:" + am.ModelName,
		"aircraft_type_name:" + am.AircraftTypeName,
		"engine_type_name:" + am.EngineTypeName,
		"family:" + am.Family,
		"manufacturer:" + am.Manufacturer,
		"status:" + statusStr,
	}
}
