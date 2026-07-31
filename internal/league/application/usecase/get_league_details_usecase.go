package usecase

import (
	"context"
	"errors"

	"github.com/jmlc643/twitbet-backend/internal/league/application/input"
	"github.com/jmlc643/twitbet-backend/internal/league/application/output"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/repository"
)

type GetLeagueDetailsUseCase struct {
	repo repository.LeagueRepository
}

func NewGetLeagueDetailsUseCase(repo repository.LeagueRepository) *GetLeagueDetailsUseCase {
	return &GetLeagueDetailsUseCase{repo: repo}
}

func (uc *GetLeagueDetailsUseCase) Execute(ctx context.Context, in input.GetLeagueDetailsInput) (output.GetLeagueDetailsOutput, error) {
	league, err := uc.repo.GetLeagueByID(ctx, in.LeagueID)
	if err != nil {
		return output.GetLeagueDetailsOutput{}, err
	}
	if league == nil {
		return output.GetLeagueDetailsOutput{}, errors.New("liga no encontrada")
	}

	participants, err := uc.repo.GetParticipantsByLeagueID(ctx, in.LeagueID)
	if err != nil {
		return output.GetLeagueDetailsOutput{}, err
	}

	if league.HideStandings && in.UserID != league.AdminID {
		participants = nil
	}

	return output.GetLeagueDetailsOutput{
		LeagueID:         league.ID,
		Name:             league.Name,
		AdminID:          league.AdminID,
		InitialBalance:   league.InitialBalance,
		MaxRecharges:     league.MaxRecharges,
		IsRankingVisible: !league.HideStandings,
		InviteCode:       league.InviteCode,
		CreatedAt:        league.CreatedAt,
		Participants:     output.EntityToParticipantRankingOutput(participants),
	}, nil
}
