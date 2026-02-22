package daily_logbook_detail

import (
	"context"
	"database/sql"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/platform/logger"
)

func (r *repository) GetDailyLogbookDetailByID(ctx context.Context, id string) (*domain.DailyLogbookDetail, error) {
	log.Info(logger.LogDailyLogbookDetailGet, "id", id)

	var entity DailyLogbookDetail
	err := r.stmtGetByID.QueryRowContext(ctx, id).Scan(
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
