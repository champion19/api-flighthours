package handlers

import (
	"testing"
	"time"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
)

func strPtrH(s string) *string                                 { return &s }
func domainPilotRolePtrH(r domain.PilotRole) *domain.PilotRole { return &r }

func TestFromDomainDailyLogbook(t *testing.T) {
	t.Run("active logbook", func(t *testing.T) {
		page := int64(42)
		logbook := &domain.DailyLogbook{
			ID:         "uuid-1",
			LogDate:    time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			EmployeeID: "emp-1",
			BookPage:   &page,
			Status:     true,
		}
		resp := FromDomainDailyLogbook(logbook, "encoded-id", "encoded-emp")
		if resp.Status != "active" {
			t.Errorf("expected status 'active', got %s", resp.Status)
		}
		if resp.ID != "encoded-id" {
			t.Errorf("expected ID 'encoded-id', got %s", resp.ID)
		}
		if resp.EmployeeID != "encoded-emp" {
			t.Errorf("expected EmployeeID 'encoded-emp', got %s", resp.EmployeeID)
		}
		if resp.LogDate != "2025-01-15" {
			t.Errorf("expected LogDate '2025-01-15', got %s", resp.LogDate)
		}
		if *resp.BookPage != 42 {
			t.Errorf("expected BookPage 42, got %d", *resp.BookPage)
		}
	})

	t.Run("inactive logbook", func(t *testing.T) {
		logbook := &domain.DailyLogbook{Status: false}
		resp := FromDomainDailyLogbook(logbook, "id", "emp")
		if resp.Status != "inactive" {
			t.Errorf("expected status 'inactive', got %s", resp.Status)
		}
	})
}

func TestCreateDailyLogbookRequest_Sanitize(t *testing.T) {
	req := &CreateDailyLogbookRequest{
		LogDate: "  2025-01-15  ",
	}
	req.Sanitize()
	if req.LogDate != "2025-01-15" {
		t.Errorf("expected trimmed '2025-01-15', got '%s'", req.LogDate)
	}
}

func TestCreateDailyLogbookRequest_ToDomain(t *testing.T) {
	t.Run("valid date", func(t *testing.T) {
		page := int64(42)
		req := &CreateDailyLogbookRequest{
			LogDate:  "2025-01-15",
			BookPage: &page,
		}
		logbook, err := req.ToDomain("emp-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if logbook.EmployeeID != "emp-1" {
			t.Errorf("expected EmployeeID 'emp-1', got %s", logbook.EmployeeID)
		}
		if *logbook.BookPage != 42 {
			t.Errorf("expected BookPage 42, got %d", *logbook.BookPage)
		}
		if !logbook.Status {
			t.Error("expected Status true for new logbook")
		}
		if logbook.ID == "" {
			t.Error("expected ID to be set")
		}
	})

	t.Run("invalid date", func(t *testing.T) {
		req := &CreateDailyLogbookRequest{LogDate: "not-a-date"}
		_, err := req.ToDomain("emp-1")
		if err == nil {
			t.Fatal("expected error for invalid date")
		}
	})
}

func TestToDailyLogbookListResponse(t *testing.T) {
	logbooks := []domain.DailyLogbook{
		{ID: "550e8400-e29b-41d4-a716-446655440000", EmployeeID: "550e8400-e29b-41d4-a716-446655440001", LogDate: time.Now(), Status: true},
		{ID: "550e8400-e29b-41d4-a716-446655440002", EmployeeID: "550e8400-e29b-41d4-a716-446655440001", LogDate: time.Now(), Status: false},
	}

	t.Run("with baseURL", func(t *testing.T) {
		resp := ToDailyLogbookListResponse(logbooks, func(s string) (string, error) {
			return "encoded-" + s[:8], nil
		}, "http://localhost:8080/api/v1")

		if resp.Total != 2 {
			t.Errorf("expected Total 2, got %d", resp.Total)
		}
		if len(resp.Logbooks) != 2 {
			t.Errorf("expected 2 logbooks, got %d", len(resp.Logbooks))
		}
		if len(resp.Links) == 0 {
			t.Error("expected collection-level links to be set")
		}
	})

	t.Run("without baseURL", func(t *testing.T) {
		resp := ToDailyLogbookListResponse(logbooks, func(s string) (string, error) {
			return "encoded", nil
		}, "")

		if len(resp.Links) != 0 {
			t.Error("expected no collection-level links when baseURL is empty")
		}
	})

	t.Run("with encoding error", func(t *testing.T) {
		resp := ToDailyLogbookListResponse(logbooks, func(s string) (string, error) {
			return "", domain.ErrInvalidID
		}, "http://localhost")

		// Should fall back to original UUID when encoding fails
		if len(resp.Logbooks) != 2 {
			t.Errorf("expected 2 logbooks even with encoding error, got %d", len(resp.Logbooks))
		}
	})
}

