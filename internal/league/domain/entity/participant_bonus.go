package entity

import (
	"time"

	"github.com/google/uuid"
)

type BonusStatus string

const (
	BonusStatusPending BonusStatus = "PENDING"
	BonusStatusUsed    BonusStatus = "USED"
)

type ParticipantBonus struct {
	ID            uuid.UUID
	LeagueID      uuid.UUID
	ParticipantID uuid.UUID
	Amount        float64
	Status        BonusStatus
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewParticipantBonus(leagueID, participantID uuid.UUID, amount float64) *ParticipantBonus {
	now := time.Now().UTC()
	return &ParticipantBonus{
		ID:            uuid.New(),
		LeagueID:      leagueID,
		ParticipantID: participantID,
		Amount:        amount,
		Status:        BonusStatusPending,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}
