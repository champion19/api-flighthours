package daily_logbook_detail

import (
	"context"
	"strings"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/databases/common"
	"github.com/champion19/api-flighthours/platform/logger"
)

func (r *repository) UpdateDailyLogbookDetail(ctx context.Context, tx output.Tx, detail domain.DailyLogbookDetail) error {
	log.Info(logger.LogDailyLogbookDetailUpdate, "data", detail.ToLogger())

	entity := FromDomain(&detail)

	sqlTx, err := common.CastTx(tx)
	if err != nil {
		log.Error(logger.LogDailyLogbookDetailUpdateError, "error", "invalid transaction type")
		return err
	}

	stmt := sqlTx.Tx.StmtContext(ctx, r.stmtUpdate)

	_, err = stmt.ExecContext(ctx,
		entity.FlightRealDate,
		entity.FlightNumber,
		entity.OriginAirportID,
		entity.DestinationAirportID,
		entity.TailNumberID,
		entity.Passengers,
		entity.OutTime,
		entity.TakeoffTime,
		entity.LandingTime,
		entity.InTime,
		entity.PilotRole,
		entity.CrewRole,
		entity.AirTime,
		entity.BlockTime,
		entity.ApproachCategory,
		entity.ApproachSubtype,
		entity.Autoland,
		entity.FlightType,
		entity.ID,
	)

	if err != nil {
		log.Error(logger.LogDailyLogbookDetailUpdateError, "error", err.Error())

		errStr := err.Error()
		if strings.Contains(errStr, "foreign key constraint") || strings.Contains(errStr, "1452") {
			if fkErr := handleDetailFKError(errStr); fkErr != nil {
				return fkErr
			}
		}
		return domain.ErrFlightCannotUpdate
	}

	log.Info(logger.LogDailyLogbookDetailUpdateOK, "id", detail.ID)
	return nil
}
