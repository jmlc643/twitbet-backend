package request

type UpdateLeagueRequest struct {
	Name           string  `json:"name" binding:"required"`
	InitialBalance float64 `json:"initial_balance" binding:"required"`
	MaxRecharges   int     `json:"max_recharges" binding:"required"`
	HideStandings  bool    `json:"hide_standings"`
}
