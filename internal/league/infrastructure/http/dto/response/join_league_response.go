package response

type JoinLeagueResponse struct {
	LeagueID   string  `json:"league_id"`
	LeagueName string  `json:"league_name"`
	Balance    float64 `json:"balance"`
}
