package daily_logbook

import (
	"testing"
	"time"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
)

func TestDailyLogbook_ToDomain(t *testing.T) {
	page := 5
	now := time.Now()
	dl := &DailyLogbook{
		ID: "dl1", LogDate: now, EmployeeID: "emp1", BookPage: &page, Status: true,
	}
	result := dl.ToDomain()

	if result.ID != "dl1" {
		t.Errorf("expected dl1, got %s", result.ID)
	}
	if result.EmployeeID != "emp1" {
		t.Errorf("expected emp1, got %s", result.EmployeeID)
	}
	if *result.BookPage != 5 {
		t.Errorf("expected BookPage 5, got %d", *result.BookPage)
	}
}

func TestDailyLogbook_FromDomain(t *testing.T) {
	page := 10
	now := time.Now()
	dm := &domain.DailyLogbook{
		ID: "dl2", LogDate: now, EmployeeID: "emp2", BookPage: &page, Status: false,
	}
	result := FromDomain(dm)

	if result.ID != "dl2" {
		t.Errorf("expected dl2, got %s", result.ID)
	}
	if result.Status != false {
		t.Error("expected Status false")
	}
}

func TestDailyLogbook_RoundTrip(t *testing.T) {
	page := 3
	now := time.Now()
	original := &domain.DailyLogbook{
		ID: "dl3", LogDate: now, EmployeeID: "emp3", BookPage: &page, Status: true,
	}
	restored := FromDomain(original).ToDomain()

	if restored.ID != original.ID {
		t.Errorf("ID mismatch")
	}
	if restored.EmployeeID != original.EmployeeID {
		t.Errorf("EmployeeID mismatch")
	}
	if *restored.BookPage != *original.BookPage {
		t.Errorf("BookPage mismatch")
	}
}

func TestDailyLogbook_NilBookPage(t *testing.T) {
	dl := &DailyLogbook{
		ID: "dl4", LogDate: time.Now(), EmployeeID: "emp4", BookPage: nil, Status: true,
	}
	result := dl.ToDomain()
	if result.BookPage != nil {
		t.Error("expected nil BookPage")
	}
}
