package message

import (
	"context"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/databases/common"
)

func (r *repository) DeleteMessage(ctx context.Context, tx output.Tx, id string) error {
	dbTx, err := common.CastTx(tx)
	if err != nil {
		return err
	}
	_, err = dbTx.ExecContext(ctx, queryMessageDelete, id)
	if err != nil {
		return domain.ErrMessageCannotDelete
	}
	return nil

}
