package handlers

import (
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/middleware"
	"github.com/champion19/api-flighthours/platform/logger"
	"github.com/gin-gonic/gin"
)

// GetFlightHoursSummary godoc
// @Summary Get flight hours summary for the authenticated employee
// @Description Returns total flight hours and breakdown by pilot role, filtered by period
// @Tags Flight Summary
// @Produce json
// @Param period query string true "Period type: monthly, bimonthly, quarterly, semiannual, annual, custom"
// @Param start_date query string false "Start date (YYYY-MM-DD), required if period=custom"
// @Param end_date query string false "End date (YYYY-MM-DD), required if period=custom"
// @Param reference_date query string false "Reference date for period calculation (defaults to today)"
// @Success 200 {object} FlightHoursSummaryResponse
// @Failure 400 {object} middleware.APIResponse
// @Failure 401 {object} middleware.APIResponse
// @Failure 500 {object} middleware.APIResponse
// @Security BearerAuth
// @Router /employees/flight-hours-summary [get]
func (h *handler) GetFlightHoursSummary() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := log.WithTraceID(traceID)

		log.Info(logger.LogFlightSummaryGet, "action", "flight_hours_summary")

		employee, ok := middleware.GetAuthenticatedUser(c)
		if !ok {
			log.Error(logger.LogFlightSummaryGetError, "error", "unauthorized")
			h.Response.Error(c, domain.MsgFlightSummaryGetErr)
			return
		}

		period := c.DefaultQuery("period", domain.PeriodMonthly)
		startDate := c.Query("start_date")
		endDate := c.Query("end_date")
		referenceDate := c.Query("reference_date")

		// Validate period
		if !domain.IsValidPeriod(period) {
			log.Warn(logger.LogFlightSummaryGetError, "invalid_period", period)
			h.Response.Error(c, domain.MsgFlightSummaryInvalid)
			return
		}

		// For custom period, start_date and end_date are required
		if period == domain.PeriodCustom {
			if startDate == "" || endDate == "" {
				log.Warn(logger.LogFlightSummaryGetError, "error", "custom period requires start_date and end_date")
				h.Response.Error(c, domain.MsgFlightSummaryInvalid)
				return
			}
		}

		summary, err := h.FlightSummaryInteractor.GetFlightHoursSummary(c.Request.Context(), traceID, employee.ID, period, startDate, endDate, referenceDate)
		if err != nil {
			log.Error(logger.LogFlightSummaryGetError, "error", err)
			h.Response.Error(c, domain.MsgFlightSummaryGetErr)
			return
		}

		response := FromDomainFlightHoursSummary(summary)

		log.Info(logger.LogFlightSummaryGetOK, "employee_id", employee.ID, "period", period)
		h.Response.SuccessWithData(c, domain.MsgFlightSummaryGetOK, response)
	}
}

// GetFlightAlerts godoc
// @Summary Get flight alerts for the authenticated employee
// @Description Returns active alerts (e.g., max consecutive hours, min landings)
// @Tags Flight Summary
// @Produce json
// @Success 200 {object} FlightAlertsResponse
// @Failure 401 {object} middleware.APIResponse
// @Failure 500 {object} middleware.APIResponse
// @Security BearerAuth
// @Router /employees/flight-alerts [get]
func (h *handler) GetFlightAlerts() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := log.WithTraceID(traceID)

		log.Info(logger.LogFlightSummaryAlerts, "action", "flight_alerts")

		employee, ok := middleware.GetAuthenticatedUser(c)
		if !ok {
			log.Error(logger.LogFlightSummaryAlertsError, "error", "unauthorized")
			h.Response.Error(c, domain.MsgFlightAlertsGetErr)
			return
		}

		alerts, err := h.FlightSummaryInteractor.GetFlightAlerts(c.Request.Context(), traceID, employee.ID)
		if err != nil {
			log.Error(logger.LogFlightSummaryAlertsError, "error", err)
			h.Response.Error(c, domain.MsgFlightAlertsGetErr)
			return
		}

		response := FromDomainFlightAlerts(alerts)

		log.Info(logger.LogFlightSummaryAlertsOK, "employee_id", employee.ID, "alert_count", len(alerts))
		h.Response.SuccessWithData(c, domain.MsgFlightAlertsGetOK, response)
	}
}

// GetRecentFlights godoc
// @Summary Get the last 5 flights for the authenticated employee
// @Description Returns the 5 most recent flights ordered by date descending
// @Tags Flight Summary
// @Produce json
// @Success 200 {array} DailyLogbookDetailResponse
// @Failure 401 {object} middleware.APIResponse
// @Failure 500 {object} middleware.APIResponse
// @Security BearerAuth
// @Router /employees/recent-flights [get]
func (h *handler) GetRecentFlights() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := log.WithTraceID(traceID)

		log.Info(logger.LogFlightSummaryRecentFlights, "action", "recent_flights")

		employee, ok := middleware.GetAuthenticatedUser(c)
		if !ok {
			log.Error(logger.LogFlightSummaryGetError, "error", "unauthorized")
			h.Response.Error(c, domain.MsgRecentFlightsGetErr)
			return
		}

		flights, err := h.FlightSummaryInteractor.GetRecentFlights(c.Request.Context(), traceID, employee.ID, 5)
		if err != nil {
			log.Error(logger.LogFlightSummaryGetError, "error", err)
			h.Response.Error(c, domain.MsgRecentFlightsGetErr)
			return
		}

		// Build response
		var responses []DailyLogbookDetailResponse
		for _, d := range flights {
			encodedID, _ := h.EncodeID(d.ID)
			encodedLogbookID, _ := h.EncodeID(d.DailyLogbookID)
			encodedRouteID, _ := h.EncodeID(d.AirlineRouteID)
			encodedAircraftID, _ := h.EncodeID(d.TailNumberID)

			response := FromDomainDailyLogbookDetail(&d, encodedID, encodedLogbookID, encodedRouteID, encodedAircraftID)
			responses = append(responses, response)
		}

		log.Info(logger.LogFlightSummaryGetOK, "employee_id", employee.ID, "count", len(responses))
		h.Response.SuccessWithData(c, domain.MsgRecentFlightsGetOK, responses)
	}
}
