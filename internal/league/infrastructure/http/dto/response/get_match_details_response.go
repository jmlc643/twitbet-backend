package response

import "time"

type GetMatchDetailsResponse struct {
	ID        string           `json:"id"`
	LeagueID  string           `json:"league_id"`
	Slug      string           `json:"slug"`
	Title     string           `json:"title"`
	StartTime time.Time        `json:"start_time"`
	Status    string           `json:"status"`
	Markets   []MarketResponse `json:"markets"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}