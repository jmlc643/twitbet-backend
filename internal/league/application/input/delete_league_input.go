package input

import "github.com/google/uuid"

type DeleteLeagueInput struct {
	LeagueID uuid.UUID
	AdminID  uuid.UUID
}
