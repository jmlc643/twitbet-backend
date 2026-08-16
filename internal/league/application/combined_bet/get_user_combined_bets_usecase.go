package combined_bet

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/application/output"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/apperror"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/repository"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/valueobject"
)

type GetUserCombinedBetsUseCase struct {
	combinedBetRepo repository.CombinedBetRepository
	leagueRepo      repository.LeagueRepository
	matchRepo       repository.MatchRepository
}

func NewGetUserCombinedBetsUseCase(combinedBetRepo repository.CombinedBetRepository, leagueRepo repository.LeagueRepository, matchRepo repository.MatchRepository) *GetUserCombinedBetsUseCase {
	return &GetUserCombinedBetsUseCase{
		combinedBetRepo: combinedBetRepo,
		leagueRepo:      leagueRepo,
		matchRepo:       matchRepo,
	}
}

func (uc *GetUserCombinedBetsUseCase) Execute(ctx context.Context, userID, leagueID uuid.UUID, status *string, startDate, endDate *time.Time, page, limit int) ([]output.CombinedBetOutput, int64, error) {
	isParticipant, err := uc.leagueRepo.IsParticipant(ctx, leagueID, userID)
	if err != nil {
		return nil, 0, err
	}
	if !isParticipant {
		return nil, 0, apperror.ErrLeagueNotFound
	}

	participant, err := uc.leagueRepo.GetParticipant(ctx, leagueID, userID)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}

	bets, total, err := uc.combinedBetRepo.GetByParticipantID(ctx, participant.ID, status, startDate, endDate, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	var marketIDs []uuid.UUID
	marketIDMap := make(map[uuid.UUID]bool)
	for _, bet := range bets {
		for _, leg := range bet.Legs {
			if !marketIDMap[leg.MarketID] {
				marketIDMap[leg.MarketID] = true
				marketIDs = append(marketIDs, leg.MarketID)
			}
		}
	}

	markets, err := uc.matchRepo.GetMarketsByIDs(ctx, marketIDs)
	if err != nil {
		return nil, 0, err
	}

	marketNameMap := make(map[uuid.UUID]string)
	marketOptionOdds := make(map[uuid.UUID]map[uuid.UUID]float64)
	for _, m := range markets {
		marketNameMap[m.ID] = m.Name
		marketOptionOdds[m.ID] = make(map[uuid.UUID]float64)
		for _, opt := range m.Options {
			marketOptionOdds[m.ID][opt.ID] = opt.CurrentOdds
		}
	}

	matches, _, err := uc.matchRepo.GetMatchesByLeagueID(ctx, leagueID, 0, 0, "")
	if err != nil {
		return nil, 0, err
	}
	matchTitleMap := make(map[uuid.UUID]string)
	for _, m := range matches {
		matchTitleMap[m.ID] = m.Title
	}

	responses := make([]output.CombinedBetOutput, 0, len(bets))
	for _, bet := range bets {
		legOutputs := make([]output.LegOutput, 0, len(bet.Legs))
		for _, leg := range bet.Legs {
			legStatus := string(leg.Status)
			if bet.Status == valueobject.CombinedStatusCashout && leg.Status == valueobject.LegStatusPending {
				legStatus = string(valueobject.LegStatusCashout)
			}

			mTitle := ""
			if leg.MatchID != nil {
				mTitle = matchTitleMap[*leg.MatchID]
			}
			mName := marketNameMap[leg.MarketID]

			legOutputs = append(legOutputs, output.LegOutput{
				ID:              leg.ID,
				MarketID:        leg.MarketID,
				MatchID:         leg.MatchID,
				MatchTitle:      mTitle,
				MarketName:      mName,
				SelectionName:   leg.SelectionName,
				OddsAtPlacement: leg.OddsAtPlacement,
				Status:          legStatus,
				SettledAt:       leg.SettledAt,
			})
		}

		var cashoutPtr *float64
		if bet.Status == valueobject.CombinedStatusAccepted || bet.Status == valueobject.CombinedStatusPending {
			currentMarketOdds := make(map[uuid.UUID]float64)
			for _, leg := range bet.Legs {
				if leg.Status == valueobject.LegStatusPending {
					if opts, ok := marketOptionOdds[leg.MarketID]; ok {
						if odd, ok := opts[leg.SelectionID]; ok {
							currentMarketOdds[leg.MarketID] = odd
						}
					}
				}
			}
			cv := bet.CalculateCashoutValue(currentMarketOdds)
			if cv > 0 {
				cashoutPtr = &cv
			}
		} else {
			cashoutPtr = bet.CashoutValue
		}

		responses = append(responses, output.CombinedBetOutput{
			ID:               bet.ID,
			UserID:           userID,
			LeagueID:         bet.LeagueID,
			Stake:            bet.Stake,
			UseBonus:         bet.UseBonus,
			TotalOdds:        bet.TotalOdds,
			PotentialWin:     bet.PotentialWin,
			Status:           string(bet.Status),
			CashoutValue:     cashoutPtr,
			CashoutExpiresAt: bet.CashoutExpiresAt,
			CreatedAt:        bet.CreatedAt,
			SettledAt:        bet.SettledAt,
			Legs:             legOutputs,
		})
	}

	return responses, total, nil
}
