package usecase

import (
	"context"
	"time"

	"github.com/jmlc643/twitbet-backend/internal/identity/application/output"
	"github.com/jmlc643/twitbet-backend/internal/identity/domain/port"
	"github.com/jmlc643/twitbet-backend/internal/identity/domain/repository"
)

func fetchUserStats(ctx context.Context, userID string, userRepo repository.UserRepository, statsCache port.UserStatsCache) output.UserStatsOutput {
	var statsOutput output.UserStatsOutput

	cachedStats, err := statsCache.GetStats(ctx, userID)
	if err == nil && cachedStats != nil {
		statsOutput = output.UserStatsOutput{
			Leagues:       cachedStats.Leagues,
			Wins:          cachedStats.Wins,
			Effectiveness: cachedStats.Effectiveness,
		}
	} else {
		stats, err := userRepo.GetUserStats(ctx, userID)
		if err == nil && stats != nil {
			statsOutput = output.UserStatsOutput{
				Leagues:       stats.Leagues,
				Wins:          stats.Wins,
				Effectiveness: stats.Effectiveness,
			}
			_ = statsCache.SetStats(ctx, userID, stats, 15*time.Minute)
		}
	}

	return statsOutput
}
