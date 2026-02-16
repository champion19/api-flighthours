package message

import (
	"context"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

func (r *repository) GetAllActive(ctx context.Context) ([]domain.Message, error) {
	rows, err := r.db.QueryContext(ctx, queryGetAllActive)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanMessages(rows)
}
