package router

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/jmlc643/twitbet-backend/internal/league/application/usecase"
	"github.com/jmlc643/twitbet-backend/internal/league/infrastructure/http/handler"
	"github.com/jmlc643/twitbet-backend/internal/league/infrastructure/persistence/postgres"
	leagueRedis "github.com/jmlc643/twitbet-backend/internal/league/infrastructure/redis"

	identityAdapter "github.com/jmlc643/twitbet-backend/internal/identity/infrastructure/adapter"
	identityMiddleware "github.com/jmlc643/twitbet-backend/internal/identity/infrastructure/http/middleware"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB, rdb *redis.Client, jwtSecret string) {
	leagueRepo := postgres.NewLeagueRepository(db)
	matchRepo := postgres.NewMatchRepository(db)
	marketPublisher := leagueRedis.NewMarketPublisher(rdb)
	betRepo := postgres.NewBetRepository(db)

	// Casos de uso de Liga
	createLeagueUC := usecase.NewCreateLeagueUseCase(leagueRepo)
	joinLeagueUC := usecase.NewJoinLeagueUseCase(leagueRepo)
	getUserLeaguesUC := usecase.NewGetUserLeaguesUseCase(leagueRepo)
	getLeagueDetailsUC := usecase.NewGetLeagueDetailsUseCase(leagueRepo)
	updateLeagueUC := usecase.NewUpdateLeagueUseCase(leagueRepo)
	deleteLeagueUC := usecase.NewDeleteLeagueUseCase(leagueRepo)
	assignAdminUC := usecase.NewAssignAdminUseCase(leagueRepo)
	removeAdminUC := usecase.NewRemoveAdminUseCase(leagueRepo)

	// Casos de uso de Partido y Mercado
	createMatchUC := usecase.NewCreateMatchUseCase(leagueRepo, matchRepo)
	createMarketUC := usecase.NewCreateMarketUseCase(leagueRepo, matchRepo)
	getLeagueMatchesUC := usecase.NewGetLeagueMatchesUseCase(matchRepo)
	getLeagueMarketsUC := usecase.NewGetLeagueMarketsUseCase(matchRepo)
	getMatchMarketsUC := usecase.NewGetMatchMarketsUseCase(matchRepo)
	getMatchDetailsUC := usecase.NewGetMatchDetailsUseCase(matchRepo)
	updateMatchStatusUC := usecase.NewUpdateMatchStatusUseCase(matchRepo, leagueRepo, marketPublisher)
	
	updateMarketStatusUC := usecase.NewUpdateMarketStatusUseCase(matchRepo, leagueRepo, marketPublisher)
	updateMarketOddsUC := usecase.NewUpdateMarketOddsUseCase(matchRepo, leagueRepo, marketPublisher)
	resolveMarketUC := usecase.NewResolveMarketUseCase(betRepo, matchRepo, leagueRepo, marketPublisher)

	// Handlers
	leagueHandler := handler.NewLeagueHandler(
		createLeagueUC,
		joinLeagueUC,
		getUserLeaguesUC,
		getLeagueDetailsUC,
		updateLeagueUC,
		deleteLeagueUC,
		assignAdminUC,
		removeAdminUC,
	)
	matchHandler := handler.NewMatchHandler(
		createMatchUC,
		createMarketUC,
		getLeagueMatchesUC,
		getLeagueMarketsUC,
		getMatchMarketsUC,
		getMatchDetailsUC,
		updateMatchStatusUC,
	)
	marketLiveHandler := handler.NewMarketLiveHandler(
		*updateMarketStatusUC,
		*updateMarketOddsUC,
		*resolveMarketUC,
	)

	placeBetUC := usecase.NewPlaceBetUseCase(betRepo, leagueRepo, matchRepo, marketPublisher)
	betHandler := handler.NewBetHandler(placeBetUC)

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
			leagueRoutes.POST("/:id/admins", leagueHandler.AssignAdmin)
			leagueRoutes.DELETE("/:id/admins/:participant_id", leagueHandler.RemoveAdmin)

			leagueRoutes.POST("/:id/matches", matchHandler.CreateMatch)
			leagueRoutes.GET("/:id/matches", matchHandler.GetMatches)

			leagueRoutes.POST("/:id/markets", matchHandler.CreateMarketForLeague)
			leagueRoutes.GET("/:id/markets", matchHandler.GetLeagueMarkets)
		}

		matchRoutes := api.Group("/matches")
		matchRoutes.Use(authMiddleware)
		{
			matchRoutes.GET("/:id", matchHandler.GetMatchDetails)
			matchRoutes.PATCH("/:id/status", matchHandler.UpdateMatchStatus)
			matchRoutes.POST("/:id/markets", matchHandler.CreateMarketForMatch)
			matchRoutes.GET("/:id/markets", matchHandler.GetMatchMarkets)
		}

		marketRoutes := api.Group("/markets")
		marketRoutes.Use(authMiddleware)
		{
			marketRoutes.PATCH("/:id/status", marketLiveHandler.UpdateStatus)
			marketRoutes.PATCH("/:id/odds", marketLiveHandler.UpdateOdds)
			marketRoutes.POST("/:id/resolve", marketLiveHandler.ResolveMarket)
		}

		betRoutes := api.Group("/bets")
		betRoutes.Use(authMiddleware)
		{
			betRoutes.POST("", betHandler.PlaceBet)
		}
	}
}