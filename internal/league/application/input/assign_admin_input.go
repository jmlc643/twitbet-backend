package input

import "github.com/google/uuid"

type AssignAdminInput struct {
	LeagueID      uuid.UUID
	OwnerID       uuid.UUID
	ParticipantID uuid.UUID
}
