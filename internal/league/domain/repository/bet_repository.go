package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
)

type BetRepository interface {
	PlaceBetAtomic(ctx context.Context, bet *entity.Bet, transaction *entity.Transaction) error
	CashoutAtomic(ctx context.Context, bet *entity.Bet, transaction *entity.Transaction) error
	GetBetByID(ctx context.Context, id uuid.UUID) (*entity.Bet, error)
	UpdateBetStatus(ctx context.Context, id uuid.UUID, status entity.BetStatus) error
	ResolveMarketAtomic(ctx context.Context, marketID uuid.UUID, winningOptionID uuid.UUID) error
	GetBetsByParticipantID(ctx context.Context, participantID uuid.UUID, status *entity.BetStatus, startDate, endDate *time.Time, limit, offset int) ([]entity.BetDetail, int64, error)
}
