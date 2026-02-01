package airline_route

import (
	"context"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
)

func (r *repository) ListAirlineRoutesByAirlineID(ctx context.Context, airlineID string) ([]domain.AirlineRoute, error) {
	filters := map[string]interface{}{
		"airline_id": airlineID,
	}
	return r.ListAirlineRoutes(ctx, filters)
}
