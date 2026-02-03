package domain

import (
	"testing"
	"time"
)

func TestAirlineEmployee_SetID(t *testing.T) {
	e := &AirlineEmployee{}
	e.SetID()

	if e.ID == "" {
		t.Error("SetID should generate a non-empty ID")
	}

	// Check it's a valid UUID format (36 chars with hyphens)
	if len(e.ID) != 36 {
		t.Errorf("expected UUID length 36, got %d", len(e.ID))
	}
}

func TestAirlineEmployee_ToLogger(t *testing.T) {
	now := time.Now()
	e := &AirlineEmployee{
		ID:        "emp-123",
		AirlineID: "airline-456",
		Bp:        "BP12345",
		StartDate: now,
		EndDate:   now.AddDate(1, 0, 0),
		Active:    true,
	}

	logs := e.ToLogger()

	if len(logs) != 6 {
		t.Errorf("expected 6 log items, got %d", len(logs))
	}

	// Verify first element contains ID
	if logs[0] != "id:emp-123" {
		t.Errorf("expected 'id:emp-123', got %q", logs[0])
	}

	// Verify airline_id
	if logs[1] != "airline_id:airline-456" {
		t.Errorf("expected 'airline_id:airline-456', got %q", logs[1])
	}

	// Verify bp
	if logs[2] != "bp:BP12345" {
		t.Errorf("expected 'bp:BP12345', got %q", logs[2])
	}

	// Verify active status
	if logs[5] != "active:true" {
		t.Errorf("expected 'active:true', got %q", logs[5])
	}
}

func TestAirlineEmployee_ToLogger_Inactive(t *testing.T) {
	e := &AirlineEmployee{
		ID:     "emp-inactive",
		Active: false,
	}

	logs := e.ToLogger()
	if logs[5] != "active:false" {
		t.Errorf("expected 'active:false', got %q", logs[5])
	}
}
