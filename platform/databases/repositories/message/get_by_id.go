package message

import (
	"context"
	"database/sql"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
)

func (r *repository) GetByID(ctx context.Context, id string) (*domain.Message, error) {
	var m domain.Message

	err := r.db.QueryRowContext(ctx, queryGetByID, id).Scan(
		&m.ID,
		&m.Code,
		&m.Type,
		&m.Category,
		&m.Module,
		&m.Title,
		&m.Content,
		&m.Active,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil

}
