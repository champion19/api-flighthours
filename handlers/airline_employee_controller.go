package handlers

import (
	"time"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/middleware"
	"github.com/champion19/api-flighthours/platform/logger"
	"github.com/gin-gonic/gin"
)

// GetEmployeeAirlineInfo godoc
// @Summary      Get authenticated employee's airline information (HU24)
// @Description  Returns the airline information (airline_id, airline_name, bp) of the authenticated employee
// @Tags         AirlineEmployees
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  EmployeeAirlineInfoResponse
// @Failure      401  {object}  middleware.ErrorResponse "Not authenticated"
// @Failure      404  {object}  middleware.ErrorResponse "Employee not found or no airline assigned"
// @Failure      500  {object}  middleware.ErrorResponse "Internal server error"
// @Router       /employee/me/airline [get]
func (h handler) GetEmployeeAirlineInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogAirlineEmployeeGet, "endpoint", "GET /employee/me/airline", "client_ip", c.ClientIP())

		// Get authenticated employee (core data)
		employee, exists := middleware.GetAuthenticatedUser(c)
		if !exists {
			log.Error(logger.LogEmployeeNotFound, "error", logger.LogErrAuthUserNotInContext, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgUnauthorized)
			return
		}

		// Get airline info for this employee
		airlineInfo, err := h.AirlineEmployeeInteractor.GetAirlineEmployeeByID(c.Request.Context(), employee.ID)

		response := EmployeeAirlineInfoResponse{}

		// If airline info exists, populate the response
		if err == nil && airlineInfo != nil && airlineInfo.AirlineID != "" {
			airline, err := h.AirlineInteractor.GetAirlineByID(c.Request.Context(), airlineInfo.AirlineID)
			if err == nil && airline != nil {
				encodedAirlineID, _ := h.EncodeID(airline.ID)
				response.AirlineID = encodedAirlineID
				response.AirlineName = airline.AirlineName
				response.AirlineCode = airline.AirlineCode
			}
			response.Bp = airlineInfo.Bp
		}

		baseURL := GetBaseURL(c)
		response.Links = []Link{
			{Href: baseURL + "/flighthours/api/v1/employee/me/airline", Rel: "self", Method: "GET"},
			{Href: baseURL + "/flighthours/api/v1/employee/me/airline", Rel: "update", Method: "PUT"},
			{Href: baseURL + "/flighthours/api/v1/employee/me", Rel: "profile", Method: "GET"},
		}

		log.Success(logger.LogAirlineEmployeeGetOK, "employee_id", employee.ID, "has_airline", response.AirlineID != "", "client_ip", c.ClientIP())
		h.Response.SuccessWithData(c, domain.MsgAirlineEmployeeGetOK, response)
	}
}

