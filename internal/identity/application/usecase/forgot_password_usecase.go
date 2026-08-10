package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/jmlc643/twitbet-backend/internal/identity/application/input"
	"github.com/jmlc643/twitbet-backend/internal/identity/domain/apperror"
	"github.com/jmlc643/twitbet-backend/internal/identity/domain/port"
	"github.com/jmlc643/twitbet-backend/internal/identity/domain/repository"
)

type ForgotPasswordUseCase struct {
	userRepo     repository.UserRepository
	otpRepo      port.OTPRepository
	emailService port.EmailService
}

func NewForgotPasswordUseCase(
	userRepo repository.UserRepository,
	otpRepo port.OTPRepository,
	emailService port.EmailService,
) *ForgotPasswordUseCase {
	return &ForgotPasswordUseCase{
		userRepo:     userRepo,
		otpRepo:      otpRepo,
		emailService: emailService,
	}
}

func (uc *ForgotPasswordUseCase) Execute(ctx context.Context, in input.ForgotPasswordInput) error {
	user, err := uc.userRepo.FindByEmail(ctx, in.Email)
	if err != nil || user == nil {
		return nil
	}

	otpCode := generateOTP()
	otpKey := fmt.Sprintf("reset:%s", user.Email)

	if err := uc.otpRepo.SaveOTP(ctx, otpKey, otpCode, 10*time.Minute); err != nil {
		fmt.Printf("Error saving OTP to Redis: %v\n", err)
		return apperror.ErrInternal
	}

	go func() {
		if err := uc.emailService.SendPasswordResetEmail(context.Background(), user.Email, otpCode); err != nil {
			fmt.Printf("Error sending password reset email to %s: %v\n", user.Email, err)
		}
	}()

	return nil
}
