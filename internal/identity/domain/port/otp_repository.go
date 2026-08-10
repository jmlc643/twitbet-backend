package port

import (
	"context"
	"time"
)

type OTPRepository interface {
	SaveOTP(ctx context.Context, key, otpCode string, ttl time.Duration) error
	GetOTP(ctx context.Context, key string) (string, error)
	DeleteOTP(ctx context.Context, key string) error
}
