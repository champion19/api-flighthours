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
		var entity DailyLogbookDetail
		err := rows.Scan(
			&entity.ID,
			&entity.DailyLogbookID,
			&entity.FlightRealDate,
			&entity.FlightNumber,
			&entity.AirlineRouteID,
			&entity.LicensePlateID,
			&entity.Passengers,
			&entity.OutTime,
			&entity.TakeoffTime,
			&entity.LandingTime,
			&entity.InTime,
			&entity.PilotRole,
			&entity.CompanionName,
			&entity.AirTime,
			&entity.BlockTime,
			&entity.DutyTime,
			&entity.ApproachType,
			&entity.FlightType,
			&entity.EmployeeLogbookID,
			&entity.LogDate,
			&entity.LicensePlate,
			&entity.ModelName,
			&entity.RouteCode,
			&entity.OriginIataCode,
			&entity.DestinationIataCode,
			&entity.AirlineCode,
		)
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
