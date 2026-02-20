package message

import (
	"context"
	"fmt"
	"strings"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/databases/common"
)

func (r *repository) UpdateMessage(ctx context.Context, tx output.Tx, message domain.Message) error {
	dbTx, err := common.CastTx(tx)
	if err != nil {
		return err
	}

	var setClauses []string
	var args []interface{}

	if message.Title != "" {
		setClauses = append(setClauses, "message_title = ?")
		args = append(args, message.Title)
	}
	if message.Content != "" {
		setClauses = append(setClauses, "message_content = ?")
		args = append(args, message.Content)
	}
	setClauses = append(setClauses, "is_active = ?")
	args = append(args, message.Active)

	if len(setClauses) == 0 {
		return domain.ErrMessageCannotUpdate
	}

	args = append(args, message.ID, message.Code)

	query := fmt.Sprintf("UPDATE system_messages SET %s WHERE id = ? and message_code = ?", strings.Join(setClauses, ", "))

	result, err := dbTx.ExecContext(ctx, query, args...)
	if err != nil {
		return domain.ErrMessageCannotUpdate
	}

	_ = result

	return nil
}
