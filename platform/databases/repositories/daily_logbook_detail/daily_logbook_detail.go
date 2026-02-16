package daily_logbook_detail

import (
	"database/sql"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

type DailyLogbookDetail struct {
	ID                  string
	DailyLogbookID      string
	FlightRealDate      string // DATE stored as string
	FlightNumber        string
	AirlineRouteID      string
	LicensePlateID      string
	Passengers          sql.NullInt64
	OutTime             sql.NullString // TIME stored as string HH:MM:SS (nullable)
	TakeoffTime         sql.NullString // nullable
	LandingTime         sql.NullString // nullable
	InTime              sql.NullString // nullable
	PilotRole           sql.NullString // nullable
	CrewRole            sql.NullString // nullable
	CompanionName       sql.NullString
	AirTime             sql.NullString // TIME stored as string HH:MM:SS (nullable)
	BlockTime           sql.NullString // TIME stored as string HH:MM:SS (nullable)
	DutyTime            sql.NullString // TIME stored as string HH:MM:SS (nullable)
	ApproachType        sql.NullString
	FlightType          sql.NullString
	EmployeeLogbookID   sql.NullString
	LogDate             sql.NullString
	LicensePlate        sql.NullString
	ModelName           sql.NullString
	RouteCode           sql.NullString
	OriginIataCode      sql.NullString
	DestinationIataCode sql.NullString
	AirlineCode         sql.NullString
}

func (d *DailyLogbookDetail) ToDomain() *domain.DailyLogbookDetail {
	detail := &domain.DailyLogbookDetail{
		ID:             d.ID,
		DailyLogbookID: d.DailyLogbookID,
		FlightRealDate: d.FlightRealDate,
		FlightNumber:   d.FlightNumber,
		AirlineRouteID: d.AirlineRouteID,
		LicensePlateID: d.LicensePlateID,
	}

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
	if d.PilotRole.Valid {
		pilotRole := domain.PilotRole(d.PilotRole.String)
		detail.PilotRole = &pilotRole
	}
	if d.CrewRole.Valid {
		crewRole := domain.CrewRole(d.CrewRole.String)
		detail.CrewRole = &crewRole
	}
	if d.AirTime.Valid {
		detail.AirTime = &d.AirTime.String
	}
	if d.BlockTime.Valid {
		detail.BlockTime = &d.BlockTime.String
	}

	if d.Passengers.Valid {
		passengers := int(d.Passengers.Int64)
		detail.Passengers = &passengers
	}

	if d.CompanionName.Valid {
		detail.CompanionName = &d.CompanionName.String
	}

	if d.DutyTime.Valid {
		detail.DutyTime = &d.DutyTime.String
	}

	if d.ApproachType.Valid {
		approachType := domain.ApproachType(d.ApproachType.String)
		detail.ApproachType = &approachType
	}

	if d.FlightType.Valid {
		detail.FlightType = &d.FlightType.String
	}

	if d.EmployeeLogbookID.Valid {
		detail.EmployeeLogbookID = &d.EmployeeLogbookID.String
	}

	if d.LogDate.Valid {
		detail.LogDate = d.LogDate.String
	}
	if d.LicensePlate.Valid {
		detail.LicensePlate = d.LicensePlate.String
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

	return detail
}

func FromDomain(d *domain.DailyLogbookDetail) *DailyLogbookDetail {
	entity := &DailyLogbookDetail{
		ID:             d.ID,
		DailyLogbookID: d.DailyLogbookID,
		FlightRealDate: d.FlightRealDate,
		FlightNumber:   d.FlightNumber,
		AirlineRouteID: d.AirlineRouteID,
		LicensePlateID: d.LicensePlateID,
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

	if d.DutyTime != nil {
		entity.DutyTime = sql.NullString{String: *d.DutyTime, Valid: true}
	}

	if d.ApproachType != nil {
		entity.ApproachType = sql.NullString{String: string(*d.ApproachType), Valid: true}
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
		&entity.AirlineRouteID,
		&entity.LicensePlateID,
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
		&entity.DutyTime,
		&entity.ApproachType,
		&entity.FlightType,
		&entity.EmployeeLogbookID,
		&entity.LogDate,
		&entity.LicensePlate,
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
