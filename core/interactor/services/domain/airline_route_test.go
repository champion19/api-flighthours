package domain

import "testing"

func TestAirlineRoute_ToLogger_Active(t *testing.T) {
	ar := &AirlineRoute{
		ID:          "route-123",
		RouteID:     "route-base-1",
		AirlineID:   "airline-456",
		AirlineCode: "AV",
		RouteCode:   "BOG-CLO",
		Status:      true,
	}

	logs := ar.ToLogger()

	if len(logs) != 6 {
		t.Errorf("expected 6 log items, got %d", len(logs))
	}

	if logs[0] != "id:route-123" {
		t.Errorf("expected 'id:route-123', got %q", logs[0])
	}

	if logs[1] != "route_id:route-base-1" {
		t.Errorf("expected 'route_id:route-base-1', got %q", logs[1])
	}

	if logs[2] != "airline_id:airline-456" {
		t.Errorf("expected 'airline_id:airline-456', got %q", logs[2])
	}

	if logs[3] != "airline_code:AV" {
		t.Errorf("expected 'airline_code:AV', got %q", logs[3])
	}

	if logs[4] != "route_code:BOG-CLO" {
		t.Errorf("expected 'route_code:BOG-CLO', got %q", logs[4])
	}

	if logs[5] != "status:active" {
		t.Errorf("expected 'status:active', got %q", logs[5])
	}
}

func TestAirlineRoute_ToLogger_Inactive(t *testing.T) {
	ar := &AirlineRoute{
		ID:     "route-inactive",
		Status: false,
	}

	logs := ar.ToLogger()

	if logs[5] != "status:inactive" {
		t.Errorf("expected 'status:inactive', got %q", logs[5])
	}
}

func TestRoute_ToLogger(t *testing.T) {
	r := &Route{
		ID:                   "route-uuid",
		OriginAirportID:      "airport-1",
		DestinationAirportID: "airport-2",
		OriginIataCode:       "BOG",
		DestinationIataCode:  "CLO",
		AirportType:          "International",
		RouteCode:            "BOG-CLO",
	}

	logs := r.ToLogger()

	if len(logs) != 7 {
		t.Errorf("expected 7 log items, got %d", len(logs))
	}

	if logs[0] != "id:route-uuid" {
		t.Errorf("expected 'id:route-uuid', got %q", logs[0])
	}

	if logs[3] != "origin_iata_code:BOG" {
		t.Errorf("expected 'origin_iata_code:BOG', got %q", logs[3])
	}

	if logs[4] != "destination_iata_code:CLO" {
		t.Errorf("expected 'destination_iata_code:CLO', got %q", logs[4])
	}

	if logs[6] != "route_code:BOG-CLO" {
		t.Errorf("expected 'route_code:BOG-CLO', got %q", logs[6])
	}
}

func TestEngine_ToLogger(t *testing.T) {
	e := &Engine{
		ID:   "engine-123",
		Name: "Turbofan",
	}

	logs := e.ToLogger()

	if len(logs) != 2 {
		t.Errorf("expected 2 log items, got %d", len(logs))
	}

	if logs[0] != "id:engine-123" {
		t.Errorf("expected 'id:engine-123', got %q", logs[0])
	}

	if logs[1] != "name:Turbofan" {
		t.Errorf("expected 'name:Turbofan', got %q", logs[1])
	}
}
