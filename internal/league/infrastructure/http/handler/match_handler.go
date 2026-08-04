package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/application/input"
	"github.com/jmlc643/twitbet-backend/internal/league/application/usecase"
	"github.com/jmlc643/twitbet-backend/internal/league/infrastructure/http/dto/request"
	"github.com/jmlc643/twitbet-backend/internal/league/infrastructure/http/dto/response"
)

type MatchHandler struct {
	createMatchUseCase  *usecase.CreateMatchUseCase
	createMarketUseCase *usecase.CreateMarketUseCase
}

func NewMatchHandler(createMatchUC *usecase.CreateMatchUseCase, createMarketUC *usecase.CreateMarketUseCase) *MatchHandler {
	return &MatchHandler{
		createMatchUseCase:  createMatchUC,
		createMarketUseCase: createMarketUC,
	}
}

func (h *MatchHandler) CreateMatch(c *gin.Context) {
	leagueIDStr := c.Param("id")
	leagueID, err := uuid.Parse(leagueIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de liga inválido"})
		return
	}

	var req request.CreateMatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	adminIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autorizado"})
		return
	}
	adminID, _ := uuid.Parse(adminIDStr.(string))

	match, err := h.createMatchUseCase.Execute(c.Request.Context(), adminID, leagueID, req.Title, req.StartTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response.NewMatchResponse(match))
}

func (h *MatchHandler) CreateMarketForLeague(c *gin.Context) {
	leagueIDStr := c.Param("id")
	leagueID, err := uuid.Parse(leagueIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de liga inválido"})
		return
	}

	var req request.CreateMarketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	adminIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autorizado"})
		return
	}
	adminID, _ := uuid.Parse(adminIDStr.(string))

	var inputOptions []input.MarketOptionInput
	for _, opt := range req.Options {
		inputOptions = append(inputOptions, input.MarketOptionInput{Name: opt.Name, Odds: opt.Odds})
	}

	market, err := h.createMarketUseCase.Execute(c.Request.Context(), adminID, leagueID, nil, req.Name, inputOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response.NewMarketResponse(market))
}

func (h *MatchHandler) CreateMarketForMatch(c *gin.Context) {
	matchIDStr := c.Param("id")
	matchID, err := uuid.Parse(matchIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de partido inválido"})
		return
	}

	var req request.CreateMarketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	adminIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autorizado"})
		return
	}
	adminID, _ := uuid.Parse(adminIDStr.(string))

	var inputOptions []input.MarketOptionInput
	for _, opt := range req.Options {
		inputOptions = append(inputOptions, input.MarketOptionInput{Name: opt.Name, Odds: opt.Odds})
	}

	market, err := h.createMarketUseCase.Execute(c.Request.Context(), adminID, uuid.Nil, &matchID, req.Name, inputOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response.NewMarketResponse(market))
}
