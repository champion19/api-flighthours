package tailnumber

import(
 domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
)



type TailNumber struct {
	ID              string `db:"id"`
	TailNumber    string `db:"tail_number"`
	AircraftModelID string `db:"aircraft_model_id"`
	AirlineID       string `db:"airline_id"`

	ModelName   string `db:"model_name"`
	AirlineName string `db:"airline_name"`
}


func (ar *TailNumber) ToDomain() *domain.TailNumber {
	return &domain.TailNumber{
		ID:              ar.ID,
		TailNumber:    ar.TailNumber,
		AircraftModelID: ar.AircraftModelID,
		AirlineID:       ar.AirlineID,
		ModelName:       ar.ModelName,
		AirlineName:     ar.AirlineName,
	}
}


func FromDomain(domainReg *domain.TailNumber) *TailNumber {
	return &TailNumber{
		ID:              domainReg.ID,
		TailNumber:    domainReg.TailNumber,
		AircraftModelID: domainReg.AircraftModelID,
		AirlineID:       domainReg.AirlineID,
		ModelName:       domainReg.ModelName,
		AirlineName:     domainReg.AirlineName,
	}
}
