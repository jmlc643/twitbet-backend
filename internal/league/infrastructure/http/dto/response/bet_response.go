package response

import "time"

type BetResponse struct {
	ID             string    `json:"id"`
	ParticipantID  string    `json:"participant_id"`
	MarketOptionID string    `json:"market_option_id"`
	Amount         float64   `json:"amount"`
	Odds           float64   `json:"odds"`
	PotentialWin   float64   `json:"potential_win"`
	Status         string    `json:"status"`
	PlacedAt       time.Time `json:"placed_at"`
}
