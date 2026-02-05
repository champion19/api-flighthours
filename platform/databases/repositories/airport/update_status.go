package airport

import (
	"context"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/databases/common"
)

// UpdateAirportStatus updates the status of an airport (active/inactive)
func (r *repository) UpdateAirportStatus(ctx context.Context, tx output.Tx, id string, status bool) error {
	sqlTx := tx.(*common.SQLTX)

	result, err := sqlTx.ExecContext(ctx, QueryUpdateStatus, status, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrAirportNotFound
	}

	return nil
}
