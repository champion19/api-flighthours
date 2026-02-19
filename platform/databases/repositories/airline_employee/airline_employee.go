package airline_employee

import (
	"time"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/platform/logger"
	"github.com/go-sql-driver/mysql"
)

// castTx is now replaced by common.CastTx — see tx_helper.go

// handleMySQLError inspects a MySQL error and returns the appropriate domain error.
// It handles foreign key violations (1452) and duplicate entries (1062).
func handleMySQLError(err error, employeeID string) error {
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		switch mysqlErr.Number {
		case 1452:
			log.Error(logger.LogDatabaseUnavailable,
				"employee_id", employeeID,
				"error", "invalid foreign key reference",
				"mysql_error", mysqlErr.Message)
			return domain.ErrInvalidForeignKey
		case 1062:
			log.Error(logger.LogDatabaseUnavailable,
				"employee_id", employeeID,
				"error", "duplicate entry",
				"mysql_error", mysqlErr.Message)
			return domain.ErrDuplicateUser
		}
	}
	log.Error(logger.LogDatabaseUnavailable, "employee_id", employeeID, "error", err)
	return err
}

// AirlineEmployee repository struct - airline-specific fields only
// This maps to the airline-related columns in the employee table
type AirlineEmployee struct {
	ID        string    `db:"id"`
	AirlineID *string   `db:"airline"`
	Bp        *string   `db:"bp"`
	StartDate time.Time `db:"start_date"`
	EndDate   time.Time `db:"end_date"`
	Active    bool      `db:"active"`
}

func (e *AirlineEmployee) ToDomain() *domain.AirlineEmployee {
	airlineID := ""
	if e.AirlineID != nil {
		airlineID = *e.AirlineID
	}
	bp := ""
	if e.Bp != nil {
		bp = *e.Bp
	}

	return &domain.AirlineEmployee{
		ID:        e.ID,
		AirlineID: airlineID,
		Bp:        bp,
		StartDate: e.StartDate,
		EndDate:   e.EndDate,
		Active:    e.Active,
	}
}

func stringPtrOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func FromDomain(domainEmployee *domain.AirlineEmployee) *AirlineEmployee {
	return &AirlineEmployee{
		ID:        domainEmployee.ID,
		AirlineID: stringPtrOrNil(domainEmployee.AirlineID),
		Bp:        stringPtrOrNil(domainEmployee.Bp),
		StartDate: domainEmployee.StartDate,
		EndDate:   domainEmployee.EndDate,
		Active:    domainEmployee.Active,
	}
}
