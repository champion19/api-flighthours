package crew_member

import (
	"database/sql"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

type CrewMember struct {
	ID         string
	EmployeeID string
	Name       string
	BP         sql.NullString
}

func (c *CrewMember) ToDomain() *domain.CrewMember {
	m := &domain.CrewMember{
		ID:         c.ID,
		EmployeeID: c.EmployeeID,
		Name:       c.Name,
	}
	if c.BP.Valid {
		m.BP = &c.BP.String
	}
	return m
}
