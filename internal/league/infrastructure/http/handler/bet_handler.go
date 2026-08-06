package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/application/usecase"
	"github.com/jmlc643/twitbet-backend/internal/league/infrastructure/http/dto/request"
	"github.com/jmlc643/twitbet-backend/internal/league/infrastructure/http/dto/response"
)

type BetHandler struct {
	placeBetUseCase *usecase.PlaceBetUseCase
}

func NewBetHandler(placeBetUseCase *usecase.PlaceBetUseCase) *BetHandler {
	return &BetHandler{placeBetUseCase: placeBetUseCase}
}

func (h *BetHandler) PlaceBet(c *gin.Context) {
	var req request.PlaceBetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user token"})
		return
	}

	leagueID, err := uuid.Parse(req.LeagueID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid league id"})
		return
	}

	marketID, err := uuid.Parse(req.MarketID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid market id"})
		return
	}

	marketOptionID, err := uuid.Parse(req.MarketOptionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid market option id"})
		return
	}

	bet, err := h.placeBetUseCase.Execute(c.Request.Context(), userID, leagueID, marketID, marketOptionID, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res := response.BetResponse{
		ID:             bet.ID.String(),
		ParticipantID:  bet.ParticipantID.String(),
		MarketOptionID: bet.MarketOptionID.String(),
		Amount:         bet.Amount,
		Odds:           bet.Odds,
		PotentialWin:   bet.PotentialWin,
		Status:         string(bet.Status),
		PlacedAt:       bet.PlacedAt,
	}

	c.JSON(http.StatusCreated, res)
}
