package crew_member

import (
	"context"
	"database/sql"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/platform/logger"
)

// GetCrewMemberByID fetches a single crew member by ID, or nil if not found.
func (r *repository) GetCrewMemberByID(ctx context.Context, id string) (*domain.CrewMember, error) {
	var entity CrewMember
	err := r.stmtGetByID.QueryRowContext(ctx, id).Scan(&entity.ID, &entity.EmployeeID, &entity.Name, &entity.BP)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Error(logger.LogCrewMemberSearchError, "id", id, "error", err)
		return nil, err
	}
	return entity.ToDomain(), nil
}
