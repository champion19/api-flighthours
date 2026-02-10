package domain


type Route struct {
	ID                   string `json:"id"`
	OriginAirportID      string `json:"origin_airport_id"`
	DestinationAirportID string `json:"destination_airport_id"`
	AirportType          string `json:"airport_type"`
	EstimatedFlightTime  string `json:"estimated_flight_time"` // Format: "HH:MM:SS"
	OriginIataCode         string `json:"origin_iata_code"`
	OriginAirportName      string `json:"origin_airport_name"`
	DestinationIataCode    string `json:"destination_iata_code"`
	DestinationAirportName string `json:"destination_airport_name"`
	RouteCode              string `json:"route_code"` // e.g., "BOG-CLO"
}


func (r *Route) ToLogger() []string {
	return []string{
		"id:" + r.ID,
		"origin_airport_id:" + r.OriginAirportID,
		"destination_airport_id:" + r.DestinationAirportID,
		"origin_iata_code:" + r.OriginIataCode,
		"destination_iata_code:" + r.DestinationIataCode,
		"airport_type:" + r.AirportType,
		"route_code:" + r.RouteCode,
	}
}
