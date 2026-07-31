package usecase

import (
	"context"

	"github.com/jmlc643/twitbet-backend/internal/league/application/input"
	"github.com/jmlc643/twitbet-backend/internal/league/application/output"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/repository"
)

type GetUserLeaguesUseCase struct {
	repo repository.LeagueRepository
}

func NewGetUserLeaguesUseCase(repo repository.LeagueRepository) *GetUserLeaguesUseCase {
	return &GetUserLeaguesUseCase{repo: repo}
}

func (uc *GetUserLeaguesUseCase) Execute(ctx context.Context, in input.GetUserLeaguesInput) (*output.GetUserLeaguesOutput, error) {
	leagues, err := uc.repo.GetUserLeagues(ctx, in.UserID)
	if err != nil {
		return nil, err
	}

	return &output.GetUserLeaguesOutput{
		Leagues: leagues,
	}, nil
}
