package combined_bet

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/application/output"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/apperror"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/repository"
)

type GetCombinedBetDetailsUseCase struct {
	combinedBetRepo repository.CombinedBetRepository
	leagueRepo      repository.LeagueRepository
}

func NewGetCombinedBetDetailsUseCase(combinedBetRepo repository.CombinedBetRepository, leagueRepo repository.LeagueRepository) *GetCombinedBetDetailsUseCase {
	return &GetCombinedBetDetailsUseCase{
		combinedBetRepo: combinedBetRepo,
		leagueRepo:      leagueRepo,
	}
}

func (uc *GetCombinedBetDetailsUseCase) Execute(ctx context.Context, betID, userID uuid.UUID) (*output.CombinedBetOutput, error) {
	bet, err := uc.combinedBetRepo.GetByID(ctx, betID)
	if err != nil {
		return nil, err
	}
	if bet == nil {
		return nil, apperror.ErrBetNotFound
	}

	participant, err := uc.leagueRepo.GetParticipant(ctx, bet.LeagueID, userID)
	if err != nil {
		return nil, err
	}
	if bet.ParticipantID != participant.ID {
		return nil, apperror.ErrUnauthorized
	}

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

	response := &output.CombinedBetOutput{
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
	}

	return response, nil
}
