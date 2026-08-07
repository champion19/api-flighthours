package interactor

import (
	"context"
	"strings"

	"github.com/champion19/api-flighthours/core/interactor/helpers"
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/input"
	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/logger"
)

// CrewMemberInteractor orchestrates a pilot's private crew roster (command crew +
// cabin crew they've flown with before), always scoped to the authenticated employee.
type CrewMemberInteractor struct {
	service input.CrewMemberService
}

func NewCrewMemberInteractor(service input.CrewMemberService) *CrewMemberInteractor {
	return &CrewMemberInteractor{service: service}
}

// SearchCrewMembers returns the pilot's own roster, optionally filtered by name.
func (i *CrewMemberInteractor) SearchCrewMembers(ctx context.Context, traceID, employeeID, query string) ([]domain.CrewMember, error) {
	log.Info(logger.LogCrewMemberSearch, "trace_id", traceID, "employee_id", employeeID, "query", query)

	members, err := i.service.SearchCrewMembers(ctx, employeeID, query)
	if err != nil {
		log.Error(logger.LogCrewMemberSearchError, "trace_id", traceID, "error", err)
		return nil, err
	}

	log.Info(logger.LogCrewMemberSearchOK, "trace_id", traceID, "count", len(members))
	return members, nil
}

// CreateCrewMember adds a person to the pilot's roster, or returns the existing one
// if a crew member with that name is already registered for this employee.
func (i *CrewMemberInteractor) CreateCrewMember(ctx context.Context, traceID, employeeID, name string, bp *string) (*domain.CrewMember, error) {
	log.Info(logger.LogCrewMemberCreate, "trace_id", traceID, "employee_id", employeeID, "name", name)

	if strings.TrimSpace(name) == "" {
		return nil, domain.ErrInvalidRequest
	}

	var result *domain.CrewMember
	err := helpers.RunWithTx(ctx, i.service, log, logger.LogCrewMemberCreateError,
		func(ctx context.Context, tx output.Tx) error {
			member, err := i.service.FindOrCreateCrewMemberTx(ctx, tx, employeeID, name, bp)
			if err != nil {
				return err
			}
			result = member
			return nil
		})
	if err != nil {
		log.Error(logger.LogCrewMemberCreateError, "trace_id", traceID, "error", err)
		return nil, err
	}

	log.Info(logger.LogCrewMemberCreateOK, "trace_id", traceID, "id", result.ID)
	return result, nil
}
