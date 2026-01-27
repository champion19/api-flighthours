package employee

import (
	"context"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/databases/common"
	"github.com/champion19/api-flighthours/platform/logger"
	"github.com/go-sql-driver/mysql"
)

func (r *repository) UpdateEmployee(ctx context.Context, tx output.Tx, employee domain.Employee) error {
	employeeToUpdate := FromDomain(employee)

	log.Debug(logger.LogRepoEmployeeUpdateStart,
		"employee_id", employeeToUpdate.ID,
		"name", employeeToUpdate.Name,
		"email", employeeToUpdate.Email,
		"active", employeeToUpdate.Active)

	dbTx, ok := tx.(*common.SQLTX)
	if !ok {
		log.Error(logger.LogRepoEmployeeUpdateTxErr, "error", "invalid transaction type")
		return domain.ErrInvalidTransaction
	}

	log.Debug(logger.LogRepoEmployeeUpdateTxOK)

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
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			case 1452:
				log.Error(logger.LogEmployeeUpdateError,
					"employee_id", employee.ID,
					"error", "invalid foreign key reference",
					"mysql_error", mysqlErr.Message)
				return domain.ErrInvalidForeignKey
			case 1406:
				log.Error(logger.LogEmployeeUpdateError,
					"employee_id", employee.ID,
					"error", "data too long",
					"mysql_error", mysqlErr.Message)
				return domain.ErrDataTooLong
			case 1062:
				log.Error(logger.LogEmployeeUpdateError,
					"employee_id", employee.ID,
					"error", "duplicate entry",
					"mysql_error", mysqlErr.Message)
				return domain.ErrDuplicateUser
			}
		}
		log.Error(logger.LogEmployeeUpdateError, "employee_id", employee.ID, "error", err)
		return domain.ErrUserCannotUpdate
	}

	rowsAffected, _ := result.RowsAffected()
	log.Debug(logger.LogRepoEmployeeUpdateQueryOK, "rows_affected", rowsAffected)

	return nil
}
