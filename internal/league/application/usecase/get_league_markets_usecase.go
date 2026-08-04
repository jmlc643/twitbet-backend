package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/repository"
)

type GetLeagueMarketsUseCase struct {
	matchRepo repository.MatchRepository
}

func NewGetLeagueMarketsUseCase(matchRepo repository.MatchRepository) *GetLeagueMarketsUseCase {
	return &GetLeagueMarketsUseCase{
		matchRepo: matchRepo,
	}
}

func (uc *GetLeagueMarketsUseCase) Execute(ctx context.Context, leagueID uuid.UUID) ([]entity.Market, error) {
	return uc.matchRepo.GetMarketsByLeagueID(ctx, leagueID)
}
