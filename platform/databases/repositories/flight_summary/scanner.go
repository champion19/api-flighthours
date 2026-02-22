package flight_summary

import (
	"database/sql"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

// scanDetail scans a full daily_logbook_detail row with all JOINed fields.
// This matches the same column order as the shared selectDetailColumns + detailJoins.
func scanDetail(rows interface {
	Scan(dest ...interface{}) error
}) (*domain.DailyLogbookDetail, error) {
	var (
		id                  string
		dailyLogbookID      string
		flightRealDate      string
		flightNumber        string
		airlineRouteID      string
		tailNumberID      string
		passengers          sql.NullInt64
		outTime             sql.NullString
		takeoffTime         sql.NullString
		landingTime         sql.NullString
		inTime              sql.NullString
		pilotRole           sql.NullString
		companionName       sql.NullString
		crewRole            sql.NullString
		airTime             sql.NullString
		blockTime           sql.NullString
		approachType        sql.NullString
		flightType          sql.NullString
		employeeLogbookID   sql.NullString
		logDate             sql.NullString
		tailNumber        sql.NullString
		modelName           sql.NullString
		routeCode           sql.NullString
		originIataCode      sql.NullString
		destinationIataCode sql.NullString
		airlineCode         sql.NullString
	)

	err := rows.Scan(
		&id, &dailyLogbookID, &flightRealDate, &flightNumber,
		&airlineRouteID, &tailNumberID, &passengers,
		&outTime, &takeoffTime, &landingTime, &inTime,
		&pilotRole, &companionName, &crewRole,
		&airTime, &blockTime,
		&approachType, &flightType, &employeeLogbookID,
		&logDate, &tailNumber, &modelName,
		&routeCode, &originIataCode, &destinationIataCode, &airlineCode,
	)
	if err != nil {
		return nil, err
	}

	detail := &domain.DailyLogbookDetail{
		ID:             id,
		DailyLogbookID: dailyLogbookID,
		FlightRealDate: flightRealDate,
		FlightNumber:   flightNumber,
		AirlineRouteID: airlineRouteID,
		TailNumberID: tailNumberID,
	}

	// Map nullable fields
	if passengers.Valid {
		p := int(passengers.Int64)
		detail.Passengers = &p
	}
	if outTime.Valid {
		detail.OutTime = &outTime.String
	}
	if takeoffTime.Valid {
		detail.TakeoffTime = &takeoffTime.String
	}
	if landingTime.Valid {
		detail.LandingTime = &landingTime.String
	}
	if inTime.Valid {
		detail.InTime = &inTime.String
	}
	if pilotRole.Valid {
		pr := domain.PilotRole(pilotRole.String)
		detail.PilotRole = &pr
	}
	if companionName.Valid {
		detail.CompanionName = &companionName.String
	}
	if crewRole.Valid {
		cr := domain.CrewRole(crewRole.String)
		detail.CrewRole = &cr
	}
	if airTime.Valid {
		detail.AirTime = &airTime.String
	}
	if blockTime.Valid {
		detail.BlockTime = &blockTime.String
	}
	if approachType.Valid {
		at := domain.ApproachType(approachType.String)
		detail.ApproachType = &at
	}
	if flightType.Valid {
		detail.FlightType = &flightType.String
	}
	if employeeLogbookID.Valid {
		detail.EmployeeLogbookID = &employeeLogbookID.String
	}
	if logDate.Valid {
		detail.LogDate = logDate.String
	}
	if tailNumber.Valid {
		detail.TailNumber = tailNumber.String
	}
	if modelName.Valid {
		detail.ModelName = modelName.String
	}
	if routeCode.Valid {
		detail.RouteCode = routeCode.String
	}
	if originIataCode.Valid {
		detail.OriginIataCode = originIataCode.String
	}
	if destinationIataCode.Valid {
		detail.DestinationIataCode = destinationIataCode.String
	}
	if airlineCode.Valid {
		detail.AirlineCode = airlineCode.String
	}

	return detail, nil
}
