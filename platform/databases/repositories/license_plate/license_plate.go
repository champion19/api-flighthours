package licenseplate

import(
 domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
)



type LicensePlate struct {
	ID              string `db:"id"`
	LicensePlate    string `db:"license_plate"`
	AircraftModelID string `db:"aircraft_model_id"`
	AirlineID       string `db:"airline_id"`

	ModelName   string `db:"model_name"`
	AirlineName string `db:"airline_name"`
}


func (ar *LicensePlate) ToDomain() *domain.LicensePlate {
	return &domain.LicensePlate{
		ID:              ar.ID,
		LicensePlate:    ar.LicensePlate,
		AircraftModelID: ar.AircraftModelID,
		AirlineID:       ar.AirlineID,
		ModelName:       ar.ModelName,
		AirlineName:     ar.AirlineName,
	}
}


func FromDomain(domainReg *domain.LicensePlate) *LicensePlate {
	return &LicensePlate{
		ID:              domainReg.ID,
		LicensePlate:    domainReg.LicensePlate,
		AircraftModelID: domainReg.AircraftModelID,
		AirlineID:       domainReg.AirlineID,
		ModelName:       domainReg.ModelName,
		AirlineName:     domainReg.AirlineName,
	}
}
