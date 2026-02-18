package aircraftmodel

import (
	"context"
	"database/sql"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
)

// GetAircraftModelByID retrieves an aircraft model by ID
func (r *repository) GetAircraftModelByID(ctx context.Context, id string) (*domain.AircraftModel, error) {
	row := r.stmtGetByID.QueryRowContext(ctx, id)

	model, err := scanAircraftModel(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrAircraftModelNotFound
		}
		return nil, err
	}

	return model.ToDomain(), nil
}
