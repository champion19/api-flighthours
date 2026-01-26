package message

import(
	"context"
	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/output"
	"github.com/champion19/flighthours-api/platform/databases/common"
)


func(r *repository)SaveMessage(ctx context.Context,tx output.Tx, message domain.Message) error{
messageTOUpdate:=FromDomain(message)

dbTx,ok:=tx.(*common.SQLTX)
if !ok{
	return domain.ErrInvalidTransaction
}
_,err:=dbTx.ExecContext(ctx,queryMessageSave,
	messageTOUpdate.ID,
	messageTOUpdate.Code,
	messageTOUpdate.Type,
	messageTOUpdate.Category,
	messageTOUpdate.Module,
	messageTOUpdate.Title,
	messageTOUpdate.Content,
	messageTOUpdate.Active,
)
if err!=nil{
	return domain.ErrMessageCannotSave
}
return nil

}
