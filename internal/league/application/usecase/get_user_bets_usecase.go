package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/apperror"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/repository"
)

type GetUserBetsUseCase struct {
	betRepo    repository.BetRepository
	leagueRepo repository.LeagueRepository
	matchRepo  repository.MatchRepository
}

func NewGetUserBetsUseCase(betRepo repository.BetRepository, leagueRepo repository.LeagueRepository, matchRepo repository.MatchRepository) *GetUserBetsUseCase {
	return &GetUserBetsUseCase{
		betRepo:    betRepo,
		leagueRepo: leagueRepo,
		matchRepo:  matchRepo,
	}
}

func (uc *GetUserBetsUseCase) Execute(ctx context.Context, userID, leagueID uuid.UUID, status *string, startDate, endDate *time.Time, page, limit int) ([]entity.BetDetail, int64, error) {
	participant, err := uc.leagueRepo.GetParticipant(ctx, leagueID, userID)
	if err != nil {
		return nil, 0, err
	}
	if participant == nil {
		return nil, 0, apperror.ErrUnauthorized
	}

	var betStatus *entity.BetStatus
	if status != nil && *status != "" {
		st := entity.BetStatus(*status)
		betStatus = &st
	}

	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}

	bets, total, err := uc.betRepo.GetBetsByParticipantID(ctx, participant.ID, betStatus, startDate, endDate, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	if bets == nil {
		bets = []entity.BetDetail{}
	}

	var marketIDs []uuid.UUID
	marketIDMap := make(map[uuid.UUID]bool)
	for i := range bets {
		if bets[i].Status == entity.BetStatusAccepted || bets[i].Status == entity.BetStatusPending {
			if !marketIDMap[bets[i].MarketID] {
				marketIDMap[bets[i].MarketID] = true
				marketIDs = append(marketIDs, bets[i].MarketID)
			}
		}
	}

	var markets []entity.Market
	if len(marketIDs) > 0 {
		markets, _ = uc.matchRepo.GetMarketsByIDs(ctx, marketIDs)
	}

	optionOddsMap := make(map[uuid.UUID]float64)
	for _, m := range markets {
		for _, opt := range m.Options {
			optionOddsMap[opt.ID] = opt.CurrentOdds
		}
	}

	for i := range bets {
		if bets[i].Status == entity.BetStatusAccepted || bets[i].Status == entity.BetStatusPending {
			currentOdds, ok := optionOddsMap[bets[i].OptionID]
			if ok {
				tempBet := &entity.Bet{
					Amount:       bets[i].Amount,
					Odds:         bets[i].Odds,
					PotentialWin: bets[i].PotentialWin,
				}
				ca := tempBet.CalculateCashoutValue(currentOdds)
				if ca > 0 {
					bets[i].CashoutAmount = &ca
				}
			} else {
				ca := bets[i].Amount * 0.90
				bets[i].CashoutAmount = &ca
			}
		}
	}

	return bets, total, nil
}