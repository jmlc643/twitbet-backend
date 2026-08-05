package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
)

type MatchRepository interface {
	CreateMatch(ctx context.Context, match *entity.Match) error
	CreateMarket(ctx context.Context, market *entity.Market) error
	GetMatchByID(ctx context.Context, id uuid.UUID) (*entity.Match, error)
	GetMatchesByLeagueID(ctx context.Context, leagueID uuid.UUID, limit, offset int, status string) ([]entity.Match, int64, error)
	GetMarketsByLeagueID(ctx context.Context, leagueID uuid.UUID) ([]entity.Market, error)
	GetMarketsByMatchID(ctx context.Context, matchID uuid.UUID) ([]entity.Market, error)
	GetMarketByID(ctx context.Context, id uuid.UUID) (*entity.Market, error)
	UpdateMarket(ctx context.Context, market *entity.Market) error
}
