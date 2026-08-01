package usecase

import (
	"context"

	"github.com/jmlc643/twitbet-backend/internal/identity/application/input"
	"github.com/jmlc643/twitbet-backend/internal/identity/domain/apperror"
	"github.com/jmlc643/twitbet-backend/internal/identity/domain/repository"
)

type UpdateProfileUseCase struct {
	userRepo repository.UserRepository
}

func NewUpdateProfileUseCase(userRepo repository.UserRepository) *UpdateProfileUseCase {
	return &UpdateProfileUseCase{
		userRepo: userRepo,
	}
}

func (uc *UpdateProfileUseCase) Execute(ctx context.Context, in input.UpdateProfileInput) error {
	user, err := uc.userRepo.FindByID(ctx, in.UserID)
	if err != nil {
		return err
	}
	if user == nil {
		return apperror.ErrUserNotFound
	}

	if in.Username != nil {
		user.Username = *in.Username
	}
	
	if in.AvatarURL != nil {
		user.AvatarURL = *in.AvatarURL
	}

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}
