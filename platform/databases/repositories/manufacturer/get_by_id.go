package manufacturer

import (
	"context"
	"database/sql"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
)

func (r *repository) GetManufacturerByID(ctx context.Context, id string) (*domain.Manufacturer, error) {
	var m Manufacturer
	err := r.stmtGetByID.QueryRowContext(ctx, id).Scan(
		&m.ID,
		&m.Name,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrManufacturerNotFound
		}
		return nil, err
	}
	return m.ToDomain(), nil
}
