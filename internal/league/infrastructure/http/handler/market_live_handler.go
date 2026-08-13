package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/application/input"
	"github.com/jmlc643/twitbet-backend/internal/league/application/usecase"
	"github.com/jmlc643/twitbet-backend/internal/league/infrastructure/http/dto/request"
)

type MarketLiveHandler struct {
	updateStatusUseCase          usecase.UpdateMarketStatusUseCase
	updateOddsUseCase            usecase.UpdateMarketOddsUseCase
	resolveMarketUseCase         usecase.ResolveMarketUseCase
	cancelMarketUseCase          usecase.CancelMarketUseCase
	updateMarketOptionStatusUC   *usecase.UpdateMarketOptionStatusUseCase
	addMarketOptionsUC           *usecase.AddMarketOptionsUseCase
}

func NewMarketLiveHandler(updateStatusUseCase usecase.UpdateMarketStatusUseCase, updateOddsUseCase usecase.UpdateMarketOddsUseCase, resolveMarketUseCase usecase.ResolveMarketUseCase, cancelMarketUseCase usecase.CancelMarketUseCase, updateMarketOptionStatusUC *usecase.UpdateMarketOptionStatusUseCase, addMarketOptionsUC *usecase.AddMarketOptionsUseCase) *MarketLiveHandler {
	return &MarketLiveHandler{
		updateStatusUseCase:        updateStatusUseCase,
		updateOddsUseCase:          updateOddsUseCase,
		resolveMarketUseCase:       resolveMarketUseCase,
		cancelMarketUseCase:        cancelMarketUseCase,
		updateMarketOptionStatusUC: updateMarketOptionStatusUC,
		addMarketOptionsUC:         addMarketOptionsUC,
	}
}

func (h *MarketLiveHandler) UpdateStatus(c *gin.Context) {
	OwnerIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no autorizado"})
		return
	}

	OwnerID, err := uuid.Parse(OwnerIDStr.(string))
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

	if err := h.updateStatusUseCase.Execute(c.Request.Context(), marketID, OwnerID, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Estado de mercado actualizado exitosamente"})
}

func (h *MarketLiveHandler) UpdateOdds(c *gin.Context) {
	OwnerIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no autorizado"})
		return
	}

	OwnerID, err := uuid.Parse(OwnerIDStr.(string))
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

	if err := h.updateOddsUseCase.Execute(c.Request.Context(), marketID, OwnerID, parsedOdds); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cuotas actualizadas exitosamente"})
}

func (h *MarketLiveHandler) ResolveMarket(c *gin.Context) {
	marketIDStr := c.Param("id")
	marketID, err := uuid.Parse(marketIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de mercado inválido"})
		return
	}

	var req request.ResolveMarketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payload inválido"})
		return
	}

	winningOptionIDs := make([]uuid.UUID, 0, len(req.WinningOptionIDs))
	for _, idStr := range req.WinningOptionIDs {
		optID, err := uuid.Parse(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de opción ganadora inválido: " + idStr})
			return
		}
		winningOptionIDs = append(winningOptionIDs, optID)
	}
	if len(winningOptionIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Se requiere al menos una opción ganadora"})
		return
	}

	if err := h.resolveMarketUseCase.Execute(c.Request.Context(), marketID, winningOptionIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mercado resuelto exitosamente"})
}

func (h *MarketLiveHandler) CancelMarket(c *gin.Context) {
	marketIDStr := c.Param("id")
	marketID, err := uuid.Parse(marketIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de mercado inválido"})
		return
	}

	var req request.CancelMarketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payload inválido o cancellation_reason faltante"})
		return
	}

	if err := h.cancelMarketUseCase.Execute(c.Request.Context(), marketID, req.CancellationReason); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mercado anulado exitosamente y saldo reembolsado"})
}

func (h *MarketLiveHandler) UpdateOptionStatus(c *gin.Context) {
	ownerIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no autorizado"})
		return
	}

	ownerID, err := uuid.Parse(ownerIDStr.(string))
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

	optionIDStr := c.Param("option_id")
	optionID, err := uuid.Parse(optionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de opción inválido"})
		return
	}

	var req request.UpdateMarketOptionStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payload inválido"})
		return
	}

	if err := h.updateMarketOptionStatusUC.Execute(c.Request.Context(), marketID, optionID, ownerID, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Estado de opción actualizado exitosamente"})
}

func (h *MarketLiveHandler) AddOptions(c *gin.Context) {
	ownerIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no autorizado"})
		return
	}

	ownerID, err := uuid.Parse(ownerIDStr.(string))
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

	var req request.AddMarketOptionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payload inválido"})
		return
	}

	var inputOptions []input.MarketOptionInput
	for _, opt := range req.Options {
		inputOptions = append(inputOptions, input.MarketOptionInput{Name: opt.Name, Odds: opt.Odds})
	}

	if err := h.addMarketOptionsUC.Execute(c.Request.Context(), marketID, ownerID, inputOptions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Opciones agregadas exitosamente"})
}