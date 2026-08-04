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
}
