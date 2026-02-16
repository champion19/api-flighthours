package aircraftmodel

import (
	"database/sql"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
)

// scanner is a minimal interface matching both *sql.Row and *sql.Rows.
type scanner interface {
	Scan(dest ...interface{}) error
}

// scanAircraftModel scans a single row into an AircraftModel, handling NullString fields.
// This helper eliminates the duplicated scan+NullString logic across list.go, get_by_family.go, and get_by_id.go.
func scanAircraftModel(s scanner) (*AircraftModel, error) {
	var model AircraftModel
	var engineTypeName sql.NullString
	var manufacturer sql.NullString

	if err := s.Scan(
		&model.ID,
		&model.ModelName,
		&model.AircraftTypeName,
		&engineTypeName,
		&model.Family,
		&manufacturer,
		&model.Status,
	); err != nil {
		return nil, err
	}

	if engineTypeName.Valid {
		model.EngineTypeName = engineTypeName.String
	}
	if manufacturer.Valid {
		model.Manufacturer = manufacturer.String
	}

	return &model, nil
}

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

// ToFamilyDomain converts the database entity to AircraftFamily domain model
func (a *AircraftModel) ToFamilyDomain() *domain.AircraftFamily {
	return &domain.AircraftFamily{
		Family:       a.Family,
		Manufacturer: a.Manufacturer,
	}
}
