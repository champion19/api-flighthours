package message

import (
	"context"
	"fmt"
	"strings"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/output"
	"github.com/champion19/flighthours-api/platform/databases/common"
)

func (r *repository) UpdateMessage(ctx context.Context, tx output.Tx, message domain.Message) error {
	dbTx, ok := tx.(*common.SQLTX)
	if !ok {
		return domain.ErrInvalidTransaction
	}

	// Build dynamic UPDATE query based on provided fields
	// Note: message_code is immutable and cannot be updated
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
	// Active is a boolean, so we always include it if the message has an ID
	setClauses = append(setClauses, "is_active = ?")
	args = append(args, message.Active)

	if len(setClauses) == 0 {
		return domain.ErrMessageCannotUpdate
	}

	// Add ID to args
	args = append(args, message.ID,message.Code)

	// Build final query
	query := fmt.Sprintf("UPDATE system_messages SET %s WHERE id = ? and message_code = ?", strings.Join(setClauses, ", "))

result, err := dbTx.ExecContext(ctx, query, args...)
	if err != nil {
		return domain.ErrMessageCannotUpdate
	}

	_ = result

	return nil
}
