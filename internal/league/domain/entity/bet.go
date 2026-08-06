package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/apperror"
)

type BetStatus string

const (
	BetStatusPending  BetStatus = "PENDING"
	BetStatusAccepted BetStatus = "ACCEPTED"
	BetStatusRejected BetStatus = "REJECTED"
	BetStatusWon      BetStatus = "WON"
	BetStatusLost     BetStatus = "LOST"
	BetStatusVoided   BetStatus = "VOIDED"
)

type Bet struct {
	ID             uuid.UUID
	ParticipantID  uuid.UUID
	MarketOptionID uuid.UUID
	Amount         float64
	Odds           float64
	PotentialWin   float64
	Status         BetStatus
	PlacedAt       time.Time
	UpdatedAt      time.Time
}

func NewBet(participantID, marketOptionID uuid.UUID, amount, odds float64) (*Bet, error) {
	if amount <= 0 {
		return nil, apperror.ErrInvalidBetAmount
	}

	now := time.Now().UTC()
	return &Bet{
		ID:             uuid.New(),
		ParticipantID:  participantID,
		MarketOptionID: marketOptionID,
		Amount:         amount,
		Odds:           odds,
		PotentialWin:   amount * odds,
		Status:         BetStatusPending,
		PlacedAt:       now,
		UpdatedAt:      now,
	}, nil
}

type BetDetail struct {
	ID           uuid.UUID
	Amount       float64
	Odds         float64
	PotentialWin float64
	Status       BetStatus
	PlacedAt     time.Time
	MatchTitle   string
	MarketID     uuid.UUID
	MarketName   string
	OptionID     uuid.UUID
	OptionName   string
}
