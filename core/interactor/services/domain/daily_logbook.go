package domain

import (
	"time"

	"github.com/google/uuid"
)

type DailyLogbook struct {
	ID           string    `json:"id"`
	LogDate      time.Time `json:"log_date"`
	EmployeeID   string    `json:"employee_id"`
	BookPage     *int      `json:"book_page,omitempty"`
	Status       bool      `json:"status"`
	TailNumberID *string   `json:"tail_number_id,omitempty"`
	TailNumber   *string   `json:"tail_number,omitempty"` // denormalized plate, read-only
	CrewRole     *CrewRole `json:"crew_role,omitempty"`   // default crew role for every flight logged under this book page
}

func (d *DailyLogbook) SetID() {
	d.ID = uuid.New().String()
}

func (d *DailyLogbook) ToLogger() []string {
	return []string{
		"id:" + d.ID,
		"employee_id:" + d.EmployeeID,
		"log_date:" + d.LogDate.Format("2006-01-02"),
	}
}

func (d *DailyLogbook) IsActive() bool {
	return d.Status
}

func (d *DailyLogbook) StatusString() string {
	if d.Status {
		return "active"
	}
	return "inactive"
}
