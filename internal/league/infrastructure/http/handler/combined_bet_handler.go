package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/application/combined_bet"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/apperror"
	"github.com/jmlc643/twitbet-backend/internal/league/infrastructure/http/dto/request"
	"github.com/jmlc643/twitbet-backend/internal/league/infrastructure/http/dto/response"
	"github.com/jmlc643/twitbet-backend/internal/league/infrastructure/http/mapper"
)

type CombinedBetHandler struct {
	placeCombinedBetUC   *combined_bet.PlaceCombinedBetUseCase
	getUserCombinedBetsUC *combined_bet.GetUserCombinedBetsUseCase
	getCombinedBetDetailsUC *combined_bet.GetCombinedBetDetailsUseCase
	cashoutCombinedBetUC *combined_bet.CashoutCombinedBetUseCase
}

func NewCombinedBetHandler(
	placeCombinedBetUC *combined_bet.PlaceCombinedBetUseCase,
	getUserCombinedBetsUC *combined_bet.GetUserCombinedBetsUseCase,
	getCombinedBetDetailsUC *combined_bet.GetCombinedBetDetailsUseCase,
	cashoutCombinedBetUC *combined_bet.CashoutCombinedBetUseCase,
) *CombinedBetHandler {
	return &CombinedBetHandler{
		placeCombinedBetUC:      placeCombinedBetUC,
		getUserCombinedBetsUC:   getUserCombinedBetsUC,
		getCombinedBetDetailsUC: getCombinedBetDetailsUC,
		cashoutCombinedBetUC:    cashoutCombinedBetUC,
	}
}

func (h *CombinedBetHandler) PlaceCombinedBet(c *gin.Context) {
	var req request.PlaceCombinedBetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		return
	}
	
	var userID uuid.UUID
	switch v := userIDValue.(type) {
	case string:
		parsed, err := uuid.Parse(v)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "ID de usuario inválido"})
			return
		}
		userID = parsed
	case uuid.UUID:
		userID = v
	}

	appInput := mapper.ToPlaceCombinedBetInput(req)

	out, err := h.placeCombinedBetUC.Execute(c.Request.Context(), userID, appInput)
	if err != nil {
		h.handleDomainError(c, err)
		return
	}

	res := mapper.ToCombinedBetResponse(*out)
	c.JSON(http.StatusCreated, res)
}

func (h *CombinedBetHandler) GetUserCombinedBets(c *gin.Context) {
	leagueIDStr := c.Query("league_id")
	if leagueIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "league_id es requerido"})
		return
	}

	leagueID, err := uuid.Parse(leagueIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "league_id inválido"})
		return
	}

	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		return
	}
	
	var userID uuid.UUID
	switch v := userIDValue.(type) {
	case string:
		parsed, err := uuid.Parse(v)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "ID de usuario inválido"})
			return
		}
		userID = parsed
	case uuid.UUID:
		userID = v
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

	out, total, err := h.getUserCombinedBetsUC.Execute(c.Request.Context(), userID, leagueID, statusPtr, startDatePtr, endDatePtr, page, limit)
	if err != nil {
		h.handleDomainError(c, err)
		return
	}

	res := make([]response.CombinedBetResponse, len(out))
	for i, o := range out {
		res[i] = mapper.ToCombinedBetResponse(o)
	}

	c.JSON(http.StatusOK, response.PaginatedCombinedBetResponse{
		Data: res,
		Meta: response.PaginationMeta{
			Total: total,
			Page:  page,
			Limit: limit,
		},
	})
}

func (h *CombinedBetHandler) GetCombinedBetDetails(c *gin.Context) {
	betID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		return
	}
	
	var userID uuid.UUID
	switch v := userIDValue.(type) {
	case string:
		parsed, err := uuid.Parse(v)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "ID de usuario inválido"})
			return
		}
		userID = parsed
	case uuid.UUID:
		userID = v
	}

	out, err := h.getCombinedBetDetailsUC.Execute(c.Request.Context(), betID, userID)
	if err != nil {
		h.handleDomainError(c, err)
		return
	}

	res := mapper.ToCombinedBetResponse(*out)
	c.JSON(http.StatusOK, res)
}

func (h *CombinedBetHandler) Cashout(c *gin.Context) {
	betID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		return
	}
	
	var userID uuid.UUID
	switch v := userIDValue.(type) {
	case string:
		parsed, err := uuid.Parse(v)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "ID de usuario inválido"})
			return
		}
		userID = parsed
	case uuid.UUID:
		userID = v
	}

	cashoutValue, err := h.cashoutCombinedBetUC.Execute(c.Request.Context(), userID, betID)
	if err != nil {
		h.handleDomainError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Cashout successful",
		"cashout_value": cashoutValue,
	})
}

func (h *CombinedBetHandler) handleDomainError(c *gin.Context, err error) {
	switch err {
	case apperror.ErrLeagueNotFound, apperror.ErrMarketNotFound, apperror.ErrMarketOptionNotFound, apperror.ErrBetNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case apperror.ErrInsufficientBalance, apperror.ErrInvalidBetAmount, apperror.ErrMarketNotActive, apperror.ErrMarketOptionBlocked, apperror.ErrInvalidMarketOptions:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case apperror.ErrUnauthorized:
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
