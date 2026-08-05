package response

type CreateLeagueResponse struct {
	ID         string `json:"id"`
	Slug       string `json:"slug"`
	InviteCode string `json:"invite_code"`
}