// AddEmployeeAirlineInfo godoc
// @Summary      Add airline information for authenticated employee (HU26)
// @Description  Adds airline information (airline_id, bp, start_date, end_date) for the authenticated employee. New employees are set as active=true by default. The 'active' field is not accepted - use activate/deactivate endpoints for that.
// @Tags         AirlineEmployees
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        airline body AddEmployeeAirlineRequest true "Airline information to add (airline_id, bp, start_date, end_date)"
// @Success      200  {object}  AddEmployeeAirlineResponse "Airline info added successfully"
// @Failure      400  {object}  middleware.ErrorResponse "Invalid request data or date format"
// @Failure      401  {object}  middleware.ErrorResponse "Not authenticated"
// @Failure      422  {object}  middleware.ErrorResponse "Invalid airline_id reference"
// @Failure      500  {object}  middleware.ErrorResponse "Internal server error"
// @Router       /employee/me/airline [put]
func (h handler) AddEmployeeAirlineInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogAirlineEmployeeAdd, "endpoint", "PUT /employee/me/airline", "client_ip", c.ClientIP())

		// Get authenticated employee (core data)
		employee, exists := middleware.GetAuthenticatedUser(c)
		if !exists {
			log.Error(logger.LogEmployeeNotFound, "error", logger.LogErrAuthUserNotInContext, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgUnauthorized)
			return
		}

		var req AddEmployeeAirlineRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error(logger.LogRegJSONParseError, "error", err, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgValInvalidReq)
			return
		}

		req.Sanitize()

		// Decode the airline_id
		airlineUUID, _ := h.resolveID(req.AirlineID)
		if airlineUUID == "" {
			log.Error(logger.LogAirlineEmployeeAddError, "error", logger.LogErrInvalidAirlineID, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgAirlineEmployeeInvalidAirline)
			return
		}

		// Verify airline exists
		airline, err := h.AirlineInteractor.GetAirlineByID(c.Request.Context(), airlineUUID)
		if err != nil || airline == nil {
			log.Error(logger.LogAirlineEmployeeAddError, "error", logger.LogErrAirlineNotFound, "airline_id", airlineUUID, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgAirlineEmployeeInvalidAirline)
			return
		}

		// Parse dates using local timezone to avoid date shift when saving to MySQL
		layout := "2006-01-02"
		startDate, err := time.ParseInLocation(layout, req.StartDate, time.Local)
		if err != nil {
			log.Error(logger.LogAirlineEmployeeAddError, "error", logger.LogErrInvalidStartDate, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgValInvalidDateFormat)
			return
		}

		var endDate time.Time
		if req.EndDate != "" {
			endDate, err = time.ParseInLocation(layout, req.EndDate, time.Local)
			if err != nil {
				log.Error(logger.LogAirlineEmployeeAddError, "error", logger.LogErrInvalidEndDate, "client_ip", c.ClientIP())
				h.Response.Error(c, domain.MsgValInvalidDateFormat)
				return
			}
		}

		// Create AirlineEmployee domain object (Active defaults to true for new employees)
		airlineEmployeeInfo := domain.AirlineEmployee{
			ID:        employee.ID,
			AirlineID: airlineUUID,
			Bp:        req.Bp,
			StartDate: startDate,
			EndDate:   endDate,
			Active:    true, // New airline employees are active by default
		}

		// Add airline info via AirlineEmployeeInteractor
		if err := h.AirlineEmployeeInteractor.AddAirlineEmployee(c.Request.Context(), employee.ID, airlineEmployeeInfo); err != nil {
			log.Error(logger.LogAirlineEmployeeAddError, "error", err, "client_ip", c.ClientIP())
			if err == domain.ErrInvalidForeignKey {
				h.Response.Error(c, domain.MsgAirlineEmployeeInvalidAirline)
				return
			}
			h.Response.Error(c, domain.MsgAirlineEmployeeSaveError)
			return
		}

		encodedAirlineID, _ := h.EncodeID(airline.ID)
		response := AddEmployeeAirlineResponse{
			AirlineID:   encodedAirlineID,
			AirlineName: airline.AirlineName,
			Bp:          req.Bp,
			StartDate:   startDate.Format(layout),
			EndDate:     endDate.Format(layout),
		}

		baseURL := GetBaseURL(c)
		response.Links = []Link{
			{Href: baseURL + "/flighthours/api/v1/employee/me/airline", Rel: "self", Method: "GET"},
			{Href: baseURL + "/flighthours/api/v1/employee/me", Rel: "profile", Method: "GET"},
		}

		log.Success(logger.LogAirlineEmployeeAddOK, "employee_id", employee.ID, "airline_id", airlineUUID, "client_ip", c.ClientIP())
		h.Response.SuccessWithData(c, domain.MsgAirlineEmployeeCreated, response)
	}
}

