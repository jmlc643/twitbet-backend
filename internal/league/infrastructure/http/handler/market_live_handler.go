package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/application/usecase"
	"github.com/jmlc643/twitbet-backend/internal/league/infrastructure/http/dto/request"
)

type MarketLiveHandler struct {
	updateStatusUseCase usecase.UpdateMarketStatusUseCase
	updateOddsUseCase   usecase.UpdateMarketOddsUseCase
}

func NewMarketLiveHandler(updateStatusUseCase usecase.UpdateMarketStatusUseCase, updateOddsUseCase usecase.UpdateMarketOddsUseCase) *MarketLiveHandler {
	return &MarketLiveHandler{
		updateStatusUseCase: updateStatusUseCase,
		updateOddsUseCase:   updateOddsUseCase,
	}
}

func (h *MarketLiveHandler) UpdateStatus(c *gin.Context) {
	adminIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no autorizado"})
		return
	}

	adminID, err := uuid.Parse(adminIDStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ID de usuario inválido"})
		return
	}

	marketIDStr := c.Param("id")
	marketID, err := uuid.Parse(marketIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de mercado inválido"})
		return
	}

	var req request.UpdateMarketStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payload inválido"})
		return
	}

	if err := h.updateStatusUseCase.Execute(c.Request.Context(), marketID, adminID, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Estado de mercado actualizado exitosamente"})
}

func (h *MarketLiveHandler) UpdateOdds(c *gin.Context) {
	adminIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no autorizado"})
		return
	}

	adminID, err := uuid.Parse(adminIDStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ID de usuario inválido"})
		return
	}

	marketIDStr := c.Param("id")
	marketID, err := uuid.Parse(marketIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de mercado inválido"})
		return
	}

	var req request.UpdateMarketOddsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payload inválido"})
		return
	}

	parsedOdds := make(map[uuid.UUID]float64)
	for optIDStr, odds := range req.OptionsOdds {
		optID, err := uuid.Parse(optIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de opción inválido: " + optIDStr})
			return
		}
		parsedOdds[optID] = odds
	}

	if err := h.updateOddsUseCase.Execute(c.Request.Context(), marketID, adminID, parsedOdds); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cuotas actualizadas exitosamente"})
}
