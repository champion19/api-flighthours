package employee

import (
	"context"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/output"
	"github.com/champion19/flighthours-api/platform/databases/common"
	"github.com/champion19/flighthours-api/platform/logger"
	"github.com/go-sql-driver/mysql"
)

func (r *repository) UpdateEmployee(ctx context.Context, tx output.Tx, employee domain.Employee) error {
	employeeToUpdate := FromDomain(employee)

	log.Debug("UpdateEmployee: Starting update",
		"employee_id", employeeToUpdate.ID,
		"name", employeeToUpdate.Name,
		"airline", employeeToUpdate.Airline,
		"email", employeeToUpdate.Email,
		"active", employeeToUpdate.Active)

	// Cast the transaction to the concrete type
	dbTx, ok := tx.(*common.SQLTX)
	if !ok {
		log.Error(logger.LogEmployeeUpdateError, "error", "invalid transaction type")
		return domain.ErrInvalidTransaction
	}

	log.Debug("UpdateEmployee: Transaction cast successful, executing query")

	// Execute the update within the transaction
	result, err := dbTx.ExecContext(ctx, QueryUpdate,
		employeeToUpdate.Name,
		employeeToUpdate.Airline,
		employeeToUpdate.Email,
		employeeToUpdate.IdentificationNumber,
		employeeToUpdate.Bp,
		employeeToUpdate.StartDate,
		employeeToUpdate.EndDate,
		employeeToUpdate.Active,
		employeeToUpdate.Role,
		employeeToUpdate.KeycloakUserID,
		employeeToUpdate.ID,
	)

	if err != nil {
		// Check for specific MySQL errors
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			case 1452:
				// Foreign key constraint fails (e.g., invalid airline)
				log.Error(logger.LogEmployeeUpdateError,
					"employee_id", employee.ID,
					"error", "invalid foreign key reference",
					"mysql_error", mysqlErr.Message)
				return domain.ErrInvalidForeignKey
			case 1406:
				// Data too long for column
				log.Error(logger.LogEmployeeUpdateError,
					"employee_id", employee.ID,
					"error", "data too long",
					"mysql_error", mysqlErr.Message)
				return domain.ErrDataTooLong
			case 1062:
				// Duplicate entry
				log.Error(logger.LogEmployeeUpdateError,
					"employee_id", employee.ID,
					"error", "duplicate entry",
					"mysql_error", mysqlErr.Message)
				return domain.ErrDuplicateUser
			}
		}
		// Generic error
		log.Error(logger.LogEmployeeUpdateError, "employee_id", employee.ID, "error", err)
		return domain.ErrUserCannotUpdate
	}

	rowsAffected, _ := result.RowsAffected()
	log.Debug("UpdateEmployee: Query executed successfully", "rows_affected", rowsAffected)

	return nil
}
