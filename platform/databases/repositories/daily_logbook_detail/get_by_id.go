package daily_logbook_detail

import (
	"context"
	"database/sql"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/platform/logger"
)

func (r *repository) GetDailyLogbookDetailByID(ctx context.Context, id string) (*domain.DailyLogbookDetail, error) {
	log.Info(logger.LogDailyLogbookDetailGet, "id", id)

	row := r.stmtGetByID.QueryRowContext(ctx, id)
	entity, err := scanDetail(row)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Warn(logger.LogDailyLogbookDetailNotFound, "id", id)
			return nil, nil
		}
		log.Error(logger.LogDailyLogbookDetailGetError, "id", id, "error", err)
		return nil, err
	}

	log.Info(logger.LogDailyLogbookDetailGetOK, "id", id)
	return entity.ToDomain(), nil
}
