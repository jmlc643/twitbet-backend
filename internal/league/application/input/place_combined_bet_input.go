package input

import "github.com/google/uuid"

type SelectionInput struct {
	MarketID     uuid.UUID
	SelectionID  uuid.UUID
	AcceptedOdds float64
}

type PlaceCombinedBetInput struct {
	LeagueID   uuid.UUID
	Stake      float64
	UseBonus   bool
	BonusID    *uuid.UUID
	Selections []SelectionInput
}
