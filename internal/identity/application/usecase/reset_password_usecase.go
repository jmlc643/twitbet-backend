package usecase

import (
	"context"
	"fmt"

	"github.com/jmlc643/twitbet-backend/internal/identity/application/input"
	"github.com/jmlc643/twitbet-backend/internal/identity/domain/apperror"
	"github.com/jmlc643/twitbet-backend/internal/identity/domain/port"
	"github.com/jmlc643/twitbet-backend/internal/identity/domain/repository"
)

type ResetPasswordUseCase struct {
	userRepo repository.UserRepository
	otpRepo  port.OTPRepository
	hasher   port.PasswordHasher
}

func NewResetPasswordUseCase(
	userRepo repository.UserRepository,
	otpRepo port.OTPRepository,
	hasher port.PasswordHasher,
) *ResetPasswordUseCase {
	return &ResetPasswordUseCase{
		userRepo: userRepo,
		otpRepo:  otpRepo,
		hasher:   hasher,
	}
}

func (uc *ResetPasswordUseCase) Execute(ctx context.Context, in input.ResetPasswordInput) error {
	user, err := uc.userRepo.FindByEmail(ctx, in.Email)
	if err != nil || user == nil {
		return apperror.ErrUserNotFound
	}

	if in.NewPassword != in.ConfirmPassword {
		return apperror.ErrPasswordsDoNotMatch
	}

	otpKey := fmt.Sprintf("reset:%s", user.Email)
	savedOTP, err := uc.otpRepo.GetOTP(ctx, otpKey)
	if err != nil || savedOTP == "" {
		return apperror.ErrOTPExpired
	}

	if savedOTP != in.OTPCode {
		return apperror.ErrOTPInvalid
	}

	hashedPassword, err := uc.hasher.HashPassword(in.NewPassword)
	if err != nil {
		return apperror.ErrInternal
	}

	user.PasswordHash = hashedPassword
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return apperror.ErrInternal
	}

	_ = uc.otpRepo.DeleteOTP(ctx, otpKey)

	return nil
}
