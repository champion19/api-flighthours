package domain

import (
	"testing"
	"time"
)

func TestDailyLogbook_SetID(t *testing.T) {
	d := &DailyLogbook{}
	d.SetID()
	if d.ID == "" {
		t.Error("expected non-empty ID after SetID()")
	}
}

func TestDailyLogbook_ToLogger(t *testing.T) {
	d := &DailyLogbook{
		ID:         "test-id",
		EmployeeID: "emp-123",
		LogDate:    time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
	}
	result := d.ToLogger()
	if len(result) != 3 {
		t.Errorf("expected 3 items, got %d", len(result))
	}
	if result[0] != "id:test-id" {
		t.Errorf("expected 'id:test-id', got %q", result[0])
	}
}

func TestDailyLogbook_IsActive(t *testing.T) {
	tests := []struct {
		name   string
		status bool
		want   bool
	}{
		{"active", true, true},
		{"inactive", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DailyLogbook{Status: tt.status}
			if got := d.IsActive(); got != tt.want {
				t.Errorf("IsActive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDailyLogbook_StatusString(t *testing.T) {
	tests := []struct {
		name   string
		status bool
		want   string
	}{
		{"active", true, "active"},
		{"inactive", false, "inactive"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DailyLogbook{Status: tt.status}
			if got := d.StatusString(); got != tt.want {
				t.Errorf("StatusString() = %q, want %q", got, tt.want)
			}
		})
	}
}
