package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/repository"
)

type CashoutBetUseCase struct {
	betRepo    repository.BetRepository
	leagueRepo repository.LeagueRepository
}

func NewCashoutBetUseCase(betRepo repository.BetRepository, leagueRepo repository.LeagueRepository) *CashoutBetUseCase {
	return &CashoutBetUseCase{
		betRepo:    betRepo,
		leagueRepo: leagueRepo,
	}
}

func (uc *CashoutBetUseCase) Execute(ctx context.Context, userID, betID uuid.UUID) (*entity.Bet, error) {
	bet, err := uc.betRepo.GetBetByID(ctx, betID)
	if err != nil {
		return nil, err
	}
	if bet == nil {
		return nil, errors.New("bet not found")
	}

	if bet.Status != entity.BetStatusAccepted && bet.Status != entity.BetStatusPending {
		return nil, errors.New("bet is not active")
	}

	if bet.IsBonusBet {
		return nil, errors.New("cannot cashout a bonus bet")
	}

	participant, err := uc.leagueRepo.GetParticipantByID(ctx, bet.ParticipantID)
	if err != nil {
		return nil, err
	}
	if participant == nil || participant.UserID != userID {
		return nil, errors.New("unauthorized to cashout this bet")
	}

	cashoutAmount := bet.Amount * 0.95

	transaction, err := entity.NewTransaction(participant.LeagueID, userID, cashoutAmount, entity.TransactionTypeCashout)
	if err != nil {
		return nil, err
	}

	bet.Status = entity.BetStatusCashout
	bet.UpdatedAt = transaction.CreatedAt

	err = uc.betRepo.CashoutAtomic(ctx, bet, transaction)
	if err != nil {
		return nil, err
	}

	return bet, nil
}
