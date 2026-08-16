package response

import (
	"time"

	"github.com/google/uuid"
)

type LegResponse struct {
	ID              uuid.UUID  `json:"id"`
	MarketID        uuid.UUID  `json:"market_id"`
	MatchID         *uuid.UUID `json:"match_id"`
	MatchTitle      string     `json:"match_title"`
	MarketName      string     `json:"market_name"`
	SelectionName   string     `json:"selection_name"`
	OddsAtPlacement float64    `json:"odds_at_placement"`
	Status          string     `json:"status"`
	SettledAt       *time.Time `json:"settled_at,omitempty"`
}

type CombinedBetResponse struct {
	ID               uuid.UUID     `json:"id"`
	UserID           uuid.UUID     `json:"user_id"`
	LeagueID         uuid.UUID     `json:"league_id"`
	Stake            float64       `json:"stake"`
	UseBonus         bool          `json:"use_bonus"`
	TotalOdds        float64       `json:"total_odds"`
	PotentialWin     float64       `json:"potential_win"`
	Status           string        `json:"status"`
	CashoutValue     *float64      `json:"cashout_value,omitempty"`
	CashoutExpiresAt *time.Time    `json:"cashout_expires_at,omitempty"`
	CreatedAt        time.Time     `json:"created_at"`
	SettledAt        *time.Time    `json:"settled_at,omitempty"`
	Legs             []LegResponse `json:"legs"`
}

type PaginatedCombinedBetResponse struct {
	Data []CombinedBetResponse `json:"data"`
	Meta PaginationMeta      `json:"meta"`
}
