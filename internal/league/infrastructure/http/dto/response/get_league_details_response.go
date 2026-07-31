package response

import "time"

type ParticipantRankingResponse struct {
	ParticipantID  string  `json:"participant_id"`
	UserID         string  `json:"user_id"`
	Username       string  `json:"username"`
	ProfilePicture *string `json:"profile_picture"`
	Balance        float64 `json:"balance"`
	Position       int     `json:"position"`
}

type GetLeagueDetailsResponse struct {
	LeagueID         string                       `json:"league_id"`
	Name             string                       `json:"name"`
	AdminID          string                       `json:"admin_id"`
	InitialBalance   float64                      `json:"initial_balance"`
	MaxRecharges     int                          `json:"max_recharges"`
	IsRankingVisible bool                         `json:"is_ranking_visible"`
	InviteCode       string                       `json:"invite_code"`
	CreatedAt        time.Time                    `json:"created_at"`
	Participants     []ParticipantRankingResponse `json:"participants"`
}
