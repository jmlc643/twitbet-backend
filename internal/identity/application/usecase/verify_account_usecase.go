package usecase

import (
	"context"
	"fmt"

	"github.com/jmlc643/twitbet-backend/internal/identity/application/input"
	"github.com/jmlc643/twitbet-backend/internal/identity/application/output"
	"github.com/jmlc643/twitbet-backend/internal/identity/domain/apperror"
	"github.com/jmlc643/twitbet-backend/internal/identity/domain/port"
	"github.com/jmlc643/twitbet-backend/internal/identity/domain/repository"
)

type VerifyAccountUseCase struct {
	userRepo     repository.UserRepository
	otpRepo      port.OTPRepository
	tokenService port.TokenService
}

func NewVerifyAccountUseCase(
	userRepo repository.UserRepository,
	otpRepo port.OTPRepository,
	tokenService port.TokenService,
) *VerifyAccountUseCase {
	return &VerifyAccountUseCase{
		userRepo:     userRepo,
		otpRepo:      otpRepo,
		tokenService: tokenService,
	}
}

func (uc *VerifyAccountUseCase) Execute(ctx context.Context, in input.VerifyAccountInput) (*output.AuthOutput, error) {
	user, err := uc.userRepo.FindByEmail(ctx, in.Email)
	if err != nil || user == nil {
		return nil, apperror.ErrUserNotFound
	}

	if user.IsVerified {
		return nil, apperror.ErrAlreadyVerified
	}

	otpKey := fmt.Sprintf("verify:%s", user.Email)
	savedOTP, err := uc.otpRepo.GetOTP(ctx, otpKey)
	if err != nil || savedOTP == "" {
		return nil, apperror.ErrOTPExpired
	}

	if savedOTP != in.OTPCode {
		return nil, apperror.ErrOTPInvalid
	}

	user.IsVerified = true
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, apperror.ErrInternal
	}

	_ = uc.otpRepo.DeleteOTP(ctx, otpKey)

	token, err := uc.tokenService.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, apperror.ErrInternal
	}

	return &output.AuthOutput{
		Token: token,
		User: output.UserOutput{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			AvatarURL: user.AvatarURL,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}
