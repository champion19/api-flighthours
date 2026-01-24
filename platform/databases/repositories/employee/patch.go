package employee

import (
	"context"


	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/output"
	"github.com/champion19/flighthours-api/platform/databases/common"
)


func (r *repository) PatchEmployee(ctx context.Context, tx output.Tx, id string, keycloakUserID string) error {

	dbTx, ok := tx.(*common.SQLTX)
	if !ok {
		return domain.ErrInvalidTransaction
	}

	result, err := dbTx.ExecContext(ctx, QueryPatch, keycloakUserID, id)
	if err != nil {
		return domain.ErrUserCannotSave
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrPersonNotFound
	}

	return nil
}
