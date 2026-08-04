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

func (uc *GetLeagueMatchesUseCase) Execute(ctx context.Context, leagueID uuid.UUID, page, limit int, status string) ([]entity.Match, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit
	return uc.matchRepo.GetMatchesByLeagueID(ctx, leagueID, limit, offset, status)
}
