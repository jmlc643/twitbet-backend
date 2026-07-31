package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/jmlc643/twitbet-backend/internal/league/application/usecase"
	"github.com/jmlc643/twitbet-backend/internal/league/infrastructure/http/handler"
	"github.com/jmlc643/twitbet-backend/internal/league/infrastructure/persistence/postgres"

	identityAdapter "github.com/jmlc643/twitbet-backend/internal/identity/infrastructure/adapter"
	identityMiddleware "github.com/jmlc643/twitbet-backend/internal/identity/infrastructure/http/middleware"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB, jwtSecret string) {
	leagueRepo := postgres.NewLeagueRepository(db)

	// Casos de uso
	createLeagueUC := usecase.NewCreateLeagueUseCase(leagueRepo)
	joinLeagueUC := usecase.NewJoinLeagueUseCase(leagueRepo)
	getUserLeaguesUC := usecase.NewGetUserLeaguesUseCase(leagueRepo)
	getLeagueDetailsUC := usecase.NewGetLeagueDetailsUseCase(leagueRepo)
	updateLeagueSettingsUC := usecase.NewUpdateLeagueSettingsUseCase(leagueRepo)

	// Handlers
	leagueHandler := handler.NewLeagueHandler(
		createLeagueUC,
		joinLeagueUC,
		getUserLeaguesUC,
		getLeagueDetailsUC,
		updateLeagueSettingsUC,
	)

	tokenService := identityAdapter.NewJWTTokenService(jwtSecret)
	authMiddleware := identityMiddleware.JWTMiddleware(tokenService)

	api := router.Group("/api/v1")
	{
		leagueRoutes := api.Group("/leagues")
		leagueRoutes.Use(authMiddleware)
		{
			leagueRoutes.GET("", leagueHandler.GetUserLeagues)
			leagueRoutes.POST("", leagueHandler.CreateLeague)
			leagueRoutes.POST("/join", leagueHandler.JoinLeague)
			leagueRoutes.GET("/:id", leagueHandler.GetLeagueDetails)
			leagueRoutes.PATCH("/:id/settings", leagueHandler.UpdateLeagueSettings)
		}
	}
}
