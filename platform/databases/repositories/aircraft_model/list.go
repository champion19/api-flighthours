package aircraftmodel

import(
	"context"
	"database/sql"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
)

// ListAircraftModels retrieves all aircraft models with optional filters
func (r *repository) ListAircraftModels(ctx context.Context, filters map[string]interface{}) ([]domain.AircraftModel, error) {
	var rows *sql.Rows
	var err error

	// Check if filtering by engine type
	if engineType, ok := filters["engine_type"]; ok {
		rows, err = r.stmtGetByEngineType.QueryContext(ctx, engineType)
	} else {
		rows, err = r.stmtGetAll.QueryContext(ctx)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var models []domain.AircraftModel
	for rows.Next() {
		var model AircraftModel
		var engineTypeName sql.NullString
		var manufacturer sql.NullString

		if err := rows.Scan(
			&model.ID,
			&model.ModelName,
			&model.AircraftTypeName,
			&engineTypeName,
			&model.Family,
			&manufacturer,
			&model.Status,
		); err != nil {
			return nil, err
		}

		if engineTypeName.Valid {
			model.EngineTypeName = engineTypeName.String
		}
		if manufacturer.Valid {
			model.Manufacturer = manufacturer.String
		}

		models = append(models, *model.ToDomain())
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return models, nil
}
