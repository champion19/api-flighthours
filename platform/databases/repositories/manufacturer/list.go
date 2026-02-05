package manufacturer

import(
	"context"
  domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
)

// ListManufacturers retrieves all manufacturers
func (r *repository) ListManufacturers(ctx context.Context) ([]domain.Manufacturer, error) {
	rows, err := r.stmtGetAll.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var manufacturers []domain.Manufacturer
	for rows.Next() {
		var m Manufacturer
		if err := rows.Scan(&m.ID, &m.Name); err != nil {
			return nil, err
		}
		manufacturers = append(manufacturers, *m.ToDomain())
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return manufacturers, nil
}

