package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/apperror"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/port"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/repository"
)

type CashoutBetUseCase struct {
	betRepo         repository.BetRepository
	leagueRepo      repository.LeagueRepository
	matchRepo       repository.MatchRepository
	marketPublisher port.MarketEventPublisher
}

func NewCashoutBetUseCase(betRepo repository.BetRepository, leagueRepo repository.LeagueRepository, matchRepo repository.MatchRepository, marketPublisher port.MarketEventPublisher) *CashoutBetUseCase {
	return &CashoutBetUseCase{
		betRepo:         betRepo,
		leagueRepo:      leagueRepo,
		matchRepo:       matchRepo,
		marketPublisher: marketPublisher,
	}
}

func (uc *CashoutBetUseCase) Execute(ctx context.Context, userID, betID uuid.UUID) (*entity.Bet, error) {
	bet, err := uc.betRepo.GetBetByID(ctx, betID)
	if err != nil {
		return nil, err
	}
	if bet == nil {
		return nil, apperror.ErrBetNotFound
	}

	if bet.Status != entity.BetStatusAccepted {
		return nil, apperror.ErrCashoutNotAvailable
	}

	if bet.IsBonusBet {
		return nil, errors.New("No se puede hacer cashout de una apuesta de bono")
	}

	participant, err := uc.leagueRepo.GetParticipantByID(ctx, bet.ParticipantID)
	if err != nil {
		return nil, err
	}
	if participant == nil || participant.UserID != userID {
		return nil, apperror.ErrUnauthorized
	}

	market, err := uc.matchRepo.GetMarketByOptionID(ctx, bet.MarketOptionID)
	if err != nil {
		return nil, err
	}
	if market != nil {
		if entity.NotResolvableMarketStatus(market) {
			return nil, apperror.ErrCashoutNotAvailable
		}
		if market.Status == string(entity.MarketStatusSuspended) || market.Status == string(entity.MarketStatusCancelled) {
			return nil, apperror.ErrCashoutNotAvailable
		}
		for _, opt := range market.Options {
			if opt.ID == bet.MarketOptionID && opt.IsBlocked() {
				return nil, apperror.ErrCashoutNotAvailable
			}
		}
	}

	currentOdds, err := uc.matchRepo.GetMarketOptionCurrentOdds(ctx, bet.MarketOptionID)
	if err != nil {
		return nil, err
	}

	cashoutAmount := bet.CalculateCashoutValue(currentOdds)
	if cashoutAmount <= 0 {
		return nil, apperror.ErrCashoutNotAvailable
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

	if uc.marketPublisher != nil {
		_ = uc.marketPublisher.PublishParticipantBalanceUpdated(ctx, participant.ID, participant.LeagueID, userID)
	}

	return bet, nil
}
