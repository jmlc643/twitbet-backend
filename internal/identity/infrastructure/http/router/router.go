package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/jmlc643/twitbet-backend/internal/identity/application/usecase"
	"github.com/jmlc643/twitbet-backend/internal/identity/domain/port"
	"github.com/jmlc643/twitbet-backend/internal/identity/infrastructure/adapter"
	"github.com/jmlc643/twitbet-backend/internal/identity/infrastructure/http/handler"
	"github.com/jmlc643/twitbet-backend/internal/identity/infrastructure/http/middleware"
	postgresRepo "github.com/jmlc643/twitbet-backend/internal/identity/infrastructure/persistence/postgres"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB, jwtSecret string, cloudinaryURL string) {
	userRepo := postgresRepo.NewUserGormRepository(db)
	hasher := adapter.NewBcryptHasher()
	tokenService := adapter.NewJWTTokenService(jwtSecret)

	var storageService port.StorageService
	if cloudinaryURL != "" {
		cloudinaryStorage, err := adapter.NewCloudinaryStorage(cloudinaryURL)
		if err == nil {
			storageService = cloudinaryStorage
		}
	}

	registerUC := usecase.NewRegisterUseCase(userRepo, hasher, tokenService)
	loginUC := usecase.NewLoginUseCase(userRepo, hasher, tokenService)
	getProfileUC := usecase.NewGetProfileUseCase(userRepo)
	updateProfileUC := usecase.NewUpdateProfileUseCase(userRepo, storageService)
	uploadAvatarUC := usecase.NewUploadAvatarUseCase(userRepo, storageService)

	authHandler := handler.NewAuthHandler(registerUC, loginUC, getProfileUC, updateProfileUC, uploadAvatarUC)
	jwtMiddleware := middleware.JWTMiddleware(tokenService)

	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		users := api.Group("/users")
		users.Use(jwtMiddleware)
		{
			users.GET("/me", authHandler.GetProfile)
			users.PUT("/me", authHandler.UpdateProfile)
			users.POST("/me/avatar", authHandler.UploadAvatar)
		}
	}
}
