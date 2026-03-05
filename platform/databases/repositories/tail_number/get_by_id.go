package tailnumber

import(
	"context"
	"database/sql"
	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
)


func (r *repository) GetTailNumberByID(ctx context.Context, id string) (*domain.TailNumber, error) {
	var ar TailNumber
	err := r.stmtGetByID.QueryRowContext(ctx, id).Scan(
		&ar.ID,
		&ar.TailNumber,
		&ar.AircraftModelID,
		&ar.AirlineID,
		&ar.ModelName,
		&ar.AirlineName,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrTailNumberNotFound
		}
		return nil, err
	}
	return ar.ToDomain(), nil
}
