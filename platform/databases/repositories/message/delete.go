package message

import(
	"context"
"github.com/champion19/flighthours-api/platform/databases/common"
"github.com/champion19/flighthours-api/core/interactor/services/domain"
"github.com/champion19/flighthours-api/core/ports/output"
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
