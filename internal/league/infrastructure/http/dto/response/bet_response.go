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

type BetDetailResponse struct {
	ID           string    `json:"id"`
	Amount       float64   `json:"amount"`
	Odds         float64   `json:"odds"`
	PotentialWin float64   `json:"potential_win"`
	Status       string    `json:"status"`
	PlacedAt     time.Time `json:"placed_at"`
	MatchTitle   string    `json:"match_title"`
	MarketName   string    `json:"market_name"`
	OptionName   string    `json:"option_name"`
}

type PaginatedBetResponse struct {
	Data []BetDetailResponse `json:"data"`
	Meta PaginationMeta      `json:"meta"`
}

type PaginationMeta struct {
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}
