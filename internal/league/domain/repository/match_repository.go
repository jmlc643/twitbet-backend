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
	GetMatchBySlug(ctx context.Context, slug string) (*entity.Match, error)
	GetMatchesByLeagueID(ctx context.Context, leagueID uuid.UUID, limit, offset int, status string) ([]entity.Match, int64, error)
	GetMarketsByLeagueID(ctx context.Context, leagueID uuid.UUID) ([]entity.Market, error)
	GetMarketsByMatchID(ctx context.Context, matchID uuid.UUID) ([]entity.Market, error)
	GetMarketsByIDs(ctx context.Context, marketIDs []uuid.UUID) ([]entity.Market, error)
	GetMarketByID(ctx context.Context, id uuid.UUID) (*entity.Market, error)
	UpdateMarket(ctx context.Context, market *entity.Market) error
	UpdateMarketAndHistory(ctx context.Context, market *entity.Market, history []entity.MarketOddsHistory) error
	AddMarketOptions(ctx context.Context, marketID uuid.UUID, options []entity.MarketOption) error
	UpdateMarketOptionStatus(ctx context.Context, marketID uuid.UUID, optionID uuid.UUID, newStatus string) error
	UpdateMatchStatusAtomic(ctx context.Context, matchID uuid.UUID, newStatus string) error
	SaveMarketOddsHistory(ctx context.Context, history []entity.MarketOddsHistory) error
}
