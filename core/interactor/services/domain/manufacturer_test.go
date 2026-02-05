package domain

import (
	"testing"
)

func TestManufacturer_ToLogger(t *testing.T) {
	t.Run("returns expected log fields", func(t *testing.T) {
		m := &Manufacturer{
			ID:   "mfg-123",
			Name: "Boeing",
		}

		result := m.ToLogger()
		if len(result) != 2 {
			t.Fatalf("expected 2 fields, got %d", len(result))
		}
		if result[0] != "id:mfg-123" {
			t.Errorf("expected 'id:mfg-123', got %s", result[0])
		}
		if result[1] != "name:Boeing" {
			t.Errorf("expected 'name:Boeing', got %s", result[1])
		}
	})

	t.Run("handles empty fields", func(t *testing.T) {
		m := &Manufacturer{}
		result := m.ToLogger()
		if len(result) != 2 {
			t.Fatalf("expected 2 fields, got %d", len(result))
		}
		if result[0] != "id:" {
			t.Errorf("expected 'id:', got %s", result[0])
		}
		if result[1] != "name:" {
			t.Errorf("expected 'name:', got %s", result[1])
		}
	})
}
