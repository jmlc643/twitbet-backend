package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
)

type MarketEventPublisher interface {
	PublishMarketCreated(ctx context.Context, market entity.Market) error
	PublishMarketStatusChanged(ctx context.Context, marketID uuid.UUID, newStatus string) error
	PublishOddsUpdated(ctx context.Context, marketID uuid.UUID, options []entity.MarketOption) error
	PublishMarketOptionsUpdated(ctx context.Context, market entity.Market) error
	PublishMatchStatusChanged(ctx context.Context, matchID uuid.UUID, newStatus string) error
	PublishMarketResolved(ctx context.Context, marketID uuid.UUID, leagueID uuid.UUID, winningOptionIDs []uuid.UUID) error
	PublishParticipantBalanceUpdated(ctx context.Context, participantID uuid.UUID, leagueID uuid.UUID, userID uuid.UUID) error
}