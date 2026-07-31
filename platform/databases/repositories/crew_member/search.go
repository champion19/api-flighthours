package crew_member

import (
	"context"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/platform/logger"
)

// SearchCrewMembers returns crew members from the pilot's own roster whose name
// or bp contains the given query (case-insensitive). An empty query returns the
// full roster. bp is matched too since it's the reliable identifier — searching
// by badge number should find the person even if the name was typed differently.
func (r *repository) SearchCrewMembers(ctx context.Context, employeeID, query string) ([]domain.CrewMember, error) {
	log.Info(logger.LogCrewMemberSearch, "employee_id", employeeID, "query", query)

	like := "%" + query + "%"
	rows, err := r.stmtSearch.QueryContext(ctx, employeeID, like, like)
	if err != nil {
		log.Error(logger.LogCrewMemberSearchError, "error", err)
		return nil, err
	}
	defer rows.Close()

	var members []domain.CrewMember
	for rows.Next() {
		var entity CrewMember
		if err := rows.Scan(&entity.ID, &entity.EmployeeID, &entity.Name, &entity.BP); err != nil {
			log.Error(logger.LogCrewMemberSearchError, "error", err)
			return nil, err
		}
		members = append(members, *entity.ToDomain())
	}

	log.Info(logger.LogCrewMemberSearchOK, "employee_id", employeeID, "count", len(members))
	return members, rows.Err()
}
