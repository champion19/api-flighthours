package message

import (
	"context"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

func (r *repository) GetByModule(ctx context.Context, module string) ([]domain.Message, error) {
	rows, err := r.db.QueryContext(ctx, queryGetByModule, module)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanMessages(rows)
}
