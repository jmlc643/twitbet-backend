package input

import "github.com/google/uuid"

type DeleteLeagueInput struct {
	LeagueID uuid.UUID
	OwnerID  uuid.UUID
}


