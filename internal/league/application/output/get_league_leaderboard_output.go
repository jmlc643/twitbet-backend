package output

import "github.com/google/uuid"

type LeaderboardEntryOutput struct {
	ParticipantID  uuid.UUID
	UserID         uuid.UUID
	Username       string
	ProfilePicture *string
	Balance        *float64
	Position       *int
	Role           string
	IsUnranked     bool
}

type GetLeagueLeaderboardOutput struct {
	LeagueID         uuid.UUID
	Status           string
	HideStandings    bool
	MinBetsToQualify int
	Leaderboard      []LeaderboardEntryOutput
}
