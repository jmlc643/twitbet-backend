package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmlc643/twitbet-backend/internal/identity/application/usecase"
	"github.com/jmlc643/twitbet-backend/internal/identity/domain/apperror"
	"github.com/jmlc643/twitbet-backend/internal/identity/infrastructure/http/dto/request"
	"github.com/jmlc643/twitbet-backend/internal/identity/infrastructure/http/mapper"
)

type AuthHandler struct {
	registerUC      *usecase.RegisterUseCase
	loginUC         *usecase.LoginUseCase
	getProfileUC    *usecase.GetProfileUseCase
	updateProfileUC *usecase.UpdateProfileUseCase
}

func NewAuthHandler(
	registerUC *usecase.RegisterUseCase,
	loginUC *usecase.LoginUseCase,
	getProfileUC *usecase.GetProfileUseCase,
	updateProfileUC *usecase.UpdateProfileUseCase,
) *AuthHandler {
	return &AuthHandler{
		registerUC:      registerUC,
		loginUC:         loginUC,
		getProfileUC:    getProfileUC,
		updateProfileUC: updateProfileUC,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := mapper.RegisterRequestToInput(req)
	output, err := h.registerUC.Execute(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, apperror.ErrUserAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, mapper.AuthOutputToResponse(output))
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := mapper.LoginRequestToInput(req)
	output, err := h.loginUC.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mapper.AuthOutputToResponse(output))
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no autorizado"})
		return
	}

	output, err := h.getProfileUC.Execute(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mapper.UserOutputToResponse(output))
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no autorizado"})
		return
	}

	var req request.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := mapper.UpdateProfileRequestToInput(req, userID.(string))
	err := h.updateProfileUC.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

