package engine

import (
	"context"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
)

// ListEngines retrieves all engines
func (r *repository) ListEngines(ctx context.Context) ([]domain.Engine, error) {
	rows, err := r.stmtGetAll.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var engines []domain.Engine
	for rows.Next() {
		var e Engine
		if err := rows.Scan(&e.ID, &e.Name); err != nil {
			return nil, err
		}
		engines = append(engines, *e.ToDomain())
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return engines, nil
}
