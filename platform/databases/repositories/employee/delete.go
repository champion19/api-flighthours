package employee

import (
	"context"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
)

func (r *repository) DeleteEmployee(ctx context.Context, tx output.Tx, id string) error {

	dbtx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = dbtx.ExecContext(ctx, QueryDelete, id)
	if err != nil {
		tx.Rollback()
		return domain.ErrUserCannotSave
	}

	if err = dbtx.Commit(); err != nil {
		return err
	}
	return nil
}
