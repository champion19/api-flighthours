package domain

// AirlineRoute represents the airline route domain model
// Links a physical route with a specific airline
type AirlineRoute struct {
	ID        string `json:"id"`
	RouteID   string `json:"route_id"`
	AirlineID string `json:"airline_id"`
	Status    bool   `json:"status"`
	OriginIataCode         string `json:"origin_iata_code,omitempty"`
	DestinationIataCode    string `json:"destination_iata_code,omitempty"`
	RouteCode              string `json:"route_code,omitempty"`
	OriginAirportName      string `json:"origin_airport_name,omitempty"`
	DestinationAirportName string `json:"destination_airport_name,omitempty"`
	AirportType            string `json:"airport_type,omitempty"`
	EstimatedFlightTime    string `json:"estimated_flight_time,omitempty"`
	AirlineCode string `json:"airline_code,omitempty"`
	AirlineName string `json:"airline_name,omitempty"`
}

// ToLogger returns a slice of strings for logging airline route information
func (ar *AirlineRoute) ToLogger() []string {
	statusStr := "inactive"
	if ar.Status {
		statusStr = "active"
	}
	return []string{
		"id:" + ar.ID,
		"route_id:" + ar.RouteID,
		"airline_id:" + ar.AirlineID,
		"airline_code:" + ar.AirlineCode,
		"route_code:" + ar.RouteCode,
		"status:" + statusStr,
	}
}
