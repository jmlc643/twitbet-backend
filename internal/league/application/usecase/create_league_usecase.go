package usecase

import (
	"context"

	"github.com/jmlc643/twitbet-backend/internal/league/application/input"
	"github.com/jmlc643/twitbet-backend/internal/league/application/output"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/repository"
)

type CreateLeagueUseCase struct {
	repo repository.LeagueRepository
}

func NewCreateLeagueUseCase(repo repository.LeagueRepository) *CreateLeagueUseCase {
	return &CreateLeagueUseCase{repo: repo}
}

func (uc *CreateLeagueUseCase) Execute(ctx context.Context, in input.CreateLeagueInput) (*output.CreateLeagueOutput, error) {
	league, err := entity.NewLeague(in.AdminID, in.Name, in.InitialBalance, in.MaxRecharges, in.HideStandings)
	if err != nil {
		return nil, err
	}

	participant, err := entity.NewParticipant(league.ID, in.AdminID, in.InitialBalance)
	if err != nil {
		return nil, err
	}

	err = uc.repo.CreateLeagueWithAdminParticipant(ctx, league, participant)
	if err != nil {
		return nil, err
	}

	return &output.CreateLeagueOutput{
		ID:         league.ID,
		InviteCode: league.InviteCode,
	}, nil
}
