package interactor

import (
	"context"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/input"
	"github.com/champion19/api-flighthours/middleware"
	"github.com/champion19/api-flighthours/platform/logger"
)

// RouteInteractor orchestrates route operations
type RouteInteractor struct {
	service input.RouteService
}

// NewRouteInteractor creates a new route interactor
func NewRouteInteractor(service input.RouteService) *RouteInteractor {
	return &RouteInteractor{
		service: service,
	}
}

// GetRouteByID retrieves a route by its ID
func (i *RouteInteractor) GetRouteByID(ctx context.Context, id string) (*domain.Route, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogRouteGet, "route_id", id)

	route, err := i.service.GetRouteByID(ctx, id)
	if err != nil {
		log.Error(logger.LogRouteGetError, "route_id", id, "error", err)
		return nil, err
	}

	log.Success(logger.LogRouteGetOK, route.ToLogger())
	return route, nil
}

// ListRoutes retrieves all routes with optional filters
func (i *RouteInteractor) ListRoutes(ctx context.Context, filters map[string]interface{}) ([]domain.Route, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogRouteList, "filters", filters)

	routes, err := i.service.ListRoutes(ctx, filters)
	if err != nil {
		log.Error(logger.LogRouteListError, "error", err)
		return nil, err
	}

	log.Success(logger.LogRouteListOK, "count", len(routes))
	return routes, nil
}
