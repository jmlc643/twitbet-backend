package request

type UpdateMarketOddsRequest struct {
	OptionsOdds map[string]float64 `json:"options_odds" binding:"required"`
}