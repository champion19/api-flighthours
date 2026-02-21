package aircraftmodel

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

func TestAircraftModel_ToDomain(t *testing.T) {
	am := &AircraftModel{
		ID: "am1", ModelName: "Boeing 737-800", AircraftTypeName: "Narrow Body",
		EngineTypeName: "Turbofan", Family: "737", Manufacturer: "Boeing", Status: true,
	}
	result := am.ToDomain()
	if result.ID != "am1" {
		t.Errorf("expected am1, got %s", result.ID)
	}
	if result.ModelName != "Boeing 737-800" {
		t.Errorf("expected Boeing 737-800, got %s", result.ModelName)
	}
	if result.Status != true {
		t.Error("expected Status true")
	}
}

func TestAircraftModel_FromDomain(t *testing.T) {
	dm := &domain.AircraftModel{
		ID: "am2", ModelName: "A320", AircraftTypeName: "Narrow Body",
		EngineTypeName: "Turbofan", Family: "A320", Manufacturer: "Airbus", Status: true,
	}
	result := FromDomain(dm)
	if result.ID != "am2" {
		t.Errorf("expected am2, got %s", result.ID)
	}
	if result.Manufacturer != "Airbus" {
		t.Errorf("expected Airbus, got %s", result.Manufacturer)
	}
}

func TestAircraftModel_RoundTrip(t *testing.T) {
	original := &domain.AircraftModel{
		ID: "am3", ModelName: "E190", AircraftTypeName: "Regional",
		EngineTypeName: "Turbofan", Family: "E-Jet", Manufacturer: "Embraer", Status: true,
	}
	restored := FromDomain(original).ToDomain()
	if restored.ID != original.ID {
		t.Errorf("ID mismatch")
	}
	if restored.ModelName != original.ModelName {
		t.Errorf("ModelName mismatch")
	}
	if restored.Manufacturer != original.Manufacturer {
		t.Errorf("Manufacturer mismatch")
	}
}

func TestAircraftModel_ToFamilyDomain(t *testing.T) {
	am := &AircraftModel{Family: "737", Manufacturer: "Boeing"}
	result := am.ToFamilyDomain()
	if result.Family != "737" {
		t.Errorf("expected 737, got %s", result.Family)
	}
	if result.Manufacturer != "Boeing" {
		t.Errorf("expected Boeing, got %s", result.Manufacturer)
	}
}

// mockScanner implements the scanner interface for testing scanAircraftModel
type mockScanner struct {
	values []interface{}
	err    error
}

func (m *mockScanner) Scan(dest ...interface{}) error {
	if m.err != nil {
		return m.err
	}
	for i, d := range dest {
		switch v := d.(type) {
		case *string:
			*v = m.values[i].(string)
		case *bool:
			*v = m.values[i].(bool)
		case *sql.NullString:
			ns := m.values[i].(sql.NullString)
			*v = ns
		}
	}
	return nil
}

func TestScanAircraftModel_Success(t *testing.T) {
	s := &mockScanner{
		values: []interface{}{
			"am1", "Boeing 737", "Narrow",
			sql.NullString{String: "Turbofan", Valid: true},
			"737",
			sql.NullString{String: "Boeing", Valid: true},
			true,
		},
	}

	result, err := scanAircraftModel(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.EngineTypeName != "Turbofan" {
		t.Errorf("expected Turbofan, got %s", result.EngineTypeName)
	}
	if result.Manufacturer != "Boeing" {
		t.Errorf("expected Boeing, got %s", result.Manufacturer)
	}
}

func TestScanAircraftModel_NullFields(t *testing.T) {
	s := &mockScanner{
		values: []interface{}{
			"am1", "ATR 72", "Turboprop",
			sql.NullString{Valid: false},
			"ATR",
			sql.NullString{Valid: false},
			true,
		},
	}

	result, err := scanAircraftModel(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.EngineTypeName != "" {
		t.Errorf("expected empty EngineTypeName, got %s", result.EngineTypeName)
	}
	if result.Manufacturer != "" {
		t.Errorf("expected empty Manufacturer, got %s", result.Manufacturer)
	}
}

func TestScanAircraftModel_Error(t *testing.T) {
	s := &mockScanner{err: errors.New("scan error")}

	_, err := scanAircraftModel(s)
	if err == nil {
		t.Fatal("expected error")
	}
}
