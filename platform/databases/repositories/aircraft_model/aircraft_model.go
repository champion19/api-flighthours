package aircraftmodel

import(
	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
)

// AircraftModel is the database entity for aircraft_model table
type AircraftModel struct {
	ID               string `db:"id"`
	ModelName        string `db:"model_name"`
	AircraftTypeName string `db:"aircraft_type_name"`
	EngineTypeName   string `db:"engine_type_name"`
	Family           string `db:"family"`
	Manufacturer     string `db:"manufacturer"`
	Status           bool   `db:"status"`
}

// ToDomain converts the database entity to domain model
func (a *AircraftModel) ToDomain() *domain.AircraftModel {
	return &domain.AircraftModel{
		ID:               a.ID,
		ModelName:        a.ModelName,
		AircraftTypeName: a.AircraftTypeName,
		EngineTypeName:   a.EngineTypeName,
		Family:           a.Family,
		Manufacturer:     a.Manufacturer,
		Status:           a.Status,
	}
}

// FromDomain converts a domain model to database entity
func FromDomain(domainModel *domain.AircraftModel) *AircraftModel {
	return &AircraftModel{
		ID:               domainModel.ID,
		ModelName:        domainModel.ModelName,
		AircraftTypeName: domainModel.AircraftTypeName,
		EngineTypeName:   domainModel.EngineTypeName,
		Family:           domainModel.Family,
		Manufacturer:     domainModel.Manufacturer,
		Status:           domainModel.Status,
	}
}
