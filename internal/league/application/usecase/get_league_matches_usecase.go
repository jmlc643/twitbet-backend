package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/repository"
)

type GetLeagueMatchesUseCase struct {
	matchRepo repository.MatchRepository
}

func NewGetLeagueMatchesUseCase(matchRepo repository.MatchRepository) *GetLeagueMatchesUseCase {
	return &GetLeagueMatchesUseCase{
		matchRepo: matchRepo,
	}
}

func (uc *GetLeagueMatchesUseCase) Execute(ctx context.Context, leagueID uuid.UUID, page, limit int, status string, includeAllMarkets bool) ([]entity.Match, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit
	matches, total, err := uc.matchRepo.GetMatchesByLeagueID(ctx, leagueID, limit, offset, status)
	if err != nil {
		return nil, 0, err
	}

	for i := range matches {
		markets, err := uc.matchRepo.GetMarketsByMatchID(ctx, matches[i].ID)
		if err == nil {
			if includeAllMarkets {
				matches[i].Markets = markets
			} else {
				if len(markets) > 3 {
					matches[i].Markets = markets[:3]
				} else {
					matches[i].Markets = markets
				}
			}
		}
	}

	return matches, total, nil
}
