package usecase

import (
	"context"
	"time"

	"github.com/jmlc643/twitbet-backend/internal/identity/application/output"
	"github.com/jmlc643/twitbet-backend/internal/identity/domain/apperror"
	"github.com/jmlc643/twitbet-backend/internal/identity/domain/port"
	"github.com/jmlc643/twitbet-backend/internal/identity/domain/repository"
)

type GetProfileUseCase struct {
	userRepo  repository.UserRepository
	statsCache port.UserStatsCache
}

func NewGetProfileUseCase(userRepo repository.UserRepository, statsCache port.UserStatsCache) *GetProfileUseCase {
	return &GetProfileUseCase{
		userRepo:  userRepo,
		statsCache: statsCache,
	}
}

func (uc *GetProfileUseCase) Execute(ctx context.Context, userID string) (*output.UserOutput, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, apperror.ErrUserNotFound
	}

	var statsOutput output.UserStatsOutput

	cachedStats, err := uc.statsCache.GetStats(ctx, userID)
	if err == nil && cachedStats != nil {
		statsOutput = output.UserStatsOutput{
			Leagues:       cachedStats.Leagues,
			Wins:          cachedStats.Wins,
			Effectiveness: cachedStats.Effectiveness,
		}
	} else {
		stats, err := uc.userRepo.GetUserStats(ctx, userID)
		if err == nil && stats != nil {
			statsOutput = output.UserStatsOutput{
				Leagues:       stats.Leagues,
				Wins:          stats.Wins,
				Effectiveness: stats.Effectiveness,
			}
			
			_ = uc.statsCache.SetStats(ctx, userID, stats, 15*time.Minute)
		}
	}

	return &output.UserOutput{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		AvatarURL: &user.AvatarURL,
		CreatedAt: user.CreatedAt,
		Stats:     statsOutput,
	}, nil
}
