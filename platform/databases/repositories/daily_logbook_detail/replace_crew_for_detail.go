package daily_logbook_detail

import (
	"context"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/databases/common"
	"github.com/champion19/api-flighthours/platform/logger"
	"github.com/google/uuid"
)

const (
	queryDeleteCrewByDetail  = `DELETE FROM daily_logbook_detail_crew WHERE daily_logbook_detail_id = ?`
	queryInsertCrewForDetail = `
		INSERT INTO daily_logbook_detail_crew (id, daily_logbook_detail_id, crew_member_id, role)
		VALUES (?, ?, ?, ?)`
)

// ReplaceCrewForDetail atomically replaces the crew assigned to a flight leg: clears whatever
// was there before and inserts the given assignments. Called within the same transaction as
// creating/updating the flight itself, so a bad assignment rolls back the whole save — same
// lax validation posture the rest of this entity already has for pilot_role/crew_role.
func (r *repository) ReplaceCrewForDetail(ctx context.Context, tx output.Tx, detailID string, assignments []domain.CrewAssignment) error {
	log.Info(logger.LogCrewAssignmentReplace, "daily_logbook_detail_id", detailID, "count", len(assignments))

	sqlTx, err := common.CastTx(tx)
	if err != nil {
		log.Error(logger.LogCrewAssignmentReplaceErr, "error", "invalid transaction type")
		return err
	}

	if _, err := sqlTx.Tx.ExecContext(ctx, queryDeleteCrewByDetail, detailID); err != nil {
		log.Error(logger.LogCrewAssignmentReplaceErr, "error", err)
		return err
	}

	for _, a := range assignments {
		id := uuid.New().String()
		if _, err := sqlTx.Tx.ExecContext(ctx, queryInsertCrewForDetail, id, detailID, a.CrewMemberID, string(a.Role)); err != nil {
			log.Error(logger.LogCrewAssignmentReplaceErr, "error", err)
			return err
		}
	}

	log.Info(logger.LogCrewAssignmentReplaceOK, "daily_logbook_detail_id", detailID, "count", len(assignments))
	return nil
}
