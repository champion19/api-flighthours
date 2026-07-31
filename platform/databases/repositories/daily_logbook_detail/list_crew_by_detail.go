package daily_logbook_detail

import (
	"context"
	"database/sql"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/platform/logger"
)

const queryCrewByDetail = `
	SELECT dldc.id, dldc.crew_member_id, cm.name, cm.bp, dldc.role
	FROM daily_logbook_detail_crew dldc
	INNER JOIN crew_member cm ON dldc.crew_member_id = cm.id
	WHERE dldc.daily_logbook_detail_id = ?
	ORDER BY dldc.created_at ASC`

// ListCrewByDetail returns the First Officer + cabin crew assigned to a single flight leg.
func (r *repository) ListCrewByDetail(ctx context.Context, detailID string) ([]domain.CrewAssignment, error) {
	log.Info(logger.LogCrewAssignmentList, "daily_logbook_detail_id", detailID)

	rows, err := r.db.QueryContext(ctx, queryCrewByDetail, detailID)
	if err != nil {
		log.Error(logger.LogCrewAssignmentListError, "error", err)
		return nil, err
	}
	defer rows.Close()

	var assignments []domain.CrewAssignment
	for rows.Next() {
		var a domain.CrewAssignment
		var role string
		var bp sql.NullString
		if err := rows.Scan(&a.ID, &a.CrewMemberID, &a.Name, &bp, &role); err != nil {
			log.Error(logger.LogCrewAssignmentListError, "error", err)
			return nil, err
		}
		if bp.Valid {
			a.BP = &bp.String
		}
		a.Role = domain.CrewMemberRole(role)
		assignments = append(assignments, a)
	}

	log.Info(logger.LogCrewAssignmentListOK, "daily_logbook_detail_id", detailID, "count", len(assignments))
	return assignments, rows.Err()
}
