package combined_bet

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/apperror"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/repository"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/valueobject"
)

type CashoutCombinedBetUseCase struct {
	combinedBetRepo repository.CombinedBetRepository
	leagueRepo      repository.LeagueRepository
	matchRepo       repository.MatchRepository
}

func NewCashoutCombinedBetUseCase(
	combinedBetRepo repository.CombinedBetRepository,
	leagueRepo repository.LeagueRepository,
	matchRepo repository.MatchRepository,
) *CashoutCombinedBetUseCase {
	return &CashoutCombinedBetUseCase{
		combinedBetRepo: combinedBetRepo,
		leagueRepo:      leagueRepo,
		matchRepo:       matchRepo,
	}
}

func (uc *CashoutCombinedBetUseCase) Execute(ctx context.Context, userID, betID uuid.UUID) (float64, error) {
	bet, err := uc.combinedBetRepo.GetByID(ctx, betID)
	if err != nil {
		return 0, err
	}
	if bet == nil {
		return 0, apperror.ErrBetNotFound
	}

	participant, err := uc.leagueRepo.GetParticipant(ctx, bet.LeagueID, userID)
	if err != nil {
		return 0, err
	}
	if bet.ParticipantID != participant.ID {
		return 0, apperror.ErrUnauthorized
	}

	if bet.Status != valueobject.CombinedStatusAccepted {
		return 0, apperror.ErrUnauthorized
	}

	currentMarketOdds := make(map[uuid.UUID]float64)
	for _, leg := range bet.Legs {
		if leg.Status == valueobject.LegStatusPending {
			market, err := uc.matchRepo.GetMarketByID(ctx, leg.MarketID)
			if err != nil {
				return 0, err
			}
			if market != nil {
				for _, opt := range market.Options {
					if opt.ID == leg.SelectionID {
						currentMarketOdds[leg.MarketID] = opt.CurrentOdds
						break
					}
				}
			}
		}
	}

	cashoutValue := bet.CalculateCashoutValue(currentMarketOdds)
	if cashoutValue <= 0 {
		return 0, apperror.ErrUnauthorized
	}

	err = uc.combinedBetRepo.UpdateCashout(ctx, bet.ID, cashoutValue)
	if err != nil {
		return 0, err
	}

	participant.Balance += cashoutValue
	_ = uc.leagueRepo.UpdateParticipant(ctx, participant)

	return cashoutValue, nil
}
