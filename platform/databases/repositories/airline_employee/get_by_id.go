package airline_employee

import (
	"context"
	"database/sql"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
)

func (r *repository) GetAirlineEmployeeByID(ctx context.Context, id string) (*domain.AirlineEmployee, error) {
	var e AirlineEmployee

	err := r.stmtGetByID.QueryRowContext(ctx, id).Scan(
		&e.ID,
		&e.AirlineID,
		&e.Bp,
		&e.StartDate,
		&e.EndDate,
		&e.Active,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrAirlineEmployeeNotFound
		}
		return nil, err
	}

	return e.ToDomain(), nil
}
