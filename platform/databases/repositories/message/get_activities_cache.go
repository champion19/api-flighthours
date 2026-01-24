package message

import (
	"context"

	cachetypes "github.com/champion19/flighthours-api/platform/cache/types"
)

func (r *repository) GetAllActiveForCache(ctx context.Context) ([]cachetypes.CachedMessage, error) {
	rows, err := r.db.QueryContext(ctx, queryGetAllActive)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []cachetypes.CachedMessage
	for rows.Next() {
		var m cachetypes.CachedMessage
		var createdAt, updatedAt interface{}
		err := rows.Scan(
			&m.ID,
			&m.Code,
			&m.Type,
			&m.Category,
			&m.Module,
			&m.Title,
			&m.Content,
			&m.Active,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}

	return messages, rows.Err()
}


