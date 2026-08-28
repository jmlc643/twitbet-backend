package port

import (
	"context"
	"time"

	"github.com/jmlc643/twitbet-backend/internal/identity/domain/entity"
)

type UserStatsCache interface {
	GetStats(ctx context.Context, userID string) (*entity.UserStats, error)
	SetStats(ctx context.Context, userID string, stats *entity.UserStats, expiration time.Duration) error
}
