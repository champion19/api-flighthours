package message

import (
	"context"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

func (r *repository) GetByType(ctx context.Context, msgType string) ([]domain.Message, error) {
	rows, err := r.db.QueryContext(ctx, queryGetByType, msgType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanMessages(rows)
}
