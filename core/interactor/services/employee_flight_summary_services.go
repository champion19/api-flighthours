package services

import (
	"context"
	"fmt"
	"time"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/logger"
)

type EmployeeFlightSummaryServiceImpl struct {
	repo output.EmployeeFlightSummaryRepository
}

func NewEmployeeFlightSummaryService(repo output.EmployeeFlightSummaryRepository) *EmployeeFlightSummaryServiceImpl {
	return &EmployeeFlightSummaryServiceImpl{repo: repo}
}

// AccumulateFlightHours updates all affected period rows when a flight detail is created or deleted.
// This runs INSIDE the same transaction as the detail create/update/delete.
func (s *EmployeeFlightSummaryServiceImpl) AccumulateFlightHours(ctx context.Context, tx output.Tx, employeeID string, detail domain.DailyLogbookDetail, isDeletion bool) error {
	log.Info(logger.LogFlightSummaryGet, "action", "accumulate",
		"employee_id", employeeID,
		"flight_date", detail.FlightRealDate,
		"is_deletion", isDeletion)

	// Parse flight date
	flightDate, err := time.Parse("2006-01-02", detail.FlightRealDate)
	if err != nil {
		log.Error(logger.LogFlightSummaryGetError, "action", "accumulate", "error", fmt.Sprintf("invalid flight date: %s", detail.FlightRealDate))
		return err
	}

	// Calculate deltas
	airTimeMinutes := parseTimeToMinutes(detail.AirTime)
	blockTimeMinutes := parseTimeToMinutes(detail.BlockTime)
	flightsDelta := 1
	landingsDelta := 0
	if domain.IsLandingRole(detail.PilotRole) {
		landingsDelta = 1
	}

	// If deletion, negate all deltas
	if isDeletion {
		airTimeMinutes = -airTimeMinutes
		blockTimeMinutes = -blockTimeMinutes
		flightsDelta = -1
		landingsDelta = -landingsDelta
	}

	// Get all affected periods for this flight date
	periods := domain.GetAffectedPeriods(flightDate)

	// Upsert each period row within the same transaction
	for _, period := range periods {
		if err := s.repo.UpsertSummary(ctx, tx, employeeID, period,
			airTimeMinutes, blockTimeMinutes, flightsDelta, landingsDelta); err != nil {
			log.Error(logger.LogFlightSummaryGetError, "action", "accumulate",
				"period_type", period.PeriodType, "error", err)
			return err
		}
	}

	log.Info(logger.LogFlightSummaryGetOK, "action", "accumulate",
		"employee_id", employeeID,
		"periods_updated", len(periods),
		"air_time_delta", airTimeMinutes)
	return nil
}

// GetSummariesByEmployee delegates to the repository
func (s *EmployeeFlightSummaryServiceImpl) GetSummariesByEmployee(ctx context.Context, employeeID, periodType string) ([]domain.EmployeeFlightSummary, error) {
	return s.repo.GetSummariesByEmployee(ctx, employeeID, periodType)
}

// GetCurrentPeriodSummary delegates to the repository
func (s *EmployeeFlightSummaryServiceImpl) GetCurrentPeriodSummary(ctx context.Context, employeeID, periodType string, year, number int) (*domain.EmployeeFlightSummary, error) {
	return s.repo.GetCurrentPeriodSummary(ctx, employeeID, periodType, year, number)
}

// GetAllSummaries delegates to the repository
func (s *EmployeeFlightSummaryServiceImpl) GetAllSummaries(ctx context.Context, employeeID string) ([]domain.EmployeeFlightSummary, error) {
	return s.repo.GetAllSummaries(ctx, employeeID)
}

// parseTimeToMinutes converts a TIME string (HH:MM:SS or HH:MM) to total minutes
func parseTimeToMinutes(timeStr *string) int {
	if timeStr == nil || *timeStr == "" {
		return 0
	}

	t, err := time.Parse("15:04:05", *timeStr)
	if err != nil {
		// Try HH:MM format
		t, err = time.Parse("15:04", *timeStr)
		if err != nil {
			return 0
		}
	}

	return t.Hour()*60 + t.Minute()
}
