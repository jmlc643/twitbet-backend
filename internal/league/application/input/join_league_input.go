package input

import "github.com/google/uuid"

type JoinLeagueInput struct {
	UserID     uuid.UUID
	InviteCode string
}
