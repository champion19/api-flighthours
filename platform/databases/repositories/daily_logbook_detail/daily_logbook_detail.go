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
	OutTime             string // TIME stored as string HH:MM:SS
	TakeoffTime         string
	LandingTime         string
	InTime              string
	PilotRole           string
	CompanionName       sql.NullString
	AirTime             string         // TIME stored as string HH:MM:SS
	BlockTime           string         // TIME stored as string HH:MM:SS
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
		OutTime:        d.OutTime,
		TakeoffTime:    d.TakeoffTime,
		LandingTime:    d.LandingTime,
		InTime:         d.InTime,
		PilotRole:      domain.PilotRole(d.PilotRole),
		AirTime:        d.AirTime,
		BlockTime:      d.BlockTime,
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
		OutTime:        d.OutTime,
		TakeoffTime:    d.TakeoffTime,
		LandingTime:    d.LandingTime,
		InTime:         d.InTime,
		PilotRole:      string(d.PilotRole),
		AirTime:        d.AirTime,
		BlockTime:      d.BlockTime,
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
