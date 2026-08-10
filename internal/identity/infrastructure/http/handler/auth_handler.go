package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jmlc643/twitbet-backend/internal/identity/application/usecase"
	"github.com/jmlc643/twitbet-backend/internal/identity/domain/apperror"
	"github.com/jmlc643/twitbet-backend/internal/identity/infrastructure/http/dto/request"
	"github.com/jmlc643/twitbet-backend/internal/identity/infrastructure/http/mapper"
)

type AuthHandler struct {
	registerUC           *usecase.RegisterUseCase
	loginUC              *usecase.LoginUseCase
	getProfileUC         *usecase.GetProfileUseCase
	updateProfileUC      *usecase.UpdateProfileUseCase
	uploadAvatarUC       *usecase.UploadAvatarUseCase
	verifyAccountUC      *usecase.VerifyAccountUseCase
	forgotPasswordUC     *usecase.ForgotPasswordUseCase
	verifyResetOtpUC     *usecase.VerifyResetOtpUseCase
	resetPasswordUC      *usecase.ResetPasswordUseCase
	changePasswordUC     *usecase.ChangePasswordUseCase
}

func NewAuthHandler(
	registerUC *usecase.RegisterUseCase,
	loginUC *usecase.LoginUseCase,
	getProfileUC *usecase.GetProfileUseCase,
	updateProfileUC *usecase.UpdateProfileUseCase,
	uploadAvatarUC *usecase.UploadAvatarUseCase,
	verifyAccountUC *usecase.VerifyAccountUseCase,
	forgotPasswordUC *usecase.ForgotPasswordUseCase,
	verifyResetOtpUC *usecase.VerifyResetOtpUseCase,
	resetPasswordUC *usecase.ResetPasswordUseCase,
	changePasswordUC *usecase.ChangePasswordUseCase,
) *AuthHandler {
	return &AuthHandler{
		registerUC:           registerUC,
		loginUC:              loginUC,
		getProfileUC:         getProfileUC,
		updateProfileUC:      updateProfileUC,
		uploadAvatarUC:       uploadAvatarUC,
		verifyAccountUC:      verifyAccountUC,
		forgotPasswordUC:     forgotPasswordUC,
		verifyResetOtpUC:     verifyResetOtpUC,
		resetPasswordUC:      resetPasswordUC,
		changePasswordUC:     changePasswordUC,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": formatValidationError(err)})
		return
	}

	input := mapper.RegisterRequestToInput(req)
	err := h.registerUC.Execute(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, apperror.ErrUserAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Cuenta creada. Se ha enviado un código de verificación a tu correo."})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": formatValidationError(err)})
		return
	}

	input := mapper.LoginRequestToInput(req)
	output, err := h.loginUC.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": formatValidationError(err)})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": formatValidationError(err)})
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

func (h *AuthHandler) UploadAvatar(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no autorizado"})
		return
	}

	file, header, err := c.Request.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "se requiere un archivo de imagen en el campo 'avatar'"})
		return
	}
	defer file.Close()

	avatarURL, err := h.uploadAvatarUC.Execute(c.Request.Context(), userID.(string), file, header.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"avatar_url": avatarURL})
}

func (h *AuthHandler) VerifyAccount(c *gin.Context) {
	var req request.VerifyAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": formatValidationError(err)})
		return
	}

	input := mapper.VerifyAccountRequestToInput(req)
	output, err := h.verifyAccountUC.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": formatValidationError(err)})
		return
	}

	c.JSON(http.StatusOK, mapper.AuthOutputToResponse(output))
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req request.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": formatValidationError(err)})
		return
	}

	input := mapper.ForgotPasswordRequestToInput(req)
	if err := h.forgotPasswordUC.Execute(c.Request.Context(), input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Si el correo está registrado, se ha enviado un código de recuperación."})
}

func (h *AuthHandler) VerifyResetOtp(c *gin.Context) {
	var req request.VerifyResetOtpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": formatValidationError(err)})
		return
	}

	input := mapper.VerifyResetOtpRequestToInput(req)
	if err := h.verifyResetOtpUC.Execute(c.Request.Context(), input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Código verificado correctamente."})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req request.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": formatValidationError(err)})
		return
	}

	input := mapper.ResetPasswordRequestToInput(req)
	if err := h.resetPasswordUC.Execute(c.Request.Context(), input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": formatValidationError(err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contraseña actualizada exitosamente."})
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no autorizado"})
		return
	}

	var req request.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": formatValidationError(err)})
		return
	}

	input := mapper.ChangePasswordRequestToInput(req, userID.(string))
	if err := h.changePasswordUC.Execute(c.Request.Context(), input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": formatValidationError(err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contraseña cambiada exitosamente."})
}

func formatValidationError(err error) string {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		for _, fe := range ve {
			switch fe.Tag() {
			case "eqfield":
				return "Las contraseñas no coinciden"
			case "required":
				return "El campo " + fe.Field() + " es obligatorio"
			case "email":
				return "El correo electrónico no es válido"
			case "min":
				return "El campo " + fe.Field() + " debe tener al menos " + fe.Param() + " caracteres"
			case "len":
				return "El campo " + fe.Field() + " debe tener exactamente " + fe.Param() + " caracteres"
			}
		}
		return "Error de validación en los datos enviados"
	}
	return err.Error()
}