package engine

import domain "github.com/champion19/api-flighthours/core/interactor/services/domain"

// Engine is the database entity for engine table
type Engine struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

// ToDomain converts the database entity to domain model
func (e *Engine) ToDomain() *domain.Engine {
	return &domain.Engine{
		ID:   e.ID,
		Name: e.Name,
	}
}

// FromDomain converts a domain model to database entity
func FromDomain(domainEngine *domain.Engine) *Engine {
	return &Engine{
		ID:   domainEngine.ID,
		Name: domainEngine.Name,
	}
}
