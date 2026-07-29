package entity

import (
	"time"

	"github.com/google/uuid"
)

type Participant struct {
	ID       uuid.UUID
	LeagueID uuid.UUID
	UserID   uuid.UUID
	Balance  float64
	JoinedAt time.Time
}

func NewParticipant(leagueID, userID uuid.UUID, balance float64) (*Participant, error) {
	return &Participant{
		ID:       uuid.New(),
		LeagueID: leagueID,
		UserID:   userID,
		Balance:  balance,
		JoinedAt: time.Now().UTC(),
	}, nil
}
