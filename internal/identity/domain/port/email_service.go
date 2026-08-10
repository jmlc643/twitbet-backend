package port

import "context"

type EmailService interface {
	SendVerificationEmail(ctx context.Context, toEmail, otpCode string) error
	SendPasswordResetEmail(ctx context.Context, toEmail, otpCode string) error
}