package domain

// AircraftFamily represents the aircraft family domain model
// It groups aircraft models by family and manufacturer
type AircraftFamily struct {
	Family       string `json:"family"`
	Manufacturer string `json:"manufacturer"`
}

func (af *AircraftFamily) ToLogger() []string {
	return []string{
		"family:" + af.Family,
		"manufacturer:" + af.Manufacturer,
	}
}
