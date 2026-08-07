package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/application/usecase"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
	"github.com/jmlc643/twitbet-backend/internal/league/infrastructure/http/dto/request"
	"github.com/jmlc643/twitbet-backend/internal/league/infrastructure/http/dto/response"
	"strconv"
	"time"
)

type BetHandler struct {
	placeBetUseCase    *usecase.PlaceBetUseCase
	getUserBetsUseCase *usecase.GetUserBetsUseCase
	cashoutBetUseCase  *usecase.CashoutBetUseCase
}

func NewBetHandler(placeBetUseCase *usecase.PlaceBetUseCase, getUserBetsUseCase *usecase.GetUserBetsUseCase, cashoutBetUseCase *usecase.CashoutBetUseCase) *BetHandler {
	return &BetHandler{
		placeBetUseCase:    placeBetUseCase,
		getUserBetsUseCase: getUserBetsUseCase,
		cashoutBetUseCase:  cashoutBetUseCase,
	}
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

	var bonusID *uuid.UUID
	if req.BonusID != nil {
		id, err := uuid.Parse(*req.BonusID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bonus id"})
			return
		}
		bonusID = &id
	}

	bet, err := h.placeBetUseCase.Execute(c.Request.Context(), userID, leagueID, marketID, marketOptionID, req.Amount, bonusID)
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

func (h *BetHandler) GetUserBets(c *gin.Context) {
	leagueIDStr := c.Param("id")
	leagueID, err := uuid.Parse(leagueIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de liga inválido"})
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

	status := c.Query("status")
	var statusPtr *string
	if status != "" {
		statusPtr = &status
	}

	var startDatePtr, endDatePtr *time.Time
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if t, err := time.Parse(time.RFC3339, startDateStr); err == nil {
			startDatePtr = &t
		}
	}
	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if t, err := time.Parse(time.RFC3339, endDateStr); err == nil {
			endDatePtr = &t
		}
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 {
		limit = 10
	}

	bets, total, err := h.getUserBetsUseCase.Execute(c.Request.Context(), userID, leagueID, statusPtr, startDatePtr, endDatePtr, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if bets == nil {
		bets = []entity.BetDetail{}
	}

	var responses []response.BetDetailResponse
	for _, b := range bets {
		responses = append(responses, response.BetDetailResponse{
			ID:           b.ID.String(),
			Amount:       b.Amount,
			Odds:         b.Odds,
			PotentialWin: b.PotentialWin,
			Status:       string(b.Status),
			PlacedAt:     b.PlacedAt,
			MatchTitle:   b.MatchTitle,
			MarketID:     b.MarketID.String(),
			MarketName:   b.MarketName,
			OptionID:     b.OptionID.String(),
			OptionName:   b.OptionName,
			CashoutAmount: b.CashoutAmount,
		})
	}
	if responses == nil {
		responses = []response.BetDetailResponse{}
	}

	c.JSON(http.StatusOK, response.PaginatedBetResponse{
		Data: responses,
		Meta: response.PaginationMeta{
			Total: total,
			Page:  page,
			Limit: limit,
		},
	})
}

func (h *BetHandler) Cashout(c *gin.Context) {
	betIDStr := c.Param("id")
	betID, err := uuid.Parse(betIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bet ID"})
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

	bet, err := h.cashoutBetUseCase.Execute(c.Request.Context(), userID, betID)
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

	c.JSON(http.StatusOK, res)
}
