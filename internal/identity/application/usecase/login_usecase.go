package usecase

import (
	"context"

	"github.com/jmlc643/twitbet-backend/internal/identity/application/input"
	"github.com/jmlc643/twitbet-backend/internal/identity/application/output"
	"github.com/jmlc643/twitbet-backend/internal/identity/domain/apperror"

	"github.com/jmlc643/twitbet-backend/internal/identity/domain/port"
	"github.com/jmlc643/twitbet-backend/internal/identity/domain/repository"
)

type LoginUseCase struct {
	userRepo     repository.UserRepository
	hasher       port.PasswordHasher
	tokenService port.TokenService
	statsCache   port.UserStatsCache
}

func NewLoginUseCase(
	userRepo repository.UserRepository,
	hasher port.PasswordHasher,
	tokenService port.TokenService,
	statsCache port.UserStatsCache,
) *LoginUseCase {
	return &LoginUseCase{
		userRepo:     userRepo,
		hasher:       hasher,
		tokenService: tokenService,
		statsCache:   statsCache,
	}
}

func (uc *LoginUseCase) Execute(ctx context.Context, in input.LoginInput) (*output.AuthOutput, error) {
	user, err := uc.userRepo.FindByEmail(ctx, in.Email)

	if err != nil || user == nil {
		return nil, apperror.ErrInvalidCredentials
	}

	if err := uc.hasher.ComparePassword(in.Password, user.PasswordHash); err != nil {
		return nil, apperror.ErrInvalidCredentials
	}

	if !user.IsVerified {
		return nil, apperror.ErrAccountNotVerified
	}

	token, err := uc.tokenService.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, apperror.ErrInternal
	}

	statsOutput := fetchUserStats(ctx, user.ID, uc.userRepo, uc.statsCache)

	return &output.AuthOutput{
		Token: token,
		User: output.UserOutput{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			AvatarURL: &user.AvatarURL,
			CreatedAt: user.CreatedAt,
			Stats:     statsOutput,
		},
	}, nil
}