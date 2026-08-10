package router

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/jmlc643/twitbet-backend/internal/identity/application/usecase"
	"github.com/jmlc643/twitbet-backend/internal/identity/domain/port"
	"github.com/jmlc643/twitbet-backend/internal/identity/infrastructure/adapter"
	"github.com/jmlc643/twitbet-backend/internal/identity/infrastructure/http/handler"
	"github.com/jmlc643/twitbet-backend/internal/identity/infrastructure/http/middleware"
	postgresRepo "github.com/jmlc643/twitbet-backend/internal/identity/infrastructure/persistence/postgres"
	"github.com/jmlc643/twitbet-backend/internal/platform/config"
	"github.com/jmlc643/twitbet-backend/internal/platform/email"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB, rdb *redis.Client, cfg *config.Config) {
	userRepo := postgresRepo.NewUserGormRepository(db)
	hasher := adapter.NewBcryptHasher()
	tokenService := adapter.NewJWTTokenService(cfg.JWTSecret)
	otpRepo := adapter.NewRedisOTPRepository(rdb)

	emailService, err := email.NewResendService(cfg, "internal/platform/email/template")
	if err != nil {
		log.Fatalf("Falló la inicialización del servicio de email (Resend): %v", err)
	}

	var storageService port.StorageService
	if cfg.CloudinaryURL != "" {
		cloudinaryStorage, err := adapter.NewCloudinaryStorage(cfg.CloudinaryURL)
		if err == nil {
			storageService = cloudinaryStorage
		}
	}

	registerUC := usecase.NewRegisterUseCase(userRepo, hasher, tokenService, otpRepo, emailService)
	loginUC := usecase.NewLoginUseCase(userRepo, hasher, tokenService)
	getProfileUC := usecase.NewGetProfileUseCase(userRepo)
	updateProfileUC := usecase.NewUpdateProfileUseCase(userRepo, storageService)
	uploadAvatarUC := usecase.NewUploadAvatarUseCase(userRepo, storageService)
	verifyAccountUC := usecase.NewVerifyAccountUseCase(userRepo, otpRepo, tokenService)
	forgotPasswordUC := usecase.NewForgotPasswordUseCase(userRepo, otpRepo, emailService)
	verifyResetOtpUC := usecase.NewVerifyResetOtpUseCase(otpRepo)
	resetPasswordUC := usecase.NewResetPasswordUseCase(userRepo, otpRepo, hasher)
	changePasswordUC := usecase.NewChangePasswordUseCase(userRepo, hasher)

	authHandler := handler.NewAuthHandler(
		registerUC, loginUC, getProfileUC, updateProfileUC, uploadAvatarUC,
		verifyAccountUC, forgotPasswordUC, verifyResetOtpUC, resetPasswordUC, changePasswordUC,
	)
	jwtMiddleware := middleware.JWTMiddleware(tokenService)

	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/verify-account", authHandler.VerifyAccount)
			auth.POST("/forgot-password", authHandler.ForgotPassword)
			auth.POST("/verify-reset-otp", authHandler.VerifyResetOtp)
			auth.POST("/reset-password", authHandler.ResetPassword)
		}

		users := api.Group("/users")
		users.Use(jwtMiddleware)
		{
			users.GET("/me", authHandler.GetProfile)
			users.PUT("/me", authHandler.UpdateProfile)
			users.POST("/me/avatar", authHandler.UploadAvatar)
			users.POST("/me/change-password", authHandler.ChangePassword)
		}
	}
}
