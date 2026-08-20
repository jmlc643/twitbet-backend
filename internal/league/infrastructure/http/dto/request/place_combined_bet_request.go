package request

import "github.com/google/uuid"

type SelectionDTO struct {
	MarketID     uuid.UUID `json:"market_id" binding:"required"`
	SelectionID  uuid.UUID `json:"selection_id" binding:"required"`
	AcceptedOdds float64   `json:"accepted_odds" binding:"required,gt=1"`
}

type PlaceCombinedBetRequest struct {
	LeagueID   uuid.UUID      `json:"league_id" binding:"required"`
	Stake      float64        `json:"stake" binding:"required,gt=0"`
	UseBonus   bool           `json:"use_bonus"`
	BonusID    *uuid.UUID     `json:"bonus_id,omitempty"`
	Selections []SelectionDTO `json:"selections" binding:"required,min=2,dive"`
}
