package request

type CreateLeagueRequest struct {
	Name           string  `json:"name" binding:"required"`
	InitialBalance float64 `json:"initial_balance" binding:"required,min=0"`
	MaxRecharges   int     `json:"max_recharges"`
	HideStandings  bool    `json:"hide_standings"`
	MinBetsToQualify int   `json:"min_bets_to_qualify"`
}
