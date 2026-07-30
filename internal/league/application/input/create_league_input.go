package input

import "github.com/google/uuid"

type CreateLeagueInput struct {
	AdminID        uuid.UUID
	Name           string
	InitialBalance float64
	MaxRecharges   int
	HideStandings  bool
}