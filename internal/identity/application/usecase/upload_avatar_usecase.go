package usecase

import (
	"context"
	"errors"
	"io"

	"github.com/jmlc643/twitbet-backend/internal/identity/domain/apperror"
	"github.com/jmlc643/twitbet-backend/internal/identity/domain/port"
	"github.com/jmlc643/twitbet-backend/internal/identity/domain/repository"
)

type UploadAvatarUseCase struct {
	userRepo       repository.UserRepository
	storageService port.StorageService
}

func NewUploadAvatarUseCase(userRepo repository.UserRepository, storageService port.StorageService) *UploadAvatarUseCase {
	return &UploadAvatarUseCase{
		userRepo:       userRepo,
		storageService: storageService,
	}
}

func (uc *UploadAvatarUseCase) Execute(ctx context.Context, userID string, file io.Reader, filename string) (string, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", apperror.ErrUserNotFound
	}

	if uc.storageService == nil {
		return "", errors.New("storage service not configured")
	}

	newAvatarURL, err := uc.storageService.UploadImage(ctx, file, filename)
	if err != nil {
		return "", err
	}

	if user.AvatarURL != "" {
		_ = uc.storageService.DeleteImage(ctx, user.AvatarURL)
	}

	user.AvatarURL = newAvatarURL
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return "", err
	}

	return newAvatarURL, nil
}
