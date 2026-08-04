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
	matchRepo := postgres.NewMatchRepository(db)

	// Casos de uso de Liga
	createLeagueUC := usecase.NewCreateLeagueUseCase(leagueRepo)
	joinLeagueUC := usecase.NewJoinLeagueUseCase(leagueRepo)
	getUserLeaguesUC := usecase.NewGetUserLeaguesUseCase(leagueRepo)
	getLeagueDetailsUC := usecase.NewGetLeagueDetailsUseCase(leagueRepo)
	updateLeagueUC := usecase.NewUpdateLeagueUseCase(leagueRepo)
	deleteLeagueUC := usecase.NewDeleteLeagueUseCase(leagueRepo)

	// Casos de uso de Partido y Mercado
	createMatchUC := usecase.NewCreateMatchUseCase(leagueRepo, matchRepo)
	createMarketUC := usecase.NewCreateMarketUseCase(leagueRepo, matchRepo)
	getLeagueMatchesUC := usecase.NewGetLeagueMatchesUseCase(matchRepo)
	getLeagueMarketsUC := usecase.NewGetLeagueMarketsUseCase(matchRepo)
	getMatchMarketsUC := usecase.NewGetMatchMarketsUseCase(matchRepo)

	// Handlers
	leagueHandler := handler.NewLeagueHandler(
		createLeagueUC,
		joinLeagueUC,
		getUserLeaguesUC,
		getLeagueDetailsUC,
		updateLeagueUC,
		deleteLeagueUC,
	)
	matchHandler := handler.NewMatchHandler(
		createMatchUC,
		createMarketUC,
		getLeagueMatchesUC,
		getLeagueMarketsUC,
		getMatchMarketsUC,
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
			leagueRoutes.PUT("/:id", leagueHandler.UpdateLeague)
			leagueRoutes.DELETE("/:id", leagueHandler.DeleteLeague)

			leagueRoutes.POST("/:id/matches", matchHandler.CreateMatch)
			leagueRoutes.GET("/:id/matches", matchHandler.GetMatches)

			leagueRoutes.POST("/:id/markets", matchHandler.CreateMarketForLeague)
			leagueRoutes.GET("/:id/markets", matchHandler.GetLeagueMarkets)
		}

		matchRoutes := api.Group("/matches")
		matchRoutes.Use(authMiddleware)
		{
			matchRoutes.POST("/:id/markets", matchHandler.CreateMarketForMatch)
			matchRoutes.GET("/:id/markets", matchHandler.GetMatchMarkets)
		}
	}
}
