package handlers

import "github.com/champion19/api-flighthours/core/interactor/services/domain"

// ==========================================
// EP1: Flight Hours Summary
// ==========================================

type FlightHoursSummaryResponse struct {
	Period        string            `json:"period"`
	StartDate     string            `json:"start_date"`
	EndDate       string            `json:"end_date"`
	TotalHours    string            `json:"total_hours"`
	TotalFlights  int               `json:"total_flights"`
	TotalLandings int               `json:"total_landings"`
	Breakdown     map[string]string `json:"breakdown"`
}

func FromDomainFlightHoursSummary(s *domain.FlightHoursSummary) FlightHoursSummaryResponse {
	breakdown := make(map[string]string)
	for role, minutes := range s.Breakdown {
		breakdown[role] = domain.FormatMinutesToHHMM(minutes)
	}

	return FlightHoursSummaryResponse{
		Period:        s.Period,
		StartDate:     s.StartDate,
		EndDate:       s.EndDate,
		TotalHours:    domain.FormatMinutesToHHMM(s.TotalMinutes),
		TotalFlights:  s.TotalFlights,
		TotalLandings: s.TotalLandings,
		Breakdown:     breakdown,
	}
}

// ==========================================
// EP2: Flight Alerts
// ==========================================

type FlightAlertResponse struct {
	Type         string `json:"type"`
	Severity     string `json:"severity"`
	Message      string `json:"message"`
	CurrentValue int    `json:"current_value"`
	Threshold    int    `json:"threshold"`
}

type FlightAlertsResponse struct {
	Alerts []FlightAlertResponse `json:"alerts"`
}

func FromDomainFlightAlerts(alerts []domain.FlightAlert) FlightAlertsResponse {
	responses := make([]FlightAlertResponse, 0, len(alerts))
	for _, a := range alerts {
		responses = append(responses, FlightAlertResponse{
			Type:         a.Type,
			Severity:     a.Severity,
			Message:      a.Message,
			CurrentValue: a.CurrentValue,
			Threshold:    a.Threshold,
		})
	}
	return FlightAlertsResponse{Alerts: responses}
}
