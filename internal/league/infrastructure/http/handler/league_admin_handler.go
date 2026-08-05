package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/application/input"
	"github.com/jmlc643/twitbet-backend/internal/league/infrastructure/http/dto/request"
)

func (h *LeagueHandler) AssignAdmin(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no autorizado"})
		return
	}

	OwnerID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "id de usuario inválido"})
		return
	}

	leagueIDStr := c.Param("id")
	leagueID, err := uuid.Parse(leagueIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id de liga inválido"})
		return
	}

	var req request.AssignAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	participantID, err := uuid.Parse(req.ParticipantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id de participante inválido"})
		return
	}

	in := input.AssignAdminInput{
		LeagueID:      leagueID,
		OwnerID:       OwnerID,
		ParticipantID: participantID,
	}

	err = h.assignAdminUC.Execute(c.Request.Context(), in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *LeagueHandler) RemoveAdmin(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no autorizado"})
		return
	}

	OwnerID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "id de usuario inválido"})
		return
	}

	leagueIDStr := c.Param("id")
	leagueID, err := uuid.Parse(leagueIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id de liga inválido"})
		return
	}

	participantIDStr := c.Param("participant_id")
	participantID, err := uuid.Parse(participantIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id de participante inválido"})
		return
	}

	in := input.RemoveAdminInput{
		LeagueID:      leagueID,
		OwnerID:       OwnerID,
		ParticipantID: participantID,
	}

	err = h.removeAdminUC.Execute(c.Request.Context(), in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