func TestToDomainDailyLogbookDetail(t *testing.T) {
	passengers := 150
	companion := "Co-Pilot Smith"
	dutyTime := "10:00"
	approachType := "PA"
	flightType := "COMMERCIAL"

	req := CreateDailyLogbookDetailRequest{
		FlightRealDate: "2025-01-15",
		FlightNumber:   "AV123",
		AirlineRouteID: "route-1",
		LicensePlateID: "lp-1",
		Passengers:     &passengers,
		OutTime:        strPtrH("08:00"),
		TakeoffTime:    strPtrH("08:15"),
		LandingTime:    strPtrH("09:30"),
		InTime:         strPtrH("09:45"),
		PilotRole:      strPtrH("PF"),
		CompanionName:  &companion,
		AirTime:        strPtrH("01:15"),
		BlockTime:      strPtrH("01:45"),
		DutyTime:       &dutyTime,
		ApproachType:   &approachType,
		FlightType:     &flightType,
	}

	detail := ToDomainDailyLogbookDetail("logbook-1", req)

	if detail.DailyLogbookID != "logbook-1" {
		t.Errorf("expected DailyLogbookID 'logbook-1', got %s", detail.DailyLogbookID)
	}
	if detail.FlightNumber != "AV123" {
		t.Errorf("expected FlightNumber 'AV123', got %s", detail.FlightNumber)
	}
	if detail.PilotRole == nil || *detail.PilotRole != domain.PilotRolePF {
		t.Errorf("expected PilotRole PF, got %v", detail.PilotRole)
	}
	if detail.ApproachType == nil || *detail.ApproachType != domain.ApproachTypePA {
		t.Error("expected ApproachType PA")
	}

	t.Run("without approach type", func(t *testing.T) {
		req2 := CreateDailyLogbookDetailRequest{
			FlightNumber: "AV456",
			PilotRole:    strPtrH("PM"),
		}
		detail2 := ToDomainDailyLogbookDetail("logbook-2", req2)
		if detail2.ApproachType != nil {
			t.Error("expected nil ApproachType")
		}
	})
}

func TestToDomainDailyLogbookDetailUpdate(t *testing.T) {
	approachType := "VOR"
	req := UpdateDailyLogbookDetailRequest{
		FlightRealDate: "2025-01-16",
		FlightNumber:   "AV789",
		AirlineRouteID: "route-2",
		LicensePlateID: "lp-2",
		OutTime:        strPtrH("10:00"),
		TakeoffTime:    strPtrH("10:15"),
		LandingTime:    strPtrH("11:30"),
		InTime:         strPtrH("11:45"),
		PilotRole:      strPtrH("PM"),
		AirTime:        strPtrH("01:15"),
		BlockTime:      strPtrH("01:45"),
		ApproachType:   &approachType,
	}

	detail := ToDomainDailyLogbookDetailUpdate("detail-1", req)

	if detail.ID != "detail-1" {
		t.Errorf("expected ID 'detail-1', got %s", detail.ID)
	}
	if detail.FlightNumber != "AV789" {
		t.Errorf("expected FlightNumber 'AV789', got %s", detail.FlightNumber)
	}
	if detail.ApproachType == nil {
		t.Fatal("expected ApproachType to be set")
	}

	t.Run("without approach type", func(t *testing.T) {
		req2 := UpdateDailyLogbookDetailRequest{
			FlightNumber: "AV456",
			PilotRole:    strPtrH("PF"),
		}
		detail2 := ToDomainDailyLogbookDetailUpdate("detail-2", req2)
		if detail2.ApproachType != nil {
			t.Error("expected nil ApproachType")
		}
	})
}

func TestUpdateDailyLogbookDetailRequest_Sanitize(t *testing.T) {
	companion := "  Co-Pilot  "
	dutyTime := "  10:00  "
	approachType := "  ILS  "
	flightType := "  COMMERCIAL  "

	req := &UpdateDailyLogbookDetailRequest{
		FlightRealDate: "  2025-01-15  ",
		FlightNumber:   "  AV123  ",
		AirlineRouteID: "  route-1  ",
		LicensePlateID: "  lp-1  ",
		OutTime:        strPtrH("  08:00  "),
		TakeoffTime:    strPtrH("  08:15  "),
		LandingTime:    strPtrH("  09:30  "),
		InTime:         strPtrH("  09:45  "),
		PilotRole:      strPtrH("  PF  "),
		CompanionName:  &companion,
		AirTime:        strPtrH("  01:15  "),
		BlockTime:      strPtrH("  01:45  "),
		DutyTime:       &dutyTime,
		ApproachType:   &approachType,
		FlightType:     &flightType,
	}
	req.Sanitize()

	if req.FlightRealDate != "2025-01-15" {
		t.Errorf("expected trimmed FlightRealDate, got '%s'", req.FlightRealDate)
	}
	if req.FlightNumber != "AV123" {
		t.Errorf("expected trimmed FlightNumber, got '%s'", req.FlightNumber)
	}
	if req.PilotRole == nil || *req.PilotRole != "PF" {
		t.Errorf("expected trimmed PilotRole, got '%v'", req.PilotRole)
	}
	if *req.CompanionName != "Co-Pilot" {
		t.Errorf("expected trimmed CompanionName, got '%s'", *req.CompanionName)
	}
}

