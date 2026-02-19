package airline_employee

import (
	"context"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/databases/common"
)

func (r *repository) AddAirlineEmployee(ctx context.Context, tx output.Tx, employee domain.AirlineEmployee) error {
	employeeToUpdate := FromDomain(&employee)

	log.Debug("AddAirlineEmployee: Starting add",
		"employee_id", employeeToUpdate.ID,
		"airline_id", employeeToUpdate.AirlineID,
		"active", employeeToUpdate.Active)

	dbTx, err := common.CastTx(tx)
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

	log.Debug("AddAirlineEmployee: Query executed successfully",
		"employee_id", employeeToUpdate.ID,
		"rows_affected", func() int64 { r, _ := result.RowsAffected(); return r }())

	return nil
}
