package airline_employee

import (
	"context"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
)

// UpdateAirlineEmployee updates airline-specific fields for an existing employee (HU25)
// The employee must already have airline info assigned
func (r *repository) UpdateAirlineEmployee(ctx context.Context, tx output.Tx, employee domain.AirlineEmployee) error {
	employeeToUpdate := FromDomain(&employee)

	log.Debug("UpdateAirlineEmployee: Starting update",
		"employee_id", employeeToUpdate.ID,
		"airline_id", employeeToUpdate.AirlineID,
		"active", employeeToUpdate.Active)

	dbTx, err := castTx(tx)
	if err != nil {
		return err
	}

	result, err := dbTx.ExecContext(ctx, QueryUpdateAirlineInfo,
		employeeToUpdate.AirlineID,
		employeeToUpdate.Bp,
		employeeToUpdate.StartDate,
		employeeToUpdate.EndDate,
		employeeToUpdate.Active,
		employeeToUpdate.ID,
	)

	if err != nil {
		return handleMySQLError(err, employee.ID)
	}

	log.Debug("UpdateAirlineEmployee: Query executed successfully",
		"employee_id", employeeToUpdate.ID,
		"rows_affected", func() int64 { r, _ := result.RowsAffected(); return r }())

	return nil
}
