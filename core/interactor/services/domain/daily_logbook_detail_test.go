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

func TestIsValidApproachCategory(t *testing.T) {
	tests := []struct {
		name     string
		category string
		want     bool
	}{
		{"RNP", "RNP", true},
		{"ILS", "ILS", true},
		{"VISUAL", "VISUAL", true},
		{"empty is valid", "", true},
		{"invalid", "APV", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidApproachCategory(tt.category); got != tt.want {
				t.Errorf("IsValidApproachCategory(%q) = %v, want %v", tt.category, got, tt.want)
			}
		})
	}
}

func TestIsValidApproachSubtype(t *testing.T) {
	tests := []struct {
		name     string
		category ApproachCategory
		subtype  string
		want     bool
	}{
		{"RNP LNAV", ApproachCategoryRNP, "LNAV", true},
		{"RNP LNAV/VNAV", ApproachCategoryRNP, "LNAV/VNAV", true},
		{"RNP AR = 0.3", ApproachCategoryRNP, "RNP AR = 0.3", true},
		{"RNP AR < 0.3", ApproachCategoryRNP, "RNP AR < 0.3", true},
		{"RNP with ILS subtype", ApproachCategoryRNP, "CAT I", false},
		{"ILS CAT I", ApproachCategoryILS, "CAT I", true},
		{"ILS CAT II", ApproachCategoryILS, "CAT II", true},
		{"ILS CAT III > 175", ApproachCategoryILS, "CAT III > 175", true},
		{"ILS CAT III < 175", ApproachCategoryILS, "CAT III < 175", true},
		{"ILS with RNP subtype", ApproachCategoryILS, "LNAV", false},
		{"VISUAL empty subtype", ApproachCategoryVisual, "", true},
		{"VISUAL with subtype", ApproachCategoryVisual, "LNAV", false},
		{"no category, empty subtype", "", "", true},
		{"no category, non-empty subtype", "", "LNAV", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidApproachSubtype(tt.category, tt.subtype); got != tt.want {
				t.Errorf("IsValidApproachSubtype(%q, %q) = %v, want %v", tt.category, tt.subtype, got, tt.want)
			}
		})
	}
}

func TestValidateApproachFields(t *testing.T) {
	rnp := ApproachCategoryRNP
	ils := ApproachCategoryILS
	visual := ApproachCategoryVisual
	lnav := "LNAV"
	catI := "CAT I"
	invalidSubtype := "VOR"
	trueVal := true
	falseVal := false

	tests := []struct {
		name     string
		category *ApproachCategory
		subtype  *string
		autoland *bool
		wantErr  bool
	}{
		{"nil everything", nil, nil, nil, false},
		{"valid RNP + LNAV", &rnp, &lnav, nil, false},
		{"valid ILS + CAT I, no autoland", &ils, &catI, nil, false},
		{"valid ILS + CAT I + autoland true", &ils, &catI, &trueVal, false},
		{"valid VISUAL, no subtype", &visual, nil, nil, false},
		{"RNP subtype mismatched with ILS category", &ils, &lnav, nil, true},
		{"unknown subtype", &rnp, &invalidSubtype, nil, true},
		{"autoland true but category RNP", &rnp, &lnav, &trueVal, true},
		{"autoland true but category VISUAL", &visual, nil, &trueVal, true},
		{"autoland false is fine outside ILS", &rnp, &lnav, &falseVal, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateApproachFields(tt.category, tt.subtype, tt.autoland)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateApproachFields() error = %v, wantErr %v", err, tt.wantErr)
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
