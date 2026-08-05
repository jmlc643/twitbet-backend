package entity

import (
	"time"

	"github.com/google/uuid"
)

type Participant struct {
	ID                uuid.UUID
	LeagueID          uuid.UUID
	UserID            uuid.UUID
	IsAdmin           bool
	Balance           float64
	RechargesConsumed int
	JoinedAt          time.Time
}

type ParticipantRanking struct {
	ParticipantID uuid.UUID
	UserID        uuid.UUID
	Username      string
	ProfilePicture *string
	Balance       float64
	Position      int
	Role          string
}

func NewParticipant(leagueID, userID uuid.UUID, balance float64, isAdmin bool) (*Participant, error) {
	return &Participant{
		ID:                uuid.New(),
		LeagueID:          leagueID,
		UserID:            userID,
		IsAdmin:           isAdmin,
		Balance:           balance,
		RechargesConsumed: 0,
		JoinedAt:          time.Now().UTC(),
	}, nil
}