package crew_member

import (
	"context"
	"database/sql"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/databases/common"
	"github.com/champion19/api-flighthours/platform/logger"
)

// FindOrCreateCrewMember adds a person to the pilot's roster, or returns the
// existing one if already registered for this employee — this is what lets
// "add if missing" and "have I flown with this person before" share the same
// simple flow. bp (badge/carné number) is the reliable identifier, so when
// provided it's checked first: the same bp always resolves to the same person
// even if the name was typed differently this time. Falls back to matching by
// exact name (case-insensitive) when bp is absent or not yet on record.
func (r *repository) FindOrCreateCrewMember(ctx context.Context, tx output.Tx, employeeID, name string, bp *string) (*domain.CrewMember, error) {
	log.Info(logger.LogCrewMemberCreate, "employee_id", employeeID, "name", name)

	sqlTx, err := common.CastTx(tx)
	if err != nil {
		log.Error(logger.LogCrewMemberCreateError, "error", "invalid transaction type")
		return nil, err
	}

	var bpValue sql.NullString
	if bp != nil && *bp != "" {
		bpValue = sql.NullString{String: *bp, Valid: true}

		getByBPStmt := sqlTx.Tx.StmtContext(ctx, r.stmtGetByEmployeeAndBP)
		var existing CrewMember
		err := getByBPStmt.QueryRowContext(ctx, employeeID, *bp).Scan(&existing.ID, &existing.EmployeeID, &existing.Name, &existing.BP)
		if err == nil {
			log.Info(logger.LogCrewMemberCreateOK, "id", existing.ID, "employee_id", employeeID, "matched_by", "bp")
			return existing.ToDomain(), nil
		}
		if err != sql.ErrNoRows {
			log.Error(logger.LogCrewMemberCreateError, "error", err)
			return nil, err
		}
	}

	member := domain.CrewMember{EmployeeID: employeeID, Name: name}
	member.SetID()

	insertStmt := sqlTx.Tx.StmtContext(ctx, r.stmtInsertIgnore)
	if _, err := insertStmt.ExecContext(ctx, member.ID, employeeID, name, bpValue); err != nil {
		log.Error(logger.LogCrewMemberCreateError, "error", err)
		return nil, err
	}

	getStmt := sqlTx.Tx.StmtContext(ctx, r.stmtGetByEmployeeAndName)
	var entity CrewMember
	err = getStmt.QueryRowContext(ctx, employeeID, name).Scan(&entity.ID, &entity.EmployeeID, &entity.Name, &entity.BP)
	if err != nil {
		log.Error(logger.LogCrewMemberCreateError, "error", err)
		return nil, err
	}

	log.Info(logger.LogCrewMemberCreateOK, "id", entity.ID, "employee_id", employeeID)
	return entity.ToDomain(), nil
}
