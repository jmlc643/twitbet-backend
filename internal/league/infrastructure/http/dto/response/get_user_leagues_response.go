package response

type LeagueSummaryResponse struct {
	LeagueID         string  `json:"league_id"`
	Slug             string  `json:"slug"`
	Name             string  `json:"name"`
	Role             string  `json:"role"`
	ParticipantCount int     `json:"participant_count"`
	Balance          float64 `json:"balance"`
}

type GetUserLeaguesResponse struct {
	Leagues []LeagueSummaryResponse `json:"leagues"`
}
