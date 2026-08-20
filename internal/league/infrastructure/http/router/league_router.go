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
	getParticipantMeUC := usecase.NewGetParticipantMeUseCase(leagueRepo)
	rechargeBalanceUC := usecase.NewRechargeBalanceUseCase(leagueRepo)
	grantLeagueBonusUC := usecase.NewGrantLeagueBonusUseCase(leagueRepo)
	getPendingBonusesUC := usecase.NewGetPendingBonusesUseCase(leagueRepo)
	updateLeagueStatusUC := usecase.NewUpdateLeagueStatusUseCase(leagueRepo)
	getLeagueLeaderboardUC := usecase.NewGetLeagueLeaderboardUseCase(leagueRepo)

	// Casos de uso de Partido y Mercado
	createMatchUC := usecase.NewCreateMatchUseCase(leagueRepo, matchRepo)
	createMarketUC := usecase.NewCreateMarketUseCase(leagueRepo, matchRepo, marketPublisher)
	getLeagueMatchesUC := usecase.NewGetLeagueMatchesUseCase(matchRepo)
	getLeagueMarketsUC := usecase.NewGetLeagueMarketsUseCase(matchRepo)
	getMatchMarketsUC := usecase.NewGetMatchMarketsUseCase(matchRepo)
	getMatchDetailsUC := usecase.NewGetMatchDetailsUseCase(matchRepo)
	updateMatchStatusUC := usecase.NewUpdateMatchStatusUseCase(matchRepo, leagueRepo, marketPublisher)
	
	combinedBetRepo := postgres.NewCombinedBetRepository(db)

	updateMarketStatusUC := usecase.NewUpdateMarketStatusUseCase(matchRepo, leagueRepo, marketPublisher)
	updateMarketOddsUC := usecase.NewUpdateMarketOddsUseCase(matchRepo, leagueRepo, marketPublisher)
	resolveMarketUC := usecase.NewResolveMarketUseCase(betRepo, combinedBetRepo, matchRepo, leagueRepo, marketPublisher)
	cancelMarketUC := usecase.NewCancelMarketUseCase(betRepo, matchRepo, leagueRepo, marketPublisher)
	updateMarketOptionStatusUC := usecase.NewUpdateMarketOptionStatusUseCase(matchRepo, leagueRepo, marketPublisher)
	addMarketOptionsUC := usecase.NewAddMarketOptionsUseCase(matchRepo, leagueRepo, marketPublisher)

	marketLiveHandler := handler.NewMarketLiveHandler(*updateMarketStatusUC, *updateMarketOddsUC, *resolveMarketUC, *cancelMarketUC, updateMarketOptionStatusUC, addMarketOptionsUC)

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
		getParticipantMeUC,
		rechargeBalanceUC,
		grantLeagueBonusUC,
		getPendingBonusesUC,
		updateLeagueStatusUC,
		getLeagueLeaderboardUC,
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
	placeBetUC := usecase.NewPlaceBetUseCase(betRepo, leagueRepo, matchRepo, marketPublisher)

	getUserBetsUC := usecase.NewGetUserBetsUseCase(betRepo, leagueRepo)
	cashoutBetUC := usecase.NewCashoutBetUseCase(betRepo, leagueRepo, matchRepo)
	betHandler := handler.NewBetHandler(placeBetUC, getUserBetsUC, cashoutBetUC)

	// Casos de uso de Apuestas Combinadas
	placeCombinedBetUC := usecase.NewPlaceCombinedBetUseCase(combinedBetRepo, leagueRepo, matchRepo, marketPublisher)
	getUserCombinedBetsUC := usecase.NewGetUserCombinedBetsUseCase(combinedBetRepo, leagueRepo, matchRepo)
	getCombinedBetDetailsUC := usecase.NewGetCombinedBetDetailsUseCase(combinedBetRepo, leagueRepo, matchRepo)
	cashoutCombinedBetUC := usecase.NewCashoutCombinedBetUseCase(combinedBetRepo, leagueRepo, matchRepo)
	
	combinedBetHandler := handler.NewCombinedBetHandler(placeCombinedBetUC, getUserCombinedBetsUC, getCombinedBetDetailsUC, cashoutCombinedBetUC)

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
			leagueRoutes.PATCH("/:id/status", leagueHandler.UpdateLeagueStatus)
			leagueRoutes.DELETE("/:id", leagueHandler.DeleteLeague)
			leagueRoutes.POST("/:id/admins", leagueHandler.AssignAdmin)
			leagueRoutes.DELETE("/:id/admins/:participant_id", leagueHandler.RemoveAdmin)

			leagueRoutes.POST("/:id/matches", matchHandler.CreateMatch)
			leagueRoutes.GET("/:id/matches", matchHandler.GetMatches)

			leagueRoutes.POST("/:id/markets", matchHandler.CreateMarketForLeague)
			leagueRoutes.GET("/:id/markets", matchHandler.GetLeagueMarkets)

			leagueRoutes.GET("/:id/leaderboard", leagueHandler.GetLeagueLeaderboard)

			leagueRoutes.GET("/:id/bets", betHandler.GetUserBets)
			leagueRoutes.GET("/:id/me", leagueHandler.GetParticipantMe)

			leagueRoutes.POST("/:id/recharge", leagueHandler.RechargeBalance)
			leagueRoutes.POST("/:id/bonuses", leagueHandler.GrantBonus)
			leagueRoutes.GET("/:id/bonuses/me", leagueHandler.GetMyBonuses)
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
			marketRoutes.POST("/:id/cancel", marketLiveHandler.CancelMarket)
			marketRoutes.POST("/:id/options", marketLiveHandler.AddOptions)
			marketRoutes.PATCH("/:id/options/:option_id/status", marketLiveHandler.UpdateOptionStatus)
		}

		betRoutes := api.Group("/bets")
		betRoutes.Use(authMiddleware)
		{
			betRoutes.POST("", betHandler.PlaceBet)
			betRoutes.POST("/:id/cashout", betHandler.Cashout)
		}

		combinedBetRoutes := api.Group("/combined-bets")
		combinedBetRoutes.Use(authMiddleware)
		{
			combinedBetRoutes.POST("", combinedBetHandler.PlaceCombinedBet)
			combinedBetRoutes.GET("", combinedBetHandler.GetUserCombinedBets)
			combinedBetRoutes.GET("/:id", combinedBetHandler.GetCombinedBetDetails)
			combinedBetRoutes.POST("/:id/cashout", combinedBetHandler.Cashout)
		}
	}
}
