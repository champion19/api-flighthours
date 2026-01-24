package message

import (
	"context"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
)

func (r *repository) GetByModule(ctx context.Context, module string) ([]domain.Message, error) {
	rows, err := r.db.QueryContext(ctx, queryGetByModule, module)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []domain.Message
	for rows.Next() {
		var m domain.Message
		err := rows.Scan(
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
			return nil, err
		}
		messages = append(messages, m)
	}

	return messages, rows.Err()
}
