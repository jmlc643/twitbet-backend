package usecase

import (
	"context"

	"github.com/jmlc643/twitbet-backend/internal/identity/application/input"
	"github.com/jmlc643/twitbet-backend/internal/identity/domain/apperror"
	"github.com/jmlc643/twitbet-backend/internal/identity/domain/port"
	"github.com/jmlc643/twitbet-backend/internal/identity/domain/repository"
)

type ChangePasswordUseCase struct {
	userRepo repository.UserRepository
	hasher   port.PasswordHasher
}

func NewChangePasswordUseCase(
	userRepo repository.UserRepository,
	hasher port.PasswordHasher,
) *ChangePasswordUseCase {
	return &ChangePasswordUseCase{
		userRepo: userRepo,
		hasher:   hasher,
	}
}

func (uc *ChangePasswordUseCase) Execute(ctx context.Context, in input.ChangePasswordInput) error {
	user, err := uc.userRepo.FindByID(ctx, in.UserID)
	if err != nil || user == nil {
		return apperror.ErrUserNotFound
	}

	if err := uc.hasher.ComparePassword(in.OldPassword, user.PasswordHash); err != nil {
		return apperror.ErrInvalidOldPassword
	}

	if in.NewPassword != in.ConfirmPassword {
		return apperror.ErrPasswordsDoNotMatch
	}

	hashedPassword, err := uc.hasher.HashPassword(in.NewPassword)
	if err != nil {
		return apperror.ErrInternal
	}

	user.PasswordHash = hashedPassword
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return apperror.ErrInternal
	}

	return nil
}
