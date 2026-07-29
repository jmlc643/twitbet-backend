package output

import "github.com/google/uuid"

type CreateLeagueOutput struct {
	ID         uuid.UUID
	InviteCode string
}