package input

import "github.com/google/uuid"

type UpdateLeagueSettingsInput struct {
	LeagueID         uuid.UUID
	AdminID          uuid.UUID
	IsRankingVisible bool
}
