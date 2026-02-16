package aircraftmodel

import (
	"context"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
)

// GetAircraftModelsByFamily retrieves all aircraft models for a specific family (HU32)
func (r *repository) GetAircraftModelsByFamily(ctx context.Context, family string) ([]domain.AircraftModel, error) {
	rows, err := r.stmtGetByFamily.QueryContext(ctx, family)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var models []domain.AircraftModel
	for rows.Next() {
		model, err := scanAircraftModel(rows)
		if err != nil {
			return nil, err
		}
		models = append(models, *model.ToDomain())
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return models, nil
}
