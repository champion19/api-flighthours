package domain

import (
	"strings"
	"testing"
)

func TestAircraftModel_ToLogger_Active(t *testing.T) {
	model := &AircraftModel{
		ID:               "model-123",
		ModelName:        "Boeing 737-800",
		AircraftTypeName: "Narrow Body",
		EngineTypeName:   "JET",
		Family:           "737",
		Manufacturer:     "Boeing",
		Status:           true,
	}

	result := model.ToLogger()

	if len(result) != 7 {
		t.Errorf("expected 7 log entries, got %d", len(result))
	}

	foundActive := false
	for _, entry := range result {
		if strings.Contains(entry, "status:active") {
			foundActive = true
		}
	}
	if !foundActive {
		t.Error("expected 'status:active' in log output")
	}

	foundModel := false
	for _, entry := range result {
		if strings.Contains(entry, "model_name:Boeing 737-800") {
			foundModel = true
		}
	}
	if !foundModel {
		t.Error("expected model name in log output")
	}
}

func TestAircraftModel_ToLogger_Inactive(t *testing.T) {
	model := &AircraftModel{
		ID:               "model-456",
		ModelName:        "Airbus A320",
		AircraftTypeName: "Narrow Body",
		EngineTypeName:   "JET",
		Family:           "A320",
		Manufacturer:     "Airbus",
		Status:           false,
	}

	result := model.ToLogger()

	foundInactive := false
	for _, entry := range result {
		if strings.Contains(entry, "status:inactive") {
			foundInactive = true
		}
	}
	if !foundInactive {
		t.Error("expected 'status:inactive' in log output")
	}
}
