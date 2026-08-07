package output

import (
	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
)

type ParticipantRankingOutput struct {
	ParticipantID  uuid.UUID
	UserID         uuid.UUID
	Username       string
	ProfilePicture *string
	Balance        float64
	Position       int
	Role           string
}

type GetLeagueDetailsOutput struct {
	LeagueID         uuid.UUID
	Slug             string
	Name             string
	OwnerID          uuid.UUID
	InitialBalance   float64
	MaxRecharges     int
	IsRankingVisible bool
	InviteCode       string
	CreatedAt        string
	ParticipantsCount int
	Participants     []ParticipantRankingOutput
}

func EntityToParticipantRankingOutput(rankings []entity.ParticipantRanking) []ParticipantRankingOutput {
	var output []ParticipantRankingOutput
	for _, r := range rankings {
		output = append(output, ParticipantRankingOutput{
			ParticipantID: r.ParticipantID,
			UserID:        r.UserID,
			Username:       r.Username,
			ProfilePicture: r.ProfilePicture,
			Balance:        r.Balance,
			Position:       r.Position,
			Role:           r.Role,
		})
	}
	return output
}