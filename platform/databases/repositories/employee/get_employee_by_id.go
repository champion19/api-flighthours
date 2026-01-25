package employee

import (
	"context"
	"database/sql"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
)

func (r *repository) GetEmployeeByID(ctx context.Context, id string) (*domain.Employee, error) {

	var e Employee
	err := r.db.QueryRowContext(ctx, QueryByID, id).Scan(
		&e.ID,
		&e.Name,
		&e.Airline,
		&e.Email,
		&e.IdentificationNumber,
		&e.Bp,
		&e.StartDate,
		&e.EndDate,
		&e.Active,
		&e.Role,
		&e.KeycloakUserID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrPersonNotFound
		}
		return nil, err
	}

	domainEmployee := e.ToDomain()
	return &domainEmployee, nil
}
