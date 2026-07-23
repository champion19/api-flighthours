package handlers

import (
	"context"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/middleware"
	"github.com/gin-gonic/gin"
)

// ListAirlineRoutes godoc
// @Summary      List all airline routes
// @Description  Returns a list of all airline-route associations with optional filters
// @Tags         Airline Routes
// @Accept       json
// @Produce      json
// @Param        airline_code query string false "Filter by airline code"
// @Param        status query string false "Filter by status (true/false/1/0)"
// @Security     BearerAuth
// @Success      200  {object}  middleware.APIResponse "Airline routes list"
// @Failure      401  {object}  middleware.ErrorResponse "Not authenticated"
// @Failure      500  {object}  middleware.ErrorResponse "Internal server error"
// @Router       /airline-routes [get]
func (h *handler) ListAirlineRoutes() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		ctx := context.Background()

		filters := make(map[string]interface{})

		if airlineCode := c.Query("airline_code"); airlineCode != "" {
			filters["airline_code"] = airlineCode
		}

		if status := c.Query("status"); status != "" {
			if status == "true" || status == "1" {
				filters["status"] = true
			} else if status == "false" || status == "0" {
				filters["status"] = false
			}
		}

		airlineRoutes, err := h.AirlineRouteInteractor.ListAirlineRoutes(ctx, traceID, filters)
		if err != nil {
			c.Error(err)
			return
		}

		baseURL := GetBaseURL(c)
		response := ToAirlineRouteListResponse(airlineRoutes, h.EncodeID, baseURL)

		h.Response.SuccessWithData(c, domain.MsgAirlineRouteListOK, response)
	}
}

// ActivateAirlineRoute godoc
// @Summary      Activate an airline route
// @Description  Sets the airline route status to active. Idempotent operation.
// @Tags         Airline Routes
// @Accept       json
// @Produce      json
// @Param        id path string true "Airline Route ID (obfuscated ID)"
// @Security     BearerAuth
// @Success      200  {object}  middleware.APIResponse "Airline route activated"
// @Failure      400  {object}  middleware.ErrorResponse "Invalid ID"
// @Failure      401  {object}  middleware.ErrorResponse "Not authenticated"
// @Failure      404  {object}  middleware.ErrorResponse "Airline route not found"
// @Failure      500  {object}  middleware.ErrorResponse "Internal server error"
// @Router       /airline-routes/{id}/activate [patch]
func (h *handler) ActivateAirlineRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		ctx := context.Background()

		id := c.Param("id")
		uuid, responseID := h.resolveID(id)

		if uuid == "" {
			c.Error(domain.ErrInvalidID)
			return
		}

		err := h.AirlineRouteInteractor.ActivateAirlineRoute(ctx, traceID, uuid)
		if err != nil {
			if err == domain.ErrAirlineRouteAlreadyActive {
				airlineRoute, _ := h.AirlineRouteInteractor.GetAirlineRouteByID(ctx, traceID, uuid)
				encodedAirlineID := ""
				if airlineRoute != nil {
					encodedAirlineID, _ = h.EncodeID(airlineRoute.AirlineID)
				}
				baseURL := GetBaseURL(c)
				links := BuildAirlineRouteStatusLinks(baseURL, responseID, true)
				h.Response.SuccessWithData(c, domain.MsgAirlineRouteActivateOK, gin.H{
					"id":         responseID,
					"airline_id": encodedAirlineID,
					"status":     true,
					"_links":     links,
				})
				return
			}
			c.Error(err)
			return
		}

		airlineRoute, _ := h.AirlineRouteInteractor.GetAirlineRouteByID(ctx, traceID, uuid)
		encodedAirlineID := ""
		if airlineRoute != nil {
			encodedAirlineID, _ = h.EncodeID(airlineRoute.AirlineID)
		}

		baseURL := GetBaseURL(c)
		links := BuildAirlineRouteStatusLinks(baseURL, responseID, true)

		h.Response.SuccessWithData(c, domain.MsgAirlineRouteActivateOK, gin.H{
			"id":         responseID,
			"airline_id": encodedAirlineID,
			"status":     true,
			"_links":     links,
		})
	}
}

