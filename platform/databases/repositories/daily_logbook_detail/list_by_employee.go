package daily_logbook_detail

import (
	"context"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/platform/logger"
)

func (r *repository) ListDailyLogbookDetailsByEmployee(ctx context.Context, employeeID string) ([]domain.DailyLogbookDetail, error) {
	log.Info(logger.LogDailyLogbookDetailList, "employee_id", employeeID)

	rows, err := r.stmtGetByEmployee.QueryContext(ctx, employeeID)
	if err != nil {
		log.Error(logger.LogDailyLogbookDetailListError, "employee_id", employeeID, "error", err)
		return nil, err
	}
	defer rows.Close()

	var details []domain.DailyLogbookDetail
	for rows.Next() {
		entity, err := scanDetail(rows)
		if err != nil {
			log.Error(logger.LogDailyLogbookDetailListError, "employee_id", employeeID, "error", err)
			return nil, err
		}
		details = append(details, *entity.ToDomain())
	}

	if err = rows.Err(); err != nil {
		log.Error(logger.LogDailyLogbookDetailListError, "employee_id", employeeID, "error", err)
		return nil, err
	}

	log.Info(logger.LogDailyLogbookDetailListOK, "employee_id", employeeID, "count", len(details))
	return details, nil
}
