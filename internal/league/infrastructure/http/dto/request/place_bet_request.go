package request

type PlaceBetRequest struct {
	LeagueID       string  `json:"league_id" binding:"required,uuid"`
	MarketID       string  `json:"market_id" binding:"required,uuid"`
	MarketOptionID string  `json:"market_option_id" binding:"required,uuid"`
	Amount         float64 `json:"amount" binding:"required,gt=0"`
	AcceptedOdds   float64 `json:"accepted_odds" binding:"required,gt=1"`
	BonusID        *string `json:"bonus_id,omitempty" binding:"omitempty,uuid"`
}
