package handlers

import (
	"net/http"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/middleware"
	"github.com/champion19/api-flighthours/platform/logger"
	"github.com/gin-gonic/gin"
)

// ListDailyLogbooks godoc
// @Summary      List daily logbooks for authenticated employee
// @Description  Returns a list of daily logbooks for the currently authenticated employee
// @Tags         DailyLogbooks
// @Produce      json
// @Param        status query bool false "Filter by status (true for active, false for inactive)"
// @Success      200  {object}  DailyLogbookListResponse
// @Failure      401  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /daily-logbooks [get]
// @Security     BearerAuth
func (h *handler) ListDailyLogbooks() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get authenticated employee from context
		employee, ok := middleware.GetAuthenticatedUser(c)
		if !ok || employee == nil {
			c.Error(domain.ErrUserNotFound)
			return
		}

		filters := make(map[string]interface{})
		if statusParam := c.Query("status"); statusParam != "" {
			if statusParam == "true" {
				filters["status"] = true
			} else if statusParam == "false" {
				filters["status"] = false
			}
		}

		logbooks, err := h.DailyLogbookInteractor.ListDailyLogbooksByEmployee(c.Request.Context(), employee.ID, filters)
		if err != nil {
			Logger.Error(logger.LogDailyLogbookListError, "employee_id", employee.ID, "error", err)
			c.Error(err)
			return
		}

		baseURL := GetBaseURL(c)
		response := ToDailyLogbookListResponse(logbooks, h.EncodeID, baseURL)

		c.JSON(http.StatusOK, response)
	}
}

// GetDailyLogbookByID godoc
// @Summary      Get daily logbook by ID
// @Description  Returns daily logbook information by ID for the authenticated employee (accepts both UUID and obfuscated ID)
// @Tags         DailyLogbooks
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Daily Logbook ID (obfuscated ID)"
// @Success      200  {object}  DailyLogbookResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /daily-logbooks/{id} [get]
// @Security     BearerAuth
func (h *handler) GetDailyLogbookByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get authenticated employee from context
		employee, ok := middleware.GetAuthenticatedUser(c)
		if !ok || employee == nil {
			c.Error(domain.ErrUserNotFound)
			return
		}

		inputID := c.Param("id")
		if inputID == "" {
			c.Error(domain.ErrInvalidID)
			return
		}

		// Resolve ID (accepts both UUID and obfuscated ID)
		logbookUUID, responseID := h.resolveID(inputID)
		if logbookUUID == "" {
			h.HandleIDDecodingError(c, inputID, domain.ErrInvalidID)
			return
		}

		logbook, err := h.DailyLogbookInteractor.GetDailyLogbookByID(c.Request.Context(), logbookUUID)
		if err != nil {
			Logger.Error(logger.LogDailyLogbookGetError, "logbook_id", logbookUUID, "error", err)
			c.Error(err)
			return
		}

		// Verify ownership
		if logbook.EmployeeID != employee.ID {
			c.Error(domain.ErrDailyLogbookUnauthorized)
			return
		}

		// Encode employee ID for response
		encodedEmployeeID, _ := h.EncodeID(employee.ID)

		baseURL := GetBaseURL(c)
		response := FromDomainDailyLogbook(logbook, responseID, encodedEmployeeID)
		response.Links = BuildDailyLogbookLinks(baseURL, responseID)

		c.JSON(http.StatusOK, response)
	}
}

// CreateDailyLogbook godoc
// @Summary      Create a new daily logbook
// @Description  Creates a new daily logbook for the authenticated employee
// @Tags         DailyLogbooks
// @Accept       json
// @Produce      json
// @Param        request body CreateDailyLogbookRequest true "Daily logbook data"
// @Success      201  {object}  DailyLogbookResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /daily-logbooks [post]
// @Security     BearerAuth
func (h *handler) CreateDailyLogbook() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get authenticated employee from context
		employee, ok := middleware.GetAuthenticatedUser(c)
		if !ok || employee == nil {
			c.Error(domain.ErrUserNotFound)
			return
		}

		var req CreateDailyLogbookRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			Logger.Error(logger.LogDailyLogbookCreateError, "error", "invalid request body")
			c.Error(domain.ErrInvalidRequest)
			return
		}

		// Sanitize input data
		req.Sanitize()

		logbook, err := req.ToDomain(employee.ID)
		if err != nil {
			Logger.Error(logger.LogDailyLogbookCreateError, "error", err)
			c.Error(err)
			return
		}

		if err := h.DailyLogbookInteractor.CreateDailyLogbook(c.Request.Context(), *logbook); err != nil {
			Logger.Error(logger.LogDailyLogbookCreateError, "error", err)
			c.Error(err)
			return
		}

		encodedID, err := h.EncodeID(logbook.ID)
		if err != nil {
			h.HandleIDEncodingError(c, logbook.ID, err)
			return
		}

		// Encode employee ID for response
		encodedEmployeeID, _ := h.EncodeID(employee.ID)

		baseURL := GetBaseURL(c)
		response := FromDomainDailyLogbook(logbook, encodedID, encodedEmployeeID)
		response.Links = BuildDailyLogbookCreatedLinks(baseURL, encodedID)

		SetLocationHeader(c, baseURL, "daily-logbooks", encodedID)
		c.JSON(http.StatusCreated, response)
	}
}


