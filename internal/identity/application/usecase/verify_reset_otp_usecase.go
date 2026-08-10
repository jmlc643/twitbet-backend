package usecase

import (
	"context"
	"fmt"

	"github.com/jmlc643/twitbet-backend/internal/identity/application/input"
	"github.com/jmlc643/twitbet-backend/internal/identity/domain/apperror"
	"github.com/jmlc643/twitbet-backend/internal/identity/domain/port"
)

type VerifyResetOtpUseCase struct {
	otpRepo port.OTPRepository
}

func NewVerifyResetOtpUseCase(otpRepo port.OTPRepository) *VerifyResetOtpUseCase {
	return &VerifyResetOtpUseCase{
		otpRepo: otpRepo,
	}
}

func (uc *VerifyResetOtpUseCase) Execute(ctx context.Context, in input.VerifyResetOtpInput) error {
	otpKey := fmt.Sprintf("reset:%s", in.Email)
	savedOTP, err := uc.otpRepo.GetOTP(ctx, otpKey)
	if err != nil || savedOTP == "" {
		return apperror.ErrOTPExpired
	}

	if savedOTP != in.OTPCode {
		return apperror.ErrOTPInvalid
	}

	return nil
}
