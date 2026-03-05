package handlers

import (
	"testing"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

func TestFromDomainFlightHoursSummary(t *testing.T) {
	t.Run("converts with breakdown", func(t *testing.T) {
		summary := &domain.FlightHoursSummary{
			Period:        "monthly",
			StartDate:     "2026-01-01",
			EndDate:       "2026-01-31",
			TotalMinutes:  150,
			TotalFlights:  5,
			TotalLandings: 3,
			Breakdown: map[string]int{
				"PF": 90,
				"PM": 60,
			},
		}
		result := FromDomainFlightHoursSummary(summary)

		if result.Period != "monthly" {
			t.Errorf("expected period 'monthly', got %q", result.Period)
		}
		if result.TotalHours != "2:30" {
			t.Errorf("expected '2:30', got %q", result.TotalHours)
		}
		if result.TotalFlights != 5 {
			t.Errorf("expected 5, got %d", result.TotalFlights)
		}
		if result.TotalLandings != 3 {
			t.Errorf("expected 3, got %d", result.TotalLandings)
		}
		if result.Breakdown["PF"] != "1:30" {
			t.Errorf("expected PF='1:30', got %q", result.Breakdown["PF"])
		}
		if result.Breakdown["PM"] != "1:00" {
			t.Errorf("expected PM='1:00', got %q", result.Breakdown["PM"])
		}
	})

	t.Run("empty breakdown", func(t *testing.T) {
		summary := &domain.FlightHoursSummary{
			Breakdown: map[string]int{},
		}
		result := FromDomainFlightHoursSummary(summary)
		if len(result.Breakdown) != 0 {
			t.Errorf("expected empty breakdown, got %d", len(result.Breakdown))
		}
	})
}

func TestFromDomainFlightAlerts(t *testing.T) {
	t.Run("converts alerts", func(t *testing.T) {
		alerts := []domain.FlightAlert{
			{
				Type:         domain.AlertTypeHourLimitMonthly,
				Severity:     domain.AlertSeverityWarning,
				Message:      "test message",
				CurrentValue: 5000,
				Threshold:    5400,
			},
			{
				Type:         domain.AlertTypeMinLandings90D,
				Severity:     domain.AlertSeverityInfo,
				Message:      "landing alert",
				CurrentValue: 2,
				Threshold:    3,
			},
		}
		result := FromDomainFlightAlerts(alerts)
		if len(result.Alerts) != 2 {
			t.Fatalf("expected 2 alerts, got %d", len(result.Alerts))
		}
		if result.Alerts[0].Type != domain.AlertTypeHourLimitMonthly {
			t.Errorf("expected type %s, got %s", domain.AlertTypeHourLimitMonthly, result.Alerts[0].Type)
		}
		if result.Alerts[1].CurrentValue != 2 {
			t.Errorf("expected current value 2, got %d", result.Alerts[1].CurrentValue)
		}
	})

	t.Run("empty alerts", func(t *testing.T) {
		result := FromDomainFlightAlerts(nil)
		if len(result.Alerts) != 0 {
			t.Errorf("expected 0 alerts, got %d", len(result.Alerts))
		}
	})
}
