package daily_logbook_detail

import (
	"context"
	"strings"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/databases/common"
	"github.com/champion19/api-flighthours/platform/logger"
)

func (r *repository) SaveDailyLogbookDetail(ctx context.Context, tx output.Tx, detail domain.DailyLogbookDetail) error {
	log.Info(logger.LogDailyLogbookDetailCreate, "data", detail.ToLogger())

	entity := FromDomain(&detail)

	sqlTx, ok := tx.(*common.SQLTX)
	if !ok {
		log.Error(logger.LogDailyLogbookDetailCreateError, "error", "invalid transaction type")
		return domain.ErrInvalidTransaction
	}

	stmt := sqlTx.Tx.StmtContext(ctx, r.stmtInsert)

	_, err := stmt.ExecContext(ctx,
		entity.ID,
		entity.DailyLogbookID,
		entity.FlightRealDate,
		entity.FlightNumber,
		entity.AirlineRouteID,
		entity.LicensePlateID,
		entity.Passengers,
		entity.OutTime,
		entity.TakeoffTime,
		entity.LandingTime,
		entity.InTime,
		entity.PilotRole,
		entity.CompanionName,
		entity.CrewRole,
		entity.AirTime,
		entity.BlockTime,
		entity.DutyTime,
		entity.ApproachType,
		entity.FlightType,
		entity.EmployeeLogbookID,
	)

	if err != nil {
		log.Error(logger.LogDailyLogbookDetailCreateError, "error", err.Error())

		errStr := err.Error()
		if strings.Contains(errStr, "foreign key constraint") || strings.Contains(errStr, "1452") {
			if strings.Contains(errStr, "daily_logbook") {
				return domain.ErrFlightInvalidLogbook
			}
			if strings.Contains(errStr, "airline_route") {
				return domain.ErrFlightInvalidRoute
			}
			if strings.Contains(errStr, "aircraft_registration") {
				return domain.ErrFlightInvalidAircraft
			}
		}
		return domain.ErrFlightCannotSave
	}

	log.Info(logger.LogDailyLogbookDetailCreateOK, "id", detail.ID)
	return nil
}
