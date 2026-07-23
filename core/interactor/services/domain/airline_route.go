package domain

// Airline route lifecycle states. An airline_route starts as Pending when
// auto-created (an employee picked an origin/destination pair that has a
// route configured but wasn't yet linked to their airline) and must be
// promoted to Active by an admin before it can be used to create flights.
const (
	AirlineRouteStatusActive   = "active"
	AirlineRouteStatusPending  = "pending"
	AirlineRouteStatusInactive = "inactive"
)

// AirlineRoute represents the airline route domain model
// Links a physical route with a specific airline
type AirlineRoute struct {
	ID                     string `json:"id"`
	RouteID                string `json:"route_id"`
	AirlineID              string `json:"airline_id"`
	Status                 string `json:"status"`
	OriginAirportID        string `json:"origin_airport_id,omitempty"`
	DestinationAirportID   string `json:"destination_airport_id,omitempty"`
	OriginIataCode         string `json:"origin_iata_code,omitempty"`
	OriginOaciCode         string `json:"origin_oaci_code,omitempty"`
	DestinationIataCode    string `json:"destination_iata_code,omitempty"`
	DestinationOaciCode    string `json:"destination_oaci_code,omitempty"`
	RouteCode              string `json:"route_code,omitempty"`
	OriginAirportName      string `json:"origin_airport_name,omitempty"`
	DestinationAirportName string `json:"destination_airport_name,omitempty"`
	AirportType            string `json:"airport_type,omitempty"`
	EstimatedFlightTime    string `json:"estimated_flight_time,omitempty"`
	AirlineCode            string `json:"airline_code,omitempty"`
	AirlineName            string `json:"airline_name,omitempty"`
}

func (ar *AirlineRoute) IsActive() bool   { return ar.Status == AirlineRouteStatusActive }
func (ar *AirlineRoute) IsPending() bool  { return ar.Status == AirlineRouteStatusPending }
func (ar *AirlineRoute) IsInactive() bool { return ar.Status == AirlineRouteStatusInactive }

// ToLogger returns a slice of strings for logging airline route information
func (ar *AirlineRoute) ToLogger() []string {
	return []string{
		"id:" + ar.ID,
		"route_id:" + ar.RouteID,
		"airline_id:" + ar.AirlineID,
		"airline_code:" + ar.AirlineCode,
		"route_code:" + ar.RouteCode,
		"status:" + ar.Status,
	}
}
