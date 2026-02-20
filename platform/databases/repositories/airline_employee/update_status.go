package airline_employee

import (
	"context"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/databases/common"
)

// UpdateAirlineEmployeeStatus updates only the active status of an airline employee (HU27/HU28)
// The employee must already have airline info assigned (airline IS NOT NULL)
// This operation is IDEMPOTENT - calling it multiple times with the same status returns success
func (r *repository) UpdateAirlineEmployeeStatus(ctx context.Context, tx output.Tx, id string, active bool) error {
	log.Debug("UpdateAirlineEmployeeStatus:Starting update",
		"employee_id", id,
		"active", active)

	dbTx, err := common.CastTx(tx)
	if err != nil {
		return err
	}

	// First, check if the employee has airline info assigned
	var hasAirline bool
	err = dbTx.QueryRowContext(ctx, "SELECT airline IS NOT NULL FROM employee WHERE id = ?", id).Scan(&hasAirline)
	if err != nil {
		log.Error("UpdateAirlineEmployeeStatus: employee not found", "employee_id", id, "error", err)
		return domain.ErrAirlineEmployeeNotFound
	}

	if !hasAirline {
		log.Warn("UpdateAirlineEmployeeStatus: Employee has no airline info assigned", "employee_id", id)
		return domain.ErrAirlineEmployeeNotFound
	}

	// Update the active status (idempotent - OK if already in desired state)
	_, err = dbTx.ExecContext(ctx, QueryUpdateStatus, active, id)
	if err != nil {
		log.Error("UpdateAirlineEmployeeStatus: failed to update status",
			"employee_id", id,
			"error", err)
		return err
	}

	log.Debug("UpdateAirlineEmployeeStatus: Status updated successfully (idempotent)",
		"employee_id", id,
		"active", active)

	return nil
}
