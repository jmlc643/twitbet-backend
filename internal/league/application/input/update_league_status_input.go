package input

import "github.com/google/uuid"

type UpdateLeagueStatusInput struct {
	LeagueID uuid.UUID
	OwnerID  uuid.UUID
	Status   string
}