// DeactivateAirlineRoute godoc
// @Summary      Deactivate an airline route
// @Description  Sets the airline route status to inactive. Idempotent operation.
// @Tags         Airline Routes
// @Accept       json
// @Produce      json
// @Param        id path string true "Airline Route ID (obfuscated ID)"
// @Security     BearerAuth
// @Success      200  {object}  middleware.APIResponse "Airline route deactivated"
// @Failure      400  {object}  middleware.ErrorResponse "Invalid ID"
// @Failure      401  {object}  middleware.ErrorResponse "Not authenticated"
// @Failure      404  {object}  middleware.ErrorResponse "Airline route not found"
// @Failure      500  {object}  middleware.ErrorResponse "Internal server error"
// @Router       /airline-routes/{id}/deactivate [patch]
func (h *handler) DeactivateAirlineRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		ctx := context.Background()

		id := c.Param("id")
		uuid, responseID := h.resolveID(id)

		if uuid == "" {
			c.Error(domain.ErrInvalidID)
			return
		}

		err := h.AirlineRouteInteractor.DeactivateAirlineRoute(ctx, traceID, uuid)
		if err != nil {
			if err == domain.ErrAirlineRouteAlreadyInactive {
				airlineRoute, _ := h.AirlineRouteInteractor.GetAirlineRouteByID(ctx, traceID, uuid)
				encodedAirlineID := ""
				if airlineRoute != nil {
					encodedAirlineID, _ = h.EncodeID(airlineRoute.AirlineID)
				}
				baseURL := GetBaseURL(c)
				links := BuildAirlineRouteStatusLinks(baseURL, responseID, false)
				h.Response.SuccessWithData(c, domain.MsgAirlineRouteDeactivateOK, gin.H{
					"id":         responseID,
					"airline_id": encodedAirlineID,
					"status":     false,
					"_links":     links,
				})
				return
			}
			c.Error(err)
			return
		}

		airlineRoute, _ := h.AirlineRouteInteractor.GetAirlineRouteByID(ctx, traceID, uuid)
		encodedAirlineID := ""
		if airlineRoute != nil {
			encodedAirlineID, _ = h.EncodeID(airlineRoute.AirlineID)
		}

		baseURL := GetBaseURL(c)
		links := BuildAirlineRouteStatusLinks(baseURL, responseID, false)

		h.Response.SuccessWithData(c, domain.MsgAirlineRouteDeactivateOK, gin.H{
			"id":         responseID,
			"airline_id": encodedAirlineID,
			"status":     false,
			"_links":     links,
		})
	}
}

// ListMyAirlineRoutes godoc
// @Summary      List routes for authenticated user's airline
// @Description  Returns the airline routes associated with the authenticated employee's airline
// @Tags         Airline Routes
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  middleware.APIResponse "My airline routes list"
// @Failure      401  {object}  middleware.ErrorResponse "Not authenticated"
// @Failure      404  {object}  middleware.ErrorResponse "Employee or airline association not found"
// @Failure      500  {object}  middleware.ErrorResponse "Internal server error"
// @Router       /employees/airline-routes [get]
func (h *handler) ListMyAirlineRoutes() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		ctx := context.Background()

		employee, exists := middleware.GetAuthenticatedUser(c)
		if !exists {
			c.Error(domain.ErrUserNotFound)
			return
		}

		airlineInfo, err := h.AirlineEmployeeInteractor.GetAirlineEmployeeByID(ctx, employee.ID)
		if err != nil || airlineInfo == nil || airlineInfo.AirlineID == "" {
			c.Error(domain.ErrAirlineEmployeeNotFound)
			return
		}

		airlineRoutes, err := h.AirlineRouteInteractor.ListMyAirlineRoutes(ctx, traceID, airlineInfo.AirlineID)
		if err != nil {
			c.Error(err)
			return
		}

		baseURL := GetBaseURL(c)
		response := ToAirlineRouteListResponse(airlineRoutes, h.EncodeID, baseURL)

		h.Response.SuccessWithData(c, domain.MsgAirlineRouteListOK, response)
	}
}

// ResolveAirlineRouteRequest is the body for POST /employees/airline-routes/resolve
type ResolveAirlineRouteRequest struct {
	OriginAirportID      string `json:"origin_airport_id" binding:"required"`
	DestinationAirportID string `json:"destination_airport_id" binding:"required"`
}

