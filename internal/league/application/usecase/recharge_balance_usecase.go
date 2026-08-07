package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/apperror"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/repository"
)

type RechargeBalanceUseCase struct {
	leagueRepo repository.LeagueRepository
}

func NewRechargeBalanceUseCase(leagueRepo repository.LeagueRepository) *RechargeBalanceUseCase {
	return &RechargeBalanceUseCase{
		leagueRepo: leagueRepo,
	}
}

func (uc *RechargeBalanceUseCase) Execute(ctx context.Context, userID, leagueID uuid.UUID) (float64, int, error) {
	league, err := uc.leagueRepo.GetLeagueByID(ctx, leagueID)
	if err != nil {
		return 0, 0, err
	}
	if league == nil {
		return 0, 0, apperror.ErrLeagueNotFound
	}

	participant, err := uc.leagueRepo.GetParticipant(ctx, leagueID, userID)
	if err != nil {
		return 0, 0, err
	}
	if participant == nil {
		return 0, 0, apperror.ErrLeagueNotFound
	}

	if participant.RechargesConsumed >= league.MaxRecharges {
		return 0, 0, errors.New("max recharges reached")
	}

	rechargeAmount := league.InitialBalance * 0.5
	participant.Balance += rechargeAmount
	participant.RechargesConsumed++

	err = uc.leagueRepo.UpdateParticipant(ctx, participant)
	if err != nil {
		return 0, 0, err
	}

	return participant.Balance, participant.RechargesConsumed, nil
}
