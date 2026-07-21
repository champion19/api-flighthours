package daily_logbook

import (
	"database/sql"
	"time"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
)

// DailyLogbook is the database entity for daily_logbook table
type DailyLogbook struct {
	ID           string         `db:"id"`
	LogDate      time.Time      `db:"log_date"`
	EmployeeID   string         `db:"employee_id"`
	BookPage     *int           `db:"book_page"`
	Status       bool           `db:"status"`
	TailNumberID sql.NullString `db:"tail_number_id"`
	TailNumber   sql.NullString `db:"tail_number"` // denormalized from LEFT JOIN tail_number
}

// ToDomain converts the database entity to domain model
func (d *DailyLogbook) ToDomain() *domain.DailyLogbook {
	logbook := &domain.DailyLogbook{
		ID:         d.ID,
		LogDate:    d.LogDate,
		EmployeeID: d.EmployeeID,
		BookPage:   d.BookPage,
		Status:     d.Status,
	}

	if d.TailNumberID.Valid {
		logbook.TailNumberID = &d.TailNumberID.String
	}
	if d.TailNumber.Valid {
		logbook.TailNumber = &d.TailNumber.String
	}

	return logbook
}

// FromDomain converts a domain model to database entity
func FromDomain(domainLogbook *domain.DailyLogbook) *DailyLogbook {
	logbook := &DailyLogbook{
		ID:         domainLogbook.ID,
		LogDate:    domainLogbook.LogDate,
		EmployeeID: domainLogbook.EmployeeID,
		BookPage:   domainLogbook.BookPage,
		Status:     domainLogbook.Status,
	}

	if domainLogbook.TailNumberID != nil {
		logbook.TailNumberID = sql.NullString{String: *domainLogbook.TailNumberID, Valid: true}
	}

	return logbook
}
