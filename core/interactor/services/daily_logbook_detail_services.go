package services

import(
	"context"
	"time"

	"github.com/champion19/api-flighthours/platform/logger"
	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
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


func (s *DailyLogbookDetailService) CreateDailyLogbookDetail(ctx context.Context, detail domain.DailyLogbookDetail) error {
	log.Info(logger.LogDailyLogbookDetailCreate, "data", detail.ToLogger())

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogDBTransactionBeginErr, "error", err)
		return err
	}

	err = s.repo.SaveDailyLogbookDetail(ctx, tx, detail)
	if err != nil {
		tx.Rollback()
		log.Error(logger.LogDailyLogbookDetailCreateError, "error", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Error(logger.LogDBTransactionCommitErr, "error", err)
		return err
	}

	log.Info(logger.LogDailyLogbookDetailCreateOK, "id", detail.ID)
	return nil
}


func (s *DailyLogbookDetailService) UpdateDailyLogbookDetail(ctx context.Context, detail domain.DailyLogbookDetail) error {
	log.Info(logger.LogDailyLogbookDetailUpdate, "data", detail.ToLogger())

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogDBTransactionBeginErr, "error", err)
		return err
	}

	err = s.repo.UpdateDailyLogbookDetail(ctx, tx, detail)
	if err != nil {
		tx.Rollback()
		log.Error(logger.LogDailyLogbookDetailUpdateError, "error", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Error(logger.LogDBTransactionCommitErr, "error", err)
		return err
	}

	log.Info(logger.LogDailyLogbookDetailUpdateOK, "id", detail.ID)
	return nil
}


func (s *DailyLogbookDetailService) DeleteDailyLogbookDetail(ctx context.Context, id string) error {
	log.Info(logger.LogDailyLogbookDetailDelete, "id", id)

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogDBTransactionBeginErr, "error", err)
		return err
	}


	if err := tx.Commit(); err != nil {
		log.Error(logger.LogDBTransactionCommitErr, "error", err)
		return err
	}

	log.Info(logger.LogDailyLogbookDetailDeleteOK, "id", id)
	return nil
}


func (s *DailyLogbookDetailService) ValidateTimeSequence(outTime, takeoffTime, landingTime, inTime string) error {
	// Parse time with flexible format (HH:MM or HH:MM:SS)
	parseTime := func(timeStr string) (time.Time, error) {
		// Try HH:MM:SS first
		t, err := time.Parse("15:04:05", timeStr)
		if err == nil {
			return t, nil
		}
		// Fallback to HH:MM
		return time.Parse("15:04", timeStr)
	}

	out, err := parseTime(outTime)
	if err != nil {
		log.Error(logger.LogDailyLogbookDetailCreateError, "error", "invalid out_time format", "value", outTime)
		return domain.ErrFlightInvalidTimeSequence
	}

	takeoff, err := parseTime(takeoffTime)
	if err != nil {
		log.Error(logger.LogDailyLogbookDetailCreateError, "error", "invalid takeoff_time format", "value", takeoffTime)
		return domain.ErrFlightInvalidTimeSequence
	}

	landing, err := parseTime(landingTime)
	if err != nil {
		log.Error(logger.LogDailyLogbookDetailCreateError, "error", "invalid landing_time format", "value", landingTime)
		return domain.ErrFlightInvalidTimeSequence
	}

	in, err := parseTime(inTime)
	if err != nil {
		log.Error(logger.LogDailyLogbookDetailCreateError, "error", "invalid in_time format", "value", inTime)
		return domain.ErrFlightInvalidTimeSequence
	}

	// Validate sequence: out < takeoff < landing < in
	if !out.Before(takeoff) {
		log.Warn(logger.LogDailyLogbookDetailCreateError, "error", "out_time must be before takeoff_time")
		return domain.ErrFlightInvalidTimeSequence
	}

	if !takeoff.Before(landing) {
		log.Warn(logger.LogDailyLogbookDetailCreateError, "error", "takeoff_time must be before landing_time")
		return domain.ErrFlightInvalidTimeSequence
	}

	if !landing.Before(in) {
		log.Warn(logger.LogDailyLogbookDetailCreateError, "error", "landing_time must be before in_time")
		return domain.ErrFlightInvalidTimeSequence
	}

	return nil
}
