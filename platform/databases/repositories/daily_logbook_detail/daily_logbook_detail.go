package daily_logbook_detail

import (
	"database/sql"
	"strings"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

type DailyLogbookDetail struct {
	ID                   string
	DailyLogbookID       string
	FlightRealDate       string // DATE stored as string
	FlightNumber         string
	OriginAirportID      string
	DestinationAirportID string
	TailNumberID         string
	Passengers           sql.NullInt64
	OutTime              sql.NullString // TIME stored as string HH:MM:SS (nullable)
	TakeoffTime          sql.NullString // nullable
	LandingTime          sql.NullString // nullable
	InTime               sql.NullString // nullable
	PilotRole            sql.NullString // nullable
	CrewRole             sql.NullString // nullable
	CompanionName        sql.NullString
	AirTime              sql.NullString // TIME stored as string HH:MM:SS (nullable)
	BlockTime            sql.NullString // TIME stored as string HH:MM:SS (nullable)
	ApproachCategory     sql.NullString
	ApproachSubtype      sql.NullString
	Autoland             sql.NullBool
	FlightType           sql.NullString
	EmployeeLogbookID    sql.NullString
	LogDate              sql.NullString
	TailNumber           sql.NullString
	ModelName            sql.NullString
	RouteCode            sql.NullString
	OriginIataCode       sql.NullString
	DestinationIataCode  sql.NullString
	AirlineCode          sql.NullString
}

func (d *DailyLogbookDetail) ToDomain() *domain.DailyLogbookDetail {
	detail := &domain.DailyLogbookDetail{
		ID:                   d.ID,
		DailyLogbookID:       d.DailyLogbookID,
		FlightRealDate:       d.FlightRealDate,
		FlightNumber:         d.FlightNumber,
		OriginAirportID:      d.OriginAirportID,
		DestinationAirportID: d.DestinationAirportID,
		TailNumberID:         d.TailNumberID,
	}

	d.mapTimeFields(detail)
	d.mapRoleFields(detail)
	d.mapOptionalFields(detail)
	d.mapDenormalizedFields(detail)

	return detail
}

func (d *DailyLogbookDetail) mapTimeFields(detail *domain.DailyLogbookDetail) {
	if d.OutTime.Valid {
		detail.OutTime = &d.OutTime.String
	}
	if d.TakeoffTime.Valid {
		detail.TakeoffTime = &d.TakeoffTime.String
	}
	if d.LandingTime.Valid {
		detail.LandingTime = &d.LandingTime.String
	}
	if d.InTime.Valid {
		detail.InTime = &d.InTime.String
	}
	if d.AirTime.Valid {
		detail.AirTime = &d.AirTime.String
	}
	if d.BlockTime.Valid {
		detail.BlockTime = &d.BlockTime.String
	}
}

func (d *DailyLogbookDetail) mapRoleFields(detail *domain.DailyLogbookDetail) {
	if d.PilotRole.Valid {
		pilotRole := domain.PilotRole(d.PilotRole.String)
		detail.PilotRole = &pilotRole
	}
	if d.CrewRole.Valid {
		crewRole := domain.CrewRole(d.CrewRole.String)
		detail.CrewRole = &crewRole
	}
}

func (d *DailyLogbookDetail) mapOptionalFields(detail *domain.DailyLogbookDetail) {
	if d.Passengers.Valid {
		passengers := int(d.Passengers.Int64)
		detail.Passengers = &passengers
	}
	if d.CompanionName.Valid {
		detail.CompanionName = &d.CompanionName.String
	}
	if d.ApproachCategory.Valid {
		approachCategory := domain.ApproachCategory(d.ApproachCategory.String)
		detail.ApproachCategory = &approachCategory
	}
	if d.ApproachSubtype.Valid {
		detail.ApproachSubtype = &d.ApproachSubtype.String
	}
	if d.Autoland.Valid {
		detail.Autoland = &d.Autoland.Bool
	}
	if d.FlightType.Valid {
		detail.FlightType = &d.FlightType.String
	}
	if d.EmployeeLogbookID.Valid {
		detail.EmployeeLogbookID = &d.EmployeeLogbookID.String
	}
}

func (d *DailyLogbookDetail) mapDenormalizedFields(detail *domain.DailyLogbookDetail) {
	if d.LogDate.Valid {
		detail.LogDate = d.LogDate.String
	}
	if d.TailNumber.Valid {
		detail.TailNumber = d.TailNumber.String
	}
	if d.ModelName.Valid {
		detail.ModelName = d.ModelName.String
	}
	if d.RouteCode.Valid {
		detail.RouteCode = d.RouteCode.String
	}
	if d.OriginIataCode.Valid {
		detail.OriginIataCode = d.OriginIataCode.String
	}
	if d.DestinationIataCode.Valid {
		detail.DestinationIataCode = d.DestinationIataCode.String
	}
	if d.AirlineCode.Valid {
		detail.AirlineCode = d.AirlineCode.String
	}
}

func FromDomain(d *domain.DailyLogbookDetail) *DailyLogbookDetail {
	entity := &DailyLogbookDetail{
		ID:                   d.ID,
		DailyLogbookID:       d.DailyLogbookID,
		FlightRealDate:       d.FlightRealDate,
		FlightNumber:         d.FlightNumber,
		OriginAirportID:      d.OriginAirportID,
		DestinationAirportID: d.DestinationAirportID,
		TailNumberID:         d.TailNumberID,
	}

	if d.OutTime != nil {
		entity.OutTime = sql.NullString{String: *d.OutTime, Valid: true}
	}
	if d.TakeoffTime != nil {
		entity.TakeoffTime = sql.NullString{String: *d.TakeoffTime, Valid: true}
	}
	if d.LandingTime != nil {
		entity.LandingTime = sql.NullString{String: *d.LandingTime, Valid: true}
	}
	if d.InTime != nil {
		entity.InTime = sql.NullString{String: *d.InTime, Valid: true}
	}
	if d.PilotRole != nil {
		entity.PilotRole = sql.NullString{String: string(*d.PilotRole), Valid: true}
	}
	if d.CrewRole != nil {
		entity.CrewRole = sql.NullString{String: string(*d.CrewRole), Valid: true}
	}
	if d.AirTime != nil {
		entity.AirTime = sql.NullString{String: *d.AirTime, Valid: true}
	}
	if d.BlockTime != nil {
		entity.BlockTime = sql.NullString{String: *d.BlockTime, Valid: true}
	}

	if d.Passengers != nil {
		entity.Passengers = sql.NullInt64{Int64: int64(*d.Passengers), Valid: true}
	}

	if d.CompanionName != nil {
		entity.CompanionName = sql.NullString{String: *d.CompanionName, Valid: true}
	}

	if d.ApproachCategory != nil {
		entity.ApproachCategory = sql.NullString{String: string(*d.ApproachCategory), Valid: true}
	}

	if d.ApproachSubtype != nil {
		entity.ApproachSubtype = sql.NullString{String: *d.ApproachSubtype, Valid: true}
	}

	if d.Autoland != nil {
		entity.Autoland = sql.NullBool{Bool: *d.Autoland, Valid: true}
	}

	if d.FlightType != nil {
		entity.FlightType = sql.NullString{String: *d.FlightType, Valid: true}
	}

	if d.EmployeeLogbookID != nil {
		entity.EmployeeLogbookID = sql.NullString{String: *d.EmployeeLogbookID, Valid: true}
	}

	return entity
}

// scanDetail scans a single row from *sql.Rows into a DailyLogbookDetail entity.
// This helper eliminates the duplicated 28-field scan block in list_by_employee.go and list_by_logbook.go.
func scanDetail(rows interface {
	Scan(dest ...interface{}) error
}) (*DailyLogbookDetail, error) {
	var entity DailyLogbookDetail
	err := rows.Scan(
		&entity.ID,
		&entity.DailyLogbookID,
		&entity.FlightRealDate,
		&entity.FlightNumber,
		&entity.OriginAirportID,
		&entity.DestinationAirportID,
		&entity.TailNumberID,
		&entity.Passengers,
		&entity.OutTime,
		&entity.TakeoffTime,
		&entity.LandingTime,
		&entity.InTime,
		&entity.PilotRole,
		&entity.CompanionName,
		&entity.CrewRole,
		&entity.AirTime,
		&entity.BlockTime,
		&entity.ApproachCategory,
		&entity.ApproachSubtype,
		&entity.Autoland,
		&entity.FlightType,
		&entity.EmployeeLogbookID,
		&entity.LogDate,
		&entity.TailNumber,
		&entity.ModelName,
		&entity.RouteCode,
		&entity.OriginIataCode,
		&entity.DestinationIataCode,
		&entity.AirlineCode,
	)
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// handleDetailFKError inspects a foreign-key error string and returns the appropriate domain error.
// Shared by SaveDailyLogbookDetail and UpdateDailyLogbookDetail.
func handleDetailFKError(errStr string) error {
	if strings.Contains(errStr, "daily_logbook") {
		return domain.ErrFlightInvalidLogbook
	}
	if strings.Contains(errStr, "origin_airport") || strings.Contains(errStr, "destination_airport") {
		return domain.ErrFlightInvalidAirport
	}
	if strings.Contains(errStr, "aircraft_registration") || strings.Contains(errStr, "tail_number") {
		return domain.ErrFlightInvalidTailNumber
	}
	return nil
}