func TestFromDomainDailyLogbookDetail(t *testing.T) {
	approachType := domain.ApproachTypePA
	flightType := "COMMERCIAL"

	detail := &domain.DailyLogbookDetail{
		ID:                  "uuid-1",
		DailyLogbookID:      "logbook-1",
		FlightRealDate:      "2025-01-15",
		FlightNumber:        "AV123",
		AirlineRouteID:      "route-1",
		LicensePlateID:      "lp-1",
		OutTime:             strPtrH("08:00"),
		TakeoffTime:         strPtrH("08:15"),
		LandingTime:         strPtrH("09:30"),
		InTime:              strPtrH("09:45"),
		PilotRole:           domainPilotRolePtrH(domain.PilotRolePF),
		AirTime:             strPtrH("01:15"),
		BlockTime:           strPtrH("01:45"),
		ApproachType:        &approachType,
		FlightType:          &flightType,
		RouteCode:           "BOG-CLO",
		OriginIataCode:      "BOG",
		DestinationIataCode: "CLO",
		AirlineCode:         "AV",
		LogDate:             "2025-01-15",
		LicensePlate:        "HK-4567",
		ModelName:           "A320",
	}

	resp := FromDomainDailyLogbookDetail(detail, "enc-id", "enc-logbook", "enc-route", "enc-aircraft")

	if resp.ID != "enc-id" {
		t.Errorf("expected ID 'enc-id', got %s", resp.ID)
	}
	if resp.FlightNumber != "AV123" {
		t.Errorf("expected FlightNumber 'AV123', got %s", resp.FlightNumber)
	}
	if resp.ApproachType == nil || *resp.ApproachType != "PA" {
		t.Error("expected ApproachType 'ILS'")
	}
	if resp.FlightType == nil || *resp.FlightType != "COMMERCIAL" {
		t.Error("expected FlightType 'COMMERCIAL'")
	}
	if resp.RouteCode != "BOG-CLO" {
		t.Errorf("expected RouteCode 'BOG-CLO', got %s", resp.RouteCode)
	}

	t.Run("without optional fields", func(t *testing.T) {
		detail2 := &domain.DailyLogbookDetail{
			ID:           "uuid-2",
			FlightNumber: "AV456",
			PilotRole:    domainPilotRolePtrH(domain.PilotRolePM),
		}
		resp2 := FromDomainDailyLogbookDetail(detail2, "enc-2", "enc-lb", "enc-rt", "enc-ac")
		if resp2.ApproachType != nil {
			t.Error("expected nil ApproachType")
		}
		if resp2.FlightType != nil {
			t.Error("expected nil FlightType")
		}
	})
}

func TestBuildDailyLogbookLinks(t *testing.T) {
	links := BuildDailyLogbookLinks("http://localhost:8080/api/v1", "logbook-1")
	if len(links) == 0 {
		t.Error("expected HATEOAS links")
	}
	// Check that self link exists
	found := false
	for _, l := range links {
		if l.Rel == "self" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected self link")
	}
}

func TestBuildDailyLogbookListLinks(t *testing.T) {
	links := BuildDailyLogbookListLinks("http://localhost:8080/api/v1")
	if len(links) == 0 {
		t.Error("expected HATEOAS links")
	}
	found := false
	for _, l := range links {
		if l.Rel == "self" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected self link")
	}
}

func TestBuildDailyLogbookCreatedLinks(t *testing.T) {
	links := BuildDailyLogbookCreatedLinks("http://localhost:8080/api/v1", "logbook-1")
	if len(links) == 0 {
		t.Error("expected HATEOAS links")
	}
}

func TestBuildDailyLogbookDetailLinksArray(t *testing.T) {
	links := BuildDailyLogbookDetailLinksArray("http://localhost:8080/api/v1", "detail-1")
	if len(links) == 0 {
		t.Error("expected HATEOAS links")
	}
}

func TestBuildDailyLogbookDetailListLinks(t *testing.T) {
	links := BuildDailyLogbookDetailListLinks("http://localhost:8080/api/v1", "logbook-1")
	if len(links) == 0 {
		t.Error("expected HATEOAS links")
	}
}
