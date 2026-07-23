package interactor

import (
	"context"

	"github.com/champion19/api-flighthours/core/interactor/helpers"
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/input"
	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/platform/logger"
	"github.com/champion19/api-flighthours/tools/utils"
)

// AirlineRouteInteractor orchestrates airline route use cases
type AirlineRouteInteractor struct {
	service      input.AirlineRouteService
	routeService input.RouteService
}

// NewAirlineRouteInteractor creates a new airline route interactor
func NewAirlineRouteInteractor(service input.AirlineRouteService, routeService input.RouteService) *AirlineRouteInteractor {
	return &AirlineRouteInteractor{
		service:      service,
		routeService: routeService,
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

	log.Info(logger.LogAirlineRouteList, "airline_id", airlineID, "endpoint", "/employees/airline-routes")

	airlineRoutes, err := i.service.ListAirlineRoutesByAirlineID(ctx, airlineID)
	if err != nil {
		log.Error(logger.LogAirlineRouteListError, "airline_id", airlineID, "error", err)
		return nil, err
	}

	log.Info(logger.LogAirlineRouteListOK, "airline_id", airlineID, "count", len(airlineRoutes))
	return airlineRoutes, nil
}

// ActivateAirlineRoute activates an airline route (HU42)
func (i *AirlineRouteInteractor) ActivateAirlineRoute(ctx context.Context, traceID, id string) (err error) {
	log := airlineRouteLog.WithTraceID(traceID)

	log.Info(logger.LogAirlineRouteActivate, "airline_route_id", id)

	err = helpers.RunWithTx(ctx, i.service, log, logger.LogAirlineRouteActivateError,
		func(ctx context.Context, tx output.Tx) error {
			if err := i.service.ActivateAirlineRouteTx(ctx, tx, id); err != nil {
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
			return nil
		})
	if err != nil {
		return err
	}

	log.Info(logger.LogAirlineRouteActivateOK, "airline_route_id", id)
	return nil
}

// DeactivateAirlineRoute deactivates an airline route (HU41)
func (i *AirlineRouteInteractor) DeactivateAirlineRoute(ctx context.Context, traceID, id string) (err error) {
	log := airlineRouteLog.WithTraceID(traceID)

	log.Info(logger.LogAirlineRouteDeactivate, "airline_route_id", id)

	err = helpers.RunWithTx(ctx, i.service, log, logger.LogAirlineRouteDeactivateError,
		func(ctx context.Context, tx output.Tx) error {
			if err := i.service.DeactivateAirlineRouteTx(ctx, tx, id); err != nil {
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
			return nil
		})
	if err != nil {
		return err
	}

	log.Info(logger.LogAirlineRouteDeactivateOK, "airline_route_id", id)
	return nil
}

// ResolveOrCreatePendingAirlineRoute finds the airline_route link between an
// existing physical route (origin/destination airport pair) and the given
// airline. If the physical route exists but isn't linked to that airline
// yet, it creates the link with status Pending — the employee can then see
// that their request is awaiting admin approval instead of a dead end, and
// the admin only has to activate it (from the existing Airline Routes
// screen) instead of creating it from scratch.
//
// Returns domain.ErrRouteNotFound when no route is configured at all for
// that origin/destination — that case still requires an admin to create the
// physical route first, same as today.
// The bool return is true when a new Pending link was created by this call,
// false when a link (of any status) already existed — callers use this to
// pick the right response message/status code.
func (i *AirlineRouteInteractor) ResolveOrCreatePendingAirlineRoute(ctx context.Context, traceID, airlineID, originAirportID, destinationAirportID string) (*domain.AirlineRoute, bool, error) {
	log := airlineRouteLog.WithTraceID(traceID)

	log.Info(logger.LogAirlineRouteResolve, "airline_id", airlineID, "origin_airport_id", originAirportID, "destination_airport_id", destinationAirportID)

	route, err := i.routeService.GetRouteByAirports(ctx, originAirportID, destinationAirportID)
	if err != nil {
		if err == domain.ErrRouteNotFound {
			log.Warn(logger.LogRouteNotFound, "origin_airport_id", originAirportID, "destination_airport_id", destinationAirportID)
		} else {
			log.Error(logger.LogAirlineRouteResolveError, "error", err)
		}
		return nil, false, err
	}

	if existing, err := i.service.GetAirlineRouteByRouteAndAirline(ctx, route.ID, airlineID); err == nil {
		log.Info(logger.LogAirlineRouteResolveOK, "route_id", route.ID, "airline_id", airlineID, "status", existing.Status, "created", false)
		return existing, false, nil
	} else if err != domain.ErrAirlineRouteNotFound {
		log.Error(logger.LogAirlineRouteResolveError, "error", err)
		return nil, false, err
	}

	newID := utils.Generate()
	err = helpers.RunWithTx(ctx, i.service, log, logger.LogAirlineRouteResolveError,
		func(ctx context.Context, tx output.Tx) error {
			saveErr := i.service.SaveAirlineRouteTx(ctx, tx, domain.AirlineRoute{
				ID:        newID,
				RouteID:   route.ID,
				AirlineID: airlineID,
				Status:    domain.AirlineRouteStatusPending,
			})
			if saveErr == domain.ErrAirlineRouteAlreadyExists {
				// A concurrent request created the same link first — that's
				// fine, we'll just re-fetch and return it below.
				return nil
			}
			return saveErr
		})
	if err != nil {
		log.Error(logger.LogAirlineRouteResolveError, "route_id", route.ID, "airline_id", airlineID, "error", err)
		return nil, false, err
	}

	result, err := i.service.GetAirlineRouteByRouteAndAirline(ctx, route.ID, airlineID)
	if err != nil {
		log.Error(logger.LogAirlineRouteResolveError, "route_id", route.ID, "airline_id", airlineID, "error", err)
		return nil, false, err
	}

	log.Info(logger.LogAirlineRouteResolveOK, "route_id", route.ID, "airline_id", airlineID, "status", result.Status, "created", true)
	return result, true, nil
}
