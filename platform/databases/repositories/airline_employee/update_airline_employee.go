package airline_employee

import (
	"context"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/databases/common"
	"github.com/champion19/api-flighthours/platform/logger"
	"github.com/go-sql-driver/mysql"
)

// UpdateAirlineEmployee updates airline-specific fields for an existing employee (HU25)
// The employee must already have airline info assigned
func (r *repository) UpdateAirlineEmployee(ctx context.Context, tx output.Tx, employee domain.AirlineEmployee) error {
	employeeToUpdate := FromDomain(&employee)

	log.Debug("UpdateAirlineEmployee: Starting update",
		"employee_id", employeeToUpdate.ID,
		"airline_id", employeeToUpdate.AirlineID,
		"active", employeeToUpdate.Active)

	// Cast the transaction to the concrete type
	dbTx, ok := tx.(*common.SQLTX)
	if !ok {
		log.Error(logger.LogDatabaseUnavailable, "error", logger.LogErrInvalidTransaction)
		return domain.ErrInvalidTransaction
	}

	// Update airline-specific fields
	result, err := dbTx.ExecContext(ctx, QueryUpdateAirlineInfo,
		employeeToUpdate.AirlineID,
		employeeToUpdate.Bp,
		employeeToUpdate.StartDate,
		employeeToUpdate.EndDate,
		employeeToUpdate.Active,
		employeeToUpdate.ID,
	)

	if err != nil {
		// Check for specific MySQL errors
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			case 1452:
				// Foreign key constraint fails (e.g., invalid airline)
				log.Error(logger.LogDatabaseUnavailable,
					"employee_id", employee.ID,
					"error", "invalid foreign key reference",
					"mysql_error", mysqlErr.Message)
				return domain.ErrInvalidForeignKey
			}
		}
		log.Error(logger.LogDatabaseUnavailable, "employee_id", employee.ID, "error", err)
		return err
	}

	log.Debug("UpdateAirlineEmployee: Query executed successfully",
		"employee_id", employeeToUpdate.ID,
		"rows_affected", func() int64 { r, _ := result.RowsAffected(); return r }())

	return nil
}
