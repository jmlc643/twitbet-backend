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
	matchRepo  repository.MatchRepository
}

func NewCashoutBetUseCase(betRepo repository.BetRepository, leagueRepo repository.LeagueRepository, matchRepo repository.MatchRepository) *CashoutBetUseCase {
	return &CashoutBetUseCase{
		betRepo:    betRepo,
		leagueRepo: leagueRepo,
		matchRepo:  matchRepo,
	}
}

func (uc *CashoutBetUseCase) Execute(ctx context.Context, userID, betID uuid.UUID) (*entity.Bet, error) {
	bet, err := uc.betRepo.GetBetByID(ctx, betID)
	if err != nil {
		return nil, err
	}
	if bet == nil {
		return nil, errors.New("Apuesta no encontrada")
	}

	if bet.Status != entity.BetStatusAccepted && bet.Status != entity.BetStatusPending {
		return nil, errors.New("La apuesta no está activa")
	}

	if bet.IsBonusBet {
		return nil, errors.New("No se puede hacer cashout de una apuesta de bono")
	}

	participant, err := uc.leagueRepo.GetParticipantByID(ctx, bet.ParticipantID)
	if err != nil {
		return nil, err
	}
	if participant == nil || participant.UserID != userID {
		return nil, errors.New("No autorizado para hacer cashout de esta apuesta")
	}

	currentOdds, err := uc.matchRepo.GetMarketOptionCurrentOdds(ctx, bet.MarketOptionID)
	if err != nil {
		return nil, err
	}

	cashoutAmount := bet.CalculateCashoutValue(currentOdds)
	if cashoutAmount <= 0 {
		return nil, errors.New("Cashout no disponible para esta apuesta")
	}

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
