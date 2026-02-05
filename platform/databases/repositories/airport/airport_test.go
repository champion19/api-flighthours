package airport

import (
	"testing"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

func TestAirport_ToDomain(t *testing.T) {
	t.Run("converts all fields correctly", func(t *testing.T) {
		city := "Bogota"
		country := "Colombia"
		iataCode := "BOG"
		airportType := "international"

		dbAirport := &Airport{
			ID:          "airport-uuid-123",
			Name:        "El Dorado International",
			City:        &city,
			Country:     &country,
			IATACode:    &iataCode,
			Status:      true,
			AirportType: &airportType,
		}

		domainAirport := dbAirport.ToDomain()

		if domainAirport.ID != dbAirport.ID {
			t.Errorf("expected ID %s, got %s", dbAirport.ID, domainAirport.ID)
		}
		if domainAirport.Name != dbAirport.Name {
			t.Errorf("expected Name %s, got %s", dbAirport.Name, domainAirport.Name)
		}
		if domainAirport.City != city {
			t.Errorf("expected City %s, got %s", city, domainAirport.City)
		}
		if domainAirport.Country != country {
			t.Errorf("expected Country %s, got %s", country, domainAirport.Country)
		}
		if domainAirport.IATACode != iataCode {
			t.Errorf("expected IATACode %s, got %s", iataCode, domainAirport.IATACode)
		}
		if domainAirport.Status != dbAirport.Status {
			t.Errorf("expected Status %v, got %v", dbAirport.Status, domainAirport.Status)
		}
		if domainAirport.AirportType != airportType {
			t.Errorf("expected AirportType %s, got %s", airportType, domainAirport.AirportType)
		}
	})

	t.Run("handles nil optional fields", func(t *testing.T) {
		dbAirport := &Airport{
			ID:     "airport-uuid-456",
			Name:   "Test Airport",
			Status: false,
		}

		domainAirport := dbAirport.ToDomain()

		if domainAirport.City != "" {
			t.Errorf("expected empty City, got %s", domainAirport.City)
		}
		if domainAirport.Country != "" {
			t.Errorf("expected empty Country, got %s", domainAirport.Country)
		}
		if domainAirport.IATACode != "" {
			t.Errorf("expected empty IATACode, got %s", domainAirport.IATACode)
		}
		if domainAirport.AirportType != "" {
			t.Errorf("expected empty AirportType, got %s", domainAirport.AirportType)
		}
	})
}

func TestFromDomain(t *testing.T) {
	t.Run("converts all fields correctly", func(t *testing.T) {
		domainAirport := &domain.Airport{
			ID:          "airport-uuid-123",
			Name:        "El Dorado International",
			City:        "Bogota",
			Country:     "Colombia",
			IATACode:    "BOG",
			Status:      true,
			AirportType: "international",
		}

		dbAirport := FromDomain(domainAirport)

		if dbAirport.ID != domainAirport.ID {
			t.Errorf("expected ID %s, got %s", domainAirport.ID, dbAirport.ID)
		}
		if dbAirport.Name != domainAirport.Name {
			t.Errorf("expected Name %s, got %s", domainAirport.Name, dbAirport.Name)
		}
		if dbAirport.City == nil || *dbAirport.City != domainAirport.City {
			t.Errorf("expected City %s, got %v", domainAirport.City, dbAirport.City)
		}
		if dbAirport.Country == nil || *dbAirport.Country != domainAirport.Country {
			t.Errorf("expected Country %s, got %v", domainAirport.Country, dbAirport.Country)
		}
		if dbAirport.IATACode == nil || *dbAirport.IATACode != domainAirport.IATACode {
			t.Errorf("expected IATACode %s, got %v", domainAirport.IATACode, dbAirport.IATACode)
		}
		if dbAirport.Status != domainAirport.Status {
			t.Errorf("expected Status %v, got %v", domainAirport.Status, dbAirport.Status)
		}
		if dbAirport.AirportType == nil || *dbAirport.AirportType != domainAirport.AirportType {
			t.Errorf("expected AirportType %s, got %v", domainAirport.AirportType, dbAirport.AirportType)
		}
	})

	t.Run("handles empty optional fields as nil", func(t *testing.T) {
		domainAirport := &domain.Airport{
			ID:     "airport-uuid-456",
			Name:   "Test Airport",
			Status: false,
		}

		dbAirport := FromDomain(domainAirport)

		if dbAirport.City != nil {
			t.Errorf("expected nil City, got %v", dbAirport.City)
		}
		if dbAirport.Country != nil {
			t.Errorf("expected nil Country, got %v", dbAirport.Country)
		}
		if dbAirport.IATACode != nil {
			t.Errorf("expected nil IATACode, got %v", dbAirport.IATACode)
		}
		if dbAirport.AirportType != nil {
			t.Errorf("expected nil AirportType, got %v", dbAirport.AirportType)
		}
	})
}
