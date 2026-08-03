package input

import "github.com/google/uuid"

type UpdateLeagueInput struct {
	LeagueID       uuid.UUID
	AdminID        uuid.UUID
	Name           string
	InitialBalance float64
	MaxRecharges   int
	HideStandings  bool
}
