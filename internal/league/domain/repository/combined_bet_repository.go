package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
)

type CombinedBetRepository interface {
	Create(ctx context.Context, bet *entity.CombinedBet) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.CombinedBet, error)
	GetByParticipantID(ctx context.Context, participantID uuid.UUID) ([]entity.CombinedBet, error)
	GetByLeagueID(ctx context.Context, leagueID uuid.UUID) ([]entity.CombinedBet, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	UpdateCashout(ctx context.Context, id uuid.UUID, cashoutValue float64) error
	ResolveCombinedBetsForMarketAtomic(ctx context.Context, marketID uuid.UUID, winningOptionIDs []uuid.UUID) error
	PlaceCombinedBetAtomic(ctx context.Context, bet *entity.CombinedBet, tx *entity.Transaction) error
}
