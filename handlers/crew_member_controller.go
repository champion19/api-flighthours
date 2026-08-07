package handlers

import (
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/middleware"
	"github.com/champion19/api-flighthours/platform/logger"
	"github.com/gin-gonic/gin"
)

// SearchCrewMembers returns the authenticated pilot's own crew roster, optionally
// filtered by name — this is also how "have I flown with this person before" is answered:
// if they're in the roster, they've been assigned to a flight before.
// @Summary Search crew members
// @Description Searches the authenticated pilot's own roster of command crew/cabin crew
// @Tags CrewMembers
// @Produce json
// @Param search query string false "Name filter (partial match)"
// @Success 200 {array} CrewMemberResponse
// @Router /crew-members [get]
// @Security BearerAuth
func (h *handler) SearchCrewMembers() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := log.WithTraceID(traceID)

		employee, ok := middleware.GetAuthenticatedUser(c)
		if !ok {
			log.Error(logger.LogCrewMemberSearchError, "error", "unauthorized")
			h.Response.Error(c, domain.MsgUnauthorized)
			return
		}

		query := c.Query("search")

		members, err := h.CrewMemberInteractor.SearchCrewMembers(c.Request.Context(), traceID, employee.ID, query)
		if err != nil {
			log.Error(logger.LogCrewMemberSearchError, "error", err)
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		responses := make([]CrewMemberResponse, 0, len(members))
		for i := range members {
			encodedID, _ := h.EncodeID(members[i].ID)
			responses = append(responses, FromDomainCrewMember(&members[i], encodedID))
		}

		log.Info(logger.LogCrewMemberSearchOK, "count", len(responses))
		h.Response.SuccessWithData(c, domain.MsgCrewMemberSearchOK, responses)
	}
}

// CreateCrewMember adds a person to the authenticated pilot's own crew roster,
// or returns the existing one if a crew member with that name already exists.
// @Summary Add crew member
// @Description Adds a person to the authenticated pilot's crew roster (or returns the existing match by name)
// @Tags CrewMembers
// @Accept json
// @Produce json
// @Param body body CreateCrewMemberRequest true "Crew member data"
// @Success 200 {object} CrewMemberResponse
// @Failure 400 {object} middleware.APIResponse
// @Router /crew-members [post]
// @Security BearerAuth
func (h *handler) CreateCrewMember() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := log.WithTraceID(traceID)

		employee, ok := middleware.GetAuthenticatedUser(c)
		if !ok {
			log.Error(logger.LogCrewMemberCreateError, "error", "unauthorized")
			h.Response.Error(c, domain.MsgUnauthorized)
			return
		}

		var req CreateCrewMemberRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error(logger.LogCrewMemberCreateError, "error", err)
			h.Response.Error(c, domain.MsgValJSONInvalid)
			return
		}
		req.Sanitize()

		if req.Name == "" {
			log.Warn(logger.LogCrewMemberCreateError, "error", "empty name")
			h.Response.Error(c, domain.MsgValInvalidReq)
			return
		}

		var bp *string
		if req.BP != "" {
			bp = &req.BP
		}

		member, err := h.CrewMemberInteractor.CreateCrewMember(c.Request.Context(), traceID, employee.ID, req.Name, bp)
		if err != nil {
			log.Error(logger.LogCrewMemberCreateError, "error", err)
			if err == domain.ErrInvalidRequest {
				h.Response.Error(c, domain.MsgValInvalidReq)
				return
			}
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		encodedID, _ := h.EncodeID(member.ID)

		log.Info(logger.LogCrewMemberCreateOK, "id", member.ID)
		h.Response.SuccessWithData(c, domain.MsgCrewMemberCreated, FromDomainCrewMember(member, encodedID))
	}
}
