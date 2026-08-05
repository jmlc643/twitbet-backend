package output

import "github.com/google/uuid"

type JoinLeagueOutput struct {
	LeagueID   uuid.UUID
	LeagueName string
	Slug       string
	Balance    float64
}
