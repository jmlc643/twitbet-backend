package usecase

import (
	"context"

	"github.com/jmlc643/twitbet-backend/internal/identity/application/input"
	"github.com/jmlc643/twitbet-backend/internal/identity/domain/apperror"
	"github.com/jmlc643/twitbet-backend/internal/identity/domain/port"
	"github.com/jmlc643/twitbet-backend/internal/identity/domain/repository"
)

type UpdateProfileUseCase struct {
	userRepo       repository.UserRepository
	storageService port.StorageService
}

func NewUpdateProfileUseCase(userRepo repository.UserRepository, storageService port.StorageService) *UpdateProfileUseCase {
	return &UpdateProfileUseCase{
		userRepo:       userRepo,
		storageService: storageService,
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
		if user.AvatarURL != "" && user.AvatarURL != *in.AvatarURL {
			if uc.storageService != nil {
				_ = uc.storageService.DeleteImage(ctx, user.AvatarURL)
			}
		}
		user.AvatarURL = *in.AvatarURL
	}

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}
