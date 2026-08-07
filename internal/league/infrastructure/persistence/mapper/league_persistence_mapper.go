package mapper

import (
	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
	"github.com/jmlc643/twitbet-backend/internal/league/infrastructure/persistence/model"
)

func EntityToLeagueModel(e *entity.League) *model.LeagueModel {
	return &model.LeagueModel{
		ID:             e.ID.String(),
		OwnerID:        e.OwnerID.String(),
		Name:           e.Name,
		Slug:           e.Slug,
		InitialBalance: e.InitialBalance,
		MaxRecharges:   e.MaxRecharges,
		HideStandings:  e.HideStandings,
		InviteCode:     e.InviteCode,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
	}
}

func EntityToParticipantModel(e *entity.Participant) *model.ParticipantModel {
	return &model.ParticipantModel{
		ID:                e.ID.String(),
		LeagueID:          e.LeagueID.String(),
		UserID:            e.UserID.String(),
		IsAdmin:           e.IsAdmin,
		Balance:           e.Balance,
		RechargesConsumed: e.RechargesConsumed,
		JoinedAt:          e.JoinedAt,
	}
}

func EntityToTransactionModel(e *entity.Transaction) *model.TransactionModel {
	return &model.TransactionModel{
		ID:        e.ID.String(),
		LeagueID:  e.LeagueID.String(),
		UserID:    e.UserID.String(),
		Amount:    e.Amount,
		Type:      string(e.Type),
		CreatedAt: e.CreatedAt,
	}
}

func EntityToBetModel(e *entity.Bet) *model.BetModel {
	var bonusID *string
	if e.ParticipantBonusID != nil {
		idStr := e.ParticipantBonusID.String()
		bonusID = &idStr
	}

	return &model.BetModel{
		ID:                 e.ID.String(),
		ParticipantID:      e.ParticipantID.String(),
		MarketOptionID:     e.MarketOptionID.String(),
		Amount:             e.Amount,
		Odds:               e.Odds,
		PotentialWin:       e.PotentialWin,
		Status:             string(e.Status),
		PlacedAt:           e.PlacedAt,
		UpdatedAt:          e.UpdatedAt,
		IsBonusBet:         e.IsBonusBet,
		ParticipantBonusID: bonusID,
	}
}

func EntityToParticipantBonusModel(e *entity.ParticipantBonus) *model.ParticipantBonusModel {
	return &model.ParticipantBonusModel{
		ID:            e.ID.String(),
		LeagueID:      e.LeagueID.String(),
		ParticipantID: e.ParticipantID.String(),
		Amount:        e.Amount,
		Status:        string(e.Status),
		CreatedAt:     e.CreatedAt,
		UpdatedAt:     e.UpdatedAt,
	}
}
