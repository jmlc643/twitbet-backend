package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/application/usecase"
	"github.com/jmlc643/twitbet-backend/internal/league/infrastructure/http/dto/request"
	"github.com/jmlc643/twitbet-backend/internal/league/infrastructure/http/mapper"
)

type LeagueHandler struct {
	createLeagueUC *usecase.CreateLeagueUseCase
	joinLeagueUC   *usecase.JoinLeagueUseCase
}

func NewLeagueHandler(createLeagueUC *usecase.CreateLeagueUseCase, joinLeagueUC *usecase.JoinLeagueUseCase) *LeagueHandler {
	return &LeagueHandler{
		createLeagueUC: createLeagueUC,
		joinLeagueUC:   joinLeagueUC,
	}
}

func (h *LeagueHandler) CreateLeague(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no autorizado"})
		return
	}

	adminID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "id de usuario inválido"})
		return
	}

	var req request.CreateLeagueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := mapper.CreateLeagueRequestToInput(req, adminID)
	output, err := h.createLeagueUC.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, mapper.CreateLeagueOutputToResponse(output))
}

func (h *LeagueHandler) JoinLeague(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no autorizado"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "id de usuario inválido"})
		return
	}

	var req request.JoinLeagueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := mapper.JoinLeagueRequestToInput(req, userID)
	output, err := h.joinLeagueUC.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mapper.JoinLeagueOutputToResponse(output))
}
