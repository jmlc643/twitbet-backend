package usecase

import (
	"context"

	"github.com/jmlc643/twitbet-backend/internal/identity/application/output"
	"github.com/jmlc643/twitbet-backend/internal/identity/domain/apperror"
	"github.com/jmlc643/twitbet-backend/internal/identity/domain/repository"
)

type GetProfileUseCase struct {
	userRepo repository.UserRepository
}

func NewGetProfileUseCase(userRepo repository.UserRepository) *GetProfileUseCase {
	return &GetProfileUseCase{
		userRepo: userRepo,
	}
}

func (uc *GetProfileUseCase) Execute(ctx context.Context, userID string) (*output.UserOutput, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, apperror.ErrUserNotFound
	}

	return &output.UserOutput{
		ID: 		user.ID,
		Username: 	user.Username,
		Email: 		user.Email,
		AvatarURL: 	user.AvatarURL,
		CreatedAt: 	user.CreatedAt,
	}, nil
}