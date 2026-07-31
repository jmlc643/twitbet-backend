package output

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
)

type ParticipantRankingOutput struct {
	ParticipantID uuid.UUID
	UserID        uuid.UUID
	Username      string
	ProfilePicture *string
	Balance       float64
	Position      int
}

type GetLeagueDetailsOutput struct {
	LeagueID         uuid.UUID
	Name             string
	AdminID          uuid.UUID
	InitialBalance   float64
	MaxRecharges     int
	IsRankingVisible bool
	InviteCode       string
	CreatedAt        time.Time
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
		})
	}
	return output
}
