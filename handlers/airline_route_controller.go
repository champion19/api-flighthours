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
