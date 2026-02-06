package manufacturer

import(
domain	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

type Manufacturer struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

func (m *Manufacturer) ToDomain() *domain.Manufacturer {
	return &domain.Manufacturer{
		ID:   m.ID,
		Name: m.Name,
	}
}


func FromDomain(domainManufacturer *domain.Manufacturer) *Manufacturer {
	return &Manufacturer{
		ID:   domainManufacturer.ID,
		Name: domainManufacturer.Name,
	}
}
