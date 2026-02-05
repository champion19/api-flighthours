package domain

import (
	"testing"
)

func TestAirport_ToLogger(t *testing.T) {
	t.Run("returns expected log fields for active airport", func(t *testing.T) {
		airport := &Airport{
			ID:          "airport-uuid-123",
			Name:        "El Dorado International",
			IATACode:    "BOG",
			Status:      true,
			AirportType: "international",
		}

		result := airport.ToLogger()

		if len(result) != 5 {
			t.Fatalf("expected 5 log fields, got %d", len(result))
		}

		expected := []string{
			"id:airport-uuid-123",
			"name:El Dorado International",
			"iata_code:BOG",
			"airport_type:international",
			"status:active",
		}

		for i, exp := range expected {
			if result[i] != exp {
				t.Errorf("expected result[%d] = %q, got %q", i, exp, result[i])
			}
		}
	})

	t.Run("returns inactive status for inactive airport", func(t *testing.T) {
		airport := &Airport{
			ID:     "airport-uuid-456",
			Name:   "Test Airport",
			Status: false,
		}

		result := airport.ToLogger()

		found := false
		for _, field := range result {
			if field == "status:inactive" {
				found = true
				break
			}
		}
		if !found {
			t.Error("expected status:inactive in log fields")
		}
	})

	t.Run("handles empty fields", func(t *testing.T) {
		airport := &Airport{}

		result := airport.ToLogger()

		if len(result) != 5 {
			t.Fatalf("expected 5 log fields, got %d", len(result))
		}

		expected := []string{
			"id:",
			"name:",
			"iata_code:",
			"airport_type:",
			"status:inactive",
		}

		for i, exp := range expected {
			if result[i] != exp {
				t.Errorf("expected result[%d] = %q, got %q", i, exp, result[i])
			}
		}
	})
}

func TestAirport_IsActive(t *testing.T) {
	tests := []struct {
		name     string
		status   bool
		expected bool
	}{
		{
			name:     "true status returns true",
			status:   true,
			expected: true,
		},
		{
			name:     "false status returns false",
			status:   false,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			airport := &Airport{Status: tt.status}
			result := airport.IsActive()
			if result != tt.expected {
				t.Errorf("expected IsActive() = %v, got %v", tt.expected, result)
			}
		})
	}
}
