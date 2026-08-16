package mapper

import (
	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/valueobject"
	"github.com/jmlc643/twitbet-backend/internal/league/infrastructure/persistence/model"
)

func ToCombinedBetModel(entity *entity.CombinedBet) *model.CombinedBetModel {
	legs := make([]model.CombinedBetLegModel, len(entity.Legs))
	for i, leg := range entity.Legs {
		legs[i] = model.CombinedBetLegModel{
			ID:              leg.ID,
			CombinedBetID:   leg.CombinedBetID,
			MarketID:        leg.MarketID,
			MatchID:         leg.MatchID,
			SelectionID:     leg.SelectionID,
			SelectionName:   leg.SelectionName,
			OddsAtPlacement: leg.OddsAtPlacement,
			Status:          string(leg.Status),
			SettledAt:       leg.SettledAt,
		}
	}

	return &model.CombinedBetModel{
		ID:                 entity.ID,
		ParticipantID:      entity.ParticipantID,
		LeagueID:           entity.LeagueID,
		Stake:              entity.Stake,
		UseBonus:           entity.UseBonus,
		TotalOdds:          entity.TotalOdds,
		PotentialWin:       entity.PotentialWin,
		Status:             string(entity.Status),
		CashoutValue:       entity.CashoutValue,
		CashoutExpiresAt:   entity.CashoutExpiresAt,
		CreatedAt:          entity.CreatedAt,
		UpdatedAt:          entity.UpdatedAt,
		SettledAt:          entity.SettledAt,
		ParticipantBonusID: entity.ParticipantBonusID,
		Legs:               legs,
	}
}

func ToCombinedBetEntity(model *model.CombinedBetModel) *entity.CombinedBet {
	legs := make([]entity.CombinedBetLeg, len(model.Legs))
	for i, leg := range model.Legs {
		legs[i] = entity.CombinedBetLeg{
			ID:              leg.ID,
			CombinedBetID:   leg.CombinedBetID,
			MarketID:        leg.MarketID,
			MatchID:         leg.MatchID,
			SelectionID:     leg.SelectionID,
			SelectionName:   leg.SelectionName,
			OddsAtPlacement: leg.OddsAtPlacement,
			Status:          valueobject.LegStatus(leg.Status),
			SettledAt:       leg.SettledAt,
		}
	}

	return &entity.CombinedBet{
		ID:                 model.ID,
		ParticipantID:      model.ParticipantID,
		LeagueID:           model.LeagueID,
		Stake:              model.Stake,
		UseBonus:           model.UseBonus,
		TotalOdds:          model.TotalOdds,
		PotentialWin:       model.PotentialWin,
		Status:             valueobject.CombinedStatus(model.Status),
		CashoutValue:       model.CashoutValue,
		CashoutExpiresAt:   model.CashoutExpiresAt,
		CreatedAt:          model.CreatedAt,
		UpdatedAt:          model.UpdatedAt,
		SettledAt:          model.SettledAt,
		Legs:               legs,
		ParticipantBonusID: model.ParticipantBonusID,
	}
}
