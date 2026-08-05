package input

import "github.com/google/uuid"

type RemoveAdminInput struct {
	LeagueID      uuid.UUID
	OwnerID       uuid.UUID
	ParticipantID uuid.UUID
}
