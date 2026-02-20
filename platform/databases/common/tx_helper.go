package common

import (
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
)

// CastTx casts an output.Tx to the concrete *SQLTX type.
// This eliminates the duplicated type assertion across all repository write operations.
func CastTx(tx output.Tx) (*SQLTX, error) {
	sqlTx, ok := tx.(*SQLTX)
	if !ok {
		return nil, domain.ErrInvalidTransaction
	}
	return sqlTx, nil
}
