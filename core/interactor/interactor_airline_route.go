package interactor

import (
	"context"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/input"
	"github.com/champion19/api-flighthours/platform/logger"
)

// AirlineRouteInteractor orchestrates airline route use cases
type AirlineRouteInteractor struct {
	service input.AirlineRouteService
}

// NewAirlineRouteInteractor creates a new airline route interactor
func NewAirlineRouteInteractor(service input.AirlineRouteService) *AirlineRouteInteractor {
	return &AirlineRouteInteractor{
		service: service,
	}
}

var airlineRouteLog = logger.NewSlogLogger()

// GetAirlineRouteByID retrieves an airline route by ID
func (i *AirlineRouteInteractor) GetAirlineRouteByID(ctx context.Context, traceID, id string) (*domain.AirlineRoute, error) {
	log := airlineRouteLog.WithTraceID(traceID)

	log.Info(logger.LogAirlineRouteGet, "airline_route_id", id)

	airlineRoute, err := i.service.GetAirlineRouteByID(ctx, id)
	if err != nil {
		if err == domain.ErrAirlineRouteNotFound {
			log.Warn(logger.LogAirlineRouteNotFound, "airline_route_id", id)
		} else {
			log.Error(logger.LogAirlineRouteGetError, "airline_route_id", id, "error", err)
		}
		return nil, err
	}

	log.Info(logger.LogAirlineRouteGetOK, "airline_route_id", id, "route_code", airlineRoute.RouteCode, "airline_code", airlineRoute.AirlineCode)
	return airlineRoute, nil
}

// ListAirlineRoutes retrieves all airline routes with optional filters
func (i *AirlineRouteInteractor) ListAirlineRoutes(ctx context.Context, traceID string, filters map[string]interface{}) ([]domain.AirlineRoute, error) {
	log := airlineRouteLog.WithTraceID(traceID)

	log.Info(logger.LogAirlineRouteList, "filters", filters)

	airlineRoutes, err := i.service.ListAirlineRoutes(ctx, filters)
	if err != nil {
		log.Error(logger.LogAirlineRouteListError, "error", err)
		return nil, err
	}

	log.Info(logger.LogAirlineRouteListOK, "count", len(airlineRoutes))
	return airlineRoutes, nil
}

func (i *AirlineRouteInteractor) ListMyAirlineRoutes(ctx context.Context, traceID, airlineID string) ([]domain.AirlineRoute, error) {
	log := airlineRouteLog.WithTraceID(traceID)

	log.Info(logger.LogAirlineRouteList, "airline_id", airlineID, "endpoint", "/airline-routes/me")

	airlineRoutes, err := i.service.ListAirlineRoutesByAirlineID(ctx, airlineID)
	if err != nil {
		log.Error(logger.LogAirlineRouteListError, "airline_id", airlineID, "error", err)
		return nil, err
	}

	log.Info(logger.LogAirlineRouteListOK, "airline_id", airlineID, "count", len(airlineRoutes))
	return airlineRoutes, nil
}

// ActivateAirlineRoute activates an airline route (HU42)
func (i *AirlineRouteInteractor) ActivateAirlineRoute(ctx context.Context, traceID, id string) error {
	log := airlineRouteLog.WithTraceID(traceID)

	log.Info(logger.LogAirlineRouteActivate, "airline_route_id", id)

	err := i.service.ActivateAirlineRoute(ctx, id)
	if err != nil {
		if err == domain.ErrAirlineRouteNotFound {
			log.Warn(logger.LogAirlineRouteNotFound, "airline_route_id", id)
			return err
		}
		if err == domain.ErrAirlineRouteAlreadyActive {
			log.Info(logger.LogAirlineRouteActivateOK, "airline_route_id", id, "status", "already_active")
			return err
		}
		log.Error(logger.LogAirlineRouteActivateError, "airline_route_id", id, "error", err)
		return err
	}

	log.Info(logger.LogAirlineRouteActivateOK, "airline_route_id", id)
	return nil
}

// DeactivateAirlineRoute deactivates an airline route (HU41)
func (i *AirlineRouteInteractor) DeactivateAirlineRoute(ctx context.Context, traceID, id string) error {
	log := airlineRouteLog.WithTraceID(traceID)

	log.Info(logger.LogAirlineRouteDeactivate, "airline_route_id", id)

	err := i.service.DeactivateAirlineRoute(ctx, id)
	if err != nil {
		if err == domain.ErrAirlineRouteNotFound {
			log.Warn(logger.LogAirlineRouteNotFound, "airline_route_id", id)
			return err
		}
		if err == domain.ErrAirlineRouteAlreadyInactive {
			log.Info(logger.LogAirlineRouteDeactivateOK, "airline_route_id", id, "status", "already_inactive")
			return err
		}
		log.Error(logger.LogAirlineRouteDeactivateError, "airline_route_id", id, "error", err)
		return err
	}

	log.Info(logger.LogAirlineRouteDeactivateOK, "airline_route_id", id)
	return nil
}
