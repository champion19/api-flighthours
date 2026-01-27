package domain

import (
	"testing"
)

func TestAirline_ToLogger(t *testing.T) {
	t.Run("returns expected log fields", func(t *testing.T) {
		airline := &Airline{
			ID:          "airline-uuid-123",
			AirlineName: "Test Airlines",
			AirlineCode: "TST",
			Status:      "active",
		}

		result := airline.ToLogger()

		if len(result) != 4 {
			t.Fatalf("expected 4 log fields, got %d", len(result))
		}

		expected := []string{
			"id:airline-uuid-123",
			"name:Test Airlines",
			"code:TST",
			"status:active",
		}

		for i, exp := range expected {
			if result[i] != exp {
				t.Errorf("expected result[%d] = %q, got %q", i, exp, result[i])
			}
		}
	})

	t.Run("handles empty fields", func(t *testing.T) {
		airline := &Airline{}

		result := airline.ToLogger()

		if len(result) != 4 {
			t.Fatalf("expected 4 log fields, got %d", len(result))
		}

		expected := []string{
			"id:",
			"name:",
			"code:",
			"status:",
		}

		for i, exp := range expected {
			if result[i] != exp {
				t.Errorf("expected result[%d] = %q, got %q", i, exp, result[i])
			}
		}
	})
}
