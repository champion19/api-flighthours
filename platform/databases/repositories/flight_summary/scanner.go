package flight_summary

import (
	"database/sql"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

// nullableString returns the string value if valid, or empty string.
func nullableString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

// nullableStringPtr returns a pointer to the string if valid, or nil.
func nullableStringPtr(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

// nullableIntPtr returns a pointer to int if valid, or nil.
func nullableIntPtr(ni sql.NullInt64) *int {
	if ni.Valid {
		v := int(ni.Int64)
		return &v
	}
	return nil
}

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
		tailNumberID        string
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
		approachCategory    sql.NullString
		approachSubtype     sql.NullString
		autoland            sql.NullBool
		flightType          sql.NullString
		employeeLogbookID   sql.NullString
		logDate             sql.NullString
		tailNumber          sql.NullString
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
		&approachCategory, &approachSubtype, &autoland, &flightType, &employeeLogbookID,
		&logDate, &tailNumber, &modelName,
		&routeCode, &originIataCode, &destinationIataCode, &airlineCode,
	)
	if err != nil {
		return nil, err
	}

	detail := &domain.DailyLogbookDetail{
		ID:                  id,
		DailyLogbookID:      dailyLogbookID,
		FlightRealDate:      flightRealDate,
		FlightNumber:        flightNumber,
		AirlineRouteID:      airlineRouteID,
		TailNumberID:        tailNumberID,
		Passengers:          nullableIntPtr(passengers),
		OutTime:             nullableStringPtr(outTime),
		TakeoffTime:         nullableStringPtr(takeoffTime),
		LandingTime:         nullableStringPtr(landingTime),
		InTime:              nullableStringPtr(inTime),
		CompanionName:       nullableStringPtr(companionName),
		AirTime:             nullableStringPtr(airTime),
		BlockTime:           nullableStringPtr(blockTime),
		FlightType:          nullableStringPtr(flightType),
		LogDate:             nullableString(logDate),
		TailNumber:          nullableString(tailNumber),
		ModelName:           nullableString(modelName),
		RouteCode:           nullableString(routeCode),
		OriginIataCode:      nullableString(originIataCode),
		DestinationIataCode: nullableString(destinationIataCode),
		AirlineCode:         nullableString(airlineCode),
	}

	// Map typed nullable fields (PilotRole, CrewRole, ApproachCategory)
	if pilotRole.Valid {
		pr := domain.PilotRole(pilotRole.String)
		detail.PilotRole = &pr
	}
	if crewRole.Valid {
		cr := domain.CrewRole(crewRole.String)
		detail.CrewRole = &cr
	}
	if approachCategory.Valid {
		ac := domain.ApproachCategory(approachCategory.String)
		detail.ApproachCategory = &ac
	}
	if approachSubtype.Valid {
		detail.ApproachSubtype = &approachSubtype.String
	}
	if autoland.Valid {
		detail.Autoland = &autoland.Bool
	}
	if employeeLogbookID.Valid {
		detail.EmployeeLogbookID = &employeeLogbookID.String
	}

	return detail, nil
}
