package employee

import (
	"context"


	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/databases/common"
	"github.com/go-sql-driver/mysql"
)

func (r *repository) Save(ctx context.Context, tx output.Tx, employee domain.Employee) error {

	employeeToSave := FromDomain(employee)

	dbTx, ok := tx.(*common.SQLTX)
	if !ok {
		return domain.ErrInvalidTransaction
	}

	_, err := dbTx.ExecContext(ctx, QuerySave,
		employeeToSave.ID,
		employeeToSave.Name,
		employeeToSave.Airline,
		employeeToSave.Email,
		employeeToSave.IdentificationNumber,
		employeeToSave.Bp,
		employeeToSave.StartDate,
		employeeToSave.EndDate,
		employeeToSave.Active,
		employeeToSave.Role,
		employeeToSave.KeycloakUserID)

	if err != nil {


		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return domain.ErrDuplicateUser
		}else{
			return domain.ErrUserCannotSave
		}
	}
	return nil
}
