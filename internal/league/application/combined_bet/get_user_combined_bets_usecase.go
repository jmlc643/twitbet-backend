package combined_bet

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/application/output"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/apperror"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/repository"
)

type GetUserCombinedBetsUseCase struct {
	combinedBetRepo repository.CombinedBetRepository
	leagueRepo      repository.LeagueRepository
}

func NewGetUserCombinedBetsUseCase(combinedBetRepo repository.CombinedBetRepository, leagueRepo repository.LeagueRepository) *GetUserCombinedBetsUseCase {
	return &GetUserCombinedBetsUseCase{
		combinedBetRepo: combinedBetRepo,
		leagueRepo:      leagueRepo,
	}
}

func (uc *GetUserCombinedBetsUseCase) Execute(ctx context.Context, userID, leagueID uuid.UUID) ([]output.CombinedBetOutput, error) {
	isParticipant, err := uc.leagueRepo.IsParticipant(ctx, leagueID, userID)
	if err != nil {
		return nil, err
	}
	if !isParticipant {
		return nil, apperror.ErrLeagueNotFound
	}

	participant, err := uc.leagueRepo.GetParticipant(ctx, leagueID, userID)
	if err != nil {
		return nil, err
	}

	bets, err := uc.combinedBetRepo.GetByParticipantID(ctx, participant.ID)
	if err != nil {
		return nil, err
	}

	responses := make([]output.CombinedBetOutput, 0, len(bets))
	for _, bet := range bets {
		legOutputs := make([]output.LegOutput, 0, len(bet.Legs))
		for _, leg := range bet.Legs {
			legOutputs = append(legOutputs, output.LegOutput{
				ID:              leg.ID,
				MarketID:        leg.MarketID,
				MatchID:         leg.MatchID,
				SelectionName:   leg.SelectionName,
				OddsAtPlacement: leg.OddsAtPlacement,
				Status:          string(leg.Status),
				SettledAt:       leg.SettledAt,
			})
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
			CashoutValue:     bet.CashoutValue,
			CashoutExpiresAt: bet.CashoutExpiresAt,
			CreatedAt:        bet.CreatedAt,
			SettledAt:        bet.SettledAt,
			Legs:             legOutputs,
		})
	}

	return responses, nil
}
