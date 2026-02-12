package daily_logbook_detail

import (
	"context"

	"github.com/champion19/api-flighthours/platform/logger"
)

// ExistsByUniqueKey checks whether a flight detail already exists with the same
// (employee_logbook_id, flight_real_date, flight_number, actual_license_plate_id).
func (r *repository) ExistsByUniqueKey(ctx context.Context, employeeLogbookID, flightRealDate, flightNumber, licensePlateID string) (bool, error) {
	var count int
	err := r.stmtExistsByUniqueKey.QueryRowContext(ctx, employeeLogbookID, flightRealDate, flightNumber, licensePlateID).Scan(&count)
	if err != nil {
		log.Error(logger.LogDailyLogbookDetailDuplicate, "error", err)
		return false, err
	}
	return count > 0, nil
}
