package message

import(
	"context"
"github.com/champion19/api-flighthours/platform/databases/common"
"github.com/champion19/api-flighthours/core/interactor/services/domain"
"github.com/champion19/api-flighthours/core/ports/output"
)


func(r *repository)DeleteMessage(ctx context.Context,tx output.Tx, id string) error {
	dbTx,ok:=tx.(*common.SQLTX)
	if !ok{
		return domain.ErrInvalidTransaction
	}
	_,err:=dbTx.ExecContext(ctx,queryMessageDelete,id)
	if err!=nil{
		return domain.ErrMessageCannotDelete
	}
	return nil




}
