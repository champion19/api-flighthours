package manufacturer

import(
domain	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

// Manufacturer is the database entity for manufacturer table
type Manufacturer struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

// ToDomain converts the database entity to domain model
func (m *Manufacturer) ToDomain() *domain.Manufacturer {
	return &domain.Manufacturer{
		ID:   m.ID,
		Name: m.Name,
	}
}

// FromDomain converts a domain model to database entity
func FromDomain(domainManufacturer *domain.Manufacturer) *Manufacturer {
	return &Manufacturer{
		ID:   domainManufacturer.ID,
		Name: domainManufacturer.Name,
	}
}
