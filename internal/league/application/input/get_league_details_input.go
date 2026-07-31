package input

import "github.com/google/uuid"

type GetLeagueDetailsInput struct {
	LeagueID uuid.UUID
	UserID   uuid.UUID
}
