package message

import (
	"context"
	"database/sql"

	cachetypes "github.com/champion19/flighthours-api/platform/cache/types"
)

// GetByCodeWithStatusForCache returns a message for cache by code without active filter
func (r *repository) GetByCodeWithStatusForCache(ctx context.Context, code string) (*cachetypes.CachedMessage, error) {
	var m cachetypes.CachedMessage
	var createdAt, updatedAt interface{}

	err := r.db.QueryRowContext(ctx, queryGetByCodeWithStatus, code).Scan(
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
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &m, nil
}
