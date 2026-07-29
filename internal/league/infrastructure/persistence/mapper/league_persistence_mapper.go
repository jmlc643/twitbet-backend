package mapper

import (
	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
	"github.com/jmlc643/twitbet-backend/internal/league/infrastructure/persistence/model"
)

func EntityToLeagueModel(e *entity.League) *model.LeagueModel {
	return &model.LeagueModel{
		ID:             e.ID.String(),
		AdminID:        e.AdminID.String(),
		Name:           e.Name,
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
		ID:       e.ID.String(),
		LeagueID: e.LeagueID.String(),
		UserID:   e.UserID.String(),
		Balance:  e.Balance,
		JoinedAt: e.JoinedAt,
	}
}
