package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/valueobject"
)

type CombinedBetLeg struct {
	ID              uuid.UUID
	CombinedBetID   uuid.UUID
	MarketID        uuid.UUID
	MatchID         *uuid.UUID
	SelectionID     uuid.UUID
	SelectionName   string
	OddsAtPlacement float64
	Status          valueobject.LegStatus
	SettledAt       *time.Time
}

func NewCombinedBetLeg(combinedBetID, marketID, selectionID uuid.UUID, matchID *uuid.UUID, selectionName string, odds float64) CombinedBetLeg {
	return CombinedBetLeg{
		ID:              uuid.New(),
		CombinedBetID:   combinedBetID,
		MarketID:        marketID,
		MatchID:         matchID,
		SelectionID:     selectionID,
		SelectionName:   selectionName,
		OddsAtPlacement: odds,
		Status:          valueobject.LegStatusPending,
		SettledAt:       nil,
	}
}
