package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
)

type MarketEventPublisher interface {
	PublishMarketStatusChanged(ctx context.Context, marketID uuid.UUID, newStatus string) error
	PublishOddsUpdated(ctx context.Context, marketID uuid.UUID, options []entity.MarketOption) error
	PublishMatchStatusChanged(ctx context.Context, matchID uuid.UUID, newStatus string) error
}
