package services

import (
	"context"
	"time"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/logger"
)

type DailyLogbookDetailService struct {
	repo output.DailyLogbookDetailRepository
}

func NewDailyLogbookDetailService(repo output.DailyLogbookDetailRepository) *DailyLogbookDetailService {
	return &DailyLogbookDetailService{
		repo: repo,
	}
}

func (s *DailyLogbookDetailService) BeginTx(ctx context.Context) (output.Tx, error) {
	return s.repo.BeginTx(ctx)
}

func (s *DailyLogbookDetailService) GetDailyLogbookDetailByID(ctx context.Context, id string) (*domain.DailyLogbookDetail, error) {
	log.Info(logger.LogDailyLogbookDetailGet, "id", id)
	return s.repo.GetDailyLogbookDetailByID(ctx, id)
}

func (s *DailyLogbookDetailService) ListDailyLogbookDetailsByLogbook(ctx context.Context, logbookID string) ([]domain.DailyLogbookDetail, error) {
	log.Info(logger.LogDailyLogbookDetailList, "logbook_id", logbookID)
	return s.repo.ListDailyLogbookDetailsByLogbook(ctx, logbookID)
}

func (s *DailyLogbookDetailService) ListDailyLogbookDetailsByEmployee(ctx context.Context, employeeID string) ([]domain.DailyLogbookDetail, error) {
	log.Info(logger.LogDailyLogbookDetailList, "employee_id", employeeID)
	return s.repo.ListDailyLogbookDetailsByEmployee(ctx, employeeID)
}

// CreateDailyLogbookDetailTx creates a daily logbook detail using an external transaction
func (s *DailyLogbookDetailService) CreateDailyLogbookDetailTx(ctx context.Context, tx output.Tx, detail domain.DailyLogbookDetail) error {
	return s.repo.SaveDailyLogbookDetail(ctx, tx, detail)
}

// UpdateDailyLogbookDetailTx updates a daily logbook detail using an external transaction
func (s *DailyLogbookDetailService) UpdateDailyLogbookDetailTx(ctx context.Context, tx output.Tx, detail domain.DailyLogbookDetail) error {
	return s.repo.UpdateDailyLogbookDetail(ctx, tx, detail)
}

// DeleteDailyLogbookDetailTx deletes a daily logbook detail using an external transaction
func (s *DailyLogbookDetailService) DeleteDailyLogbookDetailTx(ctx context.Context, tx output.Tx, id string) error {
	return s.repo.DeleteDailyLogbookDetail(ctx, tx, id)
}

func (s *DailyLogbookDetailService) ValidateTimeSequence(outTime, takeoffTime, landingTime, inTime string) error {
	out, err := parseFlightTime(outTime)
	if err != nil {
		log.Error(logger.LogDailyLogbookDetailCreateError, "error", "invalid out_time format", "value", outTime)
		return domain.ErrFlightInvalidTimeSequence
	}

	takeoff, err := parseFlightTime(takeoffTime)
	if err != nil {
		log.Error(logger.LogDailyLogbookDetailCreateError, "error", "invalid takeoff_time format", "value", takeoffTime)
		return domain.ErrFlightInvalidTimeSequence
	}

	landing, err := parseFlightTime(landingTime)
	if err != nil {
		log.Error(logger.LogDailyLogbookDetailCreateError, "error", "invalid landing_time format", "value", landingTime)
		return domain.ErrFlightInvalidTimeSequence
	}

	in, err := parseFlightTime(inTime)
	if err != nil {
		log.Error(logger.LogDailyLogbookDetailCreateError, "error", "invalid in_time format", "value", inTime)
		return domain.ErrFlightInvalidTimeSequence
	}

	// Adjust for midnight-crossing flights
	takeoff = adjustForMidnight(takeoff, out)
	landing = adjustForMidnight(landing, takeoff)
	in = adjustForMidnight(in, landing)

	// Validate sequence: out <= takeoff < landing <= in
	if !out.Before(takeoff) && !out.Equal(takeoff) {
		log.Warn(logger.LogDailyLogbookDetailCreateError, "error", "out_time must be before or equal to takeoff_time")
		return domain.ErrFlightInvalidTimeSequence
	}

	if !takeoff.Before(landing) {
		log.Warn(logger.LogDailyLogbookDetailCreateError, "error", "takeoff_time must be before landing_time")
		return domain.ErrFlightInvalidTimeSequence
	}

	if !landing.Before(in) && !landing.Equal(in) {
		log.Warn(logger.LogDailyLogbookDetailCreateError, "error", "landing_time must be before or equal to in_time")
		return domain.ErrFlightInvalidTimeSequence
	}

	return nil
}

// parseFlightTime parses a time string with flexible format (HH:MM or HH:MM:SS)
func parseFlightTime(timeStr string) (time.Time, error) {
	t, err := time.Parse("15:04:05", timeStr)
	if err == nil {
		return t, nil
	}
	return time.Parse("15:04", timeStr)
}

// adjustForMidnight handles midnight-crossing: if current is earlier than previous
// and the gap is large (>12h), add 24h to indicate next-day crossing.
func adjustForMidnight(current, previous time.Time) time.Time {
	halfDay := 12 * time.Hour
	if current.Before(previous) && previous.Sub(current) > halfDay {
		return current.Add(24 * time.Hour)
	}
	return current
}

// ExistsByUniqueKey delegates duplicate check to the repository
func (s *DailyLogbookDetailService) ExistsByUniqueKey(ctx context.Context, employeeLogbookID, flightRealDate, flightNumber, tailNumberID string) (bool, error) {
	return s.repo.ExistsByUniqueKey(ctx, employeeLogbookID, flightRealDate, flightNumber, tailNumberID)
}

// ListCrewByDetail returns the command crew + cabin crew assigned to a flight leg
func (s *DailyLogbookDetailService) ListCrewByDetail(ctx context.Context, detailID string) ([]domain.CrewAssignment, error) {
	return s.repo.ListCrewByDetail(ctx, detailID)
}

// ReplaceCrewForDetailTx replaces the crew assigned to a flight leg using an external transaction
func (s *DailyLogbookDetailService) ReplaceCrewForDetailTx(ctx context.Context, tx output.Tx, detailID string, assignments []domain.CrewAssignment) error {
	return s.repo.ReplaceCrewForDetail(ctx, tx, detailID, assignments)
}
