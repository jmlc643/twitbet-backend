package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/jmlc643/twitbet-backend/internal/identity/application/usecase"
	"github.com/jmlc643/twitbet-backend/internal/identity/infrastructure/adapter"
	"github.com/jmlc643/twitbet-backend/internal/identity/infrastructure/http/handler"
	"github.com/jmlc643/twitbet-backend/internal/identity/infrastructure/http/middleware"
	postgresRepo "github.com/jmlc643/twitbet-backend/internal/identity/infrastructure/persistence/postgres"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB, jwtSecret string) {
	userRepo := postgresRepo.NewUserGormRepository(db)
	hasher := adapter.NewBcryptHasher()
	tokenService := adapter.NewJWTTokenService(jwtSecret)

	registerUC := usecase.NewRegisterUseCase(userRepo, hasher, tokenService)
	loginUC := usecase.NewLoginUseCase(userRepo, hasher, tokenService)
	getProfileUC := usecase.NewGetProfileUseCase(userRepo)
	updateProfileUC := usecase.NewUpdateProfileUseCase(userRepo)

	authHandler := handler.NewAuthHandler(registerUC, loginUC, getProfileUC, updateProfileUC)
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
		}
	}
}
