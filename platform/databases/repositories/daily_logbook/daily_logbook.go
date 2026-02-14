package daily_logbook

import (
	"time"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
)

// DailyLogbook is the database entity for daily_logbook table
type DailyLogbook struct {
	ID         string    `db:"id"`
	LogDate    time.Time `db:"log_date"`
	EmployeeID string    `db:"employee_id"`
	BookPage   *int64    `db:"book_page"`
	Status     bool      `db:"status"`
}

// ToDomain converts the database entity to domain model
func (d *DailyLogbook) ToDomain() *domain.DailyLogbook {
	return &domain.DailyLogbook{
		ID:         d.ID,
		LogDate:    d.LogDate,
		EmployeeID: d.EmployeeID,
		BookPage:   d.BookPage,
		Status:     d.Status,
	}
}

// FromDomain converts a domain model to database entity
func FromDomain(domainLogbook *domain.DailyLogbook) *DailyLogbook {
	return &DailyLogbook{
		ID:         domainLogbook.ID,
		LogDate:    domainLogbook.LogDate,
		EmployeeID: domainLogbook.EmployeeID,
		BookPage:   domainLogbook.BookPage,
		Status:     domainLogbook.Status,
	}
}
