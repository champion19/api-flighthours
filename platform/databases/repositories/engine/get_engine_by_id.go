package engine

import (
	"context"
	"database/sql"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
)

// GetEngineByID retrieves an engine by ID
func (r *repository) GetEngineByID(ctx context.Context, id string) (*domain.Engine, error) {
	var e Engine
	err := r.stmtGetByID.QueryRowContext(ctx, id).Scan(
		&e.ID,
		&e.Name,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrEngineNotFound
		}
		return nil, err
	}
	return e.ToDomain(), nil
}
