package domain

import (
	"strconv"
	"time"

	"github.com/google/uuid"
)

type AirlineEmployee struct {
	ID        string    `json:"id"`
	AirlineID string    `json:"airline_id"`
	Bp        string    `json:"bp"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Active    bool      `json:"active"`
}

func (e *AirlineEmployee) SetID() {
	e.ID = uuid.New().String()
}

func (e *AirlineEmployee) ToLogger() []string {
	return []string{
		"id:" + e.ID,
		"airline_id:" + e.AirlineID,
		"bp:" + e.Bp,
		"start_date:" + e.StartDate.String(),
		"end_date:" + e.EndDate.String(),
		"active:" + strconv.FormatBool(e.Active),
	}
}
