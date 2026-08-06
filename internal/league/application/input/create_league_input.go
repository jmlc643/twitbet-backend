package input

import "github.com/google/uuid"

type CreateLeagueInput struct {
	OwnerID        uuid.UUID
	Name           string
	InitialBalance float64
	MaxRecharges   int
	HideStandings  bool
}