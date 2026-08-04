package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/repository"
)

type GetMatchMarketsUseCase struct {
	matchRepo repository.MatchRepository
}

func NewGetMatchMarketsUseCase(matchRepo repository.MatchRepository) *GetMatchMarketsUseCase {
	return &GetMatchMarketsUseCase{
		matchRepo: matchRepo,
	}
}

func (uc *GetMatchMarketsUseCase) Execute(ctx context.Context, matchID uuid.UUID) ([]entity.Market, error) {
	return uc.matchRepo.GetMarketsByMatchID(ctx, matchID)
}
