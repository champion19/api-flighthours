package domain

import "testing"

func TestDailyLogbookDetail_SetID(t *testing.T) {
	d := &DailyLogbookDetail{}
	d.SetID()
	if d.ID == "" {
		t.Error("expected non-empty ID after SetID()")
	}
}

func pilotRolePtr(r PilotRole) *PilotRole { return &r }

func TestDailyLogbookDetail_ToLogger(t *testing.T) {
	t.Run("with pilot role only", func(t *testing.T) {
		d := &DailyLogbookDetail{
			ID:             "detail-1",
			DailyLogbookID: "lb-1",
			FlightNumber:   "AV123",
			FlightRealDate: "2025-01-15",
			RouteCode:      "BOG-CLO",
			TailNumber:   "HK-5432",
			PilotRole:      pilotRolePtr(PilotRolePF),
		}
		result := d.ToLogger()
		if len(result) != 8 {
			t.Errorf("expected 8 items, got %d", len(result))
		}
	})

	t.Run("with crew role set", func(t *testing.T) {
		crewRole := CrewRole("captain")
		d := &DailyLogbookDetail{
			ID:             "detail-2",
			DailyLogbookID: "lb-2",
			FlightNumber:   "AV456",
			FlightRealDate: "2025-01-16",
			RouteCode:      "CLO-BOG",
			TailNumber:   "HK-1234",
			PilotRole:      pilotRolePtr(PilotRolePM),
			CrewRole:       &crewRole,
		}
		result := d.ToLogger()
		if len(result) != 8 {
			t.Errorf("expected 8 items, got %d", len(result))
		}
	})

	t.Run("with nil optional fields", func(t *testing.T) {
		d := &DailyLogbookDetail{
			ID:             "detail-3",
			DailyLogbookID: "lb-3",
			FlightNumber:   "AV789",
		}
		result := d.ToLogger()
		if len(result) != 8 {
			t.Errorf("expected 8 items, got %d", len(result))
		}
	})
}

func TestIsValidPilotRole(t *testing.T) {
	tests := []struct {
		role string
		want bool
	}{
		{"PF", true},
		{"PM", true},
		{"PFTO", true},
		{"PFL", true},
		{"INVALID", false},
		{"", false},
	}
	for _, tt := range tests {
		t.Run(tt.role, func(t *testing.T) {
			if got := IsValidPilotRole(tt.role); got != tt.want {
				t.Errorf("IsValidPilotRole(%q) = %v, want %v", tt.role, got, tt.want)
			}
		})
	}
}

func TestIsValidApproachType(t *testing.T) {
	tests := []struct {
		name         string
		approachType string
		want         bool
	}{
		{"APV", "APV", true},
		{"VISUAL", "VISUAL", true},
		{"empty is valid", "", true},
		{"invalid", "ILS", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidApproachType(tt.approachType); got != tt.want {
				t.Errorf("IsValidApproachType(%q) = %v, want %v", tt.approachType, got, tt.want)
			}
		})
	}
}

func TestIsValidCrewRole(t *testing.T) {
	tests := []struct {
		name string
		role string
		want bool
	}{
		{"captain", "captain", true},
		{"first officer", "first officer", true},
		{"empty is valid", "", true},
		{"invalid", "navigator", false},
		{"case sensitive", "Captain", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidCrewRole(tt.role); got != tt.want {
				t.Errorf("IsValidCrewRole(%q) = %v, want %v", tt.role, got, tt.want)
			}
		})
	}
}
