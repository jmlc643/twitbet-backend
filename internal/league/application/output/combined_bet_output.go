package output

import (
	"time"

	"github.com/google/uuid"
)

type LegOutput struct {
	ID              uuid.UUID
	MarketID        uuid.UUID
	MatchID         *uuid.UUID
	SelectionName   string
	OddsAtPlacement float64
	Status          string
	SettledAt       *time.Time
}

type CombinedBetOutput struct {
	ID               uuid.UUID
	UserID           uuid.UUID
	LeagueID         uuid.UUID
	Stake            float64
	UseBonus         bool
	TotalOdds        float64
	PotentialWin     float64
	Status           string
	CashoutValue     *float64
	CashoutExpiresAt *time.Time
	CreatedAt        time.Time
	SettledAt        *time.Time
	Legs             []LegOutput
}
