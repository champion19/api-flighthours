package daily_logbook_detail

import (
	"context"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/platform/logger"
)

func (r *repository) ListDailyLogbookDetailsByLogbook(ctx context.Context, logbookID string) ([]domain.DailyLogbookDetail, error) {
	log.Info(logger.LogDailyLogbookDetailList, "logbook_id", logbookID)

	rows, err := r.stmtGetByLogbook.QueryContext(ctx, logbookID)
	if err != nil {
		log.Error(logger.LogDailyLogbookDetailListError, "logbook_id", logbookID, "error", err)
		return nil, err
	}
	defer rows.Close()

	var details []domain.DailyLogbookDetail
	for rows.Next() {
		var entity DailyLogbookDetail
		err := rows.Scan(
			&entity.ID,
			&entity.DailyLogbookID,
			&entity.FlightRealDate,
			&entity.FlightNumber,
			&entity.AirlineRouteID,
			&entity.TailNumberID,
			&entity.Passengers,
			&entity.OutTime,
			&entity.TakeoffTime,
			&entity.LandingTime,
			&entity.InTime,
			&entity.PilotRole,
			&entity.CompanionName,
			&entity.CrewRole,
			&entity.AirTime,
			&entity.BlockTime,
			&entity.ApproachType,
			&entity.FlightType,
			&entity.EmployeeLogbookID,
			&entity.LogDate,
			&entity.TailNumber,
			&entity.ModelName,
			&entity.RouteCode,
			&entity.OriginIataCode,
			&entity.DestinationIataCode,
			&entity.AirlineCode,
		)
		if err != nil {
			log.Error(logger.LogDailyLogbookDetailListError, "logbook_id", logbookID, "error", err)
			return nil, err
		}
		details = append(details, *entity.ToDomain())
	}

	if err = rows.Err(); err != nil {
		log.Error(logger.LogDailyLogbookDetailListError, "logbook_id", logbookID, "error", err)
		return nil, err
	}

	log.Info(logger.LogDailyLogbookDetailListOK, "logbook_id", logbookID, "count", len(details))
	return details, nil
}
