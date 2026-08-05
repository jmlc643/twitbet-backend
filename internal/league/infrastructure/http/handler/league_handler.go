package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/application/input"
	"github.com/jmlc643/twitbet-backend/internal/league/application/usecase"
	"github.com/jmlc643/twitbet-backend/internal/league/infrastructure/http/dto/request"
	"github.com/jmlc643/twitbet-backend/internal/league/infrastructure/http/mapper"
)

type LeagueHandler struct {
	createLeagueUC     *usecase.CreateLeagueUseCase
	joinLeagueUC       *usecase.JoinLeagueUseCase
	getUserLeaguesUC   *usecase.GetUserLeaguesUseCase
	getLeagueDetailsUC *usecase.GetLeagueDetailsUseCase
	updateLeagueUC     *usecase.UpdateLeagueUseCase
	deleteLeagueUC     *usecase.DeleteLeagueUseCase
	assignAdminUC      *usecase.AssignAdminUseCase
	removeAdminUC      *usecase.RemoveAdminUseCase
}

func NewLeagueHandler(
	createLeagueUC *usecase.CreateLeagueUseCase,
	joinLeagueUC *usecase.JoinLeagueUseCase,
	getUserLeaguesUC *usecase.GetUserLeaguesUseCase,
	getLeagueDetailsUC *usecase.GetLeagueDetailsUseCase,
	updateLeagueUC *usecase.UpdateLeagueUseCase,
	deleteLeagueUC *usecase.DeleteLeagueUseCase,
	assignAdminUC *usecase.AssignAdminUseCase,
	removeAdminUC *usecase.RemoveAdminUseCase,
) *LeagueHandler {
	return &LeagueHandler{
		createLeagueUC:     createLeagueUC,
		joinLeagueUC:       joinLeagueUC,
		getUserLeaguesUC:   getUserLeaguesUC,
		getLeagueDetailsUC: getLeagueDetailsUC,
		updateLeagueUC:     updateLeagueUC,
		deleteLeagueUC:     deleteLeagueUC,
		assignAdminUC:      assignAdminUC,
		removeAdminUC:      removeAdminUC,
	}
}

func (h *LeagueHandler) CreateLeague(c *gin.Context) {
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

	var req request.CreateLeagueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := mapper.CreateLeagueRequestToInput(req, OwnerID)
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

func (h *LeagueHandler) GetUserLeagues(c *gin.Context) {
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

	in := input.GetUserLeaguesInput{
		UserID: userID,
	}

	output, err := h.getUserLeaguesUC.Execute(c.Request.Context(), in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mapper.GetUserLeaguesOutputToResponse(output))
}

func (h *LeagueHandler) GetLeagueDetails(c *gin.Context) {
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

	slug := c.Param("id")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "slug de liga inválido"})
		return
	}

	in := input.GetLeagueDetailsInput{
		Slug:   slug,
		UserID: userID,
	}

	output, err := h.getLeagueDetailsUC.Execute(c.Request.Context(), in)
	if err != nil {
		if err.Error() == "liga no encontrada" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mapper.GetLeagueDetailsOutputToResponse(output))
}

func (h *LeagueHandler) UpdateLeague(c *gin.Context) {
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

	var req request.UpdateLeagueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	in := input.UpdateLeagueInput{
		LeagueID:       leagueID,
		OwnerID:        OwnerID,
		Name:           req.Name,
		InitialBalance: req.InitialBalance,
		MaxRecharges:   req.MaxRecharges,
		HideStandings:  req.HideStandings,
	}

	err = h.updateLeagueUC.Execute(c.Request.Context(), in)
	if err != nil {
		if err.Error() == "liga no encontrada" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "solo el administrador puede modificar la liga" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *LeagueHandler) DeleteLeague(c *gin.Context) {
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

	in := input.DeleteLeagueInput{
		LeagueID: leagueID,
		OwnerID:  OwnerID,
	}

	err = h.deleteLeagueUC.Execute(c.Request.Context(), in)
	if err != nil {
		if err.Error() == "liga no encontrada" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "solo el administrador puede eliminar la liga" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}