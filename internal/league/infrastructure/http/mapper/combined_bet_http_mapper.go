package mapper

import (
	"github.com/jmlc643/twitbet-backend/internal/league/application/input"
	"github.com/jmlc643/twitbet-backend/internal/league/application/output"
	"github.com/jmlc643/twitbet-backend/internal/league/infrastructure/http/dto/request"
	"github.com/jmlc643/twitbet-backend/internal/league/infrastructure/http/dto/response"
)

func ToPlaceCombinedBetInput(req request.PlaceCombinedBetRequest) input.PlaceCombinedBetInput {
	selections := make([]input.SelectionInput, len(req.Selections))
	for i, s := range req.Selections {
		selections[i] = input.SelectionInput{
			MarketID:    s.MarketID,
			SelectionID: s.SelectionID,
		}
	}

	return input.PlaceCombinedBetInput{
		LeagueID:   req.LeagueID,
		Stake:      req.Stake,
		UseBonus:   req.UseBonus,
		BonusID:    req.BonusID,
		Selections: selections,
	}
}

func ToCombinedBetResponse(out output.CombinedBetOutput) response.CombinedBetResponse {
	legs := make([]response.LegResponse, len(out.Legs))
	for i, l := range out.Legs {
		legs[i] = response.LegResponse{
			ID:              l.ID,
			MarketID:        l.MarketID,
			MatchID:         l.MatchID,
			MatchTitle:      l.MatchTitle,
			MarketName:      l.MarketName,
			SelectionName:   l.SelectionName,
			OddsAtPlacement: l.OddsAtPlacement,
			Status:          l.Status,
			SettledAt:       l.SettledAt,
		}
	}

	return response.CombinedBetResponse{
		ID:               out.ID,
		UserID:           out.UserID,
		LeagueID:         out.LeagueID,
		Stake:            out.Stake,
		UseBonus:         out.UseBonus,
		TotalOdds:        out.TotalOdds,
		PotentialWin:     out.PotentialWin,
		Status:           out.Status,
		CashoutValue:     out.CashoutValue,
		CashoutExpiresAt: out.CashoutExpiresAt,
		CreatedAt:        out.CreatedAt,
		SettledAt:        out.SettledAt,
		Legs:             legs,
	}
}
