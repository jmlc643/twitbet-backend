package request

type JoinLeagueRequest struct {
	InviteCode string `json:"invite_code" binding:"required,len=8"`
}