package response

type ParticipantRankingResponse struct {
	ParticipantID  string  `json:"participant_id"`
	UserID         string  `json:"user_id"`
	Username       string  `json:"username"`
	ProfilePicture *string `json:"profile_picture"`
	Balance        float64 `json:"balance"`
	Position       int     `json:"position"`
	Role           string  `json:"role"`
}

type GetLeagueDetailsResponse struct {
	LeagueID         string                       `json:"league_id"`
	Slug             string                       `json:"slug"`
	Name             string                       `json:"name"`
	OwnerID          string                       `json:"owner_id"`
	InitialBalance   float64                      `json:"initial_balance"`
	MaxRecharges     int                          `json:"max_recharges"`
	IsRankingVisible bool                         `json:"is_ranking_visible"`
	InviteCode       string                       `json:"invite_code"`
	CreatedAt        string                       `json:"created_at"`
	Participants     []ParticipantRankingResponse `json:"participants"`
}