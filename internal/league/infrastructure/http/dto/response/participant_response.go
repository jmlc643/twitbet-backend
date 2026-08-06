package response

import "time"

type ParticipantMeResponse struct {
	ID                string    `json:"id"`
	LeagueID          string    `json:"league_id"`
	UserID            string    `json:"user_id"`
	IsAdmin           bool      `json:"is_admin"`
	Balance           float64   `json:"balance"`
	RechargesConsumed int       `json:"recharges_consumed"`
	JoinedAt          time.Time `json:"joined_at"`
}