// UpdateEmployeeAirlineInfo godoc
// @Summary      Edit airline information for authenticated employee (HU25)
// @Description  Updates the airline information (airline_id, bp, start_date, end_date) for the authenticated employee who already has airline info. The 'active' field is not editable - use activate/deactivate endpoints. Requires existing airline info.
// @Tags         AirlineEmployees
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        airline body UpdateEmployeeAirlineRequest true "Airline information to update (airline_id, bp, start_date, end_date)"
// @Success      200  {object}  UpdateEmployeeAirlineResponse "Airline info updated successfully"
// @Failure      400  {object}  middleware.ErrorResponse "Invalid request data or date format"
// @Failure      401  {object}  middleware.ErrorResponse "Not authenticated"
// @Failure      404  {object}  middleware.ErrorResponse "Employee has no airline info to update"
// @Failure      422  {object}  middleware.ErrorResponse "Invalid airline_id reference"
// @Failure      500  {object}  middleware.ErrorResponse "Internal server error"
// @Router       /employee/me/airline-info [put]
func (h handler) UpdateEmployeeAirlineInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogAirlineEmployeeUpdate, "endpoint", "PUT /employee/me/airline-info", "client_ip", c.ClientIP())

		// Get authenticated employee (core data)
		employee, exists := middleware.GetAuthenticatedUser(c)
		if !exists {
			log.Error(logger.LogEmployeeNotFound, "error", logger.LogErrAuthUserNotInContext, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgUnauthorized)
			return
		}

		// Verify employee has existing airline info before allowing update
		existingAirlineInfo, err := h.AirlineEmployeeInteractor.GetAirlineEmployeeByID(c.Request.Context(), employee.ID)
		if err != nil || existingAirlineInfo == nil || existingAirlineInfo.AirlineID == "" {
			log.Error(logger.LogAirlineEmployeeNotFound, "error", logger.LogErrEmployeeNoAirlineInfo, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgAirlineEmployeeNotFound)
			return
		}

		var req UpdateEmployeeAirlineRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error(logger.LogRegJSONParseError, "error", err, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgValInvalidReq)
			return
		}

		req.Sanitize()

		// Decode the airline_id
		airlineUUID, _ := h.resolveID(req.AirlineID)
		if airlineUUID == "" {
			log.Error(logger.LogAirlineEmployeeUpdateError, "error", logger.LogErrInvalidAirlineID, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgAirlineEmployeeInvalidAirline)
			return
		}

		// Verify airline exists
		airline, err := h.AirlineInteractor.GetAirlineByID(c.Request.Context(), airlineUUID)
		if err != nil || airline == nil {
			log.Error(logger.LogAirlineEmployeeUpdateError, "error", logger.LogErrAirlineNotFound, "airline_id", airlineUUID, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgAirlineEmployeeInvalidAirline)
			return
		}

		// Parse dates using local timezone to avoid date shift when saving to MySQL
		layout := "2006-01-02"
		startDate, err := time.ParseInLocation(layout, req.StartDate, time.Local)
		if err != nil {
			log.Error(logger.LogAirlineEmployeeUpdateError, "error", logger.LogErrInvalidStartDate, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgValInvalidDateFormat)
			return
		}

		var endDate time.Time
		if req.EndDate != "" {
			endDate, err = time.ParseInLocation(layout, req.EndDate, time.Local)
			if err != nil {
				log.Error(logger.LogAirlineEmployeeUpdateError, "error", logger.LogErrInvalidEndDate, "client_ip", c.ClientIP())
				h.Response.Error(c, domain.MsgValInvalidDateFormat)
				return
			}
		}

		// Create AirlineEmployee domain object (Active is not updated in HU25)
		airlineEmployeeInfo := domain.AirlineEmployee{
			ID:        employee.ID,
			AirlineID: airlineUUID,
			Bp:        req.Bp,
			StartDate: startDate,
			EndDate:   endDate,
			Active:    existingAirlineInfo.Active, // Preserve existing active status
		}

		// Update airline info via AirlineEmployeeInteractor
		if err := h.AirlineEmployeeInteractor.UpdateAirlineEmployee(c.Request.Context(), employee.ID, airlineEmployeeInfo); err != nil {
			log.Error(logger.LogAirlineEmployeeUpdateError, "error", err, "client_ip", c.ClientIP())
			if err == domain.ErrInvalidForeignKey {
				h.Response.Error(c, domain.MsgAirlineEmployeeInvalidAirline)
				return
			}
			h.Response.Error(c, domain.MsgAirlineEmployeeUpdateError)
			return
		}

		encodedAirlineID, _ := h.EncodeID(airline.ID)
		response := UpdateEmployeeAirlineResponse{
			Updated:     true,
			AirlineID:   encodedAirlineID,
			AirlineName: airline.AirlineName,
			Bp:          req.Bp,
			StartDate:   startDate.Format(layout),
			EndDate:     endDate.Format(layout),
		}

		baseURL := GetBaseURL(c)
		response.Links = []Link{
			{Href: baseURL + "/flighthours/api/v1/employee/me/airline", Rel: "self", Method: "GET"},
			{Href: baseURL + "/flighthours/api/v1/employee/me", Rel: "profile", Method: "GET"},
		}

		log.Success(logger.LogAirlineEmployeeUpdateOK, "employee_id", employee.ID, "airline_id", airlineUUID, "client_ip", c.ClientIP())
		h.Response.SuccessWithData(c, domain.MsgAirlineEmployeeUpdated, response)
	}
}
