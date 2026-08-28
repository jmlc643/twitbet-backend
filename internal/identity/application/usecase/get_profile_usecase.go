package usecase

import (
	"context"

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

	statsOutput := fetchUserStats(ctx, userID, uc.userRepo, uc.statsCache)

	return &output.UserOutput{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		AvatarURL: &user.AvatarURL,
		CreatedAt: user.CreatedAt,
		Stats:     statsOutput,
	}, nil
}