// ResolveAirlineRoute godoc
// @Summary      Resolve or auto-request the authenticated user's airline link for a route
// @Description  Given an origin/destination airport pair, returns the airline_route link for the
// @Description  authenticated employee's airline. If the physical route exists but isn't linked to
// @Description  their airline yet, creates the link with status "pending" for an admin to approve.
// @Tags         Airline Routes
// @Accept       json
// @Produce      json
// @Param        request body ResolveAirlineRouteRequest true "Origin/destination airport IDs"
// @Security     BearerAuth
// @Success      200  {object}  middleware.APIResponse "Link already existed"
// @Success      201  {object}  middleware.APIResponse "Pending link created"
// @Failure      400  {object}  middleware.ErrorResponse "Invalid request body"
// @Failure      401  {object}  middleware.ErrorResponse "Not authenticated"
// @Failure      404  {object}  middleware.ErrorResponse "No route configured for that origin/destination"
// @Failure      500  {object}  middleware.ErrorResponse "Internal server error"
// @Router       /employees/airline-routes/resolve [post]
func (h *handler) ResolveAirlineRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		ctx := context.Background()

		employee, exists := middleware.GetAuthenticatedUser(c)
		if !exists {
			c.Error(domain.ErrUserNotFound)
			return
		}

		var req ResolveAirlineRouteRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			h.Response.Error(c, domain.MsgValJSONInvalid)
			return
		}

		originUUID, _ := h.resolveID(req.OriginAirportID)
		destinationUUID, _ := h.resolveID(req.DestinationAirportID)
		if originUUID == "" || destinationUUID == "" {
			h.Response.Error(c, domain.MsgAirlineRouteInvalidRoute)
			return
		}

		airlineInfo, err := h.AirlineEmployeeInteractor.GetAirlineEmployeeByID(ctx, employee.ID)
		if err != nil || airlineInfo == nil || airlineInfo.AirlineID == "" {
			c.Error(domain.ErrAirlineEmployeeNotFound)
			return
		}

		airlineRoute, created, err := h.AirlineRouteInteractor.ResolveOrCreatePendingAirlineRoute(
			ctx, traceID, airlineInfo.AirlineID, originUUID, destinationUUID,
		)
		if err != nil {
			c.Error(err)
			return
		}

		encodedID, err := h.EncodeID(airlineRoute.ID)
		if err != nil {
			encodedID = airlineRoute.ID
		}
		encodedRouteID, err := h.EncodeID(airlineRoute.RouteID)
		if err != nil {
			encodedRouteID = airlineRoute.RouteID
		}
		encodedAirlineID, err := h.EncodeID(airlineRoute.AirlineID)
		if err != nil {
			encodedAirlineID = airlineRoute.AirlineID
		}
		encodedOriginAirportID, err := h.EncodeID(airlineRoute.OriginAirportID)
		if err != nil {
			encodedOriginAirportID = airlineRoute.OriginAirportID
		}
		encodedDestinationAirportID, err := h.EncodeID(airlineRoute.DestinationAirportID)
		if err != nil {
			encodedDestinationAirportID = airlineRoute.DestinationAirportID
		}

		response := AirlineRouteResponse{
			ID:                     encodedID,
			RouteID:                encodedRouteID,
			AirlineID:              encodedAirlineID,
			Status:                 airlineRoute.Status,
			AirlineCode:            airlineRoute.AirlineCode,
			AirlineName:            airlineRoute.AirlineName,
			OriginAirportID:        encodedOriginAirportID,
			OriginIataCode:         airlineRoute.OriginIataCode,
			OriginOaciCode:         airlineRoute.OriginOaciCode,
			DestinationAirportID:   encodedDestinationAirportID,
			DestinationIataCode:    airlineRoute.DestinationIataCode,
			DestinationOaciCode:    airlineRoute.DestinationOaciCode,
			RouteCode:              airlineRoute.RouteCode,
			OriginAirportName:      airlineRoute.OriginAirportName,
			DestinationAirportName: airlineRoute.DestinationAirportName,
			AirportType:            airlineRoute.AirportType,
			EstimatedFlightTime:    airlineRoute.EstimatedFlightTime,
		}

		if created {
			h.Response.SuccessWithData(c, domain.MsgAirlineRouteCreatedPendingOK, response)
			return
		}
		h.Response.SuccessWithData(c, domain.MsgAirlineRouteResolvedOK, response)
	}
}
