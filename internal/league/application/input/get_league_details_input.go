package input

import "github.com/google/uuid"

type GetLeagueDetailsInput struct {
	Slug     string
	UserID   uuid.UUID
}