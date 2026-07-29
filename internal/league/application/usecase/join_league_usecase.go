package usecase

import (
	"context"

	"github.com/jmlc643/twitbet-backend/internal/league/application/input"
	"github.com/jmlc643/twitbet-backend/internal/league/application/output"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/apperror"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/repository"
)

type JoinLeagueUseCase struct {
	repo repository.LeagueRepository
}

func NewJoinLeagueUseCase(repo repository.LeagueRepository) *JoinLeagueUseCase {
	return &JoinLeagueUseCase{repo: repo}
}

func (uc *JoinLeagueUseCase) Execute(ctx context.Context, in input.JoinLeagueInput) (*output.JoinLeagueOutput, error) {
	league, err := uc.repo.GetLeagueByInviteCode(ctx, in.InviteCode)
	if err != nil {
		return nil, err
	}
	if league == nil {
		return nil, apperror.ErrInvalidInviteCode
	}

	isParticipant, err := uc.repo.IsParticipant(ctx, league.ID, in.UserID)
	if err != nil {
		return nil, err
	}
	if isParticipant {
		return nil, apperror.ErrUserAlreadyJoined
	}

	participant, err := entity.NewParticipant(league.ID, in.UserID, league.InitialBalance)
	if err != nil {
		return nil, err
	}

	transaction, err := entity.NewTransaction(league.ID, in.UserID, league.InitialBalance, entity.TransactionTypeInitialBalance)
	if err != nil {
		return nil, err
	}

	err = uc.repo.JoinLeague(ctx, participant, transaction)
	if err != nil {
		return nil, err
	}

	return &output.JoinLeagueOutput{
		LeagueID:   league.ID,
		LeagueName: league.Name,
		Balance:    league.InitialBalance,
	}, nil
}
