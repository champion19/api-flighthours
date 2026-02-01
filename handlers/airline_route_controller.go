package handlers

import (
	"context"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/middleware"
	"github.com/gin-gonic/gin"
)

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
