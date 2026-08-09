package response

type LeaderboardEntryResponse struct {
	ParticipantID  string   `json:"participant_id"`
	UserID         string   `json:"user_id"`
	Username       string   `json:"username"`
	ProfilePicture *string  `json:"profile_picture"`
	Balance        *float64 `json:"balance"`
	Position       *int     `json:"position"`
	Role           string   `json:"role"`
	IsUnranked     bool     `json:"is_unranked"`
}

type GetLeagueLeaderboardResponse struct {
	LeagueID         string                     `json:"league_id"`
	Status           string                     `json:"status"`
	HideStandings    bool                       `json:"hide_standings"`
	MinBetsToQualify int                        `json:"min_bets_to_qualify"`
	Leaderboard      []LeaderboardEntryResponse `json:"leaderboard"`
}
